package autogenpostman

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// RouteScanner automatically discovers routes from Go source code
type RouteScanner struct {
	projectPath string
	routes      []DiscoveredRoute
	models      []DiscoveredModel
}

// DiscoveredRoute represents a route found in source code
type DiscoveredRoute struct {
	Method      string            `json:"method"`
	Path        string            `json:"path"`
	Handler     string            `json:"handler"`
	File        string            `json:"file"`
	LineNumber  int               `json:"lineNumber"`
	Summary     string            `json:"summary"`
	Description string            `json:"description"`
	Tags        []string          `json:"tags"`
	Middleware  []string          `json:"middleware"`
	RequestBody *string           `json:"requestBody,omitempty"`  
	Auth        bool              `json:"auth"`
}

// DiscoveredModel represents a struct/model found in source code
type DiscoveredModel struct {
	Name       string            `json:"name"`
	Fields     []FieldDefinition `json:"fields"`
	File       string            `json:"file"`
	Package    string            `json:"package"`
	JsonTags   map[string]string `json:"jsonTags"`
	Validators map[string]string `json:"validators"`
}

// FieldDefinition represents a struct field
type FieldDefinition struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	JsonTag   string `json:"jsonTag"`
	Validate  string `json:"validate"`
	Required  bool   `json:"required"`
}

// NewRouteScanner creates a new route scanner
func NewRouteScanner(projectPath string) *RouteScanner {
	return &RouteScanner{
		projectPath: projectPath,
		routes:      []DiscoveredRoute{},
		models:      []DiscoveredModel{},
	}
}

// ScanRoutes discovers routes automatically from Go source files
func (rs *RouteScanner) ScanRoutes() error {
	fmt.Printf("🔍 Auto-scanning routes in: %s\n", rs.projectPath)
	
	err := filepath.Walk(rs.projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		// Skip directories and non-Go files
		if info.IsDir() {
			// Skip vendor and common ignored directories
			name := info.Name()
			if name == "vendor" || name == "node_modules" || name == ".git" || 
			   name == "docs" || name == "test-project" || name == "example" {
				return filepath.SkipDir
			}
			return nil
		}
		
		// Only process .go files (skip tests)
		if !strings.HasSuffix(info.Name(), ".go") || strings.HasSuffix(info.Name(), "_test.go") {
			return nil
		}
		
		return rs.scanFile(path)
	})
	
	fmt.Printf("✅ Scan complete: %d routes, %d models discovered\n", len(rs.routes), len(rs.models))
	return err
}

// scanFile scans a single Go file for routes and models
func (rs *RouteScanner) scanFile(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	
	lines := strings.Split(string(content), "\n")
	
	// Scan for different route patterns
	rs.scanFiberRoutes(lines, filePath)
	rs.scanGinRoutes(lines, filePath)
	rs.scanEchoRoutes(lines, filePath)
	rs.scanMuxRoutes(lines, filePath)
	rs.scanGenericRoutes(lines, filePath)
	
	// Scan for models/DTOs
	rs.scanModels(lines, filePath)
	
	return nil
}

// scanFiberRoutes scans for Fiber framework routes
func (rs *RouteScanner) scanFiberRoutes(lines []string, filePath string) {
	patterns := map[string]string{
		// Direct app calls
		`app\.Get\s*\(\s*"([^"]+)"\s*,\s*([^)]+)\)`:    "GET",
		`app\.Post\s*\(\s*"([^"]+)"\s*,\s*([^)]+)\)`:   "POST",
		`app\.Put\s*\(\s*"([^"]+)"\s*,\s*([^)]+)\)`:    "PUT",
		`app\.Delete\s*\(\s*"([^"]+)"\s*,\s*([^)]+)\)`: "DELETE",
		`app\.Patch\s*\(\s*"([^"]+)"\s*,\s*([^)]+)\)`:  "PATCH",
		
		// Group calls (e.g., users.Post)
		`\w+\.Get\s*\(\s*"([^"]+)"\s*,\s*([^)]+)\)`:    "GET",
		`\w+\.Post\s*\(\s*"([^"]+)"\s*,\s*([^)]+)\)`:   "POST",
		`\w+\.Put\s*\(\s*"([^"]+)"\s*,\s*([^)]+)\)`:    "PUT",
		`\w+\.Delete\s*\(\s*"([^"]+)"\s*,\s*([^)]+)\)`: "DELETE",
		`\w+\.Patch\s*\(\s*"([^"]+)"\s*,\s*([^)]+)\)`:  "PATCH",
	}
	
	rs.scanWithPatterns(patterns, lines, filePath, "Fiber")
}

// scanGinRoutes scans for Gin framework routes
func (rs *RouteScanner) scanGinRoutes(lines []string, filePath string) {
	patterns := map[string]string{
		`r\.GET\s*\(\s*"([^"]+)"\s*,\s*([^)]+)\)`:       "GET",
		`r\.POST\s*\(\s*"([^"]+)"\s*,\s*([^)]+)\)`:      "POST",
		`router\.GET\s*\(\s*"([^"]+)"\s*,\s*([^)]+)\)`:  "GET",
		`router\.POST\s*\(\s*"([^"]+)"\s*,\s*([^)]+)\)`: "POST",
		`engine\.GET\s*\(\s*"([^"]+)"\s*,\s*([^)]+)\)`:  "GET",
		`engine\.POST\s*\(\s*"([^"]+)"\s*,\s*([^)]+)\)`: "POST",
	}
	
	rs.scanWithPatterns(patterns, lines, filePath, "Gin")
}

// scanEchoRoutes scans for Echo framework routes
func (rs *RouteScanner) scanEchoRoutes(lines []string, filePath string) {
	patterns := map[string]string{
		`e\.GET\s*\(\s*"([^"]+)"\s*,\s*([^)]+)\)`:   "GET",
		`e\.POST\s*\(\s*"([^"]+)"\s*,\s*([^)]+)\)`:  "POST",
		`echo\.GET\s*\(\s*"([^"]+)"\s*,\s*([^)]+)\)`: "GET",
		`echo\.POST\s*\(\s*"([^"]+)"\s*,\s*([^)]+)\)`: "POST",
	}
	
	rs.scanWithPatterns(patterns, lines, filePath, "Echo")
}

// scanMuxRoutes scans for Gorilla Mux routes  
func (rs *RouteScanner) scanMuxRoutes(lines []string, filePath string) {
	patterns := map[string]string{
		`HandleFunc\s*\(\s*"([^"]+)"\s*,\s*([^)]+)\)\.Methods\s*\(\s*"GET"\s*\)`:    "GET",
		`HandleFunc\s*\(\s*"([^"]+)"\s*,\s*([^)]+)\)\.Methods\s*\(\s*"POST"\s*\)`:   "POST",
		`Handle\s*\(\s*"([^"]+)"\s*,\s*([^)]+)\)\.Methods\s*\(\s*"GET"\s*\)`:        "GET",
		`Handle\s*\(\s*"([^"]+)"\s*,\s*([^)]+)\)\.Methods\s*\(\s*"POST"\s*\)`:       "POST",
	}
	
	rs.scanWithPatterns(patterns, lines, filePath, "Mux")
}

// scanGenericRoutes scans for generic route registration patterns
func (rs *RouteScanner) scanGenericRoutes(lines []string, filePath string) {
	// Look for RegisterRoutes patterns and grouped routes
	for i, line := range lines {
		// Group definition pattern: users := app.Group("/users")
		groupDefPattern := regexp.MustCompile(`(\w+)\s*:=\s*\w+\.Group\s*\(\s*"([^"]+)"\s*\)`)
		matches := groupDefPattern.FindStringSubmatch(line)
		if len(matches) >= 3 {
			groupVar := matches[1]
			basePath := matches[2]
			
			// Look ahead for routes in this group
			rs.scanGroupRoutes(lines, i+1, groupVar, basePath, filePath)
		}
	}
}

// scanGroupRoutes scans for routes within a route group
func (rs *RouteScanner) scanGroupRoutes(lines []string, startLine int, groupVar, basePath, filePath string) {
	for i := startLine; i < len(lines) && i < startLine+50; i++ {
		line := lines[i]
		
		// Look for group.Method patterns
		pattern := regexp.MustCompile(groupVar + `\.(\w+)\s*\(\s*"([^"]+)"\s*,\s*([^)]+)\)`)
		matches := pattern.FindStringSubmatch(line)
		
		if len(matches) >= 4 {
			method := strings.ToUpper(matches[1])
			path := matches[2]
			handler := rs.cleanHandler(matches[3])
			
			// Skip non-HTTP methods
			if !rs.isValidHTTPMethod(method) {
				continue
			}
			
			fullPath := basePath + path
			
			route := DiscoveredRoute{
				Method:      method,
				Path:        fullPath,
				Handler:     handler,
				File:        filePath,
				LineNumber:  i + 1,
				Summary:     rs.generateSummary(method, fullPath),
				Description: rs.generateDescription(method, fullPath, handler),
				Tags:        rs.extractTags(fullPath),
				Auth:        rs.detectAuth(lines, i),
			}
			
			rs.routes = append(rs.routes, route)
			fmt.Printf("   📍 %s %s -> %s (%s:%d)\n", method, fullPath, handler, filepath.Base(filePath), i+1)
		}
	}
}

// scanWithPatterns scans lines using provided regex patterns
func (rs *RouteScanner) scanWithPatterns(patterns map[string]string, lines []string, filePath, framework string) {
	for i, line := range lines {
		for pattern, method := range patterns {
			re := regexp.MustCompile(pattern)
			matches := re.FindStringSubmatch(line)
			
			if len(matches) >= 3 {
				path := matches[1]
				handler := rs.cleanHandler(matches[2])
				
				route := DiscoveredRoute{
					Method:      method,
					Path:        path,
					Handler:     handler,
					File:        filePath,
					LineNumber:  i + 1,
					Summary:     rs.generateSummary(method, path),
					Description: rs.generateDescription(method, path, handler),
					Tags:        rs.extractTags(path),
					Auth:        rs.detectAuth(lines, i),
				}
				
				rs.routes = append(rs.routes, route)
				fmt.Printf("   📍 %s %s -> %s (%s:%d) [%s]\n", method, path, handler, filepath.Base(filePath), i+1, framework)
			}
		}
	}
}

// scanModels scans for struct definitions that could be API models
func (rs *RouteScanner) scanModels(lines []string, filePath string) {
	// Focus on files likely to contain models
	fileName := strings.ToLower(filepath.Base(filePath))
	if !strings.Contains(fileName, "dto") && !strings.Contains(fileName, "model") && 
	   !strings.Contains(fileName, "domain") && !strings.Contains(fileName, "type") &&
	   !strings.Contains(fileName, "request") && !strings.Contains(fileName, "response") {
		return
	}
	
	structPattern := regexp.MustCompile(`type\s+(\w+)\s+struct\s*\{?`)
	
	for i, line := range lines {
		matches := structPattern.FindStringSubmatch(line)
		if len(matches) >= 2 {
			structName := matches[1]
			
			// Skip if it doesn't look like a public API type
			if !rs.isPublicStruct(structName) {
				continue
			}
			
			fields := rs.extractStructFields(lines, i+1)
			
			model := DiscoveredModel{
				Name:     structName,
				Fields:   fields,
				File:     filePath,
				Package:  rs.extractPackageName(lines),
				JsonTags: make(map[string]string),
			}
			
			rs.models = append(rs.models, model)
			fmt.Printf("   📦 Model: %s (%d fields) in %s\n", structName, len(fields), filepath.Base(filePath))
		}
	}
}

// Helper methods

func (rs *RouteScanner) cleanHandler(handler string) string {
	// Remove extra whitespace and clean up handler names
	handler = strings.TrimSpace(handler)
	// Remove middleware wrapping
	if strings.Contains(handler, ",") {
		parts := strings.Split(handler, ",")
		handler = strings.TrimSpace(parts[len(parts)-1])
	}
	return handler
}

func (rs *RouteScanner) isValidHTTPMethod(method string) bool {
	validMethods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD"}
	for _, valid := range validMethods {
		if method == valid {
			return true
		}
	}
	return false
}

func (rs *RouteScanner) generateSummary(method, path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) > 0 {
		resource := strings.Title(parts[0])
		return fmt.Sprintf("%s %s", strings.Title(strings.ToLower(method)), resource)
	}
	return fmt.Sprintf("%s operation", strings.Title(strings.ToLower(method)))
}

func (rs *RouteScanner) generateDescription(method, path, handler string) string {
	return fmt.Sprintf("%s operation for %s endpoint (handler: %s)", method, path, handler)
}

func (rs *RouteScanner) extractTags(path string) []string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) > 0 {
		return []string{strings.Title(parts[0])}
	}
	return []string{"API"}
}

func (rs *RouteScanner) detectAuth(lines []string, currentLine int) bool {
	// Look around current line for auth indicators
	start := max(0, currentLine-3)
	end := min(len(lines)-1, currentLine+3)
	
	for i := start; i <= end; i++ {
		line := strings.ToLower(lines[i])
		if strings.Contains(line, "auth") || strings.Contains(line, "jwt") || 
		   strings.Contains(line, "protected") || strings.Contains(line, "middleware") {
			return true
		}
	}
	return false
}

func (rs *RouteScanner) isPublicStruct(name string) bool {
	// Must start with capital letter and look like an API type
	if len(name) == 0 || name[0] < 'A' || name[0] > 'Z' {
		return false
	}
	
	// Common API struct suffixes/patterns
	apiPatterns := []string{"Dto", "DTO", "Request", "Response", "Model", "Data", "Info", "Config"}
	lowerName := strings.ToLower(name)
	
	for _, pattern := range apiPatterns {
		if strings.Contains(lowerName, strings.ToLower(pattern)) {
			return true
		}
	}
	
	return false
}

func (rs *RouteScanner) extractStructFields(lines []string, startLine int) []FieldDefinition {
	var fields []FieldDefinition
	
	for i := startLine; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		
		// End of struct
		if line == "}" {
			break
		}
		
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}
		
		// Parse field definition
		field := rs.parseFieldLine(line)
		if field != nil {
			fields = append(fields, *field)
		}
	}
	
	return fields
}

func (rs *RouteScanner) parseFieldLine(line string) *FieldDefinition {
	// Basic field parsing - can be enhanced
	parts := strings.Fields(line)
	if len(parts) < 2 {
		return nil
	}
	
	fieldName := parts[0]
	fieldType := parts[1]
	
	// Skip embedded types and methods
	if !rs.isValidFieldName(fieldName) {
		return nil
	}
	
	field := &FieldDefinition{
		Name: fieldName,
		Type: fieldType,
	}
	
	// Extract tags if present
	if len(parts) > 2 {
		tagPart := strings.Join(parts[2:], " ")
		field.JsonTag = rs.extractTag(tagPart, "json")
		field.Validate = rs.extractTag(tagPart, "validate")
		field.Required = strings.Contains(field.Validate, "required")
	}
	
	return field
}

func (rs *RouteScanner) isValidFieldName(name string) bool {
	if len(name) == 0 {
		return false
	}
	// Must start with capital letter
	return name[0] >= 'A' && name[0] <= 'Z'
}

func (rs *RouteScanner) extractTag(tagStr, tagName string) string {
	pattern := fmt.Sprintf(`%s:"([^"]*)"`, tagName)
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(tagStr)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func (rs *RouteScanner) extractPackageName(lines []string) string {
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "package ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				return parts[1]
			}
		}
	}
	return "main"
}

// GetDiscoveredRoutes returns all discovered routes
func (rs *RouteScanner) GetDiscoveredRoutes() []DiscoveredRoute {
	return rs.routes
}

// GetDiscoveredModels returns all discovered models
func (rs *RouteScanner) GetDiscoveredModels() []DiscoveredModel {
	return rs.models
}

// PrintSummary prints a summary of discovered routes and models
func (rs *RouteScanner) PrintSummary() {
	fmt.Printf("\n🎯 ROUTE DISCOVERY SUMMARY\n")
	fmt.Printf("==========================\n")
	fmt.Printf("📍 Routes discovered: %d\n", len(rs.routes))
	
	if len(rs.routes) > 0 {
		fmt.Printf("\nRoutes by method:\n")
		methodCount := make(map[string]int)
		for _, route := range rs.routes {
			methodCount[route.Method]++
		}
		for method, count := range methodCount {
			fmt.Printf("  %s: %d routes\n", method, count)
		}
	}
	
	fmt.Printf("\n📦 Models discovered: %d\n", len(rs.models))
	if len(rs.models) > 0 {
		for _, model := range rs.models {
			fmt.Printf("  %s (%d fields)\n", model.Name, len(model.Fields))
		}
	}
}

// Utility functions
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}