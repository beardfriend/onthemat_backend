package kakao

import (
	"context"
	"fmt"
	"testing"
	"time"

	"onthemat/internal/app/config"

	"github.com/chromedp/chromedp"
)

func TestKakao(t *testing.T) {
	c := config.NewConfig()
	kaka := NewKakao(c)
	d := kaka.Authorize()
	fmt.Println(d)
}

func TestGetToken(t *testing.T) {
	c := config.NewConfig()
	if err := c.Load("../../configs"); err != nil {
		t.Error(err)
		return
	}
	module := NewKakao(c)
	resp := module.Authorize()

	r := resp.Header.Peek("Location")
	loginUrl := string(r)

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.Flag("headless", false), // headless를 false로 하면 브라우저가 뜨고, true로 하면 브라우저가 뜨지않는 headless 모드로 실행됨. 기본값은 true.
	)

	contextVar, cancelFunc := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancelFunc()

	contextVar, cancelFunc = chromedp.NewContext(contextVar)
	defer cancelFunc()

	contextVar, cancelFunc = context.WithTimeout(contextVar, 50*time.Second) // timeout 값을 설정
	defer cancelFunc()

	// var value string
	var infoTip string
	err := chromedp.Run(contextVar,
		chromedp.Navigate(loginUrl),
		chromedp.WaitVisible(`#input-loginKey`),
		// chromedp.Click(`#id_email_2_label`, chromedp.NodeVisible),
		chromedp.SendKeys(`#input-loginKey`, "beardfriend@kakao.com", chromedp.ByID),
		chromedp.SendKeys("#input-password", "password", chromedp.ByID),
		// chromedp.SendKeys(`#input-loginKey`, "beardfriend@kakao.com", chromedp.ByID),
		// chromedp.SendKeys("#input-password", "password", chromedp.ByID),
		chromedp.Sleep(1*time.Second), // wait for animation to finish
		chromedp.Click(`.btn_g highlight`, chromedp.NodeVisible),
		// chromedp.Value(`#id_email_2`, &value),
	)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(infoTip)
}
