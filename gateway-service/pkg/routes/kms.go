package routes

import (
	"github.com/go-chi/chi/v5"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/mdshahjahanmiah/explore-go/error"
	"github.com/mdshahjahanmiah/explore-go/logging"
	"github.com/mdshahjahanmiah/gateway-service/pkg/handlers"
	"github.com/mdshahjahanmiah/gateway-service/pkg/services"
)

// RegisterKmsRoutes registers the KMS routes with the given router.
func RegisterKmsRoutes(r chi.Router, logger *logging.Logger, kmsService services.KmsService) {
	r.Route("/kms", func(r chi.Router) {
		opts := []kithttp.ServerOption{
			kithttp.ServerErrorEncoder(error.EncodeError),
		}

		handlePublicKey := kithttp.NewServer(
			handlers.GetPublicKeyEndpoint(logger, kmsService),
			kithttp.NopRequestDecoder,
			kithttp.EncodeJSONResponse,
			opts...,
		)

		r.Get("/public-key", handlePublicKey.ServeHTTP)
	})
}
