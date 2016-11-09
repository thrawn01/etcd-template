package etcdTemplate_test

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/thrawn01/args"
	"github.com/thrawn01/etcd-template"
)

var _ = Describe("Generate", func() {
	var dir string

	BeforeEach(func() {
		var err error
		dir, err = ioutil.TempDir("/tmp", "etcd-template-test")
		if err != nil {
			Fail(err.Error())
		}
	})

	AfterEach(func() {
		//os.RemoveAll(dir)
	})

	Describe("Generate()", func() {
		It("Should generate files from templates", func() {
			parser := args.NewParser()
			options := parser.NewOptionsFromMap(map[string]interface{}{
				"template-dir": dir,
				"output-dir":   dir,
			})

			fileName := filepath.Join(dir, "test1")
			content := []byte(`key={{ .value }}`)
			err := ioutil.WriteFile(fileName+".tpl", content, 0666)
			Expect(err).To(BeNil())

			jsonBytes, err := json.Marshal(map[string]interface{}{"value": 1})
			Expect(err).To(BeNil())

			err = etcdTemplate.Generate(options, etcdTemplate.Pair{
				Key:   "/test-key",
				Value: string(jsonBytes),
			})
			Expect(err).To(BeNil())
			content, err = ioutil.ReadFile(fileName)
			Expect(err).To(BeNil())
			Expect(string(content)).To(Equal("key=1"))
		})
	})
})
