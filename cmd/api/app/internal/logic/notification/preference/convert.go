package preference

import (
	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/types"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/notificationservice"
)

// convertPreference 把偏好矩阵的一格映射为读者端视图。
func convertPreference(in *notificationservice.NotificationPreference) types.NotificationPreferenceVO {
	return types.NotificationPreferenceVO{
		Group:     in.Group,
		Channel:   in.Channel,
		Enabled:   in.Enabled,
		Mandatory: in.Mandatory,
	}
}
