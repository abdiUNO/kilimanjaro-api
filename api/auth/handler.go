package auth

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	u "kilimanjaro-api/utils"
	"kilimanjaro-api/utils/response"
	"net/http"
)

var CreateUser = func(w http.ResponseWriter, r *http.Request) {

	user := &User{}
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

var Authenticate = func(w http.ResponseWriter, r *http.Request) {
	user := &User{}
	//decode the request body into struct and failed if any error occur
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		response.HandleError(w, u.NewError(u.EINTERNAL, "Invalid request", err))
		return
	}

	data, err := Login(user.Email, user.Password)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	response.Json(w, map[string]interface{}{
		"user": data,
	})

}

var UpdateUser = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId := params["id"]
	user := &User{}
	//decode the request body into struct and failed if any error occur
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		response.HandleError(w, u.NewError(u.EINTERNAL, "Invalid request", err))
		return
	}

	user, err := Update(userId, user)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	response.Json(w, map[string]interface{}{
		"user": user,
	})

}

type ChangePasswordBody struct {
	OldPassword string `json:",oldPassword"`
	NewPassword string `json:",newPassword"`
}

var ChangePassword = func(w http.ResponseWriter, r *http.Request) {
	token := r.Context().Value("token").(*Token)
	user := GetUser(token.UserId)

	jsonBody := &ChangePasswordBody{}
	//decode the request body into struct and failed if any error occur
	if err := json.NewDecoder(r.Body).Decode(jsonBody); err != nil {
		response.HandleError(w, u.NewError(u.EINTERNAL, "Invalid request", err))
		return
	}

	updateErr := user.UpdatePassword(jsonBody.OldPassword, jsonBody.NewPassword)

	if updateErr != nil {
		response.HandleError(w, updateErr)
		return
	}

	response.Json(w, map[string]interface{}{
		"data": "Updated user password",
	})

}

var FindUsers = func(w http.ResponseWriter, r *http.Request) {
	token := r.Context().Value("token").(*Token)
	query := r.FormValue("query")

	users, err := QueryUsers(token.UserId, query)
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
	user, dbErr := FindUserById(userId)

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

var ValidateOTP = func(w http.ResponseWriter, r *http.Request) {
	passcode := r.FormValue("code")
	params := mux.Vars(r)
	userId := params["id"]
	user, dbErr := FindUserById(userId)

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
		user.EmailVerified = true
		dbErr := GetDB().Save(&user).Error

		if dbErr != nil {
			response.HandleError(w, u.NewError(u.EINTERNAL, "could not update user", nil))
			return
		}
	}

	response.Json(w, map[string]interface{}{
		"isValid": isValid,
	})
}
