package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Tim73916/org-structure-api/internal/config"
	"github.com/Tim73916/org-structure-api/internal/database"
	"github.com/Tim73916/org-structure-api/internal/handlers"
	"github.com/Tim73916/org-structure-api/internal/handlers/middleware"
	"github.com/Tim73916/org-structure-api/internal/models"
	"github.com/Tim73916/org-structure-api/internal/repositories"
	"github.com/Tim73916/org-structure-api/internal/services"
)

func main() {
	cfg := config.Load()

	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := db.AutoMigrate(&models.Department{}, &models.Employee{}); err != nil {
		log.Fatal("Failed to migrate:", err)
	}

	deptRepo := repositories.NewDepartmentRepo(db)
	empRepo := repositories.NewEmployeeRepo(db)

	deptService := services.NewDepartmentService(deptRepo, empRepo, db)
	empService := services.NewEmployeeService(empRepo, deptRepo)

	deptHandler := handlers.NewDepartmentHandler(deptService)
	empHandler := handlers.NewEmployeeHandler(empService)

	mux := http.NewServeMux()

	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	mux.HandleFunc("/departments/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.Trim(r.URL.Path, "/")
		parts := strings.Split(path, "/")

		if len(parts) == 3 && parts[2] == "employees" && r.Method == http.MethodPost {
			empHandler.CreateEmployee(w, r)
			return
		}

		switch r.Method {
		case http.MethodPost:

			if len(parts) == 1 {
				deptHandler.CreateDepartment(w, r)
			} else {
				http.Error(w, "Not found", http.StatusNotFound)
			}
		case http.MethodGet:

			if len(parts) == 2 {
				deptHandler.GetDepartment(w, r)
			} else {
				http.Error(w, "Not found", http.StatusNotFound)
			}
		case http.MethodPatch:

			if len(parts) == 2 {
				deptHandler.UpdateDepartment(w, r)
			} else {
				http.Error(w, "Not found", http.StatusNotFound)
			}
		case http.MethodDelete:

			if len(parts) == 2 {
				deptHandler.DeleteDepartment(w, r)
			} else {
				http.Error(w, "Not found", http.StatusNotFound)
			}
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/employees", func(w http.ResponseWriter, r *http.Request) {
		path := strings.Trim(r.URL.Path, "/")
		parts := strings.Split(path, "/")

		switch r.Method {
		case http.MethodGet:

			if len(parts) == 1 {
				empHandler.ListEmployees(w, r)
			} else if len(parts) == 2 {

				empHandler.GetEmployee(w, r)
			} else {
				http.Error(w, "Not found", http.StatusNotFound)
			}
		case http.MethodPut:

			if len(parts) == 2 {
				empHandler.UpdateEmployee(w, r)
			} else {
				http.Error(w, "Not found", http.StatusNotFound)
			}
		case http.MethodDelete:

			if len(parts) == 2 {
				empHandler.DeleteEmployee(w, r)
			} else {
				http.Error(w, "Not found", http.StatusNotFound)
			}
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/employees/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.Trim(r.URL.Path, "/")
		parts := strings.Split(path, "/")

		switch r.Method {
		case http.MethodGet:
			if len(parts) == 1 {
				empHandler.ListEmployees(w, r)
			} else if len(parts) == 2 {
				empHandler.GetEmployee(w, r)
			} else {
				http.Error(w, "Not found", http.StatusNotFound)
			}
		case http.MethodPut:
			if len(parts) == 2 {
				empHandler.UpdateEmployee(w, r)
			} else {
				http.Error(w, "Not found", http.StatusNotFound)
			}
		case http.MethodDelete:
			if len(parts) == 2 {
				empHandler.DeleteEmployee(w, r)
			} else {
				http.Error(w, "Not found", http.StatusNotFound)
			}
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	handler := middleware.Logging(mux)
	handler = middleware.Recover(handler)

	srv := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("Server starting on port %s", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed:", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}
