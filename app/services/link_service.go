// Package services 友情链接业务逻辑服务
package services

import (
	"time"

	"GoHub-Service/app/models/link"
	apperrors "GoHub-Service/pkg/errors"
	"GoHub-Service/pkg/mapper"
)

// LinkService 提供友情链接查询的只读服务，主要走缓存避免频繁 DB 访问.
type LinkService struct {
	mapper mapper.Mapper[link.Link, LinkResponseDTO] // 使用泛型Mapper消除DTO转换重复
}

// NewLinkService 创建友情链接服务实例
func NewLinkService() *LinkService {
	// 定义DTO转换函数（只需一次）
	converter := func(l *link.Link) *LinkResponseDTO {
		return &LinkResponseDTO{
			ID:        l.GetStringID(),
			Name:      l.Name,
			URL:       l.URL,
			CreatedAt: l.CreatedAt,
			UpdatedAt: l.UpdatedAt,
		}
	}

	return &LinkService{
		mapper: mapper.NewSimpleMapper(converter),
	}
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

// toResponseDTO 使用Mapper将Link模型转换为响应DTO
// 优化：使用泛型Mapper消除重复代码
func (s *LinkService) toResponseDTO(l *link.Link) *LinkResponseDTO {
	return s.mapper.ToDTO(l)
}

// toResponseDTOList 使用Mapper将Link模型列表转换为响应DTO列表
// 优化：使用泛型Mapper消除重复代码，自动优化内存拷贝
func (s *LinkService) toResponseDTOList(links []link.Link) []LinkResponseDTO {
	return s.mapper.ToDTOList(links)
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
