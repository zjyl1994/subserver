package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	err := LoadDataSource("data.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	// types := []string{"plain", "base64", "sip008", "clash"}
	// for _, v := range types {
	// 	fmt.Printf("==== MODE:%s ====\n%s\n", v, GenerateSubscription(v))
	// }
	r := gin.Default()
	r.GET("/*urlPath", SubHandler)
	r.POST("/*urlPath", UpdateHandler)
	err = r.Run()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
}
