package main

// Articles is mapping to forum posts
type Articles []struct {
	ID                  int      `json:"id"`
	Title               string   `json:"title"`
	Excerpt             string   `json:"excerpt"`
	AnonymousSchool     bool     `json:"anonymousSchool"`
	AnonymousDepartment bool     `json:"anonymousDepartment"`
	Pinned              bool     `json:"pinned"`
	ForumID             string   `json:"forumId"`
	ReplyID             int      `json:"replyId"`
	CreatedAt           string   `json:"createdAt"`
	UpdatedAt           string   `json:"updatedAt"`
	CommentCount        int      `json:"commentCount"`
	LikeCount           int      `json:"likeCount"`
	WithNickname        bool     `json:"withNickname"`
	Tags                []string `json:"tags"`
	ForumName           string   `json:"forumName"`
	ForumAlias          string   `json:"forumAlias"`
	Gender              string   `json:"gender"`
	ReplyTitle          string   `json:"replyTitle"`
	ReportReason        string   `json:"reportReason"`
	Hidden              bool     `json:"hidden"`
	Media               []struct {
		URL string `json:"url"`
	} `json:"media"`
}
