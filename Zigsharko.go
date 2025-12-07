package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Discription string `json:"discription,omitempty`
	Completed   bool   `json:"completed"`
}

var (
	tasks  = make(map[int]Task)
	nextId = 1
	mu     sync.Mutex
)

func getTask(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	var t []Task
	for _, ts := range tasks {
		t = append(t, ts)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}
func delTask(w http.ResponseWriter, r *http.Request) {

	idstr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idstr)

	_, exists := tasks[id]
	mu.Lock()
	defer mu.Unlock()
	if !exists {
		http.Error(w, "Id is Not Found ", http.StatusNotFound)
		return
	}

	delete(tasks, id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("ID successfully Deleted")

}
func updateTask(w http.ResponseWriter, r *http.Request) {
	var t Task
	idstr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idstr)

	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "update field is not yet oops somethink is wrong", http.StatusNoContent)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	_, exists := tasks[id]
	if !exists {
		http.Error(w, "Update process failed , Task Id not exists", http.StatusNotFound)
		return
	}
	t.ID = id
	tasks[id] = t
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks[id])

}

func createTask(w http.ResponseWriter, r *http.Request) {
	var t Task

	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "json data not validate ", http.StatusBadRequest)
		return
	}

	mu.Lock()
	t.ID = nextId
	nextId++
	tasks[t.ID] = t
	mu.Unlock()
	fmt.Println(t)
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Contend-Type", "application/json")
	json.NewEncoder(w).Encode(t)

}
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			createTask(w, r)
		case http.MethodPut:
			updateTask(w, r)
		case http.MethodDelete:
			delTask(w, r)
		case http.MethodGet:
			getTask(w, r)
		}
	})
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		log.Println("Server started Listing port on :[8080]")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
