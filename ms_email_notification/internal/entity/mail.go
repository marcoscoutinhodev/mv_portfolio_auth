package entity

type Mail struct {
	Recipient Recipient
	Title     string
	Body      []byte
}

func NewMail(recipient Recipient, title string, body []byte) *Mail {
	return &Mail{
		Recipient: recipient,
		Title:     title,
		Body:      body,
	}
}
