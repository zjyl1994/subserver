package main

import (
	"encoding/json"
	"io/ioutil"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"github.com/xeipuuv/gojsonschema"
)

var (
	subCache       *cache.Cache
	contentTypeMap = map[string]string{
		"plain":  "text/plain",
		"base64": "text/plain",
		"sip008": "application/json",
		"clash":  "text/vnd.yaml",
	}
	subTypes     = []string{"plain", "base64", "sip008", "clash"}
	updateSchema = `{
		"$schema": "http://json-schema.org/draft-04/schema#",
		"type": "object",
		"properties": {
		  "name": {
			"type": "string"
		  },
		  "sub_option": {
			"type": "object",
			"properties": {
			  "path": {
				"type": "array",
				"items": [
				  {
					"type": "string"
				  }
				]
			  }
			},
			"required": [
			  "path"
			]
		  },
		  "update_option": {
			"type": "object",
			"properties": {
			  "path": {
				"type": "string"
			  }
			},
			"required": [
			  "path"
			]
		  },
		  "shadowsocks": {
			"type": "array",
			"items": [
			  {
				"type": "object",
				"properties": {
				  "name": {
					"type": "string"
				  },
				  "server": {
					"type": "string"
				  },
				  "port": {
					"type": "integer"
				  },
				  "password": {
					"type": "string"
				  },
				  "method": {
					"type": "string"
				  },
				  "plugin": {
					"type": "string"
				  },
				  "plugin_opts": {
					"type": "string"
				  }
				},
				"required": [
				  "name",
				  "server",
				  "port",
				  "password",
				  "method",
				  "plugin",
				  "plugin_opts"
				]
			  }
			]
		  }
		},
		"required": [
		  "name",
		  "sub_option",
		  "update_option",
		  "shadowsocks"
		]
	  }`
)

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
	result, err := gojsonschema.Validate(gojsonschema.NewStringLoader(updateSchema), gojsonschema.NewStringLoader(string(newSubJSON)))
	if err != nil {
		_ = ctx.AbortWithError(500, err)
		return
	}

	if !result.Valid() {
		jsonErrs := make([]string, len(result.Errors()))
		for i, v := range result.Errors() {
			jsonErrs[i] = v.String()
		}
		ctx.String(400, strings.Join(jsonErrs, "\n"))
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
