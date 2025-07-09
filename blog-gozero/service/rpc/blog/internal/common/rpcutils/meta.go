package rpcutils

import (
	"context"
	"fmt"

	"google.golang.org/grpc/metadata"

	"github.com/ve-weiyi/ve-blog-golang/kit/infra/restx"
)

// 用户id
func GetUserIdFromCtx(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", fmt.Errorf("metadata error")
	}

	if val, ok := md[restx.HeaderUid]; ok {
		if len(val) > 0 {
			return val[0], nil
		}
	}

	return "", fmt.Errorf("missing rpc meta '%v'", restx.HeaderUid)
}

// 游客id (即是用户的终端设备id)
func GetTerminalIdFromCtx(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", fmt.Errorf("metadata error")
	}

	if val, ok := md[restx.HeaderXTerminalId]; ok {
		if len(val) > 0 {
			return val[0], nil
		}
	}

	return "", fmt.Errorf("missing rpc meta '%v'", restx.HeaderXTerminalId)
}

func GetAppNameFromCtx(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", fmt.Errorf("metadata error")
	}

	if val, ok := md[restx.HeaderAppName]; ok {
		if len(val) > 0 {
			return val[0], nil
		}
	}

	return "", fmt.Errorf("missing rpc meta '%v'", restx.HeaderAppName)
}

func GetRemoteAgentFromCtx(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", fmt.Errorf("metadata error")
	}

	if val, ok := md[restx.HeaderRPCRemoteAgent]; ok {
		if len(val) > 0 {
			return val[0], nil
		}
	}

	return "", fmt.Errorf("missing rpc meta '%v'", restx.HeaderRPCRemoteAgent)
}

// 获取客户端IP地址
func GetRemoteIPFromCtx(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", fmt.Errorf("metadata error")
	}

	if val, ok := md[restx.HeaderRPCRemoteIP]; ok {
		if len(val) > 0 {
			return val[0], nil
		}
	}

	return "", fmt.Errorf("missing rpc meta '%v'", restx.HeaderRPCRemoteIP)
}

// 获取服务端IP地址
//func GetClientIPFromCtx(ctx context.Context) (string, error) {
//	// 从上下文中获取对等体信息
//	p, ok := peer.FromContext(ctx)
//	if !ok {
//		return "", fmt.Errorf("peer error")
//	}
//
//	// 获取客户端IP地址 host:port
//	clientIP := p.Addr.String()
//	host, _, err := net.SplitHostPort(strings.TrimSpace(clientIP))
//	if err != nil {
//		return "", err
//	}
//
//	return host, nil
//}
