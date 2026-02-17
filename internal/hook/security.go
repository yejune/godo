package hook

import "regexp"

// SecurityPolicy defines tool access control rules for PreToolUse events.
type SecurityPolicy struct {
	DenyFilePatterns         []*regexp.Regexp
	AskFilePatterns          []*regexp.Regexp
	DenyBashPatterns         []*regexp.Regexp
	AskBashPatterns          []*regexp.Regexp
	SensitiveContentPatterns []*regexp.Regexp
}

// CompilePatterns compiles a list of pattern strings into case-insensitive regexp objects.
func CompilePatterns(patterns []string) []*regexp.Regexp {
	result := make([]*regexp.Regexp, 0, len(patterns))
	for _, p := range patterns {
		re, err := regexp.Compile("(?i)" + p)
		if err != nil {
			continue
		}
		result = append(result, re)
	}
	return result
}

// DenyFilePatternStrings defines files that should NEVER be modified.
var DenyFilePatternStrings = []string{
	// Secrets and credentials
	`secrets?\.(json|ya?ml|toml)$`,
	`credentials?\.(json|ya?ml|toml)$`,
	`\.secrets/.*`,
	`secrets/.*`,
	// SSH and certificates
	`\.ssh/.*`,
	`id_rsa.*`,
	`id_ed25519.*`,
	`\.pem$`,
	`\.key$`,
	`\.crt$`,
	// Git internals
	`\.git/.*`,
	// Cloud credentials
	`\.aws/.*`,
	`\.gcloud/.*`,
	`\.azure/.*`,
	`\.kube/.*`,
	// Token files
	`\.token$`,
	`\.tokens/.*`,
	`auth\.json$`,
}

// AskFilePatternStrings defines files that require user confirmation.
var AskFilePatternStrings = []string{
	// Lock files
	`package-lock\.json$`,
	`yarn\.lock$`,
	`pnpm-lock\.ya?ml$`,
	`Gemfile\.lock$`,
	`Cargo\.lock$`,
	`poetry\.lock$`,
	`composer\.lock$`,
	`Pipfile\.lock$`,
	`uv\.lock$`,
	// Critical configs
	`tsconfig\.json$`,
	`pyproject\.toml$`,
	`Cargo\.toml$`,
	`package\.json$`,
	`docker-compose\.ya?ml$`,
	`Dockerfile$`,
	`\.dockerignore$`,
	// CI/CD configs
	`\.github/workflows/.*\.ya?ml$`,
	`\.gitlab-ci\.ya?ml$`,
	`\.circleci/.*`,
	`Jenkinsfile$`,
	// Infrastructure
	`terraform/.*\.tf$`,
	`\.terraform/.*`,
	`kubernetes/.*\.ya?ml$`,
	`k8s/.*\.ya?ml$`,
}

// DenyBashPatternStrings defines dangerous Bash commands that should NEVER be executed.
var DenyBashPatternStrings = []string{
	// Database deletion commands
	`supabase\s+db\s+reset`,
	`supabase\s+projects?\s+delete`,
	`neon\s+database\s+delete`,
	`neon\s+projects?\s+delete`,
	`pscale\s+database\s+delete`,
	`railway\s+delete`,
	`vercel\s+env\s+rm`,
	`vercel\s+projects?\s+rm`,
	// SQL dangerous commands
	`DROP\s+DATABASE`,
	`DROP\s+SCHEMA`,
	`TRUNCATE\s+TABLE`,
	// Unix dangerous file operations
	`rm\s+-rf\s+/`,
	`rm\s+-rf\s+~`,
	`rm\s+-rf\s+\*`,
	`rm\s+-rf\s+\.\*`,
	`rm\s+-rf\s+\.git\b`,
	`rm\s+-rf\s+node_modules\s*$`,
	// Windows dangerous file operations (CMD)
	`rd\s+/s\s+/q\s+[A-Za-z]:\\`,
	`rmdir\s+/s\s+/q\s+[A-Za-z]:\\`,
	`del\s+/f\s+/q\s+[A-Za-z]:\\`,
	`rd\s+/s\s+/q\s+\.git\b`,
	`del\s+/s\s+/q\s+\*\.\*`,
	`format\s+[A-Za-z]:`,
	// Windows dangerous file operations (PowerShell)
	`Remove-Item\s+.*-Recurse\s+.*-Force\s+[A-Za-z]:\\`,
	`Remove-Item\s+.*-Recurse\s+.*-Force\s+~`,
	`Remove-Item\s+.*-Recurse\s+.*-Force\s+\.git\b`,
	// Git dangerous commands
	`git\s+push\s+.*--force\s+origin\s+(main|master)`,
	`git\s+branch\s+-D\s+(main|master)`,
	// Cloud infrastructure deletion
	`terraform\s+destroy`,
	`pulumi\s+destroy`,
	`aws\s+.*\s+delete-`,
	`gcloud\s+.*\s+delete\b`,
	`az\s+group\s+delete`,
	`az\s+storage\s+account\s+delete`,
	`az\s+sql\s+server\s+delete`,
	// Docker dangerous commands
	`docker\s+system\s+prune\s+(-a|--all)`,
	`docker\s+image\s+prune\s+(-a|--all)`,
	`docker\s+container\s+prune`,
	`docker\s+volume\s+prune`,
	`docker\s+network\s+prune`,
	// Classic dangerous patterns
	`:\(\)\{\s*:\|:&\s*\};:`, // Fork bomb
	`mkfs\.`,
	`>\s*/dev/sda`,
	`dd\s+if=/dev/zero\s+of=/dev/sda`,
}

// AskBashPatternStrings defines Bash commands that require user confirmation.
var AskBashPatternStrings = []string{
	`prisma\s+migrate\s+reset`,
	`prisma\s+db\s+push\s+--force`,
	`drizzle-kit\s+push`,
	`git\s+push\s+.*--force`,
	`git\s+reset\s+--hard`,
	`git\s+clean\s+-fd`,
	`npm\s+cache\s+clean`,
	`yarn\s+cache\s+clean`,
	`pnpm\s+store\s+prune`,
}

// SensitiveContentPatternStrings defines content patterns that indicate sensitive data.
var SensitiveContentPatternStrings = []string{
	`-----BEGIN\s+(RSA\s+)?PRIVATE\s+KEY-----`,
	`-----BEGIN\s+CERTIFICATE-----`,
	`sk-[a-zA-Z0-9]{32,}`,       // OpenAI API keys
	`ghp_[a-zA-Z0-9]{36}`,       // GitHub tokens
	`gho_[a-zA-Z0-9]{36}`,       // GitHub OAuth tokens
	`glpat-[a-zA-Z0-9\-]{20}`,   // GitLab tokens
	`xox[baprs]-[a-zA-Z0-9\-]+`, // Slack tokens
	`AKIA[0-9A-Z]{16}`,          // AWS access keys
	`ya29\.[a-zA-Z0-9_\-]+`,     // Google OAuth tokens
}

// DefaultSecurityPolicy returns a SecurityPolicy with comprehensive security patterns.
func DefaultSecurityPolicy() *SecurityPolicy {
	return &SecurityPolicy{
		DenyFilePatterns:         CompilePatterns(DenyFilePatternStrings),
		AskFilePatterns:          CompilePatterns(AskFilePatternStrings),
		DenyBashPatterns:         CompilePatterns(DenyBashPatternStrings),
		AskBashPatterns:          CompilePatterns(AskBashPatternStrings),
		SensitiveContentPatterns: CompilePatterns(SensitiveContentPatternStrings),
	}
}
