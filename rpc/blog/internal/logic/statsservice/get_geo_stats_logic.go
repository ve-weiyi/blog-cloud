package statsservicelogic

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/statsrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type GetGeoStatsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGeoStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGeoStatsLogic {
	return &GetGeoStatsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取用户地理分布
func (l *GetGeoStatsLogic) GetGeoStats(in *statsrpc.GetGeoStatsRequest) (*statsrpc.GetGeoStatsResponse, error) {
	users, err := queryGeoStats(l.svcCtx.GormDB, l.svcCtx.TUserModel.TableName())
	if err != nil {
		return nil, err
	}

	visitors, err := queryGeoStats(l.svcCtx.GormDB, l.svcCtx.TGuestModel.TableName())
	if err != nil {
		return nil, err
	}

	return &statsrpc.GetGeoStatsResponse{
		Users:    users,
		Visitors: visitors,
	}, nil
}

type ipSourceCount struct {
	IpSource string `gorm:"column:ip_source"`
	Count    int64  `gorm:"column:count"`
}

func queryGeoStats(db *gorm.DB, tableName string) ([]*statsrpc.RegionStat, error) {
	sql := fmt.Sprintf(`
SELECT
    ip_source,
    COUNT(*) AS count
FROM
    %s
GROUP BY
    ip_source
ORDER BY
    count DESC
`, tableName)

	var result []*ipSourceCount
	err := db.Raw(sql).Scan(&result).Error
	if err != nil {
		return nil, err
	}

	areas := make(map[string]int64)
	for _, item := range result {
		if item.IpSource == "" {
			continue
		}
		area := findArea(item.IpSource)
		areas[area] += item.Count
	}

	var list []*statsrpc.RegionStat
	for k, v := range areas {
		list = append(list, &statsrpc.RegionStat{
			Region: k,
			Count:  v,
		})
	}

	return list, nil
}

func findArea(city string) string {
	if strings.Contains(city, "省") {
		for _, area := range provinces {
			if strings.HasPrefix(city, area) {
				return area
			}
		}
	} else if strings.Contains(city, "自治区") {
		for _, area := range autonomousRegions {
			if strings.HasPrefix(city, area) {
				return area
			}
		}
	} else if strings.Contains(city, "特别行政区") {
		for _, area := range specialAdministrativeRegions {
			if strings.HasPrefix(city, area) {
				return area
			}
		}
	} else {
		for _, area := range municipalities {
			if strings.HasPrefix(city, area) {
				return area
			}
		}
	}
	return "未知"
}

var provinces = []string{
	"河北", "山西", "辽宁", "吉林", "黑龙江",
	"江苏", "浙江", "安徽", "福建", "江西", "山东",
	"河南", "湖北", "湖南", "广东", "海南",
	"四川", "贵州", "云南", "陕西", "甘肃", "青海", "台湾",
}

var municipalities = []string{
	"北京市", "天津市", "上海市", "重庆市",
}

var autonomousRegions = []string{
	"内蒙古", "广西", "西藏", "宁夏", "新疆",
}

var specialAdministrativeRegions = []string{
	"香港", "澳门",
}
