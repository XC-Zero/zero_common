package email

import (
	"gitlab.tessan.com/data-center/tessan-erp-common/config"
	"testing"
)

func TestSendEmail(t *testing.T) {

	err := SendEmail(config.EmailConfig{
		EmailServerAddr: "smtp.163.com:25",
		SenderEmail:     "t645171033@163.com",
		EmailSecret:     "DUYFQSNJVPOOZJUP",
	}, "!!!!!!!!!", "测试!", "测试!", "1572482184@qq.com", "645171033@qq.com")
	if err != nil {
		panic(err)
	}
}
