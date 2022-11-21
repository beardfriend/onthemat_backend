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
	Name           string  `json:"name" validate:"required"`
	LogoUrl        string  `json:"logoUrl" validate:"required,urlStartHttpHttps"`
	BusinessCode   string  `json:"businessCode" validate:"required"`
	CallNumber     string  `json:"callNumber" validate:"required,phoneNumNoDash"`
	AddressRoad    string  `json:"addressRoad" validate:"required"`
	AddressDetail  *string `json:"addressDetail"`
	AddressSigungu string  `json:"addressSigungu" validate:"required"`
}

type AcademyUpdateRequestBody struct {
	Name           string  `json:"name" validate:"required"`
	LogoUrl        string  `json:"logoUrl" validate:"required,urlStartHttpHttps"`
	CallNumber     string  `json:"callNumber"`
	AddressRoad    string  `json:"addressRoad" validate:"required"`
	AddressDetail  *string `json:"addressDetail"`
	AddressSigungu string  `json:"addressSigungu" validate:"required"`
}

// ------------------- Response -------------------
