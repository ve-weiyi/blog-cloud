package metax

import (
	"context"
	"fmt"

	"google.golang.org/grpc/metadata"

	"github.com/ve-weiyi/blog-cloud/infra/biz/bizheader"
)

// 用户id
func GetApiUserIdFromCtx(ctx context.Context) (string, error) {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return "", fmt.Errorf("metadata error")
	}

	if val, ok := md[bizheader.HeaderUid]; ok {
		if len(val) > 0 {
			return val[0], nil
		}
	}

	return "", fmt.Errorf("missing rpc meta '%v'", bizheader.HeaderUid)
}

// Token
func GetApiTokenFromCtx(ctx context.Context) (string, error) {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return "", fmt.Errorf("metadata error")
	}

	if val, ok := md[bizheader.HeaderToken]; ok {
		if len(val) > 0 {
			if len(val) > 0 {
				return val[0], nil
			}
		}
	}

	return "", fmt.Errorf("missing rpc meta '%v'", bizheader.HeaderToken)
}

// 游客id (即是用户的终端设备id)
func GetApiDeviceIdFromCtx(ctx context.Context) (string, error) {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return "", fmt.Errorf("metadata error")
	}

	if val, ok := md[bizheader.HeaderXDeviceId]; ok {
		if len(val) > 0 {
			return val[0], nil
		}
	}

	return "", fmt.Errorf("missing rpc meta '%v'", bizheader.HeaderXDeviceId)
}

func GetApiAppNameFromCtx(ctx context.Context) (string, error) {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return "", fmt.Errorf("metadata error")
	}

	if val, ok := md[bizheader.HeaderAppName]; ok {
		if len(val) > 0 {
			return val[0], nil
		}
	}

	return "", fmt.Errorf("missing rpc meta '%v'", bizheader.HeaderAppName)
}

func GetApiRemoteAgentFromCtx(ctx context.Context) (string, error) {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return "", fmt.Errorf("metadata error")
	}

	if val, ok := md[bizheader.HeaderRPCRemoteAgent]; ok {
		if len(val) > 0 {
			return val[0], nil
		}
	}

	return "", fmt.Errorf("missing rpc meta '%v'", bizheader.HeaderRPCRemoteAgent)
}
