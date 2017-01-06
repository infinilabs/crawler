package parser

import (
	"github.com/PuerkitoBio/goquery"
	"log"
)

func ExampleScrape() {
	_, err := goquery.NewDocument("1.html")
	if err != nil {
		log.Fatal(err)
	}

}

func main() {
	ExampleScrape()
}
