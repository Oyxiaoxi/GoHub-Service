//Package category 模型
package category

import (
    "GoHub-Service/app/models"
    "GoHub-Service/pkg/database"
)

type Category struct {
    models.BaseModel

    Name        string `gorm:"index" json:"name,omitempty"`
    Description string `json:"description,omitempty"`
    SortOrder   int    `gorm:"type:int;default:0;index;comment:排序顺序" json:"sort_order,omitempty"`

    models.CommonTimestampsField
}

func (category *Category) Create() {
    database.DB.Create(&category)
}

func (category *Category) Save() (rowsAffected int64) {
    result := database.DB.Save(&category)
    return result.RowsAffected
}

func (category *Category) Delete() (rowsAffected int64) {
    result := database.DB.Delete(&category)
    return result.RowsAffected
}
