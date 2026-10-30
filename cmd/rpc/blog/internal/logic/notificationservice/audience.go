package notificationservicelogic

import (
	"encoding/json"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/notification"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizcode"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizerr"
)

// validateAudience 校验公告受众：目标类型必须已知，定向受众必须是合法的 JSON 字符串数组。
//
// 允许空数组 —— 草稿阶段可以先定类型、后选人，是否为空由发布时裁决。
// 格式留到发布才校验的话，失败发生在 decodeAudience 里，公告会卡在草稿且看不出原因。
func validateAudience(targetType, targetIDs string) error {
	if !notification.ValidAudience(targetType) {
		return bizerr.NewBizError(bizcode.CodeInvalidParam, "目标类型不合法")
	}
	if targetType != notification.AudienceUserIDs || targetIDs == "" {
		return nil
	}

	var ids []string
	if err := json.Unmarshal([]byte(targetIDs), &ids); err != nil {
		return bizerr.NewBizError(bizcode.CodeParamFormat, "定向用户必须是 JSON 字符串数组")
	}
	return nil
}
