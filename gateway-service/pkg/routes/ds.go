package routes

import (
	"github.com/go-chi/chi/v5"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/mdshahjahanmiah/explore-go/error"
	"github.com/mdshahjahanmiah/explore-go/logging"
	"github.com/mdshahjahanmiah/gateway-service/pkg/handlers"
	"github.com/mdshahjahanmiah/gateway-service/pkg/services"
)

// RegisterDecryptionRoutes RegisterKmsRoutes registers the KMS routes with the given router.
func RegisterDecryptionRoutes(r chi.Router, logger *logging.Logger, dsService services.DsService) {
	r.Route("/ds", func(r chi.Router) {
		opts := []kithttp.ServerOption{
			kithttp.ServerErrorEncoder(error.EncodeError),
		}

		handleCiphertext := kithttp.NewServer(
			handlers.GetCiphertextEndpoint(logger, dsService),
			handlers.NoRequestDecoder,
			kithttp.EncodeJSONResponse,
			opts...,
		)

		handleDecrypt := kithttp.NewServer(
			handlers.GetDecryptEndpoint(logger, dsService),
			handlers.DecodeDecryptRequest,
			kithttp.EncodeJSONResponse,
			opts...,
		)
		r.Get("/ciphertext", handleCiphertext.ServeHTTP)
		r.Post("/decrypt", handleDecrypt.ServeHTTP)
	})
}
