package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Task struct {
	ID          string `json:"id"`          // Уникальный идентификатор задачи
	Name        string `json:"name"`        // Название задачи
	Description string `json:"description"` // Описание задачи
	Completed   bool   `json:"completed"`   // Статус задачи (выполнена или нет)
}

var Tasks []Task

func main() {

	fmt.Println("Open file")

	file, err := os.OpenFile("TestFile.json", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)

	if err != nil {
		panic(err)
	}

	fmt.Println("запись в формате команда - задача - описание задачи: ")

	reader := bufio.NewReader(os.Stdin)

	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Ошибка чтения:", err)
		return
	}

	if strings.Contains(input, "add") {

		words := strings.Split(input, " - ")

		i := Task{ID: words[3], Name: words[1], Description: words[2], Completed: false}

		data, err := json.Marshal(i)

		if err != nil {
			fmt.Println("Ошибка:", err)
			return
		}

		_, err = file.Write(data)

		fmt.Println("Save file!")

		if err != nil {
			fmt.Println("Ошибка при записи в файл:", err)
			return
		}

	}

	if strings.Contains(input, "complete") {

		words := strings.Split(input, " ")

		//jsonData := `{"id":1}`
		var tasks []Task

		data, err := ioutil.ReadFile("TestFile.json")
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

		otvet[0].Completed = true

		//err = clearJSONfile("TestFile.json")

		data, err = json.MarshalIndent(otvet, "", "")

		_, err = file.Write(data)

		fmt.Println(otvet)

	}

	fmt.Println("Close file!")

}

func SearchFile(tasks []Task, name string) []Task {

	var founTask []Task

	for i, task := range tasks {

		if task.ID == name {
			founTask = append(founTask, tasks[i])
		}

	}

	return founTask

}

func clearJSONfile(filename string) error {

	file, err := os.Create(filename)

	if err != nil {
		return err
	}

	defer file.Close()

	emptyTasks := []Task{}

	encoder := json.NewEncoder(file)
	err = encoder.Encode(emptyTasks)

	if err != nil {
		return err
	}

	return nil

}
