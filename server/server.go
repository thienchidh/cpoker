package main

import (
	"encoding/json"
	"github.com/paulhankin/cpoker"
	"github.com/paulhankin/poker/v2/poker"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	handleRequests()
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/eval", eval).Methods("POST")
	log.Fatal(http.ListenAndServe(":8887", myRouter))
}

func eval(writer http.ResponseWriter, request *http.Request) {
	log.Println("eval")
	// get the body of our POST request
	reqBody, _ := io.ReadAll(request.Body)

	// reqBody json to list
	var list []string
	_ = json.Unmarshal(reqBody, &list)

	log.Printf("%#v\n", list)

	var cards poker.Hand
	for _, s := range list {
		c := poker.NameToCard[s]
		cards = append(cards, c)
	}

	if len(cards) != 13 {
		// respond with error
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	he := cpoker.MaxProdEvaluator{}
	h, _ := cpoker.Play(cards[:13], he)
	log.Printf("%v", h.String())

	// h string to json

	obj := map[string][]string{
		"front":  convertHand3ToListStr(h.Front),
		"middle": convertHand5ToListStr(h.Middle),
		"back":   convertHand5ToListStr(h.Back),
	}

	// respond with result
	_ = json.NewEncoder(writer).Encode(obj)
}

func convertHand5ToListStr(h [5]poker.Card) []string {
	var list []string
	for _, c := range h {
		list = append(list, c.String())
	}
	return list
}

func convertHand3ToListStr(h [3]poker.Card) []string {
	var list []string
	for _, c := range h {
		list = append(list, c.String())
	}
	return list
}
