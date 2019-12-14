package constants

const (
	hostPrefix      = "ltk-"
	ApiGatewayHost  = hostPrefix + "api-gateway"
	IAMManagerHost  = hostPrefix + "iam-manager"
	RoomManagerHost = hostPrefix + "room-manager"
	GameManagerHost = hostPrefix + "game-manager"
)

const (
	ApiGatewayPort  = 9000
	IAMManagerPort  = 9001
	RoomManagerPort = 9002
	GameManagerPort = 9003
)

const (
	IAMManagerName  = "iam-manager"
	RoomManagerName = "room-manager"
	GameManagerName = "game-manager"
)

const (
	EnvLoaderPrefix  = "ltk"
	ConfigEtcdPrefix = "/ltk/"
)

const (
	GlobalConfigKey = "global_config"
	DlockKey        = "dlock_" + GlobalConfigKey
)

// 在 pkg/pi/pi.go 中被使用
const (
	MysqlPrismaEndpointSuffix = "/ltk/mysql"
)
