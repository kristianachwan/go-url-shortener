package shortener

type Redirect struct {
	Code      string `json:"code"`
	URL       string `json:"url" validate:"empty=false & format=url"`
	CreatedAt int64  `json:"created_at"`
	Count     int64  `json:"count"`
}
