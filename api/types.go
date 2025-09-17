package api

type Api struct {
	AccessToken string
}

type AuthResponse struct {
	AccessToken string `json:"AccessToken"`
}
