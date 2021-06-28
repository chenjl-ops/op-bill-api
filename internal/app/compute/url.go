package compute

import "github.com/gin-gonic/gin"

func Url(r *gin.Engine) {
	v1 := r.Group("/bill/v1")
	{
		v1.GET("get_bill_data", getBilling)
		v1.GET("get_all_bill_data", getAllBilling)
		v1.GET("get_prediction_data", getPrediction)
		v1.GET("get_all_prediction_data", getAllPrediction)
	}
}
