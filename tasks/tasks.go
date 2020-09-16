package tasks

import (
	"encoding/json"
	"github.com/heteddy/delaytask-go/delaytask"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

type OncePingTask struct {
	delaytask.Task
	Url string `json:"Url"`
}

func (t *OncePingTask) Run() (bool, error) {
	resp, err := http.Get(t.Url)
	if err != nil {
		return false, err
	}
	t.RunAt = delaytask.TaskTime(time.Now())
	delaytask.Logger.WithFields(logrus.Fields{
		"id":      t.GetID(),
		"RunAt":   t.GetRunAt(),
		"ToRunAt": t.GetToRunAt(),
	}).Infoln("OncePingTask ToRunAt RunAt")

	defer resp.Body.Close()
	return true, nil
}

func (t *OncePingTask) ToJson() string {
	b, err := json.Marshal(t)
	if err != nil {
		return ""
	}
	return string(b)
}

type PeriodPingTask struct {
	delaytask.PeriodicTask
	Url string `json:"Url"`
}

func (t *PeriodPingTask) Run() (bool, error) {
	resp, err := http.Get(t.Url)
	defer resp.Body.Close()
	if err != nil {
		return false, err
	}
	ioutil.ReadAll(resp.Body)
	delaytask.Logger.WithFields(logrus.Fields{
		"id":  t.GetID(),
		"err": err,
	}).Infoln("PeriodPingTask Run")
	return true, nil
}
func (t *PeriodPingTask) ToJson() string {
	b, err := json.Marshal(t)
	if err != nil {
		return ""
	}
	return string(b)
}
