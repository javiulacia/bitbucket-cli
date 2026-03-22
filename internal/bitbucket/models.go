package bitbucket

type pagedResponse[T any] struct {
	Values []T    `json:"values"`
	Next   string `json:"next"`
	Page   int    `json:"page"`
	Size   int    `json:"size"`
}

type User struct {
	AccountID   string `json:"account_id,omitempty"`
	Username    string `json:"username,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
	Nickname    string `json:"nickname,omitempty"`
	UUID        string `json:"uuid,omitempty"`
}

type BranchRef struct {
	Name string `json:"name"`
}

type LinkHref struct {
	Href string `json:"href"`
}

type RepositoryLinks struct {
	HTML  LinkHref `json:"html"`
	Clone []struct {
		Href string `json:"href"`
		Name string `json:"name"`
	} `json:"clone"`
}

type Repository struct {
	UUID        string          `json:"uuid,omitempty"`
	Name        string          `json:"name"`
	FullName    string          `json:"full_name"`
	Description string          `json:"description,omitempty"`
	IsPrivate   bool            `json:"is_private"`
	Language    string          `json:"language,omitempty"`
	SCM         string          `json:"scm,omitempty"`
	MainBranch  *BranchRef      `json:"mainbranch,omitempty"`
	Links       RepositoryLinks `json:"links,omitempty"`
}

type PullRequestBranch struct {
	Branch     BranchRef   `json:"branch"`
	Repository *Repository `json:"repository,omitempty"`
}

type PullRequestParticipant struct {
	User     User `json:"user"`
	Approved bool `json:"approved"`
}

type PullRequest struct {
	ID           int                      `json:"id"`
	Title        string                   `json:"title"`
	Description  string                   `json:"description,omitempty"`
	State        string                   `json:"state"`
	CommentCount int                      `json:"comment_count,omitempty"`
	TaskCount    int                      `json:"task_count,omitempty"`
	Author       *struct{ User User }     `json:"author,omitempty"`
	Source       *PullRequestBranch       `json:"source,omitempty"`
	Destination  *PullRequestBranch       `json:"destination,omitempty"`
	Participants []PullRequestParticipant `json:"participants,omitempty"`
	Links        struct {
		HTML LinkHref `json:"html"`
	} `json:"links,omitempty"`
}
