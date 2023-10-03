package entity

type User struct {
	ID       string
	Name     string
	Email    string
	Password string
}

func NewUser(id, name, email, password string) *User {
	return &User{
		ID:       id,
		Name:     name,
		Email:    email,
		Password: password,
	}
}

func (u *User) UpdatePassword(password string) {
	u.Password = password
}
