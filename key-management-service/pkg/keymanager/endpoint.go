package keymanager

import (
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
	eError "github.com/mdshahjahanmiah/explore-go/error"
	"net/http"
)

type PublicKeyResponse struct {
	X string `json:"x"`
	Y string `json:"y"`
}

type PairingParamResponse struct {
	Params string `json:"params"`
}

func getPublicKeyEndpoint(service KeyManagementService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		publicKey := service.GetPublicKey()
		if publicKey == nil {
			return nil, eError.NewServiceError(errors.New("failed to get public key"), "Internal_Error", "NONE", http.StatusInternalServerError)
		}

		// elliptic curve cryptography (ECC), a public key is a point on the elliptic curve
		// represented by the X and Y coordinates of the point
		return PublicKeyResponse{
			X: publicKey.X().String(),
			Y: publicKey.Y().String(),
		}, nil
	}
}

func getShareKeyEndpoint(service KeyManagementService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		shares := service.GetKeyShares()
		if len(shares) < 1 {
			return nil, eError.NewServiceError(errors.New("failed to get key shares"), "Internal_Error", "NONE", http.StatusInternalServerError)
		}
		return shares, nil
	}
}

func getPairingParamsEndpoint(service KeyManagementService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		pairingParam := service.GetPairingParams()
		if pairingParam == "" {
			return nil, eError.NewServiceError(errors.New("failed to get pairing params"), "Internal_Error", "NONE", http.StatusInternalServerError)
		}
		return PairingParamResponse{
			Params: pairingParam,
		}, nil
	}
}
