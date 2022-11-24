package transport

// ------------------- Request -------------------

// ___________ Param ___________

type AcademyDetailParam struct {
	Id int `params:"id" validate:"required,numeric"`
}

type AcademyUpdateParam struct {
	Id int `params:"id" validate:"required,numeric"`
}

// ___________ Query ___________

// ___________ Body ___________

type AcademyCreateRequestBody struct {
	Info    infoForCreate `json:"info" validate:"dive"`
	YogaIDs []int         `json:"yogaIDs"`
}

type infoForCreate struct {
	Name           string  `json:"name" validate:"required"`
	LogoUrl        string  `json:"logoUrl" validate:"required,urlStartHttpHttps"`
	CallNumber     string  `json:"callNumber"`
	BusinessCode   string  `json:"businessCode" validate:"required"`
	AddressRoad    string  `json:"addressRoad" validate:"required"`
	AddressDetail  *string `json:"addressDetail"`
	AddressSigungu string  `json:"addressSigungu" validate:"required"`
}

type AcademyUpdateRequestBody struct {
	Info    info  `json:"info" validate:"dive"`
	YogaIDs []int `json:"yogaIDs"`
}

type AcademyPatchRequestBody struct {
	Info    *info  `json:"info" validate:"dive"`
	YogaIDs []*int `json:"yogaIDs" `
}

type info struct {
	Name           string  `json:"name" validate:"required"`
	LogoUrl        string  `json:"logoUrl" validate:"required,urlStartHttpHttps"`
	CallNumber     string  `json:"callNumber"`
	AddressRoad    string  `json:"addressRoad" validate:"required"`
	AddressDetail  *string `json:"addressDetail"`
	AddressSigungu string  `json:"addressSigungu" validate:"required"`
}

// ------------------- Response -------------------
