package config

import (
	"flag"
	"github.com/mdshahjahanmiah/explore-go/logging"
)

type Config struct {
	HttpAddress     string
	SecurityLevel   string
	ThresholdConfig ThresholdConfig
	LoggerConfig    logging.LoggerConfig
}

type ThresholdConfig struct {
	Enabled     bool
	Threshold   int
	TotalShares int
}

func Load() (Config, error) {
	fs := flag.NewFlagSet("", flag.ExitOnError)

	httpAddress := fs.String("http.public.address", "0.0.0.0:9001", "HTTP listen address for all specified endpoints.")
	securityLevel := fs.String("security.level", "medium", "the security level to use for the pairing parameters. Possible values are 'low', 'medium', 'high'")

	thresholdConfig := ThresholdConfig{}
	fs.BoolVar(&thresholdConfig.Enabled, "thresholdconfig.enabled", true, "whether threshold encryption is enabled or not")
	fs.IntVar(&thresholdConfig.Threshold, "thresholdconfig.threshold", 4, "the threshold number of shares required to reconstruct the secret. For instance threshold 3 means that any 3 out of the total number of shares can be used to reconstruct the secret")
	fs.IntVar(&thresholdConfig.TotalShares, "thresholdconfig.shares", 5, "the total number of shares to be generated for instance n = 5 means that the secret will be split into 5 shares")

	loggerConfig := logging.LoggerConfig{}
	fs.StringVar(&loggerConfig.CommandHandler, "logger.handler.type", "json", "handler type e.g json, otherwise default will be text type")
	fs.StringVar(&loggerConfig.LogLevel, "logger.log.level", "debug", "log level wise logging with fatal log")

	config := Config{
		HttpAddress:     *httpAddress,
		SecurityLevel:   *securityLevel,
		ThresholdConfig: thresholdConfig,
		LoggerConfig:    loggerConfig,
	}

	return config, nil
}
