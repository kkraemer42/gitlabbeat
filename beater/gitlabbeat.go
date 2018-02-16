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

	git := gitlab.NewClient(nil, bt.config.AccessToken)

	if err != nil {
		return err
	}

	bt.git = git
	git.SetBaseURL(bt.config.gitlab_address)

	ticker := time.NewTicker(bt.config.Period)
	counter := 1

	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
			logp.Info("Collecting events.")

			bt.getDeployKeys()
			//	bt.collectOrgsEvents(jobCtx, bt.config.Orgs)

		}

		logp.Info("Event sent")
		counter++
	}
}

func (bt *Gitlabbeat) getDeployKeys() {

	res, _, err := bt.git.DeployKeys.ListAllDeployKeys()
	if err != nil {
		logp.Err("Failed to collect event, got :", err)
		return
	}

	for _, label := range res {
		bt.client.Publish(bt.newLabelsEvent(label))

	}

}

func (Gitlabbeat) newLabelsEvent(label *gitlab.DeployKey) beat.Event {

	event := beat.Event{
		Timestamp: time.Now(),
		Fields: common.MapStr{
			"type":   "gitlabbeat",
			"key_id": label.ID,
			"push":   label.CanPush,
			"title":  label.Title,
		},
	}

	return event
}

func (bt *Gitlabbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
