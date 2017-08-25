package user

import (
	"strings"

	"github.com/xinzf/go-utils"

	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
)

type UserModel struct {
	Iid      string
	Name     string
	Icon     string
	Props    string
	CreateAt string
}

func (this UserModel) TableName() string {
	return "user"
}

func (user *UserModel) BeforeCreate(scope *gorm.Scope) error {
	if user.Iid == "" {
		scope.SetColumn("Iid", strings.Replace(uuid.NewV4().String(), "-", "", -1))
	}
	scope.SetColumn("CreateAt", utils.NowTime())
	return nil
}
