module github.com/polynetwork/explorer

go 1.14

require (
	github.com/astaxie/beego v1.12.3
	github.com/cosmos/cosmos-sdk v0.39.1
	github.com/ethereum/go-ethereum v1.9.15
	github.com/go-redis/redis v6.14.2+incompatible
	github.com/go-sql-driver/mysql v1.6.0
	github.com/golang/mock v1.4.0 // indirect
	github.com/joeqian10/neo-gogogo v0.0.0
	github.com/ontio/ontology v1.11.0
	github.com/pkg/errors v0.9.1
	github.com/shopspring/decimal v1.2.0
	golang.org/x/crypto v0.0.0-20200728195943-123391ffb6de // indirect
	golang.org/x/net v0.0.0-20200707034311-ab3426394381
	google.golang.org/genproto v0.0.0-20200430143042-b979b6f78d84 // indirect
	gorm.io/driver/mysql v1.1.1
	gorm.io/gorm v1.21.11
)

replace (
	github.com/cosmos/cosmos-sdk => github.com/Switcheo/cosmos-sdk v0.39.2-0.20200814061308-474a0dbbe4ba
	github.com/ethereum/go-ethereum => github.com/ethereum/go-ethereum v1.9.13
	github.com/joeqian10/neo-gogogo => github.com/blockchain-develop/neo-gogogo v0.0.0-20200824102609-fddf06a45f66
)
