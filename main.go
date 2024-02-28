package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jeypc/go-jwt-mux/controllers/authcontroller"
	"github.com/jeypc/go-jwt-mux/controllers/productcontroller"
	"github.com/jeypc/go-jwt-mux/controllers/rolescontroller"
	usersproductcontroller "github.com/jeypc/go-jwt-mux/controllers/users_productcontroller"
	usersrolecontroller "github.com/jeypc/go-jwt-mux/controllers/users_rolecontroller"
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
	r.HandleFunc("/read/user", authcontroller.ReadUser).Methods("GET")
	r.HandleFunc("/logout", authcontroller.Logout).Methods("GET")
	r.HandleFunc("/delete/user", authcontroller.DeleteUser).Methods("DELETE")
	r.HandleFunc("/users", authcontroller.UpdateUser).Methods("PUT").Queries("id", "{id}")

	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/products", data_product.Data).Methods("GET")
	api.Use(middlewares.JWTMiddleware)

	// Tambahkan rute untuk Roles
	role := r.PathPrefix("/roles").Subrouter()
	role.HandleFunc("/create", rolescontroller.CreateRole).Methods("POST")
	role.HandleFunc("/data", rolescontroller.GetRole).Methods("GET")
	role.HandleFunc("/edit", rolescontroller.UpdateRole).Methods("PUT")
	role.HandleFunc("/delete", rolescontroller.DeleteRole).Methods("DELETE")
	role.Use(middlewares.RoleAuthorizationMiddleware)

	// Route untuk Product
	product := r.PathPrefix("/product").Subrouter()
	product.HandleFunc("/create", productcontroller.CreateProduct).Methods("POST")
	product.HandleFunc("/data", productcontroller.GetProduct).Methods("GET")
	product.HandleFunc("/edit", productcontroller.UpdateProduct).Methods("PUT")
	product.HandleFunc("/delete", productcontroller.DeleteProduct).Methods("DELETE")

	// rute untuk userProduct
	userProduct := r.PathPrefix("/userProduct").Subrouter()
	userProduct.HandleFunc("/create", usersproductcontroller.CreateUserProduct).Methods("POST")
	userProduct.HandleFunc("/data", usersproductcontroller.ReadUserProduct).Methods("GET")
	userProduct.HandleFunc("/edit", usersproductcontroller.UpdateUserProduct).Methods("PUT")
	userProduct.HandleFunc("/delete", usersproductcontroller.DeleteUserProduct).Methods("DELETE")

	// rute untuk userProduct
	userRole := r.PathPrefix("/userRole").Subrouter()
	userRole.HandleFunc("/create", usersrolecontroller.CreateUserRole).Methods("POST")
	userRole.HandleFunc("/data", usersrolecontroller.ReaduserRole).Methods("GET")
	userRole.HandleFunc("/edit", usersrolecontroller.UpdateUserRole).Methods("PUT")
	userRole.HandleFunc("/delete", usersrolecontroller.DeleteuserRole).Methods("DELETE")

	r.HandleFunc("/verify-otp", authcontroller.VerifyOTP).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", r))
}
