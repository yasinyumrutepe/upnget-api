package secret

import (
	"encoding/json"
	"fmt"
	"os"
)

var (
	FullEnv map[string]map[string]any
	Env     map[string]interface{}
)

func LoadEnv(envName string, globalType bool) map[string]map[string]any {
	jsonEnv, err := os.ReadFile("./secret/.env.json")
	if err != nil {
		fmt.Println("Error loading .env file")
		panic("Error loading .env file")
	}
	if err = json.Unmarshal(jsonEnv, &FullEnv); err != nil {
		fmt.Println("Error unmarshaling .env file")
	}
	Env = FullEnv[envName]
	return FullEnv
}
