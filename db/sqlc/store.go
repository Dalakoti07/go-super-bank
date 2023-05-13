package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provide all functions to execute db queries
type Store struct {
	*Queries
	db *sql.DB
}

// NewStore creates a new store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// executes fn in database transactions and
// this function is not exported
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

// TransferTxParams contains input params of Transfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountId   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult contains result of transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

var txKey = struct{}{}

// TransferTx performs a money transfer from one account to another which includes following step
// create a transfer record, add account entries, update accounts balance within a single db transactions
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(queries *Queries) error {
		var err error

		txName := ctx.Value(txKey)

		fmt.Println(txName, "create transfer")
		result.Transfer, err = queries.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountId,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Println(txName, "create entry 1")
		result.FromEntry, err = queries.CreateEntries(ctx, CreateEntriesParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Println(txName, "create entry 2")
		result.ToEntry, err = queries.CreateEntries(ctx, CreateEntriesParams{
			AccountID: arg.ToAccountId,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// update balance
		if arg.FromAccountID < arg.ToAccountId {
			result.FromAccount, result.ToAccount, err = addMoney(
				ctx,
				queries,
				arg.FromAccountID,
				-arg.Amount,
				arg.ToAccountId,
				arg.Amount,
			)
		} else {
			// get account2 -> update its balance
			result.ToAccount, result.FromAccount, err = addMoney(
				ctx,
				queries,
				arg.ToAccountId,
				arg.Amount,
				arg.FromAccountID,
				-arg.Amount,
			)
		}
		return nil
	})
	return result, err
}

// addMoney adds amount1 to accountId1 and amount2 to amountId2
// and returns updated account1 and account2 and error
func addMoney(
	ctx context.Context,
	q *Queries,
	accountId1 int64,
	amount1 int64,
	accountId2 int64,
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(
		ctx,
		AddAccountBalanceParams{
			ID:     accountId1,
			Amount: amount1,
		})
	if err != nil {
		// a cool feature of golang if return types are named they are automatically returned
		return
	}
	account2, err = q.AddAccountBalance(
		ctx, AddAccountBalanceParams{
			ID:     accountId2,
			Amount: amount2,
		})
	return
}
