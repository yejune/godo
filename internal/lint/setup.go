package lint

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// PackageManager represents an available package manager on the system.
type PackageManager struct {
	Name    string
	Command string
}

// LinterInstallInfo holds installation methods for a linter.
type LinterInstallInfo struct {
	DisplayName    string
	Language       Language
	InstallMethods map[string][]string
	DownloadURL    string
	BuiltIn        bool // true for go vet (no install needed)
}

// AllLinterInstallInfo returns installation info for all supported linters.
func AllLinterInstallInfo() []LinterInstallInfo {
	return []LinterInstallInfo{
		{
			DisplayName: "go vet",
			Language:    LangGo,
			BuiltIn:     true,
		},
		{
			DisplayName: "ruff",
			Language:    LangPython,
			InstallMethods: map[string][]string{
				"pip":    {"pip", "install", "ruff"},
				"pip3":   {"pip3", "install", "ruff"},
				"brew":   {"brew", "install", "ruff"},
				"scoop":  {"scoop", "install", "ruff"},
				"winget": {"winget", "install", "astral-sh.ruff"},
			},
			DownloadURL: "https://docs.astral.sh/ruff/installation/",
		},
		{
			DisplayName: "tsc (TypeScript)",
			Language:    LangTypeScript,
			InstallMethods: map[string][]string{
				"npm": {"npm", "install", "-g", "typescript"},
			},
			DownloadURL: "https://www.typescriptlang.org/download",
		},
		{
			DisplayName: "eslint",
			Language:    LangJavaScript,
			InstallMethods: map[string][]string{
				"npm": {"npm", "install", "-g", "eslint"},
			},
			DownloadURL: "https://eslint.org/docs/latest/use/getting-started",
		},
		{
			DisplayName: "cargo clippy",
			Language:    LangRust,
			InstallMethods: map[string][]string{
				"rustup": {"rustup", "component", "add", "clippy"},
			},
			DownloadURL: "https://www.rust-lang.org/tools/install",
		},
	}
}

// DetectPackageManagers scans the system for available package managers.
func DetectPackageManagers() []PackageManager {
	candidates := []PackageManager{
		{"brew", "brew"},
		{"pip", "pip"},
		{"pip3", "pip3"},
		{"npm", "npm"},
		{"scoop", "scoop"},
		{"choco", "choco"},
		{"winget", "winget"},
		{"apt", "apt"},
		{"rustup", "rustup"},
	}

	var available []PackageManager
	for _, pm := range candidates {
		if _, err := exec.LookPath(pm.Command); err == nil {
			available = append(available, pm)
		}
	}
	return available
}

// ScanProjectLanguages scans the project directory for code files and returns detected languages.
func ScanProjectLanguages(projectDir string) []Language {
	seen := make(map[Language]bool)

	_ = filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		name := info.Name()
		if info.IsDir() {
			if strings.HasPrefix(name, ".") || name == "node_modules" || name == "vendor" || name == "__pycache__" || name == "target" {
				return filepath.SkipDir
			}
			return nil
		}
		lang := DetectLanguage(path)
		if lang != LangUnknown {
			seen[lang] = true
		}
		return nil
	})

	var langs []Language
	for lang := range seen {
		langs = append(langs, lang)
	}
	return langs
}

// SetupStatus holds the result of checking linter installation status.
type SetupStatus struct {
	Installed []string
	Missing   []LinterInstallInfo
}

// CheckSetupStatus checks which linters are installed and which are missing for the given languages.
func CheckSetupStatus(langs []Language) SetupStatus {
	allInfo := AllLinterInstallInfo()
	var status SetupStatus

	for _, info := range allInfo {
		needed := false
		for _, lang := range langs {
			if info.Language == lang {
				needed = true
				break
			}
		}
		if !needed {
			continue
		}
		if info.BuiltIn {
			status.Installed = append(status.Installed, info.DisplayName+" (built-in)")
			continue
		}
		if CheckLinterInstalled(info.Language) {
			status.Installed = append(status.Installed, info.DisplayName)
			continue
		}
		status.Missing = append(status.Missing, info)
	}

	return status
}

// InstallOption represents a single install method for a missing linter.
type InstallOption struct {
	Label string
	Args  []string
}

// GetInstallOptions returns available install options for a linter given available package managers.
func GetInstallOptions(info LinterInstallInfo, managers []PackageManager) []InstallOption {
	var options []InstallOption
	for _, pm := range managers {
		if args, ok := info.InstallMethods[pm.Name]; ok {
			label := strings.Join(args, " ")
			options = append(options, InstallOption{Label: label, Args: args})
		}
	}
	return options
}

// RunInstall executes a linter installation command.
func RunInstall(args []string) error {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// PrintSetupStatus prints the current linter installation status.
func PrintSetupStatus(status SetupStatus) {
	if len(status.Installed) > 0 {
		fmt.Println("Installed:")
		for _, name := range status.Installed {
			fmt.Printf("  [ok] %s\n", name)
		}
		fmt.Println()
	}

	if len(status.Missing) == 0 {
		fmt.Println("All linters are installed!")
		return
	}

	fmt.Println("Missing:")
	for _, info := range status.Missing {
		fmt.Printf("  [--] %s\n", info.DisplayName)
	}
	fmt.Println()
}
