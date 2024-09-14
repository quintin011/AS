package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"encoding/base64"
	"log"
)

func Decode(s string) []byte {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		log.Panic(err)
		return nil
	}
	return data 
}

func DecryptPrikey(cip_t []byte,pri *rsa.PrivateKey) []byte {
	hash := sha512.New()
	t, err := rsa.DecryptOAEP(hash, rand.Reader, pri, cip_t, nil)
	if err != nil {
		log.Panic(err)
	}
	return t
}

func Decrypt(t *string) (string, error) {
	bl, err := aes.NewCipher(enc_secret)
	if err != nil {
		return "",err
	}
	cip_t := Decode(*t)
	cfb := cipher.NewCFBDecrypter(bl,cip_bytes)
	rsa_t := make([]byte, len(cip_t))
	cfb.XORKeyStream(rsa_t, cip_t)

	plain := DecryptPrikey(rsa_t, prikey)
	return string(plain), nil
}
