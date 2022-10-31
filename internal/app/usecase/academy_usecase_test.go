package usecase

import (
	"context"
	"fmt"
	"testing"

	"onthemat/internal/app/config"
	"onthemat/internal/app/infrastructor"
	"onthemat/internal/app/repository"
	"onthemat/internal/app/transport"
	"onthemat/pkg/ent"
)

func TestCreate(t *testing.T) {
	ctx := context.Background()
	c := config.NewConfig()
	if err := c.Load("../../../configs"); err != nil {
		t.Error(err)
		return
	}
	db := infrastructor.NewPostgresDB()
	academyRepo := repository.NewAcademyRepository(db)
	academyUseCase := NewAcademyUsecase(academyRepo)
	if err := academyUseCase.Create(ctx, &transport.AcademyCreateRequestBody{
		Name:          "요가학원",
		BusinessCode:  "1213213",
		CallNumber:    "01043226633",
		AddressRoad:   "서울시 강남구 언주로2길 53",
		AddressSigun:  "서울시",
		AddressGun:    "서울시",
		AddressGu:     "강남구",
		AddressDong:   "논현동",
		AddressDetail: "아무곳",
		AddressX:      "123.123",
		AddressY:      "1212.11",
	}, 1); err != nil {
		if e := ent.IsConstraintError(err); e {
			fmt.Println("이미 중복된 유저가 존재합니다.")
			return
		}
		t.Error(err)
	}
}
