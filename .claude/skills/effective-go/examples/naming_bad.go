package user_utils

type UserStruct struct {
	UserId   int
	UserName string
}

func FindUserById(userId int) (*UserStruct, error) {
	return &UserStruct{UserId: userId}, nil
}
