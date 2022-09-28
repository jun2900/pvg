package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"pvg/entity"
	"pvg/model"
	"pvg/repository"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
)

func SendMailChangePassword() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var input entity.RecovEmail

		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, entity.ErrResponse{
				SourceFunction: "SendMailChangePassword",
				Detail:         fmt.Sprintf("error on parsing input, detail = %v", err),
			})
			return
		}

		user, _ := repository.GetUserByEmail(input.Email)

		if user != nil {
			randomToken := GenerateSecureToken(10)
			hashedPwd, err := bcrypt.GenerateFromPassword([]byte(randomToken), bcrypt.DefaultCost)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, entity.ErrResponse{
					SourceFunction: "SendMailChangePassword",
					Detail:         fmt.Sprintf("error on generating new pwd, detail = %v", err),
				})
				return
			}

			var updatedUser model.User
			updatedUser.Password = string(hashedPwd)

			_, err = repository.UpdateUser(user.UserID, &updatedUser)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, entity.ErrResponse{
					SourceFunction: "SendMailChangePassword",
					Detail:         fmt.Sprintf("error on updating user, detail = %v", err),
				})
				return
			}

			bodyText := fmt.Sprintf(`Hello, here is your updated password. %s`, randomToken)
			header := "Change Password"
			if err := SendMail(bodyText, header, input.Email); err != nil {
				ctx.JSON(http.StatusInternalServerError, entity.ErrResponse{
					SourceFunction: "SendMailChangePassword",
					Detail:         fmt.Sprintf("error on sending mail, detail = %v", err),
				})
				return
			}
		}

		ctx.JSON(http.StatusOK, entity.BasicResponse{
			Message: "mail sent",
		})
	}
}

func SendMail(bodyText string, header string, emailTo string) (err error) {
	var CONFIG_SMTP_HOST = os.Getenv("CONFIG_SMTP_HOST")
	var CONFIG_SMTP_PORT = os.Getenv("CONFIG_SMTP_PORT")
	var CONFIG_SENDER_NAME = os.Getenv("CONFIG_SENDER_NAME")
	var CONFIG_AUTH_EMAIL = os.Getenv("CONFIG_AUTH_EMAIL")
	var CONFIG_AUTH_PASSWORD = os.Getenv("CONFIG_AUTH_PASSWORD")

	smtpPort, _ := strconv.Atoi(CONFIG_SMTP_PORT)

	dialer := gomail.NewDialer(
		CONFIG_SMTP_HOST,
		smtpPort,
		CONFIG_AUTH_EMAIL,
		CONFIG_AUTH_PASSWORD,
	)

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", CONFIG_SENDER_NAME)
	mailer.SetHeader("To", emailTo)
	mailer.SetHeader("Subject", header)
	//mailer.Embed("./unosign-logo.png")
	mailer.SetBody("text/html", bodyText)

	if err = dialer.DialAndSend(mailer); err != nil {
		return err
	}

	return nil
}

func GenerateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
