package tweethour

import (
	"net/http"
	"time"
)

type TweetsByHour struct {
	// do not serialize the tweets
	tweets []Tweet
	From   time.Time `json:"from-hour"`
	To     time.Time `json:"to-hour"`
	Count  int       `json:"frequency"`
}

type Histogram interface {
	Get(username string) ([]TweetsByHour, Error)
}

type histogram struct {
	line     Timeline
	tokenGen TokenGen
	token    string
	client   *http.Client
}

// mainly used for unit testing
func NewHistogramWithTimeline(timeline Timeline) (*histogram, Error) {
	h := new(histogram)
	h.line = timeline
	return h, nil
}

func NewHistogram() (*histogram, Error) {
	h := new(histogram)
	h.client = &http.Client{}

	h.tokenGen = NewTokenGen(h.client)

	// get the token, when we first start up
	var err Error
	h.token, err = h.tokenGen.Generate()

	if err != nil {
		return nil, err
	}

	h.line = NewTimeline(h.client, h.token)

	return h, nil
}

func (h *histogram) Get(username string) ([]TweetsByHour, Error) {

	tByHour := make([]TweetsByHour, 24)

	now := time.Now().UTC()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	for i := 0; i < len(tByHour); i++ {
		tByHour[i].tweets = make([]Tweet, 0)
		tByHour[i].From = start.Add(time.Duration(i) * time.Hour)
		tByHour[i].To = tByHour[i].From.Add(1 * time.Hour)
	}

	tweets, err := h.todaysTweets(username)

	if err != nil {
		return nil, err
	}

	for _, t := range tweets {
		index := int(t.CreatedAt.Hour())
		tByHour[index].tweets = append(tByHour[index].tweets, t)
		tByHour[index].Count = len(tByHour[index].tweets)
	}

	return tByHour, nil
}

func (h *histogram) todaysTweets(username string) ([]Tweet, Error) {

	now := time.Now().UTC()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	// do a get first
	tweets, err := h.line.Get(username)

	if err != nil {
		return nil, err
	}

	if len(tweets) == 0 {
		return tweets, nil
	}

	getNext := true
	var index int
	var t Tweet

	end := len(tweets)

	for index, t = range tweets {

		if t.CreatedAt.Before(startOfDay) {
			getNext = false
			end = index
			break
		}
	}

	if !getNext {
		return tweets[0:end], nil
	}

	// we need to get next set of tweets, the user has been prolific
	// has tweeted more than 200 times this day

	for getNext {
		var nextTweets []Tweet

		sinceId := tweets[0].Id
		maxId := (tweets[len(tweets)-1].Id - 1)
		nextTweets, err = h.line.Next(username, maxId, sinceId)

		end = len(nextTweets)

		for index, t = range nextTweets {
			if t.CreatedAt.Before(startOfDay) {
				getNext = false
				end = index
				break
			}
		}

		// append nextTweets to tweets
		tweets = append(tweets, nextTweets[0:end]...)

	}

	return tweets[0:], nil
}
