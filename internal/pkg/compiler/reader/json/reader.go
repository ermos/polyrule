package json

import (
	"encoding/json"
	"github.com/ermos/polyrule/internal/pkg/model"
	"os"
)

type Reader struct{}

func (Reader) Read(input string) (rules map[string]model.Rule, err error) {
	file, err := os.Open(input)
	if err != nil {
		return
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&rules)
	return
}
