package usecase_test

import (
	"context"
	"mime/multipart"
	"testing"

	"onthemat/internal/app/mocks"
	"onthemat/internal/app/transport"
	"onthemat/internal/app/usecase"

	pkgMocks "onthemat/pkg/mocks"

	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UploadUCTestSuite struct {
	suite.Suite
	uploadUC            usecase.UploadUsecase
	mockImageRepository *mocks.ImageRepository
	mockS3              *pkgMocks.S3
}

// 모든 테스트 시작 전 1회
func (ts *UploadUCTestSuite) SetupSuite() {
	ts.mockImageRepository = new(mocks.ImageRepository)
	ts.mockS3 = new(pkgMocks.S3)

	ts.uploadUC = usecase.NewUploadUsecase(ts.mockImageRepository, ts.mockS3)
}

// ------------------- Test Case -------------------

func (ts *UploadUCTestSuite) TestUpload() {
	ts.Run("성공", func() {
		ts.mockS3.On("Upload", mock.AnythingOfType("string"), mock.AnythingOfType("*os.File")).Return(&manager.UploadOutput{}).Once()
		ts.mockImageRepository.On("Create", mock.Anything, mock.AnythingOfType("*ent.Image"), mock.AnythingOfType("int")).
			Return(nil)

		h := &multipart.FileHeader{
			Filename: "../../../pkg/aws/test_object/akmu.jpeg",
		}

		ts.uploadUC.Upload(context.Background(), h, &transport.UploadParams{}, 1)
	})
}

func TestUploadUCTestSuite(t *testing.T) {
	suite.Run(t, new(UploadUCTestSuite))
}
