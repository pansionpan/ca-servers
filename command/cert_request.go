package command

import "C"
import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/liwangqiang/gmsm/sm2"
	"iauto.com/asn1c-oer-lib/v2xca"
)

// CreateCertRequest 创建证书请求
func CreateCertRequest(privatekey *sm2.PrivateKey, subjectType v2xca.SubjectType, subjectName string, fromDate string, toDate string) (*v2xca.CertRequest, error) {
	if privatekey == nil {
		fmt.Println("Private key cannot be null")
		return nil, errors.New("private key cannot be null")
	}
	pubKey := privatekey.PublicKey

	if subjectType != 0 && subjectType != 1 && subjectType != 2 && subjectType != 3 && subjectType != 4 && subjectType != 5 && subjectType != 6 && subjectType != 7 {
		fmt.Println("SubjectType is wrong")
		return nil, errors.New("subjectType is wrong")
	}

	if subjectName == "" {
		fmt.Println("SubjectName is wrong")
		return nil, errors.New("subjectName is wrong")
	}

	fromTime, err := stringToTime(fromDate)
	if err != nil {
		fmt.Println("Start time is wrong")
		return nil, err
	}
	toTime, err := stringToTime(toDate)
	if err != nil {
		fmt.Println("Termination time is wrong")
		return nil, err
	}
	if !fromTime.Before(toTime) {
		fmt.Println("Termination time is less than start time")
		return nil, errors.New("termination time is less than start time")
	}

	tbsCertList := v2xca.TbsCert{
		SubjectType: subjectType,
		SubjectName: subjectName,
		SubjectAttrs: v2xca.SubjectAttribute{
			VerifyKey: (v2xca.PublicVerifyKey)(pubKey),
		},
		Validity: v2xca.Validity{
			Period: v2xca.Period{
				StartTime: &fromTime,
				EndTime:   toTime,
			},
		},
	}
	tbsCerLists := make([]v2xca.TbsCert, 0)
	tbsCerLists = append(tbsCerLists, tbsCertList)

	// 因为在生成证书请求时候，currentTime参数中会包含时间中的纳秒，但是读取证书请求中会忽略纳秒，因此此处处理忽略当前时间点的纳秒
	t1 := time.Now().Year()   //年
	t2 := time.Now().Month()  //月
	t3 := time.Now().Day()    //日
	t4 := time.Now().Hour()   //小时
	t5 := time.Now().Minute() //分钟
	t6 := time.Now().Second() //秒

	currentTimeData := time.Date(t1, t2, t3, t4, t5, t6, 0, time.Local) //获取当前时间，返回当前时间Time
	nowTime := currentTimeData.UTC()                                    // 转换为UTC时间

	csr := v2xca.CertRequest{
		Version:     2,
		CurrentTime: nowTime, // UTC
		TbsCertList: tbsCerLists,
	}
	return &csr, nil
}

// ReadCertRequest 读取证书请求
func ReadCertRequest(path string) (*v2xca.CertRequest, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return v2xca.ParseCertRequest(bytes)
}

// WriteCertRequest 存储证书请求
func WriteCertRequest(paths string, cerReq *v2xca.CertRequest) error {
	bytes := cerReq.Bytes()
	err := ioutil.WriteFile(paths, bytes, 0644)
	return err
}

// ChooseSubjectType 处理subjectType
func ChooseSubjectType(number int) (v2xca.SubjectType, error) {
	switch number {
	case 1:
		return v2xca.SubjectTypeEnrollmentCredential, nil
	case 2:
		return v2xca.SubjectTypeAuthorizationTicket, nil
	case 3:
		return v2xca.SubjectTypeAuthorizationAuthority, nil
	case 4:
		return v2xca.SubjectTypeEnrollmentAuthority, nil
	case 5:
		return v2xca.SubjectTypeRootCa, nil
	case 6:
		return v2xca.SubjectTypeCrlSigner, nil
	case 7:
		return v2xca.SubjectTypePseudonymTicket, nil
	case 8:
		return v2xca.SubjectTypePseudonymAuthority, nil
	default:
		return (v2xca.SubjectType)(number), errors.New("the subjectType is wrong")
	}
}

// stringToTime 将字符串转为本地时间
func stringToTime(timeStr string) (time.Time, error) {
	timeLayout := "20060102"
	loc, err := time.LoadLocation("UTC")
	if err != nil {
		log.Fatal(err)
	}
	theTime, err := time.ParseInLocation(timeLayout, timeStr, loc)
	if err != nil {
		return time.Now(), errors.New("the timeData is wrong")
	}
	return theTime, nil
}