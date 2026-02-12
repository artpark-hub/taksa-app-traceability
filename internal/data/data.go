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
		// DB connection usually doesn't need explicit close in GORM for long-running apps,
		// but if you had Redis/Nats, you'd close them here.
	}

	return d, cleanup, nil
}
