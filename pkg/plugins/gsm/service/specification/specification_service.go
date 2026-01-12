package specification

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/ve-weiyi/ve-blog-golang/pkg/plugins/gsm/util/prettifier"
	"github.com/ve-weiyi/ve-blog-golang/pkg/plugins/gsm/util/scraper"
)

const BaseEndpoint = "http://www.gsmarena.com/"

func getEndpoint(deviceSlug string) string {
	return fmt.Sprintf("%s%s.php", BaseEndpoint, deviceSlug)
}

func getBrand(s *goquery.Document) (brand string) {
	sectionHeadings := s.Find(".section-heading")
	popularKeyword := "Popular from "
	for i := 0; i < sectionHeadings.Length(); i++ {
		rawBrand := sectionHeadings.Eq(i).Find("a").Eq(0).Text()
		if strings.Contains(rawBrand, popularKeyword) {
			brand = strings.Replace(rawBrand, popularKeyword, "", 1)

			return
		}
	}

	return
}

func getDeviceName(s *goquery.Selection) string {
	return s.Find(".specs-phone-name-title").Text()
}

func getImageURL(s *goquery.Selection) (imageURL string, err error) {
	imageURL, imageURLExists := s.Find(".specs-photo-main").
		Find("img").Attr("src")

	if !imageURLExists {
		err = errors.New("Cannot resolve image url")
	}

	return
}

func getDeviceType(s *goquery.Document) string {
	documentTitle := s.Find("title").Text()
	var deviceType string

	if strings.Contains(strings.ToLower(documentTitle), strings.ToLower(Phone)) {
		deviceType = Phone
	} else {
		deviceType = Tablet
	}

	return deviceType
}

func getSpecificationOverviewValueDetail(s *goquery.Selection) (exists bool, value1, value2 string) {
	exists = false

	if s != nil && s.Length() > 0 {
		exists = true
		parent := s.Parent()

		value2 = strings.TrimSpace(parent.Find("strong").Text())
		parent.Find("strong").Remove()

		value1 = strings.TrimSpace(strings.Replace(parent.Text(), s.Text(), "", 1))
	}

	return
}

func getSpecificationOverviewValue(s *goquery.Selection, s2 *goquery.Selection) (value1 string, value2 string) {
	exists, value1, value2 := getSpecificationOverviewValueDetail(s)

	if !exists {
		exists, value1, value2 = getSpecificationOverviewValueDetail(s2)
	}

	return
}

func hetSpecificationOverviewValue(s *goquery.Selection, s2 *goquery.Selection) (value string) {
	if s != nil && s.Length() == 1 {
		value = "Yes"

		return
	}

	if s2 != nil && s2.Length() == 1 {
		value = "Yes"

		return
	}

	value = "No"

	return
}

func getSpecificationOverview(s *goquery.Selection) (specificationOverview SpecificationOverview) {
	launched, _ := getSpecificationOverviewValue(s.Find(".icon-launched"), nil)
	body, _ := getSpecificationOverviewValue(s.Find(".icon-mobile"), s.Find(".icon-mobile2"))
	os, _ := getSpecificationOverviewValue(s.Find(".icon-os"), nil)
	storage, _ := getSpecificationOverviewValue(s.Find(".icon-sd-card-0"), s.Find(".icon-sd-card-1"))

	touchscreen := hetSpecificationOverviewValue(s.Find(".icon-touch-0"), s.Find(".icon-touch-1"))
	resolution, size := getSpecificationOverviewValue(s.Find(".icon-touch-0"), s.Find(".icon-touch-1"))

	video, photo := getSpecificationOverviewValue(s.Find(".icon-camera-0"), s.Find(".icon-camera-1"))

	chipset, ram := getSpecificationOverviewValue(s.Find(".icon-cpu"), nil)

	technology, capacity := getSpecificationOverviewValue(s.Find(".icon-battery-0"), s.Find(".icon-battery-1"))

	specificationOverview = SpecificationOverview{
		GeneralInfo: SpecificationOverviewGeneralInfo{
			Launched: launched,
			Body:     body,
			OS:       os,
			Storage:  storage,
		},
		Display: SpecificationOverviewDisplay{
			Touchscreen: touchscreen,
			Size:        size,
			Resolution:  resolution,
		},
		Camera: SpecificationOverviewCamera{
			Photo: photo,
			Video: video,
		},
		Expansion: SpecificationOverviewExpansion{
			RAM:     ram,
			Chipset: chipset,
		},
		Battery: SpecificationOverviewBattery{
			Capacity:   capacity,
			Technology: technology,
		},
	}

	return
}

func getSpecificationDetail(s *goquery.Selection) (specificationDetail SpecificationDetail) {
	specificationDetail = SpecificationDetail{}
	tableSelector := s.Find("table")
	for i := 0; i < tableSelector.Length(); i++ {
		categorySelector := tableSelector.Eq(i).Find("th")
		for j := 0; j < categorySelector.Length(); j++ {
			// get category
			currentCategory := categorySelector.Eq(j)
			category := prettifier.PrettifyKey(strings.TrimSpace(currentCategory.Text()))

			// get top title-value pair
			title := prettifier.PrettifyKey(strings.TrimSpace(currentCategory.Siblings().Eq(0).Text()))
			value := prettifier.PrettifyValue(currentCategory.Siblings().Eq(0).Next().Text())
			specificationDetail[category] = map[string]interface{}{}
			specificationDetail[category][title] = value

			// get next title-value pair
			nextItems := categorySelector.Parent().Siblings()
			for k := 0; k < nextItems.Length(); k++ {
				title := prettifier.PrettifyKey(strings.TrimSpace(nextItems.Eq(k).Find(".ttl").Text()))
				value := prettifier.PrettifyValue(nextItems.Eq(k).Find(".nfo").Text())
				if !(title == "" && value == "") {
					specificationDetail[category][title] = value
				}
			}
		}
	}

	return
}

func GetSpecification(slug string) (specification Specification, err error) {
	scraper := scraper.NewScraper()

	endpoint := getEndpoint(slug)
	if err != nil {
		return
	}

	log.Println("endpoint", endpoint)

	document, err := scraper.Scrap(endpoint)
	if err != nil {
		return
	}

	overviewSelector := document.Find(".article-info")
	brand := getBrand(document)

	deviceName := getDeviceName(overviewSelector)

	deviceType := getDeviceType(document)

	imageURL, err := getImageURL(overviewSelector)

	if err != nil {
		return
	}

	specificationOverview := getSpecificationOverview(overviewSelector)

	detailSelector := document.Find("#specs-list")
	specificationDetail := getSpecificationDetail(detailSelector)

	specification = Specification{
		Brand:      brand,
		DeviceName: deviceName,
		DeviceType: deviceType,
		ImageURL:   imageURL,
		Overview:   specificationOverview,
		Detail:     specificationDetail,
	}

	return
}
