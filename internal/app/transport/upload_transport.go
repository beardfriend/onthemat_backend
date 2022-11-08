package transport

type UploadParams struct {
	Purpose string `params:"purpose,default=profile" validate:"required,oneof=logo profile"`
}
