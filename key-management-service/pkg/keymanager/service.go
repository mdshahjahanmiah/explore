package keymanager

import (
	"encoding/base64"
	"errors"
	"github.com/Nik-U/pbc"
	"github.com/hashicorp/vault/shamir"
	"github.com/mdshahjahanmiah/explore-go/logging"
	"github.com/mdshahjahanmiah/key-management-service/pkg/config"
)

type KeyShare struct {
	ID    int    `json:"id"`
	Share string `json:"share"`
}

type keyManagementService struct {
	encodedParams string
	publicKey     *pbc.Element
	shares        []KeyShare
	logger        *logging.Logger
}

type KeyManagementService interface {
	GetPublicKey() *pbc.Element
	GetPairingParams() string
	GetKeyShares() []KeyShare
}

// NewKeyManagementService initializes the key management service
func NewKeyManagementService(config config.Config, logger *logging.Logger) (KeyManagementService, error) {
	// Determine security parameters based on the security level from the configuration
	baseFieldSize, subgroupOrder := ToSecurityMeasures(config.SecurityLevel)

	// Generate pairing parameters using the determined security parameters
	params := pbc.GenerateA(baseFieldSize, subgroupOrder)
	if params == nil {
		logger.Fatal("failed to generate pairing parameters")
	}

	// Convert pairing parameters to a string and then encode in base64 for storage
	encodedParams := base64.StdEncoding.EncodeToString([]byte(params.String()))

	// Create a new pairing using the generated parameters
	pairing := params.NewPairing()
	if pairing == nil {
		logger.Fatal("failed to create pairing")
	}

	// Generate a random private key in the pairing's Zr field
	privateKey := pairing.NewZr().Rand()
	if privateKey == nil {
		logger.Fatal("failed to generate private key")
	}

	// Generate a random G2 element to be used as the generator
	g2Gen := pairing.NewG2().Rand()
	if g2Gen == nil {
		logger.Fatal("failed to generate G2 element")
	}

	// Generate the public key by raising the G2 generator to the power of the private key
	publicKey := pairing.NewG2().PowZn(g2Gen, privateKey)
	if publicKey == nil {
		logger.Fatal("failed to generate public key")
	}

	// Check if the public key is correctly generated
	if publicKey.Is0() {
		logger.Fatal("public key is zero")
	}

	// Convert the private key to a byte slice for further processing
	privateKeyBytes := privateKey.Bytes()
	if privateKeyBytes == nil {
		logger.Fatal("failed to convert private key to bytes")
	}

	// Prepare key shares using Shamir's Secret Sharing or other configuration
	shares, err := prepareShares(config, logger, privateKeyBytes)
	if err != nil {
		logger.Error("failed to prepare key shares", "err", err)
		return nil, err
	}

	return &keyManagementService{
		encodedParams: encodedParams,
		publicKey:     publicKey,
		shares:        shares,
		logger:        logger,
	}, nil
}

// prepareShares generates key shares from the private key bytes based on the configuration.
// If threshold sharing is enabled, it splits the private key into multiple shares using Shamir's Secret Sharing.
// Otherwise, it returns the private key as a single share.
func prepareShares(config config.Config, logger *logging.Logger, privateKeyBytes []byte) ([]KeyShare, error) {
	// Check if threshold sharing is enabled
	if config.ThresholdConfig.Enabled {
		// Validate that the threshold is not greater than the total shares
		if config.ThresholdConfig.Threshold > config.ThresholdConfig.TotalShares {
			return nil, errors.New("threshold cannot be greater than total shares")
		}
		// Validate that the threshold and total shares are greater than 0
		if config.ThresholdConfig.Threshold < 1 || config.ThresholdConfig.TotalShares < 1 {
			return nil, errors.New("threshold and total shares must be greater than 0")
		}

		// Split the private key into shares using Shamir's Secret Sharing
		sharesBytes, err := shamir.Split(privateKeyBytes, config.ThresholdConfig.TotalShares, config.ThresholdConfig.Threshold)
		if err != nil {
			return nil, errors.New("failed to split private key into shares: " + err.Error())
		}

		// Convert the shares to a slice of KeyShare structs
		shares := make([]KeyShare, len(sharesBytes))
		for i, share := range sharesBytes {
			shares[i] = KeyShare{
				ID:    i + 1,
				Share: base64.StdEncoding.EncodeToString(share),
			}
		}
		return shares, nil
	}

	// If threshold sharing is not enabled, return the single key share
	shares := []KeyShare{
		{
			ID:    1,
			Share: base64.StdEncoding.EncodeToString(privateKeyBytes),
		},
	}
	return shares, nil
}

// GetPublicKey returns the public key of the key management service.
func (kms *keyManagementService) GetPublicKey() *pbc.Element {
	return kms.publicKey
}

// GetPairingParams returns the base64-encoded pairing parameters of the key management service.
func (kms *keyManagementService) GetPairingParams() string {
	return kms.encodedParams
}

// GetKeyShares returns the key shares of the private key.
// The key shares are used in threshold cryptography, allowing a private key to be distributed among multiple parties.
func (kms *keyManagementService) GetKeyShares() []KeyShare {
	return kms.shares
}
