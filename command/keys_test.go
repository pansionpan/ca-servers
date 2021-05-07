package command

import (
	"fmt"
	"github.com/liwangqiang/gmsm/sm2"
	"log"
	"os"
	"testing"
)

func TestReadPrivateKeyFromPem(t *testing.T) {
	privKey, _ := sm2.GenerateKey()
	ok, err := sm2.WritePrivateKeytoPem("./testPriv.key", privKey, nil)
	if !ok {
		log.Fatal(err)
	}

	readPriv, err := ReadPrivateKeyFromPem("./testPriv.key")

	expected := fmt.Sprintf("%v", privKey)
	actual := fmt.Sprintf("%v", readPriv)
	_ = os.Remove("./testPriv.key")
	if expected != actual {
		t.Errorf("expected does not match the actual: ReadPRivateKeyFromPem has something wromg")
	}
}
