package testYamlPkg

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// NOTE1: struct fields must be public (Capitalized) otherwise it'll return empty
// 		  field tag is not must-have and could have options
type Person struct {
	Name string `yaml:"name"`
	Favs struct {
		Food []string `yaml:"food,omitempty"`
		Number  int `yaml:"number"`
	}
}

func ReadYamlConfig(path string){
	var result []Person
	// yamlbytes []byte
	// from string to bytes -> byteString := []byte(rawString)
	yamlbytes, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("Error reading config file: %v", err)
	} else {
		// NOTE2: Unmarshal output "result" must be a pointer (address in memory) of the variable
		err = yaml.Unmarshal(yamlbytes, &result)
		if err != nil {
			fmt.Printf("Unmarshal Error: %v\n", err.Error())
		} else {
			for _, p := range result {
				fmt.Printf("Item: %s - %v\n", p.Name, p.Favs)
			}
		}
	}
}