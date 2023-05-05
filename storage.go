package main

import (
	"database/sql"
	"fmt"
	"os"
)

var (
	host     = os.Getenv("HOST")
	port     = os.Getenv("PORT")
	user     = os.Getenv("USER")
	password = os.Getenv("PASSWORD")
	dbname   = os.Getenv("DBNAME")
	sslmode  = os.Getenv("SSLMODE")
)

var dbInfo = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode)

func connectDB() (db *sql.DB) {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		fmt.Printf("DATABASE ERROR: %v\n", err)
	}
	return
}

func createTable() error {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		fmt.Printf("DATABASE ERROR: %v\n", err)
	}
	if _, err := db.Exec(`CREATE TABLE passwords(id SERIAL PRIMARY KEY, username TEXT, chat_id INT, service TEXT, password TEXT);`); err != nil {
		fmt.Printf("%v\n", "create Table")
		return err
	}
	return nil
}

// Получает первый пароль от переданного сервиса
func Password(chat_id int, service, username string) (string, error) {
	db := connectDB()
	var service_name, pass string
	fmt.Printf("service: %v\n", service)
	fmt.Printf("chat_id: %v\n", chat_id)
	fmt.Printf("username: %v\n", username)
	row := db.QueryRow(`SELECT service, password FROM passwords WHERE service = $1 AND chat_id = $2 AND username = $3;`, service, chat_id, username)
	err := row.Scan(&service_name, &pass)
	if err != nil {
		fmt.Printf("%v\n", "Passsword")
		return "Not Found", err
	}
	return service_name + ": " + pass, nil
}

// Добавляет одну запись
func AddPassword(username string, chatid int64, service string, password string) error {
	db := connectDB()
	data := `INSERT INTO passwords(username, chat_id, service, password) VALUES($1, $2, $3, $4);`
	if _, err := db.Exec(data, username, chatid, service, password); err != nil {
		fmt.Printf("%v\n", "ADDPasssword")
		return err
	}
	return nil
}

// Удаляет одну запись из таблицы
func DelPassword(username string, chatid int64, service string) error {
	db := connectDB()
	data := `DELETE FROM passwords WHERE service = $1;`
	if _, err := db.Exec(data, service); err != nil {
		fmt.Printf("%v\n", "DELPasssword")
		return err
	}
	return nil
}

func UserData(chat_id int, username string) ([]RowServ, error) {
	db := connectDB()
	var service_name, pass string
	var data []RowServ
	rows, err := db.Query(`SELECT service, password FROM passwords WHERE chat_id = $1 AND username = $2;`, chat_id, username)
	if err != nil {
		fmt.Printf("%v\n", "USERDATA")
		return data, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&service_name, &pass); err != nil {
			return data, err
		}
		data = append(data, RowServ{service_name, pass})
	}

	if err = rows.Err(); err != nil {
		return data, err
	}
	return data, nil
}