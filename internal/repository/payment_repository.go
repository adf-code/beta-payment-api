package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/adf-code/beta-payment-api/internal/delivery/request"
	"github.com/adf-code/beta-payment-api/internal/entity"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type paymentRepo struct {
	DB *sql.DB
}

type PaymentRepository interface {
	FetchWithQueryParams(ctx context.Context, params request.PaymentListQueryParams) ([]entity.Payment, error)
	FetchByID(ctx context.Context, id uuid.UUID) (*entity.Payment, error)
	ModifyByID(ctx context.Context, id uuid.UUID, req *request.UpdatePaymentRequest) (*entity.Payment, error)
	Store(ctx context.Context, tx *sql.Tx, payment *entity.Payment) error
	Remove(ctx context.Context, id uuid.UUID) error
}

func NewPaymentRepo(db *sql.DB) PaymentRepository {
	return &paymentRepo{DB: db}
}

func (r *paymentRepo) FetchWithQueryParams(ctx context.Context, params request.PaymentListQueryParams) ([]entity.Payment, error) {
	query := "SELECT id, tag, description, amount, status, created_at, updated_at FROM payments WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	// Search
	if params.SearchField != "" && params.SearchValue != "" {
		query += fmt.Sprintf(" AND %s ILIKE $%d", params.SearchField, argIndex)
		args = append(args, "%"+params.SearchValue+"%")
		argIndex++
	}

	// Filters
	for _, f := range params.Filter {
		if len(f.Value) > 0 {
			query += fmt.Sprintf(" AND %s = ANY($%d)", f.Field, argIndex)
			args = append(args, pq.Array(f.Value))
			argIndex++
		}
	}

	// Range
	for _, r := range params.Range {
		if r.From != nil {
			query += fmt.Sprintf(" AND %s >= $%d", r.Field, argIndex)
			args = append(args, *r.From)
			argIndex++
		}
		if r.To != nil {
			query += fmt.Sprintf(" AND %s <= $%d", r.Field, argIndex)
			args = append(args, *r.To)
			argIndex++
		}
	}

	// Sort
	if params.SortField != "" && (params.SortDir == "ASC" || params.SortDir == "DESC") {
		query += fmt.Sprintf(" ORDER BY %s %s", params.SortField, params.SortDir)
	}

	// Pagination
	if params.Page > 0 && params.PerPage > 0 {
		offset := (params.Page - 1) * params.PerPage
		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
		args = append(args, params.PerPage, offset)
	}
	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []entity.Payment
	for rows.Next() {
		var p entity.Payment
		if err := rows.Scan(&p.ID, &p.Tag, &p.Description, &p.Amount, &p.Status, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		payments = append(payments, p)
	}

	return payments, nil
}

func (r *paymentRepo) FetchByID(ctx context.Context, id uuid.UUID) (*entity.Payment, error) {
	var p entity.Payment
	err := r.DB.QueryRowContext(ctx, "SELECT id, tag, description, amount, status, created_at, updated_at FROM payments WHERE id = $1 AND deleted_at is null", id).
		Scan(&p.ID, &p.Tag, &p.Description, &p.Amount, &p.Status, &p.CreatedAt, &p.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *paymentRepo) ModifyByID(ctx context.Context, id uuid.UUID, req *request.UpdatePaymentRequest) (*entity.Payment, error) {
	query := `
		UPDATE payments
		SET status = $1, updated_at = NOW()
		WHERE id = $2
		RETURNING id, tag, description, amount, status, created_at, updated_at
	`

	row := r.DB.QueryRowContext(ctx, query, req.Status, id)

	var updated entity.Payment

	err := row.Scan(
		&updated.ID,
		&updated.Tag,
		&updated.Description,
		&updated.Amount,
		&updated.Status,
		&updated.CreatedAt,
		&updated.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &updated, nil
}

func (r *paymentRepo) Store(ctx context.Context, tx *sql.Tx, payment *entity.Payment) error {
	return tx.QueryRowContext(
		ctx,
		"INSERT INTO payments (tag, description, amount) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at, status",
		payment.Tag, payment.Description, payment.Amount,
	).Scan(&payment.ID, &payment.CreatedAt, &payment.UpdatedAt, &payment.Status)
}

func (r *paymentRepo) Remove(ctx context.Context, id uuid.UUID) error {
	_, err := r.DB.ExecContext(ctx, "DELETE FROM payments WHERE id = $1", id)
	return err
}
