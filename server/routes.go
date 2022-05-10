package server

import (
	"encoding/json"
	"kilimanjaro-api/api/product"
	vendor "kilimanjaro-api/api/seller"
	"net/http"

	"kilimanjaro-api/api/auth"
)

func (s *Server) SetupRoutes() {
	s.router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		// an example API handler
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	s.router.HandleFunc("/users/new", auth.CreateUser).Methods("POST")
	s.router.HandleFunc("/users/login", auth.Authenticate).Methods("POST")
	s.router.HandleFunc("/users", auth.FindUsers).Queries("query", "{query}").Methods("GET")
	s.router.HandleFunc("/users/{id}", auth.UpdateUser).Methods("PATCH")

	s.router.HandleFunc("/users/{id}/otp-code", auth.GenerateOTP).Methods("GET")
	s.router.HandleFunc("/users/otp-code", auth.ValidateOTP).Methods("POST")

	//s.router.HandleFunc("/users/{id}/change-password", auth.ChangePassword).Methods("PATCH")

	//s.router.HandleFunc("/top_products", products.GetTopProducts).Methods("GET")

	s.router.HandleFunc("/vendors", vendor.GetAllVendors).Methods("GET")
	s.router.HandleFunc("/vendors/{id}", vendor.GetVendor).Methods("GET")
	s.router.HandleFunc("/vendors", vendor.CreateVendor).Methods("POST")

	s.router.HandleFunc("/products", product.GetAllProducts).Methods("GET")
	s.router.HandleFunc("/products/search", product.SearchProducts).Queries("q", "{q}").Methods("GET")
	s.router.HandleFunc("/products/{id}", product.GetProduct).Methods("GET")
	s.router.HandleFunc("/products", product.CreateProduct).Methods("POST")
	//s.router.HandleFunc("/users/{id}/otp-code", auth.GenerateOTP).Methods("GET")
	//s.router.HandleFunc("/users/{id}/otp-code", auth.ValidateOTP).Queries("code", "{code}").Methods("POST")
	//
	//s.router.HandleFunc("/friends", friends.GetFriends).Methods("GET")
	//s.router.HandleFunc("/block/{id}", friends.BlockUser).Methods("DELETE")
	//s.router.HandleFunc("/users/{id}/add", friends.AddFriend).Methods("PUT")
	//
	//s.router.HandleFunc("/friends/{id}/conversations", chats.CreateConversation).Methods("POST")
	//s.router.HandleFunc("/conversations", chats.GetConversations).Methods("GET")
	//s.router.HandleFunc("/conversations/{id}", chats.RemoveConversation).Methods("DELETE")
	//
	//s.router.HandleFunc("/chat/", cliques.GetGroups).Methods("GET")
	//s.router.HandleFunc("/chat/new", cliques.CreateGroup).Methods("POST")
	//s.router.HandleFunc("/chat/find", cliques.JoinGroup).Methods("POST")
	//s.router.HandleFunc("/chat/{id}/leave", cliques.LeaveGroup).Methods("PUT")
	//s.router.HandleFunc("/upload", auth.UploadProfileImage).Methods("POST")
}
