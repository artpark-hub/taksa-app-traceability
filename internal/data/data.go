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

    dsn := "host=localhost user=admin password=OnionRings$02 dbname=kavya-timescaledb port=5435 sslmode=disable TimeZone=UTC"

    // UPDATE THIS BLOCK TO ENABLE LOGGING
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info), // <--- ENABLES SQL LOGS
    })
    if err != nil {
        return nil, nil, err
    }

    return &Data{db: db}, cleanup, nil
}
