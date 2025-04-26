package handler

import (
	"github.com/CRobinDev/BCCGembira_Nusastra/internal/dto"
	"github.com/CRobinDev/BCCGembira_Nusastra/internal/service"
	"github.com/CRobinDev/BCCGembira_Nusastra/pkg/jwt"
	"github.com/CRobinDev/BCCGembira_Nusastra/pkg/response"
	"github.com/CRobinDev/BCCGembira_Nusastra/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type IPaymentHandler interface {
	CreatePayment() fiber.Handler
	UpdatePaymentStatus() fiber.Handler
}

type paymentHandler struct {
	ps  service.IPaymentService
	val validator.Validator
}

func NewPaymentHandler(ps service.IPaymentService, val validator.Validator) IPaymentHandler {
	return &paymentHandler{
		ps:  ps,
		val: val,
	}
}

func (ph *paymentHandler) CreatePayment() fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := dto.PaymentRequest{}

		err := c.BodyParser(&req)
		if err != nil {
			return err
		}

		req.Amount = 10000

		if req.Type == "premium" {
			req.Amount = 15000
		}

		userID, err := jwt.GetUser(c)
		if err != nil {
			return err
		}

		req.OrderID = uuid.NewString()
		req.UserID = userID

		valErr := ph.val.Validate(req)
		if valErr != nil {
			return valErr
		}

		resp, err := ph.ps.CreatePayment(c.Context(), req)
		if err != nil {
			return err
		}

		return response.Success(c, "payments", resp)

	}
}

func (ph *paymentHandler) UpdatePaymentStatus() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var snapReq map[string]interface{}
		if err := c.BodyParser(&snapReq); err != nil {
			return err
		}

		valErr := ph.val.Validate(snapReq)
		if valErr != nil {
			return valErr
		}

		var req dto.PaymentStatusRequest

		if orderID, ok := snapReq["order_id"].(string); ok {
			req.OrderID = orderID
		} else {
			return &response.ErrMissingStatus
		}

		if trxStatus, ok := snapReq["transaction_status"].(string); ok {
			req.TransactionStatus = trxStatus
		} else {
			return &response.ErrMissingStatus
		}

		if fraudStatus, ok := snapReq["fraud_status"].(string); ok {
			req.FraudStatus = fraudStatus
		} else {
			return &response.ErrMissingStatus
		}

		resp, err := ph.ps.UpdatePaymentStatus(c.Context(), req)
		if err != nil {
			return err
		}

		return response.Success(c, "payments", resp)
	}
}

