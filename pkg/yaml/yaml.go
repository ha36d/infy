package yaml

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

func YamlDecoder(yamlFile string) {

	yfile, err := ioutil.ReadFile(yamlFile)

	if err != nil {

		log.Fatal(err)
	}

	data := make(map[interface{}]interface{})

	err2 := yaml.Unmarshal(yfile, &data)

	if err2 != nil {

		log.Fatal(err2)
	}

	for k, v := range data {

		fmt.Printf("%s -> %d\n", k, v)
	}
}
