package data

import (
	"traceability/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/nats-io/nats.go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var ProviderSet = wire.NewSet(NewData, NewGreeterRepo, NewTraceabilityRepo)

type Data struct {
	db   *gorm.DB
	nats *nats.Conn
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

	nc, err := nats.Connect(
		c.Nats.Url,
		nats.UserInfo(c.Nats.Username, c.Nats.Password),
	)
	if err != nil {
		l.Errorf("failed to connect to nats: %v", err)
		return nil, nil, err
	}
	l.Infof("connected to nats at %s", c.Nats.Url)

	d := &Data{
		db:   db,
		nats: nc,
	}

	d.StartConsumer(loggerInput)

	cleanup := func() {
		l.Info("closing the data resources")
		nc.Close()
	}

	return d, cleanup, nil
}
