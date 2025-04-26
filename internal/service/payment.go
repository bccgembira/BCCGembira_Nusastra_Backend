package service

import (
	"context"

	"github.com/CRobinDev/BCCGembira_Nusastra/internal/dto"
	"github.com/CRobinDev/BCCGembira_Nusastra/internal/entity"
	"github.com/CRobinDev/BCCGembira_Nusastra/internal/repository"
	"github.com/CRobinDev/BCCGembira_Nusastra/pkg/helper"
	"github.com/CRobinDev/BCCGembira_Nusastra/pkg/midtrans"
	"github.com/CRobinDev/BCCGembira_Nusastra/pkg/response"
	"github.com/sirupsen/logrus"
)

type IPaymentService interface {
	CreatePayment(ctx context.Context, req dto.PaymentRequest) (dto.PaymentResponse, error)
	UpdatePaymentStatus(ctx context.Context, req dto.PaymentStatusRequest) (dto.PaymentResponse, error)
}

type paymentService struct {
	pr       repository.IPaymentRepository
	ur       repository.IUserRepository
	logger   *logrus.Logger
	midtrans midtrans.IMidtrans
}

func NewPaymentService(pr repository.IPaymentRepository, ur repository.IUserRepository, logger *logrus.Logger, midtrans midtrans.IMidtrans) IPaymentService {
	return &paymentService{
		pr:       pr,
		ur:       ur,
		logger:   logger,
		midtrans: midtrans,
	}
}

func (ps *paymentService) CreatePayment(ctx context.Context, req dto.PaymentRequest) (dto.PaymentResponse, error) {
	user, err := ps.ur.FindByID(ctx, req.UserID)
	if err != nil {
		ps.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("[paymentService.CreatePayment] Failed to find user")
		return dto.PaymentResponse{}, &response.ErrFailedFindUser
	}

	req.CustomerName = user.DisplayName
	req.CustomerEmail = user.Email

	snapResp, err := ps.midtrans.NewTransactionToken(req)
	if err != nil {
		ps.logger.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("[paymentService.CreatePayment] Failed to create payment token")
		return dto.PaymentResponse{}, err
	}

	time := helper.GetCurrentTime()
	payment := entity.Payment{
		OrderID:   req.OrderID,
		UserID:    req.UserID,
		Amount:    float64(req.Amount),
		Type:      req.Type,
		CreatedAt: time,
		UpdatedAt: time,
	}

	err = ps.pr.Save(ctx, &payment)
	if err != nil {
		return dto.PaymentResponse{}, err
	}

	ps.logger.WithFields(map[string]interface{}{
		"order_id":  payment.OrderID,
		"cust_name": user.DisplayName,
	}).Info("[paymentService.CreatePayment] Payment created successfully")

	return dto.PaymentResponse{
		SnapURL: snapResp.RedirectURL,
		OrderID: payment.OrderID,
		Status:  payment.Status,
	}, nil
}

func (ps *paymentService) UpdatePaymentStatus(ctx context.Context, req dto.PaymentStatusRequest) (dto.PaymentResponse, error) {
	var status string

	switch req.TransactionStatus {
	case "capture":
		switch req.FraudStatus {
		case "challenge":
			status = "challenge"
		case "accept":
			status = "success"
		default:
			status = "unknown"
		}

	case "settlement":
		status = "success"

	case "cancel", "expire":
		status = "failure"

	case "pending":
		status = "pending"

	case "deny":
		status = "denied"

	default:
		status = "unknown"
	}

	payment := entity.Payment{
		Status:  status,
		OrderID: req.OrderID,
	}

	rowsAffected, err := ps.pr.UpdateStatus(ctx, &payment)
	if err != nil && rowsAffected == 0 {
		ps.logger.WithFields(map[string]interface{}{
			"order_id": req.OrderID,
			"error":    err.Error(),
		}).Error("[paymentService.UpdatePaymentStatus] failed to update payment status")
		return dto.PaymentResponse{}, err
	}

	resp := dto.PaymentResponse{
		OrderID: req.OrderID,
		Status:  status,
	}

	ps.logger.WithFields(map[string]interface{}{
		"order_id": req.OrderID,
		"status":   status,
	}).Info("[paymentService.UpdatePaymentStatus] payment status updated successfully")
	return resp, nil
}
