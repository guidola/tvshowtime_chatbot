package engine

import "github.com/guidola/tvshowtime_chatbot/engine/model"
import (
	"github.com/dghubble/sling"
	"net/http"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"bytes"
	"strconv"
)

const tmvdbBaseUrl = "https://api.themoviedb.org/3/"
const popularShowsPath = "tv/popular"
const searchShowPath = "search/tv"
const showInfoPath = "tv/"
const seasonInfoPath = "tv/%d/season/%d"
type baseShowsParams struct {
	ApiKey string `url:"api_key"`
	Language string `url:"language"`
	Page int `url:"page"`
}

type searchShowParams struct {
	ApiKey string `url:"api_key"`
	Language string `url:"language"`
	Page int `url:"page"`
	Query string `url:"query"`
}

type ShowListResponse struct {
	Results []model.Show `json:"results"`
}

type tmvdbApi struct {
	ApiKey string
	Client *http.Client
	Base *sling.Sling
}


func(t *tmvdbApi) Init(key string) {

	t.ApiKey = key
	t.Client = http.DefaultClient
	t.Base = sling.New().Base(tmvdbBaseUrl).Client(t.Client)

}

func(t tmvdbApi) GetPopularShows() []model.Show {

	req, _ := t.Base.New().Get(popularShowsPath).QueryStruct(
		baseShowsParams{t.ApiKey, "en-US", 1}).Request()

	res, _ := t.Client.Do(req)
	var body []byte
	body, _ = ioutil.ReadAll(res.Body)
	pRes := &ShowListResponse{}
	json.Unmarshal(body, pRes)
	return pRes.Results[0:3]
}

func(t tmvdbApi) GetShowInfo(name string) *model.Show {
	req, _ := t.Base.New().Get(searchShowPath).QueryStruct(
		searchShowParams{t.ApiKey, "en-US", 1, name}).Request()

	res, _ := t.Client.Do(req)
	var body []byte
	body, _ = ioutil.ReadAll(res.Body)
	fmt.Println(bytes.NewBuffer(body).String())
	pRes := &ShowListResponse{}
	json.Unmarshal(body, pRes)
	if len(pRes.Results) > 0  {
		return &(pRes.Results[0])
	} else {
		return nil
	}
}

func(t tmvdbApi) GetShowFullInfo(id int) *model.Show {
	req, _ := t.Base.New().Get(showInfoPath + strconv.Itoa(id)).QueryStruct(
		baseShowsParams{t.ApiKey, "en-US", 1}).Request()


	res, _ := t.Client.Do(req)
	var body []byte
	body, _ = ioutil.ReadAll(res.Body)
	fmt.Println(bytes.NewBuffer(body).String())
	pRes := &model.Show{}
	json.Unmarshal(body, pRes)
	return pRes
}

func(t tmvdbApi) GetSeasonFullInfo(show int, season int) *model.Season {

	req, _ := t.Base.New().Get(fmt.Sprintf(seasonInfoPath, show, season)).QueryStruct(
		baseShowsParams{t.ApiKey, "en-US", 1}).Request()
	fmt.Println(req)

	res, _ := t.Client.Do(req)
	var body []byte
	body, _ = ioutil.ReadAll(res.Body)
	fmt.Println(bytes.NewBuffer(body).String())
	pRes := &model.Season{}
	json.Unmarshal(body, pRes)
	return pRes

}









