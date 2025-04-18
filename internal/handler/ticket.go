package handler

/**
 * @File: ticket.go
 * @Description:
 *
 * @Author: Timmy
 * @Create: 2025/4/18 下午10:03
 * @Software: GoLand
 * @Version:  1.0
 */

import (
	"github.com/taguo1109/go-ticket-system/internal/kafkautil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/taguo1109/go-ticket-system/internal/model"
)

func BookTicketHandler(c *gin.Context) {
	var req model.TicketRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "資料格式錯誤"})
		return
	}

	if err := kafkautil.SendTicketRequest(req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "發送失敗"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "搶票請求送出成功"})
}
