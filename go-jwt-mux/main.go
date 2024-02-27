package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jeypc/go-jwt-mux/controllers/authcontroller"
	"github.com/jeypc/go-jwt-mux/controllers/productcontroller"
	"github.com/jeypc/go-jwt-mux/controllers/rolescontroller"
	"github.com/jeypc/go-jwt-mux/controllers/userRole"
	"github.com/jeypc/go-jwt-mux/data_product"
	"github.com/jeypc/go-jwt-mux/middlewares"
	"github.com/jeypc/go-jwt-mux/models"


	
)

func main() {
	models.ConnectDatabase()
	r := mux.NewRouter()

	// Route untuk User
	r.HandleFunc("/login", authcontroller.Login).Methods("POST")
	r.HandleFunc("/register", authcontroller.Register).Methods("POST")
	r.HandleFunc("/data/user", authcontroller.ReadUser).Methods("GET")
	r.HandleFunc("/logout", authcontroller.Logout).Methods("GET")
	r.HandleFunc("/delete/user", authcontroller.DeleteUser).Methods("DELETE")
	r.HandleFunc("/users", authcontroller.UpdateUser).Methods("PUT").Queries("id", "{id}")

	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/products", data_product.Data).Methods("GET")
	api.Use(middlewares.JWTMiddleware)

	// Tambahkan rute untuk Roles menggunakan Gorilla Mux
	role := r.PathPrefix("/roles").Subrouter()
	role.HandleFunc("/create", rolescontroller.CreateRole).Methods("POST")
	role.HandleFunc("/data", rolescontroller.GetRole).Methods("GET")
	role.HandleFunc("/edit", rolescontroller.UpdateRole).Methods("PUT")
	role.HandleFunc("/delete", rolescontroller.DeleteRole).Methods("DELETE")

	// Route untuk Product
	product := r.PathPrefix("/product").Subrouter()
	product.HandleFunc("/create", productcontroller.CreateProduct).Methods("POST")
	product.HandleFunc("/data", productcontroller.GetProduct).Methods("GET")
	product.HandleFunc("/edit", productcontroller.UpdateProduct).Methods("PUT")
	product.HandleFunc("/delete", productcontroller.DeleteProduct).Methods("DELETE")

	// Tambahkan rute untuk User Roles
	role = r.PathPrefix("/user/roles").Subrouter()
	role.HandleFunc("/create", userrolecontroller.CreateRole).Methods("POST")
	role.HandleFunc("/data", userrolecontroller.GetRole).Methods("GET")
	role.HandleFunc("/edit", userrolecontroller.UpdateRole).Methods("PUT")
	role.HandleFunc("/delete", userrolecontroller.DeleteRole).Methods("DELETE")


	log.Fatal(http.ListenAndServe(":8080", r))
}
