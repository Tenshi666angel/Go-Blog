package types

type PostEntity struct {
	Id          int64
	Title       string
	Content     string
	User_id     int64
	Likes_count int
}

type PostResponse struct {
	Title       string
	Content     string
	Author      string
	Likes_count int
}

type PostRequest struct {
	Title   string
	Content string
}
