package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/api/util"
)

// @description Firestore テスト
// @version 1.0
// @Tags Firestore
// @Summary API確認
// @accept application/x-json-stream
// @Success 200 {object} gin.H {"status": "success"}
// @router /api/v1/firestore [get]
func Test(c *gin.Context) {
	err := util.AddFirestore("test", map[string]interface{}{
		"first": "ばかだねーーーーーーーお前",
		"last":  "Lovelace",
		"born":  1815,
	})

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error"})
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
