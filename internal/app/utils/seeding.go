package utils

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"onthemat/internal/app/config"
	"onthemat/internal/app/infrastructure"
	"onthemat/internal/app/model"
	"onthemat/internal/app/service"
	"onthemat/internal/app/transport"
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
			SetNickname(fake.Name()[:8]).
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
	for i := 0; i < len(ids); i++ {
		bulk[i] = t.db.Academy.Create().
			SetAddressDetail(fake.Address().Address).
			SetAddressRoad(fake.Address().Street).
			SetBusinessCode(fmt.Sprintf("%d", fake.Number(10000000, 9999999))).
			SetCallNumber(fake.Contact().Phone).
			SetName(fake.Animal()).SetUserID(ids[i]).SetAreaSigunguID(i + 1)
		t.db.User.UpdateOneID(ids[i]).SetType(model.AcademyType).Exec(context.Background())
	}
	err := t.db.Academy.CreateBulk(bulk...).Exec(context.Background())

	fmt.Println(err)
}

func (t *seeding) Teachers() {
	bulk := make([]*ent.TeacherCreate, 10)

	ids, _ := t.db.User.Query().Where(user.TypeIsNil()).Limit(10).IDs(context.Background())
	for i := 0; i < len(ids); i++ {
		yogaIds := []int{fake.Number(1, 7), fake.Number(1, 7), fake.Number(1, 7)}
		sigunguIds := []int{fake.Number(1, 40), fake.Number(1, 40), fake.Number(1, 40)}
		bulk[i] = t.db.Teacher.Create().
			SetUserID(ids[i]).
			SetProfileImageUrl(fake.URL()).
			SetName(fake.Name()[:8]).
			SetAge(fake.Number(20, 50)).
			SetIntroduce(fake.Snack()).
			AddYogaIDs(yogaIds...).
			AddSigunguIDs(sigunguIds...)
		t.db.User.UpdateOneID(ids[i]).SetType(model.TeacherType).Exec(context.Background())
	}

	err := t.db.Teacher.CreateBulk(bulk...).Exec(context.Background())
	fmt.Println(err)
}

func (t *seeding) Recruitment() {
	bulk := make([]*ent.RecruitmentCreate, 30)
	j := 1
	for i := 0; i < 30; i++ {
		if i%2 == 0 {
			j++
		}
		bulk[i] = t.db.Recruitment.Create().
			SetWriterID(j)
	}

	err := t.db.Recruitment.CreateBulk(bulk...).Exec(context.Background())
	fmt.Println(err)
	rid := 1
	for k := 0; k < 60; k++ {
		if k%2 == 0 {
			rid++
		}
		schedules := make([]*model.Schedule, 0)
		for i := 0; i < 3; i++ {
			schedules = append(schedules, &model.Schedule{
				StartDateTime: transport.TimeString(fake.DateRange(time.Now().AddDate(0, -4, 0), time.Now().AddDate(0, -3, 0))),
				EndDateTime:   transport.TimeString(fake.DateRange(time.Now(), time.Now().AddDate(0, 3, 0))),
			})
		}

		t.db.RecruitmentInstead.Create().SetRecruitmentID(rid).
			SetMinCareer(strconv.Itoa(fake.Number(1, 20))).
			SetPay(strconv.Itoa(fake.Number(10000, 50000))).
			SetSchedule(schedules).Exec(context.Background())

	}
}

func (t *seeding) YogaGroup() {
	bulk := make([]*ent.YogaGroupCreate, 4)
	bulk[0] = t.db.YogaGroup.Create().SetCategory("????????????").SetCategoryEng("ashtanga").SetDescription("???????????? ???????????????.")
	bulk[1] = t.db.YogaGroup.Create().SetCategory("??????").SetCategoryEng("hatha").SetDescription("?????? ???????????????.")
	bulk[2] = t.db.YogaGroup.Create().SetCategory("?????????").SetCategoryEng("vinyasa").SetDescription("????????? ???????????????.")
	bulk[3] = t.db.YogaGroup.Create().SetCategory("????????????").SetCategoryEng("adidas").SetDescription("???????????? ???????????????.")

	err := t.db.YogaGroup.CreateBulk(bulk...).Exec(context.Background())
	fmt.Println(err)
}

func (t *seeding) Yoga() {
	bulk := make([]*ent.YogaCreate, 8)
	bulk[0] = t.db.Yoga.Create().SetNameKor("???????????? ??????").SetLevel(5).SetYogaGroupID(1)
	bulk[1] = t.db.Yoga.Create().SetNameKor("???????????? ???????????????").SetLevel(4).SetYogaGroupID(1)
	bulk[2] = t.db.Yoga.Create().SetNameKor("???????????? ??????????????? ??????").SetLevel(3).SetYogaGroupID(1)
	bulk[3] = t.db.Yoga.Create().SetNameKor("?????? ?????????").SetLevel(3).SetYogaGroupID(2)
	bulk[4] = t.db.Yoga.Create().SetNameKor("?????? ?????????").SetLevel(2).SetYogaGroupID(2)
	bulk[5] = t.db.Yoga.Create().SetNameKor("?????? ??????").SetLevel(2).SetYogaGroupID(2)
	bulk[6] = t.db.Yoga.Create().SetNameKor("????????? ?????????").SetLevel(2).SetYogaGroupID(3)
	bulk[7] = t.db.Yoga.Create().SetNameKor("????????? ??????").SetLevel(1).SetYogaGroupID(3)
	err := t.db.Yoga.CreateBulk(bulk...).Exec(context.Background())
	fmt.Println(err)
}
