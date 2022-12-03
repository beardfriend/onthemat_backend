package model

import (
	"onthemat/internal/app/transport"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// ------------------- Default  -------------------

type DefaultTimeMixin struct {
	mixin.Schema
}

func (DefaultTimeMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Time("createdAt").
			Immutable().
			SchemaType(
				map[string]string{
					dialect.Postgres: "timestamp",
				},
			).
			GoType(transport.TimeString{}).
			Default(transport.TimeString{}.Now),

		field.Time("updatedAt").
			GoType(transport.TimeString{}).
			SchemaType(
				map[string]string{
					dialect.Postgres: "timestamp",
				},
			).
			Default(transport.TimeString{}.Now).
			UpdateDefault(transport.TimeString{}.Now),
	}
}

// ------------------- Delete -------------------

type WithDeletedTimeMixin struct {
	mixin.Schema
}

func (WithDeletedTimeMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Time("createdAt").
			Immutable().
			SchemaType(
				map[string]string{
					dialect.Postgres: "timestamp",
				},
			).
			GoType(transport.TimeString{}).
			Default(transport.TimeString{}.Now),

		field.Time("updatedAt").
			GoType(transport.TimeString{}).
			SchemaType(
				map[string]string{
					dialect.Postgres: "timestamp",
				},
			).
			Default(transport.TimeString{}.Now).
			UpdateDefault(transport.TimeString{}.Now),

		field.Time("deletedAt").
			GoType(transport.TimeString{}).
			SchemaType(
				map[string]string{
					dialect.Postgres: "timestamp",
				},
			).
			Nillable().
			Optional(),
	}
}

type (
	UserType int8

	SocialType int8
)

var (
	TeacherString    string     = "teacher"
	AcademyString    string     = "academy"
	SuperAdminString string     = "superAdmin"
	KakaoString      string     = "kakao"
	GoogleString     string     = "google"
	NaverString      string     = "naver"
	TeacherType      UserType   = 1
	AcademyType      UserType   = 2
	SuperAdminType   UserType   = 11
	KakaoSocialType  SocialType = 1
	GoogleSocialType SocialType = 2
	NaverSocialType  SocialType = 3
)

func (t *SocialType) ToString() *string {
	if t == nil {
		return nil
	}

	var result *string
	if *t == KakaoSocialType {
		result = &KakaoString
	} else if *t == GoogleSocialType {
		result = &GoogleString
	} else if *t == NaverSocialType {
		result = &NaverString
	}

	return result
}

func (t *UserType) ToString() *string {
	if t == nil {
		return nil
	}
	var result *string

	switch *t {
	case TeacherType:
		result = &TeacherString
	case AcademyType:
		result = &AcademyString
	case SuperAdminType:
		result = &SuperAdminString
	}

	return result
}

func (t *UserType) ToUserType(v *string) *UserType {
	if v == nil {
		return nil
	}

	var result *UserType

	if v == &TeacherString {
		teacher := TeacherType
		result = &teacher

	} else if v == &AcademyString {
		academy := AcademyType
		result = &academy

	} else if v == &SuperAdminString {
		superAdmin := SuperAdminType
		result = &superAdmin
	}

	return result
}

func (t *SocialType) ToSocialType(v *string) *SocialType {
	if v == nil {
		return nil
	}

	var result *SocialType

	if v == &GoogleString {
		google := GoogleSocialType
		result = &google

	} else if v == &KakaoString {
		kakao := KakaoSocialType
		result = &kakao

	} else if v == &NaverString {
		naver := NaverSocialType
		result = &naver
	}

	return result
}
