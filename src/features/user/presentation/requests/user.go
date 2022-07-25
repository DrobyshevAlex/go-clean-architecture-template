package requests

type GetProfileRequest struct {
	Id int64 `uri:"id" binding:"required"`
}

type AuthRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
