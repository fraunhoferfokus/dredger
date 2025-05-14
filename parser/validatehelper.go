func validateServerVariables(server map[string]interface{}) error { //Überprüft, dass alle in der URL enthaltenen Variablen existieren
	url, ok := server["url"].(string)
	if !ok {
		return nil // überspringt, falls keine URL vorhanden
	}

	variables, _ := server["variables"].(map[string]interface{})
	varRegex := regexp.MustCompile(`\{([^}]+)\}`)
	for _, match := range varRegex.FindAllStringSubmatch(url, -1) {
		name := match[1]
		if _, ok := variables[name]; !ok {
			return fmt.Errorf("missing variable definition for %s in server url", name)
		}
	}
	return nil
}
