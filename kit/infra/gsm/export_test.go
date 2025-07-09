package gsm

import (
	"fmt"
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/spf13/cast"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/ve-weiyi/ve-blog-golang/kit/infra/gsm/service/brand"
	"github.com/ve-weiyi/ve-blog-golang/kit/infra/gsm/service/device"
	"github.com/ve-weiyi/ve-blog-golang/kit/infra/gsm/service/specification"
	"github.com/ve-weiyi/ve-blog-golang/kit/utils/excelx"
	"github.com/ve-weiyi/ve-blog-golang/kit/utils/jsonconv"
)

// GEN 自动生成 GORM 模型结构体文件及使用示例 https://blog.csdn.net/Jeffid/article/details/126898000
const dsn = "root:mysql7914@(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"

var db *gorm.DB

func init() {
	var err error
	// 连接数据库
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			//TablePrefix: "tb_",
			// 使用单数表名，启用该选项，此时，`User` 的表名应该是 `user`
			SingularTable: true,
		},
	})
	if err != nil {
		panic(fmt.Errorf("cannot establish db connection: %w", err))
	}
	log.Println("mysql connection done")
}

type TDeviceModelInfo struct {
	Id          int64     `json:"id" gorm:"column:id" `                     // 主键
	Brand       string    `json:"brand" gorm:"column:brand" `               // 品牌
	Slug        string    `json:"slug" gorm:"column:slug" `                 // slug
	DeviceModel string    `json:"device_model" gorm:"column:device_model" ` // 设备型号
	DeviceName  string    `json:"device_name" gorm:"column:device_name" `   // 设备名称
	DeviceType  string    `json:"device_type" gorm:"column:device_type" `   // 设备类型
	DeviceId    string    `json:"device_id" gorm:"column:device_id" `       // 设备id
	Description string    `json:"description" gorm:"column:description" `   // 描述
	ImageUrl    string    `json:"image_url" gorm:"column:image_url" `       // 图片
	Capacity    string    `json:"capacity" gorm:"column:capacity" `         // 电池容量
	Data        string    `json:"data" gorm:"column:data" `                 // json
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at" `     // 创建时间
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at" `     // 更新时间
}

type TDeviceBrand struct {
	Id              int64     `json:"id" gorm:"column:id" `                               // 主键
	Name            string    `json:"name" gorm:"column:name" `                           // 名称
	Slug            string    `json:"slug" gorm:"column:slug" `                           // slug
	BrandId         string    `json:"brand_id" gorm:"column:brand_id" `                   // 品牌id
	NumberOfDevices int       `json:"number_of_devices" gorm:"column:number_of_devices" ` // 设备数量
	CreatedAt       time.Time `json:"created_at" gorm:"column:created_at" `               // 创建时间
	UpdatedAt       time.Time `json:"updated_at" gorm:"column:updated_at" `               // 更新时间
}

func checkExist(slug string, name string) (exist bool) {
	var count int64
	db.Table(name).Where("slug = ?", slug).Count(&count)
	if count > 0 {
		exist = true
	}
	return
}

func InsertBrandModel(b brand.Brand) error {
	model := &TDeviceBrand{
		Id:              0,
		Name:            b.Name,
		Slug:            b.Slug,
		BrandId:         strconv.Itoa(b.ID),
		NumberOfDevices: b.NumberOfDevices,
	}
	return db.Create(model).Error
}

func InsertDeviceModel(b brand.Brand, d device.Device, s specification.Specification) error {
	var models string
	if s.Detail["misc"] != nil {
		models = cast.ToString(s.Detail["misc"]["models"])
	} else {
		models = ""
	}

	model := &TDeviceModelInfo{
		Id:          0,
		Brand:       s.Brand,
		Slug:        d.Slug,
		DeviceModel: models,
		DeviceName:  s.DeviceName,
		DeviceType:  s.DeviceType,
		DeviceId:    strconv.Itoa(d.ID),
		Description: d.Description,
		ImageUrl:    s.ImageURL,
		Capacity:    s.Overview.Battery.Capacity,
		Data:        jsonconv.AnyToJsonNE(s),
	}

	return db.Create(model).Error
}

func Test_Sync(t *testing.T) {
	brands, err := brand.GetAllBrands()
	if err != nil {
		log.Fatal(err)
	}

	log.Println(jsonconv.AnyToJsonIndent(brands))

	//devices, err := device.GetDeviceList("apple-phones-48", 1)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//log.Println(jsonconv.AnyToJsonIndent(devices))

	//specs, err := specification.GetSpecification("apple_iphone_16_pro_max-13123")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//log.Println(jsonconv.AnyToJsonIndent(specs))

	for _, b := range brands {
		// 只同步部分设备
		//if !strings.HasPrefix(strings.ToLower(b.Name), "a") {
		//	continue
		//}

		if !checkExist(b.Slug, "t_device_brand") {
			InsertBrandModel(b)
		}

		//if b.Name != "Apple" && b.Name != "Samsung" && b.Name != "Xiaomi" && b.Name != "Huawei" && b.Name != "Google" {
		//	continue
		//}

		devices, err := GetDeviceList(b.Slug, 1)
		if err != nil {
			log.Println(err)
		}

		log.Println(jsonconv.AnyToJsonIndent(devices))

		for _, d := range devices {
			if !checkExist(d.Slug, "t_device_model_info") {
				specs, err := specification.GetSpecification(d.Slug)
				if err != nil {
					log.Println(err)
					continue
				}

				log.Println(jsonconv.AnyToJsonIndent(specs))

				err = InsertDeviceModel(b, d, specs)
				if err != nil {
					log.Println(err)
					continue
				}

				log.Println("insert device:", d.Name)
				time.Sleep(5 * time.Second)
			}
		}

		time.Sleep(10 * time.Second)
	}

}

func GetDeviceList(slug string, page int) (list []device.Device, err error) {
	devices, err := device.GetDeviceList(slug, page)
	if err != nil {
		log.Println(err)
	}

	list = append(list, devices.Items...)
	if devices.TotalPage >= page {
		ll, err := GetDeviceList(slug, page+1)
		if err != nil {
			log.Println(err)
		}

		list = append(list, ll...)
	}

	return list, nil
}

func Test_Export(t *testing.T) {
	var brands []string
	err := db.Table("t_device_model_info").Select("brand").Group("brand").Find(&brands).Error
	log.Println(err)
	log.Println(brands)

	ex := excel.excel.NewExcelExporter()
	for _, b := range brands {
		ex.NewActiveSheet(b)
		ex.SetSheetTitle([]any{"品牌", "设备名称", "电池容量", "图片", "设备类型", "设备id", "描述"})

		var ds []TDeviceModelInfo
		db.Table("t_device_model_info").Where("brand = ?", b).Find(&ds)
		log.Println(err)

		for _, d := range ds {
			ex.AddRowValue([]any{d.Brand, d.DeviceName, d.Capacity, d.ImageUrl, d.DeviceType, d.DeviceId, d.Description})
		}
	}

	ex.ExportFile("phone.xlsx")
}
