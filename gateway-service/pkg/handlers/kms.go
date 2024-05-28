package handlers

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	eError "github.com/mdshahjahanmiah/explore-go/error"
	"github.com/mdshahjahanmiah/explore-go/logging"
	"github.com/mdshahjahanmiah/gateway-service/pkg/services"
)

func GetPublicKeyEndpoint(logger *logging.Logger, kmsService services.KmsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		publicKey, err := kmsService.FetchPublicKey()
		if err != nil {
			return nil, eError.NewTransportError(err, "INTERNAL_SERVER_ERROR")
		}
		return publicKey, nil
	}
}
