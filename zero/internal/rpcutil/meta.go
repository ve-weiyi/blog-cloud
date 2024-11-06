package rpcutil

import (
	"context"
	"fmt"
	"net"
	"strings"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"

	"github.com/ve-weiyi/ve-blog-golang/kit/infra/constant"
)

func GetRPCAppName(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", fmt.Errorf("metadata error")
	}

	if val, ok := md[constant.HeaderAppName]; ok {
		if len(val) > 0 {
			return val[0], nil
		}
	}

	return "", fmt.Errorf("get rpc meta error:%v", constant.HeaderAppName)
}

func GetRPCUserId(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", fmt.Errorf("metadata error")
	}

	if val, ok := md[constant.HeaderUid]; ok {
		if len(val) > 0 {
			return val[0], nil
		}
	}

	return "", fmt.Errorf("get rpc meta error:%v", constant.HeaderUid)
}

func GetRPCUserAgent(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", fmt.Errorf("metadata error")
	}

	if val, ok := md[constant.HeaderRPCUserAgent]; ok {
		if len(val) > 0 {
			return val[0], nil
		}
	}

	return "", fmt.Errorf("get rpc meta error:%v", constant.HeaderRPCUserAgent)
}

// 获取客户端IP地址
func GetRPCClientIP(ctx context.Context) (string, error) {
	// 从上下文中获取对等体信息
	p, ok := peer.FromContext(ctx)
	if !ok {
		return "", fmt.Errorf("peer error")
	}

	// 获取客户端IP地址 host:port
	clientIP := p.Addr.String()
	host, _, err := net.SplitHostPort(strings.TrimSpace(clientIP))
	if err != nil {
		return "", err
	}

	return host, nil
}
