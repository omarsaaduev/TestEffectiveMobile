package handler

import (
	"TestEffectiveMobile/cmd/internal/model"
	"TestEffectiveMobile/cmd/internal/service"
	_ "TestEffectiveMobile/docs"
	"encoding/json"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
	"strconv"
)

// Интерфейс для обработки запросов с людьми
type PersonHandler interface {
	GetPersons(w http.ResponseWriter, r *http.Request)
	AddPerson(w http.ResponseWriter, r *http.Request)
	UpdatePerson(w http.ResponseWriter, r *http.Request)
	DeletePerson(w http.ResponseWriter, r *http.Request)
}

// Реализация обработчика для людей
type PersonHandlerImpl struct {
	service service.PersonService
}

// Конструктор для создания обработчика
func NewPersonHandler(service service.PersonService) *PersonHandlerImpl {
	return &PersonHandlerImpl{service: service}
}

// Структура стандартного ответа об ошибке
type ErrorResponse struct {
	Detail string `json:"detail"`
}

// Структура стандартного ответа об успехе
type SuccessResponse struct {
	Message string `json:"message"`
}

// Получение всех людей с пагинацией и фильтрами
// @Summary Получить список людей
// @Description Получение всех людей с пагинацией и фильтрами
// @Tags Person
// @Accept json
// @Produce json
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Количество на странице" default(10)
// @Param name query string false "Фильтр по имени"
// @Param gender query string false "Фильтр по полу"
// @Param nationality query string false "Фильтр по национальности"
// @Success 200 {array} model.Person
// @Failure 500 {object} ErrorResponse
// @Router /persons [get]
func (h *PersonHandlerImpl) GetPersons(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page <= 0 {
		page = 1
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		limit = 10
	}
	name := r.URL.Query().Get("name")
	gender := r.URL.Query().Get("gender")
	nationality := r.URL.Query().Get("nationality")

	persons, err := h.service.GetPersons(page, limit, name, gender, nationality)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Не удалось получить список людей")
		return
	}

	h.respondWithJSON(w, http.StatusOK, persons)
}

// Добавление нового человека
// @Summary Добавить нового человека
// @Description Добавляет нового человека в БД с обогащением данными
// @Tags Person
// @Accept json
// @Produce json
// @Param person body model.Person true "Данные нового человека"
// @Success 201 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /person [post]
func (h *PersonHandlerImpl) AddPerson(w http.ResponseWriter, r *http.Request) {
	var person model.Person
	if err := json.NewDecoder(r.Body).Decode(&person); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Некорректный формат данных")
		return
	}

	err := h.service.AddPerson(person)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Не удалось добавить человека")
		return
	}

	h.respondWithJSON(w, http.StatusCreated, SuccessResponse{Message: "Человек успешно добавлен"})
}

// Обновление данных человека
// @Summary Обновить человека
// @Description Обновляет данные человека по ID
// @Tags Person
// @Accept json
// @Produce json
// @Param id path int true "ID человека"
// @Param person body model.Person true "Обновленные данные"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /person/{id} [put]
func (h *PersonHandlerImpl) UpdatePerson(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Некорректный ID")
		return
	}

	var person model.Person
	if err := json.NewDecoder(r.Body).Decode(&person); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Некорректный формат данных")
		return
	}

	err = h.service.UpdatePerson(id, person)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Не удалось обновить данные")
		return
	}

	h.respondWithJSON(w, http.StatusOK, SuccessResponse{Message: "Данные успешно обновлены"})
}

// Удаление человека
// @Summary Удалить человека
// @Description Удаляет человека по ID
// @Tags Person
// @Accept json
// @Produce json
// @Param id path int true "ID человека"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /person/{id} [delete]
func (h *PersonHandlerImpl) DeletePerson(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Некорректный ID")
		return
	}

	err = h.service.DeletePerson(id)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Не удалось удалить человека")
		return
	}

	h.respondWithJSON(w, http.StatusOK, SuccessResponse{Message: "Человек успешно удалён"})
}

// Универсальный метод для ответа с JSON и статусом
func (h *PersonHandlerImpl) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		slog.Error("Ошибка кодирования JSON: ", err)
	}
}

// Универсальный метод для ответа с ошибкой
func (h *PersonHandlerImpl) respondWithError(w http.ResponseWriter, code int, message string) {
	h.respondWithJSON(w, code, ErrorResponse{Detail: message})
}
