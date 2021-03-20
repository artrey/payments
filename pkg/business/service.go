package business

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Service struct {
	pool *pgxpool.Pool
}

type Payment struct {
	Id       string
	SenderId int64
	Amount   int64
}

func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

func (s *Service) CreatePayment(ctx context.Context, senderId, amount int64) (string, error) {
	paymentId := uuid.New().String()
	_, err := s.pool.Exec(ctx, `
		INSERT INTO payments(Id, SenderId, Amount) values($1, $2, $3)
	`, paymentId, senderId, amount)
	if err != nil {
		return "", err
	}
	return paymentId, nil
}

func (s *Service) GetUserPayments(ctx context.Context, userId int64) ([]*Payment, error) {
	return s.extractPayments(ctx, `
		SELECT Id, SenderId, Amount FROM payments WHERE SenderId = $1
	`, userId)
}

func (s *Service) GetAllPayments(ctx context.Context) ([]*Payment, error) {
	return s.extractPayments(ctx, `
		SELECT Id, SenderId, Amount FROM payments
	`)
}

func (s *Service) extractPayments(ctx context.Context, sql string, params ...interface{}) ([]*Payment, error) {
	rows, err := s.pool.Query(ctx, sql, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	payments := make([]*Payment, 0)
	for rows.Next() {
		var p Payment
		err = rows.Scan(&p.Id, &p.SenderId, &p.Amount)
		if err != nil {
			return nil, err
		}
		payments = append(payments, &p)
	}

	return payments, nil
}
