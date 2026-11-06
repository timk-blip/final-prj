package api

import (
	"encoding/json"
	"final-prj/pkg/db"
	"log"
	"net/http"
)

type TasksResp struct {
	Tasks []db.Task `json:"tasks"`
}

const tasksLimit = 50

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.FormValue("search")
	tasks, err := db.Tasks(tasksLimit, search)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(TasksResp{Tasks: tasks}); err != nil {
		log.Printf("encode response: %v", err)
	}
}
