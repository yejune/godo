package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// ValidName matches lowercase letters, digits, and hyphens (must start with a letter).
var ValidName = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)

// CreateAgent creates a new agent definition file.
func CreateAgent(name string) error {
	agentDir := filepath.Join(".claude", "agents", "do")
	agentFile := filepath.Join(agentDir, name+".md")

	if FileExists(agentFile) {
		return fmt.Errorf("agent already exists at %s", agentFile)
	}

	if err := os.MkdirAll(agentDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", agentDir, err)
	}

	titleName := ToTitleCase(name)

	content := fmt.Sprintf(`---
name: %s
description: >
  TODO: Describe when Claude should delegate to this agent.
model: inherit
tools:
  - Read
  - Write
  - Edit
  - Grep
  - Glob
  - Bash
---

# %s Agent

## Expertise
TODO: Define the agent's domain expertise.

## Instructions
TODO: Add specific instructions for this agent.
`, name, titleName)

	if err := os.WriteFile(agentFile, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", agentFile, err)
	}

	fmt.Printf("Created agent: %s\n", agentFile)
	return nil
}

// CreateSkill creates a new skill definition file.
func CreateSkill(name string) error {
	skillName := "do-" + name
	skillDir := filepath.Join(".claude", "skills", skillName)
	skillFile := filepath.Join(skillDir, "SKILL.md")

	if FileExists(skillFile) {
		return fmt.Errorf("skill already exists at %s", skillFile)
	}

	if err := os.MkdirAll(skillDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", skillDir, err)
	}

	today := time.Now().Format("2006-01-02")

	content := fmt.Sprintf(`---
name: %s
description: >
  TODO: Describe what this skill provides.
license: Apache-2.0
compatibility: Designed for Claude Code
user-invocable: false
metadata:
  version: "1.0.0"
  category: "domain"
  status: "experimental"
  updated: "%s"
  tags: "TODO"
---

# %s

## Overview
TODO: Describe the skill's purpose and capabilities.

## Usage
TODO: Add usage examples and patterns.
`, skillName, today, skillName)

	if err := os.WriteFile(skillFile, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", skillFile, err)
	}

	fmt.Printf("Created skill: %s\n", skillFile)
	return nil
}

// ToTitleCase converts "my-agent-name" to "My Agent Name".
func ToTitleCase(name string) string {
	parts := strings.Split(name, "-")
	for i, p := range parts {
		if len(p) > 0 {
			parts[i] = strings.ToUpper(p[:1]) + p[1:]
		}
	}
	return strings.Join(parts, " ")
}

// FileExists returns true if the path exists.
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
