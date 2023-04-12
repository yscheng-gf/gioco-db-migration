package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"github.com/yscheng-gf/gioco-db-migration/internal/config"
	"github.com/yscheng-gf/gioco-db-migration/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var configFile = flag.String("f", "etc/config.yaml", "the config file")
var envFile = flag.String("env", "etc/.env", "the env file")

func main() {
	c := new(config.Conf)
	mustLoadConfig(c)

	mongoCli, err := mongo.NewClient(
		options.Client().SetHosts(strings.Split(c.Mongo.Host, ",")),
	)
	if err != nil {
		panic(err)
	}

	mongoCli.Connect(context.Background())

	cursor, err := mongoCli.Database("GIOCO_PLUS").Collection("sys_operators").Find(context.TODO(), bson.M{"status": "1"})
	if err != nil {
		fmt.Println(err)
		return
	}

	defer cursor.Close(context.TODO())

	ops := new([]models.Operator)
	err = cursor.All(context.TODO(), ops)
	if err != nil {
		fmt.Println(err)
		return
	}

	pgDb, err := gorm.Open(postgres.Open(fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", c.Postgres.Host, c.Postgres.Port, c.Postgres.User, c.Postgres.Password, c.Postgres.DBName, c.Postgres.SslMode)))
	if err != nil {
		panic(fmt.Errorf("gorm error: %w", err))
	}

	for _, op := range *ops {
		pgDb.Exec(fmt.Sprintf("ALTER TABLE %s ALTER COLUMN balance TYPE numeric(24, 8);", op.Code+"_member_wallets"))
		pgDb.Exec(fmt.Sprintf("ALTER TABLE %s ALTER COLUMN before_balance TYPE numeric(24, 8);", op.Code+"_member_transactions"))
		pgDb.Exec(fmt.Sprintf("ALTER TABLE %s ALTER COLUMN amount TYPE numeric(24, 8);", op.Code+"_member_transactions"))
		pgDb.Exec(fmt.Sprintf("ALTER TABLE %s ALTER COLUMN after_balance TYPE numeric(24, 8);", op.Code+"_member_transactions"))

		fmt.Printf("%s done.\n", op.Code)
	}

}

func mustLoadConfig(c *config.Conf) {
	if err := godotenv.Load(*envFile); err != nil {
		log.Fatal("Error loading .env file")
	}

	viper.AddConfigPath(".")
	viper.SetConfigName(*configFile)
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("error load config file: %w", err))
	}

	for _, k := range viper.AllKeys() {
		value := viper.GetString(k)
		if strings.HasPrefix(value, "${") && strings.HasSuffix(value, "}") {
			viper.Set(k, getEnvOrPanic(strings.TrimSuffix(strings.TrimPrefix(value, "${"), "}")))
		}
	}

	if err := viper.Unmarshal(c); err != nil {
		panic(fmt.Errorf("error viper unmarshal: %w", err))
	}
}

func getEnvOrPanic(env string) string {
	res := os.Getenv(env)
	if len(res) == 0 {
		panic("Mandatory env variable not found:" + env)
	}
	return res
}
