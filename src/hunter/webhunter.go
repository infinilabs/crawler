/** 
 * User: Medcl
 * Date: 13-7-8
 * Time: 下午5:42 
 */
package hunter

import (
	"net/http"
	"log"
	"io/ioutil"
	"regexp"
)


type  Task struct{
  Url,Request,Response []byte
}

func fetchUrl(url []byte,success chan Task,failure chan string){
	resource := string(url)
	defer func () {
		failure <- resource
	}()


	resp, err := http.Get(resource)
	if err != nil {
		log.Println("we have an error!: ", err)
		return
	}
	defer resp.Body.Close()
	log.Printf("getting %v\n", resource)
	body, _ := ioutil.ReadAll(resp.Body)
	task := Task{url,nil, body}
//	log.Printf("Response %v\n",string(body))
	success <- task

}


func ThrottledCrawl(curl chan []byte, success chan Task, failure chan string, visited map[string]int) {
	maxGos := 10
	numGos := 0
	for {
		if numGos > maxGos {
			<-failure
			numGos -= 1
		}
		url := string(<-curl)
		if _, ok := visited[url]; !ok {
			go fetchUrl([]byte(url), success, failure)
			numGos += 1
			log.Print("gos +1, ",numGos);
		}
		visited[url] += 1
	}
}

func Seed(curl chan []byte,seed string) {
	curl <- []byte(seed)
}

func GetUrls(curl chan []byte, task Task, regex *regexp.Regexp) {
	log.Print("parsing external links:",string(task.Url))
	matches := regex.FindAllSubmatch(task.Response, -1)
	for _, match := range matches {
		curl <- match[1]
	}
}
