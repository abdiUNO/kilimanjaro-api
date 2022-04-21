package products

import (
	"github.com/gorilla/mux"
	"kilimanjaro-api/utils/response"
	"log"
	"net/http"
)

var GetProduct = func(w http.ResponseWriter, r *http.Request) {
	//token := r.Context().Value("token").(*Token)
	params := mux.Vars(r)
	podId := params["id"]
	log.Println(podId)
	podcast, err := FindProductById(podId)
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

	products, err := FindAllProducts()

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

	products, err := QueryProduct(query)
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
