package captcha

import (
	"errors"

	"pinnacle-primary-be/core/logx"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/mojocn/base64Captcha"

	captcha20230305 "github.com/alibabacloud-go/captcha-20230305/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
)

type (
	Config struct {
		EnableAliDriver bool   `yaml:"EnableAliDriver"`
		SecretID        string `yaml:"SecretID"`
		SecretKey       string `yaml:"SecretKey"`
		Endpoint        string `yaml:"Endpoint;default:captcha.cn-shanghai.aliyuncs.com"`
	}
	Option func(opt *Captcha)
	Result struct {
		Id         string `json:"id"`
		Base64Blog string `json:"raw"`
		Answer     string `json:"answer,omitempty"`
	}
	Captcha struct {
		store      base64Captcha.Store
		driver     base64Captcha.Driver
		captcha    *base64Captcha.Captcha
		aliCaptcha *captcha20230305.Client
	}
)

func NewCaptchaService(opts ...Option) *Captcha {
	// 生成默认数字
	//driver := base64Captcha.DefaultDriverDigit
	// 此尺寸的调整需要根据网站进行调试，链接：
	// https://captcha.mojotv.cn/
	c := &Captcha{
		store:  base64Captcha.DefaultMemStore,
		driver: base64Captcha.DefaultDriverDigit,
	}
	for _, opt := range opts {
		opt(c)
	}
	c.captcha = base64Captcha.NewCaptcha(c.driver, c.store)
	return c
}

func (c *Captcha) GenerateCaptcha() (*Result, error) {
	id, b64s, _, err := c.captcha.Generate()
	if err != nil {
		logx.Error("Register GetCaptchaPhoto get base64Captcha has err: ", err)
		return nil, err
	}

	return &Result{
		Id:         id,
		Base64Blog: b64s,
	}, nil
}

func (c *Captcha) VerifyCaptcha(id string, value string) bool {
	verifyResult := c.store.Verify(id, value, true)
	return verifyResult
}

func (c *Captcha) VerifyAlibabaCaptcha(value string) bool {
	request := &captcha20230305.VerifyCaptchaRequest{}
	request.CaptchaVerifyParam = tea.String(value)
	captchaVerifyResult, err := func() (captchaVerifyResult bool, _e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		resp, _err := c.aliCaptcha.VerifyCaptcha(request)
		if _err != nil {
			return true, _err
		}
		logx.Debug(resp.Body)
		if resp.Body.Result.VerifyResult != nil {
			return *resp.Body.Result.VerifyResult, nil
		}

		return true, errors.New("VerifyCaptchaResult is nil")
	}()

	if err != nil {
		logx.Error("VerifyAlibabaCaptcha has err: ", err)
		return false
	}

	return captchaVerifyResult
}

func (c *Captcha) GetAnswer(id string) string {
	return c.store.Get(id, false)
}

func WithStore(store base64Captcha.Store) Option {
	return func(c *Captcha) {
		c.store = store
	}
}

func WithAliDriver(captcha Config) Option {
	return func(c *Captcha) {
		if !captcha.EnableAliDriver {
			return
		}
		config := &openapi.Config{
			AccessKeyId:     tea.String(captcha.SecretID),
			AccessKeySecret: tea.String(captcha.SecretKey),
			Endpoint:        tea.String(captcha.Endpoint),
		}
		client, _err := captcha20230305.NewClient(config)
		if _err != nil {
			logx.Error("New Alibaba Captcha Client has err: ", _err)
			return
		}
		c.aliCaptcha = client
		logx.Debug("New Alibaba Captcha Client has success")
	}
}

func WithDriver(driver base64Captcha.Driver) Option {
	return func(c *Captcha) {
		c.driver = driver
	}
}
