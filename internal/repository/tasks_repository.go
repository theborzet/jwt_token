package repository

import (
	"github.com/theborzet/time-tracker/internal/models"
)

func (r *ApiRepository) GetUserTasks(userID int, start, end string) ([]*models.TaskTimeSpent, error) {
	var tasks []*models.TaskTimeSpent
	query := `SELECT taskName, EXTRACT(EPOCH FROM (endTime - startTime))/60 as timeSpent 
              FROM tasks WHERE userId = $1 AND startTime >= $2 AND endTime <= $3 
              ORDER BY timeSpent DESC`
	rows, err := r.db.Queryx(query, userID, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task models.TaskTimeSpent
		if err := rows.StructScan(&task); err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}

	return tasks, nil
}

func (r *ApiRepository) StartTask(userID int, taskName string, startTime string) error {
	_, err := r.db.Exec("INSERT INTO tasks (userId, taskName, startTime) VALUES ($1, $2, $3)", userID, taskName, startTime)
	return err
}

func (r *ApiRepository) EndTask(userID int, endTime string) error {
	_, err := r.db.Exec("UPDATE tasks SET endTime = $1 WHERE userId = $2 AND endTime IS NULL", endTime, userID)
	return err
}
