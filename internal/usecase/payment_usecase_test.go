package usecase_test

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/adf-code/beta-payment-api/internal/delivery/request"
	"github.com/adf-code/beta-payment-api/internal/entity"
	mailMocks "github.com/adf-code/beta-payment-api/internal/pkg/mail/mocks"
	repoMocks "github.com/adf-code/beta-payment-api/internal/repository/mocks"
	"github.com/adf-code/beta-payment-api/internal/usecase"
	"github.com/rs/zerolog"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAllPayments(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mockRepo := new(repoMocks.PaymentRepository)
	mockEmail := new(mailMocks.SendGridClient)
	logger := zerolog.Nop()

	paymentUC := usecase.NewPaymentUseCase(mockRepo, db, logger, mockEmail)

	expected := []entity.Payment{
		{ID: uuid.New(), Title: "Go Programming", Author: "Alice", Year: 2020},
	}

	mockRepo.On("FetchWithQueryParams", mock.Anything, mock.Anything).Return(expected, nil)
	result, err := paymentUC.GetAll(context.TODO(), request.PaymentListQueryParams{})

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestGetPaymentByID(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mockRepo := new(repoMocks.PaymentRepository)
	mockEmail := new(mailMocks.SendGridClient)
	logger := zerolog.Nop()

	paymentUC := usecase.NewPaymentUseCase(mockRepo, db, logger, mockEmail)

	id := uuid.New()
	expectedPayment := &entity.Payment{ID: id, Title: "Clean Code", Author: "Robert C. Martin", Year: 2008}

	mockRepo.On("FetchByID", mock.Anything, id).Return(expectedPayment, nil)

	result, err := paymentUC.GetByID(context.TODO(), id)

	assert.NoError(t, err)
	assert.Equal(t, expectedPayment, result)
	mockRepo.AssertExpectations(t)
}

func TestCreatePayment(t *testing.T) {
	// Step 1: Setup mock DB & transaction
	db, sqlMock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlMock.ExpectBegin()
	sqlMock.ExpectCommit()

	// Step 2: Mock repository & email
	mockRepo := new(repoMocks.PaymentRepository)
	mockEmail := new(mailMocks.SendGridClient)
	logger := zerolog.Nop()

	payment := entity.Payment{
		ID:     uuid.New(),
		Title:  "Test Payment",
		Author: "Test Author",
		Year:   2024,
	}

	// Setup repo expectations (ignore actual DB op)
	mockRepo.On("Store", mock.Anything, mock.AnythingOfType("*sql.Tx"), &payment).Return(nil)
	mockEmail.On("SendPaymentCreatedEmail", payment).Return(nil)

	// Step 3: Call usecase
	paymentUC := usecase.NewPaymentUseCase(mockRepo, db, logger, mockEmail)
	result, err := paymentUC.Create(context.TODO(), payment)

	// Step 4: Assertions
	assert.NoError(t, err)
	assert.Equal(t, payment.Title, result.Title)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
	mockRepo.AssertExpectations(t)
	mockEmail.AssertExpectations(t)
}

func TestDeletePayment(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mockRepo := new(repoMocks.PaymentRepository)
	mockEmail := new(mailMocks.SendGridClient)
	logger := zerolog.Nop()

	paymentUC := usecase.NewPaymentUseCase(mockRepo, db, logger, mockEmail)

	id := uuid.New()
	mockRepo.On("Remove", mock.Anything, id).Return(nil)

	err = paymentUC.Delete(context.TODO(), id)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
