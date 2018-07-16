package common

import (
	// external
	"github.com/kataras/iris"

	// internal
	"github.com/sniperkit/snk.golang.vuejs-multi-backend/model"
)

func SendErrorJSON(message string, ctx iris.Context) {
	ctx.JSON(iris.Map{
		"errCode": model.ERROR,
		"message": message,
		"data":    iris.Map{},
	})
}
