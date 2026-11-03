package api

import (
	"encoding/json"
	"final-prj/pkg/db"
	"net/http"
)

type TasksResp struct {
	Tasks []db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.FormValue("search")
	tasks, err := db.Tasks(50, search)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(TasksResp{Tasks: tasks})

}
