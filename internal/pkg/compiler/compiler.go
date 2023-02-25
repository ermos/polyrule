package compiler

import (
	"fmt"
	"github.com/ermos/polyrule/internal/pkg/compiler/lang/js"
	"github.com/ermos/polyrule/internal/pkg/compiler/reader/json"
	"github.com/ermos/polyrule/internal/pkg/compiler/reader/toml"
	"github.com/ermos/polyrule/internal/pkg/model"
	"os"
	"path/filepath"
	"strings"
)

var compiler = map[string]Lang{
	"js":         js.Lang{},
	"javascript": js.Lang{},
}

var reader = map[string]Reader{
	"json": json.Reader{},
	"toml": toml.Reader{},
}

func GetAvailableLang() []string {
	keys := make([]string, 0, len(compiler))
	for name := range compiler {
		keys = append(keys, name)
	}
	return keys
}

func Compile(lang, input, output string) (err error) {
	if compiler[lang] == nil {
		return fmt.Errorf("%s output is not supported", lang)
	}

	info, err := os.Stat(input)
	if err != nil {
		return
	}

	var inputs []string
	if info.IsDir() {
		inputs, err = getFiles(input)
		if err != nil {
			return
		}
	} else {
		inputs = append(inputs, input)
	}

	list := make(map[string]map[string]model.Rule)

	for _, i := range inputs {
		rules := make(map[string]model.Rule)

		inputNamespace := strings.Replace(filepath.Base(i), filepath.Ext(i), "", 1)

		inputType := strings.TrimLeft(filepath.Ext(i), ".")
		if reader[inputType] == nil {
			return fmt.Errorf("%s input is not supported", inputType)
		}

		rules, err = reader[inputType].Read(i)
		if err != nil {
			return
		}

		list[inputNamespace] = rules
	}

	content, err := compiler[lang].Compile(list)
	if err != nil {
		return
	}

	return writeFile(output, content)
}

func getFiles(dirPath string) (files []string, err error) {
	err = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return
}

func writeFile(path, content string) (err error) {
	file, err := os.Create(path)
	if err != nil {
		return
	}
	defer file.Close()

	_, err = file.WriteString(content)
	return err
}
