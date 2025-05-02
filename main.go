package main

import (
	"encoding/csv"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
)

type Product struct {
	Name     string
	Image    string
	Cost     int
	Diametr  string
	Season   string
	Seasontr string
}


func check(err error) {
	if err != nil {
		fmt.Printf("ОШИБКА: %v \n", err)
		os.Exit(1)
	}
}


func mainPageHandler(w http.ResponseWriter, r *http.Request) {
	result, err := MinMaxPrice("data.csv")
	check(err)

	products, err := ParseCSV("data.csv")
	check(err)

	data := struct {
		Products []Product
		MinPrice uint16
		MaxPrice uint16
	}{
		Products: products,
		MinPrice: uint16(result[0]),
		MaxPrice: uint16(result[1]),
	}
	tmpl, err := template.ParseFiles("templates/index.html")
	check(err)
	tmpl.Execute(w, data)
}

func ParseCSV(filename string) ([]Product, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var products []Product
	for _, row := range rows[1:] {
		cost, _ := strconv.Atoi(row[2])
		var season string
		if row[4] == "зима" {
			season = "winter"
		}
		if row[4] == "лето" {
			season = "summer"
		}
		imgpath := fmt.Sprintf("static/img/tires/%s", row[1])
		products = append(products, Product{
			Name:     row[0],
			Image:    imgpath,
			Cost:     cost,
			Diametr:  row[3],
			Season:   season,
			Seasontr: row[4],
		})
	}
	return products, nil
}
func MinMaxPrice(filename string) ([]int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return []int{0, 0}, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'
	rows, err := reader.ReadAll()
	if err != nil {
		return []int{0, 0}, err
	}

	// Инициализируем min и max с первым значением (пропуская заголовок)
	if len(rows) < 2 {
		return []int{0, 0}, nil // Если нет данных, кроме заголовка
	}

	firstCost, err := strconv.Atoi(rows[1][2]) // Первая строка с данными (индекс 1)
	if err != nil {
		return []int{0, 0}, err
	}

	min, max := firstCost, firstCost

	// Перебираем строки, начиная со второй (индекс 2)
	for _, row := range rows[2:] {
		cost, err := strconv.Atoi(row[2]) // cost находится в индексе 2
		if err != nil {
			return []int{0, 0}, err
		}

		if cost < min {
			min = cost
		}
		if cost > max {
			max = cost
		}
	}

	return []int{min, max}, nil
}

func main() {
	// Обработчик статических файлов
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Обработка главной страницы
	http.HandleFunc("/", mainPageHandler)
	// Запуск сервера
	fmt.Println("Сервер запущен")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Сервер остановлен")
	}
}
