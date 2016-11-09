package etcdTemplate

import (
	"encoding/json"
	"html/template"

	"fmt"
	"io/ioutil"
	"strings"

	"os"
	"path"

	"github.com/pkg/errors"
	"github.com/thrawn01/args"
)

func Generate(options *args.Options, pair Pair) error {

	srcDir := options.String("template-dir")
	outDir := options.String("output-dir")
	if !options.IsSet("output-dir") {
		outDir = srcDir
	}

	// UnMarshal the pair.Value from JSON
	values := make(map[string]interface{})
	err := json.Unmarshal([]byte(pair.Value), &values)
	if err != nil {
		return errors.Wrap(err, "JSON unmarshal error")
	}

	// Search for files in the template directory
	templates, err := getTemplates(srcDir)
	if err != nil {
		return errors.Wrap(err, "Template dir listing error")
	}

	tmpl := template.Must(template.ParseFiles(templates...))

	// Apply the template using the map from JSON
	for _, template := range templates {
		output := strings.Replace(path.Join(outDir, path.Base(template)), ".tpl", "", -1)
		fd, err := os.Create(output)
		if err != nil {
			fd.Close()
			return errors.Wrap(err, fmt.Sprintf("open %s error", template))
		}
		// Execute the template
		if err = tmpl.ExecuteTemplate(fd, path.Base(template), &values); err != nil {
			fd.Close()
			return err
		}
		fd.Sync()
		fd.Close()
	}
	return nil
}

func getTemplates(directory string) ([]string, error) {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".tpl") {
			result = append(result, path.Join(directory, file.Name()))
		}
	}
	return result, nil
}
