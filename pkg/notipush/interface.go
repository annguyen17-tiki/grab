package notipush

type IService interface {
	Send(input *FcmMessageInput) error
}

type FcmMessageInput struct {
	Title string
	Body  string
	Token string
	Link  string
}
