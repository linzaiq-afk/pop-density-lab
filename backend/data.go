package main

import (
	"math"
	"sort"
)

// YearlyStat описывает состояние объекта за конкретный год.
type YearlyStat struct {
	Year       int     `json:"year"`
	Population int64   `json:"population"` // человек
	AreaKm2    float64 `json:"area_km2"`   // км²
	Density    float64 `json:"density"`    // человек на км²
}

// Entity — страна или город с историческим рядом наблюдений.
type Entity struct {
	ID      string       `json:"id"`
	Name    string       `json:"name"`
	Type    string       `json:"type"` // "country" | "city"
	Country string       `json:"country"`
	History []YearlyStat `json:"history"`
}

// calcDensity считает плотность населения ρ = N / S и округляет её до сотых.
func calcDensity(population int64, areaKm2 float64) float64 {
	if population < 0 || areaKm2 <= 0 {
		return 0
	}
	return math.Round(float64(population)/areaKm2*100) / 100
}

// getDataset возвращает демонстрационный датасет (округлённые оценки
// численности населения и площади, 1960–2026). Данные предназначены для
// учебной лабораторной работы, а не для официальной статистики.
func getDataset() []Entity {
	raw := []Entity{
		{
			ID: "ru", Name: "Россия", Type: "country", Country: "Россия",
			History: []YearlyStat{
				{Year: 1960, Population: 119897000, AreaKm2: 17098246},
				{Year: 1980, Population: 138126000, AreaKm2: 17098246},
				{Year: 2000, Population: 146597000, AreaKm2: 17098246},
				{Year: 2020, Population: 146171000, AreaKm2: 17125191},
				{Year: 2026, Population: 145800000, AreaKm2: 17125191},
			},
		},
		{
			ID: "cn", Name: "Китай", Type: "country", Country: "Китай",
			History: []YearlyStat{
				{Year: 1960, Population: 667070000, AreaKm2: 9596961},
				{Year: 1980, Population: 981235000, AreaKm2: 9596961},
				{Year: 2000, Population: 1262645000, AreaKm2: 9596961},
				{Year: 2020, Population: 1411100000, AreaKm2: 9596961},
				{Year: 2026, Population: 1405000000, AreaKm2: 9596961},
			},
		},
		{
			ID: "jp", Name: "Япония", Type: "country", Country: "Япония",
			History: []YearlyStat{
				{Year: 1960, Population: 92501000, AreaKm2: 377975},
				{Year: 1980, Population: 116782000, AreaKm2: 377975},
				{Year: 2000, Population: 126843000, AreaKm2: 377975},
				{Year: 2020, Population: 126261000, AreaKm2: 377975},
				{Year: 2026, Population: 123200000, AreaKm2: 377975},
			},
		},
		{
			ID: "us", Name: "США", Type: "country", Country: "США",
			History: []YearlyStat{
				{Year: 1960, Population: 180671000, AreaKm2: 9833517},
				{Year: 1980, Population: 227225000, AreaKm2: 9833517},
				{Year: 2000, Population: 282162000, AreaKm2: 9833517},
				{Year: 2020, Population: 331527000, AreaKm2: 9833517},
				{Year: 2026, Population: 340000000, AreaKm2: 9833517},
			},
		},
		{
			ID: "in", Name: "Индия", Type: "country", Country: "Индия",
			History: []YearlyStat{
				{Year: 1960, Population: 445000000, AreaKm2: 3287263},
				{Year: 1980, Population: 697000000, AreaKm2: 3287263},
				{Year: 2000, Population: 1056000000, AreaKm2: 3287263},
				{Year: 2020, Population: 1396000000, AreaKm2: 3287263},
				{Year: 2026, Population: 1450000000, AreaKm2: 3287263},
			},
		},
		{
			ID: "de", Name: "Германия", Type: "country", Country: "Германия",
			History: []YearlyStat{
				{Year: 1960, Population: 72815000, AreaKm2: 357588},
				{Year: 1980, Population: 78288000, AreaKm2: 357588},
				{Year: 2000, Population: 81500000, AreaKm2: 357588},
				{Year: 2020, Population: 83200000, AreaKm2: 357588},
				{Year: 2026, Population: 83500000, AreaKm2: 357588},
			},
		},
		{
			ID: "ng", Name: "Нигерия", Type: "country", Country: "Нигерия",
			History: []YearlyStat{
				{Year: 1960, Population: 45138000, AreaKm2: 923768},
				{Year: 1980, Population: 73423000, AreaKm2: 923768},
				{Year: 2000, Population: 122877000, AreaKm2: 923768},
				{Year: 2020, Population: 208327000, AreaKm2: 923768},
				{Year: 2026, Population: 232000000, AreaKm2: 923768},
			},
		},
		{
			ID: "br", Name: "Бразилия", Type: "country", Country: "Бразилия",
			History: []YearlyStat{
				{Year: 1960, Population: 72179000, AreaKm2: 8515767},
				{Year: 1980, Population: 121150000, AreaKm2: 8515767},
				{Year: 2000, Population: 174018000, AreaKm2: 8515767},
				{Year: 2020, Population: 212559000, AreaKm2: 8515767},
				{Year: 2026, Population: 218000000, AreaKm2: 8515767},
			},
		},
		{
			ID: "msk", Name: "Москва", Type: "city", Country: "Россия",
			History: []YearlyStat{
				{Year: 1960, Population: 6133000, AreaKm2: 878.7},
				{Year: 1980, Population: 8111000, AreaKm2: 878.7},
				{Year: 2000, Population: 9932000, AreaKm2: 1081.0},
				{Year: 2020, Population: 12678000, AreaKm2: 2561.5}, // после расширения Новой Москвы
				{Year: 2026, Population: 13250000, AreaKm2: 2561.5},
			},
		},
		{
			ID: "spb", Name: "Санкт-Петербург", Type: "city", Country: "Россия",
			History: []YearlyStat{
				{Year: 1960, Population: 3321000, AreaKm2: 1439},
				{Year: 1980, Population: 4588000, AreaKm2: 1439},
				{Year: 2000, Population: 4661000, AreaKm2: 1439},
				{Year: 2020, Population: 5384000, AreaKm2: 1439},
				{Year: 2026, Population: 5600000, AreaKm2: 1439},
			},
		},
		{
			ID: "kzn", Name: "Казань", Type: "city", Country: "Россия",
			History: []YearlyStat{
				{Year: 1960, Population: 667000, AreaKm2: 280.0},
				{Year: 1980, Population: 1002000, AreaKm2: 340.0},
				{Year: 2000, Population: 1101000, AreaKm2: 425.3},
				{Year: 2020, Population: 1257000, AreaKm2: 588.9},
				{Year: 2026, Population: 1320000, AreaKm2: 614.0},
			},
		},
		{
			ID: "tyo", Name: "Токио", Type: "city", Country: "Япония",
			History: []YearlyStat{
				{Year: 1960, Population: 9683000, AreaKm2: 2187.0},
				{Year: 1980, Population: 11618000, AreaKm2: 2187.0},
				{Year: 2000, Population: 12064000, AreaKm2: 2188.0},
				{Year: 2020, Population: 14047000, AreaKm2: 2194.0},
				{Year: 2026, Population: 14210000, AreaKm2: 2194.0},
			},
		},
		{
			ID: "nyc", Name: "Нью-Йорк", Type: "city", Country: "США",
			History: []YearlyStat{
				{Year: 1960, Population: 7782000, AreaKm2: 778.2},
				{Year: 1980, Population: 7072000, AreaKm2: 778.2},
				{Year: 2000, Population: 8008000, AreaKm2: 778.2},
				{Year: 2020, Population: 8804000, AreaKm2: 778.2},
				{Year: 2026, Population: 8900000, AreaKm2: 778.2},
			},
		},
		{
			ID: "del", Name: "Дели", Type: "city", Country: "Индия",
			History: []YearlyStat{
				{Year: 1960, Population: 2644000, AreaKm2: 1484},
				{Year: 1980, Population: 5729000, AreaKm2: 1484},
				{Year: 2000, Population: 13850000, AreaKm2: 1484},
				{Year: 2020, Population: 30093000, AreaKm2: 1484},
				{Year: 2026, Population: 34000000, AreaKm2: 1484},
			},
		},
		{
			ID: "lag", Name: "Лагос", Type: "city", Country: "Нигерия",
			History: []YearlyStat{
				{Year: 1960, Population: 1000000, AreaKm2: 1171},
				{Year: 1980, Population: 2600000, AreaKm2: 1171},
				{Year: 2000, Population: 7300000, AreaKm2: 1171},
				{Year: 2020, Population: 14862000, AreaKm2: 1171},
				{Year: 2026, Population: 17000000, AreaKm2: 1171},
			},
		},
		{
			ID: "ber", Name: "Берлин", Type: "city", Country: "Германия",
			History: []YearlyStat{
				{Year: 1960, Population: 3274000, AreaKm2: 891.7},
				{Year: 1980, Population: 3048000, AreaKm2: 891.7},
				{Year: 2000, Population: 3382000, AreaKm2: 891.7},
				{Year: 2020, Population: 3664000, AreaKm2: 891.7},
				{Year: 2026, Population: 3900000, AreaKm2: 891.7},
			},
		},
		{
			ID: "sao", Name: "Сан-Паулу", Type: "city", Country: "Бразилия",
			History: []YearlyStat{
				{Year: 1960, Population: 3667000, AreaKm2: 1521},
				{Year: 1980, Population: 8493000, AreaKm2: 1521},
				{Year: 2000, Population: 10434000, AreaKm2: 1521},
				{Year: 2020, Population: 12325000, AreaKm2: 1521},
				{Year: 2026, Population: 12900000, AreaKm2: 1521},
			},
		},
	}

	for i := range raw {
		sort.Slice(raw[i].History, func(a, b int) bool {
			return raw[i].History[a].Year < raw[i].History[b].Year
		})
		for j := range raw[i].History {
			stat := &raw[i].History[j]
			stat.Density = calcDensity(stat.Population, stat.AreaKm2)
		}
	}
	return raw
}
