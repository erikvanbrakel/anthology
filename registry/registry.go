package registry

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/erikvanbrakel/anthology/models"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"io"
)

type Registry interface {
	GetModuleData(namespace, name, provider, version string) (reader *bytes.Buffer, err error)
	ListModules(namespace, name, provider string, offset, limit int) (modules []models.Module, total int, err error)
	PublishModule(namespace, name, provider, version string, data io.Reader) (err error)
	GetProviderMetaData(namespace, name, version, OS, arch string) (providerMetaData models.ProviderDownload, err error)
	GetProviderData(namespace, name, version, OS, arch, file string) (reader *bytes.Buffer, err error)
	ListProviders(namespace, name string, offset, limit int) (modules []models.Provider, total int, err error)
	PublishProvider(namespace, name, version, OS, arch string, data io.Reader) (err error)
}

func sha256File(f io.Reader) (string, error) {

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func getGPGKeys(pubringFile io.Reader) ([]models.GPGKeys, error) {

	pubKeyRing, err := openpgp.ReadArmoredKeyRing(pubringFile)

	if err != nil {
		return nil, err
	}

	var gpgKeys []models.GPGKeys

	for _, key := range pubKeyRing {
		// Fingerprint
		fingerprint := hex.EncodeToString(key.PrimaryKey.Fingerprint[:])

		// Public Key
		var pubKeyBuf bytes.Buffer
		err = key.Serialize(&pubKeyBuf)

		if err != nil {
			return nil, err
		}

		armored := bytes.NewBuffer(nil)
		encBuf, err := armor.Encode(armored, openpgp.PublicKeyType, nil)
		if err != nil {
			return nil, err
		}
		_, err = encBuf.Write(pubKeyBuf.Bytes())
		if err != nil {
			return nil, err
		}
		encBuf.Close()

		gpgKeys = append(gpgKeys, models.GPGKeys{
			KeyID:      fingerprint,
			ASCIIArmor: armored.String(),
		})
	}
	return gpgKeys, nil
}
