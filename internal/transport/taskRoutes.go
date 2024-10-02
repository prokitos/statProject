package transport

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func getTask(c *gin.Context) {

	c.JSON(http.StatusOK, "okey")
}

func insertTask(c *gin.Context) {
	time.Sleep(time.Second * 1)
	c.JSON(http.StatusOK, "okey")
}

func deleteTask(c *gin.Context) {

	c.JSON(http.StatusOK, "okey")
}

func updateTask(c *gin.Context) {

	c.JSON(http.StatusOK, "okey")
}
