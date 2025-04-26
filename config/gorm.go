package config

import (
	"os"
	"time"

	"github.com/CRobinDev/BCCGembira_Nusastra/internal/entity"
	"github.com/CRobinDev/BCCGembira_Nusastra/pkg/log"
	"github.com/CRobinDev/BCCGembira_Nusastra/pkg/response"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDB() (*gorm.DB, error) {
	gormLogger := newLogger(log.NewLogger())
	db, err := gorm.Open(postgres.Open(os.Getenv("POSTGRES_DSN")), &gorm.Config{
		PrepareStmt: true,
		Logger:      gormLogger,
	})
	if err != nil {
		return nil, &response.ErrConnectDatabase
	}
	return db, nil
}

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&entity.User{},
		&entity.Chat{},
	)
	if err != nil {
		return &response.ErrMigrateDatabase
	}
	return nil
}

func newLogger(log *logrus.Logger) logger.Interface {
	return logger.New(
		logrusWriter{log},
		logger.Config{
			SlowThreshold:             800 * time.Millisecond,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)
}

type logrusWriter struct {
	Log *logrus.Logger
}

func (w logrusWriter) Printf(format string, args ...interface{}) {
	w.Log.Infof(format, args...)
}
