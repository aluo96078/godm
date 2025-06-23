package odm

// With 用於指定在查詢時需要預先載入的關聯（類似 Laravel 的 with()）
// 例如：.With("posts", "comments") 會觸發對 posts 和 comments 的 $lookup
//
// With specifies which relations to eager-load during the query.
// Similar to Laravel's with(), e.g., .With("posts", "comments")
// will trigger $lookup for posts and comments collections.
func (m *GODM) With(relations ...string) *GODM {
	if m.WithRelations == nil {
		m.WithRelations = []string{}
	}
	m.WithRelations = append(m.WithRelations, relations...)
	return m
}

// SetRelationConfig 設定關聯查詢的配置，用於搭配 With() 執行 $lookup
// SetRelationConfig sets relation configurations for use with With() lookups.
func (m *GODM) SetRelationConfig(configs map[string]RelationConfig) *GODM {
	if m.RelationConfigs == nil {
		m.RelationConfigs = map[string]RelationConfig{}
	}
	for k, v := range configs {
		m.RelationConfigs[k] = v
	}
	return m
}
