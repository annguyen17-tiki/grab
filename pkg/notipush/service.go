package notipush

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"

	"google.golang.org/api/option"
)

type service struct {
	conf   *config
	client *messaging.Client
}

func New() (IService, error) {
	conf, err := loadConfig()
	if err != nil {
		return nil, err
	}

	svc := &service{}
	svc.conf = conf

	ctx := context.Background()
	app, err := firebase.NewApp(ctx, nil, option.WithCredentialsFile(svc.conf.CredentialsFile))
	if err != nil {
		return nil, fmt.Errorf("error initializing app: %v", err)
	}

	svc.client, err = app.Messaging(ctx)
	if err != nil {
		return nil, fmt.Errorf("error initializing message: %v", err)
	}

	return svc, nil
}

func (svc *service) Send(input *FcmMessageInput) error {
	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: input.Title,
			Body:  input.Body,
		},
		Webpush: &messaging.WebpushConfig{
			FCMOptions: &messaging.WebpushFCMOptions{
				Link: input.Link,
			},
		},
		Token: input.Token,
	}

	response, err := svc.client.Send(context.Background(), message)
	if err != nil {
		return err
	}

	fmt.Printf("successfully sent message: %s", response)
	return nil
}
