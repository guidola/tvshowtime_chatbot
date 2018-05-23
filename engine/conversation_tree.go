package engine

import (
	"strconv"
	"github.com/xrash/smetrics"
	"fmt"
	"strings"
)

const IDLE_STATE 		= 0x00
const SHOW_STATE 		= 0x01
const QUANTITY_STATE 	= 0x02
const SEASON_STATE		= 0x03
const EPISODE_STATE		= 0x04

type intent func(*context) []byte
type directorOptions []string
type parameterPrefixOptions []string

type conversationPath struct {
	directors []directorOptions
	prefixes []parameterPrefixOptions
	executor intent
}

type conversationTree map[int][]conversationPath

func defaultExecutor(ctx *context) []byte {
	return generateRandomIDKWhatYouMeanMessage()
}

func executeSayHiBack(ctx *context) []byte {
	ctx.state = IDLE_STATE
	return generateRandomHiMessage()
}

func executeSendToFuckOff(ctx *context) []byte {
	ctx.state = IDLE_STATE
	return generateRandomFuckOffMessage()
}

func executeTellMyMood(ctx *context) []byte {
	ctx.state = IDLE_STATE
	return generateRandomMoodMessage()
}

func executeGetPopularShows(ctx *context) []byte {
	ctx.state = IDLE_STATE
	return generateGetPopularShowsMessage(ctx.api.GetPopularShows())
}

func executeGetShowInfo(ctx *context) []byte {
	ctx.state = IDLE_STATE
	if ctx.insuficientInfo {
		return generateInsuficientInfoMessage()
	}

	showName := ctx.parameters[0]
	ctx.parameters = []string{}
	show := ctx.api.GetShowInfo(showName)

	if show == nil {
		return generateShowNotFoundMessage(showName)
	}

	ctx.show = *ctx.api.GetShowFullInfo(show.Id)
	ctx.state = SHOW_STATE

	return generateGetShowInfoMessage(*show)
}

func executeGetSeasonInfo(ctx *context) []byte {
	ctx.state = SHOW_STATE
	if ctx.insuficientInfo {
		return generateInsuficientInfoMessage()
	}

	s := ctx.parameters[0]
	season, err := strconv.Atoi(ctx.parameters[0])
	ctx.parameters = []string{}

	if  err != nil || season > ctx.show.SeasonCount {
		return generateSeasonNotFoundMessage(s)
	}

	ctx.show.Seasons[season] = *ctx.api.GetSeasonFullInfo(ctx.show.Id, season)

	ctx.state = SEASON_STATE
	ctx.season = season
	return generateGetSeasonInfoMessage(ctx.show.Seasons[season])
}

func executeGetEpisodeInfo(ctx *context) []byte {
	ctx.state = SEASON_STATE
	if ctx.insuficientInfo {
		return generateInsuficientInfoMessage()
	}

	s := ctx.parameters[0]
	episode, err := strconv.Atoi(ctx.parameters[0])
	ctx.parameters = []string{}

	var max_ratio = 0.0
	var max_index = -1
	if err != nil {
		//we were given a name we gotta get the id for that episode

		for ind, epi := range ctx.show.Seasons[ctx.season].Episodes {
			va := strings.ToLower(epi.Name)
			ratio := smetrics.JaroWinkler(va, s, 0.7, 2)
			fmt.Println(ratio)
			if ratio > max_ratio {
				max_ratio = ratio
				max_index = ind
			}
		}

	}

	if (err == nil && (episode <= 0 || episode > len(ctx.show.Seasons[ctx.season].Episodes) - 1)) || (max_ratio < likelihoodThreshold && max_index != -1)  {
		return generateEpisodeNotFoundMessage(ctx.show.Seasons[ctx.season].Name, s)
	}

	if max_index != -1 {
		ctx.episode = max_index
	} else {
		ctx.episode = episode - 1
	}

	ctx.state = EPISODE_STATE
	return generateGetEpisodeInfoMessage(ctx.show.Seasons[ctx.season].Name, ctx.show.Seasons[ctx.season].Episodes[ctx.episode])
}

func executeGetEpisodeList(ctx *context) []byte {
	ctx.state = SEASON_STATE
	return generateGetEpisodeListMessage(ctx.show.Name, ctx.show.Seasons[ctx.season].Name, ctx.show.Seasons[ctx.season].Episodes)
}

func executeGetShowEpisodeCount(ctx *context) []byte {
	ctx.state = QUANTITY_STATE
	return generateEpisodeCountMessage(ctx.show.Name, ctx.show.EpisodeCount)
}

func executeGetShowSeasonCount(ctx *context) []byte {
	ctx.state = QUANTITY_STATE
	return generateSeasonCountMessage(ctx.show.Name, ctx.show.SeasonCount)
}

func executeGetShowStatus(ctx *context) []byte {
	ctx.state = SHOW_STATE
	return generateShowStatusMessage(ctx.show.Name, ctx.show.Status)
}

func executeGetMostPopularEpisode(ctx *context) []byte {
	ctx.state = SEASON_STATE
	var popu_epi string
	var max_popu = 0.0
	for _, epi := range ctx.show.Seasons[ctx.season].Episodes {
		if epi.Popularity > max_popu {
			max_popu = epi.Popularity
			popu_epi = epi.Name
		}
	}

	return generateMostPopularEpisodeMessage(ctx.show.Seasons[ctx.season].Name, popu_epi)
}

func executeGetMostValuedEpisode(ctx *context) []byte {
	ctx.state = SEASON_STATE
	var popu_epi string
	var max_popu = 0.0
	for _, epi := range ctx.show.Seasons[ctx.season].Episodes {
		if epi.Score > max_popu {
			max_popu = epi.Score
			popu_epi = epi.Name
		}
	}

	return generateMostValuedEpisodeMessage(ctx.show.Seasons[ctx.season].Name, popu_epi)
}

var cTree = conversationTree{
	IDLE_STATE: []conversationPath{
		{directors: []directorOptions{
				{
					"popular",
					"trend",
					"fire",
				},
				{
					"show",
					"tv-show",
				},
			},
		executor: executeGetPopularShows,
		},
		{directors: []directorOptions{
			{
				"weather",
				"clima",
				"temperature",
			},
			{
				"how",
				"which",
				"is",
			},
		},
			executor: executeSendToFuckOff,
		},
		{directors: []directorOptions{
				{
					"how",
				},
				{
					"are",
					"mood",
				},
				{
					"you",
				},
			},
		executor: executeTellMyMood,
		},
		{directors: []directorOptions{
				{
					"hi",
					"hello",
					"morn",
					"afternoon",
					"hey",
				},
			},
		executor: executeSayHiBack,
		},
		{directors: []directorOptions{
				{
					"tell",
					"explain",
					"inform",
				},
				{
					"about",
				},
			},
		executor: executeGetShowInfo,
		prefixes: []parameterPrefixOptions{
				{
					"about",
				},
			},
		},

	},
	SHOW_STATE: []conversationPath{
		{directors: []directorOptions{
				{
					"how",
					"what",
				},
				{
					"much",
					"mani",
					"count",
					"number",
				},
				{
					"episod",
				},
			},
		executor: executeGetShowEpisodeCount,
		},
		{directors: []directorOptions{
				{
					"which",
					"what",
					"is",
				},
				{
					"status",
					"state",
					"product",
					"return",
					"plan",
					"end",
					"pilot",
					"cancel",
				},
			},
		executor: executeGetShowStatus,
		},
		{directors: []directorOptions{
				{
					"how",
					"what",
				},
				{
					"much",
					"mani",
					"count",
					"number",
				},
				{
					"season",
				},
			},
			executor: executeGetShowSeasonCount,
		},
		{directors: []directorOptions{
				{
					"tell",
					"explain",
					"inform",
				},
				{
					"about",
				},
				{
					"season",
				},
			},
		executor: executeGetSeasonInfo,
		prefixes: []parameterPrefixOptions{
				{
					"season",
				},
			},
		},

	},
	QUANTITY_STATE: []conversationPath{
		{directors: []directorOptions{
				{
					"and",
					"nd",
					"&",

				},
				{
					"season",
				},
			},
		executor: executeGetShowSeasonCount,
		},
		{directors: []directorOptions{
				{
					"and",
					"nd",
					"&",
				},
				{
					"episod",
				},
			},
		executor: executeGetShowEpisodeCount,
		},
	},
	SEASON_STATE: []conversationPath{
		{directors: []directorOptions{
				{
					"tell",
					"explain",
					"inform",
				},
				{
					"about",
				},
			},
		executor: executeGetEpisodeInfo,
		prefixes: []parameterPrefixOptions{
				{
					"episod",
					"about",
				},
			},
		},
		{directors: []directorOptions{
				{
					"list",
					"enumerat",
				},
				{
					"episod",
				},
			},
		executor: executeGetEpisodeList,
		},
		{directors: []directorOptions{
				{
					"popular",
				},
				{
					"episod",
				},
			},
			executor: executeGetMostPopularEpisode,
		},
		{directors: []directorOptions{
				{
					"valued",
					"best",
					"coolest",
					"greatest",
				},
				{
					"episod",
				},
			},
			executor: executeGetMostValuedEpisode,
		},
	},
}
