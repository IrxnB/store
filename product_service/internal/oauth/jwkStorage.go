package oauth

import (
	"crypto/rsa"
	"fmt"
	"sync"
	"time"
)

type JwkStorage struct {
	mx  sync.RWMutex
	key *rsa.PublicKey
	ja  JwkWebApi
}

func NewJwkStorage(api JwkWebApi) (storage *JwkStorage, err error) {
	storage = &JwkStorage{ja: api}
	err = storage.updateKey()
	if err != nil {
		return nil, err
	}

	go func() {
		timer := time.NewTicker(1 * time.Minute)
		for range timer.C {
			storage.updateKey()
		}
	}()
	return
}

func (js *JwkStorage) GetPublic() (*rsa.PublicKey, error) {
	js.mx.RLock()
	defer js.mx.RUnlock()

	if js.key == nil {
		return nil, fmt.Errorf("no key")
	}

	return js.key, nil
}

func (js *JwkStorage) updateKey() error {
	key, err := js.ja.GetJWK()

	if err != nil {
		return err
	}

	js.mx.Lock()
	defer js.mx.Unlock()

	js.key = &rsa.PublicKey{N: key.N, E: key.E}
	return nil
}
