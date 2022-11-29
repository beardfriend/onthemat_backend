package request

type UploadParams struct {
	Purpose string `params:"purpose" validate:"required,oneof=logo profile"`
}
