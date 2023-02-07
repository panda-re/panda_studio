package main

import (
	"encoding/json"
	"fmt"

	"github.com/panda-re/panda_studio/internal/db/models"
)

type abc struct {
	Type string
}

func main() {
	jsonArrRaw := `[
		{
			"id": "63d5955ed14c76798cf58c58",
			"name": "Test Program",
			"instructions": [
				{
					"type": "command",
					"command": "touch hello123.txt"
				},
				{
					"type": "start_recording",
					"recording_name": "test_recording123"
				}
			]
		}
	]`

	var elements []models.InteractionProgram

	err := json.Unmarshal([]byte(jsonArrRaw), &elements)
	if err != nil {
		panic(err)
	}

	for _, prog := range elements {
		fmt.Printf("%+v\n", prog)
		for _, instr := range prog.Instructions {
			fmt.Printf("%T %+v\n", instr, instr)
		}
	}

	serialized, err := json.MarshalIndent(elements, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(serialized))
}