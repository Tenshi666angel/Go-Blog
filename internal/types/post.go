package types

type PostEntity struct {
	Id         int64
	Title      string
	Content    string
	UserId     int64
	LikesCount int
	AppID      string
}

type PostResponse struct {
	Title      string
	Content    string
	Author     string
	LikesCount int
	AppID      string
}

type PostRequest struct {
	Title   string
	Content string
}

type Like struct {
	PostID   string
	Username string
}
