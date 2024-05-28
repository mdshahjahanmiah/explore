package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/mdshahjahanmiah/explore-go/di"
	eHttp "github.com/mdshahjahanmiah/explore-go/http"
	"github.com/mdshahjahanmiah/explore-go/logging"
	"github.com/mdshahjahanmiah/gateway-service/pkg/config"
	"github.com/mdshahjahanmiah/gateway-service/pkg/routes"
	"github.com/mdshahjahanmiah/gateway-service/pkg/services"
	"go.uber.org/dig"
	"log/slog"
	"time"
)

func main() {
	c := di.New()

	c.Provide(func() (config.Config, error) {
		conf, err := config.Load()
		if err != nil {
			slog.Error("failed to load configuration", "err", err)
			return config.Config{}, err
		}
		return conf, nil
	})

	c.Provide(func(conf config.Config) (*logging.Logger, error) {
		logger, err := logging.NewLogger(conf.LoggerConfig)
		if err != nil {
			slog.Error("initializing logger", "err", err)
			return nil, err
		}

		return logger, nil
	})

	c.Provide(func(conf config.Config, logger *logging.Logger) services.KmsService {
		kmsService, err := services.NewKmsService(conf.KmsHttpAddress, 10*time.Second)
		if err != nil {
			logger.Fatal("initializing kms service", "err", err)
		}
		return kmsService
	})

	c.Provide(func(conf config.Config, logger *logging.Logger, kmsService services.KmsService) services.DsService {
		keyShares, err := kmsService.GetShares()
		if err != nil {
			logger.Fatal("getting key shares", "err", err)
		}
		return services.NewDsService(conf.DsHttpAddress, 10*time.Second, keyShares)
	})

	c.Provide(func(logger *logging.Logger, kmsService services.KmsService, dsService services.DsService) eHttp.Endpoint {
		r := chi.NewRouter()
		routes.RegisterRoutes(r, logger, kmsService, dsService)
		return eHttp.Endpoint{Pattern: "/*", Handler: r}
	}, dig.Group("endpoint"))

	c.Provide(func(config config.Config) *eHttp.ServerConfig {
		return &eHttp.ServerConfig{
			HttpAddress: config.HttpAddress,
		}
	})

	c.ProvideMonitoringEndpoints("endpoint")

	// Invoke the server with the provided dependencies
	c.Invoke(func(in struct {
		dig.In
		Conf         config.Config
		ServerConfig *eHttp.ServerConfig
		Endpoints    []eHttp.Endpoint `group:"endpoint"`
	}) {
		// Use NewServer to create the server
		server := eHttp.NewServer(in.ServerConfig, in.Endpoints, nil)
		c.Provide(func() di.StartCloser { return server }, dig.Group("startclose"))
	})

	c.Start()
}
