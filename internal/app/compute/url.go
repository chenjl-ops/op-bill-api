package compute

import "github.com/gin-gonic/gin"

func Url(r *gin.Engine) {
	v1 := r.Group("/bill/v1")
	{
		v1.GET("get_bill_data", GetBilling)
		v1.GET("get_budget_data", GetBudget)
		v1.GET("test_baidu_bill", TestBaiduBill)
	}
}

