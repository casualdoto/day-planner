package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"day-planner/internal/bridge"
	"day-planner/internal/models"
)

func RunServer(b *bridge.Bridge) error {
	addr := os.Getenv("DAY_PLANNER_ADDR")
	if addr == "" {
		addr = "127.0.0.1:8080"
	}

	mux := http.NewServeMux()
	registerAPI(mux, b)
	registerFrontend(mux)

	log.Printf("Day Planner backend listening on http://%s", addr)
	return http.ListenAndServe(addr, withCORS(mux))
}

func registerAPI(mux *http.ServeMux, b *bridge.Bridge) {
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, map[string]string{"status": "ok"})
	})

	mux.HandleFunc("GET /api/day", func(w http.ResponseWriter, r *http.Request) {
		result, err := b.GetDayEntry(r.URL.Query().Get("date"))
		writeResult(w, result, err)
	})
	mux.HandleFunc("POST /api/day/plan", func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			Date     string `json:"date"`
			PlanText string `json:"planText"`
		}
		if !readBody(w, r, &input) {
			return
		}
		writeResult(w, nil, b.SaveDayPlan(input.Date, input.PlanText))
	})
	mux.HandleFunc("POST /api/day/result", func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			Date       string `json:"date"`
			ResultText string `json:"resultText"`
		}
		if !readBody(w, r, &input) {
			return
		}
		writeResult(w, nil, b.SaveDayResult(input.Date, input.ResultText))
	})

	mux.HandleFunc("POST /api/tasks/create", func(w http.ResponseWriter, r *http.Request) {
		var input models.CreateTaskInput
		if !readBody(w, r, &input) {
			return
		}
		result, err := b.CreateTask(input)
		writeResult(w, result, err)
	})
	mux.HandleFunc("POST /api/tasks/update", func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			ID    int64                  `json:"id"`
			Input models.UpdateTaskInput `json:"input"`
		}
		if !readBody(w, r, &input) {
			return
		}
		result, err := b.UpdateTask(input.ID, input.Input)
		writeResult(w, result, err)
	})
	mux.HandleFunc("POST /api/tasks/delete", func(w http.ResponseWriter, r *http.Request) {
		var input idInput
		if !readBody(w, r, &input) {
			return
		}
		writeResult(w, nil, b.DeleteTask(input.ID))
	})
	mux.HandleFunc("POST /api/tasks/toggle", func(w http.ResponseWriter, r *http.Request) {
		var input idInput
		if !readBody(w, r, &input) {
			return
		}
		result, err := b.ToggleTaskCompleted(input.ID)
		writeResult(w, result, err)
	})
	mux.HandleFunc("POST /api/tasks/list", func(w http.ResponseWriter, r *http.Request) {
		var filter models.TaskFilter
		if !readBody(w, r, &filter) {
			return
		}
		result, err := b.ListTasks(filter)
		writeResult(w, result, err)
	})

	mux.HandleFunc("GET /api/tags", func(w http.ResponseWriter, r *http.Request) {
		result, err := b.ListTags()
		writeResult(w, result, err)
	})
	mux.HandleFunc("POST /api/tags/create", func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			Name  string `json:"name"`
			Color string `json:"color"`
		}
		if !readBody(w, r, &input) {
			return
		}
		result, err := b.CreateTag(input.Name, input.Color)
		writeResult(w, result, err)
	})
	mux.HandleFunc("POST /api/tags/update", func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			ID    int64  `json:"id"`
			Name  string `json:"name"`
			Color string `json:"color"`
		}
		if !readBody(w, r, &input) {
			return
		}
		result, err := b.UpdateTag(input.ID, input.Name, input.Color)
		writeResult(w, result, err)
	})
	mux.HandleFunc("POST /api/tags/delete", func(w http.ResponseWriter, r *http.Request) {
		var input idInput
		if !readBody(w, r, &input) {
			return
		}
		writeResult(w, nil, b.DeleteTag(input.ID))
	})

	mux.HandleFunc("POST /api/calendar/create", func(w http.ResponseWriter, r *http.Request) {
		var input models.CreateCalendarEventInput
		if !readBody(w, r, &input) {
			return
		}
		result, err := b.CreateCalendarEvent(input)
		writeResult(w, result, err)
	})
	mux.HandleFunc("POST /api/calendar/update", func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			ID    int64                           `json:"id"`
			Input models.CreateCalendarEventInput `json:"input"`
		}
		if !readBody(w, r, &input) {
			return
		}
		result, err := b.UpdateCalendarEvent(input.ID, input.Input)
		writeResult(w, result, err)
	})
	mux.HandleFunc("POST /api/calendar/delete", func(w http.ResponseWriter, r *http.Request) {
		var input idInput
		if !readBody(w, r, &input) {
			return
		}
		writeResult(w, nil, b.DeleteCalendarEvent(input.ID))
	})
	mux.HandleFunc("POST /api/calendar/toggle", func(w http.ResponseWriter, r *http.Request) {
		var input idInput
		if !readBody(w, r, &input) {
			return
		}
		result, err := b.ToggleCalendarEventCompleted(input.ID)
		writeResult(w, result, err)
	})
	mux.HandleFunc("GET /api/calendar", func(w http.ResponseWriter, r *http.Request) {
		result, err := b.ListCalendarEvents(r.URL.Query().Get("startDate"), r.URL.Query().Get("endDate"))
		writeResult(w, result, err)
	})
}

func registerFrontend(mux *http.ServeMux) {
	dist := filepath.Join("frontend", "dist")
	index := filepath.Join(dist, "index.html")
	if _, err := os.Stat(index); err != nil {
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Frontend dev server should run at http://localhost:5173", http.StatusNotFound)
		})
		return
	}
	files := http.FileServer(http.Dir(dist))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Join(dist, filepath.Clean(r.URL.Path))
		if info, err := os.Stat(path); err == nil && !info.IsDir() {
			files.ServeHTTP(w, r)
			return
		}
		http.ServeFile(w, r, index)
	})
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "http://localhost:5173" || origin == "http://127.0.0.1:5173" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

type idInput struct {
	ID int64 `json:"id"`
}

func readBody(w http.ResponseWriter, r *http.Request, target any) bool {
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(target); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return false
	}
	return true
}

func writeResult(w http.ResponseWriter, result any, err error) {
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, result)
}

func writeJSON(w http.ResponseWriter, value any) {
	w.Header().Set("Content-Type", "application/json")
	if value == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if err := json.NewEncoder(w).Encode(value); err != nil {
		writeError(w, http.StatusInternalServerError, err)
	}
}

func writeError(w http.ResponseWriter, status int, err error) {
	if err == nil {
		err = errors.New("unknown error")
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprint(err)})
}
