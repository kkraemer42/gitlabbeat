package beater

import (
	"fmt"
	"os"
	"strconv"
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

// New Creates beater
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

//Run runs the beat
func (bt *Gitlabbeat) Run(b *beat.Beat) error {
	logp.Info("gitlabbeat is running! Hit CTRL-C to stop it.")

	var err error
	bt.client, err = b.Publisher.Connect()

	git := gitlab.NewClient(nil, os.Getenv("ACCESSTOKEN"))

	if err != nil {
		return err
	}

	bt.git = git
	git.SetBaseURL(os.Getenv("GITLABADDRESS"))

	projectID, idErr := strconv.Atoi(os.Getenv("PROJECTID"))
	if idErr != nil {
		return idErr
	}

	var period, error = time.ParseDuration(os.Getenv("COLLECTIONPERIOD"))
	if error != nil {
		return err
	}
	ticker := time.NewTicker(period)
	//	ticker := time.NewTicker(bt.config.Period)

	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
			logp.Info("Collecting events.")

			bt.getMergeRequests(projectID)
			bt.getIssues(projectID)
			bt.getCommits(projectID)
			bt.getProjects()
			bt.getUsers()
		}

		logp.Info("Event sent")
	}
}

func (bt *Gitlabbeat) getMergeRequests(projectID int) {

	var page = 0
	var bool = true
	for bool != false {

		result, response, err := bt.git.MergeRequests.ListProjectMergeRequests(projectID, &gitlab.ListProjectMergeRequestsOptions{
			ListOptions: gitlab.ListOptions{Page: page, PerPage: 100},
			Scope:       gitlab.String("all"),
		})
		if err != nil {
			logp.Err("Failed to collect mergeRequests, got :", err)
			return
		}

		if response.NextPage != 0 {
			page++
		} else {
			bool = false
		}

		for _, mergeRequest := range result {
			bt.client.Publish(bt.newMergeRequestEvent(mergeRequest))

		}

	}

}

func (bt *Gitlabbeat) getProjects() {

	var page = 0
	var bool = true
	for bool != false {

		result, response, err := bt.git.Projects.ListProjects(&gitlab.ListProjectsOptions{
			ListOptions: gitlab.ListOptions{Page: page, PerPage: 100},
		})
		if err != nil {
			logp.Err("Failed to collect projects, got :", err)
			return
		}

		if response.NextPage != 0 {
			page++
		} else {
			bool = false
		}

		for _, project := range result {
			bt.client.Publish(bt.newProjectEvent(project))

		}

	}

}

func (bt *Gitlabbeat) getCommits(projectID int) {

	var page = 0
	var bool = true
	for bool != false {

		result, response, err := bt.git.Commits.ListCommits(projectID, &gitlab.ListCommitsOptions{
			ListOptions: gitlab.ListOptions{Page: page, PerPage: 100},
		})
		if err != nil {
			logp.Err("Failed to collect commits, got :", err)
			return
		}

		if response.NextPage != 0 {
			page++
		} else {
			bool = false
		}

		for _, commit := range result {
			bt.client.Publish(bt.newCommitEvent(commit))

		}

	}

}

func (bt *Gitlabbeat) getIssues(projectID int) {

	var page = 0
	var bool = true
	for bool != false {

		result, response, err := bt.git.Issues.ListProjectIssues(projectID, &gitlab.ListProjectIssuesOptions{
			ListOptions: gitlab.ListOptions{Page: page, PerPage: 100},
			Scope:       gitlab.String("all"),
		})
		if err != nil {
			logp.Err("Failed to collect issues, got :", err)
			return
		}

		if response.NextPage != 0 {
			page++
		} else {
			bool = false
		}

		for _, issue := range result {
			bt.client.Publish(bt.newIssueEvent(issue))

		}

	}

}

func (bt *Gitlabbeat) getUsers() {

	var page = 0
	var bool = true
	for bool != false {

		result, response, err := bt.git.Users.ListUsers(&gitlab.ListUsersOptions{
			ListOptions: gitlab.ListOptions{Page: page, PerPage: 100},
		})
		if err != nil {
			logp.Err("Failed to collect users, got :", err)
			return
		}

		if response.NextPage != 0 {
			page++
		} else {
			bool = false
		}

		for _, user := range result {
			bt.client.Publish(bt.newUserEvent(user))

		}

	}

}

func (bt *Gitlabbeat) getPipelines() {

	var page = 0
	var bool = true
	for bool != false {

		result, response, err := bt.git.Pipelines.ListProjectPipelines(1287, &gitlab.ListProjectPipelinesOptions{
			ListOptions: gitlab.ListOptions{Page: page, PerPage: 20},
		})
		if err != nil {
			logp.Err("Failed to collect pipelines, got :", err)
			return
		}

		if response.NextPage != 0 {
			page++
		} else {
			bool = false
		}

		for _, pipelineID := range result {

			pipeline, _, err := bt.git.Pipelines.GetPipeline(1287, pipelineID.ID)

			if err != nil {
				logp.Err("Couldn't load pipeline details, got :", err)
				return
			}
			bt.client.Publish(bt.newPipelineEvent(pipeline))
		}

	}

}

func (Gitlabbeat) newPipelineEvent(pipeline *gitlab.Pipeline) beat.Event {

	event := beat.Event{
		Timestamp: time.Now(),
		Fields: common.MapStr{
			"type":                "gitlabbeat",
			"pipeline_status":     pipeline.Status,
			"pipeline_startDate":  pipeline.StartedAt,
			"pipeline_yamlErrors": pipeline.YamlErrors,
			"pipeline_user":       pipeline.User,
			"pipeline_duration":   pipeline.Duration,
			"pipeline_finishDate": pipeline.FinishedAt,
			"pipeline_ID":         pipeline.ID,
		},
	}

	return event
}

func (Gitlabbeat) newMergeRequestEvent(mergeRequest *gitlab.MergeRequest) beat.Event {

	event := beat.Event{
		Timestamp: time.Now(),
		Fields: common.MapStr{
			"type":                      "gitlabbeat",
			"mergeRequest_Id":           mergeRequest.ID,
			"mergeRequest_CreationDate": mergeRequest.CreatedAt,
			"mergeRequest_Title":        mergeRequest.Title,
			"mergeRequest_Assignee":     mergeRequest.Assignee,
		},
	}

	return event
}

func (Gitlabbeat) newCommitEvent(commit *gitlab.Commit) beat.Event {

	event := beat.Event{
		Timestamp: time.Now(),
		Fields: common.MapStr{
			"type":        "gitlabbeat",
			"commit_date": commit.CommittedDate,
			"commit_id":   commit.ID,
		},
	}

	return event

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

func (Gitlabbeat) newProjectEvent(project *gitlab.Project) beat.Event {

	event := beat.Event{
		Timestamp: time.Now(),
		Fields: common.MapStr{
			"type":          "gitlabbeat",
			"project_id":    project.ID,
			"project_name":  project.Name,
			"project_owner": project.Owner,
		},
	}

	return event
}

func (Gitlabbeat) newIssueEvent(issue *gitlab.Issue) beat.Event {

	event := beat.Event{
		Timestamp: time.Now(),
		Fields: common.MapStr{
			"type":             "gitlabbeat",
			"issue_Id":         issue.ID,
			"issue_Title":      issue.Title,
			"issue_Milestone":  issue.Milestone,
			"issue_State":      issue.State,
			"issue_LastUpdate": issue.UpdatedAt,
			"issue_Assignees":  issue.Assignees,
			"issue_Project":    issue.ProjectID,
			"issue_DueDate":    issue.DueDate,
		},
	}

	return event
}

func (bt *Gitlabbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
