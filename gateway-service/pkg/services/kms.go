package services

import (
	"encoding/json"
	httpclient "github.com/mdshahjahanmiah/gateway-service/pkg/client"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

// KmsService is a service that interacts with the KMS
type kmsService struct {
	keyShares []KeyShareResponse
	client    *httpclient.Client
	kmsURL    string
	mutex     sync.Mutex
}

type KmsService interface {
	FetchPublicKey() (PublicKeyResponse, error)
	GetShares() ([]KeyShareResponse, error)
}

// PublicKeyResponse is the response from the KMS for the public key
type PublicKeyResponse struct {
	X string `json:"x"`
	Y string `json:"y"`
}

// KeyShareResponse is the response from the KMS for the key shares
type KeyShareResponse struct {
	ID    int    `json:"id"`
	Share string `json:"share"`
}

// NewKmsService creates a new KmsService with the specified KMS URL and timeout
func NewKmsService(kmsURL string, timeout time.Duration) (KmsService, error) {
	client := httpclient.NewHttpClient(timeout)
	service := &kmsService{
		client: client,
		kmsURL: kmsURL,
	}

	// Fetch and cache key shares during initialization
	err := service.fetchAndCacheKeyShares()
	if err != nil {
		return nil, err
	}

	return service, nil
}

// FetchPublicKey fetches the public key from the KMS
func (kms *kmsService) FetchPublicKey() (PublicKeyResponse, error) {
	var publicKey []byte
	resp, err := kms.client.Get(kms.kmsURL + "/public-key")
	if err != nil {
		return PublicKeyResponse{}, err
	}
	defer resp.Body.Close()
	publicKey, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return PublicKeyResponse{}, err
	}

	var response PublicKeyResponse

	// Unmarshal the JSON data into the struct
	err = json.Unmarshal(publicKey, &response)
	if err != nil {
		return PublicKeyResponse{}, err
	}

	return response, err
}

// GetShares fetches the key shares from the KMS
func (kms *kmsService) GetShares() ([]KeyShareResponse, error) {
	kms.mutex.Lock()
	defer kms.mutex.Unlock()

	if len(kms.keyShares) == 0 {
		return nil, errors.New("no key shares available")
	}
	return kms.keyShares, nil
}

// fetchAndCacheKeyShares fetches key shares from the KMS and caches them
func (k *kmsService) fetchAndCacheKeyShares() error {
	resp, err := k.client.Get(k.kmsURL + "/key-shares")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to fetch key shares from KMS")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var keyShares []KeyShareResponse
	if err = json.Unmarshal(body, &keyShares); err != nil {
		return err
	}

	k.mutex.Lock()
	k.keyShares = keyShares
	k.mutex.Unlock()

	return nil
}
