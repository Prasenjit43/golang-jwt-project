package helpers

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
)

func CheckUserType(context *gin.Context, role string) (err error) {
	userType := context.GetString("userType")
	err = nil

	if userType != role {
		err = errors.New("unauthorized to access this resource")
		return err
	}

	return err

}

func MatchUserTypeWithUserId(context *gin.Context, userId string) (err error) {
	userType := context.GetString("userType")
	uid := context.GetString("uid")
	err = nil
	fmt.Println("userType :", userType)
	fmt.Println("uid :", uid)

	if userType == "USER" && uid != userId {
		err = errors.New("unauthorized to access this resource")
		return err
	}

	err = CheckUserType(context, userType)
	return err

}
