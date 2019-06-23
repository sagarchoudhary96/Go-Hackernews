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
			Type: graphql.NewList(storyType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {

				res, err := fetchStories()
				if err != nil {
					return nil, err
				}
				return res, nil
			},
			Description: "Fetch top 20 posts",
		},
		"post": &graphql.Field{
			Type: storyType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type:        graphql.Int,
					Description: "Search for post by ID",
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id := p.Args["id"].(int)
				res, err := fetchStoryByID(int64(id))
				if err != nil {
					return nil, err
				}

				return res, nil
			},
			Description: "Fetch Single Post",
		},
	},
})

// Define StoryType
var storyType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Story",
	Fields: graphql.Fields{
		"By": &graphql.Field{
			Type: graphql.String,
		},
		"ID": &graphql.Field{
			Type: graphql.String,
		},
		"Title": &graphql.Field{
			Type: graphql.String,
		},
		"Score": &graphql.Field{
			Type: graphql.String,
		},
		"URL": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query: queryType,
})

// create story struct
type Story struct {
	By    string `json:"by"`
	ID    int    `json:"id"`
	Score int    `json:"score"`
	Title string `json:"title"`
	URL   string `json:"url"`
}

func fetchStoryByID(id int64) (Story, error) {

	// create story url
	storyURL := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%v.json?print=pretty", id)

	// fetch user story
	response, err := http.Get(storyURL)

	if err != nil {
		return Story{}, err
	}

	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	story := Story{}
	json.Unmarshal(body, &story)
	return story, nil

}
func fetchStories() ([]Story, error) {
	// make request to get stories ids
	response, err := http.Get("https://hacker-news.firebaseio.com/v0/topstories.json?print=pretty")
	if err != nil {
		fmt.Println(err)
	}

	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	storyIds := make([]int64, 0)
	json.Unmarshal(body, &storyIds)

	stories := make([]Story, 0)

	for i := 0; i < 20; i++ {
		story, err := fetchStoryByID(storyIds[i])

		if err != nil {
			return nil, err
		}

		stories = append(stories, story)
	}
	return stories, err
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
