package model

type Show struct {
	Id int 				`json:"id"`
	Name string 		`json:"name"`
	Overview string 	`json:"overview"`
	AirDate string 		`json:"first_air_date"`
	Score float64 		`json:"vote_average"`
	Status string		`json:"status"`
	EpisodeCount int	`json:"number_of_episodes"`
	SeasonCount int		`json:"number_of_seasons"`
	Seasons []Season	`json:"seasons"`
}
