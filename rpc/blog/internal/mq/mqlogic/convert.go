package mqlogic

// stringToPtr 空串返回 nil，避免把无效内容写入数据库。
func stringToPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
