package service

import (
	"TestEffectiveMobile/cmd/internal/model"
	"TestEffectiveMobile/cmd/internal/repository"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

// Интерфейс сервиса для работы с людьми
type PersonService interface {
	AddPerson(person model.Person) error
	GetPersons(page, limit int, name, gender, nationality string) ([]model.Person, error)
	UpdatePerson(id int, person model.Person) error
	DeletePerson(id int) error
}

// Реализация сервиса для работы с людьми
type PersonServiceImpl struct {
	repo repository.PersonRepository
}

// Конструктор для создания нового сервиса
func NewPersonService(repo repository.PersonRepository) *PersonServiceImpl {
	return &PersonServiceImpl{repo: repo}
}

// Добавление нового человека с обогащением данных из внешних API
func (s *PersonServiceImpl) AddPerson(person model.Person) error {
	slog.Info("Получение данных для имени: ", person.Name)

	// Обогащение данными из внешних API
	age, err := getAge(person.Name)

	if err != nil {
		return err
	}
	fmt.Println("Полученное имя", age)

	gender, err := getGender(person.Name)
	if err != nil {
		return err
	}

	nationality, err := getNationality(person.Name)
	if err != nil {
		return err
	}

	// Сохранение в базе данных
	err = s.repo.SavePerson(person.Name, person.Surname, person.Patronymic, age, gender, nationality)
	if err != nil {
		slog.Error("Ошибка сохранения человека: ", err)
		return err
	}
	slog.Info("Человек успешно добавлен в базу данных.")
	return nil
}

func (s *PersonServiceImpl) GetPersons(page, limit int, name, gender, nationality string) ([]model.Person, error) {
	slog.Info("Получение людей с фильтрами: ", name, gender, nationality)

	// Валидация параметров пагинации
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	// Получаем людей из репозитория с фильтрами
	return s.repo.GetAllPersons(page, limit, name, gender, nationality)
}

// Обновление данных о человеке
func (s *PersonServiceImpl) UpdatePerson(id int, person model.Person) error {
	return s.repo.UpdatePerson(id, person.Name, person.Surname, person.Patronymic, person.Age, person.Gender, person.Nationality)
}

// Удаление человека по ID
func (s *PersonServiceImpl) DeletePerson(id int) error {
	return s.repo.DeletePerson(id)
}

// Вспомогательные функции для получения данных из внешних API

func getAge(name string) (int, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.agify.io/?name=%s", name))
	if err != nil {
		slog.Error("Ошибка получения возраста: ", err)
		return 0, err
	}
	defer resp.Body.Close()

	var person model.Person
	if err := json.NewDecoder(resp.Body).Decode(&person); err != nil {
		slog.Error("Ошибка декодирования ответа по возрасту: ", err)
		return 0, err
	}

	return person.Age, nil
}

func getGender(name string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.genderize.io/?name=%s", name))
	if err != nil {
		slog.Error("Ошибка получения пола: ", err)
		return "", err
	}
	defer resp.Body.Close()

	var person model.Person
	if err := json.NewDecoder(resp.Body).Decode(&person); err != nil {
		slog.Error("Ошибка декодирования ответа по полу: ", err)
		return "", err
	}

	return person.Gender, nil
}

func getNationality(name string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.nationalize.io/?name=%s", name))
	if err != nil {
		slog.Error("Ошибка получения национальности: ", "error", err)
		return "", err
	}
	defer resp.Body.Close()

	var data model.NationalizeResponse
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		slog.Error("Ошибка декодирования ответа по национальности",
			slog.String("error", err.Error()),
		)
		return "", err
	}

	if len(data.Country) == 0 {
		slog.Info("Не удалось определить национальность для имени", "name", name)
		return "", nil
	}

	nationality := data.Country[0].CountryID
	slog.Info("Определена национальность", "name", name, "nationality", nationality)
	return nationality, nil
}
