# Architecture Document Template [HARD]

For complex tasks, write to `.do/jobs/{YY}/{MM}/{DD}/{title-kebab-case}/architecture.md`.
Written after analysis.md is complete. Refer to complexity-check.md for complexity criteria.

```markdown
# {Task Title} Architecture Design

## Overview

(1-2 sentence summary: what, why, how)

(ASCII architecture diagram)

---

## 1. Directory Structure

(File/folder tree — with role comments for each file)

---

## 2. Core Interfaces

(Define core types, interfaces, contracts at code level)

(Include JSDoc/docstring style descriptions for each interface)

---

## 3. Error Handling

(Error hierarchy, error code enums, error wrapping strategy)

---

## 4. Component Implementations

### 4.1 {Component A}

(Implementation details: class/function signatures, internal logic, config mapping)

### 4.2 {Component B}

(Implementation details)

---

## 5. Integration Layer

(Framework integration, DI configuration, module structure, etc.)

---

## 6. Configuration

### 6.1 Package/Build Configuration

(package.json, tsconfig, build tool settings, etc.)

### 6.2 Public API

(Externally exposed interfaces, re-export structure)

---

## 7. Approach Comparison

### Approach A: {approach name} (recommended)

(description)

| Aspect | Assessment |
|--------|-----------|
| Complexity | (low/medium/high) |
| Extensibility | (assessment) |
| Testability | (assessment) |

### Approach B: {alternative approach}

(description + rejection reason)

### Conclusion: {chosen approach} selected

(Selection rationale summary)

---

## 8. Testing Strategy

### Unit Tests
(Mock-based unit test targets and examples)

### Integration Tests
(Real infrastructure integration test targets and examples)

### Test Matrix

| Layer | Method | Infrastructure |
|-------|--------|---------------|
| (layer) | (test type) | (required infrastructure) |

---

## 9. Implementation Order

(Phase-by-phase implementation order, numbered at file level)

---

## 10. Risk Mitigation

| Risk | Impact | Mitigation |
|------|--------|------------|
| (risk) | (impact) | (mitigation) |
```

## Architecture Template Rules
- [HARD] Overview must include ASCII diagram — visually represent system structure
- [HARD] Core Interfaces must be at implementable code level — no pseudocode
- [HARD] Approach Comparison must compare at least 2 approaches — specify selection rationale
- [HARD] Implementation Order numbered at file level for direct conversion to Plan/Checklist
- [HARD] Testing Strategy must separate Unit/Integration — include test target file paths
- [HARD] Cross-check that all requirements (MUST/SHOULD) from analysis.md are reflected in architecture.md
- Sections may be added/reduced based on project characteristics, but core sections above must be retained
