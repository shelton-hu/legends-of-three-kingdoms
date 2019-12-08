package constants

const (
	hostPrefix         = "ltk-"
	ApiGatewayHost     = hostPrefix + "api-gateway"
	IAMManagerHost     = hostPrefix + "iam-manager"
	ProcessManagerHost = hostPrefix + "process-manager"
)

const (
	ApiGatewayPort     = 9000
	IAMManagerPort     = 9001
	ProcessManagerPort = 9002
)

const (
	IAMManagerName     = "iam-manager"
	ProcessManagerName = "process-manager"
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
