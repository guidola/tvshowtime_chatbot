package model

type Episode struct {
	Id int 				`json:"episode_number"`
	Name string 		`json:"name"`
	Overview string 	`json:"overview"`
	AirDate string 		`json:"air_date"`
	Score float64 		`json:"vote_average"`
	Popularity float64	`json:"vote_count"`
}
