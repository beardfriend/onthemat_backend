package naver

import (
	"fmt"
	"testing"

	"onthemat/internal/app/config"
)

func TestNaver(t *testing.T) {
	c := config.NewConfig()
	if err := c.Load("../../configs"); err != nil {
		t.Error(err)
		return
	}

	nM := NewNaver(c)
	url := nM.Authorize()
	fmt.Println(url)
	// nM.GetToken()
}
