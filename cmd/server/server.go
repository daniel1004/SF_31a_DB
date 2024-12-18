package main

import (
	"GoNews/pkg/api"
	"GoNews/pkg/storage"
	"GoNews/pkg/storage/memdb"
	mongo "GoNews/pkg/storage/mondoDB"
	"GoNews/pkg/storage/postgres"
	"log"
	"net/http"
)

const (
	databaseName   = "test"
	collectionName = "posts"
)

// Сервер GoNews.
type server struct {
	db  storage.Interface
	api *api.API
}

func main() {
	// Создаём объект сервера.
	var srv server

	// Создаём объекты баз данных.
	//
	// БД в памяти.
	db := memdb.New()

	// Реляционная БД PostgreSQL.
	db2, err := postgres.New("postgresql://postgres:admin@localhost:5432/TEST_32")
	if err != nil {
		log.Fatal(err)
	}

	//// Документная БД MongoDB.
	db3, err := mongo.New("mongodb://172.17.0.2:27017/", databaseName, collectionName)
	if err != nil {
		log.Fatal(err)
	}
	_, _, _ = db, db2, db3

	// Инициализируем хранилище сервера конкретной БД.
	srv.db = db3

	// Создаём объект API и регистрируем обработчики.
	srv.api = api.New(srv.db)

	// Запускаем веб-сервер на порту 8080 на всех интерфейсах.
	// Предаём серверу маршрутизатор запросов,
	// поэтому сервер будет все запросы отправлять на маршрутизатор.
	// Маршрутизатор будет выбирать нужный обработчик.
	http.ListenAndServe(":8080", srv.api.Router())
}
