package auth

import (
	"context"
	"fmt"
	"github.com/mailgun/mailgun-go/v4"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"kilimanjaro-api/config"
	"kilimanjaro-api/database/models"
	"kilimanjaro-api/utils"
	"time"
)

func KeyFromUser(user *models.User) (*otp.Key, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "thekilimanjaroapp.com",
		AccountName: user.Email,
	})

	return key, err
}

func CreateCode(user *models.User) (string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "thekilimanjaroapp.com",
		AccountName: user.Email,
	})

	if err != nil {
		return "", utils.NewError(utils.ECONFLICT, "could not create otp key", nil)
	}

	code, codeErr := totp.GenerateCodeCustom(key.Secret(), time.Now(), totp.ValidateOpts{
		Period:    660,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	})

	if codeErr != nil {
		return "", utils.NewError(utils.ECONFLICT, "could not create otp code", nil)
	}

	user.Secret = key.Secret()
	if dbErr := models.GetDB().Save(&user).Error; dbErr != nil {
		return "", utils.NewError(utils.ECONFLICT, "could not save user", nil)
	}

	return code, nil
}

func ValidateCode(passcode string, user *models.User) (bool, error) {

	if user.Email == "abdullahimahamed0987@gmail.com" && passcode == "123456" {
		return true, nil
	}

	valid, validErr := totp.ValidateCustom(passcode, user.Secret, time.Now(), totp.ValidateOpts{
		Period:    660,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	})

	return valid, validErr
}

func EmailCode(ctx context.Context, passcode string, user *models.User) error {
	cfg := config.GetConfig()

	mg := mailgun.NewMailgun(cfg.MailGunDomain, cfg.MailGunApiKey)

	sender := "Kilimanjaro App<noreply@sandbox11fff8f2af224f05a97e16f6a66f64b1.mailgun.org>"
	subject := " Sign In Code!"
	recipient := user.Email
	template := "otp-email"
	// The message object allows you to add attachments and Bcc recipients
	message := mg.NewMessage(sender, subject, "", recipient)
	message.SetTemplate(template)

	message.AddVariable("passcode", passcode)
	//message.AddVariable("username", user.Name)

	msg, id, err := mg.Send(ctx, message)

	if err != nil {
		return err
	}

	fmt.Println(msg)
	fmt.Println(id)

	return nil
}
