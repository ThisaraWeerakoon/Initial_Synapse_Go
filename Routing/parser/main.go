package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// --- Input Struct Definitions (from your request) ---

type API struct {
	Context     string // Used for basePath calculation
	Name        string // Ignored as requested
	Version     string // Used for info.version and potentially basePath
	VersionType string // Used for basePath calculation ('url' or 'context')
	Resources   []Resource
	Position    Position // Ignored as requested
}

type Resource struct {
	Methods       []string        // e.g., GET, POST, PUT
	URITemplate   URITemplateInfo
	InSequence    Sequence // Ignored as requested
	FaultSequence Sequence // Ignored as requested
}

type URITemplateInfo struct {
	FullTemplate    string            // e.g., /hospital/{category}/path/to/{date}?name={varName}&age={varAge}
	PathTemplate    string            // e.g., /hospital/{category}/path/to/{date}
	PathParameters  []string          // e.g., category, date
	QueryParameters map[string]string // e.g., {name:varName, age:varAge} - maps query param name to template variable name
}

// --- Placeholder types for ignored fields ---
type Position struct{}
type Sequence struct{}


// This function calculates the base path based on the API context and versioning type.
func (api *API) calculateBasePath() string {
	basePath := api.Context

	// Remove trailing slash from context if present
	if len(basePath) > 1 && strings.HasSuffix(basePath, "/") {
		basePath = basePath[:len(basePath)-1]
	}

    // Ensure basePath starts with a slash if not empty
    if basePath != "" && !strings.HasPrefix(basePath, "/") {
        basePath = "/" + basePath
    }


	// Handle versioning based on versionType
	if api.Version != "" && api.VersionType != "" {
		switch api.VersionType {
		case "url":
			// For URL type, add version as a path segment
            // Ensure no double slash if basePath is "/"
            if basePath == "/" {
                 basePath = "/" + api.Version
            } else if basePath == "" {
                 basePath = "/" + api.Version
            } else {
                 basePath = basePath + "/" + api.Version
            }
		case "context":
			// For context type, replace {version} placeholder if it exists
			versionPattern := "{version}"
			basePath = strings.Replace(basePath, versionPattern, api.Version, 1)
		}
	}
	return basePath
}


// GenerateOpenAPISpec creates an OpenAPI 3.0.x specification document
// as a map[string]interface{} based on the input API definition.
func (api *API) GenerateOpenAPISpec(hostname string, port int) (map[string]interface{}, error) {
	// --- 1. Basic OpenAPI Structure ---
	spec := make(map[string]interface{})
	spec["openapi"] = "3.0.3" // Specify OpenAPI version

	// --- 2. Info Object ---
	info := make(map[string]interface{})
    // Use API Name or a default if context is empty.
    title := api.Name
    if title == "" {
        title = "API Documentation"
    }
    info["title"] = title
	info["version"] = api.Version 
	spec["info"] = info

	// --- 3. Servers Object ---
	basePath := api.calculateBasePath()
	serverURL := fmt.Sprintf("http://%s:%d%s", hostname, port, basePath)
    // Validate the URL 
    _, err := url.Parse(serverURL)
    if err != nil {
        return nil, fmt.Errorf("generated server URL is invalid '%s': %w", serverURL, err)
    }

	servers := []map[string]interface{}{
		{
			"url": serverURL,
		},
	}
	spec["servers"] = servers

	// --- 4. Paths Object ---
	paths := make(map[string]interface{})

	for _, resource := range api.Resources {
		pathTemplate := resource.URITemplate.PathTemplate
		// OpenAPI paths MUST start with a '/'
		if !strings.HasPrefix(pathTemplate, "/") {
			pathTemplate = "/" + pathTemplate
		}

		// Get or create the Path Item Object for this path
		pathItem, exists := paths[pathTemplate]
		if !exists {
			pathItem = make(map[string]interface{})
			paths[pathTemplate] = pathItem
		}
		pathItemMap := pathItem.(map[string]interface{}) // Type assertion

		// --- 5. Parameters (Path and Query) ---
		parameters := make([]interface{}, 0)

		// Path Parameters
		for _, paramName := range resource.URITemplate.PathParameters {
			parameters = append(parameters, map[string]interface{}{
				"name":        paramName,
				"in":          "path",
				"required":    true, // Path parameters are always required
				"description": fmt.Sprintf("Path parameter: %s", paramName), // Optional description
				"schema": map[string]interface{}{
					"type": "string", // Defaulting to string.
				},
			})
		}

		// Query Parameters
		for queryParamName := range resource.URITemplate.QueryParameters {
			parameters = append(parameters, map[string]interface{}{
				"name":        queryParamName,
				"in":          "query",
				"required":    true, 
				"description": fmt.Sprintf("Query parameter: %s", queryParamName), // Optional description
				"schema": map[string]interface{}{
					"type": "string", // Defaulting to string.
				},
			})
		}


		// --- 6. Operations (GET, POST, PUT, etc.) ---
		for _, method := range resource.Methods {
			httpMethod := strings.ToLower(method) // OpenAPI methods are lowercase (get, post, put)

			operation := make(map[string]interface{})
			// Add a basic summary 
			operation["summary"] = fmt.Sprintf("%s operation for %s", strings.ToUpper(httpMethod), pathTemplate)

            if len(parameters) > 0 {
                 operation["parameters"] = parameters
            }

			// --- 7. Request Body (Needed for POST, PUT, PATCH, etc.) ---
			if httpMethod == "post" || httpMethod == "put" || httpMethod == "patch" {
				requestBody := make(map[string]interface{})
				requestBody["description"] = "Request body payload" // Placeholder description
				requestBody["required"] = true // Typically true for POST/PUT

				content := make(map[string]interface{})
	
				content["application/json"] = map[string]interface{}{
					"schema": map[string]interface{}{
						"type": "object",
                        "properties": map[string]interface{}{
                            "message": map[string]interface{}{
                                "type": "string",
                                "example": "Placeholder - Define actual schema based on API needs",
                            },
                        },
					},
				}
				requestBody["content"] = content
				operation["requestBody"] = requestBody
			}

			// --- 8. Responses ---
			// Adding minimal default responses.
			responses := make(map[string]interface{})
			responses["200"] = map[string]interface{}{ // Successful response
				"description": "OK", // Basic description
                "content": map[string]interface{}{
					"application/json": map[string]interface{}{
						"schema": map[string]interface{}{
							"type": "object",
                            "properties": map[string]interface{}{
                                "status": map[string]interface{}{
                                    "type": "string",
                                    "example": "Success",
                                },
                            },
						},
					},
                },
			}

			responses["default"] = map[string]interface{}{ // Default response for errors
				"description": "Unexpected error",
                 "content": map[string]interface{}{
					"application/json": map[string]interface{}{
						"schema": map[string]interface{}{
                            "type": "object",
                            "properties": map[string]interface{}{
                                "error": map[string]interface{}{
                                    "type": "string",
                                    "example": "Error details",
                                },
                            },
						},
					},
                },
			}
			operation["responses"] = responses

			// Assign the operation to the correct method in the path item
			pathItemMap[httpMethod] = operation
		} // End methods loop
	} // End resources loop

	spec["paths"] = paths

	return spec, nil
}

// --- Main function for Example Usage ---

func main() {
	// Example API definition based on your description
	exampleAPI := API{
		Context:     "/api", // Base part of the path
		Version:     "2.0",  // Version string
		VersionType: "url",  // Version added as part of the URL path
		Resources: []Resource{
			{
				Methods: []string{"GET", "POST", "PUT"}, // Methods for this resource path
				URITemplate: URITemplateInfo{
					// Note: Query parameters are part of FullTemplate but handled separately in spec generation
					FullTemplate:    "/hospital/{category}/path/to/{date}?name={varName}&age={varAge}",
					PathTemplate:    "/hospital/{category}/path/to/{date}", // Path part for OpenAPI path key
					PathParameters:  []string{"category", "date"},         // Extracted from {} in PathTemplate
					QueryParameters: map[string]string{"name": "varName", "age": "varAge"}, // Query params ?key=value
				},
			},
            // Add more resources here if the API has multiple paths
            // {
            //    Methods: []string{"GET"},
            //    URITemplate: URITemplateInfo {
            //        PathTemplate: "/health",
            //        PathParameters: []string{},
            //        QueryParameters: map[string]string{},
            //    },
            // },
		},
	}

	hostname := "localhost"
	port := 8290

	// Generate the OpenAPI spec
	swaggerSpecMap, err := exampleAPI.GenerateOpenAPISpec(hostname, port)
	if err != nil {
		fmt.Println("Error generating Swagger spec:", err)
		return
	}

	// Convert the map to JSON for output/saving
	// Using MarshalIndent for readability
	jsonData, err := json.MarshalIndent(swaggerSpecMap, "", "  ") // Use "  " for indentation
	if err != nil {
		fmt.Println("Error marshalling spec to JSON:", err)
		return
	}

	// Print the resulting JSON specification
	fmt.Println(string(jsonData))

    // You could also write this jsonData to a file (e.g., swagger.json or openapi.json)
    // ioutil.WriteFile("openapi.json", jsonData, 0644)
}