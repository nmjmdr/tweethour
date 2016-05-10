package tweethour

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const pageSize = 200

const timelineUrl = "https://api.twitter.com/1.1/statuses/user_timeline.json?user_id="
const bearer = "Bearer "

const createdAt = "created_at"
const id = "id_str"
const text = "text"

type Tweet struct {
	CreatedAt time.Time
	Text      string
	Id        uint64
}

type Timeline interface {
	Get(username string) ([]Tweet, Error)
	Next(username string, lastId uint64, sinceId uint64) ([]Tweet, Error)
}

type timeline struct {
	client *http.Client
	token  string
}

func NewTimeline(client *http.Client, token string) Timeline {
	t := new(timeline)
	t.token = token
	t.client = client
	return t
}

func (t *timeline) makeRequest(url string) (*http.Response, Error) {

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add(authHeader, bearer+t.token)

	if err != nil {
		return nil, NewTimelineError(err)
	}

	var response *http.Response
	response, err = t.client.Do(req)

	if err != nil {
		return nil, NewTimelineError(err)
	}

	if response.StatusCode != http.StatusOK {
		return nil, NewStatusError(response.StatusCode, http.StatusText(response.StatusCode))
	}

	return response, nil
}

func (t *timeline) parse(response *http.Response) ([]Tweet, Error) {

	var v []map[string]interface{}
	err := json.NewDecoder(response.Body).Decode(&v)

	if err != nil {
		return nil, NewTimelineError(err)
	}

	if v == nil {
		return nil, NewTimelineError(errors.New("Could not decode the response for timeline"))
	}

	var created, tweetId, tweetText string

	tweets := make([]Tweet, 0)

	for _, r := range v {

		created, err = getValue(createdAt, r)
		if err != nil {
			return nil, NewTimelineError(err)
		}
		tweetId, err = getValue(id, r)
		if err != nil {
			return nil, NewTimelineError(err)
		}
		tweetText, err = getValue(text, r)
		if err != nil {
			return nil, NewTimelineError(err)
		}

		var t *Tweet
		t, err := makeTweet(created, tweetId, tweetText)
		if err != nil {
			return nil, NewTimelineError(err)
		}
		tweets = append(tweets, (*t))
	}

	return tweets, nil
}

func makeTweet(created string, id string, text string) (*Tweet, error) {
	t := Tweet{}
	date, err := time.Parse(time.RubyDate, created)
	if err != nil {
		return nil, err
	}

	t.CreatedAt = date

	var idval uint64
	idval, err = strconv.ParseUint(id, 10, 64)

	if err != nil {
		return nil, err
	}
	t.Id = idval

	t.Text = text

	return &t, nil

}

func (t *timeline) Get(username string) ([]Tweet, Error) {

	url := t.getFirstRequestUrl(username)
	res, err := t.makeRequest(url)

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var tweets []Tweet
	tweets, err = t.parse(res)

	if err != nil {
		return nil, err
	}

	return tweets, nil
}

func (t *timeline) getFirstRequestUrl(username string) string {
	reqUrl := fmt.Sprintf("%s%s&screen_name=%s&count=%d&trim_user=true", timelineUrl, username, username, pageSize)
	return reqUrl
}

func (t *timeline) getNextRequestUrl(username string, maxId uint64, sinceId uint64) string {
	reqUrl := fmt.Sprintf("%s%s&screen_name=%s&count=%d&trim_user=true&since_id=%d&max_id=%d", timelineUrl, username, username, pageSize, sinceId, maxId)
	return reqUrl
}

func (t *timeline) Next(username string, lastId uint64, sinceId uint64) ([]Tweet, Error) {
	url := t.getNextRequestUrl(username, lastId, sinceId)
	res, err := t.makeRequest(url)

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var tweets []Tweet
	tweets, err = t.parse(res)

	if err != nil {
		return nil, err
	}

	return tweets, nil
}
