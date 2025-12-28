// Package services 友情链接业务逻辑服务
package services

import (
	"GoHub-Service/app/models/link"
	"time"
)

// LinkService 友情链接服务
type LinkService struct{}

// NewLinkService 创建友情链接服务实例
func NewLinkService() *LinkService {
	return &LinkService{}
}

// LinkResponseDTO 友情链接响应DTO
type LinkResponseDTO struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// LinkListResponseDTO 友情链接列表响应DTO
type LinkListResponseDTO struct {
	Links []LinkResponseDTO `json:"links"`
}

// toResponseDTO 将Link模型转换为响应DTO
func (s *LinkService) toResponseDTO(l *link.Link) *LinkResponseDTO {
	return &LinkResponseDTO{
		ID:        l.GetStringID(),
		Name:      l.Name,
		URL:       l.URL,
		CreatedAt: l.CreatedAt,
		UpdatedAt: l.UpdatedAt,
	}
}

// toResponseDTOList 将Link模型列表转换为响应DTO列表
func (s *LinkService) toResponseDTOList(links []link.Link) []LinkResponseDTO {
	dtos := make([]LinkResponseDTO, len(links))
	for i, l := range links {
		dtos[i] = LinkResponseDTO{
			ID:        l.GetStringID(),
			Name:      l.Name,
			URL:       l.URL,
			CreatedAt: l.CreatedAt,
			UpdatedAt: l.UpdatedAt,
		}
	}
	return dtos
}

// GetAllCached 获取所有缓存的友情链接
func (s *LinkService) GetAllCached() *LinkListResponseDTO {
	links := link.AllCached()
	return &LinkListResponseDTO{
		Links: s.toResponseDTOList(links),
	}
}
