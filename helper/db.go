package helper

import (
	"context"
	"fmt"
	"strings"

	"github.com/yscheng-gf/gioco-db-migration/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitMongo(c *config.Conf) *mongo.Client {
	mongoCli, err := mongo.Connect(
		context.TODO(),
		options.Client().SetHosts(strings.Split(c.Mongo.Host, ",")),
	)
	if err != nil {
		panic(err)
	}
	return mongoCli
}

func InitPgDb(c *config.Conf) *gorm.DB {
	pgDb, err := gorm.Open(
		postgres.Open(fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", c.Postgres.Host, c.Postgres.Port, c.Postgres.User, c.Postgres.Password, c.Postgres.DBName, c.Postgres.SslMode)),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		},
	)
	if err != nil {
		panic(fmt.Errorf("gorm error: %w", err))
	}
	return pgDb
}
