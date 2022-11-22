package http

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"onthemat/internal/app/mocks"
	"onthemat/internal/app/transport"
	"onthemat/internal/app/utils"
	"onthemat/pkg/ent"
	pkgMock "onthemat/pkg/mocks"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UserHandlerTestSuite struct {
	suite.Suite
	mockUserUseCase *mocks.UserUseCase
	mockMiddleware  *mocks.MiddleWare
	mockValidator   *pkgMock.Validator
	fiber           *fiber.App

	req *http.Request
}

// 모든 테스트 시작 전 1회
func (ts *UserHandlerTestSuite) SetupSuite() {
	ts.mockUserUseCase = new(mocks.UserUseCase)
	ts.mockMiddleware = new(mocks.MiddleWare)

	ts.fiber = fiber.New()
	ts.mockValidator = new(pkgMock.Validator)
	NewUserHandler(ts.mockMiddleware, ts.mockUserUseCase, ts.fiber)
}

// Body가 없거나 QueryString 혹은 Param이 고정인 경우 여기서 공통으로 사용.
func (ts *UserHandlerTestSuite) BeforeTest(suiteName, testName string) {
	switch testName {
	case "TestGetMe":
		ts.req = httptest.NewRequest(fiber.MethodGet, "/user/me", nil)
	case "TestAddYoga":
		ts.req = httptest.NewRequest(fiber.MethodPost, "/user/yoga", nil)

	}
}

func (ts *UserHandlerTestSuite) TearDownTest() {
	ts.req = nil
}

// ------------------- Test Case -------------------

func (ts *UserHandlerTestSuite) TestGetMe() {
	ts.Run("success", func() {
		ts.mockMiddleware.On("Auth", mock.Anything).Return(nil).Once()
		ts.mockUserUseCase.On("GetMe", mock.Anything, mock.Anything).
			Return(&ent.User{
				ID:              1,
				IsEmailVerified: true,
			}, nil).
			Once()

		resp, _ := ts.fiber.Test(ts.req)
		result := utils.MakeRespWithDataForTest[transport.UserMeResponse](resp.Body)
		fmt.Println(result)
		ts.Equal(resp.StatusCode, http.StatusOK)
	})
}

func TestUserHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(UserHandlerTestSuite))
}
