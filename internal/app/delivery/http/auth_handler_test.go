package http

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"onthemat/internal/app/common"
	"onthemat/internal/app/mocks"
	"onthemat/internal/app/usecase"
	"onthemat/internal/app/utils"
	pkgMock "onthemat/pkg/mocks"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type AuthHDTestSuite struct {
	suite.Suite
	mockAuthUseCase *mocks.AuthUseCase
	mockValidator   *pkgMock.Validator
	fiber           *fiber.App

	req *http.Request
}

// 모든 테스트 시작 전 1회
func (ts *AuthHDTestSuite) SetupSuite() {
	ts.mockAuthUseCase = new(mocks.AuthUseCase)

	ts.fiber = fiber.New()
	ts.mockValidator = new(pkgMock.Validator)
	NewAuthHandler(ts.mockAuthUseCase, ts.mockValidator, ts.fiber)
}

// Body가 없거나 QueryString 혹은 Param이 고정인 경우 여기서 공통으로 사용.
func (ts *AuthHDTestSuite) BeforeTest(suiteName, testName string) {
	switch testName {
	case "TestKakao":
		ts.req = httptest.NewRequest(fiber.MethodGet, "/auth/kakao", nil)
	case "TestGoogle":
		ts.req = httptest.NewRequest(fiber.MethodGet, "/auth/google", nil)
	case "TestNaver":
		ts.req = httptest.NewRequest(fiber.MethodGet, "/auth/naver", nil)
	case "TestKakaoCallback":
		ts.req = httptest.NewRequest(fiber.MethodGet, "/auth/kakao/callback", nil)
		q := ts.req.URL.Query()
		q.Add("code", "kakaocallbackQuery")
		ts.req.URL.RawQuery = q.Encode()

	}
}

func (ts *AuthHDTestSuite) TearDownTest() {
	ts.req = nil
}

// ------------------- Test Case -------------------

func (ts *AuthHDTestSuite) TestKakao() {
	ts.Run("success", func() {
		ts.mockAuthUseCase.On("KakaoRedirectUrl", mock.Anything).Return("http://www.kakao.com").Once()

		// http.Response
		resp, _ := ts.fiber.Test(ts.req)
		redirectUrl := resp.Header.Get("Location")
		ts.Equal(redirectUrl, "http://www.kakao.com")
		ts.Equal(resp.StatusCode, fiber.StatusFound)
	})
}

func (ts *AuthHDTestSuite) TestGoogle() {
	ts.Run("success", func() {
		ts.mockAuthUseCase.On("GoogleRedirectUrl", mock.Anything).Return("http://www.google.com").Once()

		// http.Response
		resp, _ := ts.fiber.Test(ts.req)
		redirectUrl := resp.Header.Get("Location")
		ts.Equal(redirectUrl, "http://www.google.com")
		ts.Equal(resp.StatusCode, fiber.StatusFound)
	})
}

func (ts *AuthHDTestSuite) TestNaver() {
	ts.Run("success", func() {
		ts.mockAuthUseCase.On("NaverRedirectUrl", mock.Anything).Return("http://www.naver.com").Once()

		// http.Response
		resp, _ := ts.fiber.Test(ts.req)
		redirectUrl := resp.Header.Get("Location")
		ts.Equal(redirectUrl, "http://www.naver.com")
		ts.Equal(resp.StatusCode, fiber.StatusFound)
	})
}

func (ts *AuthHDTestSuite) TestKakaoCallback() {
	ts.Run("success", func() {
		// mock
		ts.mockAuthUseCase.On("SocialLogin",
			mock.Anything,
			mock.AnythingOfType("model.SocialType"),
			mock.AnythingOfType("string"),
		).Return(&usecase.LoginResult{
			AccessToken: "accessToken",
		}, nil).Once()

		// http.Response
		resp, _ := ts.fiber.Test(ts.req)
		result := utils.MakeRespWithDataForTest[usecase.LoginResult](resp.Body)

		ts.Equal(resp.StatusCode, http.StatusOK)
		ts.Equal(result.Result.AccessToken, "accessToken")
	})

	ts.Run("일시적인 에러", func() {
		ts.mockAuthUseCase.On("SocialLogin",
			mock.Anything,
			mock.AnythingOfType("model.SocialType"),
			mock.AnythingOfType("string"),
		).Return(nil, errors.New("something")).Once()

		// http.Response
		resp, _ := ts.fiber.Test(ts.req)

		result := utils.MakeErrorForTests(resp.Body)

		ts.Equal(resp.StatusCode, http.StatusInternalServerError)
		ts.Equal(result.ErrCode, http.StatusInternalServerError)
	})
}

func (ts *AuthHDTestSuite) TestSignUp() {
	ts.Run("Body를 보내지 않았을 때", func() {
		req := httptest.NewRequest(fiber.MethodPost, "/auth/signup", nil)
		req.Header.Set("Content-Type", "application/json")
		resp, _ := ts.fiber.Test(req, -1)
		respErr := utils.MakeErrorForTests(resp.Body)
		ts.Equal(resp.StatusCode, http.StatusBadRequest)
		ts.Equal(respErr.ErrCode, common.ErrJsonMissing)
	})
}

func TestAuthHDTestSuite(t *testing.T) {
	suite.Run(t, new(AuthHDTestSuite))
}
