package repository

import (
	"context"
	"testing"

	"onthemat/internal/app/infrastructor"
	"onthemat/pkg/ent"
)

func TestCreate(t *testing.T) {
	context := context.Background()
	db := infrastructor.NewPostgresDB()
	userRepo := NewUserRepository(db)
	academyRepo := NewAcademyRepository(db)
	newUser, err := userRepo.Create(context, &ent.User{
		Email:       "asd@naver.com",
		Password:    "asdsdads",
		SocialKey:   nil,
		PhoneNum:    nil,
		Nickname:    nil,
		TermAgreeAt: nil,
	})
	if err != nil {
		t.Error(err)
	}
	if err := academyRepo.Create(context, &ent.Acadmey{
		Name:          "name",
		BusinessCode:  "123",
		CallNumber:    "01043226633",
		AddressRoad:   "도로명 주소",
		AddressSigun:  "시군",
		AddressGu:     "구",
		AddressDong:   "동",
		AddressDetail: "상세주소",
		AddressX:      "123.123",
		AddressY:      "123.123",
	}, newUser.ID); err != nil {
		t.Error(err)
	}

	err = academyRepo.Update(context, &ent.Acadmey{
		AddressX: "asd",
		AddressY: "",
	}, newUser.ID)
	if err != nil {
		t.Error(err)
	}
}
