package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yscheng-gf/gioco-db-migration/internal/config"
	"github.com/yscheng-gf/gioco-db-migration/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	configFile  = "etc/config.yaml"
	devEnvFile  = "etc/.env"
	prodEnvFile = "etc/prod.env"
)

var (
	env *string
	c   = &config.Conf{}

	rootCmd = cobra.Command{
		Use: "gioco-db-migrate",
		Long: `
   _______                         ____                 _                  __     
  / ____(_)___  _________     ____/ / /_     ____ ___  (_)___ __________ _/ /____ 
 / / __/ / __ \/ ___/ __ \   / __  / __ \   / __  __ \/ / __  / ___/ __  / __/ _ \
/ /_/ / / /_/ / /__/ /_/ /  / /_/ / /_/ /  / / / / / / / /_/ / /  / /_/ / /_/  __/
\____/_/\____/\___/\____/   \__,_/_.___/  /_/ /_/ /_/_/\__, /_/   \__,_/\__/\___/ 
                                                      /____/                      
`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	migrateCmd = cobra.Command{
		Use:   "migrate",
		Short: "This is migrate postgresDB column balance digital",
		Long:  "This is migrate postgresDB column balance digital",
		Run: func(cmd *cobra.Command, args []string) {
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

			fmt.Println("Migrate Environment: " + cmd.Flags().Lookup("environment").Value.String())
			for _, op := range *ops {
				pgDb.Exec(fmt.Sprintf("ALTER TABLE %s ALTER COLUMN balance TYPE numeric(24, 8);", op.Code+"_member_wallets"))
				pgDb.Exec(fmt.Sprintf("ALTER TABLE %s ALTER COLUMN before_balance TYPE numeric(24, 8);", op.Code+"_member_transactions"))
				pgDb.Exec(fmt.Sprintf("ALTER TABLE %s ALTER COLUMN amount TYPE numeric(24, 8);", op.Code+"_member_transactions"))
				pgDb.Exec(fmt.Sprintf("ALTER TABLE %s ALTER COLUMN after_balance TYPE numeric(24, 8);", op.Code+"_member_transactions"))

				fmt.Printf("%s done.\n", op.Code)
			}
		},
	}
)

func main() {
	env = migrateCmd.Flags().StringP("environment", "e", "dev", "This flag for setting db environment. Allow: [\"dev\", \"prod\"]")
	rootCmd.AddCommand(&migrateCmd)
	rootCmd.Execute()
}

func mustLoadConfig(c *config.Conf) {
	envFile := devEnvFile
	if *env == "prod" {
		envFile = prodEnvFile
	}
	if err := godotenv.Load(envFile); err != nil {
		log.Fatal("Error loading .env file")
	}

	viper.AddConfigPath(".")
	viper.SetConfigName(configFile)
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