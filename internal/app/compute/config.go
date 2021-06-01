package compute

// source资源数据 product_name
var sourceBillProductNames = [...]string{
	"DDoS高防IP ADAS",
	"EIP带宽包 EIP_BP",
	"Elasticsearch",
	"IPv6公网网关",
	"NAT网关",
	"云数据库RDS",
	"云数据库RDS RDS_PROXY",
	"云数据库RDS RDS_REPLICA",
	"云数据库SCS for Redis",
	"内容分发网络 CDN",
	"函数计算 CFC",
	"密钥管理服务 KMS",
	"对等连接 PEERCONN",
	"对象存储 BOS",
	"弹性公网IP EIP",
	"弹性裸金属服务器 BBC",
	"文件存储 CFS",
	"日志服务 BLS",
	"服务网卡 SNIC",
	"本地DNS服务 LD",
	"海外CDN CDN_ABOAD",
	"移动域名解析 HTTPDNS",
	"视频创作分发平台",
	"负载均衡 BLB",
	"音视频直播 LSS",
}

// source账单折扣率
var sourceBillTex = map[string]float64 {
	"Elasticsearch": 0.5,
	"DDoS高防IP ADAS": 0.8,
	"负载均衡 BLB": 0.8,
	"对等连接 PEERCONN": 0.8,
	"NAT网关": 0.8,
	"内容分发网络 CDN": 1,
	"海外CDN CDN_ABOAD": 1,
	"视频创作分发平台": 1,
	"音视频直播 LSS": 1,
	"对象存储 BOS": 1,
}

// share资源数据 product_name
var shareBillProductNames = [...]string{
	"DDoS高防IP ADAS",
	"Elasticsearch",
	"NAT网关",
	"SSL证书",
	"云数据库RDS",
	"云数据库SCS for Redis",
	"弹性公网IP EIP",
	"弹性裸金属服务器 BBC",
	"数据可视化 Sugar",
	"移动域名解析 HTTPDNS",
}
