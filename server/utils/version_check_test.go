package utils

import (
	"testing"
)

// 测试 IsValidVersion 函数
func TestIsValidVersion(t *testing.T) {
	// 测试用例
	cases := []struct {
		version string
		want    bool
	}{
		{"V1.0.0", true},
		{"V1", false},
		{"V1.0.0.0.0", false},
		{"1.0.0", false},
		{"V1.0.0-beta", false},
		{"V1234.5678.91011", false},
	}

	// 遍历测试用例
	for _, c := range cases {
		got := IsValidVersion(c.version)
		if got != c.want {
			t.Errorf("IsValidVersion(%q) == %v, want %v", c.version, got, c.want)
		}
	}
}

// 测试 CompareVersions 函数
func TestCompareVersions(t *testing.T) {
	// 测试用例
	cases := []struct {
		newVersion string
		oldVersion string
		want       int
	}{
		{"V1.0.0", "V1.0.0", VersionEqual},
		{"V2.0.0", "V1.0.0", VersionGreater},
		{"V1.0.0", "V2.0.0", VersionLess},
		{"V1.0.0", "V1.0", VersionGreater},
		{"V1.0", "V1.0.0", VersionLess},
		{"V1.0.1", "V1.0.0", VersionGreater},
		{"V1.0.0", "V1.0.1", VersionLess},
		{"V1.0.0", "V1.0.0.0", VersionLess},
		{"V1.0.0.0", "V1.0.0", VersionGreater},
		{"V1.0.0", "V1.0.0-beta", VersionInvalid},
		{"V1.0.0-beta", "V1.0.0", VersionInvalid},
		{"V1.0.0", "1.0.0", VersionInvalid},
		{"1.0.0", "V1.0.0", VersionInvalid},
		{"v2.0.23", "v1.1", VersionGreater},
		{"", "V1.1", VersionInvalid},
		{"1.0.23", "", VersionInvalid},
		{"V1.0.23", "一二三", VersionInvalid},
	}

	// 遍历测试用例
	for _, c := range cases {
		got := CompareVersions(c.newVersion, c.oldVersion)
		if got != c.want {
			t.Errorf("CompareVersions(%q, %q) == %v, want %v", c.newVersion, c.oldVersion, got, c.want)
		}
	}
}
