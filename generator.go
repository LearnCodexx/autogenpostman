package autogenpostman

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

// Config is the high-level generator configuration.
type Config struct {
	WorkingDir       string
	SwaggerInputPath string
	OutputPath       string
	CollectionName   string
	Pretty           bool
	Swag             SwagConfig
	Postman          PostmanConfig
}

// AutoConfig is a high-level configuration for one-call generation.
// It tries to generate swagger with swag first, then falls back to existing OpenAPI files.
type AutoConfig struct {
	WorkingDir        string
	OutputPath        string
	CollectionName    string
	Pretty            bool
	MainFile          string
	SwagOutputDir     string
	SwaggerInputPath  string
	SwaggerCandidates []string
	Postman           PostmanConfig
}

// SwagConfig controls whether swagger docs should be generated via swag.
type SwagConfig struct {
	Enabled         bool
	MainFile        string
	OutputDir       string
	ParseDependency bool
	ParseInternal   bool
	InstanceName    string
	UseGoRun        bool
}

// PostmanConfig controls how openapi-to-postmanv2 is called.
type PostmanConfig struct {
	UseLocalCLI bool
	CLIPath     string
	Options     map[string]string
}

// Generator executes the generation workflow.
type Generator struct {
	runner    commandRunner
	lookPath  func(file string) (string, error)
	fileExist func(path string) bool
}

type commandRunner interface {
	Run(ctx context.Context, dir string, name string, args ...string) error
}

type execRunner struct{}

func (execRunner) Run(ctx context.Context, dir string, name string, args ...string) error {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir

	out, err := cmd.CombinedOutput()
	if err != nil {
		trimmed := strings.TrimSpace(string(out))
		if trimmed == "" {
			return fmt.Errorf("run %q failed: %w", name, err)
		}
		return fmt.Errorf("run %q failed: %w: %s", name, err, trimmed)
	}
	return nil
}

// New returns a generator with the default process runner.
func New() *Generator {
	return &Generator{
		runner:   execRunner{},
		lookPath: exec.LookPath,
		fileExist: func(path string) bool {
			_, err := os.Stat(path)
			return err == nil
		},
	}
}

// Generate runs swag (optional) and converts OpenAPI/Swagger to a Postman collection.
func (g *Generator) Generate(ctx context.Context, cfg Config) error {
	if g == nil {
		return errors.New("generator is nil")
	}
	if g.runner == nil {
		g.runner = execRunner{}
	}
	if g.lookPath == nil {
		g.lookPath = exec.LookPath
	}
	if g.fileExist == nil {
		g.fileExist = func(path string) bool {
			_, err := os.Stat(path)
			return err == nil
		}
	}

	if err := cfg.normalize(); err != nil {
		return err
	}

	if cfg.Swag.Enabled {
		if err := g.runSwag(ctx, cfg); err != nil {
			return err
		}
	}

	if err := g.runConvert(ctx, cfg); err != nil {
		return err
	}

	if cfg.CollectionName != "" {
		if err := renameCollection(cfg.OutputPath, cfg.CollectionName); err != nil {
			return err
		}
	}

	return nil
}

// Generate is a package-level convenience API.
func Generate(ctx context.Context, cfg Config) error {
	return New().Generate(ctx, cfg)
}

// GenerateAuto is a high-level package API intended for cross-application reuse.
// Behavior:
// 1) If SwaggerInputPath is provided, it is used directly (no swag run).
// 2) Otherwise, if swag exists in PATH, swagger is generated via swag.
// 3) If swag is unavailable, generator falls back to existing OpenAPI candidate files.
func (g *Generator) GenerateAuto(ctx context.Context, cfg AutoConfig) error {
	fmt.Printf("🚀 Starting autogenpostman generation...\n")
	
	if g == nil {
		return errors.New("generator is nil")
	}
	if g.runner == nil || g.lookPath == nil || g.fileExist == nil {
		*g = *New()
	}

	workingDir, err := resolveWorkingDir(cfg.WorkingDir)
	if err != nil {
		return err
	}
	fmt.Printf("📁 Working directory: %s\n", workingDir)

	if cfg.OutputPath == "" {
		cfg.OutputPath = filepath.Join("docs", "postman_collection.json")
	}
	fmt.Printf("📄 Output file: %s\n", cfg.OutputPath)
	
	if cfg.MainFile == "" {
		cfg.MainFile = g.findMainFile(workingDir)
		if cfg.MainFile == "" {
			// Try to auto-scaffold if no main.go found
			if err := g.tryAutoScaffold(workingDir); err == nil {
				fmt.Printf("🔧 Auto-created basic project structure\n")
				cfg.MainFile = "main.go"
			} else {
				return fmt.Errorf("cannot find main.go file in common locations\n\n💡 Solutions:\n1. Run from project root directory\n2. Create main.go with basic API structure\n3. Use explicit path: --main-file=path/to/main.go\n4. Auto-scaffold: go run github.com/learncodexx/autogenpostman/cmd/postman@latest")
			}
		}
	} else {
		fmt.Printf("📍 Using specified main file: %s\n", cfg.MainFile)
	}

	if cfg.SwagOutputDir == "" {
		// Use docs as default, more standard
		cfg.SwagOutputDir = filepath.Join(workingDir, "docs")
		// Create the directory if it doesn't exist
		if err := os.MkdirAll(cfg.SwagOutputDir, 0755); err != nil {
			return fmt.Errorf("create swag output dir: %w", err)
		}
	}
	fmt.Printf("📂 Swagger output dir: %s\n", cfg.SwagOutputDir)

	if len(cfg.SwaggerCandidates) == 0 {
		cfg.SwaggerCandidates = g.getSwaggerCandidates(workingDir)
	}
	lowLevelCfg := Config{
		WorkingDir:     workingDir,
		OutputPath:     cfg.OutputPath,
		CollectionName: cfg.CollectionName,
		Pretty:         cfg.Pretty,
		Postman:        cfg.Postman,
	}

	if cfg.SwaggerInputPath != "" {
		fmt.Printf("📄 Using provided OpenAPI file: %s\n", cfg.SwaggerInputPath)
		lowLevelCfg.SwaggerInputPath = cfg.SwaggerInputPath
		return g.Generate(ctx, lowLevelCfg)
	}

	fmt.Printf("🔄 Attempting auto-generation with swag...\n")
	lowLevelCfg.Swag = SwagConfig{
		Enabled:         true,
		MainFile:        cfg.MainFile,
		OutputDir:       cfg.SwagOutputDir,
		ParseDependency: true,
		ParseInternal:   true,
		UseGoRun:        true,
	}
	swagErr := g.Generate(ctx, lowLevelCfg)
	if swagErr == nil {
		fmt.Printf("🎉 Generation completed successfully!\n")
		return nil
	}

	fmt.Printf("⚠️  Swagger generation failed, trying existing OpenAPI files...\n")

	for _, candidate := range cfg.SwaggerCandidates {
		candidatePath := candidate
		if !filepath.IsAbs(candidatePath) {
			candidatePath = filepath.Join(workingDir, candidatePath)
		}
		if g.fileExist(candidatePath) {
			fmt.Printf("ℹ️  Using existing OpenAPI file: %s\n", candidate)
			lowLevelCfg.SwaggerInputPath = candidatePath
			lowLevelCfg.Swag = SwagConfig{}
			return g.Generate(ctx, lowLevelCfg)
		}
	}

	suggestions := []string{
		"1. Add swagger annotations to your main.go:",
		"   // @title Your API Name",
		"   // @version 1.0", 
		"   // @host localhost:8080",
		"   // @BasePath /api/v1",
		"2. Install swag: go install github.com/swaggo/swag/cmd/swag@latest",
		"3. Or create OpenAPI file at: " + cfg.SwaggerCandidates[0],
		"4. Or specify existing file: --swagger-input=path/to/openapi.yaml",
	}

	return fmt.Errorf("auto swagger generation failed (%v) and no OpenAPI file found\n\n💡 Solutions:\n%s", 
		swagErr, strings.Join(suggestions, "\n"))
}

// GenerateAuto is a package-level convenience API for one-call generation.
func GenerateAuto(ctx context.Context, cfg AutoConfig) error {
	return New().GenerateAuto(ctx, cfg)
}

// QuickGenerate is the most convenient one-liner for generating Postman collections
// It handles everything automatically: scaffolding, annotation validation, generation
func QuickGenerate(projectDir string) error {
	return QuickGenerateNamed(projectDir, "API Collection")
}

// QuickGenerateNamed generates Postman collection with custom name
func QuickGenerateNamed(projectDir, collectionName string) error {
	ctx := context.Background()
	cfg := AutoConfig{
		WorkingDir:     projectDir,
		OutputPath:     "docs/postman_collection.json",
		CollectionName: collectionName,
		Pretty:         true,
	}
	
	generator := New()
	err := generator.GenerateAuto(ctx, cfg)
	if err != nil {
		// If generation fails, try to provide helpful guidance
		if strings.Contains(err.Error(), "cannot find main.go") {
			fmt.Printf("🚨 No main.go found. Creating basic structure...\n")
			if scaffoldErr := generator.tryAutoScaffold(projectDir); scaffoldErr == nil {
				fmt.Printf("✨ Created basic main.go. Please add your API endpoints and try again.\n")
				fmt.Printf("🛠️ Next steps:\n1. Add your API routes to main.go\n2. Run: go mod init your-project-name\n3. Run: go get github.com/gin-gonic/gin\n4. Run generation again\n")
			}
		}
		return err
	}
	
	fmt.Printf("✅ Success! Postman collection: docs/postman_collection.json\n")
	return nil
}

func (g *Generator) runSwag(ctx context.Context, cfg Config) error {
	// Validate main file exists
	mainFilePath := filepath.Join(cfg.WorkingDir, cfg.Swag.MainFile)
	if !g.fileExist(mainFilePath) {
		return fmt.Errorf("main file not found: %s (try --main-file flag to specify correct path)", cfg.Swag.MainFile)
	}

	// Validate swagger annotations (lenient mode)
	if err := g.validateMainFileLenient(mainFilePath); err != nil {
		return fmt.Errorf("main file validation failed: %w", err)
	}

	args := []string{"init", "-g", cfg.Swag.MainFile, "-o", cfg.Swag.OutputDir}
	if cfg.Swag.ParseDependency {
		args = append(args, "--parseDependency")
	}
	if cfg.Swag.ParseInternal {
		args = append(args, "--parseInternal")
	}
	if cfg.Swag.InstanceName != "" {
		args = append(args, "--instanceName", cfg.Swag.InstanceName)
	}

	cmdName := "swag"
	cmdArgs := args
	if cfg.Swag.UseGoRun {
		if _, err := g.lookPath("swag"); err != nil {
			fmt.Printf("⚠️  swag not found in PATH, using 'go run github.com/swaggo/swag/cmd/swag@latest'\n")
			cmdName = "go"
			cmdArgs = append([]string{"run", "github.com/swaggo/swag/cmd/swag@latest"}, args...)
		} else {
			fmt.Printf("✅ Using swag from PATH\n")
		}
	}

	fmt.Printf("🔄 Generating swagger from: %s\n", cfg.Swag.MainFile)
	if err := g.runner.Run(ctx, cfg.WorkingDir, cmdName, cmdArgs...); err != nil {
		return fmt.Errorf("generate swagger with swag: %w\n💡 Troubleshooting:\n- Check swagger annotations in %s\n- Ensure @title, @version, @host, @BasePath are present\n- Verify import path: import _ \"yourmod/docs\"", err, cfg.Swag.MainFile)
	}

	// Validate swagger.json was generated
	swaggerPath := filepath.Join(cfg.Swag.OutputDir, "swagger.json")
	if !g.fileExist(swaggerPath) {
		return fmt.Errorf("swagger.json not generated at %s - check swagger annotations", swaggerPath)
	}

	fmt.Printf("✅ Swagger generated: %s\n", swaggerPath)
	return nil
}

func (g *Generator) runConvert(ctx context.Context, cfg Config) error {
	// Validate swagger input exists
	if !g.fileExist(cfg.SwaggerInputPath) {
		return fmt.Errorf("swagger input file not found: %s", cfg.SwaggerInputPath)
	}

	// Create output directory if it doesn't exist
	outputDir := filepath.Dir(cfg.OutputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	args := []string{"--yes", "openapi-to-postmanv2", "-s", cfg.SwaggerInputPath, "-o", cfg.OutputPath}
	if cfg.Pretty {
		args = append(args, "-p")
	}

	opts := buildOptions(cfg.Postman.Options)
	for _, opt := range opts {
		args = append(args, "-O", opt)
	}

	name := "npx"
	if cfg.Postman.UseLocalCLI {
		name = cfg.Postman.CLIPath
		if name == "" {
			name = "openapi2postmanv2"
		}
		args = []string{"-s", cfg.SwaggerInputPath, "-o", cfg.OutputPath}
		if cfg.Pretty {
			args = append(args, "-p")
		}
		for _, opt := range opts {
			args = append(args, "-O", opt)
		}
	}

	// Check if npx/node is available
	if name == "npx" {
		if _, err := g.lookPath("npx"); err != nil {
			return fmt.Errorf("npx not found in PATH: %w\n💡 Install Node.js:\n- Ubuntu/Debian: sudo apt install nodejs npm\n- macOS: brew install node\n- Windows: Download from nodejs.org", err)
		}
	}

	fmt.Printf("🔄 Converting OpenAPI to Postman: %s → %s\n", cfg.SwaggerInputPath, cfg.OutputPath)
	if err := g.runner.Run(ctx, cfg.WorkingDir, name, args...); err != nil {
		return fmt.Errorf("convert openapi to postman: %w\n💡 Troubleshooting:\n- Verify swagger.json is valid JSON\n- Check Node.js/npm installation\n- Try: npx openapi-to-postmanv2 --version", err)
	}

	// Validate output file was created
	if _, err := os.Stat(cfg.OutputPath); err != nil {
		return fmt.Errorf("postman collection not generated at %s: %w", cfg.OutputPath, err)
	}

	fmt.Printf("✅ Postman collection generated: %s\n", cfg.OutputPath)
	return nil
}

func buildOptions(in map[string]string) []string {
	if len(in) == 0 {
		return nil
	}
	keys := make([]string, 0, len(in))
	for k := range in {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	out := make([]string, 0, len(keys))
	for _, k := range keys {
		out = append(out, k+"="+in[k])
	}
	return out
}

func (cfg *Config) normalize() error {
	workingDir, err := resolveWorkingDir(cfg.WorkingDir)
	if err != nil {
		return err
	}
	cfg.WorkingDir = workingDir

	if cfg.OutputPath == "" {
		cfg.OutputPath = filepath.Join(cfg.WorkingDir, "postman_collection.json")
	} else if !filepath.IsAbs(cfg.OutputPath) {
		cfg.OutputPath = filepath.Join(cfg.WorkingDir, cfg.OutputPath)
	}

	if cfg.Swag.Enabled {
		if cfg.Swag.MainFile == "" {
			cfg.Swag.MainFile = "cmd/main.go"
		}
		if cfg.Swag.OutputDir == "" {
			cfg.Swag.OutputDir = "docs"
		}
	}

	if cfg.SwaggerInputPath == "" {
		if !cfg.Swag.Enabled {
			return errors.New("swagger_input_path is required when swag is disabled")
		}
		cfg.SwaggerInputPath = filepath.Join(cfg.Swag.OutputDir, "swagger.json")
	}

	if !filepath.IsAbs(cfg.SwaggerInputPath) {
		cfg.SwaggerInputPath = filepath.Join(cfg.WorkingDir, cfg.SwaggerInputPath)
	}

	if cfg.Postman.CLIPath == "" {
		cfg.Postman.CLIPath = "openapi2postmanv2"
	}

	return nil
}

func resolveWorkingDir(workingDir string) (string, error) {
	if workingDir != "" {
		return workingDir, nil
	}
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("resolve working directory: %w", err)
	}
	return wd, nil
}

func renameCollection(path string, name string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read generated postman file: %w", err)
	}

	var c map[string]any
	if err := json.Unmarshal(content, &c); err != nil {
		return fmt.Errorf("parse generated postman file: %w", err)
	}

	info, ok := c["info"].(map[string]any)
	if !ok {
		info = map[string]any{}
		c["info"] = info
	}
	info["name"] = name

	updated, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("encode postman file: %w", err)
	}
	updated = append(updated, '\n')

	if err := os.WriteFile(path, updated, 0o644); err != nil {
		return fmt.Errorf("write postman file: %w", err)
	}

	return nil
}

// findMainFile attempts to locate the main.go file in common locations
func (g *Generator) findMainFile(workingDir string) string {
	candidates := []string{
		"main.go",
		"cmd/main.go", 
		"cmd/api/main.go",
		"cmd/server/main.go",
		"cmd/app/main.go",
		"app/main.go",
		"api/main.go",
		"server/main.go",
		"service/main.go",
		"cmd/service/main.go",
	}
	
	fmt.Printf("🔍 Auto-detecting main.go file...\n")
	for _, candidate := range candidates {
		fullPath := filepath.Join(workingDir, candidate)
		if g.fileExist(fullPath) {
			fmt.Printf("✅ Found main.go at: %s\n", candidate)
			return candidate
		}
		fmt.Printf("   ❌ Not found: %s\n", candidate)
	}
	
	fmt.Printf("⚠️  Main.go not found in common locations.\n")
	return ""
}

// getSwaggerCandidates returns potential swagger file locations
func (g *Generator) getSwaggerCandidates(workingDir string) []string {
	baseDirs := []string{
		"docs",
		"api/docs", 
		"cmd/postman",
		"swagger",
		"openapi",
		".",
	}
	
	files := []string{
		"swagger.json",
		"openapi.yaml", 
		"openapi.yml",
		"swagger.yaml",
		"swagger.yml",
		"api.yaml",
		"api.yml",
		"spec.yaml",
		"spec.yml",
	}
	
	var candidates []string
	for _, dir := range baseDirs {
		for _, file := range files {
			candidates = append(candidates, filepath.Join(dir, file))
		}
	}
	
	return candidates
}

// tryAutoScaffold attempts to create a basic project structure with minimal main.go
func (g *Generator) tryAutoScaffold(workingDir string) error {
	mainPath := filepath.Join(workingDir, "main.go")
	if g.fileExist(mainPath) {
		return fmt.Errorf("main.go already exists")
	}
	
	// Create basic main.go with minimal swagger annotations
	basicMainGo := `package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

// @title		API
// @version		1.0
// @description	Auto-generated API
// @host		localhost:8080
// @BasePath	/api/v1
func main() {
	r := gin.Default()
	
	// @Summary		Health check
	// @Tags		health
	// @Success		200	{object}	map[string]string
	// @Router		/health [get]
	r.GET("/api/v1/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	
	r.Run(":8080")
}
`
	err := os.WriteFile(mainPath, []byte(basicMainGo), 0644)
	if err != nil {
		return fmt.Errorf("create basic main.go: %w", err)
	}
	
	fmt.Printf("🎆 Created basic main.go with swagger annotations\n")
	return nil
}

// validateMainFileLenient checks swagger annotations but allows missing ones with warnings
func (g *Generator) validateMainFileLenient(mainPath string) error {
	content, err := os.ReadFile(mainPath)
	if err != nil {
		return fmt.Errorf("read main file %s: %w", mainPath, err)
	}

	contentStr := string(content)
	
	// Check for required annotations
	required := []struct {
		annotation string
		message    string
		defaultVal string
	}{
		{"@title", "Missing @title annotation", "API"},
		{"@version", "Missing @version annotation", "1.0"},
		{"@host", "Missing @host annotation", "localhost:8080"},
		{"@BasePath", "Missing @BasePath annotation", "/api/v1"},
	}
	
	missing := []string{}
	warnings := []string{}
	for _, req := range required {
		if !strings.Contains(contentStr, req.annotation) {
			missing = append(missing, req.message)
			warnings = append(warnings, fmt.Sprintf("// %s %s", req.annotation, req.defaultVal))
		}
	}
	
	// Only show warnings for missing annotations, don't fail
	if len(missing) > 0 {
		fmt.Printf("⚠️  Missing swagger annotations (will use defaults):\n")
		for _, warning := range missing {
			fmt.Printf("   - %s\n", warning)
		}
		fmt.Printf("\n💡 To fix, add these to %s:\n%s\n\n", mainPath, strings.Join(warnings, "\n"))
	}
	
	// Check for docs import (non-fatal)
	if !strings.Contains(contentStr, `/docs"`) && !strings.Contains(contentStr, "/docs\"") {
		fmt.Printf("⚠️  Docs import missing - swag may fail (add: import _ \"yourmodule/docs\")\n")
	}
	
	// Always return nil for lenient validation
	return nil
}
