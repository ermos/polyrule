package compiler

import (
	"fmt"
	"github.com/ermos/polyrule/internal/pkg/compiler/lang/js"
	"github.com/ermos/polyrule/internal/pkg/compiler/lang/php"
	"github.com/ermos/polyrule/internal/pkg/compiler/reader/json"
	"github.com/ermos/polyrule/internal/pkg/compiler/reader/toml"
	"github.com/ermos/polyrule/internal/pkg/model"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

var compiler = map[string]Lang{
	"js":         js.Lang{},
	"javascript": js.Lang{},
	"php":        php.Lang{},
}

var reader = map[string]Reader{
	"json": json.Reader{},
	"toml": toml.Reader{},
}

var isBackupCreated = false
var withBackup = false

func GetAvailableLang() []string {
	keys := make([]string, 0, len(compiler))
	for name := range compiler {
		keys = append(keys, name)
	}
	return keys
}

func Compile(cmd *cobra.Command, lang, input, output string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Compilation failed :", r)
			err = restoreBackup(output)
			if err != nil {
				fmt.Println("Backup restoration failed : ", r)
				fmt.Println("You can find your backup here : ", fmt.Sprintf("%s.backup", output))
			}
		}
	}()

	clean, err := cmd.Flags().GetBool("clean")
	if err != nil {
		panic(err)
	}

	if clean {
		withBackup = true
	}

	if compiler[lang] == nil {
		panic(fmt.Sprintf("%s output is not supported", lang))
	}

	err = createBackup(output)
	if err != nil {
		panic(err)
	}

	info, err := os.Stat(input)
	if err != nil {
		panic(err)
	}

	var inputs []string
	if info.IsDir() {
		inputs, err = getFiles(input)
		if err != nil {
			panic(err)
		}
	} else {
		inputs = append(inputs, input)
	}

	for _, i := range inputs {
		var content string

		path := strings.Replace(filepath.Dir(i), input, "", 1)
		name := strings.Replace(filepath.Base(i), filepath.Ext(i), "", 1)
		rules := make(map[string]model.Rule)

		inputType := strings.TrimLeft(filepath.Ext(i), ".")
		if reader[inputType] == nil {
			panic(fmt.Sprintf("%s input is not supported", inputType))
		}

		rules, err = reader[inputType].Read(i)
		if err != nil {
			panic(err)
		}

		content, err = compiler[lang].Compile(cmd, path, name, rules)
		if err != nil {
			panic(err)
		}

		err = writeFile(
			filepath.Join(
				output,
				path,
				fmt.Sprintf("%s.%s", name, compiler[lang].GetExtension()),
			),
			content,
		)
		if err != nil {
			panic(err)
		}
	}

	removeBackup(output)

	return
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
	err = os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		return
	}

	file, err := os.Create(path)
	if err != nil {
		return
	}
	defer file.Close()

	_, err = file.WriteString(content)
	return err
}

func createBackup(path string) (err error) {
	if withBackup {
		err = os.Rename(path, fmt.Sprintf("%s.backup", path))
		if err != nil {
			return
		}
		isBackupCreated = true
	}
	return
}

func restoreBackup(path string) (err error) {
	if withBackup && isBackupCreated {
		err = os.RemoveAll(path)
		if err != nil {
			return
		}

		return os.Rename(fmt.Sprintf("%s.backup", path), path)
	}
	return
}

func removeBackup(path string) {
	if withBackup && isBackupCreated {
		_ = os.RemoveAll(fmt.Sprintf("%s.backup", path))
	}
}
