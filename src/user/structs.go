package user

type UserBase struct {
	Id         int    `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	Avatar     string `json:"avatar,omitempty"`
	OpenId     string `json:"open_id,omitempty"`
	SessionKey string `json:"session_key,omitempty"`
}

type UserReq struct {
	UserBase
	Code string `json:"code,omitempty"`
}

type UserResp struct {
	UserBase
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	IsNewUser bool   `json:"is_new_user"`
}
