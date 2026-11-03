package api

import (
	"encoding/json"
	"errors"
	"final-prj/pkg/db"
	"net/http"
	"time"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}
	if err := checkTask(&task); err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}
	id, err := db.AddTask(task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int64{"id": id})
}

func checkTask(task *db.Task) error {
	now := time.Now()
	today := now.Format(DateLayout)

	if task.Title == "" {
		return errors.New("title is empty")
	}
	if task.Date == "" {
		task.Date = today
	}
	if _, err := time.Parse(DateLayout, task.Date); err != nil {
		return err
	}
	if task.Repeat != "" {
		next, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
		if task.Date < today {
			task.Date = next
		}
	} else if task.Date < today {
		task.Date = today
	}
	return nil
}

func writeError(w http.ResponseWriter, err error, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeError(w, errors.New("id is empty"), http.StatusBadRequest)
		return
	}
	task, err := db.GetTask(id)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			writeError(w, err, http.StatusNotFound)
			return
		}
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func putTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}
	if task.ID == "" {
		writeError(w, errors.New("id is empty"), http.StatusBadRequest)
		return
	}
	if err := checkTask(&task); err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}
	if err := db.UpdateTask(&task); err != nil {
		if errors.Is(err, db.ErrNotFound) {
			writeError(w, err, http.StatusNotFound)
			return
		}
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct{}{})
}

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeError(w, errors.New("id is empty"), http.StatusBadRequest)
		return
	}
	task, err := db.GetTask(id)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			writeError(w, err, http.StatusNotFound)
			return
		}
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			if errors.Is(err, db.ErrNotFound) {
				writeError(w, err, http.StatusNotFound)
				return
			}
			writeError(w, err, http.StatusInternalServerError)
			return
		}
	} else {
		next, err := NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			writeError(w, err, http.StatusInternalServerError)
			return
		}
		task.Date = next
		if err := db.UpdateTask(task); err != nil {
			if errors.Is(err, db.ErrNotFound) {
				writeError(w, err, http.StatusNotFound)
				return
			}
			writeError(w, err, http.StatusInternalServerError)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct{}{})
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeError(w, errors.New("id is empty"), http.StatusBadRequest)
		return
	}
	if err := db.DeleteTask(id); err != nil {
		if errors.Is(err, db.ErrNotFound) {
			writeError(w, err, http.StatusNotFound)
			return
		}
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct{}{})
}
