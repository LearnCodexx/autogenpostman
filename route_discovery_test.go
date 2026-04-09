package autogenpostman

import (
	"testing"
)

func TestRouteScanner_ScanFiberRoutes(t *testing.T) {
	scanner := NewRouteScanner(".")
	
	testLines := []string{
		`app.Post("/users/signin", ua.SignIn)`,
		`app.Get("/users/profile", ua.GetProfile)`,
		`users.Post("/signup", ua.SignUp)`,
		`userGet.Post("/all", ua.GetAllUser)`,
	}
	
	scanner.scanFiberRoutes(testLines, "test.go")
	
	routes := scanner.GetDiscoveredRoutes()
	
	// Should find at least the direct app.* routes
	if len(routes) < 2 {
		t.Errorf("Expected at least 2 routes, got %d", len(routes))
	}
	
	// Check first route
	found := false
	for _, route := range routes {
		if route.Method == "POST" && route.Path == "/users/signin" && route.Handler == "ua.SignIn" {
			found = true
			break
		}
	}
	
	if !found {
		t.Error("Expected to find POST /users/signin -> ua.SignIn route")
	}
}

func TestRouteScanner_ScanModels(t *testing.T) {
	scanner := NewRouteScanner(".")
	
	testLines := []string{
		`package dto`,
		``,
		`type UserDto struct {`,
		`	Name     string ` + "`" + `json:"name" validate:"required"` + "`",
		`	Email    string ` + "`" + `json:"email" validate:"required,email"` + "`",  
		`	Password string ` + "`" + `json:"password" validate:"required,min=6"` + "`",
		`}`,
		``,
		`type LoginRequest struct {`,
		`	Username string ` + "`" + `json:"username"` + "`",
		`	Password string ` + "`" + `json:"password"` + "`",
		`}`,
	}
	
	scanner.scanModels(testLines, "dto/user_dto.go")
	
	models := scanner.GetDiscoveredModels()
	
	if len(models) != 2 {
		t.Errorf("Expected 2 models, got %d", len(models))
	}
	
	// Check UserDto model
	found := false
	for _, model := range models {
		if model.Name == "UserDto" {
			found = true
			if len(model.Fields) != 3 {
				t.Errorf("Expected UserDto to have 3 fields, got %d", len(model.Fields))
			}
			
			// Check if name field exists and has correct properties
			nameFieldFound := false
			for _, field := range model.Fields {
				if field.Name == "Name" {
					nameFieldFound = true
					if field.JsonTag != "name" {
						t.Errorf("Expected Name field JsonTag to be 'name', got '%s'", field.JsonTag)
					}
					if !field.Required {
						t.Error("Expected Name field to be required")
					}
				}
			}
			
			if !nameFieldFound {
				t.Error("Expected to find Name field in UserDto")
			}
			break
		}
	}
	
	if !found {
		t.Error("Expected to find UserDto model")
	}
}

func TestRouteScanner_GenerateSwaggerParts(t *testing.T) {
	scanner := NewRouteScanner(".")
	
	// Add a test route
	route := DiscoveredRoute{
		Method:      "POST",
		Path:        "/users/signin",
		Handler:     "ua.SignIn", 
		Summary:     "User sign in",
		Description: "Authenticate user",
		Tags:        []string{"Users"},
		Auth:        false,
	}
	
	scanner.routes = []DiscoveredRoute{route}
	
	// Add a test model
	model := DiscoveredModel{
		Name: "UserDto",
		Fields: []FieldDefinition{
			{
				Name:     "Email",
				Type:     "string",
				JsonTag:  "email", 
				Required: true,
			},
		},
	}
	
	scanner.models = []DiscoveredModel{model}
	
	// Test basic functionality
	routes := scanner.GetDiscoveredRoutes()
	models := scanner.GetDiscoveredModels()
	
	if len(routes) != 1 {
		t.Errorf("Expected 1 route, got %d", len(routes))
	}
	
	if len(models) != 1 {
		t.Errorf("Expected 1 model, got %d", len(models))
	}
	
	if routes[0].Path != "/users/signin" {
		t.Errorf("Expected route path '/users/signin', got '%s'", routes[0].Path)
	}
}

func TestGenerateOperationId(t *testing.T) {
	g := &Generator{}
	
	route := DiscoveredRoute{
		Method: "POST",
		Path:   "/users/signin",
	}
	
	operationId := g.generateOperationId(route)
	expected := "postUsers"
	
	if operationId != expected {
		t.Errorf("Expected operation ID '%s', got '%s'", expected, operationId)
	}
}

func TestConvertGoTypeToSwagger(t *testing.T) {
	g := &Generator{}
	
	testCases := []struct {
		goType   string
		expected string
	}{
		{"string", "string"},
		{"int", "integer"},  
		{"int64", "integer"},
		{"bool", "boolean"},
		{"float64", "number"},
		{"[]string", "array"},
		{"Unknown", "string"},
	}
	
	for _, tc := range testCases {
		result := g.convertGoTypeToSwagger(tc.goType)
		if result != tc.expected {
			t.Errorf("Expected %s -> %s, got %s", tc.goType, tc.expected, result)
		}
	}
}

// Benchmark tests
func BenchmarkRouteScanning(b *testing.B) {
	scanner := NewRouteScanner(".")
	
	testLines := []string{
		`app.Post("/users/signin", ua.SignIn)`,
		`app.Get("/users/profile", ua.GetProfile)`,
		`app.Put("/users/:id", ua.UpdateUser)`,
		`app.Delete("/users/:id", ua.DeleteUser)`,
		`users.Post("/signup", ua.SignUp)`,
	}
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		scanner.routes = []DiscoveredRoute{} // Reset
		scanner.scanFiberRoutes(testLines, "test.go")
	}
}