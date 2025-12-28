// Package models 模型通用属性和方法
//
// 本包提供了所有模型的基础结构和通用方法
//
// 使用说明：
//
// 1. BaseModel：所有模型都应嵌入此结构，提供统一的ID字段
//
// 2. CommonTimestampsField：提供创建时间和更新时间字段
//
// 3. 模型定义示例：
//    type User struct {
//        models.BaseModel
//        Name  string `gorm:"column:name;type:varchar(255)" json:"name"`
//        Email string `gorm:"column:email;type:varchar(255);unique" json:"email"`
//        models.CommonTimestampsField
//    }
//
// 4. 模型方法建议：
//    - 基础CRUD：Create(), Save(), Delete() 在模型文件中定义
//    - 查询方法：Get(), GetBy(), All(), Paginate() 在util文件中定义
//    - 钩子方法：BeforeCreate(), AfterCreate() 等在hooks文件中定义
package models

import (
    "time"
    "github.com/spf13/cast"
)

// BaseModel 模型基类，提供统一的ID字段
// 所有业务模型都应该嵌入此结构
type BaseModel struct {
    ID uint64 `gorm:"column:id;primaryKey;autoIncrement;" json:"id,omitempty"`
}

// CommonTimestampsField 时间戳字段
// GORM会自动管理这两个字段的值
type CommonTimestampsField struct {
    CreatedAt time.Time `gorm:"column:created_at;index;" json:"created_at,omitempty"`
    UpdatedAt time.Time `gorm:"column:updated_at;index;" json:"updated_at,omitempty"`
}

// GetStringID 获取 ID 的字符串格式
// 用于需要字符串类型ID的场景，如JWT token生成
func (a BaseModel) GetStringID() string {
    return cast.ToString(a.ID)
}
