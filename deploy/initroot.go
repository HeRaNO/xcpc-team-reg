package deploy

import (
	"context"
	"fmt"

	"github.com/HeRaNO/xcpc-team-reg/model"
)

func InitRootUser() {
	var pwdtoken, email string
	n, err := fmt.Scanf("%s %s", &pwdtoken, &email)
	if err != nil {
		panic(err)
	}
	if n != 2 {
		panic("[ERROR] read pwdtoken and email error")
	}
	rootReg := model.UserRegister{
		Name:     "root",
		Email:    email,
		PwdToken: pwdtoken,
	}
	err = model.CreateNewUser(context.Background(), rootReg, 1)
	if err != nil {
		panic(err)
	}
}
