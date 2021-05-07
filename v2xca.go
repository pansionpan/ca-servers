/*
main 包是命令行工具主入口，负责调用各个函数
命令行关键词：
1.sm2 	有关公私钥的操作函数
2.req	有关证书请求的操作函数
3.ca	有关根证书或中级证书的操作函数
*/
package main

import (
	"caserver-cmd/command"
	"os"

	"github.com/urfave/cli"
)

func main() {
	// 实例化命令行工具cli
	app := cli.NewApp()
	// 设定应用名字－－ca-server-cmd
	app.Name = "ca-server-cmd"
	// 设定版本号
	app.Version = "1.0.0"
	// 创建相关函数
	app.Commands = []cli.Command{
		{
			Name:  "sm2",
			Usage: "create keys or print keys",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "genkey",
					Usage: "generate private key",
				},
				cli.StringFlag{
					Name:  "out",
					Usage: "the file path",
				},
				cli.BoolFlag{
					Name:  "pubout",
					Usage: "create public key from private key",
				},
				cli.StringFlag{
					Name:  "in",
					Usage: "insert a private key",
				},
			},
			Action: command.CreateKeys,
		},
		{
			Name:  "req",
			Usage: "create cer request or print the specified cer request",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "new",
					Usage: "create new request",
				},
				cli.BoolFlag{
					Name:  "print",
					Usage: "print result to the console",
				},
				cli.StringFlag{
					Name:  "key",
					Usage: "specified private key",
				},
				cli.StringFlag{
					Name:  "out",
					Usage: "output file path",
				},
				cli.StringFlag{
					Name:  "in",
					Usage: "specified csr",
				},
			},
			Action: command.CreateCertRequestCommand,
		},
		{
			Name:  "ca",
			Usage: "create root certificate or print root certificate",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "selfsign",
					Usage: "create selfsign certificate",
				},
				cli.BoolFlag{
					Name:  "subca",
					Usage: "create sub certificate",
				},
				cli.BoolFlag{
					Name:  "print",
					Usage: "write certificate to console",
				},
				cli.StringFlag{
					Name:  "in",
					Usage: "specify certificate request",
				},
				cli.StringFlag{
					Name:  "signkey",
					Usage: "specify private key file",
				},
				cli.StringFlag{
					Name:  "out",
					Usage: "output file path",
				},
				cli.StringFlag{
					Name:  "cert",
					Usage: "specify issue cert",
				},
			},
			Action: command.CreateRootAndSubCert,
		},
	}
	// 接受os.Args启动程序(接受命令行参数，运行命令行启动)
	app.Run(os.Args)
}
