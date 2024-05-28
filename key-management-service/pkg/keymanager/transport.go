package keymanager

import (
	"github.com/go-chi/chi/v5"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/mdshahjahanmiah/explore-go/error"
	"github.com/mdshahjahanmiah/explore-go/http"
)

func MakeHandler(service KeyManagementService) http.Endpoint {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(error.EncodeError),
	}

	handlePublicKey := kithttp.NewServer(
		getPublicKeyEndpoint(service),
		kithttp.NopRequestDecoder,
		kithttp.EncodeJSONResponse,
		opts...,
	)

	handleShare := kithttp.NewServer(
		getShareKeyEndpoint(service),
		kithttp.NopRequestDecoder,
		kithttp.EncodeJSONResponse,
		opts...,
	)

	handlePairingParam := kithttp.NewServer(
		getPairingParamsEndpoint(service),
		kithttp.NopRequestDecoder,
		kithttp.EncodeJSONResponse,
		opts...,
	)

	r := chi.NewRouter()
	r.Method("GET", "/public-key", handlePublicKey)
	r.Method("GET", "/key-shares", handleShare)
	r.Method("GET", "/pairing-param", handlePairingParam)

	return http.Endpoint{Pattern: "/*", Handler: r}
}
