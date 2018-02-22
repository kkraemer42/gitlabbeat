package beater

import (
	"fmt"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"

	"github.com/xanzy/go-gitlab"

	"github.com/kkraemer42/gitlabbeat/config"
)

type Gitlabbeat struct {
	done   chan struct{}
	config config.Config
	git    *gitlab.Client
	client beat.Client
}

// Creates beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	c := config.DefaultConfig
	if err := cfg.Unpack(&c); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &Gitlabbeat{
		done:   make(chan struct{}),
		config: c,
	}
	return bt, nil
}

// The function Run runs the beat
func (bt *Gitlabbeat) Run(b *beat.Beat) error {
	logp.Info("gitlabbeat is running! Hit CTRL-C to stop it.")

	var err error
	bt.client, err = b.Publisher.Connect()

	//git := gitlab.NewClient(nil, bt.config.AccessToken)
	git := gitlab.NewClient(nil, "k8JGGxHJRjAxANoGqimh")

	if err != nil {
		return err
	}

	bt.git = git
	//git.SetBaseURL(bt.config.GitlabAdress)
	git.SetBaseURL("https://gitlab.ballpark.altemista.cloud/api/v4")

	ticker := time.NewTicker(bt.config.Period)

	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
			logp.Info("Collecting events.")

			//bt.getMergeRequests()
			bt.getIssues()
			//bt.getUsers()

		}

		logp.Info("Event sent")
	}
}

func (bt *Gitlabbeat) getMergeRequests() {

	res, _, err := bt.git.MergeRequests.ListMergeRequests(&gitlab.ListMergeRequestsOptions{})
	if err != nil {
		logp.Err("Failed to collect Merge Requests, got :", err)
		return
	}

	for _, mergeRequest := range res {
		bt.client.Publish(bt.newMergeRequestEvent(mergeRequest))

	}

}

/*
func (bt *Gitlabbeat) getPipelines() {

	res, _, err := bt.git.Projects.ListProjects(nil, nil)
	if err != nil {
		logp.Err("Failed to collect Merge Requests, got :", err)
		return
	}

	for _, projects := range res {
		for _, pipeline := range projects {


			bt.client.Publish(bt.newPipelineEvent(pipeline))

		}

		bt.client.Publish(bt.newPipelineEvent(pipeline))

	}

}

*/

func (bt *Gitlabbeat) getIssues() {

	/*
		opt := &gitlab.ListIssuesOptions{
			Scope: gitlab.String("all"),
		}
	*/

	var page = 0
	var bool = true
	for bool != true {

		res, _, err := bt.git.Issues.ListIssues(&gitlab.ListIssuesOptions{
			Scope: gitlab.String("all"),
			PerPage: gitlab.String("100"),
			Page: gitlab.String(page),
			
	
		})
		if err != nil {
			logp.Err("Failed to collect event, got :", err)
			return
		}

		if(bt.git.BaseURL != 0){
			page++
		} else {
			bool = false
		}

	}


	

	for _, issue := range res {
		bt.client.Publish(bt.newIssueEvent(issue))

	}

}

func (bt *Gitlabbeat) getUsers() {

	res, _, err := bt.git.Users.ListUsers(&gitlab.ListUsersOptions{})
	if err != nil {
		logp.Err("Failed to collect Users, got :", err)
		return
	}

	for _, user := range res {
		bt.client.Publish(bt.newUserEvent(user))

	}

}

func (Gitlabbeat) newUserEvent(user *gitlab.User) beat.Event {

	event := beat.Event{
		Timestamp: time.Now(),
		Fields: common.MapStr{
			"type":             "gitlabbeat",
			"user_id":          user.ID,
			"user_mail":        user.Email,
			"user_adminstatus": user.IsAdmin,
			"user_state":       user.State,
		},
	}

	return event
}

func (Gitlabbeat) newMergeRequestEvent(mergeRequest *gitlab.MergeRequest) beat.Event {

	event := beat.Event{
		Timestamp: time.Now(),
		Fields: common.MapStr{
			"type":                     "gitlabbeat",
			"mergeRequestId":           mergeRequest.ID,
			"mergeRequestCreationDate": mergeRequest.CreatedAt,
			"mergeRequestTitle":        mergeRequest.Title,
			"mergeRequestAssignee":     mergeRequest.Assignee,
		},
	}

	return event
}

func (Gitlabbeat) newPipelineEvent(pipeline *gitlab.Pipeline) beat.Event {

	event := beat.Event{
		Timestamp: time.Now(),
		Fields: common.MapStr{
			"type":              "gitlabbeat",
			"pipelineId":        pipeline.ID,
			"pipelineStartTime": pipeline.StartedAt,
			"pipelineStatus":    pipeline.Status,
			"pipelineDuration":  pipeline.Duration,
		},
	}

	return event
}

func (Gitlabbeat) newIssueEvent(issue *gitlab.Issue) beat.Event {

	event := beat.Event{
		Timestamp: time.Now(),
		Fields: common.MapStr{
			"type":            "gitlabbeat",
			"issueId":         issue.ID,
			"issueTitle":      issue.Title,
			"issueMilestone":  issue.Milestone,
			"issueState":      issue.State,
			"issueLastUpdate": issue.UpdatedAt,
			"issueAssignees":  issue.Assignees,
			"issueProject":    issue.ProjectID,
			"issueDueDate":    issue.DueDate,
		},
	}

	return event
}

func (bt *Gitlabbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
