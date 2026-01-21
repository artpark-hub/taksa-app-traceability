package data

import (
	"traceability/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger" // <--- ADD THIS IMPORT
)

var ProviderSet = wire.NewSet(NewData, NewGreeterRepo, NewTraceabilityRepo)

type Data struct {
	db *gorm.DB
}

func NewData(c *conf.Data, loggerInput log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(loggerInput).Info("closing the data resources")
	}

	dsn := c.Database.Source

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, nil, err
	}

	return &Data{db: db}, cleanup, nil
}
