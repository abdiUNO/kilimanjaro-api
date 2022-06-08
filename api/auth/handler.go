package auth

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"kilimanjaro-api/database/models"
	u "kilimanjaro-api/utils"
	"kilimanjaro-api/utils/response"
	"net/http"
)

var CreateUser = func(w http.ResponseWriter, r *http.Request) {

	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user) //decode the request body into struct and failed if any error occur
	if err != nil {
		response.HandleError(w, u.NewError(u.EINTERNAL, "Invalid request", err))
		return
	}

	if validErr := user.Validate(); validErr != nil {
		response.HandleError(w, validErr)
		return
	}

	data, ormErr := user.Create()
	if ormErr != nil {
		response.HandleError(w, u.NewError(u.EINTERNAL, "Internal server err", ormErr))
		return
	}

	response.Json(w, map[string]interface{}{
		"user": data,
	})

}

var GetVendorProducts = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId := params["id"]

	products, err := models.FindProductsByVendor(userId)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	response.Json(w, map[string]interface{}{
		"products": products,
	})
}

//find user by email
//   if user found, return user
//create user if user not found
//   if user created, return user and token

var Authenticate = func(w http.ResponseWriter, r *http.Request) {
	logger := r.Context().Value("logger").(*log.Entry)

	user := &models.User{}
	//decode the request body into struct and failed if any error occur
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		log.Debug(err.Error())
		response.HandleError(w, u.NewError(u.EINTERNAL, "Invalid request", err))
		return
	}

	data, err := models.Login(user.Email)
	if err != nil {
		if data.Email == "" {
			data, err = user.Create()
			if err != nil {
				log.Error(err)
				response.HandleError(w, u.NewError(u.EINTERNAL, "Internal server err", err))
				return
			}
		} else {
			response.HandleError(w, err)
			return
		}
	}

	if data != nil {
		logger.Info(data.Email)
		code, err := CreateCode(data)
		logger.Info(code)
		if err != nil {
			log.Error(err)
			response.HandleError(w, u.NewError(u.EINTERNAL, "could not create code", err))
			return
		}

		logger.Info(data.UserType)
		err = EmailCode(r.Context(), code, data)
		if err != nil {
			fmt.Println(err.Error())
			response.HandleError(w, u.NewError(u.EINTERNAL, "could not send otp email", err))
			return
		}
	}

	response.Json(w, map[string]interface{}{
		"user": data,
	})

}

var UpdateUser = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId := params["id"]
	user := &models.User{}
	//decode the request body into struct and failed if any error occur
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		response.HandleError(w, u.NewError(u.EINTERNAL, "Invalid request", err))
		return
	}

	user, err := models.Update(userId, user)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	response.Json(w, map[string]interface{}{
		"user": user,
	})

}

var FindUsers = func(w http.ResponseWriter, r *http.Request) {
	token := r.Context().Value("token").(*models.Token)
	query := r.FormValue("query")

	users, err := models.QueryUsers(token.UserId, query)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	response.Json(w, map[string]interface{}{
		"users": users,
	})

}

var GenerateOTP = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId := params["id"]
	user, dbErr := models.FindUserById(userId)

	if dbErr != nil {
		response.HandleError(w, u.NewError(u.ENOTFOUND, "could not find user", dbErr))
		return
	}

	code, err := CreateCode(user)
	if err != nil {
		response.HandleError(w, u.NewError(u.EINTERNAL, "could not create code", err))
		return
	}

	err = EmailCode(r.Context(), code, user)
	if err != nil {
		fmt.Println(err.Error())
		response.HandleError(w, u.NewError(u.EINTERNAL, "could not send otp email", err))
		return
	}

	response.Json(w, map[string]interface{}{
		"codeSent": true,
	})
}

type ValidateRequest struct {
	Code   string `json:",code"`
	UserID string `json:",userId"`
}

var ValidateOTP = func(w http.ResponseWriter, r *http.Request) {
	logger := r.Context().Value("logger").(*log.Entry)

	formData := &ValidateRequest{}
	if err := json.NewDecoder(r.Body).Decode(formData); err != nil {
		log.Debug(err.Error())
		response.HandleError(w, u.NewError(u.EINTERNAL, "Invalid request", err))
		return
	}

	passcode := formData.Code
	userId := formData.UserID

	logger.Info(userId)

	user, dbErr := models.FindUserById(userId)

	logger.Info(passcode)

	if dbErr != nil {
		response.HandleError(w, u.NewError(u.ENOTFOUND, "could not find user", dbErr))
		return
	}

	isValid, err := ValidateCode(passcode, user)

	if err != nil {
		response.HandleError(w, u.NewError(u.EINTERNAL, "could not validate code", err))
		return
	}

	if isValid == true {

		user.JwtToken = models.GenerateJwtToken(user)

		if user.EmailVerified == false {
			user.EmailVerified = true
			dbErr := models.GetDB().Save(&user).Error

			if dbErr != nil {
				response.HandleError(w, u.NewError(u.EINTERNAL, "could not update user", nil))
				return
			}
		}
	}

	response.Json(w, map[string]interface{}{
		"isValid": isValid,
		"user":    user,
	})
}
