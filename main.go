package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	mgo "gopkg.in/mgo.v2"

	"github.com/PuerkitoBio/goquery"
)

type Connection struct {
	DB *mgo.Database
}

type News struct {
	Title    string
	Link     string
	HasTrump bool
	Date     time.Time
}

const (
	MongoDb = "trumps"
	Substr  = "Trump"
)

var (
	conn      Connection
	numTrumps = 0
)

func main() {
	host := os.Getenv("MONGO_DB_HOST")
	session, err := mgo.Dial(host)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	conn = Connection{session.DB(MongoDb)}

	url := "http://omni.se"
	if err = scrape(url); err != nil {
		log.Println(err)
		return
	}

	fmt.Printf("Number of trumps found: %d\n", numTrumps)

}

func scrape(url string) error {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Printf("Error sraping site %s\n")
		return err
	}

	doc.Find("article a.article-link h1").Each(func(index int, item *goquery.Selection) {
		linkTag := item
		link, _ := linkTag.Attr("href")
		linkText := linkTag.Text()

		n := News{
			Title:    linkText,
			Link:     link,
			HasTrump: false,
			Date:     time.Now(),
		}

		if strings.Contains(linkText, Substr) {
			n.HasTrump = true
			numTrumps++
		}

		if err := conn.DB.C("news").Insert(&n); err != nil {
			log.Printf("Error inserting news %s\n", err.Error())
			return
		}
	})

	return nil
}
