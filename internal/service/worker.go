package service

import (
	"github.com/annguyen17-tiki/grab/internal/model"
	"github.com/annguyen17-tiki/grab/pkg/notipush"
	"github.com/gocraft/work"
)

func (svc *service) startWorkers() {
	pool := work.NewWorkerPool(struct{}{}, 3, model.RedisNamespace, svc.redisPool)
	pool.JobWithOptions(model.FCMWorkerTopic, work.JobOptions{MaxFails: 3, MaxConcurrency: 3}, svc.pushWebNoti)
	pool.Start()
}

func (svc *service) pushWebNoti(job *work.Job) error {
	payload, err := model.CreateWorkerPayload(job.Args)
	if err != nil {
		return err
	}

	notiPusher, err := notipush.New()
	if err != nil {
		return err
	}

	token, err := svc.store.Token().Get(payload.AccountID)
	if err != nil {
		return err
	}

	if err = notiPusher.Send(&notipush.FcmMessageInput{
		Token: token,
		Title: payload.Title,
		Body:  payload.Body,
		Link:  payload.Link,
	}); err != nil {
		return err
	}

	return nil
}
