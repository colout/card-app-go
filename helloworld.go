package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"

	"bitbucket.org/pkg/inflect"
	"github.com/gorilla/mux"
)

//Type for card object

type Card struct {
	Name      string
	Id        string
	Power     string
	Toughness string
	Cost      string
	Colors    []string
	Text      string
	Types     []string
	SubTypes  []string

	//Custom
	PrimaryColor string //not in original json...will add manually
	TypeDesc     string
}

func (b *Card) MarshalJSON() ([]byte, error) {

	//Set Primary Color
	getPrimaryColor := "none"
	if len(b.Colors) > 0 {
		getPrimaryColor = b.Colors[0]
	}

	//Set TypeDesc - friendly description based on type / subtime

	getTypeDesc := inflect.Titleize(strings.Trim(fmt.Sprint(b.Types), "[]"))
	getSubType := inflect.Titleize(strings.Trim(fmt.Sprint(b.SubTypes), "[]"))
	if len(getSubType) != 0 {
		getTypeDesc += " - " + getSubType
	}

	return json.Marshal(struct {
		Name      string   `json:"name"`
		Id        string   `json:"id"`
		Power     string   `json:"power,omitempty"`
		Toughness string   `json:"toughness,omitempty"`
		Cost      string   `json:"cost,omitempty"`
		Colors    []string `json:"colors,omitempty"`
		Text      string   `json:"text"`
		Types     []string `json:"types,omitempty"`
		SubTypes  []string `json:"subtypes,omitempty"`

		//Custom
		PrimaryColor string `json:"primarycolor,omitempty"`
		TypeDesc     string `json:"typedesc,omitempty"`
	}{
		Name:      b.Name,
		Id:        b.Id,
		Power:     b.Power,
		Toughness: b.Toughness,
		Cost:      b.Cost,
		Colors:    b.Colors,
		Text:      b.Text,
		Types:     b.Types,
		SubTypes:  b.SubTypes,

		//Custom
		PrimaryColor: getPrimaryColor,
		TypeDesc:     getTypeDesc,
	})
}

//Useful Functions:

// Remove single item from slice
func RemoveFromSlice(slice []Card, i int) {
	slice = append(slice[:i], slice[i+1:]...)
}

//Functions
func main() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/c/mydeck", DrawCard)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func DrawCard(w http.ResponseWriter, r *http.Request) {

	//URL to make API call
	url := "https://api.deckbrew.com/mtg/cards?page=2"

	// API call to URL
	resp, err := http.Get(url)

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	//Create array of all cards from API
	var deckcard []Card

	//This is going to be the 5 card hand
	var handcard []Card

	//Convert json to type
	json.Unmarshal(body, &deckcard)

	//Remove 5 random cards
	for i := 1; i <= 5; i++ {
		//item to pop
		var r = rand.Intn(len(deckcard) - 1)

		//Add to hand, remove from deck
		handcard = append(handcard, deckcard[r])
		RemoveFromSlice(deckcard, r)
	}

	//Marshall that shit
	handcardjson, _ := json.MarshalIndent(handcard, "", "    ")
	handcardjsonstring := string(handcardjson)

	fmt.Fprintln(w, "{\"mydeck\":"+handcardjsonstring+"}")
}
