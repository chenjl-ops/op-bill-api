package prediction

import (
	"github.com/gin-gonic/gin"
)

func Url(r *gin.Engine) {
	v1 := r.Group("/prediction/v1")
	{
		v1.GET("create_table", createTable)
		v1.GET("insert_baidu_bill_data", insertBaiduBillData)
	}
}
