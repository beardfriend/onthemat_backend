package google

import (
	"fmt"
	"testing"

	"onthemat/internal/app/config"
)

func TestGetToken(t *testing.T) {
	c := config.NewConfig()
	if err := c.Load("../../configs"); err != nil {
		t.Error(err)
		return
	}
	googleModule := NewGoogle(c)

	// 코드 값 가져온 결과
	resp := googleModule.GetToken("4/0AfgeXvtCkAy9NWE5yZ0TylMhFEeOe-fOS12y-1M3EdZUu7Ih4D6x1Sw2ytGPif6gntCE_g")
	if resp.StatusCode() != 200 {
		t.Error("에러")
		return
	}
	fmt.Println(resp)
}
