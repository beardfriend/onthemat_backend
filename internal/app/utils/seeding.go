package utils

import (
	"context"
	"fmt"
	"time"

	"onthemat/internal/app/config"
	"onthemat/internal/app/infrastructure"
	"onthemat/internal/app/model"
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
		hased := t.as.HashPassword("asd123456!", t.c.Secret.Password)
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
	bulk := make([]*ent.AcademyCreate, 20)

	ids, _ := t.db.User.Query().Where(user.TypeIsNil()).Limit(20).IDs(context.Background())
	fmt.Println(ids)
	for i := 0; i < len(ids); i++ {
		bulk[i] = t.db.Academy.Create().
			SetAddressDetail(fake.Address().Address).
			SetAddressRoad(fake.Address().Street).
			SetBusinessCode(fmt.Sprintf("%d", fake.Number(10000000, 9999999))).
			SetCallNumber(fake.Contact().Phone).
			SetName(fake.Animal()).SetUserID(ids[i])
		t.db.User.UpdateOneID(ids[i]).SetType(model.AcademyType).Exec(context.Background())
	}
	err := t.db.Academy.CreateBulk(bulk...).Exec(context.Background())

	fmt.Println(err)
}
