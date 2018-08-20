package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	// "strconv"
	"strings"
	// "time"

	gitlab "github.com/xanzy/go-gitlab"
	"golang.org/x/oauth2"

	"github.com/mattermost/mattermost-server/mlog"
	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
)

const (
	API_ERROR_ID_NOT_CONNECTED = "not_connected"
	GITLAB_ICON_URL            = "https://gitlab.com/gitlab-com/gitlab-artwork/raw/master/logo/logo-square.png"
	GITLAB_USERNAME            = "GitLab Plugin"
)

type APIErrorResponse struct {
	ID         string `json:"id"`
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

func writeError(w http.ResponseWriter, err error) {
	writeAPIError(w, &APIErrorResponse{
		Message:    err.Error(),
		StatusCode: http.StatusInternalServerError,
	})
}

func writeAPIError(w http.ResponseWriter, err *APIErrorResponse) {
	b, _ := json.Marshal(err)
	w.WriteHeader(err.StatusCode)
	w.Write(b)
}

func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	if err := p.IsValid(); err != nil {
		http.Error(w, "This plugin is not configured.", http.StatusNotImplemented)
	}

	w.Header().Set("Content-Type", "application/json")

	switch path := r.URL.Path; path {
	// case "/webhook":
	// 	p.handleWebhook(w, r)
	case "/oauth/connect":
		p.connectUserToGitLab(w, r)
	case "/oauth/complete":
		p.completeConnectUserToGitLab(w, r)
	case "/api/v1/connected":
		p.getConnected(w, r)
	// case "/api/v1/todo":
	// 	p.postToDo(w, r)
	// case "/api/v1/reviews":
	// 	p.getReviews(w, r)
	case "/api/v1/merge_requests/created":
		p.getYourMergeRequests(w, r)
	case "/api/v1/merge_requests/assigned":
		p.getAssignedMergeRequests(w, r)
	case "/api/v1/issues/assigned":
		p.getAssignedIssues(w, r)
	// case "/api/v1/mentions":
	// 	p.getMentions(w, r)
	// case "/api/v1/unreads":
	// 	p.getUnreads(w, r)
	// case "/api/v1/settings":
	// 	p.updateSettings(w, r)
	case "/api/v1/todos":
		p.getTodos(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (p *Plugin) connectUserToGitLab(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("Mattermost-User-ID")
	if userID == "" {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}

	conf := p.getOAuthConfig()

	state := fmt.Sprintf("%v_%v", model.NewId(), userID)

	p.API.KVSet(state, []byte(state))

	url := conf.AuthCodeURL(state, oauth2.AccessTypeOffline)

	http.Redirect(w, r, url, http.StatusFound)
}

func (p *Plugin) completeConnectUserToGitLab(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	conf := p.getOAuthConfig()

	code := r.URL.Query().Get("code")
	if len(code) == 0 {
		http.Error(w, "missing authorization code", http.StatusBadRequest)
		return
	}

	state := r.URL.Query().Get("state")

	if storedState, err := p.API.KVGet(state); err != nil {
		http.Error(w, "missing stored state", http.StatusBadRequest)
		return
	} else if string(storedState) != state {
		http.Error(w, "invalid state", http.StatusBadRequest)
		return
	}

	userID := strings.Split(state, "_")[1]

	p.API.KVDelete(state)

	tok, err := conf.Exchange(ctx, code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	gitlabClient := p.gitlabConnect(*tok)
	gitlabUser, _, err := gitlabClient.Users.CurrentUser(gitlab.WithContext(ctx))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if user, err := p.API.GetUser(userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		if user.Props == nil {
			user.Props = model.StringMap{}
		}
		user.Props["gitlab_user"] = gitlabUser.Username
		_, err = p.API.UpdateUser(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	userInfo := &GitLabUserInfo{
		UserID:         userID,
		Token:          tok,
		GitLabUserId:   gitlabUser.ID,
		GitLabUsername: gitlabUser.Username,
		// LastToDoPostAt: model.GetMillis(),
		Settings: &UserSettings{
			// 	SidebarButtons: SETTING_BUTTONS_TEAM,
			// 	DailyReminder:  true,
			// 	Notifications:  true,
		},
	}

	if err := p.storeGitLabUserInfo(userInfo); err != nil {
		mlog.Error(err.Error())
		http.Error(w, "Unable to connect user to GitLab", http.StatusInternalServerError)
		return
	}

	if err := p.storeGitLabToUserIDMapping(gitlabUser.Username, userID); err != nil {
		mlog.Error(err.Error())
	}

	// Post intro post
	// message := fmt.Sprintf("#### Welcome to the Mattermost GitHub Plugin!\n You've connected your Mattermost account to [%s](%s) on GitHub. Read about the features of this plugin below:\n\n##### Daily Reminders\nThe first time you log in each day, you will get a post right here letting you know what messages you need to read and what pull requests are awaiting your review.\nTurn off reminders with `/github settings reminders off`.\n\n##### Notifications\nWhen someone mentions you, requests your review, comments on or modifies one of your pull requests/issues, or assigns you, you'll get a post here about it.\nTurn off notifications with `/github settings notifications off`.\n\n##### Sidebar Buttons\nCheck out the buttons in the left-hand sidebar of Mattermost.\n* The first button tells you how many pull requests are awaiting your review\n* The second tracks the number of unread messages you have\n* The third will refresh the numbers\n\nClick on them!\n\n##### Slash Commands\n"+strings.Replace(COMMAND_HELP, "|", "`", -1), gitlabUser.GetLogin(), gitlabUser.GetHTMLURL())
	message := fmt.Sprintf("#### Welcome to the Mattermost GitLab Plugin!\nYou've connected your Mattermost account to %s on GitLab. It presumably does things.", gitlabUser.Username) // TODO improve this message
	p.CreateBotDMPost(userID, message, "custom_gitlab_welcome")

	p.API.PublishWebSocketEvent(
		WS_EVENT_CONNECT,
		map[string]interface{}{
			"connected":        true,
			"gitlab_url":       p.GitLabURL,
			"gitlab_username":  userInfo.GitLabUsername,
			"gitlab_user_id":   userInfo.GitLabUserId,
			"gitlab_client_id": p.GitLabOAuthClientID,
		},
		&model.WebsocketBroadcast{UserId: userID},
	)

	html := `
<!DOCTYPE html>
<html>
	<head>
		<script>
			window.close();
		</script>
	</head>
	<body>
		<p>Completed connecting to GitLab.</p>
	</body>
</html>
`

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

type ConnectedResponse struct {
	Connected      bool          `json:"connected"`
	GitLabURL      string        `json:"gitlab_url"`
	GitLabUsername string        `json:"gitlab_username"`
	GitLabUserId   int           `json:"gitlab_user_id"`
	GitLabClientID string        `json:"gitlab_client_id"`
	Settings       *UserSettings `json:"settings"`
}

func (p *Plugin) getConnected(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("Mattermost-User-ID")
	if userID == "" {
		writeAPIError(w, &APIErrorResponse{ID: "", Message: "Not authorized.", StatusCode: http.StatusUnauthorized})
		return
	}

	resp := &ConnectedResponse{Connected: false, GitLabURL: p.GitLabURL}

	info, _ := p.getGitLabUserInfo(userID)
	if info != nil && info.Token != nil {
		resp.Connected = true
		resp.GitLabUsername = info.GitLabUsername
		resp.GitLabUserId = info.GitLabUserId
		resp.GitLabClientID = p.GitLabOAuthClientID
		resp.Settings = info.Settings

		// if info.Settings.DailyReminder && r.URL.Query().Get("reminder") == "true" {
		// 	lastPostAt := info.LastToDoPostAt

		// 	var timezone *time.Location
		// 	offset, _ := strconv.Atoi(r.Header.Get("X-Timezone-Offset"))
		// 	timezone = time.FixedZone("local", -60*offset)

		// 	// Post to do message if it's the next day and been more than an hour since the last post
		// 	now := model.GetMillis()
		// 	nt := time.Unix(now/1000, 0).In(timezone)
		// 	lt := time.Unix(lastPostAt/1000, 0).In(timezone)
		// 	if nt.Sub(lt).Hours() >= 1 && (nt.Day() != lt.Day() || nt.Month() != lt.Month() || nt.Year() != lt.Year()) {
		// 		p.PostToDo(info)
		// 		info.LastToDoPostAt = now
		// 		p.storeGitLabUserInfo(info)
		// 	}
		// }
	}

	b, _ := json.Marshal(resp)
	w.Write(b)
}

// func (p *Plugin) getMentions(w http.ResponseWriter, r *http.Request) {
// 	userID := r.Header.Get("Mattermost-User-ID")
// 	if userID == "" {
// 		http.Error(w, "Not authorized", http.StatusUnauthorized)
// 		return
// 	}

// 	ctx := context.Background()

// 	var githubClient *github.Client
// 	username := ""

// 	if info, err := p.getGitLabUserInfo(userID); err != nil {
// 		writeAPIError(w, err)
// 		return
// 	} else {
// 		githubClient = p.githubConnect(*info.Token)
// 		username = info.GitLabUsername
// 	}

// 	result, _, err := githubClient.Search.Issues(ctx, getMentionSearchQuery(username, p.GitLabOrg), &github.SearchOptions{})
// 	if err != nil {
// 		mlog.Error(err.Error())
// 	}

// 	resp, _ := json.Marshal(result.Issues)
// 	w.Write(resp)
// }

// func (p *Plugin) getUnreads(w http.ResponseWriter, r *http.Request) {
// 	userID := r.Header.Get("Mattermost-User-ID")
// 	if userID == "" {
// 		http.Error(w, "Not authorized", http.StatusUnauthorized)
// 		return
// 	}

// 	ctx := context.Background()

// 	var githubClient *github.Client

// 	if info, err := p.getGitLabUserInfo(userID); err != nil {
// 		writeAPIError(w, err)
// 		return
// 	} else {
// 		githubClient = p.githubConnect(*info.Token)
// 	}

// 	notifications, _, err := githubClient.Activity.ListNotifications(ctx, &github.NotificationListOptions{})
// 	if err != nil {
// 		mlog.Error(err.Error())
// 	}

// 	filteredNotifications := []*github.Notification{}
// 	for _, n := range notifications {
// 		if n.GetReason() == "subscribed" {
// 			continue
// 		}

// 		if p.checkOrg(n.GetRepository().GetOwner().GetLogin()) != nil {
// 			continue
// 		}

// 		filteredNotifications = append(filteredNotifications, n)
// 	}

// 	resp, _ := json.Marshal(filteredNotifications)
// 	w.Write(resp)
// }

func (p *Plugin) getFromGitLab(
	w http.ResponseWriter,
	r *http.Request,
	doGetFromGitLab func(*gitlab.Client) (interface{}, error),
) {
	// Get current user ID from header
	userID := r.Header.Get("Mattermost-User-ID")
	if userID == "" {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}

	// Get a GitLab client
	var gitlabClient *gitlab.Client
	if info, err := p.getGitLabUserInfo(userID); err != nil {
		writeAPIError(w, err)
		return
	} else {
		gitlabClient = p.gitlabConnect(*info.Token)
	}

	results, err := doGetFromGitLab(gitlabClient)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, _ := json.Marshal(results)
	w.Write(resp)
}

func (p *Plugin) getYourMergeRequests(w http.ResponseWriter, r *http.Request) {
	p.getFromGitLab(w, r, func(gitlabClient *gitlab.Client) (interface{}, error) {
		options := &gitlab.ListMergeRequestsOptions{
			State: stringToPointer("opened"),
			Scope: stringToPointer("created_by_me"),
		}

		results, _, err := gitlabClient.MergeRequests.ListMergeRequests(options, gitlab.WithContext(r.Context()))
		return results, err
	})
}

func (p *Plugin) getAssignedMergeRequests(w http.ResponseWriter, r *http.Request) {
	p.getFromGitLab(w, r, func(gitlabClient *gitlab.Client) (interface{}, error) {
		options := &gitlab.ListMergeRequestsOptions{
			State: stringToPointer("opened"),
			Scope: stringToPointer("assigned_to_me"),
		}

		results, _, err := gitlabClient.MergeRequests.ListMergeRequests(options, gitlab.WithContext(r.Context()))
		return results, err
	})
}

func (p *Plugin) getAssignedIssues(w http.ResponseWriter, r *http.Request) {
	p.getFromGitLab(w, r, func(gitlabClient *gitlab.Client) (interface{}, error) {
		options := &gitlab.ListIssuesOptions{
			State: stringToPointer("opened"),
			Scope: stringToPointer("assigned_to_me"),
		}

		results, _, err := gitlabClient.Issues.ListIssues(options, gitlab.WithContext(r.Context()))
		return results, err
	})
}

func (p *Plugin) getTodos(w http.ResponseWriter, r *http.Request) {
	p.getFromGitLab(w, r, func(gitlabClient *gitlab.Client) (interface{}, error) {
		options := &gitlab.ListTodosOptions{
			State: stringToPointer("pending"),
		}

		results, _, err := gitlabClient.Todos.ListTodos(options, gitlab.WithContext(r.Context()))
		return results, err
	})
}

// func (p *Plugin) postToDo(w http.ResponseWriter, r *http.Request) {
// 	userID := r.Header.Get("Mattermost-User-ID")
// 	if userID == "" {
// 		writeAPIError(w, &APIErrorResponse{ID: "", Message: "Not authorized.", StatusCode: http.StatusUnauthorized})
// 		return
// 	}

// 	var githubClient *github.Client
// 	username := ""

// 	if info, err := p.getGitLabUserInfo(userID); err != nil {
// 		writeAPIError(w, err)
// 		return
// 	} else {
// 		githubClient = p.githubConnect(*info.Token)
// 		username = info.GitLabUsername
// 	}

// 	text, err := p.GetToDo(context.Background(), username, githubClient)
// 	if err != nil {
// 		mlog.Error(err.Error())
// 		writeAPIError(w, &APIErrorResponse{ID: "", Message: "Encountered an error getting the to do items.", StatusCode: http.StatusUnauthorized})
// 		return
// 	}

// 	if err := p.CreateBotDMPost(userID, text, "custom_git_todo"); err != nil {
// 		writeAPIError(w, &APIErrorResponse{ID: "", Message: "Encountered an error posting the to do items.", StatusCode: http.StatusUnauthorized})
// 	}

// 	w.Write([]byte("{\"status\": \"OK\"}"))
// }

// func (p *Plugin) updateSettings(w http.ResponseWriter, r *http.Request) {
// 	userID := r.Header.Get("Mattermost-User-ID")
// 	if userID == "" {
// 		http.Error(w, "Not authorized", http.StatusUnauthorized)
// 		return
// 	}

// 	var settings *UserSettings
// 	json.NewDecoder(r.Body).Decode(&settings)
// 	if settings == nil {
// 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}

// 	info, err := p.getGitLabUserInfo(userID)
// 	if err != nil {
// 		writeAPIError(w, err)
// 		return
// 	}

// 	info.Settings = settings

// 	if err := p.storeGitLabUserInfo(info); err != nil {
// 		mlog.Error(err.Error())
// 		http.Error(w, "Encountered error updating settings", http.StatusInternalServerError)
// 	}

// 	resp, _ := json.Marshal(info.Settings)
// 	w.Write(resp)
// }
