package main

import (
	"encoding/csv"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
)

type Product struct {
	Name         string
	Image        string
	Cost         int
	Diametr      string
	Season       string
	Seasontr     string
	Width        string
	Profile      string
	Manufacturer string
}

type Data struct {
	Products []Product
	MinPrice uint16
	MaxPrice uint16
	Widths   []int
	Profiles []int
	Diameters []int
}

func check(err error) {
	if err != nil {
		log.Fatal(fmt.Sprintf("ОШИБКА: %v \n", err))
	}
}

func mainPageHandler(w http.ResponseWriter, r *http.Request) {
	data := AllValues()

	tmpl, err := template.ParseFiles("templates/index.html")
	check(err)
	tmpl.Execute(w, data)
}

func ParseCSV(filename string) ([]Product, error) {
	file, err := os.Open(filename)
	check(err)

	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'
	rows, err := reader.ReadAll()
	check(err)

	var products []Product
	for _, row := range rows[1:] {
		cost, _ := strconv.Atoi(row[6])

		var season string
		switch row[4] {
		case "зима":
			season = "winter"
		case "лето":
			season = "summer"
		}

		imgpath := fmt.Sprintf("static/img/tires/%s", row[7])

		products = append(products, Product{
			Name:         row[0],  // название
			Image:        imgpath, // путь к изображению
			Cost:         cost,    // цена
			Diametr:      row[3],  // диаметр (4-й столбец, индекс 3)
			Season:       season,  // сезон (англ.)
			Seasontr:     row[4],  // сезон (рус.)
			Width:        row[1],  // ширина
			Profile:      row[2],  // профиль
			Manufacturer: row[5],  // производитель
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

	// Если нет данных, кроме заголовка
	if len(rows) < 2 {
		return []int{0, 0}, nil
	}

	firstCost, err := strconv.Atoi(rows[1][6])
	if err != nil {
		return []int{0, 0}, err
	}

	min, max := firstCost, firstCost

	// Перебираем строки, начиная со второй (индекс 2)
	for _, row := range rows[2:] {
		cost, err := strconv.Atoi(row[6])
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

func UniqueNumericColumn(filename string, columnIndex int) ([]int, error) {
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

	// Храним уникальные числа
	uniqueValues := make(map[int]struct{})

	// Пропускаем заголовок и обрабатываем строки
	for _, row := range rows[1:] {
		if len(row) <= columnIndex {
			continue // Пропускаем строки без нужного столбца
		}

		valueStr := row[columnIndex]
		value, err := strconv.Atoi(valueStr)
		if err != nil {
			continue // Пропускаем некорректные числа
		}

		uniqueValues[value] = struct{}{}
	}

	// Преобразуем map в слайс
	result := make([]int, 0, len(uniqueValues))
	for val := range uniqueValues {
		result = append(result, val)
	}

	// Сортируем
	sort.Ints(result)

	return result, nil
}
func AllValues() Data {
	result, err := MinMaxPrice("data.csv")
	check(err)

	products, err := ParseCSV("data.csv")
	check(err)
	widths, err := UniqueNumericColumn("data.csv", 1)
	profiles, err := UniqueNumericColumn("data.csv", 2)
	diameters, err := UniqueNumericColumn("data.csv", 3)
	data := Data{
		Products: products,
		MinPrice: uint16(result[0]),
		MaxPrice: uint16(result[1]),
		Widths:   widths,
		Profiles: profiles,
		Diameters: diameters,
	}
	return data
}
func main() {
	// Обработчик статических файлов
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Обработка главной страницы
	http.HandleFunc("/", mainPageHandler)
	// Запуск сервера
	port := "127.0.0.1:8080"
	fmt.Printf("Сервер запущен %s\n", port)
	http.ListenAndServe(port, nil)
}
