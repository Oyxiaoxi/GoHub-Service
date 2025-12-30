// Package services 互动业务逻辑
package services

import (
	"GoHub-Service/app/repositories"
	apperrors "GoHub-Service/pkg/errors"
)

// InteractionService 提供点赞、收藏、关注、浏览等互动能力
type InteractionService struct {
	repo      repositories.InteractionRepository
	topicRepo repositories.TopicRepository
	userRepo  repositories.UserRepository
}

// NewInteractionService 创建互动服务实例
func NewInteractionService() *InteractionService {
	return &InteractionService{
		repo:      repositories.NewInteractionRepository(),
		topicRepo: repositories.NewTopicRepository(),
		userRepo:  repositories.NewUserRepository(),
	}
}

func (s *InteractionService) LikeTopic(userID, topicID string) *apperrors.AppError {
	if _, err := s.topicRepo.GetByID(topicID); err != nil {
		return apperrors.WrapError(err, "获取话题失败")
	}
	if err := s.repo.LikeTopic(userID, topicID); err != nil {
		return apperrors.WrapError(err, "点赞话题失败")
	}
	return nil
}

func (s *InteractionService) UnlikeTopic(userID, topicID string) *apperrors.AppError {
	if _, err := s.topicRepo.GetByID(topicID); err != nil {
		return apperrors.WrapError(err, "获取话题失败")
	}
	if err := s.repo.UnlikeTopic(userID, topicID); err != nil {
		return apperrors.WrapError(err, "取消点赞失败")
	}
	return nil
}

func (s *InteractionService) FavoriteTopic(userID, topicID string) *apperrors.AppError {
	if _, err := s.topicRepo.GetByID(topicID); err != nil {
		return apperrors.WrapError(err, "获取话题失败")
	}
	if err := s.repo.FavoriteTopic(userID, topicID); err != nil {
		return apperrors.WrapError(err, "收藏话题失败")
	}
	return nil
}

func (s *InteractionService) UnfavoriteTopic(userID, topicID string) *apperrors.AppError {
	if _, err := s.topicRepo.GetByID(topicID); err != nil {
		return apperrors.WrapError(err, "获取话题失败")
	}
	if err := s.repo.UnfavoriteTopic(userID, topicID); err != nil {
		return apperrors.WrapError(err, "取消收藏失败")
	}
	return nil
}

func (s *InteractionService) FollowUser(followerID, followeeID string) *apperrors.AppError {
	if _, err := s.userRepo.GetByID(followeeID); err != nil {
		return apperrors.WrapError(err, "获取用户失败")
	}
	if err := s.repo.FollowUser(followerID, followeeID); err != nil {
		if err.Error() == "cannot follow yourself" {
			return apperrors.ValidationError("不能关注自己", map[string]interface{}{"user_id": followeeID})
		}
		return apperrors.WrapError(err, "关注用户失败")
	}
	return nil
}

func (s *InteractionService) UnfollowUser(followerID, followeeID string) *apperrors.AppError {
	if _, err := s.userRepo.GetByID(followeeID); err != nil {
		return apperrors.WrapError(err, "获取用户失败")
	}
	if err := s.repo.UnfollowUser(followerID, followeeID); err != nil {
		if err.Error() == "cannot follow yourself" {
			return apperrors.ValidationError("不能取消关注自己", map[string]interface{}{"user_id": followeeID})
		}
		return apperrors.WrapError(err, "取消关注失败")
	}
	return nil
}

func (s *InteractionService) AddTopicView(topicID string) *apperrors.AppError {
	if _, err := s.topicRepo.GetByID(topicID); err != nil {
		return apperrors.WrapError(err, "获取话题失败")
	}
	if err := s.repo.IncrementTopicView(topicID); err != nil {
		return apperrors.WrapError(err, "增加浏览量失败")
	}
	return nil
}
