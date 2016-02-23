package parser
import (
	"log"
	"github.com/PuerkitoBio/goquery"
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