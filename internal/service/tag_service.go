package service

import (
	"errors"
	"job-platform-go2/internal/dto"
	"job-platform-go2/internal/model"
	"job-platform-go2/internal/repository"
)

type TagService struct {
	TagRepo *repository.TagRepository
}

// GetAllTagsGrouped 获取所有标签(按分类分组)
func (s *TagService) GetAllTagsGrouped() (*dto.GroupedTagsResponse, error) {
	// 1. 从 Repository 获取数据
	tagsByCategory, categories, err := s.TagRepo.GetAllTagsGrouped()
	if err != nil {
		return nil, errors.New("获取标签失败")
	}

	// 2. 构建响应
	resp := &dto.GroupedTagsResponse{
		GroupedTags: make([]dto.TagCategoryGroup, 0),
	}

	// 3. 遍历所有分类
	for _, category := range categories {
		group := dto.TagCategoryGroup{
			CategoryID:   category.ID,
			CategoryName: category.Name,
			Tags:         make([]dto.TagItem, 0),
		}

		// 获取该分类下的所有标签
		if tags, ok := tagsByCategory[category.ID]; ok {
			for _, tag := range tags {
				group.Tags = append(group.Tags, dto.TagItem{
					TagID:   tag.ID,
					TagName: tag.Name,
				})
			}
		}

		resp.GroupedTags = append(resp.GroupedTags, group)
	}

	return resp, nil
}

// CreateTag 创建自定义标签
func (s *TagService) CreateTag(req dto.CreateTagRequest, userID int) (*dto.CreateTagResponse, error) {
	// 1. 检查标签是否已存在
	exists, err := s.TagRepo.CheckTagExists(req.Name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("标签已存在")
	}

	// 2. 创建标签
	tag := &model.Tag{
		Name:       req.Name,
		CategoryID: req.CategoryID,
		CreatedBy:  userID,
	}

	if err := s.TagRepo.CreateTag(tag); err != nil {
		return nil, errors.New("创建标签失败")
	}

	// 3. 构建响应
	return &dto.CreateTagResponse{
		TagID:      tag.ID,
		TagName:    tag.Name,
		CategoryID: tag.CategoryID,
	}, nil
}
