package repository

import (
	"fmt"

	"github.com/theborzet/time-tracker/internal/models"
)

func (r *ApiRepository) GetUsers(filters map[string]string) ([]*models.User, error) {
	var users []*models.User

	query := "SELECT id, passportNumber, passportSerie, surname, name, patronymic, address FROM users WHERE 1=1"
	args := make([]interface{}, 0)
	argCount := 1

	for k, v := range filters {
		query += fmt.Sprintf(" AND %s LIKE $%d", k, argCount)
		args = append(args, "%"+v+"%")
		argCount++
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		r.logger.Printf("Error querying users: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.PassportNumber, &user.PassportSerie, &user.Surname, &user.Name, &user.Patronymic, &user.Address)
		if err != nil {
			r.logger.Printf("Error scanning user row: %v", err)
			return nil, err
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		r.logger.Printf("Error iterating over user rows: %v", err)
		return nil, err
	}

	return users, nil
}

func (r *ApiRepository) CreateUser(user *models.User) error {
	_, err := r.db.Exec(`INSERT INTO users (passportNumber, passportSerie, surname, name, patronymic, address) 
                         VALUES ($1, $2, $3, $4, $5, $6)`,
		user.PassportNumber, user.PassportSerie, user.Surname, user.Name, user.Patronymic, user.Address)
	if err != nil {
		r.logger.Printf("Error creating user: %v", err)
	}
	return err
}

func (r *ApiRepository) UpdateUser(user *models.User) error {
	result, err := r.db.Exec(`UPDATE users SET passportNumber=$1, passportSerie=$2, surname=$3, name=$4, patronymic=$5, address=$6 
                         WHERE id=$7`,
		user.PassportNumber, user.PassportSerie, user.Surname, user.Name, user.Patronymic, user.Address, user.ID)
	if err != nil {
		r.logger.Printf("Error updating user: %v", err)
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
	return err
}

func (r *ApiRepository) DeleteUser(id int) error {
	_, err := r.db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		r.logger.Printf("Error deleting user: %v", err)
	}
	return err
}
