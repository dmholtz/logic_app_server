package server

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserStore interface {
	Signup(c Credentials) error
	Login(c Credentials) (string, error)
	Logout(token string) error
}

type MyUserStore struct {
}

func (us *MyUserStore) Signup(credential Credentials) error {
	return nil
}

func (us *MyUserStore) Login(credentials Credentials) (string, error) {
	return "", nil
}

func (us *MyUserStore) Logout(token string) error {
	return nil
}
