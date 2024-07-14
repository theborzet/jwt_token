package repository

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/theborzet/time-tracker/internal/models"
)

func TestGetUsers(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewApiRepository(sqlxDB)

	expectedUsers := []*models.User{
		{ID: 1, PassportNumber: "1234", PassportSerie: "567890", Surname: "Иванов", Name: "Иван", Patronymic: "Иванович", Address: "г. Москва"},
		{ID: 2, PassportNumber: "9877", PassportSerie: "543210", Surname: "Иванов", Name: "Петр", Patronymic: "Петрович", Address: "г. Санкт-Петербург"},
	}
	filter := map[string]string{"surname": "Иванов"}

	rows := sqlmock.NewRows([]string{"id", "passportNumber", "passportSerie", "surname", "name", "patronymic", "address"})
	for _, user := range expectedUsers {
		rows.AddRow(user.ID, user.PassportNumber, user.PassportSerie, user.Surname, user.Name, user.Patronymic, user.Address)
	}

	mock.ExpectQuery(`SELECT id, passportNumber, passportSerie, surname, name, patronymic, address FROM users WHERE 1=1 AND surname = \?$`).
		WithArgs("Иванов").
		WillReturnRows(rows)

	users, err := repo.GetUsers(filter)
	assert.NoError(t, err)
	assert.NotNil(t, users)
	assert.Equal(t, len(expectedUsers), len(users)) // Проверяем количество возвращенных пользователей

	for i := range users {
		assert.Equal(t, expectedUsers[i].Surname, users[i].Surname)
		assert.Equal(t, expectedUsers[i].PassportNumber, users[i].PassportNumber)
		assert.Equal(t, expectedUsers[i].PassportSerie, users[i].PassportSerie)
		assert.Equal(t, expectedUsers[i].Name, users[i].Name)
		assert.Equal(t, expectedUsers[i].Patronymic, users[i].Patronymic)
		assert.Equal(t, expectedUsers[i].Address, users[i].Address)

	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
func TestGetUserByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewApiRepository(sqlxDB)

	expectedUser := &models.User{ID: 1, PassportNumber: "1234", PassportSerie: "233445", Surname: "Иванов", Name: "Иван", Patronymic: "Иванович", Address: "г. Москва"}

	rows := sqlmock.NewRows([]string{"id", "passportNumber", "passportSerie", "surname", "name", "patronymic", "address"}).
		AddRow(expectedUser.ID, expectedUser.PassportNumber, expectedUser.PassportSerie, expectedUser.Surname, expectedUser.Name, expectedUser.Patronymic, expectedUser.Address)

	mock.ExpectQuery(`SELECT \* FROM users WHERE id = \$1$`).
		WithArgs(1).
		WillReturnRows(rows)

	user, err := repo.GetUserByID(1)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.Surname, user.Surname)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewApiRepository(sqlxDB)

	user := &models.User{
		PassportNumber: "1234",
		PassportSerie:  "567890",
		Surname:        "Иванов",
		Name:           "Иван",
		Patronymic:     "Иванович",
		Address:        "г. Москва",
	}

	mock.ExpectExec(`INSERT INTO users (.+) VALUES (.+)$`).
		WithArgs(user.PassportNumber, user.PassportSerie, user.Surname, user.Name, user.Patronymic, user.Address).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.CreateUser(user)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewApiRepository(sqlxDB)

	user := &models.User{
		ID:             1,
		PassportNumber: "1234",
		PassportSerie:  "567890",
		Surname:        "Иванов",
		Name:           "Иван",
		Patronymic:     "Иванович",
		Address:        "г. Москва",
	}

	mock.ExpectExec(`UPDATE users SET passport_number=\$1, passport_serie=\$2, surname=\$3, name=\$4, patronymic=\$5, address=\$6 WHERE id=\$7`).
		WithArgs(user.PassportNumber, user.PassportSerie, user.Surname, user.Name, user.Patronymic, user.Address, user.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.UpdateUser(user)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDeleteUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewApiRepository(sqlxDB)

	mock.ExpectExec(`DELETE FROM users WHERE id = \$1`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.DeleteUser(1)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
