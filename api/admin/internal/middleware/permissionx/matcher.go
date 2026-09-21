package permissionx

import "strings"

// 规范化 HTTP 方法（大小写/通配符）
func normalizeMethod(method string) string {
	method = strings.TrimSpace(method)
	if method == "" {
		return ""
	}
	method = strings.ToUpper(method)
	if method == "*" || method == "ALL" {
		return "*"
	}
	return method
}

// 规范化请求路径（去除 query、尾随斜杠、补齐前导斜杠）
func normalizePath(path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		return ""
	}
	if idx := strings.IndexByte(path, '?'); idx >= 0 {
		path = path[:idx]
	}
	if path != "/" {
		path = strings.TrimSuffix(path, "/")
	}
	if path != "" && !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return path
}

// matchPermission 支持通配符路径匹配（策略中的 pattern 匹配请求中的 target）
// pattern 示例: "GET:/api/v1/users/:id"  或  "GET:/api/v1/users/*"
// target  示例: "GET:/api/v1/users/123"
func matchPermission(pattern, target string) bool {
	if pattern == target {
		return true
	}
	// 方法部分必须完全一致（或 pattern 方法为 *）
	pm, pp, ok1 := splitPermID(pattern)
	tm, tp, ok2 := splitPermID(target)
	if !ok1 || !ok2 {
		return false
	}
	if pm != "*" && pm != tm {
		return false
	}
	return pathMatch(pp, tp)
}

// splitPermID 将 "METHOD:PATH" 拆分为 (method, path, ok)
func splitPermID(id string) (method, path string, ok bool) {
	idx := strings.IndexByte(id, ':')
	if idx < 0 {
		return "", "", false
	}
	return id[:idx], id[idx+1:], true
}

// pathMatch 支持 :param 和 * 通配符的路径匹配
func pathMatch(pattern, path string) bool {
	pp := strings.Split(strings.Trim(pattern, "/"), "/")
	tp := strings.Split(strings.Trim(path, "/"), "/")

	i, j := 0, 0
	for i < len(pp) && j < len(tp) {
		seg := pp[i]
		if seg == "*" {
			return true // * 匹配剩余所有段
		}
		if strings.HasPrefix(seg, ":") {
			// :param 匹配任意单段
			i++
			j++
			continue
		}
		if seg != tp[j] {
			return false
		}
		i++
		j++
	}
	return i == len(pp) && j == len(tp)
}
