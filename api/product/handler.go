package product

import (
	"encoding/json"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"kilimanjaro-api/database/models"
	u "kilimanjaro-api/utils"
	"kilimanjaro-api/utils/response"

	"net/http"
)

var CreateProduct = func(w http.ResponseWriter, r *http.Request) {
	token := r.Context().Value("token").(*models.Token)
	logger := r.Context().Value("logger").(*log.Entry)

	user, _ := models.FindUserById(token.UserId)

	product := &models.Product{}
	product.VendorID = user.ID

	jsonErr := json.NewDecoder(r.Body).Decode(product) //decode the request body into struct and failed if any error occur
	if jsonErr != nil {
		logger.Errorln(jsonErr.Error())

		response.HandleError(w, u.NewError(u.EINTERNAL, "Invalid request", jsonErr))
		return
	}

	data, productErr := product.Create()

	data.Vendor = user

	if productErr != nil {
		logger.Errorln(productErr.Error())

		response.HandleError(w, u.NewError(u.EINTERNAL, "Internal server err", productErr))
		return
	}

	response.Json(w, map[string]interface{}{
		"product": data,
	})
}

var GetProduct = func(w http.ResponseWriter, r *http.Request) {
	//token := r.Context().Value("token").(*Token)
	params := mux.Vars(r)
	podId := params["id"]
	podcast, err := models.FindProductById(podId)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	response.Json(w, map[string]interface{}{
		"podcast": podcast,
	})

}

var GetAllProducts = func(w http.ResponseWriter, r *http.Request) {

	logger := r.Context().Value("logger").(*log.Entry)
	logger.Warn("GET ALL PRODUCTS")
	//user := auth.GetUser(token.UserId)

	products, err := models.FindAllProducts()

	if err != nil {
		response.HandleError(w, err)
		return
	}

	response.Json(w, map[string]interface{}{
		"products": products,
	})
}

var SearchProducts = func(w http.ResponseWriter, r *http.Request) {
	//token := r.Context().Value("token").(*Token)
	query := r.FormValue("q")

	products, err := models.QueryProduct(query)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	response.Json(w, map[string]interface{}{
		"products": products,
	})

}

//var GetTopProducts = func(w http.ResponseWriter, r *http.Request) {
//	//token := r.Context().Value("token").(*Token)
//	products, err := TopProducts()
//	if err != nil {
//		log.Println(err)
//		response.HandleError(w, err)
//		return
//	}
//
//	response.Json(w, map[string]interface{}{
//		"products": products,
//	})
//
//}
