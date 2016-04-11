package beater

import (
	"fmt"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/cfgfile"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"

	"github.com/s-fujimoto/deviobeat/config"
    
	"net/http"
	"io/ioutil"
	"github.com/bitly/go-simplejson"
)

type Deviobeat struct {
	beatConfig *config.Config
	done       chan struct{}
	period     time.Duration
}

// Creates beater
func New() *Deviobeat {
	return &Deviobeat{
		done: make(chan struct{}),
	}
}

/// *** Beater interface methods ***///

func (bt *Deviobeat) Config(b *beat.Beat) error {

	// Load beater beatConfig
	err := cfgfile.Read(&bt.beatConfig, "")
	if err != nil {
		return fmt.Errorf("Error reading config file: %v", err)
	}

	return nil
}

func (bt *Deviobeat) Setup(b *beat.Beat) error {

	// Setting default period if not set
	if bt.beatConfig.Deviobeat.Period == "" {
		bt.beatConfig.Deviobeat.Period = "1s"
	}

	var err error
	bt.period, err = time.ParseDuration(bt.beatConfig.Deviobeat.Period)
	if err != nil {
		return err
	}

	return nil
}

func (bt *Deviobeat) Run(b *beat.Beat) error {
	logp.Info("deviobeat is running! Hit CTRL-C to stop it.")

	ticker := time.NewTicker(bt.period)
	counter := 1
	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}
        
        url := bt.beatConfig.Deviobeat.URL
        resp, _ := http.Get(url)
        body, _ := ioutil.ReadAll(resp.Body)
        json, _ := simplejson.NewJson(body)
        doccount, _ := json.Get("count").Int()

		event := common.MapStr{
			"@timestamp": common.Time(time.Now()),
			"type":       b.Name,
			"counter":    counter,
            "doc_count":  doccount,
		}
		b.Events.PublishEvent(event)
		logp.Info("Event sent")
		counter++
	}
}

func (bt *Deviobeat) Cleanup(b *beat.Beat) error {
	return nil
}

func (bt *Deviobeat) Stop() {
	close(bt.done)
}
