# gioco-db-migration

## First
Before using, you need to create two files named `etc/.env` and `etc/.prod.env`.
Then fill two files.

Run following command:
```sh
make env
```

## Docker
Build
```sh
make docker
```

## Usage
As use docker, run following command:

```sh
docker run --rm -it gioco-db-migration
```

Or you have `golang` environment:
```sh
go run main.go
```

## Summary
Migrate postgres db numeric(24, 2) -> numeric(24, 8)

Affected fields
- op_member_transactions
  - before_balance
  - amount
  - after_balance
- op_member_wallets
  - balance
