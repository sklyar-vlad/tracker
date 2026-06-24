package dto

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func ToTokenResponse(accessToken, refreshToken string) TokenResponse {
	resp := TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return resp
}
