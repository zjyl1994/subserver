package main

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
)

var subCache *cache.Cache
var contentTypeMap = map[string]string{
	"plain":  "text/plain",
	"base64": "text/plain",
	"sip008": "application/json",
	"clash":  "text/vnd.yaml",
}
var subTypes = []string{"plain", "base64", "sip008", "clash"}

func init() {
	subCache = cache.New(5*time.Minute, 10*time.Minute)
}
func SubHandler(ctx *gin.Context) {
	urlpath := ctx.Param("urlPath")
	subtype := ctx.Query("type")
	if !stringContains(urlpath, datasource.SubOption.Path) ||
		!stringContains(subtype, subTypes) {
		ctx.AbortWithStatus(404)
		return
	}
	var subContent string
	if s, ok := subCache.Get(subtype); ok {
		subContent = s.(string)
	} else {
		subContent = GenerateSubscription(subtype)
		subCache.Set(subtype, subContent, cache.DefaultExpiration)
	}
	ctx.String(200, contentTypeMap[subtype], subContent)
}

func UpdateHandler(ctx *gin.Context) {
	urlpath := ctx.Param("urlPath")
	if urlpath != datasource.UpdateOption.Path {
		ctx.AbortWithStatus(404)
		return
	}
	newSubJSON, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		_ = ctx.AbortWithError(500, err)
		return
	}
	err = json.Unmarshal(newSubJSON, &datasource)
	if err != nil {
		_ = ctx.AbortWithError(500, err)
		return
	}
	err = ioutil.WriteFile("data.json", newSubJSON, 0644)
	if err != nil {
		_ = ctx.AbortWithError(500, err)
		return
	}
	ctx.String(200, "done")
}

func stringContains(s string, arr []string) bool {
	for _, a := range arr {
		if a == s {
			return true
		}
	}
	return false
}
