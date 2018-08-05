package store

type ImageStruct struct {
	FileName string `json:"filename"`
	Image    []byte `json:"image"`
}

type ImageRequest struct {
	Token  string        `json:"token"`
	Images []ImageStruct `json:"images"`
}

type User struct {
	Name string `json:"user"`
	Pass string `json:"pass"`
	Role string `json:"role,omitempty"`
}

type Users []User

type JwtToken struct {
	Token    string `json:"token"`
	UserName string `json:"username"`
}

type Session struct {
	ID    string `json:"_id,omitempty"`
	Name  string `json:"name"`
	Token string `json:"token"`
}

type Sessions []Session

type ImageStructResponse struct {
	FileName string `json:"file_name"`
	Image    string `json:"image_b64"`
}

type ImagesResponse []ImageStructResponse
