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

	if cfg.OutputPath == "" {
		cfg.OutputPath = filepath.Join("docs", "postman_collection.json")
	}
	if cfg.MainFile == "" {
		cfg.MainFile = g.findMainFile(workingDir)
	}
	if cfg.SwagOutputDir == "" {
		// Use docs as default, more standard
		cfg.SwagOutputDir = filepath.Join(workingDir, "docs")
		// Create the directory if it doesn't exist
		if err := os.MkdirAll(cfg.SwagOutputDir, 0755); err != nil {
			return fmt.Errorf("create swag output dir: %w", err)
		}
	}
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
		lowLevelCfg.SwaggerInputPath = cfg.SwaggerInputPath
		return g.Generate(ctx, lowLevelCfg)
	}

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
		return nil
	}

	for _, candidate := range cfg.SwaggerCandidates {
		candidatePath := candidate
		if !filepath.IsAbs(candidatePath) {
			candidatePath = filepath.Join(workingDir, candidatePath)
		}
		if g.fileExist(candidatePath) {
			lowLevelCfg.SwaggerInputPath = candidatePath
			lowLevelCfg.Swag = SwagConfig{}
			return g.Generate(ctx, lowLevelCfg)
		}
	}

	return fmt.Errorf("auto swagger generation failed (%v) and no OpenAPI file found; provide SwaggerInputPath or create one of: %s", swagErr, strings.Join(cfg.SwaggerCandidates[:3], ", "))
}

// GenerateAuto is a package-level convenience API for one-call generation.
func GenerateAuto(ctx context.Context, cfg AutoConfig) error {
	return New().GenerateAuto(ctx, cfg)
}

func (g *Generator) runSwag(ctx context.Context, cfg Config) error {
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
			cmdName = "go"
			cmdArgs = append([]string{"run", "github.com/swaggo/swag/cmd/swag@latest"}, args...)
		}
	}

	if err := g.runner.Run(ctx, cfg.WorkingDir, cmdName, cmdArgs...); err != nil {
		return fmt.Errorf("generate swagger with swag: %w", err)
	}
	return nil
}

func (g *Generator) runConvert(ctx context.Context, cfg Config) error {
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

	if err := g.runner.Run(ctx, cfg.WorkingDir, name, args...); err != nil {
		return fmt.Errorf("convert openapi to postman: %w", err)
	}
	if _, err := os.Stat(cfg.OutputPath); err != nil {
		return fmt.Errorf("convert openapi to postman: output file not found at %s", cfg.OutputPath)
	}
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
	}
	
	for _, candidate := range candidates {
		fullPath := filepath.Join(workingDir, candidate)
		if g.fileExist(fullPath) {
			return candidate
		}
	}
	
	// Default fallback
	return "main.go"
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
