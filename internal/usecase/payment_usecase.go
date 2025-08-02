package usecase

import (
	"context"
	"database/sql"
	"github.com/adf-code/beta-payment-api/internal/delivery/request"
	"github.com/adf-code/beta-payment-api/internal/entity"
	"github.com/adf-code/beta-payment-api/internal/repository"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type PaymentUseCase interface {
	GetAll(ctx context.Context, params request.PaymentListQueryParams) ([]entity.Payment, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Payment, error)
	UpdateByID(ctx context.Context, id uuid.UUID, req *request.UpdatePaymentRequest) (*entity.Payment, error)
	Create(ctx context.Context, payment entity.Payment) (*entity.Payment, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type paymentUseCase struct {
	paymentRepo repository.PaymentRepository
	db          *sql.DB
	logger      zerolog.Logger
}

func NewPaymentUseCase(paymentRepo repository.PaymentRepository, db *sql.DB, logger zerolog.Logger) PaymentUseCase {
	return &paymentUseCase{
		paymentRepo: paymentRepo,
		db:          db,
		logger:      logger,
	}
}

func (uc *paymentUseCase) GetAll(ctx context.Context, params request.PaymentListQueryParams) ([]entity.Payment, error) {
	uc.logger.Info().Str("usecase", "GetAll").Msg("⚙️ Fetching all payments")
	return uc.paymentRepo.FetchWithQueryParams(ctx, params)
}

func (uc *paymentUseCase) GetByID(ctx context.Context, id uuid.UUID) (*entity.Payment, error) {
	uc.logger.Info().Str("usecase", "GetByID").Msg("⚙️ Fetching payment by ID")
	return uc.paymentRepo.FetchByID(ctx, id)
}

func (uc *paymentUseCase) UpdateByID(ctx context.Context, id uuid.UUID, req *request.UpdatePaymentRequest) (*entity.Payment, error) {
	uc.logger.Info().Str("usecase", "UpdateByID").Msg("⚙️ Fetching all payments")
	return uc.paymentRepo.ModifyByID(ctx, id, req)
}

func (uc *paymentUseCase) Create(ctx context.Context, payment entity.Payment) (*entity.Payment, error) {
	uc.logger.Info().Str("usecase", "Create").Msg("⚙️ Store payment")
	tx, err := uc.db.Begin()
	if err != nil {
		uc.logger.Error().Err(err).Msg("❌ Failed to begin transaction")
		return nil, err
	}

	err = uc.paymentRepo.Store(ctx, tx, &payment)
	if err != nil {
		tx.Rollback()
		uc.logger.Error().Err(err).Msg("❌ Failed to store payment, rolling back")
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		uc.logger.Error().Err(err).Msg("❌ Failed to commit transaction")
		return nil, err
	}

	uc.logger.Info().Str("payment_id", payment.ID.String()).Msg("✅ Payment created")
	return &payment, nil
}

func (uc *paymentUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	uc.logger.Info().Str("usecase", "Delete").Msg("⚙️ Remove payment")
	return uc.paymentRepo.Remove(ctx, id)
}
