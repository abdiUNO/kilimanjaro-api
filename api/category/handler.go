package category

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"kilimanjaro-api/database/models"
	u "kilimanjaro-api/utils"
	"kilimanjaro-api/utils/response"
	"log"
	"net/http"
)

var CreateCategory = func(w http.ResponseWriter, r *http.Request) {
	category := &models.Category{}

	jsonErr := json.NewDecoder(r.Body).Decode(category) //decode the request body into struct and failed if any error occur
	if jsonErr != nil {
		response.HandleError(w, u.NewError(u.EINTERNAL, "Invalid request", jsonErr))
		return
	}

	data, categoryErr := category.Create()

	if categoryErr != nil {
		response.HandleError(w, u.NewError(u.EINTERNAL, "Internal server err", categoryErr))
		return
	}

	response.Json(w, map[string]interface{}{
		"category": data,
	})
}

var GetCategory = func(w http.ResponseWriter, r *http.Request) {
	//token := r.Context().Value("token").(*Token)
	params := mux.Vars(r)
	categoryId := params["id"]
	log.Println(categoryId)
	category, err := models.FindCategoryById(categoryId)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	response.Json(w, map[string]interface{}{
		"category": category,
	})

}

var GetAllCategories = func(w http.ResponseWriter, r *http.Request) {
	categories, err := models.FindAllCategories()
	if err != nil {
		response.HandleError(w, err)
		return
	}

	response.Json(w, map[string]interface{}{
		"categories": categories,
	})

}
