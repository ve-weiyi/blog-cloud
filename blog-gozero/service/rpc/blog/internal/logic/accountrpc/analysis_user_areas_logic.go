package accountrpclogic

import (
	"context"
	"fmt"
	"strings"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/accountrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AnalysisUserAreasLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAnalysisUserAreasLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AnalysisUserAreasLogic {
	return &AnalysisUserAreasLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询用户分布区域
func (l *AnalysisUserAreasLogic) AnalysisUserAreas(in *accountrpc.AnalysisUserAreasReq) (*accountrpc.AnalysisUserAreasResp, error) {

	tableName := l.svcCtx.TUserModel.TableName()
	if in.UserType == 1 {
		tableName = l.svcCtx.TVisitorModel.TableName()
	}

	type UserArea struct {
		IpSource string `gorm:"column:ip_source"`
		Count    int    `gorm:"column:count"`
	}

	sql := fmt.Sprintf(`
SELECT 
    ip_source,
    COUNT(*) AS count
FROM 
    %s
GROUP BY 
    ip_source
ORDER BY 
    count DESC;
`, tableName)

	var result []*UserArea
	err := l.svcCtx.Gorm.Raw(sql).Scan(&result).Error
	if err != nil {
		return nil, err
	}

	// 统计地区
	areas := make(map[string]int64)
	for _, item := range result {
		if item.IpSource == "" {
			continue
		}

		area := findArea(item.IpSource)
		if _, ok := areas[area]; ok {
			areas[area] += int64(item.Count)
		} else {
			areas[area] = int64(item.Count)
		}
	}

	// 转换
	var list []*accountrpc.UserArea
	for k, v := range areas {
		list = append(list, &accountrpc.UserArea{
			Area:  k,
			Count: v,
		})
	}

	return &accountrpc.AnalysisUserAreasResp{
		List: list,
	}, nil
}

// findProvince 查找城市所属的区域
func findArea(city string) string {
	if strings.Contains(city, "省") {
		// 提取份名称
		for _, area := range provinces {
			if strings.HasPrefix(city, area) {
				return area
			}
		}
	} else if strings.Contains(city, "自治区") {
		//
		for _, area := range autonomousRegions {
			if strings.HasPrefix(city, area) {
				return area
			}
		}
	} else if strings.Contains(city, "特别行政区") {
		//
		for _, area := range specialAdministrativeRegions {
			if strings.HasPrefix(city, area) {
				return area
			}
		}
	} else {
		// 直辖市
		for _, area := range municipalities {
			if strings.HasPrefix(city, area) {
				return area
			}
		}
	}

	return "未知" // 如果未找到匹配的份
}

// 省份
var provinces = []string{
	"河北", "山西", "辽宁", "吉林", "黑龙江",
	"江苏", "浙江", "安徽", "福建", "江西", "山东",
	"河南", "湖北", "湖南", "广东", "海南",
	"四川", "贵州", "云南", "陕西", "甘肃", "青海", "台湾",
}

// 直辖市
var municipalities = []string{
	"北京市", "天津市", "上海市", "重庆市",
}

// 自治区
var autonomousRegions = []string{
	"内蒙古", "广西", "西藏", "宁夏", "新疆",
}

// 特别行政区
var specialAdministrativeRegions = []string{
	"香港", "澳门",
}
