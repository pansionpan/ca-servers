package command

import (
	"fmt"
	"github.com/liwangqiang/gmsm/sm2"
	"iauto.com/asn1c-oer-lib/v2xca"
	"os"

	"log"
	"testing"
)

/*
测试CreateCertRequest方法:
	1.测试当私钥参数为空时
	2.测试当证书请求中subjectType为空时
	3.测试当证书请求中subjectName为空时
	4.测试当证书请求中起始时间fromData为空时
	5.测试当证书请求中终止时间toData为空时
*/

func TestStringToTime(t *testing.T) {
	input := "20050101"
	got, _ := stringToTime(input)
	gotUTC := got.UnixNano()/1e9
	inputUTC := int64(1104537600)
	if gotUTC != inputUTC {
		t.Fatal("the input does not match the got")
	}
}

func TestCreateCertRequestWithoutPrivatekey(t *testing.T) {
	expected := "private key cannot be null"

	subjectName := "TestWithNoPrivateKey"
	fromDate := "20100101"
	toData := "20300101"

	_, err := CreateCertRequest(nil, v2xca.SubjectTypeEnrollmentCredential, subjectName, fromDate, toData)
	errContent := fmt.Sprintf("%v", err)
	if errContent != expected {
		t.Errorf("expected does not match the actual: PrivateKey is null")
	}
}

func TestCreateCertRequestWithoutSubjectType(t *testing.T) {
	expected := "subjectType is wrong"

	privatekey := createPrivateKey()
	subjectName := "TestWithNoPrivateKey"
	fromDate := "20100101"
	toData := "20300101"

	_, err := CreateCertRequest(privatekey, 9, subjectName, fromDate, toData)
	errContent := fmt.Sprintf("%v", err)
	if errContent != expected {
		t.Errorf("expected does not match the actual: SubjectType is null")
	}
}

func TestCreateCertRequestWithoutSubjectName(t *testing.T) {
	expected := "subjectName is wrong"

	privatekey := createPrivateKey()
	subjectName := ""
	fromDate := "20100101"
	toData := "20300101"

	_, err := CreateCertRequest(privatekey, v2xca.SubjectTypeRootCa, subjectName, fromDate, toData)
	errContent := fmt.Sprintf("%v", err)
	if errContent != expected {
		t.Errorf("expected does not match the actual: SubjectName is null")
	}
}

func TestCreateCertRequestWithoutFromData(t *testing.T) {
	expected := "the timeData is wrong"

	privatekey := createPrivateKey()
	subjectName := "TestWithNoFromData"
	fromDate := ""
	toData := "20300101"

	_, err := CreateCertRequest(privatekey, v2xca.SubjectTypeRootCa, subjectName, fromDate, toData)
	errContent := fmt.Sprintf("%v", err)
	if errContent != expected {
		t.Errorf("expected does not match the actual: FromData is null")
	}
}

func TestCreateCertRequestWithoutToData(t *testing.T) {
	expected := "the timeData is wrong"

	privatekey := createPrivateKey()
	subjectName := "TestWithNoToData"
	fromDate := "20100101"
	toData := ""

	_, err := CreateCertRequest(privatekey, v2xca.SubjectTypeRootCa, subjectName, fromDate, toData)
	errContent := fmt.Sprintf("%v", err)
	if errContent != expected {
		t.Errorf("expected does not match the actual: ToData is null")
	}
}

//TestReadCertRequestAndWriteCertRequest 测试读取证书请求方法与存储证书请求方法
func TestReadCertRequestAndWriteCertRequest(t *testing.T) {
	privatekey := createPrivateKey()
	subjectName := "TestReadCertRequest"
	fromDate := "20200101"
	toData := "20300101"
	req, _ := CreateCertRequest(privatekey, v2xca.SubjectTypeEnrollmentCredential, subjectName, fromDate, toData)

	err := WriteCertRequest("./testCerReq.csr", req)
	if err != nil {
		log.Fatal(err)
	}

	expected := req
	actual, err := ReadCertRequest("./testCerReq.csr")
	if err != nil {
		log.Fatal(err)
	}
	expecteds := fmt.Sprintf("%v", expected)
	actuals := fmt.Sprintf("%v", actual)
	_ = os.Remove("./testCerReq.csr")
	if expecteds != actuals {
		log.Fatal("TestReadCertRequestAndWriteCertRequest方法验证失败")
	}
}

func createPrivateKey() *sm2.PrivateKey {
	privKey, err := sm2.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}
	return privKey
}
