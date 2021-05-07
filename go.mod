module caserver-cmd

go 1.12

require (
	github.com/buger/goterm v0.0.0-20181115115552-c206103e1f37
	github.com/liwangqiang/gmsm v1.2.1
	github.com/urfave/cli v1.21.0
	iauto.com/asn1c-oer-lib/v2xca v0.0.0-00010101000000-000000000000
)

replace iauto.com/asn1c-oer-lib/v2xca => ../asn1c-oer-lib/v2xca
