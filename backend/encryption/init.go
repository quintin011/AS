package encryption

import (
	//"crypto/rand"
	//"crypto/rsa"
	//"crypto/x509"
	//"encoding/pem"
	"log"
	//"os"
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/matelang/jwt-go-aws-kms/v2/jwtkms"
)

var (
	cip_bytes  = []byte{53, 77, 119, 106, 121, 122, 103, 98, 36, 36, 36, 74, 113, 89, 116, 106}
	enc_secret = []byte{79, 37, 42, 114, 42, 118, 65, 85, 67, 82, 54, 78, 88, 48, 49, 33}
	// pubkey *rsa.PublicKey
	// prikey *rsa.PrivateKey
	svc    *kms.Client
	keyId  string
	kmscfg *jwtkms.Config
)

func init() {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion("ap-southeast-2"))
	if err != nil {
		log.Panic(err)
	}

	svc = kms.NewFromConfig(cfg)
	keyId = "arn:aws:kms:ap-southeast-2:170164929417:key/ec055ebf-a2e4-46e5-872f-df8d08014675"
	kmscfg = jwtkms.NewKMSConfig(svc, keyId, false)
	// _, prierr := os.Stat("private.pem")
	// _, puberr := os.Stat("public.pem")
	// if os.IsNotExist(puberr) && os.IsNotExist(prierr) {
	// 	prikey,pubkey = GenKey(2048)
	// 	prib,pubb := PrikeyBytes(prikey), PubkeyBytes(pubkey)
	// 	err := os.WriteFile("private.pem", prib, 0644)
	// 	if err != nil {
	// 		log.Panic(err)
	// 	}
	// 	err = os.WriteFile("public.pem", pubb, 0644)
	// 	if err != nil {
	// 		log.Panic(err)
	// 	}
	// } else if os.IsNotExist(puberr) && os.IsExist(prierr) {
	// 	log.Panic("please import private/public key under root directory")
	// } else if os.IsExist(puberr) && os.IsNotExist(prierr) {
	// 	log.Panic("please import private/public key under root directory")
	// } else {
	// 	prib, err := os.ReadFile("private.pem")
	// 	if err != nil {
	// 		log.Panic(err)
	// 	}
	// 	prikey = BytePrikey(prib)
	// 	pubb, err := os.ReadFile("public.pem")
	// 	if err != nil {
	// 		log.Panic(err)
	// 	}
	// 	pubkey = BytePubkey(pubb)

	// }
	// prikey = BytePrikey(prib)
	// pubkey = BytePubkey(pubb)
}

// func GenKey(bits int) (*rsa.PrivateKey, *rsa.PublicKey) {
// 	prikey,err := rsa.GenerateKey(rand.Reader, bits)
// 	if err != nil {
// 		log.Panic(err)
// 	}
// 	return prikey, &prikey.PublicKey
// }

// func PrikeyBytes(pri *rsa.PrivateKey) []byte {
// 	B := pem.EncodeToMemory(
// 		&pem.Block{
// 			Type: "RSA PRIVATE KEY",
// 			Bytes:x509.MarshalPKCS1PrivateKey(pri),
// 		})
// 	return B
// }

// func PubkeyBytes(pub *rsa.PublicKey) []byte {
// 	pubASN, err := x509.MarshalPKIXPublicKey(pub)
// 	if err != nil {
// 		log.Panic(err)
// 	}
// 	B := pem.EncodeToMemory(&pem.Block{
// 		Type:  "RSA PUBLIC KEY",
// 		Bytes: pubASN,
// 	})
// 	return B
// }

// func BytePrikey(pri []byte) *rsa.PrivateKey {
// 	bl, _ := pem.Decode(pri)
// 	key, err := x509.ParsePKCS1PrivateKey(bl.Bytes)
// 	if err != nil {
// 		log.Panic(err)
// 	}
// 	return key
// }

// func BytePubkey(pub []byte) *rsa.PublicKey {
// 	bl, _ := pem.Decode(pub)
// 	keys,err:= x509.ParsePKIXPublicKey(bl.Bytes)
// 	if err != nil {
// 		log.Panic(err)
// 	}
// 	key := keys.(*rsa.PublicKey)
// 	return key
// }
