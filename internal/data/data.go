package data

import (
	"traceability/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var ProviderSet = wire.NewSet(NewData, NewTraceabilityRepo)

type Data struct {
	db *gorm.DB
}

func NewData(c *conf.Data, loggerInput log.Logger) (*Data, func(), error) {
	l := log.NewHelper(loggerInput)

	dsn := c.Database.Source
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, nil, err
	}

	d := &Data{
		db: db,
	}

	cleanup := func() {
		l.Info("closing the data resources")
		sqlDB, err := db.DB()
		if err != nil {
			l.Errorf("failed to get underlying sql.DB: %v", err)
			return
		}
		if err := sqlDB.Close(); err != nil {
			l.Errorf("failed to close database: %v", err)
		}
	}

	return d, cleanup, nil
}
