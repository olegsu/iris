package server

import (
	"github.com/gin-gonic/gin"
	"github.com/olegsu/iris/pkg/logger"
)

func StartServer(logger logger.Logger) {
	logger.Debug("Server started")
	r := gin.Default()
	r.Run()
}
