package transport

// ------------------- Request -------------------

// ___________ Body ___________

type UserYogaUpdateBody struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
