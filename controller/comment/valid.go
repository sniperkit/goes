package comment

import (
	"strings"
	"unicode/utf8"

	// internal
	"github.com/sniperkit/snk.golang.vuejs-multi-backend/config"
	"github.com/sniperkit/snk.golang.vuejs-multi-backend/controller/common"
	"github.com/sniperkit/snk.golang.vuejs-multi-backend/model"

	// external
	"github.com/sniperkit/iris"
)

func commentValid(comment *model.Comment, ctx iris.Context) {
	if err := ctx.ReadJSON(comment); err != nil {
		common.SendErrorJSON("参数错误", ctx)
		return
	}

	if comment.ArticleID != 0 {
		var article model.Article
		if model.DB.First(&article, comment.ArticleID).RecordNotFound() {
			common.SendErrorJSON("无效的评论文章 ID", ctx)
			return
		}
	} else {
		common.SendErrorJSON("非法文章id", ctx)
		return
	}

	if comment.ParentID != 0 {
		var parentComment model.Comment
		if err := model.DB.First(&parentComment, comment.ParentID).Error; err != nil {
			common.SendErrorJSON("无效的评论id", ctx)
			return
		}
	} else {
		common.SendErrorJSON("非法评论id", ctx)
		return
	}

	comment.Content = strings.TrimSpace(comment.Content)

	if comment.Content == "" {
		common.SendErrorJSON("评论不能为空", ctx)
		return
	}

	if utf8.RuneCountInString(comment.Content) > config.ServerConfig.MaxCommentLength {
		common.SendErrorJSON("评论字数超过限制", ctx)
		return
	}
}
