package device

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/ve-weiyi/ve-blog-golang/kit/infra/gsm/util/scraper"

	"github.com/PuerkitoBio/goquery"
)

const BaseEndpoint = "http://www.gsmarena.com/"

func splitBrandSlug(brandSlug string) (id int, name string, err error) {
	tmpID := strings.Split(brandSlug, "-")
	lastIndex := len(tmpID) - 1
	tmpID2 := tmpID[lastIndex]

	name = strings.Join(tmpID[:lastIndex], "-")
	id, err = strconv.Atoi(tmpID2)

	return
}

func getEndpoint(brandSlug string, page int) (endpoint string, err error) {
	if page <= 1 {
		return fmt.Sprintf("%s%s.php", BaseEndpoint, brandSlug), nil
	}

	id, name, err := splitBrandSlug(brandSlug)

	if err != nil {
		return
	}

	endpoint = fmt.Sprintf("%s%s-f-%d-0-p%d.php", BaseEndpoint, name, id, page)

	return
}

func getSlug(s *goquery.Selection) (slug string, err error) {
	tmpSlug, exists := s.Attr("href")

	if exists {
		slug = strings.TrimSuffix(tmpSlug, ".php")
	} else {
		err = errors.New("Cannot resolve slug")
	}

	return
}

func getID(s *goquery.Selection) (id int, err error) {
	slug, err := getSlug(s)

	if err != nil {
		return
	}

	tmpID := strings.Split(slug, "-")
	lastIndex := len(tmpID) - 1
	tmpID2 := tmpID[lastIndex]

	id, err = strconv.Atoi(tmpID2)

	if err != nil {
		err = errors.New("Cannot resolve id")
	}

	return
}

func getDescription(s *goquery.Selection) (description string, err error) {
	description, exists := s.Attr("title")

	if !exists {
		err = errors.New("Cannot resolve description")
	}

	return
}

func getImage(s *goquery.Selection) (image string, err error) {
	image, exists := s.Attr("src")

	if !exists {
		err = errors.New("Cannot resolve image")
	}

	return
}

func getName(s *goquery.Selection) (name string, err error) {
	span := s.Find("span")

	if span == nil {
		err = errors.New("Cannot resolve name")
	}

	name = span.Text()

	return
}

func getTotalPage(s *goquery.Selection) (totalPage int, err error) {
	pages := []int{}
	for i := 0; i < s.Length(); i++ {
		page, convErr := strconv.Atoi(s.Eq(i).Text())

		if convErr != nil {
			continue
		}

		pages = append(pages, page)
	}
	totalPage = max(pages)

	if totalPage == 0 {
		totalPage = 1
	}

	return
}

func max(slice []int) int {
	var m int
	for _, e := range slice {
		if e > m {
			m = e
		}
	}
	return m
}

func newDevice(s *goquery.Selection) (device Device, err error) {
	img := s.Find("img")
	strong := s.Find("strong")

	id, err := getID(s)

	if err != nil {
		return
	}

	name, err := getName(strong)

	if err != nil {
		return
	}

	slug, err := getSlug(s)

	if err != nil {
		return
	}

	image, err := getImage(img)

	if err != nil {
		return
	}

	description, err := getDescription(img)

	if err != nil {
		return
	}

	device = Device{
		ID:          id,
		Name:        name,
		Slug:        slug,
		Image:       image,
		Description: description,
	}

	return
}

func GetDeviceList(brandSlug string, page int) (deviceList DeviceList, err error) {
	scraper := scraper.NewScraper()

	endpoint, err := getEndpoint(brandSlug, page)
	if err != nil {
		return
	}

	log.Println("endpoint", endpoint)

	document, err := scraper.Scrap(endpoint)
	if err != nil {
		return
	}

	links := document.Find(".makers").Find("a")
	pages := document.Find(".nav-pages").Find("strong, a")

	totalPage, err := getTotalPage(pages)
	if err != nil {
		return
	}

	deviceList = DeviceList{
		CurrentPage: page,
		TotalPage:   totalPage,
	}

	for i := 0; i < links.Length(); i++ {
		device, err := newDevice(links.Eq(i))

		if err != nil {
			return DeviceList{}, err
		}

		deviceList.AddDevice(device)
	}

	return
}
