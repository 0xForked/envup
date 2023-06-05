package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"sync"
)

func init() {
	loadEnv()
}

func main() {
	var (
		envKey   string
		envValue string
	)
	var rootCMD = &cobra.Command{Long: `Env Updater CLI`}
	// usage go run main.go update-env --key "UPDATE_ME_PLEASE" --value "MY_NEW_VALUE"
	var updaterEnvCMD = &cobra.Command{
		Use:   "update-env",
		Short: `updater environment: update-env --key {env_key} --value {env_value}"`,
		Run: func(cmd *cobra.Command, args []string) {
			if envKey == "" && envValue == "" {
				log.Println(`please set --key and --value, e.g: go run main.go --key "MY_KEY" --value "MY_VALUE" `)
				return
			}
			Instance.updateEnvHelper(envKey, envValue)
		},
	}
	updaterEnvCMD.Flags().StringVar(&envKey, "key", "", "environment key")
	updaterEnvCMD.Flags().StringVar(&envValue, "value", "", "environment value")
	rootCMD.AddCommand(updaterEnvCMD)
	if err := rootCMD.Execute(); err != nil {
		log.Panicf("ERR_EXEC %s", err)
	}
}

/// CONFIGURATION

var (
	cfgOnce  sync.Once
	Instance *AppConfig
)

type AppConfig struct {
	UpdateMePlease string `mapstructure:"UPDATE_ME_PLEASE"`
}

func loadEnv() {
	log.Printf("Load configuration file . . . .")
	// read env handler
	readEnv := func() {
		// load env file
		viper.SetConfigFile(".env")
		// find environment file
		viper.AutomaticEnv()
		// error handling for specific case
		if err := viper.ReadInConfig(); err != nil {
			// specified error message
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				// Config file not found; ignore error if desired
				log.Println(".env file not found!, please copy .example.env and paste as .env")
			}
			// general error message
			log.Printf("ENV_ERROR: %s", err.Error())
		}
		// extract config to struct
		if err := viper.Unmarshal(&Instance); err != nil {
			panic(fmt.Sprintf("ENV_ERROR: %s", err.Error()))
		}
	}
	// instance
	cfgOnce.Do(func() {

		// read env
		readEnv()
		// subs to event
		viper.OnConfigChange(func(e fsnotify.Event) {
			log.Printf("update configuration data . . . .")
			readEnv()
		})
		// watch file update
		viper.WatchConfig()
		// notify that config file is ready
		log.Println("configuration file: ready")
	})
}

func (cfg *AppConfig) updateEnvHelper(key, value any) {
	if err := viper.ReadInConfig(); err != nil {
		log.Println("READ", err.Error())
	}

	viper.Set(key.(string), value)

	viper.SetConfigType("dotenv")

	if err := viper.WriteConfig(); err != nil {
		log.Println("WRITE", err.Error())
	}
}
