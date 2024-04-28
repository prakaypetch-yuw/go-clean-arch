package db

import (
	"context"
	"fmt"

	"github.com/prakaypetch-yuw/go-clean-arch/config"
	"github.com/prakaypetch-yuw/go-clean-arch/internal/domain/entity"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type GormContextType string

var gormContextKey GormContextType = "gorm"

func ProvideDB(cfg config.Config) (*gorm.DB, func(), error) {
	db, err := gorm.Open(mysql.Open(dsn(cfg)), &gorm.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("Database initialize failed")
		return nil, func() {
			// Do nothing
		}, err
	}

	conn, err := db.DB()
	if err != nil {
		return nil, func() {
			// Do nothing
		}, err
	}

	cleanup := func() {
		err := conn.Close()
		if err != nil {
			return
		}
	}

	if err := conn.Ping(); err != nil {
		return nil, cleanup, err
	}

	err = db.AutoMigrate(&entity.User{})
	if err != nil {
		return nil, cleanup, err
	}

	log.Info().Msg("Successfully connected to the database")
	return db, cleanup, nil
}

func dsn(cfg config.Config) string {
	dbCfg := cfg.DB
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&tls=false", dbCfg.User, dbCfg.Password,
		dbCfg.Host, dbCfg.Port, dbCfg.Name)
}

func GormFromContext(ctx context.Context) *gorm.DB {
	out, ok := ctx.Value(gormContextKey).(*gorm.DB)
	if !ok {
		return nil
	}
	return out
}
