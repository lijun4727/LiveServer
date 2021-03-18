package common

import "time"

const (
	TokenExpireTime            = 24 * time.Hour * 7
	AccountRpcPort             = ":50051"
	LoginRpcPort               = ":50052"
	ShopManagePort             = ":50053"
	ImageMangePort             = ":50054"
	MysqlConnCmd               = "lion:472780330@tcp(10.10.128.7:3306)/live_server?charset=utf8"
	AccountManageRpcAddress    = "127.0.0.1" + AccountRpcPort
	LoginRpcAddress            = "127.0.0.1" + LoginRpcPort
	ShopRpcAddress             = "127.0.0.1" + ShopManagePort	
	SignalPort                 = ":8082"
	MediaTransactionPort       = ":8083"
	ApiReadTimeOut             = 30 * time.Second
	ApiWriteTimeOut            = 30 * time.Second
	RpcTimeout                 = 20 * time.Second
	RedisSlaveAddress          = "127.0.0.1:6379"
	RedisSlavePassword         = "l"
	RedisMasterAddress         = "127.0.0.1:6379"
	RedisMasterPassword        = "l"
	RedisReadTimeOut           = 10 * time.Second
	RedisWriteTimeOut          = 10 * time.Second
	RedisAccountTimeOut        = time.Hour * 24 * 7
	ShopNameMaxChar            = 10
	ShopDescMaxChar            = 50
	RabbitmqHandleOrderAddress = "amqp://lijun:l@192.168.1.6:5672/"
	RabbitmqHandleOrderQuene   = "create_order"
	RabbitmqHandleOrderTimeout = time.Minute
	MongoConnCmd               = "mongodb://lion:123456@192.168.1.6:27017"
	AccountManageSvrName       = "account_manage"
)
