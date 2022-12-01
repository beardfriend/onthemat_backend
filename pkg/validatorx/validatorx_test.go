package validatorx_test

import (
	"testing"

	"onthemat/internal/app/utils"
	"onthemat/pkg/validatorx"

	"github.com/stretchr/testify/assert"
)

type TestMustFieldStruct struct {
	Id          *int    `json:"id"`
	AgencyName  *string `json:"agencyName" validate:"must"`
	ImageUrl    *string `json:"imageUrl"`
	Description *string `json:"description" validate:"must"`
}

func TestMustFieldCheckValidation(t *testing.T) {
	validator := validatorx.NewValidatorx().AddCheckMustFieldIfIdFieldExistValidation("must").SetExtractTagName().Init()
	st := &TestMustFieldStruct{
		Id:         utils.Int(1),
		AgencyName: utils.String("name"),
		ImageUrl:   nil,
	}

	validator.ValidateStruct(st)
}

type TestRegistStruct struct {
	Password string `json:"password_aa" validate:"cusPassword,min=5,max=20"`
	Url      string `json:"url" validate:"cusUrl"`
	Phone    string `query:"aa_asd" validate:"cusPhone"`
}

func TestRegistValidation(t *testing.T) {
	validator := validatorx.NewValidatorx().
		AddPasswordAtLeastOneCharNumValidation("cusPassword").
		AddPhoneNumValidation("cusPhone").
		AddUrlValidation("cusUrl").
		SetExtractTagName().
		Init()

	t.Run("모두 다 정상적으로 Validation 되는 지", func(t *testing.T) {
		s := TestRegistStruct{Password: "awesome1", Url: "http://www.naver.com", Phone: "01043226633"}
		err := validator.ValidateStruct(s)

		assert.Equal(t, len(err), 0)
	})

	t.Run("모두 다 에러일 경우", func(t *testing.T) {
		s := TestRegistStruct{Password: "awesome1#./", Url: "httpaaa://www.naver.com", Phone: "043226633"}
		err := validator.ValidateStruct(s)
		assert.Equal(t, len(err), 3)
	})
}

type TestPasswordStruct struct {
	Password string `validate:"cusPassword,min=5,max=20"`
}

// 1자 이상의 문자열(특수문자)과 숫자를 포함한 5자 20이하의 패스워드 ... 가능한 특수문자는 !@$%^&*
func TestPassword(t *testing.T) {
	validator := validatorx.NewValidatorx().AddPasswordAtLeastOneCharNumValidation("cusPassword").Init()

	t.Run("정상 입력된 경우", func(t *testing.T) {
		t.Run("소문자 + 숫자", func(t *testing.T) {
			s := TestPasswordStruct{Password: "awesome1"}
			if err := validator.ValidateStruct(s); err != nil {
				t.Error("에러")
			}
		})

		t.Run("대소문자 + 숫자", func(t *testing.T) {
			s := TestPasswordStruct{Password: "Awesome1"}
			if err := validator.ValidateStruct(s); err != nil {
				t.Error("에러")
			}
		})

		t.Run("소문자 + 숫자 + 특수문자", func(t *testing.T) {
			s := TestPasswordStruct{Password: "awesome1!"}
			if err := validator.ValidateStruct(s); err != nil {
				t.Error("에러")
			}
		})

		t.Run("대소문자 + 숫자 + 특수문자", func(t *testing.T) {
			s := TestPasswordStruct{Password: "Awesomeasfsafwq1!"}
			if err := validator.ValidateStruct(s); err != nil {
				t.Error("에러")
			}
		})

		t.Run("대소문자 + 숫자 + 특수문자", func(t *testing.T) {
			s := TestPasswordStruct{Password: "!@Awesome1!"}
			if err := validator.ValidateStruct(s); err != nil {
				t.Error("에러")
			}
		})
	})

	t.Run("잘못 입력된 경우", func(t *testing.T) {
		t.Run("5자 미만", func(t *testing.T) {
			s := TestPasswordStruct{Password: "asd1"}
			if err := validator.ValidateStruct(s); err == nil {
				t.Error("에러")
			}
		})

		t.Run("20자 이상", func(t *testing.T) {
			s := TestPasswordStruct{Password: "asdasfasfsafaffaffaffaffasfasfa123f"}
			if err := validator.ValidateStruct(s); err == nil {
				t.Error("에러")
			}
		})

		t.Run("한글", func(t *testing.T) {
			s := TestPasswordStruct{Password: "ㅁㄴㄹㄹㄴ"}
			if err := validator.ValidateStruct(s); err == nil {
				t.Error("에러")
			}
		})

		t.Run("소문자만", func(t *testing.T) {
			s := TestPasswordStruct{Password: "asdff"}
			if err := validator.ValidateStruct(s); err == nil {
				t.Error("에러")
			}
		})

		t.Run("대문자만", func(t *testing.T) {
			s := TestPasswordStruct{Password: "ASDFF"}
			if err := validator.ValidateStruct(s); err == nil {
				t.Error("에러")
			}
		})

		t.Run("숫자만", func(t *testing.T) {
			s := TestPasswordStruct{Password: "12345"}
			if err := validator.ValidateStruct(s); err == nil {
				t.Error("에러")
			}
		})

		t.Run("금지된 특수문자", func(t *testing.T) {
			s := TestPasswordStruct{Password: "asdf123#"}
			if err := validator.ValidateStruct(s); err == nil {
				t.Error("에러")
			}
		})
	})
}

type TestUrlStruct struct {
	Url string `validate:"cusUrl"`
}

func TestUrl(t *testing.T) {
	validator := validatorx.NewValidatorx().AddUrlValidation("cusUrl").Init()

	t.Run("정상 입려된 경우", func(t *testing.T) {
		t.Run("http로 시작", func(t *testing.T) {
			s := TestUrlStruct{Url: "http://www.naver.com"}
			if err := validator.ValidateStruct(s); err != nil {
				t.Error("에러")
			}
		})

		t.Run("https로 시작", func(t *testing.T) {
			s := TestUrlStruct{Url: "https://www.naver.com"}
			if err := validator.ValidateStruct(s); err != nil {
				t.Error("에러")
			}
		})

		t.Run("www 없는 경우", func(t *testing.T) {
			s := TestUrlStruct{Url: "https://github.io"}
			if err := validator.ValidateStruct(s); err != nil {
				t.Error("에러")
			}
		})

		t.Run("쿼리스트링 및 하위경로", func(t *testing.T) {
			s := TestUrlStruct{Url: "https://www.naver.com/adsd?query=asd&ff=qwe"}
			if err := validator.ValidateStruct(s); err != nil {
				t.Error("에러")
			}
		})
		t.Run("쿼리스트링 및 하위경로", func(t *testing.T) {
			s := TestUrlStruct{Url: "https://www.naver.com/adsd/asd?ffb=asd"}
			if err := validator.ValidateStruct(s); err != nil {
				t.Error("에러")
			}
		})

		t.Run("#이 있을 때", func(t *testing.T) {
			s := TestUrlStruct{Url: "https://www.naver.com/adsd/asd#태그"}
			if err := validator.ValidateStruct(s); err != nil {
				t.Error("에러")
			}
		})
	})

	t.Run("잘못 입력된 경우", func(t *testing.T) {
		t.Run("/ 가 하나 없을 때", func(t *testing.T) {
			s := TestUrlStruct{Url: "http:/www.naver.com"}
			if err := validator.ValidateStruct(s); err == nil {
				t.Error("에러")
			}
		})

		t.Run("/ 가 두 개  없을 때", func(t *testing.T) {
			s := TestUrlStruct{Url: "http:www.naver.com"}
			if err := validator.ValidateStruct(s); err == nil {
				t.Error("에러")
			}
		})

		t.Run(": 이 없을 때", func(t *testing.T) {
			s := TestUrlStruct{Url: "httpwww.naver.com"}
			if err := validator.ValidateStruct(s); err == nil {
				t.Error("에러")
			}
		})

		t.Run("http or https 로 시작하지 않을 때", func(t *testing.T) {
			s := TestUrlStruct{Url: "htts://www.naver.com"}
			if err := validator.ValidateStruct(s); err == nil {
				t.Error("에러")
			}
		})
	})
}

type TestPhoneStruct struct {
	Phone string `validate:"cusPhone"`
}

func TestPhone(t *testing.T) {
	validator := validatorx.NewValidatorx().AddPhoneNumValidation("cusPhone").Init()

	t.Run("정상 입력된 경우", func(t *testing.T) {
		t.Run("010", func(t *testing.T) {
			s := TestPhoneStruct{Phone: "01043226633"}
			if err := validator.ValidateStruct(s); err != nil {
				t.Error("에러")
			}
		})

		t.Run("011", func(t *testing.T) {
			s := TestPhoneStruct{Phone: "01143226633"}
			if err := validator.ValidateStruct(s); err != nil {
				t.Error("에러")
			}
		})

		t.Run("10자리", func(t *testing.T) {
			s := TestPhoneStruct{Phone: "0116338517"}
			if err := validator.ValidateStruct(s); err != nil {
				t.Error("에러")
			}
		})
	})
	t.Run("잘못 입력된 경우", func(t *testing.T) {
		t.Run("문자열", func(t *testing.T) {
			s := TestPhoneStruct{Phone: "0104322663a"}
			if err := validator.ValidateStruct(s); err == nil {
				t.Error("에러")
			}
		})

		t.Run("특수문자", func(t *testing.T) {
			s := TestPhoneStruct{Phone: "0114322663!"}
			if err := validator.ValidateStruct(s); err == nil {
				t.Error("에러")
			}
		})

		t.Run("13자리", func(t *testing.T) {
			s := TestPhoneStruct{Phone: "010432266333"}
			if err := validator.ValidateStruct(s); err == nil {
				t.Error("에러")
			}
		})

		t.Run("10자리", func(t *testing.T) {
			s := TestPhoneStruct{Phone: "010432266"}
			if err := validator.ValidateStruct(s); err == nil {
				t.Error("에러")
			}
		})
	})
}
