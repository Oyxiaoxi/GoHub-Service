//Package topic 模型
package topic

import (
    "GoHub-Service/app/models"
    "GoHub-Service/app/models/category"
    "GoHub-Service/app/models/user"
    "GoHub-Service/pkg/database"

    "github.com/spf13/cast"
)

type Topic struct {
    models.BaseModel

    Title      string `json:"title,omitempty"`
    Body       string `json:"body,omitempty"`
    UserID     string `gorm:"index" json:"user_id,omitempty"`
    CategoryID string `gorm:"index" json:"category_id,omitempty"`

    // 通过 user_id 关联用户
    User user.User `json:"user"`

    // 通过 category_id 关联分类
    Category category.Category `json:"category"`

    models.CommonTimestampsField
}

func (topic *Topic) Create() {
    database.DB.Create(&topic)
}

func (topic *Topic) Save() (rowsAffected int64) {
    result := database.DB.Save(&topic)
    return result.RowsAffected
}

func (topic *Topic) Delete() (rowsAffected int64) {
    result := database.DB.Delete(&topic)
    return result.RowsAffected
}

// GetID 实现Model接口
func (topic *Topic) GetID() uint64 {
    return topic.ID
}

// GetOwnerID 实现OwnershipChecker接口，返回资源所有者ID
func (topic *Topic) GetOwnerID() string {
    return topic.UserID
}

// GetStringID 获取字符串格式的ID
func (topic *Topic) GetStringID() string {
    return cast.ToString(topic.ID)
}
