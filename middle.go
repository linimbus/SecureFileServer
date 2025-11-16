package main

import (
	"encoding/json"
	"strings"

	"github.com/astaxie/beego/context"
)

func unauthorizedResponse(ctx *context.Context, message string) {
	ctx.Output.SetStatus(401)
	ctx.Output.Header("Content-Type", "application/json")

	response := map[string]interface{}{
		"error":   "Unauthorized",
		"message": message,
		"code":    401,
	}
	jsonData, _ := json.Marshal(response)
	ctx.Output.Body(jsonData)
	ctx.ResponseWriter.WriteHeader(401)
}

func TokenAuthFilter(ctx *context.Context) {
	if !AppConf.TokenEnable {
		return
	}

	authHeader := ctx.Request.Header.Get("Authorization")
	if authHeader == "" {
		unauthorizedResponse(ctx, "Missing Authorization header")
		return
	}

	if strings.ToLower(authHeader) != AppConf.TokenSecret {
		unauthorizedResponse(ctx, "Authorization header does not match token secret")
		return
	}
}
