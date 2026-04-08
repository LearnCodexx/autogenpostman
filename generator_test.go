package generatepostmanfile

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

type calledCmd struct {
	dir  string
	name string
	args []string
}

type fakeRunner struct {
	calls []calledCmd
	err   error
}

func (f *fakeRunner) Run(_ context.Context, dir string, name string, args ...string) error {
	f.calls = append(f.calls, calledCmd{dir: dir, name: name, args: args})
	if f.err == nil {
		maybeWriteOutputFile(tinyCollectionJSON(), args...)
	}
	return f.err
}

type selectiveRunner struct {
	calls []calledCmd
}

func (s *selectiveRunner) Run(_ context.Context, dir string, name string, args ...string) error {
	s.calls = append(s.calls, calledCmd{dir: dir, name: name, args: args})
	if name == "go" {
		return errors.New("cannot run go swag")
	}
	maybeWriteOutputFile(tinyCollectionJSON(), args...)
	return nil
}

func maybeWriteOutputFile(content []byte, args ...string) {
	for i := 0; i < len(args)-1; i++ {
		if args[i] != "-o" {
			continue
		}
		out := args[i+1]
		if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
			return
		}
		_ = os.WriteFile(out, content, 0o644)
		return
	}
}

func tinyCollectionJSON() []byte {
	return []byte(`{"info":{"name":"tmp"},"item":[]}`)
}

func TestNormalizeRequiresSwaggerPathWhenSwagDisabled(t *testing.T) {
	cfg := Config{WorkingDir: "/tmp", Swag: SwagConfig{Enabled: false}}
	if err := cfg.normalize(); err == nil {
		t.Fatal("expected error when swagger input path is empty and swag is disabled")
	}
}

func TestGenerateRunsSwagAndNpx(t *testing.T) {
	runner := &fakeRunner{}
	g := &Generator{runner: runner}

	cfg := Config{
		WorkingDir:       "/tmp/project",
		OutputPath:       "postman.json",
		SwaggerInputPath: "docs/swagger.json",
		Pretty:           true,
		Swag: SwagConfig{
			Enabled:         true,
			MainFile:        "cmd/main.go",
			OutputDir:       "docs",
			ParseDependency: true,
			ParseInternal:   true,
			InstanceName:    "usersvc",
		},
		Postman: PostmanConfig{
			Options: map[string]string{
				"folderStrategy": "Tags",
			},
		},
	}

	err := g.Generate(context.Background(), cfg)
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}

	if len(runner.calls) != 2 {
		t.Fatalf("expected 2 command calls, got %d", len(runner.calls))
	}

	first := runner.calls[0]
	if first.name != "swag" {
		t.Fatalf("expected first command swag, got %s", first.name)
	}

	second := runner.calls[1]
	if second.name != "npx" {
		t.Fatalf("expected second command npx, got %s", second.name)
	}

	if second.args[0] != "--yes" || second.args[1] != "openapi-to-postmanv2" {
		t.Fatalf("unexpected npx args prefix: %#v", second.args)
	}
}

func TestGenerateRunsLocalCLI(t *testing.T) {
	runner := &fakeRunner{}
	g := &Generator{runner: runner}

	cfg := Config{
		WorkingDir:       "/tmp/project",
		OutputPath:       "postman.json",
		SwaggerInputPath: "docs/swagger.json",
		Postman: PostmanConfig{
			UseLocalCLI: true,
			CLIPath:     "openapi2postmanv2",
		},
	}

	err := g.Generate(context.Background(), cfg)
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}

	if len(runner.calls) != 1 {
		t.Fatalf("expected 1 command call, got %d", len(runner.calls))
	}

	if runner.calls[0].name != "openapi2postmanv2" {
		t.Fatalf("expected local CLI command, got %s", runner.calls[0].name)
	}
}

func TestGeneratePropagatesRunnerError(t *testing.T) {
	runner := &fakeRunner{err: errors.New("boom")}
	g := &Generator{runner: runner}

	cfg := Config{
		WorkingDir:       "/tmp/project",
		OutputPath:       "postman.json",
		SwaggerInputPath: "docs/swagger.json",
	}

	err := g.Generate(context.Background(), cfg)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestGenerateAutoUsesSwagWhenAvailable(t *testing.T) {
	runner := &fakeRunner{}
	g := &Generator{
		runner: runner,
		lookPath: func(file string) (string, error) {
			if file == "swag" {
				return "/usr/local/bin/swag", nil
			}
			return "", errors.New("not found")
		},
		fileExist: func(path string) bool { return false },
	}

	err := g.GenerateAuto(context.Background(), AutoConfig{
		WorkingDir: "/tmp/project",
		OutputPath: "docs/postman.json",
		Pretty:     true,
	})
	if err != nil {
		t.Fatalf("GenerateAuto returned error: %v", err)
	}

	if len(runner.calls) != 2 {
		t.Fatalf("expected 2 command calls, got %d", len(runner.calls))
	}
	if runner.calls[0].name != "swag" {
		t.Fatalf("expected first command swag, got %s", runner.calls[0].name)
	}
	if runner.calls[1].name != "npx" {
		t.Fatalf("expected second command npx, got %s", runner.calls[1].name)
	}
}

func TestGenerateAutoFallbackToExistingOpenAPIWhenSwagMissing(t *testing.T) {
	dir := t.TempDir()
	openapiPath := filepath.Join(dir, "docs", "openapi.yaml")
	if err := os.MkdirAll(filepath.Dir(openapiPath), 0o755); err != nil {
		t.Fatalf("mkdir docs: %v", err)
	}
	if err := os.WriteFile(openapiPath, []byte("openapi: 3.0.0\n"), 0o644); err != nil {
		t.Fatalf("write openapi fixture: %v", err)
	}

	runner := &selectiveRunner{}
	g := &Generator{
		runner: runner,
		lookPath: func(string) (string, error) {
			return "", errors.New("swag not found")
		},
		fileExist: func(path string) bool {
			_, err := os.Stat(path)
			return err == nil
		},
	}

	err := g.GenerateAuto(context.Background(), AutoConfig{
		WorkingDir: dir,
		OutputPath: "docs/postman.json",
	})
	if err != nil {
		t.Fatalf("GenerateAuto returned error: %v", err)
	}

	if len(runner.calls) != 2 {
		t.Fatalf("expected 2 command calls, got %d", len(runner.calls))
	}
	if runner.calls[0].name != "go" {
		t.Fatalf("expected first command go (swag fallback), got %s", runner.calls[0].name)
	}
	if runner.calls[1].name != "npx" {
		t.Fatalf("expected second command npx, got %s", runner.calls[1].name)
	}

	args := runner.calls[1].args
	foundSource := false
	for i := 0; i < len(args)-1; i++ {
		if args[i] == "-s" && args[i+1] == openapiPath {
			foundSource = true
			break
		}
	}
	if !foundSource {
		t.Fatalf("expected converter source path %s, args=%v", openapiPath, args)
	}
}

func TestGenerateAutoReturnsErrorWhenNoSwagAndNoOpenAPI(t *testing.T) {
	runner := &fakeRunner{err: errors.New("boom")}
	g := &Generator{
		runner: runner,
		lookPath: func(string) (string, error) {
			return "", errors.New("swag not found")
		},
		fileExist: func(string) bool { return false },
	}

	err := g.GenerateAuto(context.Background(), AutoConfig{
		WorkingDir: "/tmp/project",
		OutputPath: "docs/postman.json",
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if len(runner.calls) != 1 {
		t.Fatalf("expected one command call (go swag), got %d", len(runner.calls))
	}
	if runner.calls[0].name != "go" {
		t.Fatalf("expected go command, got %s", runner.calls[0].name)
	}
}

func TestRenameCollection(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "collection.json")

	original := map[string]any{
		"info": map[string]any{"name": "old"},
		"item": []any{},
	}
	content, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal fixture: %v", err)
	}
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	if err := renameCollection(path, "new-collection"); err != nil {
		t.Fatalf("renameCollection error: %v", err)
	}

	updated, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read updated file: %v", err)
	}

	var got map[string]any
	if err := json.Unmarshal(updated, &got); err != nil {
		t.Fatalf("unmarshal updated file: %v", err)
	}

	info, ok := got["info"].(map[string]any)
	if !ok {
		t.Fatalf("missing info object: %#v", got)
	}
	if info["name"] != "new-collection" {
		t.Fatalf("unexpected name: %#v", info["name"])
	}
	if _, exists := got["item"]; !exists {
		t.Fatalf("expected item to be preserved, got: %#v", got)
	}
}
