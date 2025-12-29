// Package services 友情链接业务逻辑服务
package services

import (
	"GoHub-Service/app/models/link"
	apperrors "GoHub-Service/pkg/errors"
	"time"
)

// LinkService 提供友情链接查询的只读服务，主要走缓存避免频繁 DB 访问.
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

// GetAllCached 仅从缓存拉取链接列表，缓存缺失时返回 NotFound 供上层决定回源或直接响应.
func (s *LinkService) GetAllCached() (*LinkListResponseDTO, *apperrors.AppError) {
	links := link.AllCached()
	if links == nil {
		return nil, apperrors.NotFoundError("友情链接")
	}
	return &LinkListResponseDTO{
		Links: s.toResponseDTOList(links),
	}, nil
}
