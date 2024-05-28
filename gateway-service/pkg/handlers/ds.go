package handlers

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	eError "github.com/mdshahjahanmiah/explore-go/error"
	"github.com/mdshahjahanmiah/explore-go/logging"
	"github.com/mdshahjahanmiah/gateway-service/pkg/services"
	"net/http"
)

// GetDecryptEndpoint returns the endpoint for decrypting a ciphertext
func GetDecryptEndpoint(logger *logging.Logger, dsService services.DsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(services.DecryptRequest)

		decryptedMessage, err := dsService.Decrypt(req.Ciphertext)
		if err != nil {
			logger.Error("failed to decrypt message", "error", err)
			return nil, eError.NewTransportError(err, "INTERNAL_SERVER_ERROR")
		}

		return services.DecryptResponse{
			DecryptedMessage: decryptedMessage,
		}, nil
	}
}

// GetCiphertextEndpoint decodes a decrypt request from an HTTP request
func GetCiphertextEndpoint(logger *logging.Logger, dsService services.DsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		ciphertextResponse, err := dsService.Ciphertext()
		if err != nil {
			logger.Error("failed to generate cipher text message", "error", err)
			return nil, eError.NewTransportError(err, "INTERNAL_SERVER_ERROR")
		}

		return ciphertextResponse, nil
	}
}

// DecodeDecryptRequest decodes a decrypt request from an HTTP request
func DecodeDecryptRequest(ctx context.Context, request *http.Request) (interface{}, error) {
	decoder := json.NewDecoder(request.Body)

	var decryptRequest services.DecryptRequest
	err := decoder.Decode(&decryptRequest)
	if err != nil {
		return nil, eError.NewServiceError(err, "decode decrypt request", "payload", http.StatusBadRequest)
	}

	return services.DecryptRequest{
		Ciphertext: decryptRequest.Ciphertext,
	}, nil
}

// NoRequestDecoder decodes a ciphertext request from an HTTP request
func NoRequestDecoder(ctx context.Context, request *http.Request) (interface{}, error) {
	return nil, nil
}
