// Package services 友情链接业务逻辑服务
package services

import (
	"GoHub-Service/app/models/link"
)

// LinkService 友情链接服务
type LinkService struct{}

// NewLinkService 创建友情链接服务实例
func NewLinkService() *LinkService {
	return &LinkService{}
}

// GetAllCached 获取所有缓存的友情链接
func (s *LinkService) GetAllCached() []link.Link {
	return link.AllCached()
}
