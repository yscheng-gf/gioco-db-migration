package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yscheng-gf/gioco-db-migration/cmd"
	"github.com/yscheng-gf/gioco-db-migration/internal/config"
)

const (
	configFile     = "etc/config.yaml"
	devEnvFile     = "etc/.env"
	stagingEnvFile = "etc/.staging.env"
	prodEnvFile    = "etc/.prod.env"
)

var (
	env string

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
)

func main() {
	rootCmd.PersistentFlags().StringVarP(&env, "environment", "e", "", "This flag for setting db environment. Allow: [\"dev\", \"staging\", \"prod\"]")
	if err := rootCmd.ParseFlags(os.Args); err != nil {
		log.Fatalln(err)
	}
	c := &config.Conf{}
	mustLoadConfig(c)
	rootCmd.AddCommand(cmd.NewTransLogPrecisionMigrateCmd(c))
	rootCmd.AddCommand(cmd.NewTransLogFeeMigrateCmd(c))
	rootCmd.AddCommand(cmd.NewMigrateRecommenderCmd(c))
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}

func mustLoadConfig(c *config.Conf) {
	var envFile string
	switch env {
	case "prod":
		envFile = prodEnvFile
	case "staging":
		envFile = stagingEnvFile
	case "dev":
		envFile = devEnvFile
	default:
		log.Fatalf("Error loading [%s] .env file\n", env)
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
