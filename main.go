package main

import (
	"fmt"
	fb "github.com/huandu/facebook"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type User struct {
	name string
	id   string
}

type Post struct {
	id           string
	message      string
	created_time time.Time
	update_time  time.Time
	from         User
}

const (
	GroupID = "1601597320127277"
	Secreet = ""
)

// 1601597320127277/feed/?metadata=1
// 1601597320127277/feed/?fields=id,created_time,updated_time,from,message&limit=1000&since=today&until=tomorrow

func GroupPostReader(group_id, secrete string,
	ticker *time.Ticker,
	post chan Post,
	done chan bool) {
	defer log.Println("Ticker stopped")

	request_url := fmt.Sprintf("%s/feed/", group_id)

loop:
	for {
		select {
		case <-ticker.C:
			log.Println("Time tick")

			response, _ := fb.Get(request_url, fb.Params{
				"fields":       "id,created_time,updated_time,from,message",
				"limit":        1000,
				"since":        "today",
				"until":        "tomorrow",
				"access_token": secrete,
			})

			log.Println(response["data"])
		case <-done:
			log.Println("Done swithc case")
			break loop
		}
	}

}

// func cleanup() (done chan bool) {
// 	done = make(chan bool)
//
// 	signals := make(chan os.Signal, 2)
// 	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
//
// 	go func() {
// 		sig := <-signals
// 		log.Println(sig)
// 		log.Println("Clining up")
// 		done <- true
// 	}()
//
// 	return done
// }

func main() {
	log.Println("Start of application")

	done := make(chan bool)
	ticker := time.NewTicker(2 * time.Second)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signals

		log.Println("Before ticker stop and closing done channel")
		ticker.Stop()
		close(done)
		log.Println("After ticker stop and closing done channel")
		time.Sleep(10 * time.Second)

		os.Exit(1)
	}()

	channel_post := make(chan Post)
	go GroupPostReader(GroupID, Secreet, ticker, channel_post, done)

	for post := range channel_post {
		log.Println(post)
	}
}
