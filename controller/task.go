package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/daoone/hammer/logic"
)

func TaskList(c *gin.Context){
	unionId, err := c.Get("unionId")
	status, err2 := c.GetQuery("status")
	if !err || !err2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": true})
	}
	c.JSON(http.StatusOK,logic.TaskList(unionId.(string),status))
}
