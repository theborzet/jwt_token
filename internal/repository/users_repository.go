package repository

import (
	"fmt"

	"github.com/theborzet/time-tracker/internal/models"
)

func (r *ApiRepository) GetUsers(filter map[string]string) ([]*models.User, error) {
	var users []*models.User

	query := "SELECT id, passportNumber, passportSerie, surname, name, patronymic, address FROM users WHERE 1=1"
	args := make([]interface{}, 0)

	for k, v := range filter {
		query += fmt.Sprintf(" AND %s = ?", k)
		args = append(args, v)
	}

	rows, err := r.db.Queryx(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.PassportNumber, &user.PassportSerie, &user.Surname, &user.Name, &user.Patronymic, &user.Address)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *ApiRepository) GetUserByID(id int) (*models.User, error) {
	var user models.User
	err := r.db.Get(&user, "SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *ApiRepository) CreateUser(user *models.User) error {
	_, err := r.db.Exec(`INSERT INTO users (passportNumber, passportSerie, surname, name, patronymic, address) 
                         VALUES ($1, $2, $3, $4, $5, $6)`,
		user.PassportNumber, user.PassportSerie, user.Surname, user.Name, user.Patronymic, user.Address)
	return err
}

func (r *ApiRepository) UpdateUser(user *models.User) error {
	_, err := r.db.Exec(`UPDATE users SET passportNumber=$1, passportSerie=$2, surname=$3, name=$4, patronymic=$5, address=$6 
                         WHERE id=$7`,
		user.PassportNumber, user.PassportSerie, user.Surname, user.Name, user.Patronymic, user.Address, user.ID)
	return err
}

func (r *ApiRepository) DeleteUser(id int) error {
	_, err := r.db.Exec("DELETE FROM users WHERE id = $1", id)
	return err
}
