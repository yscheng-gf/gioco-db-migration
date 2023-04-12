# gioco-db-migration

遷移postgres db numeric(24, 2) -> numeric(24, 8)
受影響欄位
- op_member_transactions
  - before_balance
  - amount
  - after_balance
- op_member_wallets
  - balance
