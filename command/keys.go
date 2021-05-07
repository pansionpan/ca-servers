package command

import (
	"errors"
	"github.com/liwangqiang/gmsm/sm2"
	"iauto.com/asn1c-oer-lib/v2xca"
	"log"
)

// ReadPrivateKeyFromPem 读取路径中的私钥
func ReadPrivateKeyFromPem(path string) (*v2xca.PrivateKey, error) {
	sm2PrivKey, err := sm2.ReadPrivateKeyFromPem(path, nil)
	if err != nil {
		return nil, errors.New("读取私钥文件失败")
	}

	return (*v2xca.PrivateKey)(sm2PrivKey), nil
}

// CreatePublicKeyAndPrivateKey 生成公钥私钥对
func CreatePublicKeyAndPrivateKey() (*v2xca.PrivateKey, *v2xca.PublicVerifyKey) {
	privKey, err := sm2.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}
	pubKey := v2xca.PublicVerifyKey(privKey.PublicKey)

	return (*v2xca.PrivateKey)(privKey), &pubKey
}
