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

func (t *seeding) YogaGroup() {
	bulk := make([]*ent.YogaGroupCreate, 4)
	bulk[0] = t.db.YogaGroup.Create().SetCategory("아쉬탕가").SetCategoryEng("ashtanga").SetDescription("아쉬탕가 요가입니다.")
	bulk[1] = t.db.YogaGroup.Create().SetCategory("하타").SetCategoryEng("hatha").SetDescription("하타 요가입니다.")
	bulk[2] = t.db.YogaGroup.Create().SetCategory("빈야사").SetCategoryEng("vinyasa").SetDescription("빈야사 요가입니다.")
	bulk[3] = t.db.YogaGroup.Create().SetCategory("아디다스").SetCategoryEng("adidas").SetDescription("아디다스 요가입니다.")

	err := t.db.YogaGroup.CreateBulk(bulk...).Exec(context.Background())
	fmt.Println(err)
}

func (t *seeding) Yoga() {
	bulk := make([]*ent.YogaCreate, 8)
	bulk[0] = t.db.Yoga.Create().SetNameKor("아쉬탕가 레드").SetLevel(5).SetYogaGroupID(1)
	bulk[1] = t.db.Yoga.Create().SetNameKor("아쉬탕가 프라이머리").SetLevel(4).SetYogaGroupID(1)
	bulk[2] = t.db.Yoga.Create().SetNameKor("아쉬탕가 프라이머리 하프").SetLevel(3).SetYogaGroupID(1)
	bulk[3] = t.db.Yoga.Create().SetNameKor("하타 플로우").SetLevel(3).SetYogaGroupID(2)
	bulk[4] = t.db.Yoga.Create().SetNameKor("하타 테라피").SetLevel(2).SetYogaGroupID(2)
	bulk[5] = t.db.Yoga.Create().SetNameKor("하타 기초").SetLevel(2).SetYogaGroupID(2)
	bulk[6] = t.db.Yoga.Create().SetNameKor("빈야사 플로우").SetLevel(2).SetYogaGroupID(3)
	bulk[7] = t.db.Yoga.Create().SetNameKor("빈야사 기초").SetLevel(1).SetYogaGroupID(3)
	err := t.db.Yoga.CreateBulk(bulk...).Exec(context.Background())
	fmt.Println(err)
}
