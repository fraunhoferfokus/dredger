package main

import (
	fs "dredger/fileUtils"
	"fmt"
	"log"
	"os"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

/*
func main() {
	// Read your AsyncAPI spec (YAML or JSON)
	data, err := os.ReadFile("asyncapi.yaml")
	if err != nil {
		panic(fmt.Errorf("failed to read asyncapi.yaml: %w", err))
	}

	// Parse the document (strict mode: true)
	doc, err := parser.Parse(data, parser.WithStrict(true))
	if err != nil {
		panic(fmt.Errorf("failed to parse AsyncAPI spec: %w", err))
	}

	// Access parsed data
	info := doc.Info()
	fmt.Println("Title:    ", info.Title())
	fmt.Println("Version:  ", info.Version())
	fmt.Println("Channels: ")

	for name, ch := range doc.Channels().All() {
		fmt.Printf("- %s: %s\n", name, ch.Description())
	}
}
*/

func main() {
	//TODO Logging-Mechanismus später hinzufügen
	schemaFile := "./examples/schemas/jsonSchema.json"
	instanceFile := "./examples/asyncapiv3-min.json"

	// Read file content
	data, err := os.ReadFile(instanceFile)
	if err != nil {
		log.Fatal(err)
	}

	err = parser(data, os.Stdout)

	/*// Parse the document (note: parser detects JSON or YAML automatically)
	doc, err := parseschema.NewParser(data)
	if err != nil {
		log.Fatal(err)
	}*/

	// Access document information
	info := doc.Info()
	fmt.Println("AsyncAPI Title:", info.Title())
	fmt.Println("   Version:       ", info.Version())
	fmt.Println("   Description:   ", info.Description())

	// Print available channels
	fmt.Println("Channels:")
	for name, channel := range doc.Channels().All() {
		fmt.Printf("- %s: %s\n", name, channel.Description())
	}

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
