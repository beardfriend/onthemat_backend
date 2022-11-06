package service

import (
	"testing"

	"onthemat/internal/app/config"
	"onthemat/pkg/google"
)

func TestGetGoogleInfo(t *testing.T) {
	c := config.NewConfig()
	if err := c.Load("../../../configs"); err != nil {
		t.Error(err)
		return
	}
	gM := google.NewGoogle(c)
	authSvc := NewAuthService(nil, gM)
	if err := authSvc.GetGoogleInfo("4/0AfgeXvuUwVdNqgtAB-bignyytBY6Hab2dskz-MzkLqK1BdRfsXiMRQOUMqdUmBVlsFFqaQ"); err != nil {
		t.Error(err)
		return
	}
}
