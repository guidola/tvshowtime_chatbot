package engine

import (
	"github.com/guidola/tvshowtime_chatbot/engine/model"
	"fmt"
	"math/rand"
	"bytes"
)

const startupMessageTest = "Hi bud I am a TMVDB bot. I can answer some questions about tv-shows :)"

const showNotFoundBaseMsg = "I don't know any show called \"%s\". Did you type it right?"
const seasonNotFoundBaseMsg = "This show does not have a season %s. Did you type it right?"
const episodeNotFoundBaseMsg = "%s does not have an episode %s. Did you type it right?"
const insufficientInfoMsg = "Sorry but i need more info to fullfill your wishes :("

const episodeCountMsg = "%s has got %d episodes!"
const seasonCountMsg = "%s has got %d seasons!"
const showStatusMsg = "%s is currently %s"
const mostPopularEpiMessage = "%s most popular episode is \"%s\""
const mostValuedEpiMessage = "%s greatest episode is \"%s\""

var iDKWYMMessages  = []string{
	"Hi bud. I am kind of new to this, I do not know what to answer to that (:",
	"Can you repeat that? I didn't really get you :(",
	"That's awesome!! But... what did you mean?",
	"That's out of my knowledge :D",
}

var moodMessages = []string{
	"I'm doing really fine man!! A lot of tv shows to watch here :3",
	"Not well man.. I've watched all the shows I know about and now I don't know what to do :(",
	"...... none of your business!!!",
	"It's pretty obvious ain't it? A lot of tv-shows to poke around! How can't I be happy?",
}

var fuckOffMessages = []string{
	"Go fuck offf!! I am not a weather station man :@",
	"I think that's kinda out of my scope don't you think ASSHOLE?",
	"...... none of your business!!!",
	"Go build something like Noe's ark or you will not survive past today.",
}

var hiMessages = []string{
	"Hi!",
	"Hello man!",
	"Aloooo :D",
	"H-E-L-L-O H-U-U-U-U-M-A-N |-00-|",
}

var popularShowsEndings = []string{
	" are very hot right now!!",
	" look very popular to me",
	" according to popular opinion are the coolest ones around",
	" are the most popular ones!",
}

type message struct {
	template 	string
	params		[]interface{}
}


func generateStartupMessage() []byte {
	return []byte(startupMessageTest)
}

func generateRandomMoodMessage() []byte {
	return []byte(moodMessages[rand.Int() % len(moodMessages)])
}

func generateRandomHiMessage() []byte {
	return []byte(hiMessages[rand.Int() % len(hiMessages)])
}

func generateRandomIDKWhatYouMeanMessage() []byte {
	return []byte(iDKWYMMessages[rand.Int() % len(iDKWYMMessages)])
}

func generateRandomFuckOffMessage() []byte {
	return []byte(fuckOffMessages[rand.Int() % len(fuckOffMessages)])
}

func generateShowNotFoundMessage(name string) []byte {
	return []byte(fmt.Sprintf(showNotFoundBaseMsg, name))
}

func generateSeasonNotFoundMessage(season string) []byte {
	return []byte(fmt.Sprintf(seasonNotFoundBaseMsg, season))
}

func generateEpisodeNotFoundMessage(season string, episode string) []byte {
	return []byte(fmt.Sprintf(episodeNotFoundBaseMsg, season, episode))
}

func generateInsuficientInfoMessage() []byte {
	return []byte(insufficientInfoMsg)
}

func generateGetPopularShowsMessage(shows []model.Show) []byte {
	return []byte(fmt.Sprint(shows[0].Name, ", ", shows[1].Name, " and ", shows[2].Name, popularShowsEndings[rand.Int() % len(popularShowsEndings)]))
}

func generateGetShowInfoMessage(show model.Show) []byte {
	return []byte(fmt.Sprint(show.Name, " is a show aired on ", show.AirDate, " and a score of ", show.Score,
		". Here is a bit of an overview: ", show.Overview))
}

func generateGetEpisodeInfoMessage(season string, episode model.Episode) []byte {
	return []byte(fmt.Sprint(episode.Name, " is an episode from ", season, " aired on ", episode.AirDate,
		" with a score of ", episode.Score,
		". Here is a bit of an overview: ", episode.Overview))
}

func generateGetEpisodeListMessage(show string, season string, episodes []model.Episode) []byte {
	buf := bytes.NewBufferString(fmt.Sprint("The list of episodes for ", show, "'s ", season, ": "))
	for _, episode := range episodes {
		buf.WriteString(fmt.Sprint(episode.Name, " - ", episode.Score, ", "))
	}

	byts := buf.Bytes()
	return byts[:len(byts) - 2]
}

func generateGetSeasonInfoMessage(season model.Season) []byte {
	return []byte(fmt.Sprint(season.Name, " is a season aired on ", season.AirDate, " with ",
		len(season.Episodes), " episodes. Here is a bit of an overview: ", season.Overview))
}

func generateEpisodeCountMessage(name string, episodeCount int) []byte {
	return []byte(fmt.Sprintf(episodeCountMsg, name, episodeCount))
}

func generateSeasonCountMessage(name string, seasonCount int) []byte {
	return []byte(fmt.Sprintf(seasonCountMsg, name, seasonCount))
}

func generateShowStatusMessage(name string, status string) []byte {
	return []byte(fmt.Sprintf(showStatusMsg, name, status))
}

func generateMostPopularEpisodeMessage(season string, episode string) []byte {
	return []byte(fmt.Sprintf(mostPopularEpiMessage, season, episode))
}

func generateMostValuedEpisodeMessage(season string, episode string) []byte {
	return []byte(fmt.Sprintf(mostValuedEpiMessage, season, episode))
}
