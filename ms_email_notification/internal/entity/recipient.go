package entity

type Recipient struct {
	Name  string
	Email string
}

func NewRecipient(name, email string) *Recipient {
	return &Recipient{
		Name:  name,
		Email: email,
	}
}
