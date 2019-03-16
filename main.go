package main

// required package
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

// Define query type
var queryType = graphql.NewObject(graphql.ObjectConfig{
	Name: "HackerNews",
	Fields: graphql.Fields{
		"latestPost": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				fetchStories()
				return "Hello World", nil
			},
		},
	},
})

var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query: queryType,
})

// create story struct
type Story struct {
	by    string
	id    int
	score int
	title string
	url   string
}

// func fetchStoryById(id string) Story {

// }
func fetchStories() {
	// make request to this
	response, err := http.Get("https://hacker-news.firebaseio.com/v0/topstories.json?print=pretty")
	if err != nil {
		fmt.Println(err)
	}

	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	stories := make([]int64, 0)
	json.Unmarshal(body, &stories)
	fmt.Println(stories)

	//for fist 20 make request to fetch stories
}

func main() {
	// create handler
	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	http.Handle("/graphql", h)
	fmt.Println("Listening on port: 3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
