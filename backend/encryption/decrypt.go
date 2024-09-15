package encryption

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
)

func Decode(s string) []byte {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		log.Panic(err)
		return nil
	}
	return data
}

func DecryptPrikey(cip_t []byte) []byte {
	// hash := sha512.New()
	// t, err := rsa.DecryptOAEP(hash, rand.Reader, pri, cip_t, nil)
	// if err != nil {
	// 	log.Panic(err)
	// }
	kmsout, err := svc.Decrypt(context.Background(), &kms.DecryptInput{
		CiphertextBlob:      cip_t,
		KeyId:               &keyId,
		EncryptionAlgorithm: types.EncryptionAlgorithmSpecRsaesOaepSha256,
	})
	if err != nil {
		log.Panic(err)
	}
	return kmsout.Plaintext
}

func Decrypt(t *string) (string, error) {
	bl, err := aes.NewCipher(enc_secret)
	if err != nil {
		return "", err
	}
	cip_t := Decode(*t)
	cfb := cipher.NewCFBDecrypter(bl, cip_bytes)
	rsa_t := make([]byte, len(cip_t))
	cfb.XORKeyStream(rsa_t, cip_t)

	plain := DecryptPrikey(rsa_t)
	return string(plain), nil
}
