# gioco-db-migration

## First
 Before using, you need to create two files named `etc/.env` and `etc/prod.env`.

## Docker
build
```sh
make docker
```

## Usage
run flowing command:
```sh
docker run --rm -it gioco-db-migration
```

## Summary
遷移postgres db numeric(24, 2) -> numeric(24, 8)

Affected fields
- op_member_transactions
  - before_balance
  - amount
  - after_balance
- op_member_wallets
  - balance
