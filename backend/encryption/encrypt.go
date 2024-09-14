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

func Encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func EncryptPubkey(t []byte, pub *rsa.PublicKey) []byte {
	hash := sha512.New()
	cip_t, err := rsa.EncryptOAEP(hash,rand.Reader, pub,t,nil)
	if err != nil {
		log.Panic(err)
	}
	return cip_t
}

func Encrypt(t []byte) (string, error) {
	rsa_t := EncryptPubkey(t, pubkey)
	bl, err := aes.NewCipher(enc_secret)
	if err != nil {
		return "", err
	}
	cfb := cipher.NewCFBEncrypter(bl,cip_bytes)
	cip_t := make([]byte, len(rsa_t))
	cfb.XORKeyStream(cip_t, rsa_t)
	return Encode(cip_t), nil
}
