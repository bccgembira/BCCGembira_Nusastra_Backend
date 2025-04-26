package repository

import (
	"context"

	"github.com/CRobinDev/BCCGembira_Nusastra/internal/entity"
	"github.com/CRobinDev/BCCGembira_Nusastra/pkg/response"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IPaymentRepository interface {
	Save(ctx context.Context, payment *entity.Payment) error
	UpdateStatus(ctx context.Context, payment *entity.Payment) (int64, error)
}

type paymentRepository struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func NewPaymentRepository(db *gorm.DB, logger *logrus.Logger) IPaymentRepository {
	return &paymentRepository{
		db:     db,
		logger: logger,
	}
}

func (pr *paymentRepository) Save(ctx context.Context, payment *entity.Payment) error {
	tr := pr.db.WithContext(ctx).Begin()

	if err := tr.Create(&payment).Error; err != nil {
		tr.Rollback()
		pr.logger.Errorf("failed to create user: %v", err)
		return &response.ErrSavePayment
	}

	if err := tr.Commit().Error; err != nil {
		tr.Rollback()
		pr.logger.Errorf("failed to commit user creation: %v", err)
		return &response.ErrSavePayment
	}

	return nil
}

func (pr *paymentRepository) UpdateStatus(ctx context.Context, payment *entity.Payment) (int64, error) {
	tr := pr.db.WithContext(ctx).Begin()

	result := tr.Model(&payment).Where("order_id = ?", payment.OrderID).Updates(map[string]interface{}{
		"status":     payment.Status,
		"updated_at": payment.UpdatedAt,
	})

	if result.Error != nil {
		tr.Rollback()
		pr.logger.Errorf("failed to update payment status: %v", result.Error)
		return 0, &response.ErrUpdateStatus
	}

	if result.Error == gorm.ErrRecordNotFound {
		pr.logger.Warnf("no payment found for order_id: %v", payment.OrderID)
		return result.RowsAffected, &response.ErrUpdateStatus
	}
	
	if result.Error != nil {
		pr.logger.Errorf("failed to update status: %v", result.Error)
		return result.RowsAffected, &response.ErrUpdateStatus
	}

	if err := tr.Commit().Error; err != nil {
		tr.Rollback()
		pr.logger.Errorf("failed to commit payment status update: %v", err)
		return result.RowsAffected, &response.ErrUpdateStatus
	}

	return result.RowsAffected, nil
}
