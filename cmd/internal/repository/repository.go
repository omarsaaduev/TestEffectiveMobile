package repository

import (
	"TestEffectiveMobile/cmd/internal/model"
	"database/sql"
	"fmt"
	"log/slog"
)

type PersonRepository interface {
	SavePerson(name, surname, patronymic string, age int, gender, nationality string) error
	DeletePerson(id int) error
	UpdatePerson(id int, name, surname, patronymic string, age int, gender, nationality string) error
	GetPerson(id int) (*model.Person, error)
	GetAllPersons(page, limit int, name, gender, nationality string) ([]model.Person, error)
}

type PersonRepositoryPgSQL struct {
	db *sql.DB
}

func NewPersonRepositoryPgSQL(db *sql.DB) *PersonRepositoryPgSQL {
	return &PersonRepositoryPgSQL{db: db}
}

func (r *PersonRepositoryPgSQL) SavePerson(name, surname, patronymic string, age int, gender, nationality string) error {
	_, err := r.db.Exec("INSERT INTO persons (name, surname, patronymic, age, gender, nationality) VALUES ($1, $2, $3, $4, $5, $6)", name, surname, patronymic, age, gender, nationality)
	return err
}

func (r *PersonRepositoryPgSQL) DeletePerson(id int) error {
	_, err := r.db.Exec("DELETE FROM persons WHERE id = $1", id)
	return err
}

func (r *PersonRepositoryPgSQL) UpdatePerson(id int, name, surname, patronymic string, age int, gender, nationality string) error {
	_, err := r.db.Exec("UPDATE persons SET name=$1, surname=$2, patronymic=$3, age=$4, gender=$5, nationality=$6 WHERE id=$7",
		name, surname, patronymic, age, gender, nationality, id)
	return err
}

func (r *PersonRepositoryPgSQL) GetPerson(id int) (*model.Person, error) {
	row := r.db.QueryRow("SELECT id, name, surname, patronymic, age, gender, nationality FROM persons WHERE id=$1", id)
	var p model.Person
	err := row.Scan(&p.ID, &p.Name, &p.Surname, &p.Patronymic, &p.Age, &p.Gender, &p.Nationality)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PersonRepositoryPgSQL) GetAllPersons(page, limit int, name, gender, nationality string) ([]model.Person, error) {
	var people []model.Person

	// Строим SQL запрос с фильтрами
	query := "SELECT id, name, surname, patronymic, age, gender, nationality FROM persons WHERE 1=1"

	// Добавляем фильтры
	if name != "" {
		query += fmt.Sprintf(" AND name ILIKE '%%%s%%'", name)
	}
	if gender != "" {
		query += fmt.Sprintf(" AND gender = '%s'", gender)
	}
	if nationality != "" {
		query += fmt.Sprintf(" AND nationality = '%s'", nationality)
	}

	// Пагинация
	query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, (page-1)*limit)

	rows, err := r.db.Query(query)
	if err != nil {
		slog.Error("Ошибка выполнения запроса: ", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var person model.Person
		if err := rows.Scan(&person.ID, &person.Name, &person.Surname, &person.Patronymic, &person.Age, &person.Gender, &person.Nationality); err != nil {
			slog.Error("Ошибка при сканировании строки: ", err)
			return nil, err
		}
		people = append(people, person)
	}

	if err := rows.Err(); err != nil {
		slog.Error("Ошибка при обработке строк: ", err)
		return nil, err
	}

	return people, nil
}
