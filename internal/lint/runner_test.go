package lint

import "testing"

func Test_DetectLanguage_go(t *testing.T) {
	if got := DetectLanguage("main.go"); got != LangGo {
		t.Errorf("got %q, want %q", got, LangGo)
	}
}

func Test_DetectLanguage_python(t *testing.T) {
	cases := []string{"app.py", "stub.pyi"}
	for _, f := range cases {
		if got := DetectLanguage(f); got != LangPython {
			t.Errorf("%s: got %q, want %q", f, got, LangPython)
		}
	}
}

func Test_DetectLanguage_typescript(t *testing.T) {
	cases := []string{"app.ts", "comp.tsx", "mod.mts", "mod.cts"}
	for _, f := range cases {
		if got := DetectLanguage(f); got != LangTypeScript {
			t.Errorf("%s: got %q, want %q", f, got, LangTypeScript)
		}
	}
}

func Test_DetectLanguage_javascript(t *testing.T) {
	cases := []string{"app.js", "comp.jsx", "mod.mjs", "mod.cjs"}
	for _, f := range cases {
		if got := DetectLanguage(f); got != LangJavaScript {
			t.Errorf("%s: got %q, want %q", f, got, LangJavaScript)
		}
	}
}

func Test_DetectLanguage_rust(t *testing.T) {
	if got := DetectLanguage("lib.rs"); got != LangRust {
		t.Errorf("got %q, want %q", got, LangRust)
	}
}

func Test_DetectLanguage_unknown(t *testing.T) {
	cases := []string{"readme.md", "config.yaml", "image.png", ""}
	for _, f := range cases {
		if got := DetectLanguage(f); got != LangUnknown {
			t.Errorf("%s: got %q, want %q", f, got, LangUnknown)
		}
	}
}

func Test_IsCodeFile(t *testing.T) {
	if !IsCodeFile("main.go") {
		t.Error("main.go should be a code file")
	}
	if IsCodeFile("readme.md") {
		t.Error("readme.md should not be a code file")
	}
}

func Test_GroupFilesByLanguage(t *testing.T) {
	files := []string{"a.go", "b.go", "c.py", "d.md", "e.ts"}
	groups := GroupFilesByLanguage(files)

	if len(groups[LangGo]) != 2 {
		t.Errorf("Go files: got %d, want 2", len(groups[LangGo]))
	}
	if len(groups[LangPython]) != 1 {
		t.Errorf("Python files: got %d, want 1", len(groups[LangPython]))
	}
	if len(groups[LangTypeScript]) != 1 {
		t.Errorf("TypeScript files: got %d, want 1", len(groups[LangTypeScript]))
	}
	if _, exists := groups[LangUnknown]; exists {
		t.Error("unknown files should not be grouped")
	}
}

func Test_AllLinters_returns_all_languages(t *testing.T) {
	linters := AllLinters()
	if len(linters) != 5 {
		t.Fatalf("expected 5 linters, got %d", len(linters))
	}

	langs := make(map[Language]bool)
	for _, l := range linters {
		langs[l.Language] = true
	}
	for _, expected := range []Language{LangGo, LangPython, LangTypeScript, LangJavaScript, LangRust} {
		if !langs[expected] {
			t.Errorf("missing linter for %q", expected)
		}
	}
}

func Test_LinterForLanguage_found(t *testing.T) {
	info, ok := LinterForLanguage(LangGo)
	if !ok {
		t.Fatal("expected to find linter for Go")
	}
	if info.Command != "go" {
		t.Errorf("Command: got %q, want %q", info.Command, "go")
	}
}

func Test_LinterForLanguage_not_found(t *testing.T) {
	_, ok := LinterForLanguage(LangUnknown)
	if ok {
		t.Error("should not find linter for unknown language")
	}
}

func Test_RunLinter_unknown_returns_nil(t *testing.T) {
	diags := RunLinter(LangUnknown, []string{"test.txt"}, ".")
	if diags != nil {
		t.Errorf("expected nil for unknown language, got %v", diags)
	}
}

func Test_ParseGoVetOutput_valid(t *testing.T) {
	output := "main.go:10:5: unreachable code\nutils.go:20:1: unused variable\n"
	diags := ParseGoVetOutput(output)
	if len(diags) != 2 {
		t.Fatalf("expected 2 diagnostics, got %d", len(diags))
	}
	if diags[0].File != "main.go" {
		t.Errorf("File: got %q, want %q", diags[0].File, "main.go")
	}
	if diags[0].Line != 10 {
		t.Errorf("Line: got %d, want 10", diags[0].Line)
	}
	if diags[0].Message != "unreachable code" {
		t.Errorf("Message: got %q", diags[0].Message)
	}
}

func Test_ParseGoVetOutput_empty(t *testing.T) {
	diags := ParseGoVetOutput("")
	if len(diags) != 0 {
		t.Errorf("expected 0 diagnostics for empty output, got %d", len(diags))
	}
}
