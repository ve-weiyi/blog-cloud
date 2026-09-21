package middlewarex

import (
	"fmt"
	"net/http"

	"github.com/ve-weiyi/vkit/x/cryptox"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/infra/biz/bizcode"
	"github.com/ve-weiyi/blog-cloud/infra/biz/bizerr"
	"github.com/ve-weiyi/blog-cloud/infra/biz/bizheader"
	"github.com/ve-weiyi/blog-cloud/infra/responsex"
)

type DeviceTokenMiddleware struct {
}

func NewDeviceTokenMiddleware() *DeviceTokenMiddleware {
	return &DeviceTokenMiddleware{}
}

func (m *DeviceTokenMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logx.Debugf("DeviceTokenMiddleware Handle")
		ts := r.Header.Get(bizheader.HeaderTimestamp)
		tm := r.Header.Get(bizheader.HeaderXDeviceId)
		tk := r.Header.Get(bizheader.HeaderXDeviceToken)

		// 没有客户端id
		if tm == "" {
			// 拦截
			//responsex.Response(r, w, nil, bizerr.NewBizError(bizcode.CodeInvalidParam, fmt.Sprintf("request header field '%v' is missing", bizheader.HeaderXDeviceId)))
			// 默认放行
			next.ServeHTTP(w, r)
			return
		}

		// 请求头缺少参数
		if ts == "" {
			responsex.Response(r, w, nil, bizerr.NewBizError(bizcode.CodeParamMissing, fmt.Sprintf("request header field '%v' is missing", bizheader.HeaderTimestamp)))
			return
		}

		if tk == "" {
			responsex.Response(r, w, nil, bizerr.NewBizError(bizcode.CodeParamMissing, fmt.Sprintf("request header field '%v' is missing", bizheader.HeaderXDeviceToken)))
			return
		}

		// 判断 token = md5(tm,ts)
		if tk != cryptox.Sha256v(tm, ts) {
			responsex.Response(r, w, nil, bizerr.NewBizError(bizcode.CodeRequestSignInvalid, "无效请求,游客签名错误"))
			return
		}

		next.ServeHTTP(w, r)
	}
}
