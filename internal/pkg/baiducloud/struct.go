package baiducloud

// BaiduCloud 百度云认证数据
type BaiduCloud struct {
	Url             string
	Headers         map[string]string
	Method          string
	AccessKeyID     string
	AccessKeySecret string
}
