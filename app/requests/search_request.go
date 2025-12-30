package requests

import (
    "github.com/gin-gonic/gin"
    "github.com/thedevsaddam/govalidator"
)

// SearchRequest 通用搜索请求
 type SearchRequest struct {
    Keyword string `json:"keyword,omitempty" valid:"keyword"`
}

// SearchValidation 搜索关键词验证
func SearchValidation(data interface{}, c *gin.Context) map[string][]string {
    rules := govalidator.MapData{
        "keyword": []string{"required", "min_cn:1", "max_cn:50"},
    }
    messages := govalidator.MapData{
        "keyword": []string{
            "required:关键词必填",
            "min_cn:关键词不能为空",
            "max_cn:关键词不能超过50字",
        },
    }
    return validate(data, rules, messages)
}
