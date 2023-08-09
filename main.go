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
	configFile  = "etc/config.yaml"
	devEnvFile  = "etc/.env"
	prodEnvFile = "etc/.prod.env"
)

var (
	env *string

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
	env = rootCmd.PersistentFlags().StringP("environment", "e", "dev", "This flag for setting db environment. Allow: [\"dev\", \"prod\"]")
	c := &config.Conf{}
	mustLoadConfig(c)
	rootCmd.AddCommand(cmd.NewTransLogPrecisionMigrateCmd(c))
	rootCmd.AddCommand(cmd.NewTransLogFeeMigrateCmd(c))
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
