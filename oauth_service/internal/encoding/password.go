package encoding

import "golang.org/x/crypto/bcrypt"

func Encode(password string) (string, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}

	return string(hashBytes), nil
}

func Check(password, PasswordHash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(PasswordHash), []byte(password)) == nil
}
