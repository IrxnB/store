package jwt

import "math/big"

type JWK struct {
	N *big.Int `json:"N"`
	E int      `json:"E"`
}
