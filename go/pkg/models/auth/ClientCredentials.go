package auth

type ClientCredentials struct {
	Password string `json:"password"`
	Email    string `json:"email"`
	Url      string `json:"url"`
	AnonKey  string `json:"anon_key"`
}
