package jwt

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func generateSecret() string {
	secret := make([]byte, 64)
	_, err := rand.Read(secret)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(secret)
}

func main() {
	secret := generateSecret()
	fmt.Println("Generated JWT Secret:", secret)
}
