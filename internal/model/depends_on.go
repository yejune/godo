package model

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// DependsOn represents typed dependency declarations in agent frontmatter.
// Each field corresponds to a dependency category. Nil/empty means no
// dependencies of that type. Pointer receiver on Frontmatter allows
// omitempty to work correctly (nil = no depends_on field in YAML).
type DependsOn struct {
	Phases         []string       `yaml:"phases,omitempty,flow"`
	Artifacts      []ArtifactDep  `yaml:"artifacts,omitempty"`
	Agents         []AgentDep     `yaml:"agents,omitempty"`
	Env            []string       `yaml:"env,omitempty,flow"`
	Services       []ServiceDep   `yaml:"services,omitempty"`
	ChecklistItems []string       `yaml:"checklist_items,omitempty,flow"`
}

// ArtifactDep represents a file dependency with optional required flag.
// Supports hybrid YAML: scalar string (shorthand) or object {path, required}.
//
// Scalar form:  "plan.md"           -> ArtifactDep{Path: "plan.md", Required: true}
// Object form:  {path: "a.md", required: false} -> ArtifactDep{Path: "a.md", Required: false}
type ArtifactDep struct {
	Path     string `yaml:"path"`
	Required bool   `yaml:"required"`
}

// UnmarshalYAML implements custom YAML unmarshaling for ArtifactDep.
// Handles both scalar string and object forms.
func (a *ArtifactDep) UnmarshalYAML(value *yaml.Node) error {
	switch value.Kind {
	case yaml.ScalarNode:
		// Shorthand: "plan.md" -> {Path: "plan.md", Required: true}
		a.Path = value.Value
		a.Required = true
		return nil
	case yaml.MappingNode:
		// Object form: {path: "plan.md", required: false}
		// Use an alias type to avoid infinite recursion
		type artifactDepAlias ArtifactDep
		aux := &artifactDepAlias{Required: true} // default required=true
		if err := value.Decode(aux); err != nil {
			return fmt.Errorf("invalid artifact dependency: %w", err)
		}
		*a = ArtifactDep(*aux)
		return nil
	default:
		return fmt.Errorf("artifact dependency must be a string or object, got %v", value.Kind)
	}
}

// AgentDep represents a dependency on another agent's checklist completion.
// Supports hybrid YAML: scalar string (shorthand) or object {name, items}.
//
// Scalar form:  "expert-backend"    -> AgentDep{Name: "expert-backend", Items: nil}
// Object form:  {name: "expert-backend", items: ["#1"]} -> AgentDep{Name: "expert-backend", Items: ["#1"]}
type AgentDep struct {
	Name  string   `yaml:"name"`
	Items []string `yaml:"items,omitempty,flow"`
}

// UnmarshalYAML implements custom YAML unmarshaling for AgentDep.
// Handles both scalar string and object forms.
func (a *AgentDep) UnmarshalYAML(value *yaml.Node) error {
	switch value.Kind {
	case yaml.ScalarNode:
		// Shorthand: "expert-backend" -> {Name: "expert-backend"}
		a.Name = value.Value
		a.Items = nil
		return nil
	case yaml.MappingNode:
		// Object form: {name: "expert-backend", items: ["#1", "#2"]}
		type agentDepAlias AgentDep
		aux := &agentDepAlias{}
		if err := value.Decode(aux); err != nil {
			return fmt.Errorf("invalid agent dependency: %w", err)
		}
		*a = AgentDep(*aux)
		return nil
	default:
		return fmt.Errorf("agent dependency must be a string or object, got %v", value.Kind)
	}
}

// ServiceDep represents a dependency on a Docker Compose service.
// Supports hybrid YAML: scalar string (shorthand) or object {name, healthcheck}.
//
// Scalar form:  "postgres"          -> ServiceDep{Name: "postgres", Healthcheck: false}
// Object form:  {name: "postgres", healthcheck: true} -> ServiceDep{Name: "postgres", Healthcheck: true}
type ServiceDep struct {
	Name        string `yaml:"name"`
	Healthcheck bool   `yaml:"healthcheck"`
}

// UnmarshalYAML implements custom YAML unmarshaling for ServiceDep.
// Handles both scalar string and object forms.
func (s *ServiceDep) UnmarshalYAML(value *yaml.Node) error {
	switch value.Kind {
	case yaml.ScalarNode:
		// Shorthand: "postgres" -> {Name: "postgres", Healthcheck: false}
		s.Name = value.Value
		s.Healthcheck = false
		return nil
	case yaml.MappingNode:
		// Object form: {name: "postgres", healthcheck: true}
		type serviceDepAlias ServiceDep
		aux := &serviceDepAlias{}
		if err := value.Decode(aux); err != nil {
			return fmt.Errorf("invalid service dependency: %w", err)
		}
		*s = ServiceDep(*aux)
		return nil
	default:
		return fmt.Errorf("service dependency must be a string or object, got %v", value.Kind)
	}
}
