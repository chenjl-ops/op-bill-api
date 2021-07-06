package billing

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/shopspring/decimal"
	"github.com/thoas/go-funk"
	"op-bill-api/internal/app/middleware/logger"
	"op-bill-api/internal/pkg/config"
	"op-bill-api/internal/pkg/mysql"
	"op-bill-api/internal/pkg/requests"
	"strconv"
	"strings"
	"time"
)

var DpUser map[string]map[string]string

//DpUser := make(map[string]map[string]string, 0)

// GetMonthFirstDate 获取月第一天和最后一天
// 获取月第一天 shift为偏移量 正为下月 负为上月
func GetMonthFirstDate(t time.Time, shift int) string {
	t = t.AddDate(0, shift, -t.Day()+1)
	return t.Format(timeFormat)

}

// GetMonthLastDate 获取月最后一天 shift为偏移量 正为下月 负为上月
func GetMonthLastDate(t time.Time, shift int) string {
	t = t.AddDate(0, shift, -t.Day())
	return t.Format(timeFormat)
}

// GetMonthDate 获取月第一天和最后一天
func GetMonthDate() map[string]string {
	data := make(map[string]string)

	t := time.Now()
	data["thisMonthFirstDate"] = GetMonthFirstDate(t, 0)  // 本月第一天
	data["thisMonthLastDate"] = GetMonthLastDate(t, 1)    // 本月最后一天
	data["lastMonthFirstDate"] = GetMonthFirstDate(t, -1) //上月第一天
	data["lastMonthLastDate"] = GetMonthLastDate(t, 0)    //上月最后一天

	return data
}

// GetTexData 获取折扣率数据
func GetTexData(name string) (SourceBillTex, error) {
	var data SourceBillTex
	_, err := mysql.Engine.Where("name = ?", name).Get(&data)
	return data, err
}

// InsertOrUpdateTexData 更新折扣率，通过name和tex action标识动作
func InsertOrUpdateTexData(name string, tex float64, action string) bool {
	data := SourceBillTex{Name: name, Tex: tex}
	var err error
	if action == "POST" { // 支持新增和update减少客户端逻辑
		has, _ := mysql.Engine.Exist(&SourceBillTex{Name: name})
		if has {
			_, err = mysql.Engine.Update(&data, &SourceBillTex{Name: name})
		} else {
			_, err = mysql.Engine.Insert(&data)
		}
	} else if action == "PUT" {
		_, err = mysql.Engine.Update(&data, &SourceBillTex{Name: name})

	} else {
		return false
	}

	// 判断操作结果
	if err != nil {
		return false
	} else {
		return true
	}
}

// CheckBillData 校验账单数据
func CheckBillData(month string, isShare bool) bool {
	has, _ := mysql.Engine.Exist(&BillData{
		Month:   month,
		IsShare: isShare,
	})
	return has
}

// GetBillData 查询账单数据
func GetBillData(month string, isShare bool) (BillData, error) {
	var data BillData
	_, err := mysql.Engine.Where("month = ?", month).And("isShare = ?", isShare).Get(&data)
	return data, err
}

// InsertBillData 数据写入
func InsertBillData(month string, isShare bool, data map[string]map[string]float64) error {
	x := BillData{
		Month:   month,
		IsShare: isShare,
		Data:    data,
	}

	_, err := mysql.Engine.Insert(&x)
	return err
}

// GetAllBillData 获取账单全量数据 资金/损益 口径
// TODO 分页功能
func GetAllBillData(isShare bool) ([]BillData, error) {
	data := make([]BillData, 0)
	err := mysql.Engine.Where("isShare = ?", isShare).Find(&data)
	if err != nil {
		logrus.Println("获取账单全量数据异常: ", err)
	}
	return data, err
}

// GetBaiduShareBillData 获取损益账单数据详情
func GetBaiduShareBillData(month string) ([]config.ShareBill, error) {
	var data []config.ShareBill
	err := mysql.Engine.Where("month = ?", month).Find(&data)
	return data, err
}

// GetBaiduSourceBillData 获取资金口径账单数据详情
func GetBaiduSourceBillData(month string) ([]config.SourceBill, error) {
	var data []config.SourceBill
	err := mysql.Engine.Where("month = ?", month).Find(&data)
	return data, err
}

// 临时需求
// 计算部门分摊花费
func getBaiduSourceBillData(month string, sellType string) ([]config.SourceBill, error) {
	sellTypes := map[string]string{"prepay": "预付费", "postpay": "后付费"}
	var data []config.SourceBill
	err := mysql.Engine.Where("month = ?", month).Where("sell_type = ?", sellTypes[sellType]).Find(&data)
	return data, err
}

// 请求CMDB接口获取bcc数据
func getBcc(id string) map[string]interface{} {
	url := fmt.Sprintf("http://op-do-cmdb-api.prod-devops.k8s.chj.cloud/op-do-cmdb-api/v1-0/commondata/all/md_cloud_virtual_machine?cloud_id=%s", id)
	var data []map[string]interface{}
	err := requests.Request(url, &data)
	if err != nil {
		logrus.Println(err)
		return nil
	} else {
		return data[0]
	}
}

func getAllBcc() []interface{} {
	//url := fmt.Sprintf("http://op-do-cmdb-api.prod-devops.k8s.chj.cloud/op-do-cmdb-api/v1-0/commondata/md_cloud_virtual_machine?cloud_id=%s", id)
	url := fmt.Sprintf("http://op-do-cmdb-api.prod-devops.k8s.chj.cloud/op-do-cmdb-api/v1-0/commondata/md_cloud_virtual_machine?size=2000")
	data := make(map[string]interface{})
	result := make(map[string]interface{})
	headers := map[string]string{"Content-Type": "application/json", "authorization": "iLimzW8YM/MKOWXg1EMt7N5xTIfo4vfnKTNWqDgmrgk="}
	err := requests.RequestMethod(url, "GET", headers, data, &result)
	if err != nil {
		fmt.Println(err)
	}
	return result["data"].([]interface{})
}

func getVolume(id string) map[string]interface{} {
	url := fmt.Sprintf("http://op-do-cmdb-api.prod-devops.k8s.chj.cloud/op-do-cmdb-api/v1-0/commondata/all/md_disk?cloud_id=%s", id)
	var data []map[string]interface{}
	err := requests.Request(url, &data)
	if err != nil {
		logrus.Println(err)
		return nil
	} else {
		return data[0]
	}
}

func getUserDps(username string) map[string]string {
	url := fmt.Sprintf("http://op-do-cmdb-api.prod-devops.k8s.chj.cloud/op-do-cmdb-api/v1-0/commondata/all/user_info?username=%s", username)
	var data []map[string]string
	err := requests.Request(url, &data)
	if err != nil {
		logrus.Println(err)
		return nil
	} else {
		return data[0]
	}
}

func getK8sId() map[string]string {
	url := fmt.Sprintf("http://op-do-cmdb-api.dev.k8s.chj.cloud/op-do-cmdb-api/v1-0/commondata/all/md_k8s_cluster")
	var data []map[string]string
	err := requests.Request(url, &data)
	if err != nil {
		logrus.Println(err)
	}

	// {"prod_devops"}
	ids := make(map[string]string, 0)
	for _, v := range data {
		ids[fmt.Sprintf("%s_%s", v["env"], v["cluster"])] = v["_id"]
	}
	return ids
}

func getK8sNodes(id string) []string {
	url := fmt.Sprintf("http://op-do-cmdb-api.dev.k8s.chj.cloud/op-do-cmdb-api/v1-0/commondata/all/md_k8s_node?belongs_to_k8s=%s", id)
	var data []map[string]string
	err := requests.Request(url, &data)
	if err != nil {
		logrus.Println(err)
	}

	responseData := make([]string, 0)
	for _, v := range data {
		if v["provider_id"] != "" {
			responseData = append(responseData, strings.Replace(v["provider_id"], "cce://", "", 1))
		}
	}
	return responseData
}

// 获取k8s集群对应机器 百度CCE优先
func getK8sClusterInfo() map[string][]string {
	ids := getK8sId()

	data := make(map[string][]string)
	for k, v := range ids {
		data[k] = getK8sNodes(v)
	}

	return data
}

func getBccDp(id string, ch chan<- map[string]map[string]string) {
	data := make(map[string]map[string]string)
	bccData := getBcc(id)
	//dp := getUserDps(bccData["owner"].(string))
	if funk.Contains(funk.Keys(DpUser), bccData["owner"].(string)) {
		data[id] = DpUser[bccData["owner"].(string)]
		ch <- data
	}
}

//func getBccVo(id string, ch chan<- map[string][]string) {
//	bccData := getAllBcc()
//
//	var disks []string
//	if bccData["data_disks"] != nil {
//		for _, v := range bccData["data_disks"].([]interface{}) {
//			disks = append(disks, v.(map[string]interface{})["cloud_id"].(string))
//		}
//	}
//	for _, v := range bccData["system_disk"].([]interface{}) {
//		disks = append(disks, v.(map[string]interface{})["cloud_id"].(string))
//	}
//
//	data := make(map[string][]string)
//	data[bccData["cloud_id"].(string)] = disks
//
//	ch <- data
//}

func getVoDp(id string, ch chan<- map[string]map[string]string) {
	data := make(map[string]map[string]string, 0)
	voData := getVolume(id)
	if voData != nil {
		//dp := getUserDps(voData["owner"].(string))
		data[id] = DpUser[voData["owner"].(string)]
		ch <- data
	}
}

func getDepartmentBill(month string) map[string]map[string]float64 {
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
	//
	var apiNameBillNameMap = map[string]string{
		"DDoS高防服务":        "DDoS高防IP ADAS",
		"EIP带宽包":          "EIP带宽包 EIP_BP",
		"IPv6公网网关":        "IPv6公网网关",
		"NAT网关":           "NAT网关",
		"SSL证书服务":         "SSL证书",
		"云服务器":            "云服务器 BCC",
		"云磁盘":             "云磁盘 CDS",
		"代理关系型数据库":        "云数据库RDS RDS_PROXY",
		"关系型数据库":          "云数据库RDS",
		"内容分发网络":          "内容分发网络 CDN",
		"函数计算":            "函数计算 CFC",
		"只读关系型数据库":        "云数据库RDS RDS_REPLICA",
		"密钥管理服务":          "密钥管理服务 KMS",
		"对等连接":            "对等连接 PEERCONN",
		"对象存储":            "对象存储 BOS",
		"弹性公网IP":          "弹性公网IP EIP",
		"文件存储":            "文件存储 CFS",
		"服务网卡":            "服务网卡 SNIC",
		"本地DNS服务":         "本地DNS服务 LD",
		"海外CDN":           "海外CDN CDN_ABOAD",
		"物理服务器":           "弹性裸金属服务器 BBC",
		"百度Elasticsearch": "Elasticsearch",
		"百度日志服务":          "日志服务 BLS",
		"移动解析":            "移动域名解析 HTTPDNS",
		"简单缓存服务":          "云数据库SCS for Redis",
		"负载均衡":            "负载均衡 BLB",
		"音视频点播":           "视频创作分发平台",
		"音视频直播":           "音视频直播 LSS",
		"音视频转码":           "音视频处理 MCP",
	}

	// 获取所有部门
	var allDps []map[string]string
	err := requests.Request("http://op-do-cmdb-api.prod-devops.k8s.chj.cloud/op-do-cmdb-api/v1-0/commondata/all/user_info", &allDps)
	if err != nil {
		logrus.Println("获取部门信息异常: ", err)
	}

	// 初始化部门
	// {"dpNamePath": dpName}
	DpUser = make(map[string]map[string]string, 0)
	dps := make(map[string]string, 0)
	for _, v := range allDps {
		//logrus.Println("部门信息: ", v)
		DpUser[v["username"]] = v
		if !funk.Contains(funk.Keys(dps), v["departmentNamePath"]) {
			dps[v["departmentNamePath"]] = v["departmentName"]
		}
	}
	//logrus.Println("部门: ", dps)

	// 获取bcc 对应磁盘
	bccVos := make(map[string][]string, 0)
	bccAllData := getAllBcc()

	for _, v := range bccAllData {
		tempX := make([]string, 0)

		if v.(map[string]interface{})["data_disks"] != nil {
			for _, vv := range v.(map[string]interface{})["data_disks"].([]interface{}) {
				tempX = append(tempX, vv.(map[string]interface{})["cloud_id"].(string))
			}
			//tempX = append(tempX, v.(map[string]interface{})["data_disks"].([]string)...)
		}
		for _, vv := range v.(map[string]interface{})["system_disk"].([]interface{}) {
			tempX = append(tempX, vv.(map[string]interface{})["cloud_id"].(string))
		}
		//tempX = append(tempX, v.(map[string]interface{})["system_disk"].([]string)...)
		bccVos[v.(map[string]interface{})["cloud_id"].(string)] = tempX
	}

	// 获取所有Bcc ID
	var bccIds []string
	bccDPS := make(map[string]map[string]string, 0)
	var bccData []map[string]interface{}
	err = requests.Request("http://op-do-cmdb-api.prod-devops.k8s.chj.cloud/op-do-cmdb-api/v1-0/commondata/all/md_cloud_virtual_machine", &bccData)
	if err != nil {
		logrus.Println("获取Bcc异常: ", err)
	}
	for _, v := range bccData {
		bccIds = append(bccIds, v["cloud_id"].(string))
		bccDPS[v["cloud_id"].(string)] = DpUser[v["owner"].(string)]
	}
	//logrus.Println("DPS: ", bccDPS)
	//logrus.Println("bcc 总长度: ", len(bccIds), bccIds)

	//var wg sync.WaitGroup
	//
	//bccChan := make(chan map[string]map[string]string)
	////bccVoChan := make(chan map[string][]string)
	//// 并发请求相关bcc对应部门关系
	//
	////for i := 0; i < 2; i++ {
	////	wg.Add(1)
	////	go func() {
	////		defer wg.Done()
	//for _, j := range bccIds {
	//	//logrus.Println("bcc id:", j)
	//	go getBccDp(j, bccChan)
	//	//go getBccVo(j, bccVoChan)
	//}
	//	}()
	//}
	//wg.Wait()

	//for _, i := range bccIds {
	//	go getBccDp(i, bccChan)
	//	go getBccVo(i, bccVoChan)
	//}

	// 获取所有磁盘ID
	var voIds []string
	voDPS := make(map[string]map[string]string, 0)
	var voData []map[string]interface{}
	err = requests.Request("http://op-do-cmdb-api.prod-devops.k8s.chj.cloud/op-do-cmdb-api/v1-0/commondata/all/md_disk", &voData)
	if err != nil {
		logrus.Println("获取Volume异常: ", err)
	}
	for _, v := range voData {
		voIds = append(voIds, v["cloud_id"].(string))
		voDPS[v["cloud_id"].(string)] = DpUser[v["owner"].(string)]
	}
	//voChan := make(chan map[string]map[string]string)

	//for j := 0; j < 2; j++ {
	//	wg.Add(1)
	//	go func() {
	//		defer wg.Done()
	//for _, i := range voIds {
	//	go getVoDp(i, voChan)
	//}
	//	}()
	//}
	// 初始化bcc对应关系，cds对应关系. bcc对应部门，bcc对应硬盘，硬盘对应部门
	//bccDps := make(map[string]map[string]string, 0)
	//
	//voDps := make(map[string]map[string]string, 0)
	//for _, i := range bccIds {
	//	data := <-bccChan
	//	bccDps[i] = data[i]
	//
	//	//voData := <-bccVoChan
	//	//bccVos[i] = voData[i]
	//}
	//
	//for _, i := range voIds {
	//	data := <-voChan
	//	voDps[i] = data[i]
	//}
	//
	//var billData map[string]map[string]float64
	//for _, v := range shareBillProductNames {
	//	billData[v] = make(map[string]float64, 0)
	//}

	// 初始化部门花费金额
	billData := make(map[string]float64, 0)
	for k, _ := range dps {
		//logrus.Println("部门信息: ", k)
		billData[k] = 0.00
	}
	//logrus.Println("初始化dp: ", billData)

	// 获取损益账单数据
	shareBillData := make([]config.ShareBill, 0)
	if err := mysql.Engine.Where("month = ?", month).Find(&shareBillData); err != nil {
		logger.Log.Error("查询损益数据异常: ", err)
	}
	// 获取资金账单数据 后付费
	sourceData := make([]config.SourceBill, 0)
	if err := mysql.Engine.Where("month = ?", "2021-06-01_2021-06-30").Where("sell_type = ?", "后付费").Find(&sourceData); err != nil {
		logger.Log.Error("查询资金数据异常: ", err)
	}

	// 获取k8s集群对应bcc相关信息
	k8sData := getK8sClusterInfo()
	//logrus.Println("k8sData: ", k8sData)

	// 初始k8s集群花费
	k8sBillData := make(map[string]float64, 0)
	for k, _ := range k8sData {
		k8sBillData[k] = 0.00
	}

	// 初始化其他花费
	otherBillData := make(map[string]float64, 0)
	for _, v := range shareBillProductNames {
		otherBillData[v] = 0.00
	}

	// 计算损益账单到部门
	for _, v := range shareBillData {
		x, _ := strconv.ParseFloat(v.ShareCope, 64)

		if v.ProductName == "云服务器 BCC" {
			isOk := false
			for z, y := range k8sData {
				if funk.Contains(y, v.AssetId) {
					k8sBillData[z] = k8sBillData[z] + x
					isOk = true
				}
			}
			if !isOk {
				billData[bccDPS[v.AssetId]["departmentNamePath"]] = billData[bccDPS[v.AssetId]["departmentNamePath"]] + x // 每个部门bcc花费
				//if bccDPS[v.AssetId]["departmentNamePath"] != "" {
				//	billData[bccDPS[v.AssetId]["departmentNamePath"]] = billData[bccDPS[v.AssetId]["departmentNamePath"]] + x // 每个部门bcc花费
				//} else {
				//	logrus.Println("空部门: ", v.AssetId)
				//}
			}
		} else if v.ProductName == "云磁盘 CDS" {
			isOk := false
			for z, y := range k8sData {
				for _, i := range y {
					if funk.Contains(bccVos[i], v.AssetId) {
						k8sBillData[z] = k8sBillData[z] + x
						isOk = true
					}
				}
			}
			if !isOk {
				billData[voDPS[v.AssetId]["departmentNamePath"]] = billData[voDPS[v.AssetId]["departmentNamePath"]] + x // 每个部门cds花费
				//if voDPS[v.AssetId]["departmentNamePath"] != "" {
				//	billData[voDPS[v.AssetId]["departmentNamePath"]] = billData[voDPS[v.AssetId]["departmentNamePath"]] + x // 每个部门cds花费
				//} else {
				//	logrus.Println("空部门: ", v.AssetId)
				//}
			}
		} else {
			otherBillData[v.ProductName] = otherBillData[v.ProductName] + x
		}
	}

	sourceBillData := make(map[string]float64, 0)

	for _, v := range apiNameBillNameMap{
		sourceBillData[v] = 0.00
	}

	for _, v := range sourceData {
		x, _ := strconv.ParseFloat(v.OrderCost, 64)
		texData, err := getTex(apiNameBillNameMap[v.ProductName])
		vv := 0.00
		if err != nil {
			vv, _ = decimal.NewFromFloat(x * texData.Tex).Round(2).Float64()
		} else {
			// 默认3.9折
			vv, _ = decimal.NewFromFloat(x * 0.39).Round(2).Float64()
		}

		if apiNameBillNameMap[v.ProductName] == "云服务器 BCC" {
			isOk := false
			for z, y := range k8sData {
				if funk.Contains(y, v.AssetId) {
					k8sBillData[z] = k8sBillData[z] + vv
					isOk = true
				}
			}
			if !isOk {
				billData[bccDPS[v.AssetId]["departmentNamePath"]] = billData[bccDPS[v.AssetId]["departmentNamePath"]] + vv // 每个部门bcc花费
			}
		} else if apiNameBillNameMap[v.ProductName] == "云磁盘 CDS" {
			isOk := false
			for z, y := range k8sData {
				for _, i := range y {
					if funk.Contains(bccVos[i], v.AssetId) {
						k8sBillData[z] = k8sBillData[z] + vv
						isOk = true
					}
				}
			}
			if !isOk {
				billData[voDPS[v.AssetId]["departmentNamePath"]] = billData[voDPS[v.AssetId]["departmentNamePath"]] + vv // 每个部门cds花费
			}
		}
		sourceBillData[v.ProductName] = sourceBillData[v.ProductName] + vv
	}

	data := map[string]map[string]float64{"dp": billData, "k8s": k8sBillData, "other": otherBillData, "sourceBill": sourceBillData}
	return data

}

func getTex(name string) (SourceBillTex, error) {
	var texData SourceBillTex
	_, err := mysql.Engine.Where("name = ?", name).Get(&texData)
	return texData, err
}
