package about

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tiktok/src/constant/strings"
	"tiktok/src/gateway/models"
	"time"
)

func Handle(c *gin.Context) {
	var req models.AboutReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status_code": strings.GateWayErrorCode,
			"status_msg":  strings.GateWayError,
		})
	}
	res := models.AboutRes{
		Echo:      req.Echo,
		TimeStamp: time.Now().Unix(),
	}
	c.JSON(http.StatusOK, res)
}
