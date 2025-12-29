//Package link 模型
package link

import (
    "GoHub-Service/app/models"
    "GoHub-Service/pkg/database"
    "GoHub-Service/pkg/cache"
    "GoHub-Service/pkg/helpers"
    
    "time"

)

type Link struct {
    models.BaseModel

    Name string `gorm:"index" json:"name,omitempty"`
    URL  string `gorm:"index" json:"url,omitempty"`

    models.CommonTimestampsField
}

func (link *Link) Create() {
    database.DB.Create(&link)
}

func (link *Link) Save() (rowsAffected int64) {
    result := database.DB.Save(&link)
    return result.RowsAffected
}

func (link *Link) Delete() (rowsAffected int64) {
    result := database.DB.Delete(&link)
    return result.RowsAffected
}

func AllCached() (links []Link) {
    // 设置缓存 key
    cacheKey := "links:all"
    // 设置过期时间
    expireTime := 120 * time.Minute
    // 取数据
    cache.GetObject(cacheKey, &links)

    // 如果数据为空
    if helpers.Empty(links) {
        // 查询数据库
        links = All()
        if helpers.Empty(links) {
            return links
        }
        // 设置缓存
        cache.Set(cacheKey, links, expireTime)
    }
    return
}
