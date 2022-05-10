package product

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"kilimanjaro-api/database/models"
	u "kilimanjaro-api/utils"
	"kilimanjaro-api/utils/response"
	"log"
	"net/http"
)

var CreateProduct = func(w http.ResponseWriter, r *http.Request) {
	product := &models.Product{}
	err := json.NewDecoder(r.Body).Decode(product) //decode the request body into struct and failed if any error occur
	if err != nil {
		response.HandleError(w, u.NewError(u.EINTERNAL, "Invalid request", err))
		return
	}

	log.Println("TEST")

	vendor, _ := models.FindVendorById(product.VendorID)

	data, productErr := product.Create()
	data.Vendor = vendor

	if productErr != nil {
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
	log.Println(podId)
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

	//token := r.Context().Value("token").(*auth.Token)
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
