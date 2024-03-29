package request

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/server/global"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/sqlx"
	"github.com/ve-weiyi/ve-blog-golang/server/utils/iputil"
)

// 请求上下文,一般存放请求头参数
type Context struct {
	context.Context `json:"-" header:"-"`
	Token           string `json:"token" header:"token" example:""`
	UID             int    `json:"uid" header:"-" example:""`
	Username        string `json:"username" header:"-" example:""`
	IpAddress       string `json:"ip_address" header:"-" example:""`
	Agent           string `json:"agent" header:"-" example:""`
}

func (s *Context) GetContext() context.Context {
	return s.Context
}

func (s *Context) GetIpSource() string {
	ip := s.IpAddress
	location, err := iputil.GetIpInfoByBaidu(ip)
	if err != nil {
		global.LOG.Println("GetIpInfoByBaidu:", err)
		return "未知ip"
	}
	return location.Location
}

// PageQuery Paging common input parameter structure
type PageQuery struct {
	sqlx.PageLimit
	Sorts      []*sqlx.Sort      `json:"sorts" form:"sorts"`           // 排序
	Conditions []*sqlx.Condition `json:"conditions" form:"conditions"` // 使用条件语句查询
}

func (s *PageQuery) PageClause() (int, int) {
	return s.PageLimit.Limit(), s.PageLimit.Offset()
}

func (s *PageQuery) OrderClause() string {
	return sqlx.OrderClause(s.Sorts)
}

func (s *PageQuery) ConditionClause() (string, []interface{}) {
	return sqlx.ConditionClause(s.Conditions)
}
