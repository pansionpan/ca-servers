##ca-server-cmd命令行工具1.0.0版
##功能
* `sm2`
	1.生成`私钥`
	2.从指定`私钥`中取出对应的`公钥`
* `req`
	1.生成`证书请求`
	2.打印证书请求
* `ca`
	1.生成`自签证书`
	2.生成`中级证书`
	3.打印证书

##命令

* v2xca sm2 -genkey -out privkey.key                		 	                   生成 sm2 private key
* v2xca sm2 -pubout -in privkey.key -out pubkey.key		  		                   从 private key 中取出 public key
* v2xca req -new -key user.key -out user.csr                      		           生成证书请求
* v2xca req -print -in user.csr                                   		           创建公钥与私钥
* v2xca ca  -selfsign -in rootca.csr -signkey priv.key -out cert.oer               生成证书请求
* v2xca ca  -subca -in subca.csr -signkey priv.key -cert rootca.oer -out cert.oer  生成二级证书
* v2xca ca  -print -in cert.oer                                                    打印证书


##运行

1. v2xca
```golang
NAME:
   ca-server-cmd - A new cli application

USAGE:
   v2xca [global options] command [command options] [arguments...]

VERSION:
   1.0.0

COMMANDS:
   sm2  		    create Private key
   req  		    create cer request
   ca      		    create root certificate
   h                Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

2. v2xca sm2
```golang
NAME:
   v2xca sm2 - create private key or public key

USAGE:
   v2xca sm2 [command options] [arguments...]

OPTIONS:
   -genkey                generate private key
   -out value  		      the file path
   -pubout                create public key from private key
   -in  value             insert a private key
```

3. v2xca req
```golang
NAME:
   v2xca req - create cer request or print the specified cer request

USAGE:
   v2xca req [command options] [arguments...]

OPTIONS:
   -new                  create new request
   -print                print result to the console
   -key value            specified private key
   -out value            output file path
   -in  value  		     specified csr
```

4. v2xca ca
```golang
NAME:
   v2xca ca - create root certificate or print root certificate

USAGE:
   v2xca ca [command options] [arguments...]

OPTIONS:
   -selfsign            create selfsign certificate
   -subca               create sub certificate
   -print               write certificate to console
   -in      value       specify certificate request
   -signkey value       specify private key file
   -out     value       output file path
   -cert    value       specify issue cert
```


