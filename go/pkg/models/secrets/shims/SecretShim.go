package shims

type SecretShim interface {
	GetSecret() (ClientCredentials, error)

	SetSecret(credentials ClientCredentials) error

	RemoveSecret() error
}

type ClientCredentials struct {
	Password string `json:"password"`
	Email    string `json:"email"`
	Url      string `json:"url"`
	AnonKey  string `json:"anon_key"`
}
