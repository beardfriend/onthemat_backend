package utils

import (
	"context"
	"fmt"
	"time"

	"onthemat/internal/app/config"
	"onthemat/internal/app/infrastructure"
	"onthemat/internal/app/service"
	"onthemat/pkg/ent"
	"onthemat/pkg/ent/user"

	fake "github.com/brianvoe/gofakeit/v6"
)

type seeding struct {
	db *ent.Client
	as service.AuthService
	c  *config.Config
}

func NewSeeding() *seeding {
	c := config.NewConfig()
	if err := c.Load("./configs"); err != nil {
		panic(err)
	}

	cli := infrastructure.NewPostgresDB(c)
	as := service.NewAuthService(nil, nil, nil, nil)
	return &seeding{
		db: cli,
		as: as,
		c:  c,
	}
}

func (t *seeding) Users() {
	bulk := make([]*ent.UserCreate, 40)

	for i := 0; i < 40; i++ {
		hased := t.as.HashPassword("asd1234", t.c.Secret.Password)
		bulk[i] = t.db.User.Create().
			SetEmail(fake.Email()).
			SetIsEmailVerified(true).
			SetNickname(fake.Name()).
			SetPhoneNum(fake.Phone()).
			SetPassword(hased).
			SetTermAgreeAt(time.Now())
	}
	err := t.db.User.CreateBulk(bulk...).Exec(context.Background())
	fmt.Println(err)
}

func (t *seeding) Academies() {
	bulk := make([]*ent.AcademyCreate, 40)

	ids, _ := t.db.User.Query().Where(user.TypeIsNil()).IDs(context.Background())
	for i := 0; i < len(ids); i++ {
		bulk[i] = t.db.Academy.Create().
			SetAddressSigun(fake.Address().State).
			SetAddressDetail(fake.Address().Address).
			SetAddressDong(fake.Address().City).
			SetAddressGu(fake.Address().Street).
			SetAddressRoad(fake.Address().Street).
			SetAddressX("123.3").
			SetAddressY("1234.4").
			SetBusinessCode(fmt.Sprintf("%d", fake.Number(10000000, 9999999))).
			SetCallNumber(fake.Contact().Phone).
			SetName(fake.Animal()).SetUserID(ids[i])
	}
	err := t.db.Academy.CreateBulk(bulk...).Exec(context.Background())
	fmt.Println(err)
}
