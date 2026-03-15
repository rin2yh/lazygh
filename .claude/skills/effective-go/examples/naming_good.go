package users

type User struct {
	ID   int
	Name string
}

func Find(id int) (*User, error) {
	return &User{ID: id, Name: "alice"}, nil
}
