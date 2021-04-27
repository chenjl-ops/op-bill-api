package billing

import (
	"github.com/gin-gonic/gin"
	"op-bill-api/internal/app/compute"
)

func Url(r *gin.Engine) {
	v1 := r.Group("/billing/v1")
	{
		v1.GET("create_table", createTable)
		v1.GET("insert_bill_data", insertData)
		v1.GET("get_month_data", getMonthData)
		v1.GET("get_bill_data", compute.GetBilling)
	}
}
