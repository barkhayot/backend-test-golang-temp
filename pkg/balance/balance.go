package balance

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrInsufficientBalance = errors.New("insufficient balance")

type Service struct {
	db *pgxpool.Pool
}

func NewService(db *pgxpool.Pool) (*Service, error) {
	if db == nil {
		return nil, errors.New("database connection is nil")
	}

	return &Service{db: db}, nil
}

func (r *Service) Charge(
	ctx context.Context,
	userID int64,
	amount float64,
) (*BalanceHistory, error) {

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var from float64
	err = tx.QueryRow(ctx,
		`SELECT balance FROM users WHERE id = $1 FOR UPDATE`,
		userID,
	).Scan(&from)
	if err != nil {
		return nil, err
	}

	if from < amount {
		return nil, ErrInsufficientBalance
	}

	to := from - amount

	_, err = tx.Exec(ctx,
		`UPDATE users SET balance = $1 WHERE id = $2`,
		to, userID,
	)
	if err != nil {
		return nil, err
	}

	var h BalanceHistory
	err = tx.QueryRow(ctx,
		`INSERT INTO 
			user_balance_history(
				user_id, 
				balance_from, 
				balance_to, 
				amount
			)
		 VALUES ($1, $2, $3, $4)
		 RETURNING 
		 	id, 
			created_at
		`,
		userID, from, to, amount,
	).Scan(&h.ID, &h.CreatedAt)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	h.UserID = userID
	h.BalanceFrom = from
	h.BalanceTo = to
	h.Amount = amount

	return &h, nil
}

func (r *Service) GetHistory(
	ctx context.Context,
	userID int64,
) ([]BalanceHistory, error) {

	rows, err := r.db.Query(ctx,
		`SELECT 
			id, 
			user_id, 
			balance_from, 
			balance_to, 
			amount, 
			created_at
		 FROM 
		 	user_balance_history
		 WHERE 
		 	user_id = $1
		 ORDER BY 
		 	created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []BalanceHistory
	for rows.Next() {
		var h BalanceHistory
		if err := rows.Scan(
			&h.ID,
			&h.UserID,
			&h.BalanceFrom,
			&h.BalanceTo,
			&h.Amount,
			&h.CreatedAt,
		); err != nil {
			return nil, err
		}
		res = append(res, h)
	}

	return res, nil
}
