package entity

import (
	"github.com/adf-code/beta-payment-api/internal/valueobject"
	"github.com/google/uuid"
	"time"
)

type Payment struct {
	ID          uuid.UUID            `json:"id"`
	Tag         string               `json:"tag"`
	Description string               `json:"description"`
	Amount      valueobject.BigFloat `json:"amount"`
	Status      string               `json:"status"`
	CreatedAt   *time.Time           `json:"created_at"`
	UpdatedAt   *time.Time           `json:"updated_at"`
}
