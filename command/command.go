/*
command 包用于存储命令行工具可以用到的所有函数，
其函数包括：
	1.sm2                 创建私钥或从私钥中取出公钥
	2.req   			  生成证书请求
	3.ca                  生成根证书或二级证书
主要包含四个文件，分别为：
	1.v2xca sm2 -genkey -out privkey.key                                               生成 sm2 private key
	2.v2xca sm2 -pubout -in privkey.key -out pubkey.key                                从 private key 中取出 public key
	3.v2xca req -new -key user.key -out user.csr                                       生成证书请求
	4.v2xca req -print -in user.csr                                                    打印证书请求
	5.v2xca ca -selfsign -in rootca.csr -signkey priv.key -out cert.oer                生成自签证书
	6.v2xca ca -subca -in subca.csr -signkey priv.key -cert rootca.oer -out cert.oer   生成中级证书
	7.v2xca ca -print -in cert.oer                                                     打印证书
*/
package command

import (
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/liwangqiang/gmsm/sm2"
	"github.com/urfave/cli"
	"iauto.com/asn1c-oer-lib/v2xca"
	"io/ioutil"
	"os"
)

func CreateKeys(c *cli.Context) error {
	filePath := c.String("out")     // -out 输出文件路径
	privPath := c.String("in")      // -in  存储private key的文件
	existGenkey := c.Bool("genkey") // -genkey
	existPubout := c.Bool("pubout") // -pubout

	if existGenkey && privPath != "" {
		fmt.Println("Parameter mismatch!")
		return errors.New("parameter mismatch")
	}

	// 当参数为-genkey时
	if existGenkey && !existPubout {
		privatekey, _ := CreatePublicKeyAndPrivateKey()

		if filePath == "" {
			der, err := sm2.MarshalSm2PrivateKey((*sm2.PrivateKey)(privatekey), nil)
			if err != nil {
				fmt.Println("Failed to create private key")
				return err
			}
			var block *pem.Block
			block = &pem.Block{
				Type:    "PRIVATE KEY",
				Bytes:   der,
			}

			file := "test.key"
			files, err := os.Create(file)
			if err != nil {
				fmt.Println("Fail to show private key")
				return errors.New("fail to show private key")
			}
			defer files.Close()
			err = pem.Encode(files, block)
			if err != nil {
				fmt.Println("Fail to show private key")
				return errors.New("fail to show private key")
			}

			data, err := ioutil.ReadFile("test.key")
			if err != nil {
				fmt.Println("Fail to show private key")
				return errors.New("fail to show private key")
			}

			_ = os.Remove("test.key")
			fmt.Println(string(data))
			return nil
		}
		if len(filePath) < 4 {
			fmt.Println("Fail to name: name is too short")
			return errors.New("fail to name: name is too short")
		}
		if first := filePath[len(filePath)-4:]; first != ".key"{
			fmt.Println("Private key naming format error")
			return errors.New("private key naming format error")
		}

		_, err := sm2.WritePrivateKeytoPem(filePath, (*sm2.PrivateKey)(privatekey), nil)
		if err != nil {
			fmt.Println("Failed to output private key file")
			return err
		}
		fmt.Println("Successful private key generation")
		return nil
	} else if !existGenkey && existPubout {
		// -pubout
		if privPath == "" {
			fmt.Println("The private key is null, Missing option: '-in'.See 'v2xca sm2 -help'")
			return errors.New("the private key is null")
		}

		privKey, err := ReadPrivateKeyFromPem(privPath)
		if err != nil {
			fmt.Println("Private key fetch error")
			return err
		}
		pubKey := v2xca.PublicVerifyKey(privKey.PublicKey)

		if filePath == "" {
			der, err := sm2.MarshalSm2PublicKey((*sm2.PublicKey)(&pubKey))
			if err != nil {
				fmt.Println("Failed to create public key")
				return err
			}
			var block *pem.Block
			block = &pem.Block{
				Type:    "PUBLIC KEY",
				Bytes:   der,
			}

			file := "test.key"
			files, err := os.Create(file)
			if err != nil {
				fmt.Println("Fail to show private key")
				return errors.New("fail to show private key")
			}
			defer files.Close()
			err = pem.Encode(files, block)
			if err != nil {
				fmt.Println("Fail to show private key")
				return errors.New("fail to show private key")
			}

			data, err := ioutil.ReadFile("test.key")
			if err != nil {
				fmt.Println("Fail to show private key")
				return errors.New("fail to show private key")
			}

			_ = os.Remove("test.key")
			fmt.Println(string(data))
			return nil
		}
		if len(filePath) < 4 {
			fmt.Println("Fail to name: name is too short")
			return errors.New("fail to name: name is too short")
		}
		if first := filePath[len(filePath)-4:]; first != ".key" {
			fmt.Println("Public key naming format error")
			return errors.New("public key naming format error")
		}

		_, err = sm2.WritePublicKeytoPem(filePath, (*sm2.PublicKey)(&pubKey), nil)
		if err != nil {
			fmt.Println("Failed to output public key file")
			return err
		}
		fmt.Println("Successful public key generation")
		return nil
	}

	fmt.Println("Error parameters, Missing option, See 'v2xca sm2 -help'")
	return errors.New("error parameters")
}

// CreateCertRequsetCommand 生成证书请求
func CreateCertRequestCommand(c *cli.Context) error {
	existNew := c.Bool("new")     // 是否存在new关键字
	existPrint := c.Bool("print") // 是否存在print
	keyPath := c.String("key")    // 指定私钥位置
	outPath := c.String("out")    // 输出文件位置
	inPath := c.String("in")      // 指定打印的 csr 文件

	if (existNew && existPrint) || (!existNew && !existPrint) {
		fmt.Println("Error parameters, Missing option, See 'v2xca req -help'")
		return errors.New("error parameters")
	}

	if existNew && inPath != "" {
		fmt.Println("Parameter mismatch!")
		return errors.New("parameter mismatch")
	}

	// 当参数为-new时
	if existNew && !existPrint {
		if keyPath == "" {
			fmt.Println("The private key is null")
			return errors.New("the private key is empty")
		}
		// 先读取私钥位置
		userpriv, err := sm2.ReadPrivateKeyFromPem(keyPath, nil)
		if err != nil {
			fmt.Println("Private key reading error")
			return err
		}

		if outPath == "" && keyPath != "" {
			userinput, err := GetUserInput()
			if err != nil {
				fmt.Println(err)
				return err
			}
			num, err := ChooseSubjectType(userinput.Number)
			if err != nil {
				fmt.Println(err)
				return err
			}
			cer, err := CreateCertRequest(userpriv, num, userinput.SubjectName, userinput.FromDate, userinput.ToDate)
			if err != nil {
				fmt.Println("Failed to create certificate request")
				return err
			}
			v2xca.PrintCertRequest(cer)
			return nil
		}

		if keyPath != "" && outPath != "" {
			if len(outPath) < 4 {
				fmt.Println("Fail to name: name is too short")
				return errors.New("fail to name: name is too short")
			}
			if first := outPath[len(outPath)-4:]; first != ".csr" {
				fmt.Println("Public key naming format error")
				return errors.New("public key naming format error")
			}
			userinput, err := GetUserInput()
			if err != nil {
				fmt.Println(err)
				return err
			}
			num, err := ChooseSubjectType(userinput.Number)
			if err != nil {
				fmt.Println(err)
				return err
			}
			cer, err := CreateCertRequest(userpriv, num, userinput.SubjectName, userinput.FromDate, userinput.ToDate)
			if err != nil {
				fmt.Println("Failed to create certificate request")
				return err
			}

			err = WriteCertRequest(outPath, cer)
			if err != nil {
				fmt.Println("Output Certificate Request Failed")
				return err
			}
			fmt.Println("Certificate Request Stored")
			return nil
		}
		fmt.Println("Error parameters, Missing option, See 'v2xca req -help'")
		return errors.New("error parameters")
	}

	if !existNew && existPrint {
		if inPath == "" {
			fmt.Println("The specify certificate request is null")
			return errors.New("the specify certificate request is null")
		}

		cer, err := ReadCertRequest(inPath)
		if err != nil {
			fmt.Println("Failed to read certificate request file")
			return err
		}
		v2xca.PrintCertRequest(cer)
		return nil
	}

	fmt.Println("Error parameter, Missing option, See 'v2xca req -help'")
	return errors.New("error parameter")
}

// CreateRootCaAndPrivatekey 生成证书
func CreateRootAndSubCert(c *cli.Context) error {
	existSelf := c.Bool("selfsign") // 是否是自签证书
	existSub := c.Bool("subca")     // 是否是二级证书
	existPrint := c.Bool("print")   // 是否是打印输出
	inPath := c.String("in")        // 指定证书请求文件
	keyPath := c.String("signkey")  // 指定签名的私钥文件
	outPath := c.String("out")      // 指定输出文件路径
	certPath := c.String("cert")    // 指定 issuer 的证书

	if existSelf && certPath != "" {
		fmt.Println("Parameter mismatch")
		return errors.New("parameter mismatch")
	}

	if existSelf && !existSub && !existPrint {
		// 生成自签证书
		if inPath == "" || keyPath == "" {
			fmt.Println("Certificate request or private key cannot be empty")
			return errors.New("certificate request file or private key file cannot be empty")
		}

		// 读取自签证书请求文件位置
		selfCertReq, err := ReadCertRequest(inPath)
		if err != nil {
			fmt.Println("Failed to read certificate request")
			return err
		}

		// 读取私钥文件位置
		selfPrivKey, err := ReadPrivateKeyFromPem(keyPath)
		if err != nil {
			fmt.Println("Failed to read private key")
			return err
		}

		if selfCertReq.TbsCertList[0].SubjectType != v2xca.SubjectTypeRootCa {
			fmt.Println("Certificate request does not accord with root certificate requirements")
			return errors.New("certificate request does not accord with root certificate requirements")
		}

		public1 := fmt.Sprintf("%v", (sm2.PublicKey)(selfCertReq.TbsCertList[0].SubjectAttrs.VerifyKey))
		public2 := fmt.Sprintf("%v", selfPrivKey.PublicKey)
		if public1 != public2 {
			fmt.Println("The content of the certificate request does not match the content of the private key")
			return errors.New("the content of the certificate request does not match the content of the private key")
		}

		resp, err := v2xca.CreateCertificate(&selfCertReq.TbsCertList[0], nil, selfPrivKey)
		if err != nil {
			fmt.Println("Failed to create root certificate")
			return err
		}

		if outPath == "" {
			v2xca.PrintCert(resp)
			return nil
		}
		if len(outPath) < 4 {
			fmt.Println("Fail to name: name is too short")
			return errors.New("fail to name: name is too short")
		}
		if first := outPath[len(outPath)-4:]; first != ".oer" {
			fmt.Println("Certificate naming format error")
			return errors.New("certificate naming format error")
		}

		errs := v2xca.WriteCert(resp, outPath)
		if errs != nil {
			fmt.Println("Output Certificate Failed")
			return errs
		}
		fmt.Println("Create root certificate successfully")
		return nil
	}

	if !existSelf && existSub && !existPrint {
		// 生成二级证书
		if inPath == "" || keyPath == "" || certPath == "" {
			fmt.Println("Certificate request or private key or parent certificate cannot not be null")
			return errors.New("certificate request or private key or parent certificate cannot not be null")
		}

		subCert, err := v2xca.ReadCert(certPath)
		if err != nil {
			fmt.Println("Failed to read parent certificate request")
			return err
		}

		subCertReq, err := ReadCertRequest(inPath)
		if err != nil {
			fmt.Println("Failed to read certificate request")
			return err
		}

		// 读取私钥文件位置
		subPrivKey, err := ReadPrivateKeyFromPem(keyPath)
		if err != nil {
			fmt.Println("Failed to read private key")
			return err
		}

		if subCertReq.TbsCertList[0].SubjectType == v2xca.SubjectTypeRootCa {
			fmt.Println("The specific certificate-request's subjectType cannot be 'RootCA' ")
			return errors.New("the specific certificate-request's subjectType cannot be 'RootCA' ")
		}

		public1 := fmt.Sprintf("%v", (sm2.PublicKey)(subCertReq.TbsCertList[0].SubjectAttrs.VerifyKey))
		public2 := fmt.Sprintf("%v", subPrivKey.PublicKey)
		if public1 != public2 {
			fmt.Println("The content of the certificate request does not match the content of the private key")
			return errors.New("the content of the certificate request does not match the content of the private key")
		}

		resp, err := v2xca.CreateCertificate(&subCertReq.TbsCertList[0], subCert, subPrivKey)
		if err != nil {
			fmt.Println("Failed to create sub certificate")
			return err
		}

		if outPath == "" {
			v2xca.PrintCert(resp)
			return nil
		}
		if len(outPath) < 4 {
			fmt.Println("Fail to name: name is too short")
			return errors.New("fail to name: name is too short")
		}
		if first := outPath[len(outPath)-4:]; first != ".oer" {
			fmt.Println("Certificate Naming Format Error")
			return errors.New("certificate Naming Format Error")
		}

		errs := v2xca.WriteCert(resp, outPath)
		if errs != nil {
			fmt.Println("Output Certificate Failed")
			return errs
		}
		fmt.Println("Create sub certificate successfully")
		return nil
	}

	if !existSub && !existSelf && existPrint {
		if outPath != "" {
			fmt.Println("Wrong option, See 'v2xca ca -help'")
			return errors.New("wrong option, See 'v2xca ca -help'")
		}

		if inPath == "" {
			fmt.Println("The specify certificate file is empty")
			return errors.New("the specify certificate file is empty")
		}

		cer, err := v2xca.ReadCert(inPath)
		if err != nil {
			fmt.Println("failed to read certificate file")
			return err
		}
		v2xca.PrintCert(cer)
		return nil
	}
	fmt.Println("Error parameter, Missing option, See 'v2xca ca -help'")
	return errors.New("error parameter")
}
