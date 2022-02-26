package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferParams) (TransferResult, error)
	TransferTxPure(ctx context.Context, args TransferParams) (TransferResult, error)
}

type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// extend Store functionality to execute transactions
func (s *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction error: %v, rollback Error: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

type TransferResult struct {
	FromAccount Account
	ToAccount   Account
	Transfer    Transfer
	FromEntry   Entry
	ToEntry     Entry
}

type TransferParams struct {
	FromAccountID int64
	ToAccountID   int64
	Amount        int64
}

// Transfer amount transaction from AccountA to AccountB using
// 1- check if AccountA has enough balance (AccountA.amount >= amount)
// 2- create transfer record to AccountB with amount amount
// 3- create Entry record on AccountA with -amount ie: negative amount
// 4- create Entry record on AccountB with +amount
// 5- subtract the amount from the AccountA (AccountA.balance - amount)
// 6- add the amount to the AccountB (AccountB.balance + amount)
func (store *SQLStore) TransferTx(ctx context.Context, arg TransferParams) (TransferResult, error) {
	var result TransferResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		// 1- check enough balance
		enoughParam := EnoughAccountBalanceParams{
			Balance: arg.Amount,
			ID:      arg.FromAccountID,
		}

		// func (q *Queries) EnoughAccountBalance(ctx context.Context, arg EnoughAccountBalanceParams) (bool, error) {
		enoughBalance, err := q.EnoughAccountBalance(ctx, enoughParam)
		if err != nil {
			return err
		}
		if !enoughBalance {
			return fmt.Errorf("not enough balance")
		}

		// 2- Create transfer record
		tsfrParams := CreateTransferParams(arg)
		result.Transfer, err = q.CreateTransfer(ctx, tsfrParams)
		if err != nil {
			return err
		}

		// 3- Create FromAccount entry record
		fromEntryParam := CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		}
		result.FromEntry, err = q.CreateEntry(ctx, fromEntryParam)
		if err != nil {
			return err
		}

		// 4- Create toAccount entry record
		toEntryParam := CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		}
		result.ToEntry, err = q.CreateEntry(ctx, toEntryParam)
		if err != nil {
			return err
		}

		if arg.FromAccountID > arg.ToAccountID {
			// 5- Subtract amount from fromAccount balance
			result.FromAccount, err = addMoney(q, ctx, arg.FromAccountID, -arg.Amount)
			if err != nil {
				return err
			}
			// 6- Add amount to toAccount balance
			result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{ID: arg.ToAccountID, Amount: arg.Amount})
			if err != nil {
				return err
			}
		} else {
			// 6- Add amount to toAccount balance
			result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{ID: arg.ToAccountID, Amount: arg.Amount})
			if err != nil {
				return err
			}
			// 5- Subtract amount from fromAccount balance
			result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{ID: arg.FromAccountID, Amount: -arg.Amount})
			if err != nil {
				return err
			}
		}

		return err
	})

	return result, err
}

func addMoney(q *Queries, ctx context.Context, accountId int64, amount int64) (Account, error) {
	return q.AddAccountBalance(ctx, AddAccountBalanceParams{ID: accountId, Amount: amount})
}

func (store *SQLStore) TransferTxPure(ctx context.Context, args TransferParams) (TransferResult, error) {

	var result TransferResult
	// Create a helper function for preparing failure results.
	fail := func(err error) (TransferResult, error) {
		return result, fmt.Errorf("TransferTxPure: %v", err)
	}

	// Get a Tx for making transaction requests.
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return fail(err)
	}
	// Defer a rollback in case anything fails.
	defer tx.Rollback()

	// 1- check if AccountA has enough balance (AccountA.amount >= amount)
	var enough bool
	if err = tx.QueryRowContext(ctx, "SELECT (balance >= $1) from accounts where id = $2",
		args.Amount, args.FromAccountID).Scan(&enough); err != nil {
		if err == sql.ErrNoRows {
			return fail(fmt.Errorf("invalid account"))
		}
		return fail(err)
	}
	if !enough {
		return fail(fmt.Errorf("not enough balance"))
	}

	// 2- create transfer record to AccountB with amount amount
	_, err = tx.ExecContext(ctx, "INSERT INTO transfers (from_account_id, to_account_id, amount) VALUES ($1, $2, $3)",
		args.FromAccountID, args.ToAccountID, args.Amount)
	if err != nil {
		return fail(err)
	}

	// 3- create Entry record on AccountA with -amount ie: negative amount
	_, err = tx.ExecContext(ctx, "INSERT INTO entries (account_id, amount) VALUES ($1, $2)",
		args.FromAccountID, -args.Amount)
	if err != nil {
		return fail(err)
	}

	// 4- create Entry record on AccountB with +amount
	_, err = tx.ExecContext(ctx, "INSERT INTO entries (account_id, amount) VALUES ($1, $2)",
		args.ToAccountID, args.Amount)
	if err != nil {
		return fail(err)
	}

	// 5- subtract the amount from the AccountA (AccountA.balance - amount)
	_, err = tx.ExecContext(ctx, "UPDATE accounts SET balance = balance - $1 WHERE id = $2",
		args.Amount, args.FromAccountID)
	if err != nil {
		return fail(err)
	}
	// 6- add the amount to the AccountB (AccountB.balance + amount)
	_, err = tx.ExecContext(ctx, "UPDATE accounts SET balance = balance + $1 WHERE id = $2",
		args.Amount, args.ToAccountID)
	if err != nil {
		return fail(err)
	}

	result.FromAccount, _ = store.Queries.GetAccount(ctx, args.FromAccountID)
	result.ToAccount, _ = store.Queries.GetAccount(ctx, args.ToAccountID)
	// Create a new row in the album_order table.
	if err != nil {
		return fail(err)
	}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		return fail(err)
	}

	// Return the order ID.
	return result, nil
}
