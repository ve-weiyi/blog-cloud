package client

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-stomp/stomp/v3"
	"github.com/go-stomp/stomp/v3/frame"
)

const (
	// 最大心跳超时值，约 11.5 天
	maxHeartBeat = 999999999
)

var (
	// 心跳头部值的正则表达式
	heartBeatRegexp = regexp.MustCompile("^[0-9]{1,9},[0-9]{1,9}$")
)

func frameToBytes(f *frame.Frame) []byte {
	var buf bytes.Buffer
	frame.NewWriter(&buf).Write(f)
	return buf.Bytes()
}

func uint64ToString(n uint64) string {
	return strconv.FormatUint(n, 10)
}

func stringToUint64(s string) uint64 {
	n, _ := strconv.ParseUint(s, 10, 64)
	return n
}

// parseHeartBeat 解析 heart-beat 头部 "cx,cy"
func parseHeartBeat(hb string) (cx, cy time.Duration) {
	if hb == "" {
		return 0, 0
	}

	// 验证格式
	if !heartBeatRegexp.MatchString(hb) {
		return 0, 0
	}

	parts := strings.Split(hb, ",")
	cxMs, _ := strconv.ParseUint(parts[0], 10, 32)
	cyMs, _ := strconv.ParseUint(parts[1], 10, 32)

	// 检查最大值
	if cxMs > maxHeartBeat {
		cxMs = maxHeartBeat
	}
	if cyMs > maxHeartBeat {
		cyMs = maxHeartBeat
	}

	return time.Duration(cxMs) * time.Millisecond, time.Duration(cyMs) * time.Millisecond
}

// formatHeartBeat 格式化 heart-beat 头部
func formatHeartBeat(cx, cy time.Duration) string {
	return fmt.Sprintf("%d,%d", cx/time.Millisecond, cy/time.Millisecond)
}

// negotiateVersion 协商 STOMP 协议版本
func negotiateVersion(acceptVersion string) (string, error) {
	if acceptVersion == "" {
		return stomp.V10.String(), nil
	}

	supportedVersions := []stomp.Version{stomp.V12, stomp.V11, stomp.V10}
	clientVersions := strings.Split(acceptVersion, ",")

	// 选择最高的共同支持版本
	for _, supported := range supportedVersions {
		for _, client := range clientVersions {
			if strings.TrimSpace(client) == supported.String() {
				return supported.String(), nil
			}
		}
	}

	return "", errUnsupportedVersion
}
