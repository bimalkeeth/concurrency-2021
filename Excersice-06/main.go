package main

import (
	"fmt"
	"github.com/caser/gophernews"
	"os"
	"sync"
)
import "github.com/jzelinskie/geddit"

var redditSession *geddit.LoginSession
var hackerNewsClient *gophernews.Client

func init() {

	hackerNewsClient = gophernews.NewClient()
	var err error
	redditSession, err = geddit.NewLoginSession("g_d_bot", "K417k4FTua52", "gdAgent v0")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type Story struct {
	title  string
	url    string
	author string
	source string
}

func getHnStoryDetail(id int, c chan<- Story, wg *sync.WaitGroup) {
	defer wg.Done()
	story, err := hackerNewsClient.GetStory(id)
	if err != nil {
		return
	}
	newStory := Story{
		title:  story.Title,
		url:    story.URL,
		author: story.By,
		source: "Hacker News",
	}
	c <- newStory
}

func newHnStories(c chan<- Story) {

	defer close(c)

	changes, err := hackerNewsClient.GetChanges()
	if err != nil {
		fmt.Println(err)
	}
	var wg sync.WaitGroup
	for _, id := range changes.Items {
		wg.Add(1)
		go getHnStoryDetail(id, c, &wg)
	}
	wg.Wait()
}

func newRedditStories(c chan<- Story) {

	defer close(c)
	sort := geddit.PopularitySort(geddit.NewSubmissions)
	var listenOption geddit.ListingOptions
	submission, err := redditSession.SubredditSubmissions("programming", sort, listenOption)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, story := range submission {

		newStory := Story{
			title:  story.Title,
			url:    story.URL,
			author: story.Author,
			source: "Reddit /r/programming",
		}
		c <- newStory
	}
}

func outputToConsole(c <-chan Story) {
	for {
		s := <-c
		fmt.Printf("%s: %s\nby %s on %s\n\n", s.title, s.url, s.author, s.source)
	}
}

func outputToFile(c <-chan Story, file *os.File) {
	for {
		s := <-c
		fmt.Fprintf(file, "%s: %s\nby %s on %s\n\n", s.title, s.url, s.author, s.source)
	}
}

func main() {
	fromHn := make(chan Story, 8)
	fromReddit := make(chan Story, 8)
	toFile := make(chan Story, 8)
	toConsole := make(chan Story, 8)

	go newHnStories(fromHn)
	go newRedditStories(fromReddit)

	file, err := os.Create("stories.txt")
	if err != nil {
		fmt.Println("error occurred")
		os.Exit(1)
	}
	go outputToConsole(toConsole)
	go outputToFile(toFile, file)

	hOpen := true
	redditOpen := true

	for hOpen || redditOpen {
		select {
		case story, open := <-fromHn:
			if open {
				toFile <- story
				toConsole <- story
			} else {
				hOpen = false
			}
		case story, open := <-fromReddit:
			if open {
				toFile <- story
				toConsole <- story
			} else {
				hOpen = false
			}
		}
	}

}
