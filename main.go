package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
	"github.com/yscheng-gf/gioco-db-migration/internal/config"
	"github.com/yscheng-gf/gioco-db-migration/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	configFile  = "etc/config.yaml"
	devEnvFile  = "etc/.env"
	prodEnvFile = "etc/.prod.env"
)

var (
	env     *string
	digital *int
	c       = &config.Conf{}

	rootCmd = cobra.Command{
		Use: "gioco-db-migration",
		Long: `
   _______                         ____                 _                  __  _           
  / ____(_)___  _________     ____/ / /_     ____ ___  (_)___ __________ _/ /_(_)___  ____ 
 / / __/ / __ \/ ___/ __ \   / __  / __ \   / __  __ \/ / __  / ___/ __  / __/ / __ \/ __ \
/ /_/ / / /_/ / /__/ /_/ /  / /_/ / /_/ /  / / / / / / / /_/ / /  / /_/ / /_/ / /_/ / / / /
\____/_/\____/\___/\____/   \__,_/_.___/  /_/ /_/ /_/_/\__, /_/   \__,_/\__/_/\____/_/ /_/ 
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

			pgDb, err := gorm.Open(
				postgres.Open(fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", c.Postgres.Host, c.Postgres.Port, c.Postgres.User, c.Postgres.Password, c.Postgres.DBName, c.Postgres.SslMode)),
				&gorm.Config{
					Logger: logger.Default.LogMode(logger.Silent),
				},
			)
			if err != nil {
				panic(fmt.Errorf("gorm error: %w", err))
			}

			p := mpb.New(mpb.WithWidth(64))
			bar := p.AddBar(int64(len(*ops)),
				mpb.PrependDecorators(
					decor.Name("Migrate", decor.WCSyncSpaceR),
					decor.Spinner([]string{}, decor.WCSyncSpaceR),
					decor.Percentage(decor.WCSyncSpaceR),
				),
				mpb.AppendDecorators(
					decor.OnComplete(
						decor.EwmaETA(decor.ET_STYLE_GO, 30, decor.WCSyncWidth), "done",
					),
				),
			)

			fmt.Println("Migrate Environment: " + *env)
			for _, op := range *ops {
				start := time.Now()
				pgDb.Exec(fmt.Sprintf("ALTER TABLE %s ALTER COLUMN balance TYPE numeric(24, %d);", op.Code+"_member_wallets", *digital))
				pgDb.Exec(fmt.Sprintf("ALTER TABLE %s ALTER COLUMN before_balance TYPE numeric(24, %d);", op.Code+"_member_transactions", *digital))
				pgDb.Exec(fmt.Sprintf("ALTER TABLE %s ALTER COLUMN amount TYPE numeric(24, %d);", op.Code+"_member_transactions", *digital))
				pgDb.Exec(fmt.Sprintf("ALTER TABLE %s ALTER COLUMN after_balance TYPE numeric(24, %d);", op.Code+"_member_transactions", *digital))

				// fmt.Printf("%s done.\n", op.Code)
				bar.EwmaIncrement(time.Since(start))
			}
			p.Wait()
		},
	}
)

func main() {
	env = migrateCmd.Flags().StringP("environment", "e", "dev", "This flag for setting db environment. Allow: [\"dev\", \"prod\"]")
	digital = migrateCmd.Flags().IntP("digital", "d", 8, "This flag for setting Nnd decimal place.")
	rootCmd.AddCommand(&migrateCmd)
	rootCmd.Execute()
}

func mustLoadConfig(c *config.Conf) {
	var envFile string
	switch *env {
	case "prod":
		envFile = prodEnvFile
	case "dev":
		envFile = devEnvFile
	default:
		log.Fatalln("Error loading .env file")
	}

	if err := godotenv.Load(envFile); err != nil {
		log.Fatalf("Error loading %s file", envFile)
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
