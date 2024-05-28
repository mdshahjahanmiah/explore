package decrypt

import (
	"github.com/go-chi/chi/v5"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/mdshahjahanmiah/explore-go/error"
	"github.com/mdshahjahanmiah/explore-go/http"
	"github.com/mdshahjahanmiah/explore-go/logging"
)

func MakeHandler(logger *logging.Logger, service Service) http.Endpoint {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(error.EncodeError),
	}

	handleGetCiphertext := kithttp.NewServer(
		getCiphertextEndpoint(logger, service),
		kithttp.NopRequestDecoder,
		kithttp.EncodeJSONResponse,
		opts...,
	)

	handlePartialDecryption := kithttp.NewServer(
		getPartialDecryptEndpoint(logger, service),
		decodeDecryptRequest,
		kithttp.EncodeJSONResponse,
		opts...,
	)

	r := chi.NewRouter()
	r.Method("GET", "/ciphertext", handleGetCiphertext)
	r.Method("POST", "/partial-decrypt", handlePartialDecryption)

	return http.Endpoint{Pattern: "/*", Handler: r}
}
