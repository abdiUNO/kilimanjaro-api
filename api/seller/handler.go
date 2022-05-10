package vendor

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"kilimanjaro-api/database/models"
	u "kilimanjaro-api/utils"
	"kilimanjaro-api/utils/response"
	"net/http"
)

var CreateVendor = func(w http.ResponseWriter, r *http.Request) {
	token := r.Context().Value("token").(*models.Token)
	user, err := models.FindUserById(token.UserId)

	vendor := &models.Vendor{}
	err = json.NewDecoder(r.Body).Decode(vendor) //decode the request body into struct and failed if any error occur
	if err != nil {
		response.HandleError(w, u.NewError(u.EINTERNAL, "Invalid request", err))
		return
	}

	vendor.Email = user.Email
	vendor.Phone = user.Phone

	vendor, err = vendor.Create()
	if err != nil {
		response.HandleError(w, u.NewError(u.EINTERNAL, "Internal server err", err))
		return
	}

	response.Json(w, map[string]interface{}{
		"vendor": vendor,
	})
}

var GetVendor = func(w http.ResponseWriter, r *http.Request) {
	//token := r.Context().Value("token").(*Token)
	params := mux.Vars(r)
	vendorID := params["id"]
	vendor, err := models.FindVendorById(vendorID)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	response.Json(w, map[string]interface{}{
		"vendor": vendor,
	})

}

var GetAllVendors = func(w http.ResponseWriter, r *http.Request) {
	//token := r.Context().Value("token").(*Token)
	vendors, err := models.FindAllVendors()
	if err != nil {
		response.HandleError(w, err)
		return
	}

	response.Json(w, map[string]interface{}{
		"vendors": vendors,
	})

}
