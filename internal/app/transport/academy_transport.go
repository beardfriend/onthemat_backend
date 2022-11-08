package transport

// ------------------- Request -------------------

// ___________ Body ___________

type AcademyCreateRequestBody struct {
	Name          string `json:"name"`
	LogoUrl       string `json:"logoUrl"`
	BusinessCode  string `json:"businessCode"`
	CallNumber    string `json:"callNumber"`
	AddressRoad   string `json:"addressRoad"`
	AddressSigun  string `json:"addressSiGun"`
	AddressGu     string `json:"addressGu"`
	AddressDong   string `json:"addressDong"`
	AddressDetail string `json:"addressDetail"`
	AddressX      string `json:"addressX"`
	AddressY      string `json:"addressY"`
}

// ------------------- Response -------------------

type AcademyDetailRepsonse struct {
	ID            int    `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	Email         string `json:"email,omitempty"`
	BusinessCode  string `json:"business_code,omitempty"`
	CallNumber    string `json:"call_number,omitempty"`
	AddressRoad   string `json:"address_road,omitempty"`
	AddressDetail string `json:"address_detail,omitempty"`
	AddressSigun  string `json:"address_sigun,omitempty"`
	AddressX      string `json:"address_x,omitempty"`
	AddressY      string `json:"address_y,omitempty"`
}
