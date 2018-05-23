package model

type Season struct {
	Id int 				`json:"season_number"`
	Name string 		`json:"name"`
	Overview string 	`json:"overview"`
	AirDate string 		`json:"air_date"`
	Episodes []Episode 	`json:"episodes"`
	Score float64		`json:"vote_average"`
}
