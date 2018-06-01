package main

// Forums response data structure
type Forums []struct {
	ID                string `json:"id"`
	Alias             string `json:"alias"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	SubscriptionCount int    `json:"subscriptionCount"`
	Subscribed        bool   `json:"subscribed"`
	Read              bool   `json:"read"`
	CreatedAt         string `json:"createdAt"`
	UpdatedAt         string `json:"updatedAt"`
	CanPost           bool   `json:"canPost"`
	IgnorePost        bool   `json:"ignorePost"`
	Invisible         bool   `json:"invisible"`
	IsSchool          bool   `json:"isSchool"`
	FullyAnonymous    bool   `json:"fullyAnonymous"`
	CanUseNickname    bool   `json:"canUseNickname"`
	PostThumbnail     struct {
		Size string `json:"size"`
	} `json:"postThumbnail"`
	ShouldCategorized bool     `json:"shouldCategorized"`
	TitlePlaceholder  string   `json:"titlePlaceholder"`
	Subcategories     []string `json:"subcategories"`
}

// PostMeta single post metadata
type PostMeta struct {
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

// Post response data structure
type Post struct {
	ID                   int      `json:"id"`
	Title                string   `json:"title"`
	Content              string   `json:"content"`
	AnonymousSchool      bool     `json:"anonymousSchool"`
	AnonymousDepartment  bool     `json:"anonymousDepartment"`
	Pinned               bool     `json:"pinned"`
	ForumID              string   `json:"forumId"`
	ReplyID              int      `json:"replyId"`
	CreatedAt            string   `json:"createdAt"`
	UpdatedAt            string   `json:"updatedAt"`
	CommentCount         int      `json:"commentCount"`
	LikeCount            int      `json:"likeCount"`
	Tags                 []string `json:"tags"`
	WithNickname         bool     `json:"withNickname"`
	ReportReason         string   `json:"reportReason"`
	HiddenByAuthor       bool     `json:"hiddenByAuthor"`
	ForumName            string   `json:"forumName"`
	ForumAlias           string   `json:"forumAlias"`
	Gender               string   `json:"gender"`
	ReplyTitle           string   `json:"replyTitle"`
	PersonaSubscriptable bool     `json:"personaSubscriptable"`
	Hidden               bool     `json:"hidden"`
	Media                []struct {
		URL string `json:"url"`
	} `json:"media"`
}

// Comment response data structure
type Comments []struct {
	ID             string `json:"id"`
	Anonymous      bool   `json:"anonymous"`
	PostID         int    `json:"postId"`
	CreatedAt      string `json:"createdAt"`
	UpdatedAt      string `json:"updatedAt"`
	Floor          int    `json:"floor"`
	Content        string `json:"content"`
	LikeCount      int    `json:"likeCount"`
	WithNickname   bool   `json:"withNickname"`
	HiddenByAuthor bool   `json:"hiddenByAuthor"`
	Gender         string `json:"gender"`
	School         string `json:"school"`
	Department     string `json:"department"`
	Host           bool   `json:"host"`
	ReportReason   string `json:"reportReason"`
	Hidden         bool   `json:"hidden"`
	InReview       bool   `json:"inReview"`
}
