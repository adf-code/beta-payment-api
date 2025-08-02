package mocks

import (
	"context"
	"database/sql"
	"github.com/adf-code/beta-payment-api/internal/delivery/request"
	"github.com/adf-code/beta-payment-api/internal/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type PaymentRepository struct {
	mock.Mock
}

func (m *PaymentRepository) FetchWithQueryParams(ctx context.Context, params request.PaymentListQueryParams) ([]entity.Payment, error) {
	args := m.Called(ctx, params)
	return args.Get(0).([]entity.Payment), args.Error(1)
}

func (m *PaymentRepository) FetchByID(ctx context.Context, id uuid.UUID) (*entity.Payment, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entity.Payment), args.Error(1)
}

func (m *PaymentRepository) Store(ctx context.Context, tx *sql.Tx, payment *entity.Payment) error {
	args := m.Called(ctx, tx, payment)
	return args.Error(0)
}

func (m *PaymentRepository) Remove(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
