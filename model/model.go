package model

type User struct {
	Uid       int64  `json:"uid"`
	Username  string `json:"username"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Password  string `json:"password"`
	Token     string `json:"token"`
	Role      string `json:"role"`
}

type ResponseResult struct {
	Error  string      `json:"error"`
	Result interface{} `json:"result"`
}

type Movie struct {
	Movieid    string  `json:"movieid"`
	Title      string  `json:"title"`
	Year       string  `json:"year"`
	Rating     float64 `json:"rating"`
	Totalusers int     `json:"totalusers"`
}
