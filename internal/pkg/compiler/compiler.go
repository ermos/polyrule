package compiler

import (
	"fmt"
	"github.com/ermos/polyrule/internal/pkg/compiler/lang/golang"
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

var compilers = map[string]Lang{
	"js":         js.Lang{},
	"javascript": js.Lang{},
	"php":        php.Lang{},
	"go":         golang.Lang{},
}

var reader = map[string]Reader{
	"json": json.Reader{},
	"toml": toml.Reader{},
}

var isBackupCreated = false
var withBackup = false

func GetAvailableLang() []string {
	keys := make([]string, 0, len(compilers))
	for name := range compilers {
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

	if compilers[lang] == nil {
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

	inputs, err := getInputs(input, info)
	if err != nil {
		panic(err)
	}

	err = process(cmd, lang, inputs, output, input, info)
	if err != nil {
		panic(err)
	}

	removeBackup(output)

	return
}

func process(
	cmd *cobra.Command,
	lang string,
	inputs []string,
	output, basePath string,
	basePathInfo os.FileInfo,
) (err error) {
	for _, i := range inputs {
		var content string

		path := strings.Replace(filepath.Dir(i), basePath, "", 1)
		name := strings.Replace(filepath.Base(i), filepath.Ext(i), "", 1)
		rules := make(map[string]model.Rule)

		inputType := strings.TrimLeft(filepath.Ext(i), ".")
		if reader[inputType] == nil {
			return fmt.Errorf("%s input is not supported", inputType)
		}

		log.Verbose("read file", i)
		rules, err = reader[inputType].Read(i)
		if err != nil {
			return
		}

		for _, r := range rules {
			if !types.IsValidType(r.Type) {
				return fmt.Errorf("undefined type \"%s\" found in %s", r.Type, i)
			}
		}

		log.Verbose("compile file", i)
		compiler := compilers[lang].New(cmd, output, path, name, rules).(Lang)

		content, err = compiler.Compile()
		if err != nil {
			return
		}

		writePath := output
		if basePathInfo.IsDir() {
			writePath = compiler.OutputDirPath()
		}

		err = writeFile(writePath, content)
		if err != nil {
			return
		}
	}

	return
}

func getInputs(basePath string, basePathInfo os.FileInfo) (inputs []string, err error) {
	if basePathInfo.IsDir() {
		log.Verbose("get files from", basePath, "directory")
		return getFiles(basePath)
	}
	return append(inputs, basePath), nil
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
		if os.IsNotExist(err) {
			return nil
		}
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
