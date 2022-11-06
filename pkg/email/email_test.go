package email

import (
	"testing"

	"onthemat/internal/app/config"
)

func TestSendEmail(t *testing.T) {
	c := config.NewConfig()
	if err := c.Load("../../configs"); err != nil {
		t.Error(err)
		return
	}

	mailModule := NewEmail(c)
	to := []string{"beardfriend21@gmail.com"}
	subject := "본문 제목\r\n"
	blank := "\r\n"
	body := "본문 냉용\r\n"
	msg := []byte(subject + blank + body)
	if err := mailModule.Send(to, msg); err != nil {
		t.Error(err)
		return
	}
}
