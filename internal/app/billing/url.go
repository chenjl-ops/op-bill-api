package billing

import (
	"github.com/gin-gonic/gin"
)

func Url(r *gin.Engine) {
	v1 := r.Group("/billing/v1")
	{
		v1.GET("create_table", createTable)
		v1.GET("insert_bill_data", insertData)
		v1.GET("get_month_data", getMonthData)
	}
}
