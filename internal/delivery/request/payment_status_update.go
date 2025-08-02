package request

import (
	"errors"
)

type UpdatePaymentRequest struct {
	Status string `json:"status"`
}

func (r *UpdatePaymentRequest) Validate() error {
	if r.Status == "" {
		return errors.New("status is required")
	}
	return nil
}
