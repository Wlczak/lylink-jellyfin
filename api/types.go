package api

type Api struct {
	Username    string
	AccessToken string
}

type AuthResponse struct {
	AccessToken string `json:"AccessToken"`
}
