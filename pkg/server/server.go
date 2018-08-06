package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func StartServer() {
	fmt.Println("Server started")
	r := gin.Default()
	r.Run()
}
