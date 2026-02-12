package dto

// ==================== 标签相关 DTO ====================

// CreateTagRequest 创建自定义标签请求
type CreateTagRequest struct {
	Name       string `json:"name" binding:"required"`        // 新标签的名称
	CategoryID int    `json:"category_id" binding:"required"` // 所属分类ID
}

// CreateTagResponse 创建标签响应
type CreateTagResponse struct {
	TagID      int    `json:"tag_id"`
	TagName    string `json:"tag_name"`
	CategoryID int    `json:"category_id"`
}

// GroupedTagsResponse 分组标签响应
type GroupedTagsResponse struct {
	GroupedTags []TagCategoryGroup `json:"grouped_tags"`
}

// TagCategoryGroup 标签分类组
type TagCategoryGroup struct {
	CategoryID   int       `json:"category_id"`
	CategoryName string    `json:"category_name"`
	Tags         []TagItem `json:"tags"`
}

// TagItem 标签项
type TagItem struct {
	TagID   int    `json:"tag_id"`
	TagName string `json:"tag_name"`
}
