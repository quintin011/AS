package encryption

import (
	"context"
	"crypto/aes"
	"crypto/cipher"

	// "crypto/rand"
	// "crypto/rsa"
	// "crypto/sha512"
	"encoding/base64"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
)

func Encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func EncryptPubkey(t []byte) []byte {
	// hash := sha512.New()
	// cip_t, err := rsa.EncryptOAEP(hash,rand.Reader, pub,t,nil)
	// if err != nil {
	// 	log.Panic(err)
	// }
	kmsout, err := svc.Encrypt(context.Background(), &kms.EncryptInput{
		KeyId:               &keyId,
		Plaintext:           t,
		EncryptionAlgorithm: types.EncryptionAlgorithmSpecRsaesOaepSha256,
	})
	if err != nil {
		log.Panic(err)
	}
	return kmsout.CiphertextBlob
}

func Encrypt(t []byte) (string, error) {
	rsa_t := EncryptPubkey(t)
	bl, err := aes.NewCipher(enc_secret)
	if err != nil {
		return "", err
	}
	cfb := cipher.NewCFBEncrypter(bl, cip_bytes)
	cip_t := make([]byte, len(rsa_t))
	cfb.XORKeyStream(cip_t, rsa_t)
	return Encode(cip_t), nil
}
