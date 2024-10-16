package main

import (
	"bufio"
	"encoding/json"
	"fmt"
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

	fmt.Println("запись в формате команда - задача - описание задачи: ")

	reader := bufio.NewReader(os.Stdin)

	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Ошибка чтения:", err)
		return
	}

	if strings.Contains(input, "add") {

		addTask(input)

	}

	if strings.Contains(input, "complete") {

		completeTask(input)

	}

	if strings.Contains(input, "delete") {

		deleteTask(input)

	}

	if strings.Contains(input, "show") {

		showTask(input)

	}

	showTask(input)

}

func SearchFile(tasks []Task, name string) *Task {

	for i, task := range tasks {

		if task.ID == name {
			return &tasks[i]
			break
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

func addTask(input string) {

	words := strings.Split(input, " - ")

	i := Task{
		ID:          words[3],
		Name:        words[1],
		Description: words[2],
		Completed:   false,
	}

	var tasks []Task

	dataFile, err := os.ReadFile("TestFile.json")

	err = json.Unmarshal(dataFile, &tasks)

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

}

func completeTask(input string) {

	words := strings.Split(input, " ")

	var tasks []Task

	data, err := os.ReadFile("TestFile.json")
	if err != nil {
		fmt.Println("Ошибка чтения файла:", err)
		return
	}

	err = json.Unmarshal(data, &tasks)

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

	} else {
		fmt.Println("ID не найден")
	}

}

func deleteTask(input string) {

	words := strings.Split(input, " ")

	var tasks []Task

	data, err := os.ReadFile("TestFile.json")
	if err != nil {
		fmt.Println("Ошибка чтения файла:", err)
		return
	}

	err = json.Unmarshal(data, &tasks)

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

	clearJSONfile("TestFile.json", DelTasks)

}

func showTask(input string) {

	words := strings.Split(input, " ")

	var tasks []Task

	data, err := os.ReadFile("TestFile.json")
	if err != nil {
		fmt.Println("Ошибка чтения файла:", err)
		return
	}

	err = json.Unmarshal(data, &tasks)

	if err != nil {
		fmt.Println("Ошибка десериализации JSON:", err)
		return
	}

	for i := 0; i < len(tasks); i++ {
		if tasks[i].ID == words[1] {
			fmt.Println(tasks[i].Name, tasks[i].Description, tasks[i].Completed)
		}
	}

}
