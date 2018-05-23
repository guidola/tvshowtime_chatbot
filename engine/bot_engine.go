package engine

import (
	"github.com/kljensen/snowball"
	"fmt"
	"strings"
	"github.com/xrash/smetrics"
	"github.com/guidola/tvshowtime_chatbot/engine/model"
	"github.com/SimonWaldherr/golibs/xmath"
)

const likelihoodThreshold = 0.85

type context struct {
	state int 		// defines the state of the conversation we're at. Considered when processing the input.
					// allows to fetch the proper intent to generate from a multi-level dictionary which actually defines
					// the conversation tree

	parameters []string 		// current conversation flow relevant parameters that may be used for queries or
								// generating responses
	insuficientInfo bool		// states if the input had enough parameters to perform the query.
	input string
	show model.Show
	season int
	episode int
	words []string
	api		tmvdbApi
}

type BotEngine struct {
	context context
}


func simplifyInput(s string) []string {

	// split the input in words
	words := strings.Split(s, " ")

	// steem the input so we're left with the word's roots
	var steemedWords []string
	for _, word := range words {
		stemmed, _ := snowball.Stem(word, "english", true)
		steemedWords = append(steemedWords, stemmed)
	}

	return steemedWords
}


func evaluateLikelihood(directors []directorOptions, words []string) float64 {

	// compare each director against all the words in the input, stick with the highest match for that director
	likelihoods := []float64{}
	for _, director := range directors {

		var max_ratio = 0.0
		for _, option := range director {
			for _, word := range words {
				ratio := smetrics.JaroWinkler(option, word, 0.7, 2)
				if ratio > max_ratio {
					max_ratio = ratio
				}
			}
		}

		likelihoods = append(likelihoods, max_ratio)
	}

	return xmath.Geometric(likelihoods)
}

func extractAndSetContext(context *context, prefixes []parameterPrefixOptions, words []string) {

	context.insuficientInfo = false

	for _, prefix := range prefixes {

		var max_ratio = 0.0
		var max_index = 0
		for _, option := range prefix {
			for word_i, word := range words {
				ratio := smetrics.JaroWinkler(option, word, 0.7, 2)
				if ratio > max_ratio {
					max_ratio = ratio
					max_index = word_i
				}
			}
		}

		if max_ratio > likelihoodThreshold && max_index < len(words) - 1 {
			context.parameters = append(context.parameters, strings.TrimSpace(strings.Split(context.input, context.words[max_index])[1]))
			a := 0
			a = a
		} else {
			context.insuficientInfo = true
			return
		}
	}


}

func computeIntent(context *context, words []string) intent {

	currentConvState := cTree[context.state]

	var maxMatchRatio float64 = 0
	var mostLikelyPathIndex int
	for index, path := range currentConvState {
		match_ratio := evaluateLikelihood(path.directors, words)
		if match_ratio > maxMatchRatio {
			mostLikelyPathIndex = index
			maxMatchRatio = match_ratio
		}
	}
	fmt.Println(maxMatchRatio)
	if maxMatchRatio > likelihoodThreshold {
		extractAndSetContext(context, currentConvState[mostLikelyPathIndex].prefixes, words)
		return currentConvState[mostLikelyPathIndex].executor
	} else if context.state == QUANTITY_STATE {

		context.state = SHOW_STATE
		ret := computeIntent(context, words)
		context.state = QUANTITY_STATE

		return ret

	} else if context.state == SEASON_STATE {

		context.state = SHOW_STATE
		ret := computeIntent(context, words)
		context.state = SEASON_STATE

		return ret

	} else if context.state == EPISODE_STATE {

		context.state = SEASON_STATE
		ret := computeIntent(context, words)
		context.state = EPISODE_STATE

		return ret

	} else if context.state != IDLE_STATE {

			stt := context.state
			context.state = IDLE_STATE
			ret := computeIntent(context, words)
			context.state = stt

			return ret


	} else {

		return defaultExecutor
	}
}


func (e *BotEngine) compute(in string) []byte {

	words := simplifyInput(in)
	fmt.Println(words)
	e.context.input = in
	e.context.words = strings.Split(in, " ")

	return computeIntent(&(e.context), words)(&(e.context))
}
