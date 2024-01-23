package encoding

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
)

type KeyFileStorage struct {
	path            string
	privateFileName string
	publicFileName  string
	private         *rsa.PrivateKey
	public          *rsa.PublicKey
}

func NewKeyStorage(path, privateFileName, publicFileName string) (storage *KeyFileStorage, err error) {
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return nil, err
	}
	storage = &KeyFileStorage{path: path, privateFileName: privateFileName, publicFileName: publicFileName}

	err = storage.readKeys()
	if err != nil {
		err = storage.generateKeys()
		if err != nil {
			return nil, err
		}
	}

	return storage, nil
}

func (ks *KeyFileStorage) GetPublic() (*rsa.PublicKey, error) {
	if ks.public == nil {
		return &rsa.PublicKey{}, fmt.Errorf("no key")
	}
	return ks.public, nil
}
func (ks *KeyFileStorage) GetPrivate() (*rsa.PrivateKey, error) {
	if ks.private == nil {
		return &rsa.PrivateKey{}, fmt.Errorf("no key")
	}
	return ks.private, nil
}

func (ks *KeyFileStorage) generateKeys() (err error) {
	ks.private, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	ks.public = &ks.private.PublicKey

	err = ks.writePrivate()
	if err != nil {
		return err
	}
	err = ks.writePublic()
	if err != nil {
		return err
	}
	return nil
}

func (ks *KeyFileStorage) readKeys() (err error) {
	ks.private, err = ks.readPrivate()
	if err != nil {
		return err
	}

	ks.public, err = ks.readPublic()
	if err != nil {
		return err
	}
	return nil
}

func (ks *KeyFileStorage) readPrivate() (*rsa.PrivateKey, error) {
	keyFile, err := os.ReadFile(ks.getFilePath(ks.privateFileName))
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyFile)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("error decoding key")
	}

	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func (ks *KeyFileStorage) readPublic() (*rsa.PublicKey, error) {
	keyFile, err := os.ReadFile(ks.getFilePath(ks.publicFileName))
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyFile)
	if block == nil || block.Type != "RSA PUBLIC KEY" {
		return nil, fmt.Errorf("error decoding key")
	}

	key, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func (ks *KeyFileStorage) writePrivate() error {
	outFile, err := os.Create(ks.getFilePath(ks.privateFileName))
	if err != nil {
		return err
	}

	defer outFile.Close()

	bytes := x509.MarshalPKCS1PrivateKey(ks.private)
	block := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: bytes}

	return pem.Encode(outFile, block)
}

func (ks *KeyFileStorage) writePublic() error {

	bytes := x509.MarshalPKCS1PublicKey(ks.public)
	block := &pem.Block{Type: "RSA PUBLIC KEY", Bytes: bytes}

	outFile, err := os.Create(ks.getFilePath(ks.publicFileName))
	if err != nil {
		return err
	}

	defer outFile.Close()
	return pem.Encode(outFile, block)
}

func (ks *KeyFileStorage) getFilePath(name string) string {
	return filepath.Join(ks.path, name)
}
