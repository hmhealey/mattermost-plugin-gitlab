package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	// "strings"

	"github.com/mattermost/mattermost-server/mlog"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"

	// "github.com/google/go-github/github"
	"github.com/xanzy/go-gitlab"
	"golang.org/x/oauth2"
)

const (
	GITLAB_TOKEN_KEY        = "_gitlabtoken"
	GITLAB_USERNAME_KEY     = "_gitlabusername"
	WS_EVENT_CONNECT        = "connect"
	WS_EVENT_DISCONNECT     = "disconnect"
	WS_EVENT_REFRESH        = "refresh"
	SETTING_BUTTONS_TEAM    = "team"
	SETTING_BUTTONS_CHANNEL = "channel"
	SETTING_BUTTONS_OFF     = "off"
	SETTING_NOTIFICATIONS   = "notifications"
	SETTING_REMINDERS       = "reminders"
	SETTING_ON              = "on"
	SETTING_OFF             = "off"
)

type Plugin struct {
	plugin.MattermostPlugin
	gitlabClient *gitlab.Client

	BotUserID string

	GitLabURL string
	GitLabOAuthClientID     string
	GitLabOAuthClientSecret string

	// GitLabOrg               string
	Username                string
	// WebhookSecret           string
	EncryptionKey string
}

func (p *Plugin) gitlabConnect(token oauth2.Token) *gitlab.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&token)
	tc := oauth2.NewClient(ctx, ts)

	return gitlab.NewOAuthClient(tc, token.AccessToken)
}

func (p *Plugin) OnActivate() error {
	if err := p.IsValid(); err != nil {
		return err
	}

	// p.API.RegisterCommand(getCommand())

	if user, err := p.API.GetUserByUsername(p.Username); err != nil {
		return fmt.Errorf("Unable to find user with configured username: %v", p.Username)
	} else {
		p.BotUserID = user.Id
	}

	return nil
}

func (p *Plugin) IsValid() error {
	if p.GitLabURL == "" {
		return fmt.Errorf("Must have a GitLab URL specified")
	}

	if p.GitLabOAuthClientID == "" {
		return fmt.Errorf("Must have a GitLab oauth client id")
	}

	if p.GitLabOAuthClientSecret == "" {
		return fmt.Errorf("Must have a GitLab oauth client secret")
	}

	if p.EncryptionKey == "" {
		return fmt.Errorf("Must have an encryption key")
	}

	if p.Username == "" {
		return fmt.Errorf("Need a user to make posts as")
	}

	config := p.API.GetConfig()

	if config.ServiceSettings.SiteURL == nil || *config.ServiceSettings.SiteURL == "" {
		return fmt.Errorf("ServiceSettings.SiteURL must be configured")
	}

	return nil
}

func (p *Plugin) getOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     p.GitLabOAuthClientID,
		ClientSecret: p.GitLabOAuthClientSecret,
		Scopes:       []string{"api"}, // TODO restrict the scope
		RedirectURL:  fmt.Sprintf("%s/plugins/gitlab/oauth/complete", *p.API.GetConfig().ServiceSettings.SiteURL),
		Endpoint: oauth2.Endpoint{
			AuthURL:  fmt.Sprintf("%s/oauth/authorize", p.GitLabURL),
			TokenURL: fmt.Sprintf("%s/oauth/token", p.GitLabURL),
		},
	}
}

type GitLabUserInfo struct {
	UserID         string
	Token          *oauth2.Token
	GitLabUserId   int
	GitLabUsername string
	// LastToDoPostAt int64
	Settings       *UserSettings
}

type UserSettings struct {
	// SidebarButtons string `json:"sidebar_buttons"`
	// DailyReminder  bool   `json:"daily_reminder"`
	// Notifications  bool   `json:"notifications"`
}

func (p *Plugin) storeGitLabUserInfo(info *GitLabUserInfo) error {
	encryptedToken, err := encrypt([]byte(p.EncryptionKey), info.Token.AccessToken)
	if err != nil {
		return err
	}

	info.Token.AccessToken = encryptedToken

	jsonInfo, err := json.Marshal(info)
	if err != nil {
		return err
	}

	if err := p.API.KVSet(info.UserID+GITLAB_TOKEN_KEY, jsonInfo); err != nil {
		return err
	}

	return nil
}

func (p *Plugin) getGitLabUserInfo(userID string) (*GitLabUserInfo, *APIErrorResponse) {
	var userInfo GitLabUserInfo

	if infoBytes, err := p.API.KVGet(userID + GITLAB_TOKEN_KEY); err != nil || infoBytes == nil {
		return nil, &APIErrorResponse{ID: API_ERROR_ID_NOT_CONNECTED, Message: "Must connect user account to GitLab first.", StatusCode: http.StatusBadRequest}
	} else if err := json.Unmarshal(infoBytes, &userInfo); err != nil {
		return nil, &APIErrorResponse{ID: "", Message: "Unable to parse token.", StatusCode: http.StatusInternalServerError}
	}

	unencryptedToken, err := decrypt([]byte(p.EncryptionKey), userInfo.Token.AccessToken)
	if err != nil {
		mlog.Error(err.Error())
		return nil, &APIErrorResponse{ID: "", Message: "Unable to decrypt access token.", StatusCode: http.StatusInternalServerError}
	}

	userInfo.Token.AccessToken = unencryptedToken

	return &userInfo, nil
}

func (p *Plugin) storeGitLabToUserIDMapping(gitlabUsername, userID string) error {
	if err := p.API.KVSet(gitlabUsername+GITLAB_USERNAME_KEY, []byte(userID)); err != nil {
		return fmt.Errorf("Encountered error saving gitlab username mapping")
	}
	return nil
}

func (p *Plugin) getGitLabToUserIDMapping(gitlabUsername string) string {
	userID, _ := p.API.KVGet(gitlabUsername + GITLAB_USERNAME_KEY)
	return string(userID)
}

// func (p *Plugin) disconnectGitHubAccount(userID string) {
// 	userInfo, _ := p.getGitLabUserInfo(userID)
// 	if userInfo == nil {
// 		return
// 	}

// 	p.API.KVDelete(userID + GITLAB_TOKEN_KEY)
// 	p.API.KVDelete(userInfo.GitHubUsername + GITLAB_USERNAME_KEY)

// 	if user, err := p.API.GetUser(userID); err == nil && user.Props != nil && len(user.Props["gitlab_user"]) > 0 {
// 		delete(user.Props, "gitlab_user")
// 		p.API.UpdateUser(user)
// 	}

// 	p.API.PublishWebSocketEvent(
// 		WS_EVENT_DISCONNECT,
// 		nil,
// 		&model.WebsocketBroadcast{UserId: userID},
// 	)
// }

func (p *Plugin) CreateBotDMPost(userID, message, postType string) *model.AppError {
	channel, err := p.API.GetDirectChannel(userID, p.BotUserID)
	if err != nil {
		mlog.Error("Couldn't get bot's DM channel", mlog.String("user_id", userID))
		return err
	}

	post := &model.Post{
		UserId:    p.BotUserID,
		ChannelId: channel.Id,
		Message:   message,
		Type:      postType,
		Props: map[string]interface{}{
			"from_webhook":      "true",
			"override_username": GITLAB_USERNAME,
			"override_icon_url": GITLAB_ICON_URL,
		},
	}

	if _, err := p.API.CreatePost(post); err != nil {
		mlog.Error(err.Error())
		return err
	}

	return nil
}

// func (p *Plugin) PostToDo(info *GitLabUserInfo) {
// 	text, err := p.GetToDo(context.Background(), info.GitHubUsername, p.githubConnect(*info.Token))
// 	if err != nil {
// 		mlog.Error(err.Error())
// 		return
// 	}

// 	p.CreateBotDMPost(info.UserID, text, "custom_git_todo")
// }

// func (p *Plugin) GetToDo(ctx context.Context, username string, githubClient *github.Client) (string, error) {
// 	issueResults, _, err := githubClient.Search.Issues(ctx, getReviewSearchQuery(username, p.GitLabOrg), &github.SearchOptions{})
// 	if err != nil {
// 		return "", err
// 	}

// 	notifications, _, err := githubClient.Activity.ListNotifications(ctx, &github.NotificationListOptions{})
// 	if err != nil {
// 		return "", err
// 	}

// 	yourPrs, _, err := githubClient.Search.Issues(ctx, getYourPrsSearchQuery(username, p.GitLabOrg), &github.SearchOptions{})
// 	if err != nil {
// 		return "", err
// 	}

// 	yourAssignments, _, err := githubClient.Search.Issues(ctx, getYourAssigneeSearchQuery(username, p.GitLabOrg), &github.SearchOptions{})
// 	if err != nil {
// 		return "", err
// 	}

// 	text := "##### Unread Messages\n"

// 	notificationCount := 0
// 	notificationContent := ""
// 	for _, n := range notifications {
// 		if n.GetReason() == "subscribed" {
// 			continue
// 		}

// 		if n.GetRepository() == nil {
// 			p.API.LogError("Unable to get repository for notification in todo list. Skipping.")
// 			continue
// 		}

// 		if p.checkOrg(n.GetRepository().GetOwner().GetLogin()) != nil {
// 			continue
// 		}

// 		switch n.GetSubject().GetType() {
// 		case "RepositoryVulnerabilityAlert":
// 			message := fmt.Sprintf("[Vulnerability Alert for %v](%v)", n.GetRepository().GetFullName(), fixGithubNotificationSubjectURL(n.GetSubject().GetURL()))
// 			notificationContent += fmt.Sprintf("* %v\n", message)
// 		default:
// 			url := fixGithubNotificationSubjectURL(n.GetSubject().GetURL())
// 			notificationContent += fmt.Sprintf("* %v\n", url)
// 		}

// 		notificationCount++
// 	}

// 	if notificationCount == 0 {
// 		text += "You don't have any unread messages.\n"
// 	} else {
// 		text += fmt.Sprintf("You have %v unread messages:\n", notificationCount)
// 		text += notificationContent
// 	}

// 	text += "##### Review Requests\n"

// 	if issueResults.GetTotal() == 0 {
// 		text += "You have don't have any pull requests awaiting your review.\n"
// 	} else {
// 		text += fmt.Sprintf("You have %v pull requests awaiting your review:\n", issueResults.GetTotal())

// 		for _, pr := range issueResults.Issues {
// 			text += fmt.Sprintf("* %v\n", pr.GetHTMLURL())
// 		}
// 	}

// 	text += "##### Your Open Pull Requests\n"

// 	if yourPrs.GetTotal() == 0 {
// 		text += "You have don't have any open pull requests.\n"
// 	} else {
// 		text += fmt.Sprintf("You have %v open pull requests:\n", yourPrs.GetTotal())

// 		for _, pr := range yourPrs.Issues {
// 			text += fmt.Sprintf("* %v\n", pr.GetHTMLURL())
// 		}
// 	}

// 	text += "##### Your Assigments\n"

// 	if yourAssignments.GetTotal() == 0 {
// 		text += "You have don't have any assignments.\n"
// 	} else {
// 		text += fmt.Sprintf("You have %v assignments:\n", yourAssignments.GetTotal())

// 		for _, assign := range yourAssignments.Issues {
// 			text += fmt.Sprintf("* %v\n", assign.GetHTMLURL())
// 		}
// 	}

// 	return text, nil
// }

// func (p *Plugin) checkOrg(org string) error {
// 	configOrg := strings.TrimSpace(p.GitLabOrg)
// 	if configOrg != "" && configOrg != org {
// 		return fmt.Errorf("Only repositories in the %v organization are supported", configOrg)
// 	}

// 	return nil
// }

// func (p *Plugin) sendRefreshEvent(userID string) {
// 	p.API.PublishWebSocketEvent(
// 		WS_EVENT_REFRESH,
// 		nil,
// 		&model.WebsocketBroadcast{UserId: userID},
// 	)
// }
