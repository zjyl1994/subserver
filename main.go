package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	err := LoadDataSource("/app/data/data.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	r := gin.Default()
	r.GET("/*urlPath", SubHandler)
	r.POST("/*urlPath", UpdateHandler)
	err = r.Run(":8080")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
}
