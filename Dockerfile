# Build stage
FROM golang:1.17.8-alpine3.15 as builder
WORKDIR /app
COPY . .

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
RUN go build -o main main.go

# RUN stage
FROM alpine:3.15
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/app.env .
# COPY --from=builder /app/migrate .
COPY --from=builder go/bin/migrate .
COPY /db/migrate /app/db/migrate
COPY /scripts /scripts
RUN chmod +x /scripts/*
# RUN apk add postgresql-client

EXPOSE 3009
CMD [ "/app/main" ]
ENTRYPOINT [ "/scripts/start.sh" ]
