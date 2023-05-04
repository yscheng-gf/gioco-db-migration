# gioco-db-migration

## Usage
使用前須先新增 `etc/.env` 和 `etc/prod.env` 這兩個檔案

## Docker
build
```sh
docker build --rm -t gioco-db-migration .
```

執行
```sh
docker run --rm gioco-db-migration
```

## Summary
遷移postgres db numeric(24, 2) -> numeric(24, 8)
受影響欄位
- op_member_transactions
  - before_balance
  - amount
  - after_balance
- op_member_wallets
  - balance