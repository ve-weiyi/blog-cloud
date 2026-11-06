package quota

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/ai/aictx"
	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/ai/aiquota"
	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/svc"
	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/types"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizcode"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizerr"
)

type GetAiQuotaLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询今日剩余额度
func NewGetAiQuotaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAiQuotaLogic {
	return &GetAiQuotaLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAiQuotaLogic) GetAiQuota() (resp *types.AiQuotaResp, err error) {
	identity, ok := aictx.IdentityFrom(l.ctx)
	if !ok {
		return nil, bizerr.NewBizError(bizcode.CodeUnauthenticated, "缺少身份信息")
	}

	limit := aiquota.Limit(l.svcCtx.Config.AiConf.Quota, identity.Guest)
	key := aiquota.DailyKey(aiquota.Subject(identity), time.Now().Format("20060102"))
	used, err := aiquota.Used(l.ctx, l.svcCtx.RedisClient, key)
	if err != nil {
		l.Errorf("读取 AI 配额失败: %v", err)
		used = 0
	}

	remaining := int64(-1) // -1 表示不限额度
	if limit > 0 {
		remaining = limit - used
		if remaining < 0 {
			remaining = 0
		}
	}

	return &types.AiQuotaResp{
		DailyQuota: limit,
		Used:       used,
		Remaining:  remaining,
	}, nil
}
