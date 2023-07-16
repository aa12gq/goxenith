// 数据(库)层常用的公共定义
package model

import (
	"gorm.io/plugin/soft_delete"

	"time"
)

type Model struct {
	Id        uint64                `gorm:"primarykey" json:"id"`
	CreatedAt *time.Time            `json:"created_at"`
	UpdatedAt *time.Time            `json:"updated_at"`
	DeletedAt *time.Time            `gorm:"default:null"`
	IsDel     soft_delete.DeletedAt `gorm:"column:_delete_;softDelete:flag,DeletedAtField:DeletedAt" json:"_delete_"`
}

// SoftDelete 自身表有id的可以内嵌该结构体实现软删除
type SoftDelete struct {
	DeletedAt *time.Time            `gorm:"default:null"`
	IsDel     soft_delete.DeletedAt `gorm:"column:_delete_;softDelete:flag,DeletedAtField:DeletedAt" json:"_delete_"` // 软删除
}

type BaseTime struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// EntMigrateOptions 使用entgo orm框架进行migrate的配置选项
type EntMigrateOptions struct {
	// 是否删除无用字段
	DropColumn bool
	// 是否删除无用索引
	DropIndex bool
	// 是否创建外键
	CreateForeignKey bool
	// 配置文件所在目录或路径
	ConfigPath string
	// 显示详细日志，如：打印sql日志等
	Verbose bool
	// 同Verbose
	Debug bool
	// 迁移执行超时时间，单位：秒。大于等于0的整数，等于0时，永不超时。
	Timeout uint
}
