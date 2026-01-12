package apiutils

// ExtractFields 从列表中提取字段值，去重并过滤空值
func ExtractFields[T any](list []T, extractor func(T) string) []string {
	seen := make(map[string]struct{})
	result := make([]string, 0, len(list))

	for _, item := range list {
		if id := extractor(item); id != "" {
			if _, exists := seen[id]; !exists {
				seen[id] = struct{}{}
				result = append(result, id)
			}
		}
	}
	return result
}

// ExtractMultiFields 从列表中提取多个字段值，去重并过滤空值
func ExtractMultiFields[T any](list []T, extractor func(T) []string) []string {
	seen := make(map[string]struct{})
	result := make([]string, 0, len(list))

	for _, item := range list {
		for _, id := range extractor(item) {
			if id != "" {
				if _, exists := seen[id]; !exists {
					seen[id] = struct{}{}
					result = append(result, id)
				}
			}
		}
	}
	return result
}

// BatchQuery 批量提取ID并查询数据的通用函数
func BatchQuery[T any, R any](list []T, extractor func(T) string, query func([]string) (R, error)) (R, error) {
	ids := ExtractFields(list, extractor)
	return query(ids)
}

// BatchQueryMulti 批量提取多个ID并查询数据的通用函数
func BatchQueryMulti[T any, R any](list []T, extractor func(T) []string, query func([]string) (R, error)) (R, error) {
	ids := ExtractMultiFields(list, extractor)
	return query(ids)
}
