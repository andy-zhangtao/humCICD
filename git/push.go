/*
 * Copyright (c) 2018.
 * andy-zhangtao <ztao8607@gmail.com>
 */

package git

//Write by zhangtao<ztao8607@gmail.com> . In 2018/3/13.

type GitHubPush_repository struct {
	Id                int                         `json:"id"`
	Name              string                      `json:"name"`
	Full_name         string                      `json:"full_name"`
	Owner             GitHubPush_repository_owner `json:"owner"`
	Private           bool                        `json:"private"`
	Html_url          string                      `json:"html_url"`
	Fork              bool                        `json:"fork"`
	Url               string                      `json:"url"`
	Forks_url         string                      `json:"forks_url"`
	Keys_url          string                      `json:"keys_url"`
	Collaborators_url string                      `json:"collaborators_url"`
	Teams_url         string                      `json:"teams_url"`
	Hooks_url         string                      `json:"hooks_url"`
	Issue_events_url  string                      `json:"issue_events_url"`
	Events_url        string                      `json:"events_url"`
	Assignees_url     string                      `json:"assignees_url"`
	Branches_url      string                      `json:"branches_url"`
	Tags_url          string                      `json:"tags_url"`
	Blobs_url         string                      `json:"blobs_url"`
	Git_tags_url      string                      `json:"git_tags_url"`
	Git_refs_url      string                      `json:"git_refs_url"`
	Trees_url         string                      `json:"trees_url"`
	Statuses_url      string                      `json:"statuses_url"`
	Languages_url     string                      `json:"languages_url"`
	Stargazers_url    string                      `json:"stargazers_url"`
	Contributors_url  string                      `json:"contributors_url"`
	Subscribers_url   string                      `json:"subscribers_url"`
	Subscription_url  string                      `json:"subscription_url"`
	Commits_url       string                      `json:"commits_url"`
	Git_commits_url   string                      `json:"git_commits_url"`
	Comments_url      string                      `json:"comments_url"`
	Issue_comment_url string                      `json:"issue_comment_url"`
	Contents_url      string                      `json:"contents_url"`
	Compare_url       string                      `json:"compare_url"`
	Merges_url        string                      `json:"merges_url"`
	Archive_url       string                      `json:"archive_url"`
	Downloads_url     string                      `json:"downloads_url"`
	Issues_url        string                      `json:"issues_url"`
	Pulls_url         string                      `json:"pulls_url"`
	Milestones_url    string                      `json:"milestones_url"`
	Notifications_url string                      `json:"notifications_url"`
	Labels_url        string                      `json:"labels_url"`
	Releases_url      string                      `json:"releases_url"`
	Deployments_url   string                      `json:"deployments_url"`
	//Created_at        uint32                      `json:"created_at"`
	Updated_at string `json:"updated_at"`
	//Pushed_at         uint32                      `json:"pushed_at"`
	Git_url           string `json:"git_url"`
	Ssh_url           string `json:"ssh_url"`
	Clone_url         string `json:"clone_url"`
	Svn_url           string `json:"svn_url"`
	Size              int    `json:"size"`
	Stargazers_count  int    `json:"stargazers_count"`
	Watchers_count    int    `json:"watchers_count"`
	Has_issues        bool   `json:"has_issues"`
	Has_projects      bool   `json:"has_projects"`
	Has_downloads     bool   `json:"has_downloads"`
	Has_wiki          bool   `json:"has_wiki"`
	Has_pages         bool   `json:"has_pages"`
	Forks_count       int    `json:"forks_count"`
	Archived          bool   `json:"archived"`
	Open_issues_count int    `json:"open_issues_count"`
	Forks             int    `json:"forks"`
	Open_issues       int    `json:"open_issues"`
	Watchers          int    `json:"watchers"`
	Default_branch    string `json:"default_branch"`
	Stargazers        int    `json:"stargazers"`
	Master_branch     string `json:"master_branch"`
}
type GitHubPush_pusher struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
type GitHubPush_head_commit_committer struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
}
type GitHubPush_head_commit_removed struct {
}
type GitHubPush_head_commit struct {
	Id        string                           `json:"id"`
	Tree_id   string                           `json:"tree_id"`
	Distinct  bool                             `json:"distinct"`
	Message   string                           `json:"message"`
	Timestamp string                           `json:"timestamp"`
	Url       string                           `json:"url"`
	Author    GitHubPush_head_commit_author    `json:"author"`
	Committer GitHubPush_head_commit_committer `json:"committer"`
	Added     []string                         `json:"added"`
	Removed   []string                         `json:"removed"`
	Modified  []string                         `json:"modified"`
}
type GitHubPush_head_commit_modified struct {
}
type GitHubPush_repository_owner struct {
	Name                string `json:"name"`
	Email               string `json:"email"`
	Login               string `json:"login"`
	Id                  int    `json:"id"`
	Avatar_url          string `json:"avatar_url"`
	Gravatar_id         string `json:"gravatar_id"`
	Url                 string `json:"url"`
	Html_url            string `json:"html_url"`
	Followers_url       string `json:"followers_url"`
	Following_url       string `json:"following_url"`
	Gists_url           string `json:"gists_url"`
	Starred_url         string `json:"starred_url"`
	Subscriptions_url   string `json:"subscriptions_url"`
	Organizations_url   string `json:"organizations_url"`
	Repos_url           string `json:"repos_url"`
	Events_url          string `json:"events_url"`
	Received_events_url string `json:"received_events_url"`
	Type                string `json:"type"`
	Site_admin          bool   `json:"site_admin"`
}
type GitHubPush_sender struct {
	Login               string `json:"login"`
	Id                  int    `json:"id"`
	Avatar_url          string `json:"avatar_url"`
	Gravatar_id         string `json:"gravatar_id"`
	Url                 string `json:"url"`
	Html_url            string `json:"html_url"`
	Followers_url       string `json:"followers_url"`
	Following_url       string `json:"following_url"`
	Gists_url           string `json:"gists_url"`
	Starred_url         string `json:"starred_url"`
	Subscriptions_url   string `json:"subscriptions_url"`
	Organizations_url   string `json:"organizations_url"`
	Repos_url           string `json:"repos_url"`
	Events_url          string `json:"events_url"`
	Received_events_url string `json:"received_events_url"`
	Type                string `json:"type"`
	Site_admin          bool   `json:"site_admin"`
}
type GitHubPush struct {
	Ref         string                 `json:"ref"`
	Before      string                 `json:"before"`
	After       string                 `json:"after"`
	Created     bool                   `json:"created"`
	Deleted     bool                   `json:"deleted"`
	Forced      bool                   `json:"forced"`
	Compare     string                 `json:"compare"`
	Commits     []GitHubPush_commits   `json:"commits"`
	Head_commit GitHubPush_head_commit `json:"head_commit"`
	Repository  GitHubPush_repository  `json:"repository"`
	Pusher      GitHubPush_pusher      `json:"pusher"`
	Sender      GitHubPush_sender      `json:"sender"`
}
type GitHubPush_commits struct {
	Id        string `json:"id"`
	Tree_id   string `json:"tree_id"`
	Distinct  bool   `json:"distinct"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
	Url       string `json:"url"`
}
type GitHubPush_head_commit_author struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
}
type GitHubPush_head_commit_added struct {
}
