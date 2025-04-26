package dto

import "github.com/google/uuid"

type PaymentRequest struct {
	CustomerName  string    `json:"customer_name,omitempty"`
	CustomerEmail string    `json:"customer_email,omitempty"`
	UserID        uuid.UUID `json:"user_id,omitempty"`
	OrderID       string    `json:"order_id" validate:"required"`
	Amount        int64     `json:"amount" validate:"required"`
	Type          string    `json:"type" validate:"required,oneof=freeze premium"`
}

type PaymentResponse struct {
	SnapURL string `json:"snap_url,omitempty"`
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
}

type PaymentStatusRequest struct {
	OrderID           string    `json:"order_id,omitempty"`
	TransactionStatus string    `json:"transaction_status,omitempty"`
	FraudStatus       string    `json:"fraud_status,omitempty"`
}
