package serve

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// Request to roll multiple dice
type DiceRoll struct {
	// the number of sides for each die rolled
	Sides []int `json:"sides"`
}

type Dice struct {
	Dice []Die `json:"dice"`
}


func ServeDice() error {
	http.HandleFunc("/roll_dice", handleDice)
	return http.ListenAndServe(diceAddr(), nil)
}

func handleDice(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	roll := &DiceRoll{}

	switch r.Method {
	case "POST":
		decoder := json.NewDecoder(r.Body)

		err := decoder.Decode(roll)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad roll_dice request"))
			return
		}


	case "GET":
		roll.Sides = []int{6, 6}
	}

	dice := make([]Die, len(roll.Sides))

	status := http.StatusOK

	for i := 0; i < len(roll.Sides); i++ {
		dice[i].Sides = roll.Sides[i]
		d, err := rollDie(roll.Sides[i])
		if err != nil {
			status = http.StatusInternalServerError
		} else {
			dice[i].Value = d.Value
		}
	}

	result := Dice{
		Dice: dice,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	encoder.Encode(result)
}

func rollDie(sides int) (Die, error) {
	roll := DieRoll{Sides: sides}

	d := Die{}

	rollBytes, err := json.Marshal(roll)

	if err != nil {
		return d, err
	}

	client := http.DefaultClient

	response, err := client.Post(fmt.Sprintf("http://%s/roll_die",
		dieAddr()),
		"application/json",
		bytes.NewReader(rollBytes))

	decoder := json.NewDecoder(response.Body)

	err = decoder.Decode(&d)

	return d, err
}

func diceAddr() string {
	addr, ok := os.LookupEnv("DICE_ADDRESS")

	if !ok {
		addr = "localhost:6667"
	}

	return addr
}
