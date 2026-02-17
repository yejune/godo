package hook

import (
	"testing"
)

func Test_CompilePatterns_valid(t *testing.T) {
	patterns := []string{`\.env$`, `secrets\.json$`}
	compiled := CompilePatterns(patterns)
	if len(compiled) != 2 {
		t.Fatalf("expected 2 patterns, got %d", len(compiled))
	}
}

func Test_CompilePatterns_skips_invalid(t *testing.T) {
	patterns := []string{`\.env$`, `[invalid`, `ok\.txt$`}
	compiled := CompilePatterns(patterns)
	if len(compiled) != 2 {
		t.Fatalf("expected 2 valid patterns (invalid skipped), got %d", len(compiled))
	}
}

func Test_CompilePatterns_case_insensitive(t *testing.T) {
	compiled := CompilePatterns([]string{`\.PEM$`})
	if len(compiled) != 1 {
		t.Fatal("expected 1 pattern")
	}
	if !compiled[0].MatchString("server.pem") {
		t.Error("case-insensitive match should match 'server.pem'")
	}
	if !compiled[0].MatchString("SERVER.PEM") {
		t.Error("case-insensitive match should match 'SERVER.PEM'")
	}
}

func Test_DefaultSecurityPolicy_deny_file_patterns(t *testing.T) {
	policy := DefaultSecurityPolicy()

	denyFiles := []string{
		"secrets.json",
		"credentials.yaml",
		".ssh/id_rsa",
		"id_ed25519.pub",
		"server.pem",
		"tls.key",
		"ca.crt",
		".git/config",
		".aws/credentials",
		"auth.json",
	}

	for _, f := range denyFiles {
		matched := false
		for _, re := range policy.DenyFilePatterns {
			if re.MatchString(f) {
				matched = true
				break
			}
		}
		if !matched {
			t.Errorf("expected deny match for %q", f)
		}
	}
}

func Test_DefaultSecurityPolicy_ask_file_patterns(t *testing.T) {
	policy := DefaultSecurityPolicy()

	askFiles := []string{
		"package-lock.json",
		"yarn.lock",
		"Cargo.lock",
		"tsconfig.json",
		"package.json",
		"docker-compose.yml",
		"Dockerfile",
		".github/workflows/ci.yml",
	}

	for _, f := range askFiles {
		matched := false
		for _, re := range policy.AskFilePatterns {
			if re.MatchString(f) {
				matched = true
				break
			}
		}
		if !matched {
			t.Errorf("expected ask match for %q", f)
		}
	}
}

func Test_DefaultSecurityPolicy_deny_bash_patterns(t *testing.T) {
	policy := DefaultSecurityPolicy()

	denyCmds := []string{
		"git push --force origin main",
		"terraform destroy",
		"docker system prune -a",
	}

	for _, cmd := range denyCmds {
		matched := false
		for _, re := range policy.DenyBashPatterns {
			if re.MatchString(cmd) {
				matched = true
				break
			}
		}
		if !matched {
			t.Errorf("expected deny match for bash command %q", cmd)
		}
	}
}

func Test_DefaultSecurityPolicy_ask_bash_patterns(t *testing.T) {
	policy := DefaultSecurityPolicy()

	askCmds := []string{
		"git push --force",
		"git clean -fd",
		"prisma migrate reset",
	}

	for _, cmd := range askCmds {
		matched := false
		for _, re := range policy.AskBashPatterns {
			if re.MatchString(cmd) {
				matched = true
				break
			}
		}
		if !matched {
			t.Errorf("expected ask match for bash command %q", cmd)
		}
	}
}

func Test_DefaultSecurityPolicy_sensitive_content_patterns(t *testing.T) {
	policy := DefaultSecurityPolicy()

	sensitiveContent := []string{
		"-----BEGIN RSA PRIVATE KEY-----",
		"-----BEGIN PRIVATE KEY-----",
		"-----BEGIN CERTIFICATE-----",
		"sk-abcdefghijklmnopqrstuvwxyz012345",
		"ghp_abcdefghijklmnopqrstuvwxyz0123456789",
		"AKIAIOSFODNN7EXAMPLE",
	}

	for _, content := range sensitiveContent {
		matched := false
		for _, re := range policy.SensitiveContentPatterns {
			if re.MatchString(content) {
				matched = true
				break
			}
		}
		if !matched {
			t.Errorf("expected sensitive match for %q", content)
		}
	}
}

func Test_DefaultSecurityPolicy_safe_files_not_denied(t *testing.T) {
	policy := DefaultSecurityPolicy()

	safeFiles := []string{
		"main.go",
		"internal/cli/root.go",
		"README.md",
		"Makefile",
	}

	for _, f := range safeFiles {
		for _, re := range policy.DenyFilePatterns {
			if re.MatchString(f) {
				t.Errorf("safe file %q should not match deny pattern %s", f, re.String())
			}
		}
	}
}
