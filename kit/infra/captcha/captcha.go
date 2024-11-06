package captcha

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/mojocn/base64Captcha"
	"github.com/redis/go-redis/v9"
)

// 验证码仓库
type CaptchaHolder struct {
	randSource *rand.Rand // 随机数种子

	DefaultHeight  int     // 默认高度 40
	DefaultWidth   int     // 默认宽度 80
	DefaultLength  int     // 默认长度,位数 6
	DefaultMaxSkew float64 // 默认倾斜因子 0.7
	DefaultDotRate float64 // 默认干扰点比率 25%

	store base64Captcha.Store

	DriverAudio   *base64Captcha.DriverAudio
	DriverString  *base64Captcha.DriverString
	DriverChinese *base64Captcha.DriverChinese
	DriverMath    *base64Captcha.DriverMath
	DriverDigit   *base64Captcha.DriverDigit
}

type Option func(*CaptchaHolder)

func WithRedisStore(rdb *redis.Client) Option {
	return func(o *CaptchaHolder) {
		o.store = NewRedisStore(rdb)
	}
}

func NewCaptchaHolder(options ...Option) *CaptchaHolder {
	ch := &CaptchaHolder{
		randSource:     rand.New(rand.NewSource(time.Now().UnixNano())),
		DefaultHeight:  80,
		DefaultWidth:   240,
		DefaultLength:  6,
		DefaultMaxSkew: 0.7,
		DefaultDotRate: 0.20,
		store:          base64Captcha.DefaultMemStore,
		DriverAudio:    base64Captcha.DefaultDriverAudio,
		DriverString:   base64Captcha.NewDriverString(80, 240, 0, 0, 5, "1234567890abcdefghijklmnopqrstuvwxyz", nil, nil, nil),
		DriverChinese:  base64Captcha.NewDriverChinese(80, 240, 0, 0, 5, "1234567890abcdefghijklmnopqrstuvwxyz", nil, nil, nil),
		DriverMath:     base64Captcha.NewDriverMath(80, 240, 0, 0, nil, nil, nil),
		DriverDigit:    base64Captcha.NewDriverDigit(40, 80, 6, 0.7, 10),
	}

	// 应用选项
	for _, option := range options {
		option(ch)
	}

	return ch
}

// 生成随机验证码
func (rs *CaptchaHolder) GetCodeCaptcha(key string) (code string, err error) {
	var randomInt string
	// 生成随机6位整数
	for i := 0; i < rs.DefaultLength; i++ {
		randomInt = randomInt + strconv.Itoa(rs.randSource.Intn(10))
	}

	err = rs.store.Set(key, randomInt)
	if err != nil {
		return "", err
	}

	return randomInt, nil
}

func (rs *CaptchaHolder) GetImageCaptcha(CaptchaType string, height int, width int, length int) (string, string, error) {
	var driver base64Captcha.Driver

	if height == 0 {
		height = rs.DefaultHeight
	}
	if width == 0 {
		width = rs.DefaultWidth
	}
	if length == 0 {
		length = rs.DefaultLength
	}

	var dotCount = int(float64(height)*rs.DefaultDotRate + float64(width)*rs.DefaultDotRate)

	//create base64 encoding captcha
	switch CaptchaType {
	case "audio":
		driver = rs.DriverAudio
	case "string":
		driver = rs.DriverString.ConvertFonts()
	case "math":
		driver = rs.DriverMath.ConvertFonts()
	case "chinese":
		driver = rs.DriverChinese.ConvertFonts()
	case "digit":
		driver = base64Captcha.NewDriverDigit(height, width, length, rs.DefaultMaxSkew, dotCount)
	default:
		driver = base64Captcha.NewDriverDigit(height, width, length, rs.DefaultMaxSkew, dotCount)
	}

	c := base64Captcha.NewCaptcha(driver, rs.store)
	id, b64s, _, err := c.Generate()

	return id, b64s, err
}

// 验证验证码
func (rs *CaptchaHolder) VerifyCaptcha(id string, answer string) bool {
	return rs.store.Verify(id, answer, true)
}
