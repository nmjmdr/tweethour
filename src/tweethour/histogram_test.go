package tweethour

import (
	"testing"
	"time"
)

type mockTimeline struct {
	getTweets  []Tweet
	nextTweets []Tweet

	// these are applied to nextTweets, so that Next function returns in multiple steps
	step     int
	lastStep int
}

func (mt *mockTimeline) Get(username string) ([]Tweet, Error) {
	return mt.getTweets, nil
}

func (mt *mockTimeline) Next(username string, lastId uint64, sinceId uint64) ([]Tweet, Error) {

	if mt.step >= len(mt.nextTweets) {
		mt.step = len(mt.nextTweets)
	}

	tweets := mt.nextTweets[mt.lastStep:mt.step]

	mt.lastStep = mt.step
	mt.step = (mt.step + mt.step)

	return tweets, nil
}

func Test_Get(t *testing.T) {
	mockLine := new(mockTimeline)
	h, err := NewHistogramWithTimeline(mockLine)

	if err != nil {
		t.Fatal(err)
	}

	mockLine.getTweets = make([]Tweet, 24)

	// a tweet per hour, except for the last

	now := time.Now().UTC()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	for i := 0; i < len(mockLine.getTweets); i++ {
		mockLine.getTweets[i].CreatedAt = start.Add(time.Duration(i) * time.Hour).UTC()

		if i == len(mockLine.getTweets)-1 {
			mockLine.getTweets[i].CreatedAt = time.Now().Add(-24 * time.Hour).UTC()
		}
	}

	tByHour, err := h.Get("user")

	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i > len(mockLine.getTweets); i-- {
		if tByHour[i].Count != 1 {
			t.Fatal("Did not get the expected number of tweets per hour, should have been 1 tweet")
		}
	}

	if tByHour[len(mockLine.getTweets)-1].Count != 0 {
		t.Fatal("Did not get the expected number of tweets per hour, should have been 0 tweets")
	}

}

func Test_TodayTweetsFirstGet(t *testing.T) {

	mockLine := new(mockTimeline)
	h, err := NewHistogramWithTimeline(mockLine)

	if err != nil {
		t.Fatal(err)
	}

	mockLine.getTweets = make([]Tweet, 4)

	// setup tweets such that only the first three tweets are marked
	// as todays tweets
	for i := 0; i < len(mockLine.getTweets); i++ {
		mockLine.getTweets[i].CreatedAt = time.Now().UTC()

		if i == len(mockLine.getTweets)-1 {
			mockLine.getTweets[i].CreatedAt = time.Now().Add(-24 * time.Hour).UTC()
		}
	}

	tweets, err := h.todaysTweets("user")

	if len(tweets) != 3 {
		t.Fatal("number of tweets is not equal to expected value")
	}

}

func Test_TodayTweetsGetNext(t *testing.T) {

	mockLine := new(mockTimeline)
	h, err := NewHistogramWithTimeline(mockLine)

	if err != nil {
		t.Fatal(err)
	}

	totalTweets := 8
	mockLine.getTweets = make([]Tweet, totalTweets/2)

	// setup tweets such that only the first three tweets are marked
	// as todays tweets
	for i := 0; i < len(mockLine.getTweets); i++ {
		mockLine.getTweets[i].CreatedAt = time.Now().UTC()
	}

	mockLine.nextTweets = make([]Tweet, totalTweets/2)

	mockLine.step = len(mockLine.nextTweets)

	// setup tweets such that only the first three tweets are marked
	// as todays tweets
	for i := 0; i < len(mockLine.nextTweets); i++ {
		mockLine.nextTweets[i].CreatedAt = time.Now().UTC()

		if i == len(mockLine.nextTweets)-1 {
			mockLine.nextTweets[i].CreatedAt = time.Now().Add(-24 * time.Hour).UTC()
		}
	}

	tweets, err := h.todaysTweets("user")

	if len(tweets) != totalTweets-1 {
		t.Fatal("number of tweets is not equal to expected value")
	}

}

func Test_TodayTweetsGetNextInSteps(t *testing.T) {

	mockLine := new(mockTimeline)
	h, err := NewHistogramWithTimeline(mockLine)

	if err != nil {
		t.Fatal(err)
	}

	totalTweets := 50
	mockLine.getTweets = make([]Tweet, totalTweets/5)

	// setup tweets such that only the first three tweets are marked
	// as todays tweets
	for i := 0; i < len(mockLine.getTweets); i++ {
		mockLine.getTweets[i].CreatedAt = time.Now().UTC()
	}

	mockLine.nextTweets = make([]Tweet, totalTweets-len(mockLine.getTweets))
	mockLine.step = 10
	mockLine.lastStep = 0

	// setup tweets such that only the first three tweets are marked
	// as todays tweets
	for i := 0; i < len(mockLine.nextTweets); i++ {
		mockLine.nextTweets[i].CreatedAt = time.Now().UTC()

		if i == len(mockLine.nextTweets)-1 {
			mockLine.nextTweets[i].CreatedAt = time.Now().Add(-24 * time.Hour).UTC()
		}
	}

	tweets, err := h.todaysTweets("user")

	if len(tweets) != totalTweets-1 {
		t.Fatal("number of tweets is not equal to expected value")
	}

}
