package compiler

import (
	"fmt"
	"github.com/ermos/polyrule/internal/pkg/compiler/lang/js"
	"github.com/ermos/polyrule/internal/pkg/compiler/lang/php"
	"github.com/ermos/polyrule/internal/pkg/compiler/reader/json"
	"github.com/ermos/polyrule/internal/pkg/compiler/reader/toml"
	"github.com/ermos/polyrule/internal/pkg/log"
	"github.com/ermos/polyrule/internal/pkg/model"
	"github.com/ermos/polyrule/internal/pkg/types"
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
			log.Print("Compilation failed :", r)
			err = restoreBackup(output)
			if err != nil {
				log.Print("Backup restoration failed : ", r)
				log.Print("You can find your backup here : ", fmt.Sprintf("%s.backup", output))
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
		log.Verbose("get files from", input, "directory")
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

		log.Verbose("read file", i)
		rules, err = reader[inputType].Read(i)
		if err != nil {
			panic(err)
		}

		for _, r := range rules {
			if !types.IsValidType(r.Type) {
				panic(fmt.Sprintf("undefined type \"%s\" found in %s", r.Type, i))
			}
		}

		log.Verbose("compile file", i)
		content, err = compiler[lang].Compile(cmd, path, name, rules)
		if err != nil {
			panic(err)
		}

		writePath := filepath.Join(output)
		if info.IsDir() {
			writePath = filepath.Join(
				output,
				path,
				fmt.Sprintf("%s.%s", name, compiler[lang].GetExtension()),
			)
		}

		err = writeFile(writePath, content)
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
	log.Verbose("create output directory if not exists for", path)
	err = os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		return
	}

	log.Verbose("create output file", path)
	file, err := os.Create(path)
	if err != nil {
		return
	}
	defer file.Close()

	log.Verbose("write output content in", path)
	_, err = file.WriteString(content)
	return err
}

func createBackup(path string) (err error) {
	if withBackup {
		log.Verbose("create back-up directory :", fmt.Sprintf("%s.backup", path))
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
		log.Verbose("restore back-up :", fmt.Sprintf("%s.backup", path))
		if err = os.RemoveAll(path); err != nil {
			return
		}
		return os.Rename(fmt.Sprintf("%s.backup", path), path)
	}
	return
}

func removeBackup(path string) {
	if withBackup && isBackupCreated {
		log.Verbose("remove back-up :", fmt.Sprintf("%s.backup", path))
		_ = os.RemoveAll(fmt.Sprintf("%s.backup", path))
	}
}
