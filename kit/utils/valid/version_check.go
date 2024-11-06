package valid

import (
	"regexp"
	"strconv"
	"strings"
)

const (
	VersionGreater = 1
	VersionEqual   = 0
	VersionLess    = -1
	VersionInvalid = -2
)

// Validate if version is valid or not
// Valid version could be V1.0, V1.2.3, or V1.2.3.4
func IsValidVersion(version string) bool {
	validVersion := regexp.MustCompile(`^(V)?\d{1,4}(\.\d{1,4}){1,3}$`)
	return validVersion.MatchString(version)
}

// 版本号比较
func CompareVersions(newVersion, oldVersion string) int {
	newVersion = strings.ToUpper(newVersion)
	oldVersion = strings.ToUpper(oldVersion)
	// 检查版本号是否合法，如果不合法则打印日志并返回VersionLess
	if !IsValidVersion(oldVersion) || !IsValidVersion(newVersion) {
		return VersionInvalid
	}

	// 去除版本号中的前缀V
	newVersion = strings.ReplaceAll(newVersion, "V", "")
	oldVersion = strings.ReplaceAll(oldVersion, "V", "")

	// 将版本号按照.进行分割
	newComponents := strings.Split(newVersion, ".")
	oldComponents := strings.Split(oldVersion, ".")

	// 比较每个组件的大小
	for i := 0; i < len(newComponents) && i < len(oldComponents); i++ {
		newNum, err := strconv.Atoi(newComponents[i])
		if err != nil {
			return VersionInvalid
		}
		oldNum, err := strconv.Atoi(oldComponents[i])
		if err != nil {
			return VersionInvalid
		}
		// 如果新版本的当前组件小于旧版本的当前组件，则返回VersionLess
		if newNum < oldNum {
			return VersionLess
		} else if newNum > oldNum {
			// 如果新版本的当前组件大于旧版本的当前组件，则返回VersionGreater
			return VersionGreater
		}
	}

	// 三位数的版本号比两位数的版本号大, 如V1.0.0 > V1.0
	// 如果新版本的组件数小于旧版本的组件数，则返回VersionLess
	if len(newComponents) < len(oldComponents) {
		return VersionLess
	} else if len(newComponents) > len(oldComponents) {
		// 如果新版本的组件数大于旧版本的组件数，则返回VersionGreater
		return VersionGreater
	}

	// 如果以上条件都不满足，则返回VersionEqual
	return VersionEqual
}
