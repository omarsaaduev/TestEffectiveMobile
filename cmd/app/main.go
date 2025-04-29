package main

import (
	"TestEffectiveMobile/cmd/internal/config"
	"TestEffectiveMobile/cmd/internal/handler"
	"TestEffectiveMobile/cmd/internal/logger"
	"TestEffectiveMobile/cmd/internal/repository"
	"TestEffectiveMobile/cmd/internal/service"
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	const op = "cmd.app.main"

	// Загрузка конфигурации
	cfg := config.NewConfig()

	// Настройка логирования
	if err := logger.SetupGlobalLogger(&cfg.Log); err != nil {
		fmt.Printf("Не удалось настроить logger: %v\n", err)
		os.Exit(1)
	}

	slog.Info("Запуск приложения", "environment", cfg.Env)

	// Установка соединения с базой данных Psql
	db, err := sql.Open("postgres", cfg.DB.GetConnectionString())
	if err != nil {
		slog.Error(op, "Ошибка подключения к базе данных:", err)
		os.Exit(1)
	}
	defer db.Close()

	// Проверка соединения с базой данных
	if err = db.Ping(); err != nil {
		slog.Error(op, "Ошибка при проверке соединения с базой данных:", err)
		os.Exit(1)
	}

	// Создаем репозиторий
	repo := repository.NewPersonRepositoryPgSQL(db)

	// Создаем сервис
	ps := service.NewPersonService(repo)

	// Настройка маршрутов с использованием Gorilla Mux
	r := mux.NewRouter()

	// Регистрация маршрутов
	handler.SetupRoutes(r, handler.NewPersonHandler(ps))

	// Старт сервера
	slog.Warn(fmt.Sprintf("Сервер запущен и прослушивает порт %s\n", cfg.Server.Port))
	if err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.Server.Port), r); err != nil {
		log.Fatalf("ошибка запуска сервера: %v", err)
	}
}
