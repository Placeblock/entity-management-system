package rest

type Response struct {
	Title  string `json:"title"`
	Detail string `json:"detail"`
	Status int    `json:"-"`
}
