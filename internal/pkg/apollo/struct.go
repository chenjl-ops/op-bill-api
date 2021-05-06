package apollo

//apollo自身结构配置
type Apollo struct {
	AppID           string
	Cluster         string
	NameSpaceName   string
	ApolloServerUrl string
}

//apollo内容配置
type Specification struct {
	// ServerRunPort        string `envconfig:"SERVER_RUN_PORT" mapstructure:"server_run_port"`
	BaiduAccessKeyID     string `envconfig:"BAIDU_ACCESS_KEY_ID" mapstructure:"baidu_access_key_id"`
	BaiduAccessKeySecret string `envconfig:"BAIDU_ACCESS_KEY_SECRET" mapstructure:"baidu_access_key_secret"`
	InsertMysqlSum       int    `envconfig:"INSERT_MYSQL_SUM" mapstructure:"insert_mysql_sum"`
	MysqlUserName        string `envconfig:"MYSQL_USERNAME" mapstructure:"mysql_db_user"`
	MysqlPassword        string `envconfig:"MYSQL_PASSWORD" mapstructure:"mysql_db_passwd"`
	MysqlHost            string `envconfig:"MYSQL_HOST" mapstructure:"mysql_db_host"`
	MysqlPort            int    `envconfig:"MYSQL_PORT" mapstructure:"mysql_db_port"`
	MysqlDBName          string `envconfig:"MYSQL_DBNAME" mapstructure:"mysql_db_name"`
	RedisMasterName      string `envconfig:"REDIS_MASTER_NAME" mapstructure:"redis_cluster"`
	RedisSentinelAddress string `envconfig:"REDIS_SENTINEL_ADDRESS" mapstructure:"redis_sentinels"`
	RedisPasswd          string `envconfig:"REDIS_PASSWD" mapstructure:"redis_password"`
	CmdbAppUrl           string `envconfig:"CMDB_APP_URL" mapstructure:"cmdb_app_url"`
	CmdbEcsUrl           string `envconfig:"CMDB_ECS_URL" mapstructure:"cmdb_ecs_url"`
	CmdbVolumeUrl        string `envconfig:"CMDB_VOLUME_URL" mapstructure:"cmdb_volume_url"`
	CmdbAppInstanceUrl   string `envconfig:"CMDB_APP_INSTANCE_URL" mapstructure:"cmdb_app_instance_url"`
	AppInfoUrl           string `envconfig:"APP_INFO_URL" mapstructure:"app_info_url"`
	DependentName        string `envconfig:"DEPENDENT_NAME" mapstructure:"dependent_name"`
	ShareBillUrl         string `envconfig:"SHARE_BILL_URL" mapstructure:"share_bill_url"`
	SourceBillUrl        string `envconfig:"SOURCE_BILL_URL" mapstructure:"source_bill_url"`
	DownloadPath         string `envconfig:"DOWNLOAD_PATH" mapstructure:"download_path"`
}
