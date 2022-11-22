package request

type YogaGroupCreateBody struct {
	Category    string `json:"category" valiate:"required"`
	CategoryEng string `json:"categoryEng"`
	Description string `json:"description"`
}

type YogaCreateBody struct {
	YogaGroupId int     `json:"yogaGroupId"`
	NameKor     string  `json:"nameKor"`
	NameEng     *string `json:"nameEng"`
	Level       *int    `json:"level"`
	Description *string `json:"string"`
}
