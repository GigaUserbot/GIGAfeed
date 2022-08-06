package main

type Event struct {
	Ref string `json:"ref,omitempty"`
	// The pull request itself.
	PullRequest PullRequest `json:"pull_request,omitempty"`
	// The user who pushed the commits.
	Pusher Pusher `json:"pusher,omitempty"`
	// The action performed. Can be created, edited, or deleted.
	Action string `json:"action,omitempty"`
	// The commit comment resource.
	Comment Comment `json:"comment,omitempty"`
	// Commits
	Commits []Commit `json:"commits,omitempty"`
	// The discussion resource.
	Discussion Discussion `json:"discussion,omitempty"`
	// The issue the comment belongs to.
	Issue Issue `json:"issue,omitempty"`
	// The repository where the event occurred.
	Repository Repository `json:"repository"`
	// The user that triggered the event.
	Sender User `json:"sender"`
	// URL that shows the changes in this ref update, from the before commit to the after commit.
	// For a newly created ref that is directly based on the default branch, this is the comparison between the head of the default branch and the after commit.
	// Otherwise, this shows all commits until the after commit.
	Compare string `json:"compare,omitempty"`
}

type Commit struct {
	Message string `json:"message"`
	Url     string `json:"url"`
}

type Pusher struct {
	Name  string `json:"name"`
	Email string `json:"email,omitempty"`
}

type PullRequest struct {
	Title    string `json:"title"`
	Url      string `json:"html_url"`
	IssueUrl string `json:"issue_url,omitempty"`
	Number   int    `json:"number"`
}

type Issue struct {
	Title  string `json:"title"`
	Url    string `json:"html_url"`
	Number int    `json:"number"`
	User   User   `json:"user,omitempty"`
}

type Repository struct {
	Name    string `json:"name"`
	Private bool   `json:"private"`
	Url     string `json:"html_url"`
}

type Discussion struct {
	Title  string `json:"title"`
	Url    string `json:"html_url"`
	Number int    `json:"number"`
	User   User   `json:"user,omitempty"`
}

type Comment struct {
	Url  string `json:"html_url,omitempty"`
	User User   `json:"user,omitempty"`
}

type User struct {
	Name string `json:"login,omitempty"`
	Url  string `json:"html_url,omitempty"`
	Type string `json:"type,omitempty"`
}
