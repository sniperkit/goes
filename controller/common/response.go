package common

import (
	// external
	"github.com/kataras/iris"

	// internal
	"github.com/sniperkit/goes/model"
)

func SendErrorJSON(message string, ctx iris.Context) {
	ctx.JSON(iris.Map{
		"errCode": model.ERROR,
		"message": message,
		"data":    iris.Map{},
	})
}
