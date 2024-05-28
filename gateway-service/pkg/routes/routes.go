package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/mdshahjahanmiah/explore-go/logging"
	"github.com/mdshahjahanmiah/gateway-service/pkg/services"
)

// RegisterRoutes registers the KMS and decryption routes with the given router.
func RegisterRoutes(r chi.Router, logger *logging.Logger, kmsService services.KmsService, dsService services.DsService) {
	RegisterKmsRoutes(r, logger, kmsService)
	RegisterDecryptionRoutes(r, logger, dsService)
}
