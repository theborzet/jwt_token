package models

type User struct {
	ID             int    `db:"id" json:"id"`
	PassportNumber string `db:"passportNumber" json:"passportNumber"`
	PassportSerie  string `db:"passportSerie" json:"passportSerie"`
	Surname        string `db:"surname" json:"surname"`
	Name           string `db:"name" json:"name"`
	Patronymic     string `db:"patronymic" json:"patronymic"`
	Address        string `db:"address" json:"address"`
}

type Task struct {
	ID        int    `db:"id" json:"id"`
	UserID    int    `db:"userId" json:"userId"`
	TaskName  string `db:"taskName" json:"taskName"`
	StartTime string `db:"startTime" json:"startTime"`
	EndTime   string `db:"endTime" json:"endTime"`
}

type TaskTimeSpent struct {
	TaskName  string  `db:"taskName" json:"taskName"`
	TimeSpent float64 `db:"timeSpent" json:"timeSpent"`
}
