package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pvg/entity"
	"pvg/model"
	"pvg/repository"
	"strconv"
	"time"

	"pvg/util"

	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input entity.UserRegisterInput

		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, entity.ErrResponse{
				SourceFunction: "CreateUser",
				Detail:         fmt.Sprintf("error on parsing input, detail = %v", err),
			})
			return
		}

		resp, err := util.CallCheckUser(entity.CheckUserExist{
			Username:    input.Username,
			Email:       input.Email,
			PhoneNumber: input.PhoneNumber,
		})

		if err != nil {
			ctx.JSON(http.StatusBadRequest, entity.ErrResponse{
				SourceFunction: "CreateUser",
				Detail:         fmt.Sprintf("error on calling check user exist, detail = %v", err),
			})
			return
		}

		var checkUserResp entity.BasicResponse
		json.Unmarshal(resp, &checkUserResp)

		if checkUserResp.Message != "user does not exist" {
			ctx.JSON(http.StatusBadRequest, entity.ErrResponse{
				SourceFunction: "CreateUser",
				Detail:         "user already exist",
			})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, entity.ErrResponse{
				SourceFunction: "CreateUser",
				Detail:         fmt.Sprintf("error on hashing password, detail = %v", err),
			})
			return
		}

		birthDateParsed, err := time.Parse("2006-01-02", input.Birthday)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, entity.ErrResponse{
				SourceFunction: "CreateUser",
				Detail:         fmt.Sprintf("error on parsing birthdate, detail = %v", err),
			})
			return
		}

		_, err = repository.CreateUser(&model.User{
			Username:    input.Username,
			Password:    string(hashedPassword),
			Firstname:   input.FirstName,
			Lastname:    input.LastName,
			PhoneNumber: input.PhoneNumber,
			Email:       input.Email,
			Birthday:    birthDateParsed,
		})

		if err != nil {
			ctx.JSON(http.StatusBadRequest, entity.ErrResponse{
				SourceFunction: "CreateUser",
				Detail:         err.Error(),
			})
			return
		}

		otpSecret := util.RandomSecret(20)
		passcode, err := totp.GenerateCodeCustom(otpSecret, time.Now(), totp.ValidateOpts{
			Period:    180,
			Skew:      1,
			Digits:    4,
			Algorithm: otp.AlgorithmSHA512,
		})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, entity.ErrResponse{
				SourceFunction: "CreateUser",
				Detail:         fmt.Sprintf("error on generating otp, details = %v", err),
			})
			return
		}

		otp, err := repository.GetOtpByEmail(input.Email)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, entity.ErrResponse{
				SourceFunction: "CreateUser",
				Detail:         fmt.Sprintf("error on retrieving otp = %v", err),
			})
			return
		}

		if otp == nil {
			_, err := repository.AddOtp(&model.Otp{
				OtpExpire:      time.Now().Add(time.Minute * time.Duration(2)),
				ActivationCode: passcode,
				Email:          input.Email,
				OtpSecret:      otpSecret,
			})
			if err != nil {
				ctx.JSON(http.StatusBadRequest, entity.ErrResponse{
					SourceFunction: "CreateUser",
					Detail:         fmt.Sprintf("error on retrieving otp, details = %v", err),
				})
				return
			}
		} else {
			_, err := repository.UpdateOtp(otp.IDOtp, &model.Otp{
				OtpExpire:      time.Now().Add(time.Minute * time.Duration(2)),
				ActivationCode: passcode,
				Email:          input.Email,
				OtpSecret:      otpSecret,
			})
			if err != nil {
				ctx.JSON(http.StatusBadRequest, entity.ErrResponse{
					SourceFunction: "CreateUser",
					Detail:         fmt.Sprintf("error on updating otp, details = %v", err),
				})
				return
			}
		}

		bodyText := fmt.Sprintf(`Hello, here is your otp. %s`, passcode)
		header := "User Registration"
		if err := SendMail(bodyText, header, input.Email); err != nil {
			ctx.JSON(http.StatusInternalServerError, entity.ErrResponse{
				SourceFunction: "CreateUser",
				Detail:         fmt.Sprintf("error on sending mail, detail = %v", err),
			})
			return
		}

		ctx.JSON(http.StatusOK, entity.BasicResponse{
			Message: "mail sent",
		})

	}
}

func GetSpecificUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.Param("userId")

		userIdInt, err := strconv.Atoi(userId)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, entity.ErrResponse{
				SourceFunction: "GetSpecificUser",
				Detail:         fmt.Sprintf("cannot parse user id, detail = %v", err),
			})
			return
		}

		result, err := repository.GetSpecificUser(int32(userIdInt))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, entity.ErrResponse{
				SourceFunction: "GetSpecificUser",
				Detail:         fmt.Sprintf("cannot get user, detail = %v", err),
			})
			return
		}

		ctx.JSON(http.StatusOK, result)

	}
}

func GetAllUsers() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		users, err := repository.GetAllUsers()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, entity.ErrResponse{
				SourceFunction: "GetAllUsers",
				Detail:         fmt.Sprintf("cannot get users, detail = %v", err),
			})
			return
		}

		ctx.JSON(http.StatusOK, users)
	}
}

func UpdateUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.Param("userId")

		userIdInt, err := strconv.Atoi(userId)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, entity.ErrResponse{
				SourceFunction: "UpdateUser",
				Detail:         fmt.Sprintf("cannot parse user id, detail = %v", err),
			})
			return
		}

		var input entity.UserUpdateInput

		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, entity.ErrResponse{
				SourceFunction: "UpdateUser",
				Detail:         fmt.Sprintf("error on parsing input, detail = %v", err),
			})
			return
		}

		user, err := repository.GetSpecificUser(int32(userIdInt))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, entity.ErrResponse{
				SourceFunction: "UpdateUser",
				Detail:         err.Error(),
			})
			return
		}

		var updatedUser model.User

		if len(input.UpdatePassword) > 0 {
			if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.CurrentPassword)); err != nil {
				ctx.JSON(http.StatusBadRequest, entity.ErrResponse{
					SourceFunction: "UpdateUser",
					Detail:         fmt.Sprintf("invalid password, detail = %v", err),
				})
				return
			}

			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.UpdatePassword), bcrypt.DefaultCost)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, entity.ErrResponse{
					SourceFunction: "UpdateUser",
					Detail:         fmt.Sprintf("error on hashing password, detail = %v", err),
				})
				return
			}
			updatedUser.Password = string(hashedPassword)
		}

		if len(input.Birthday) > 0 {
			birthDateParsed, err := time.Parse("2006-01-02", input.Birthday)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, entity.ErrResponse{
					SourceFunction: "UpdateUser",
					Detail:         fmt.Sprintf("error on parsing birthdate, detail = %v", err),
				})
				return
			}
			updatedUser.Birthday = birthDateParsed
		}

		updatedUser.Username = input.Username
		updatedUser.Firstname = input.FirstName
		updatedUser.Lastname = input.LastName
		updatedUser.PhoneNumber = input.PhoneNumber
		updatedUser.Email = input.Email

		_, err = repository.UpdateUser(int32(userIdInt), &updatedUser)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, entity.ErrResponse{
				SourceFunction: "UpdateUser",
				Detail:         fmt.Sprintf("error on hashing password, detail = %v", err),
			})
			return
		}

		ctx.JSON(http.StatusOK, entity.BasicResponse{
			Message: "user updated",
		})
	}
}

func DeleteUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.Param("userId")

		userIdInt, err := strconv.Atoi(userId)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, entity.ErrResponse{
				SourceFunction: "DeleteUser",
				Detail:         fmt.Sprintf("cannot parse user id, detail = %v", err),
			})
			return
		}

		if err := repository.DeleteUser(int32(userIdInt)); err != nil {
			ctx.JSON(http.StatusBadRequest, entity.ErrResponse{
				SourceFunction: "DeleteUser",
				Detail:         fmt.Sprintf("cannot delete, detail = %v", err),
			})
			return
		}

		ctx.JSON(http.StatusOK, entity.BasicResponse{
			Message: "user deleted",
		})
	}
}

func CheckUserExist() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input entity.CheckUserExist

		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, entity.ErrResponse{
				SourceFunction: "CheckUserExist",
				Detail:         fmt.Sprintf("error on parsing input, detail = %v", err),
			})
			return
		}

		user, err := repository.GetUserByUsernameOrEmailOrPhoneNumber(input.Username, input.Email, input.PhoneNumber)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, entity.ErrResponse{
				SourceFunction: "CheckUserExist",
				Detail:         fmt.Sprintf("error on retreiving user, detail = %v", err),
			})
			return
		}

		if user == nil {
			ctx.JSON(http.StatusOK, entity.BasicResponse{
				Message: "user does not exist",
			})
			return
		}

		ctx.JSON(http.StatusOK, entity.BasicResponse{
			Message: "user exist",
		})
	}
}

func VerifyUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input entity.OtpInput

		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, entity.ErrResponse{
				SourceFunction: "VerifyUser",
				Detail:         fmt.Sprintf("error on parsing input, detail = %v", err),
			})
			return
		}

		otpRec, err := repository.GetOtpByEmail(input.Email)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, entity.ErrResponse{
				SourceFunction: "VerifyUser",
				Detail:         fmt.Sprintf("error on retrieving otp, detail = %v", err),
			})
			return
		}

		if otpRec == nil {
			ctx.JSON(http.StatusBadRequest, entity.ErrResponse{
				SourceFunction: "VerifyUser",
				Detail:         "otp record not found",
			})
			return
		}

		result, err := totp.ValidateCustom(input.Passcode, otpRec.OtpSecret, time.Now(), totp.ValidateOpts{
			Period:    180,
			Skew:      1,
			Digits:    4,
			Algorithm: otp.AlgorithmSHA512,
		})

		if result {
			if err := repository.DeleteOtp(otpRec.IDOtp); err != nil {
				ctx.JSON(http.StatusInternalServerError, entity.ErrResponse{
					SourceFunction: "VerifyUser",
					Detail:         fmt.Sprintf("error on deleting otp, detail = %v", err),
				})
				return
			}
		}
		if err != nil {
			ctx.JSON(http.StatusBadRequest, entity.ErrResponse{
				SourceFunction: "VerifyUser",
				Detail:         fmt.Sprintf("error on validating otp, detail = %v", err),
			})
			return
		}

		user, err := repository.GetUserByEmail(input.Email)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, entity.ErrResponse{
				SourceFunction: "VerifyUser",
				Detail:         fmt.Sprintf("error on retreiving user, detail = %v", err),
			})
			return
		}

		if user == nil {
			ctx.JSON(http.StatusBadRequest, entity.ErrResponse{
				SourceFunction: "VerifyUser",
				Detail:         "user not found",
			})
			return
		}

		_, err = repository.UpdateUser(user.UserID, &model.User{
			Status: "REGISTERED",
		})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, entity.ErrResponse{
				SourceFunction: "VerifyUser",
				Detail:         fmt.Sprintf("error on updating user, details = %v", err),
			})
			return
		}

		ctx.JSON(http.StatusOK, entity.BasicResponse{
			Message: "user registered",
		})
	}
}
