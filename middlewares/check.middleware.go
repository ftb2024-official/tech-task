package middlewares

import (
	"strconv"

	util "tech_task/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type SimpleJSON struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
}

func CheckParamMiddleware(ctx *gin.Context) {
	id, ok := ctx.Params.Get("id")
	if !ok {
		ctx.JSON(400, gin.H{"error": "Parameter `id` is required"})
		ctx.Abort()
	}

	if util.IsValidUUID(id) {
		ctx.Next()
	} else {
		ctx.JSON(400, gin.H{"error": "Parameter `id` is not uuid"})
		ctx.Abort()
	}
}

func CheckNameMiddleware(ctx *gin.Context) {
	var simpleJSON SimpleJSON
	err := ctx.BindJSON(&simpleJSON)
	if err != nil {
		logrus.Printf("Unable to marshal request body due to (%v)", err)
	}

	if util.OnlyLetters(simpleJSON.Name) && util.OnlyLetters(simpleJSON.Surname) {
		ctx.Set("name", simpleJSON.Name)
		ctx.Set("surname", simpleJSON.Surname)
		ctx.Next()
	} else {
		ctx.JSON(400, gin.H{"error": "Invalid name / surname"})
		ctx.Abort()
	}
	ctx.Next()
}

func CheckQueryParamMiddleware(ctx *gin.Context) {
	if limitStr, ok := ctx.GetQuery("limit"); ok {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			ctx.JSON(400, gin.H{"error": "Invalid limit query parameter\n"})
			ctx.Set("limit", 1)
			ctx.Next()
		}
		if limit <= 0 {
			ctx.Set("limit", 1)
			ctx.Next()
		}
		ctx.Set("limit", limit)
		ctx.Next()
	}

	ctx.Set("limit", 1)
	ctx.Next()
}
