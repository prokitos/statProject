package transport

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetTask(c *gin.Context) {

	c.JSON(http.StatusOK, "okey")
}

func InsertTask(c *gin.Context) {
	time.Sleep(time.Second * 1)
	c.JSON(http.StatusOK, "okey")
}

func DeleteTask(c *gin.Context) {

	c.JSON(http.StatusOK, "okey")
}

func UpdateTask(c *gin.Context) {

	c.JSON(http.StatusOK, "okey")
}
