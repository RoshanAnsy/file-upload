package types

type FileResponse struct {
	FileName  string `json:"fileName"`
	UUID      string `json:"uuid"`
	TempPath  string `json:"tempPath"`
	FullPath  string `json:"fullPath"`
	AccessURL string `json:"accessUrl"`
}

type User struct {
	Name     string `json:"name,omitempty"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ResetPasswordRequest struct {
	Email       string `json:"email"`
	NewPassword string `json:"newPassword"`
}

type GenerateApiKeyRequest struct {
	KeyName string `json:"keyName"`
}