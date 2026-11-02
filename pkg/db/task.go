package db

import (
	"database/sql"
	"errors"
	"strconv"
	"time"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func Tasks(limit int, search string) ([]Task, error) {
	today := time.Now().Format("20060102")

	var rows *sql.Rows
	var err error

	if search == "" {
		rows, err = DB.Query(`SELECT id, date, title, comment, repeat FROM scheduler
			WHERE date >= ? ORDER BY date LIMIT ?`, today, limit)
	} else if t, e := time.Parse("02.01.2006", search); e == nil {
		rows, err = DB.Query(`SELECT id, date, title, comment, repeat FROM scheduler
			WHERE date = ? ORDER BY date LIMIT ?`, t.Format("20060102"), limit)
	} else {
		pattern := "%" + search + "%"
		rows, err = DB.Query(`SELECT id, date, title, comment, repeat FROM scheduler
			WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?`, pattern, pattern, limit)
	}
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tasks := []Task{}
	for rows.Next() {
		var t Task
		var id int64
		if err := rows.Scan(&id, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, err
		}
		t.ID = strconv.FormatInt(id, 10)
		tasks = append(tasks, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

func GetTask(id string) (*Task, error) {
	var t Task
	var id64 int64
	err := DB.QueryRow(`SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`, id).
		Scan(&id64, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		return nil, err
	}
	t.ID = strconv.FormatInt(id64, 10)
	return &t, nil
}

func UpdateTask(task *Task) error {
	res, err := DB.Exec(`UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`,
		task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("task not found")
	}
	return nil
}
func DeleteTask(id string) error {
	res, err := DB.Exec(`DELETE FROM scheduler WHERE id = ?`, id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("task not found")
	}
	return nil
}
