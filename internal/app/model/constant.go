package model

type (
	UserType int8

	SocialType int8
)

var (
	TeacherString    string     = "teacher"
	AcademyString    string     = "academy"
	KakaoString      string     = "kakao"
	GoogleString     string     = "google"
	NaverString      string     = "naver"
	TeacherType      UserType   = 1
	AcademyType      UserType   = 2
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
	if *t == TeacherType {
		result = &TeacherString
	} else if *t == AcademyType {
		result = &AcademyString
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
