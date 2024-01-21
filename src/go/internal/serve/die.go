package serve

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"os"
)

// Representation of a single die
type Die struct {
	Sides int `json:"sides"`
	Value int `json:"value"`
}

// A request to roll a single die
type DieRoll struct {
	Sides int `json:"sides"`
}

func ServeDie() error {
	http.HandleFunc("/roll_die", handleDie)
	return http.ListenAndServe(dieAddr(), nil)
}

func handleDie(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	switch r.Method {
	case "POST":
		decoder := json.NewDecoder(r.Body)
		roll := &DieRoll{}

		err := decoder.Decode(roll)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad roll_die request"))
			return
		}

		encoder.Encode(Die{
			Sides: roll.Sides,
			Value: rollN(roll.Sides),
		})

	case "GET":
		sides := 6
		encoder.Encode(&Die{
			Sides: sides,
			Value: rollN(sides),
		})
	}
}

func rollN(sides int) int {
	return rand.Intn(sides) + 1
}

func dieAddr() string {
	addr, ok := os.LookupEnv("DIE_ADDRESS")

	if !ok {
		addr = "localhost:6666"
	}

	return addr
}
