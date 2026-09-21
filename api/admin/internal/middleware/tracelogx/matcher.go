package tracelogx

import "strings"

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

func splitPathSegments(path string) []string {
	path = strings.Trim(path, "/")
	if path == "" {
		return []string{}
	}
	return strings.Split(path, "/")
}

func matchPathSegments(pattern []string, path []string) bool {
	if len(pattern) == 0 {
		return len(path) == 0
	}

	if pattern[0] == "**" {
		if len(pattern) == 1 {
			return true
		}
		for i := 0; i <= len(path); i++ {
			if matchPathSegments(pattern[1:], path[i:]) {
				return true
			}
		}
		return false
	}

	if len(path) == 0 {
		return false
	}

	if !segmentMatch(pattern[0], path[0]) {
		return false
	}
	return matchPathSegments(pattern[1:], path[1:])
}

func segmentMatch(patternSeg string, pathSeg string) bool {
	return patternSeg == "*" || isParamSegment(patternSeg) || patternSeg == pathSeg
}

func isParamSegment(seg string) bool {
	return (strings.HasPrefix(seg, ":") && len(seg) > 1) ||
		(strings.HasPrefix(seg, "{") && strings.HasSuffix(seg, "}") && len(seg) > 2)
}
