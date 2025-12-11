package main

import (
	// "context"
	"encoding/json"
	// "log"
	// "net/http"
	// "os"
	// "os/signal"
	// "time"

	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
)
var validate = validator.New()
type User struct{
	Id 	  string  `json:"id"`
	Name  string  `json:"username" validate:"required,min=3"`
	Email string  `json:"email" validate:"required,email"`
}

var users = map[string]User{}

func jsonContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
func welcomePage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello welcome to the page"))
}

func getUsers(w http.ResponseWriter , r *http.Request){
	userList := make([]User,0,len(users))
	for _ , user :=  range users{
		userList = append(userList,user)
	}
	writeRes(w,200,userList)
}

func writeRes(w http.ResponseWriter ,code int,data interface{}){
	w.WriteHeader(code)
  	json.NewEncoder(w).Encode(map[string]interface{}{
    "data": data,
  })
}
func writeJSONError(w http.ResponseWriter, code int, message, detail string) {
  w.WriteHeader(code)
  json.NewEncoder(w).Encode(map[string]interface{}{
    "error": map[string]interface{}{
      "code":    code,
      "message": message,
      "details": detail,
    },
  })
}
func setUser(w http.ResponseWriter , r *http.Request){
	var u User

	if err:=json.NewDecoder(r.Body).Decode(&u);err!=nil{
		writeJSONError(w,http.StatusBadRequest,"Invalidate user data",err.Error())
		return
	}

	if err := validate.Struct(&u);err!=nil{
		writeJSONError(w,http.StatusBadRequest,"Data Missing",err.Error())
		return
	}

	_ ,exists := users[u.Id]
	if exists {
		writeJSONError(w,http.StatusBadRequest,"user Id already in","id")
		return
	}

	users[u.Id] = u
	writeRes(w,200,u)
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(jsonContentType)

	r.Route("/v1/users", func(r chi.Router) {
		r.Get("/", welcomePage)
		r.Get("/alluser",getUsers)
		r.Post("/setUser",setUser)

	})
	srv := &http.Server{
		Addr:    ":3434",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen : %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("server shutting down forcely")
	}
	log.Printf("server shutdown ")
}
