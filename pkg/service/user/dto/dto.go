package dto

import "github.com/ardihikaru/graphql-example-part-1/pkg/graph/user/model"

type User model.User

// ToModel transforms struct
func (u *User) ToModel() *model.User {
	return &model.User{
		UserID:   u.UserID,
		UserNm:   u.UserNm,
		IsAdmin:  u.IsAdmin,
		StatusCd: u.StatusCd,
	}
}
