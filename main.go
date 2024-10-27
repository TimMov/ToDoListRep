package main

import (
	"bufio"
	"database/sql"
	_ "database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"os"
	"strings"
)

type Task struct {
	ID          string `json:"id"`          // Уникальный идентификатор задачи
	Name        string `json:"name"`        // Название задачи
	Description string `json:"description"` // Описание задачи
	Completed   bool   `json:"completed"`   // Статус задачи (выполнена или нет)
}

func main() {

	fmt.Println("Выберите команду: ")
	fmt.Println("add - задача - описание задачи - ID")
	fmt.Println("addProfile Имя")
	fmt.Println("complete ID")
	fmt.Println("delete ID")
	fmt.Println("show ID")
	fmt.Println("exit")

	reader := bufio.NewReader(os.Stdin)

	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Ошибка чтения:", err)
		return
	}

	input = strings.TrimSpace(input)

	command := strings.Split(input, " ")[0]

	dataFile, err := os.ReadFile("TestFile.json")

	if err != nil {
		fmt.Println("Ошибка: ", err)
		return
	}

	var tasks []Task
	var db *sql.DB

	if command != "exit" {
		db, err = openDataBase()

		if err != nil {
			fmt.Println(err)
			return
		}
	}

	switch command {
	case "add":
		addTask(input, dataFile, tasks, db)
	case "addProfile":
		addProfile(input, db)
	case "complete":
		completeTask(input, dataFile, tasks, db)
	case "delete":
		deleteTask(input, dataFile, tasks, db)
	case "show":
		showTask(input, dataFile, tasks, db)
	case "exit":
		return
	default:
		fmt.Println("Неизвестная команда. Доступные команды: add, complete, delete, show, addProfile, exit")
	}

}

func SearchFile(tasks []Task, name string) *Task {

	for i, task := range tasks {
		if task.ID == name {
			return &tasks[i]
		}
	}

	return nil

}

func clearJSONfile(filename string, Massive []Task) error {

	file, err := os.Create(filename)

	if err != nil {
		return err
	}

	defer file.Close()

	emptyTasks := Massive

	encoder := json.NewEncoder(file)
	err = encoder.Encode(emptyTasks)

	if err != nil {
		return err
	}

	return nil

}

func addTask(input string, dataFile []byte, tasks []Task, db *sql.DB) {

	words := strings.Split(input, " - ")

	if len(words) != 4 {
		fmt.Println("Неккоректный ввод команды")
		fmt.Println("Попробуйте ввести в формате: add - задача - описание задачи - ID")

		return
	}

	i := Task{
		ID:          words[3],
		Name:        words[1],
		Description: words[2],
		Completed:   false,
	}

	err := json.Unmarshal(dataFile, &tasks)

	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	CheckTasks := SearchFile(tasks, words[3])

	if CheckTasks != nil {
		fmt.Println("ID же есть:", CheckTasks.ID)
		return
	}

	var founTasks []Task

	founTasks = append(founTasks, i)

	for Num := 0; Num < len(tasks); Num++ {
		founTasks = append(founTasks, tasks[Num])
	}

	err = clearJSONfile("TestFile.json", founTasks)

	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	fmt.Println("Save file!")

	query := `insert into tasks (userid, task, decription, complete) VALUES ($1, $2, $3, false)`

	_, err = db.Exec(query, 1, words[1], words[2])

	if err != nil {
		panic(err)
	}

}

func completeTask(input string, dataFile []byte, tasks []Task, db *sql.DB) {

	words := strings.Split(input, " ")

	if len(words) != 2 {

		fmt.Println("Неккоректный ввод команды")
		fmt.Println("Попробуйте ввести в формате: completeTask ID")

		return

	}

	err := json.Unmarshal(dataFile, &tasks)

	if err != nil {
		fmt.Println("Ошибка десериализации JSON:", err)
		return
	}

	otvet := SearchFile(tasks, words[1])

	if otvet != nil {
		for Num := 0; Num < len(tasks); Num++ {
			if tasks[Num].ID == otvet.ID {
				tasks[Num].Completed = true
				break
			}
		}

		err = clearJSONfile("TestFile.json", tasks)

		fmt.Println("Задача ", otvet.ID, " выполнена")

	} else {
		fmt.Println("ID не найден")
		return
	}

	query := `update tasks set complete = true where id = $1 and complete = false`

	_, err = db.Exec(query, words[1])

	if err != nil {
		panic(err)
	}

}

func deleteTask(input string, dataFile []byte, tasks []Task, db *sql.DB) {

	words := strings.Split(input, " ")

	if len(words) != 2 {

		fmt.Println("Неккоректный ввод команды")
		fmt.Println("Попробуйте ввести в формате: deleteTask ID")

		return

	}

	err := json.Unmarshal(dataFile, &tasks)

	if err != nil {
		fmt.Println("Ошибка десериализации JSON:", err)
		return
	}

	var DelTasks []Task

	for i := 0; i < len(tasks); i++ {
		if tasks[i].ID != words[1] {
			DelTasks = append(DelTasks, tasks[i])
		}
	}

	err = clearJSONfile("TestFile.json", DelTasks)

	if err != nil {
		fmt.Println("Ошибка: ", err)
		return
	}

	query := `delete from tasks where ID = $1`

	_, err = db.Exec(query, words[1])

	if err != nil {
		panic(err)
	}

	fmt.Println("Удалено")

}

func showTask(input string, dataFile []byte, tasks []Task, db *sql.DB) {

	words := strings.Split(input, " ")

	if len(words) != 2 {

		fmt.Println("Неккоректный ввод команды")
		fmt.Println("Попробуйте ввести в формате: showTask ID")

		return

	}

	err := json.Unmarshal(dataFile, &tasks)

	if err != nil {
		fmt.Println("Ошибка десериализации JSON:", err)
		return
	}

	for i := 0; i < len(tasks); i++ {
		if tasks[i].ID == words[1] {
			fmt.Println(tasks[i].Name, tasks[i].Description, tasks[i].Completed)
		}
	}

	query := `select ID, task, decription, complete from tasks where ID = $1`

	rows, err := db.Query(query, words[1])

	if rows.Next() {

		err := rows.Scan(&tasks[0].ID, &tasks[0].Name, &tasks[0].Description, &tasks[0].Completed)

		if err != nil {
			panic(err)
		}

		fmt.Println("ID: ", tasks[0].ID)
		fmt.Println("Наименование: ", tasks[0].Name)
		fmt.Println("Описание: ", tasks[0].Description)
		fmt.Println("Выполнено: ", tasks[0].Completed)

	}

	if err != nil {
		panic(err)
	}

}

func addProfile(input string, db *sql.DB) {

	words := strings.Split(input, " ")

	if len(words) != 2 {

		fmt.Println("Неккоректный ввод команды")
		fmt.Println("Попробуйте ввести в формате: addProfile Имя пользователя")

		return

	}

	query := `INSERT INTO users (name) VALUES ($1)`

	_, err := db.Exec(query, words[1])

	if err != nil {
		panic(err)
	}

}

func openDataBase() (*sql.DB, error) {

	conn := "host=localhost port=5432 user=admin password=root dbname=todolist sslmode=disable"

	db, err := sql.Open("postgres", conn)

	return db, err

}
