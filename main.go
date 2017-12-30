package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/graphql-go/graphql"
)

type secondLayer struct {
	UID  string `json:"uid"`
	Type string `json:"type"`
}

type firstLayer struct {
	Name        string      `json:"name"`
	Address     string      `json:"address"`
	Location    string      `json:"location"`
	SecondLayer secondLayer `json:"second_layer"`
}

var (
	querySchema graphql.Schema
)

func init() {
	objSecondLayer := graphql.NewObject(graphql.ObjectConfig{
		Name:        "SecondLayer",
		Description: "second layer.",
		Fields: graphql.Fields{
			"uid": &graphql.Field{
				Type:        graphql.String,
				Description: "uid",
			},
			"type": &graphql.Field{
				Type:        graphql.String,
				Description: "type.",
			},
		},
	})

	objFirstLayer := graphql.NewObject(graphql.ObjectConfig{
		Name:        "FirstLayer",
		Description: "first layer.",
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type:        graphql.String,
				Description: "name",
			},
			"address": &graphql.Field{
				Type:        graphql.String,
				Description: "address.",
			},
			"location": &graphql.Field{
				Type:        graphql.String,
				Description: "location.",
			},
			"second_layer": &graphql.Field{
				Type:        objSecondLayer,
				Description: "second_layer.",
			},
		},
	})

	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"get_map": &graphql.Field{
				Type: objFirstLayer,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					result := map[string]interface{}{"name": "Map Name", "address": "Map Address", "location": "Map Location", "second_layer": map[string]interface{}{"uid": "Map UID", "type": "Map Type"}}
					objJSON, _ := json.Marshal(result)
					firstObj := firstLayer{}
					json.Unmarshal(objJSON, &firstObj)
					return firstObj, nil
				},
			},
			"get_struct": &graphql.Field{
				Type: objFirstLayer,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					secondObj := secondLayer{UID: "Struct UID", Type: "Struct Type"}
					firstObj := firstLayer{Name: "Struct Name", Address: "Struct Address", Location: "Struct Location", SecondLayer: secondObj}
					return firstObj, nil
				},
			},
		},
	})

	querySchema, _ = graphql.NewSchema(graphql.SchemaConfig{
		Query: queryType,
	})
}

func queryGqlHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var (
		query string
	)

	if r.Method == "GET" {
		query = r.URL.Query().Get("query")
	} else {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		query = string(body)
	}

	result := graphql.Do(graphql.Params{
		Schema:        querySchema,
		RequestString: query,
	})
	json.NewEncoder(w).Encode(result)

}

func main() {
	log.Println("Test with Get      : curl -XPOST -d 'query {get_map{name,address,location,second_layer{uid,type}},get_struct{name,address,location,second_layer{uid,type}}}' http://localhost:8080/api")
	http.HandleFunc("/api", queryGqlHandler)

	http.ListenAndServe(":8080", nil)
}
