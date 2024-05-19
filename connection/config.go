package connection

import (
	"errors"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	Name     string
	Host     string
	User     string
	Password string
	Port     int
}

func defineFlags(flagSet *flag.FlagSet, workingDir string) {
	flagSet.String("database_host", "localhost", "Database host.")
	flagSet.String("database_user", "postgres", "Database user")
	flagSet.String("database_password", "", "Database password")
	flagSet.String("database_name", "postgres", "Database name")
	flagSet.String("database_port", "5432", "Database port")
}

func getConfig() (*Config, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	flags := flag.NewFlagSet("config", flag.ExitOnError)
	defineFlags(flags, cwd)

	if err = flags.Parse(os.Args[1:]); err != nil {
		return nil, err
	}

	viper.SetConfigName("ponzu") // ponzu config file
	viper.SetConfigType("props")
	viper.AddConfigPath(cwd)

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil && errors.As(err, &viper.ConfigFileNotFoundError{}) {
		log.Info("config file not found. will default to provided flags")
		err = nil
	}

	if err = viper.BindPFlags(flags); err != nil {
		return nil, err
	}

	return &Config{
		Host:     viper.GetString("database_host"),
		User:     viper.GetString("database_user"),
		Password: viper.GetString("database_password"),
		Port:     viper.GetInt("database_port"),
		Name:     viper.GetString("database_name"),
	}, nil
}
