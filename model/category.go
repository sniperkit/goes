package model

import (
	"time"
)

// 文章分类
type Category struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `sql:"index" json:"deletedAt"`
	Name      string     `json:"name"`
	Sequence  int        `json:"sequence"`
	ParentID  int        `json:"parentId"`
	Status    int        `json:"status"`
}

// 开启
const CategoryStatusOpen = 1

// 关闭
const CategoryStatusClose = 2
