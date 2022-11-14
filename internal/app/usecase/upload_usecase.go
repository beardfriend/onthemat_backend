package usecase

import (
	"context"
	"crypto/sha256"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"time"

	ex "onthemat/internal/app/common"
	"onthemat/internal/app/repository"

	"onthemat/internal/app/transport"
	"onthemat/pkg/aws"
	"onthemat/pkg/ent"
	"onthemat/pkg/ent/image"
	"onthemat/pkg/validatorx"
)

type UploadUsecase interface {
	Upload(ctx context.Context, file *multipart.FileHeader, params *transport.UploadParams, userId int) error
}

type uploadUseCase struct {
	imageRepo repository.ImageRepository
	s3        *aws.S3
}

func NewUploadUsecase(imageRepo repository.ImageRepository, s3 *aws.S3) UploadUsecase {
	return &uploadUseCase{
		imageRepo: imageRepo,
		s3:        s3,
	}
}

func (u *uploadUseCase) Upload(ctx context.Context, file *multipart.FileHeader, params *transport.UploadParams, userId int) error {
	fileExt := filepath.Ext(file.Filename)

	isUsable, _ := validatorx.ImageExtensionValidator(fileExt)
	if !isUsable {
		return ex.NewBadRequestError("사용할 수 없는 확장자입니다.")
	}

	// make Filename
	millisecondTime := time.Now().UnixMilli()
	hashedFileName := sha256.Sum256([]byte(fmt.Sprintf("onthemat_%d_%s%s", millisecondTime, file.Filename, fileExt)))
	key := fmt.Sprintf("%s/%x%s", params.Purpose, hashedFileName, fileExt)

	fileBody, _ := file.Open()

	resp := u.s3.Upload(key, fileBody)

	u.imageRepo.Create(ctx, &ent.Image{
		Name:        key,
		Path:        resp.Location,
		Size:        int(file.Size),
		ContentType: file.Header.Get("Content-Type"),
		Type:        image.Type(params.Purpose),
	}, userId)

	return nil
}
