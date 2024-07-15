package repository

import (
	"fmt"

	"github.com/theborzet/time-tracker/internal/models"
)

func (r *ApiRepository) GetUserTasks(userID int, start, end string) ([]*models.TaskTimeSpent, error) {
	var tasks []*models.TaskTimeSpent
	query := `SELECT taskName, EXTRACT(EPOCH FROM (endTime - startTime))/60 as timeSpent 
              FROM tasks WHERE userId = $1 AND startTime >= $2 AND endTime <= $3 
              ORDER BY timeSpent DESC`
	rows, err := r.db.Query(query, userID, start, end)
	if err != nil {
		r.logger.Printf("Error querying user tasks: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task models.TaskTimeSpent
		if err := rows.Scan(&task.TaskName, &task.TimeSpent); err != nil {
			r.logger.Printf("Error scanning row: %v\n", err)
			return nil, err
		}
		tasks = append(tasks, &task)
	}
	if err := rows.Err(); err != nil {
		r.logger.Printf("Error after scanning rows: %v\n", err)
		return nil, err
	}

	return tasks, nil
}

func (r *ApiRepository) StartTask(userID int, taskName string, startTime string) error {
	_, err := r.db.Exec("INSERT INTO tasks (userId, taskName, startTime) VALUES ($1, $2, $3)", userID, taskName, startTime)
	if err != nil {
		r.logger.Printf("Error starting task: %v\n", err)
		return err
	}
	return nil
}

func (r *ApiRepository) EndTask(userID int, taskName string, endTime string) error {
	result, err := r.db.Exec("UPDATE tasks SET endTime = $1 WHERE userId = $2 AND taskName = $3 AND endTime IS NULL", endTime, userID, taskName)
	if err != nil {
		r.logger.Printf("Error ending task: %v\n", err)
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Printf("Error getting rows affected: %v\n", err)
		return err
	}
	if rowsAffected == 0 {
		r.logger.Printf("No rows were updated for task ending\n")
		return fmt.Errorf("no rows were updated")
	}
	return nil
}
