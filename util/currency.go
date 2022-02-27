package util

func ValidCurrency(currency string) bool {
	for _, cur := range AllowedCurrencies() {
		if cur == currency {
			return true
		}
	}
	return false
}

func AllowedCurrencies() [4]string {
	return [4]string{
		"USD",
		"SAR",
		"LE",
		"EUR",
	}
}
