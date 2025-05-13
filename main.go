package main

import (
	fs "dredger/fileUtils"
	"fmt"
	"log"
	"os"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

func main() {
	//TODO Logging-Mechanismus später hinzufügen
	schemaFile := "./examples/schemas/jsonSchema.json"
	instanceFile := "./examples/asyncapiv3-min.json"

	if fs.Check_YAML(schemaFile) {
		fmt.Println(fs.ConvertYAMLIntoJSON(schemaFile))
	}

	c := jsonschema.NewCompiler()
	sch, err := c.Compile(schemaFile)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Open(instanceFile)
	if err != nil {
		log.Fatal(err)
	}
	inst, err := jsonschema.UnmarshalJSON(f)
	if err != nil {
		log.Fatal(err)
	}

	err = sch.Validate(inst)
	fmt.Println("valid:", err == nil)

}
