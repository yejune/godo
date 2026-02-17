# Analysis Document Template [HARD]

For complex tasks, write to `.do/jobs/{YY}/{MM}/{DD}/{title-kebab-case}/analysis.md`.
Refer to complexity-check.md for complexity criteria.

```markdown
# {Task Title} Analysis

**Analysis Date**: {YYYY-MM-DD}
**Target Project**: {project name}
**Review Scope**: {analysis target scope}

---

## 1. Current System Analysis

### 1.1 Components

| Item | Value |
|------|-------|
| **Core Modules** | (module/service names) |
| **Tech Stack** | (languages, frameworks, libraries) |
| **Data Storage** | (DB, cache, filesystem, etc.) |
| **External Dependencies** | (third-party APIs, services, etc.) |

### 1.2 Data Flow / Core Logic

(Data flow, state machines, core business logic of current system with code snippets)

### 1.3 Current Configuration / Infrastructure

(docker-compose, environment variables, connection settings, infrastructure details)

---

## 2. Change Scope Analysis

### 2.1 Areas Requiring Change

(Current code → post-change code comparison per area)

### 2.2 What Must Be Preserved

(Business logic, data structures, contracts that must be preserved regardless of change scope)

### 2.3 Current Features Unsupported in Target Technology

| Feature | Current | Target | Migration Method |
|---------|---------|--------|-----------------|
| (feature name) | (current support) | (target support) | (alternative or custom implementation) |

---

## 3. Requirements Summary

### 3.1 Must Have (MUST)

| Feature | Description | Priority |
|---------|-------------|----------|
| (feature name) | (description) | P0/P1 |

### 3.2 Should Have (SHOULD)

| Feature | Description | Alternative |
|---------|-------------|-------------|
| (feature name) | (description) | (alternative if any) |

### 3.3 Could Have (COULD)

| Feature | Description | Exclusion Reason |
|---------|-------------|-----------------|
| (feature name) | (description) | (why outside v1 scope) |

### 3.4 Won't Have (WON'T)

| Feature | Reason |
|---------|--------|
| (feature name) | (exclusion reason) |

---

## 4. Technology Options Comparison

### 4.1 Candidates

| Candidate | Type | Pros | Cons |
|-----------|------|------|------|
| (Candidate A) | (type) | (pros) | (cons) |
| (Candidate B) | (type) | (pros) | (cons) |

### 4.2 Recommended Choice: {chosen candidate}

**Reasons:**
- (rationale 1)
- (rationale 2)

---

## 5. Migration Strategy

### 5.1 Phased Transition

| Phase | Content | Impact Scope |
|-------|---------|-------------|
| Phase 1 | (content) | (impact scope) |
| Phase 2 | (content) | (impact scope) |

### 5.2 Risks and Mitigations

| Risk | Impact | Mitigation |
|------|--------|------------|
| (risk description) | HIGH/MEDIUM/LOW | (mitigation method) |

---

## 6. Conclusion

### Key Changes Summary
1. (change 1)
2. (change 2)

### Recommended Implementation Strategy
1. (strategy 1)
2. (strategy 2)

---

**Author**: {agent name}
**Review Status**: Analysis complete
**Next Step**: Architecture design
```

## Analysis Template Rules
- [HARD] Section 1 (Current System Analysis) must be based on code reverse-engineering — extract from actual code, not guesses
- [HARD] Section 3 (Requirements) must use MoSCoW classification (MUST/SHOULD/COULD/WON'T)
- [HARD] Section 4 (Technology Options) must compare at least 2 candidates
- [HARD] Section 5 (Risks) must specify impact level (HIGH/MEDIUM/LOW)
- [HARD] All code snippets must be excerpted from actual project code — no hypothetical code
- Sections may be added/reduced based on project characteristics, but core sections above must be retained
