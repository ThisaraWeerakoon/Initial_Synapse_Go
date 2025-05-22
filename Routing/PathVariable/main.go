package main

import (
	"regexp"
)

func main(){
	// Example route pattern
	path := "/product/{id}/{iddddd}/orders/{orderId}"
	re := regexp.MustCompile(`\{([^}]+)\}`)
	matches := re.FindAllStringSubmatch(path, -1)
	variables := make([]string, 0, len(matches))
	for _, match := range matches {
		if len(match) > 1 {
			variables = append(variables, match[1])
		}
	}
	// Print the extracted variables
	for _, variable := range variables {
		println(variable)
	}

}

