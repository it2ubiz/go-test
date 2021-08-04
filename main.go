package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

type responseOutput struct {
	H string  `json:"H"`
	K float32 `json:"K"`
}

type requestInput struct {
	A bool    `json:"A"`
	B bool    `json:"B"`
	C bool    `json:"C"`
	D float32 `json:"D"`
	E int     `json:"E"`
	F int     `json:"F"`
}

var custom string
var output responseOutput

// ConditionCheck is Conditions check function - to make output values from input
func ConditionCheck(input requestInput) responseOutput {
	//Basic rules
	var vH string
	var vK float32
	if input.A && input.B && !input.C {
		if custom == "2" {
			vH = "T"
		} else {
			vH = "M"
		}
	} else if input.A && input.B && input.C {
		vH = "P"
	} else if !input.A && input.B && input.C {
		vH = "T"
	} else if input.A && !input.B && input.C && custom == "2" {
		vH = "M"
	} else {
		log.Println("Input error case: A is " + strconv.FormatBool(input.A) + " B is " + strconv.FormatBool(input.B) + " C is " + strconv.FormatBool(input.C))
		output = responseOutput{}
	}

	switch vH {
	case "M":
		if custom == "2" {
			vK = float32(input.F) + input.D + (input.D * float32(input.E) / 100)
		} else {
			vK = input.D + (input.D * float32(input.E) / 10)
		}
	case "P":
		if custom == "1" {
			vK = 2*input.D + (input.D * float32(input.E) / 100)
		} else {
			vK = input.D + (input.D * (float32(input.E) - float32(input.F)) / 25.5)
		}
	case "T":
		vK = input.D - (input.D * float32(input.F) / 30)
	}
	output = responseOutput{H: vH, K: vK}
	return output
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody requestInput
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if err != nil {
		panic(err)
	}
	log.Println(reqBody)
	cond := ConditionCheck(reqBody)
	if (cond == responseOutput{}) {
		http.Error(w, "Incorrect input", http.StatusInternalServerError)
	} else {
		json.NewEncoder(w).Encode(cond)
	}
}

func handleRequests(router *mux.Router) {
	router.HandleFunc("/", mainHandler).Methods("POST", "OPTIONS")
}

func main() {
	router := mux.NewRouter()
	handleRequests(router)
	args := os.Args[1:]
	log.Println(args)
	if len(args) != 0 {
		custom = args[0]
	}

	if custom == "1" || custom == "2" {
		log.Println("Server started with custom: " + custom)
	} else {
		log.Println("Server started with basic rules")
	}
	log.Fatal(http.ListenAndServe(":8000", router))
}
