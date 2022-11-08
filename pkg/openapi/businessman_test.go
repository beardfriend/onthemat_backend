package openapi

import (
	"fmt"
	"testing"

	"onthemat/internal/app/config"
)

func TestBusineessMan(t *testing.T) {
	c := config.NewConfig()
	if err := c.Load("../../configs"); err != nil {
		t.Error(err)
		return
	}
	businessModule := NewBusinessMan(c)
	resp := businessModule.GetStatus("1138621886")
	fmt.Println(resp.StatusCode())
	fmt.Println(resp)
}
