# Persona System Architecture Design

## Overview

The persona system provides AI character identity (WHO speaks) and output style (HOW they speak) as two independent axes, producing 4 x 3 = 12 behavioral combinations from 7 definition files. The current implementation hardcodes persona data in Go source (`buildPersona()`, `spinnerStemsYoungF`), making it impossible to add/modify personas without recompiling godo. This redesign externalizes all persona content into markdown files under `personas/do/`, enabling file-based persona management while preserving the existing hook injection mechanism.

```
                         ┌─────────────────────────────────┐
                         │     settings.local.json          │
                         │  DO_PERSONA=young-f              │
                         │  DO_USER_NAME=승민               │
                         └──────────┬──────────────────────┘
                                    │
                         ┌──────────▼──────────────────────┐
                         │         godo binary              │
                         │                                  │
                         │  ┌─────────────────────────┐    │
                         │  │  PersonaLoader           │    │
                         │  │  (reads personas/do/)    │    │
                         │  └────────┬────────────────┘    │
                         │           │                      │
                         │  ┌────────▼────────────────┐    │
                         │  │  PersonaRenderer         │    │
                         │  │  (templates + variables) │    │
                         │  └────────┬────────────────┘    │
                         │           │                      │
                         └───────────┼──────────────────────┘
                                     │
               ┌─────────────────────┼──────────────────────┐
               │                     │                      │
    ┌──────────▼──────┐   ┌──────────▼──────┐   ┌──────────▼──────┐
    │  SessionStart   │   │ UserPromptSubmit │   │  PostToolUse    │
    │  hook           │   │  hook            │   │  hook (.*)      │
    │                 │   │                  │   │                 │
    │ spinner verbs   │   │ honorific +      │   │ honorific +     │
    │ apply to        │   │ tone reminder    │   │ tone reminder   │
    │ settings.json   │   │ (survives /clear)│   │ (every tool)    │
    └─────────────────┘   └──────────────────┘   └─────────────────┘
```

---

## 1. Directory Structure

```
personas/do/
├── persona.yaml                    # Manifest: lists all persona assets
├── CLAUDE.md                       # Do's top-level CLAUDE.md (persona overlay)
│
├── characters/                     # WHO speaks (4 files)
│   ├── young-f.md                  # 밝은 20대 여성 천재 개발자
│   ├── young-m.md                  # 자신감 넘치는 20대 남성 천재 개발자
│   ├── senior-f.md                 # 30년 경력 레전드 50대 여성 개발자
│   └── senior-m.md                 # 업계 전설 50대 남성 시니어 아키텍트
│
├── styles/                         # HOW they speak (3 files)
│   ├── sprint.md                   # 말 최소화, 즉시 실행
│   ├── pair.md                     # 협업적 톤, 공동 의사결정
│   └── direct.md                   # 군더더기 없이 필요한 것만
│
├── spinners/                       # Loading animation verbs (per character)
│   ├── young-f.yaml                # Playful stems + emoji
│   ├── young-m.yaml                # Formal stems
│   ├── senior-f.yaml               # Professional stems
│   └── senior-m.yaml               # Authoritative stems
│
├── agents/                         # Do persona-specific agents
│   └── (existing agent .md files)
├── skills/                         # Do persona-specific skills
│   └── (existing skill .md files)
├── rules/                          # Do persona-specific rules
│   └── (existing rule .md files)
├── commands/                       # Do persona-specific commands
│   └── do/
│       └── style.md
├── hooks/                          # Do persona hook config
│   └── (hook settings overrides)
└── settings.json                   # Do persona settings overrides
```

---

## 2. Core Interfaces / Data Formats

### 2.1 Character File Format (`characters/*.md`)

Each character file uses YAML frontmatter for machine-readable fields and markdown body for the full system prompt injection content.

```markdown
---
id: young-f
name: "Do"
honorific_template: "{{name}}선배"
honorific_default: "선배"
tone: "반말+존댓말 혼합 (~할게요, ~했어요, ~해볼까요?)"
character_summary: "밝고 에너지 넘치는 20대 여성 천재 개발자"
relationship: "후배가 선배에게 캐주얼한 존중"
---

# Character: young-f

## Identity

밝고 에너지 넘치는 20대 여성 천재 개발자. 복잡한 문제도 가볍게 풀어냄.
어떤 기술 스택이든 자유자재. 에너지 넘치고 적극적인 캐릭터.

## Speech Patterns

### Do Mode Examples
- "바로 시작할게요!"
- "병렬로 돌릴게요!"
- "이 정도는 한방에!"

### Focus Mode Examples
- "차근차근 해볼게요"
- "하나씩 확인할게요"
- "여기 원인 보이네요!"

### Team Mode Examples
- "팀 구성할게요!"
- "같이 해봐요!"
- "다 같이 달려볼까요?"

## Behavioral Guidelines

(Additional character-specific behavioral notes for longer system prompt injection.
This section is loaded only for SessionStart full injection, not for per-tool reminders.)
```

**Frontmatter fields (machine-readable):**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | Yes | Persona identifier matching DO_PERSONA value |
| `name` | string | Yes | Display name for the AI character |
| `honorific_template` | string | Yes | Template with `{{name}}` placeholder |
| `honorific_default` | string | Yes | Fallback when DO_USER_NAME is empty |
| `tone` | string | Yes | One-line tone description for reminder injection |
| `character_summary` | string | Yes | One-line character description |
| `relationship` | string | No | Describes the user-AI relationship dynamic |

**Template variable:** `{{name}}` is replaced with `DO_USER_NAME` at runtime.

### 2.2 Style File Format (`styles/*.md`)

Style files define HOW the AI responds, independent of character.

```markdown
---
id: pair
name: "Pair Programming"
description: "협업적 톤, 공동 의사결정"
---

# Style: pair (Pair Programming)

## Response Principles

- 친절한 동료처럼 대화
- 공동 의사결정 (제안 후 확인)
- 적절한 설명과 맥락 제공

## Format Rules

- 코드 변경 전 접근 방식 설명
- 대안이 있으면 옵션 제시
- 작업 완료 후 요약 제공

## Tone Markers

- "이렇게 해볼까요?"
- "두 가지 방법이 있는데..."
- "제 생각엔 이게 나을 것 같아요"
```

### 2.3 Spinner File Format (`spinners/*.yaml`)

```yaml
# spinners/young-f.yaml
persona: young-f
suffix_pattern:
  cycle: 3                    # Rotate through suffixes every N entries
  suffixes:
    - "중이에요!"
    - "중이에요 선배!"
    - "중!"
stems:
  - stem: "거의 다 해내는"
    emoji: ""
  - stem: "갓 구운 답변 준비하는"
    emoji: "🍞"
  - stem: "향긋하게 우려내는"
    emoji: "☕"
  - stem: "머리 풀가동으로 계산하는"
    emoji: "🧮"
  # ... (all existing stems from spinner.go)
```

```yaml
# spinners/senior-m.yaml
persona: senior-m
suffix_pattern:
  cycle: 1                    # Always the same suffix
  suffixes:
    - "중입니다."
stems:
  - stem: "해내는"
  - stem: "실행하는"
  - stem: "실현하는"
  # ... (all existing stems from spinnerStemsDefault)
```

### 2.4 Persona Manifest (`persona.yaml`)

Updated to reference character/style/spinner files:

```yaml
name: do
version: "3.0.0"
description: "Do persona - execution-first adaptive orchestrator"

brand: do
brand_dir: ".do"
brand_cmd: "/do"

claude_md: CLAUDE.md

# Character files (4 personas)
characters:
  - characters/young-f.md
  - characters/young-m.md
  - characters/senior-f.md
  - characters/senior-m.md

# Style files (3 styles)
styles:
  - styles/sprint.md
  - styles/pair.md
  - styles/direct.md

# Spinner verb definitions
spinners:
  - spinners/young-f.yaml
  - spinners/young-m.yaml
  - spinners/senior-f.yaml
  - spinners/senior-m.yaml

# Existing asset references
agents:
  - agents/manager-ddd.md
  # ... existing agent files
skills: []
rules: []
commands:
  - commands/do/style.md

# Hook configuration
hooks:
  SessionStart:
    - command: "godo hook session-start"
      timeout: 10000
  UserPromptSubmit:
    - command: "godo hook user-prompt-submit"
      timeout: 5000
  PreToolUse:
    - command: "godo hook pre-tool"
      timeout: 5000
      matcher: "Write|Edit|Bash"
  PostToolUse:
    - command: "godo hook post-tool-use"
      timeout: 5000
      matcher: ".*"

# Settings overrides
settings:
  outputStyle: pair

# Installer for godo binary
installer:
  binary: godo
  init_cmd: "brew install yejune/tap/godo"

# Slot content (for core template filling)
slot_content:
  QUALITY_FRAMEWORK: content/quality-framework.md

# Agent patches (modify core agents for Do persona)
agent_patches:
  expert-backend:
    append_skills:
      - do-foundation-checklist
    remove_skills:
      - moai-foundation-quality
```

### 2.5 Go Interface: PersonaLoader

```go
// PersonaData holds parsed character data from a character .md file.
type PersonaData struct {
    ID                string   // "young-f", "young-m", etc.
    Name              string   // Display name
    HonorificTemplate string   // "{{name}}선배"
    HonorificDefault  string   // "선배"
    Tone              string   // One-line tone description
    CharacterSummary  string   // One-line character description
    Relationship      string   // Relationship dynamic
    FullContent       string   // Full markdown body (for full injection)
}

// SpinnerData holds parsed spinner configuration.
type SpinnerData struct {
    Persona       string
    SuffixCycle   int
    Suffixes      []string
    Stems         []SpinnerStem
}

type SpinnerStem struct {
    Stem  string
    Emoji string
}

// LoadCharacter reads and parses a character .md file.
// Returns parsed PersonaData with frontmatter fields and body content.
func LoadCharacter(personaDir, personaType string) (*PersonaData, error)

// LoadSpinner reads and parses a spinner .yaml file.
func LoadSpinner(personaDir, personaType string) (*SpinnerData, error)

// BuildHonorific applies the user name to the honorific template.
// If userName is empty, returns HonorificDefault.
func (p *PersonaData) BuildHonorific(userName string) string

// BuildReminder produces the short one-line reminder for PostToolUse/UserPromptSubmit.
// Format: 반드시 "{honorific}"로 호칭할 것. 말투: {tone}
func (p *PersonaData) BuildReminder(userName string) string

// BuildSpinnerVerbs generates the full spinner verb list from SpinnerData.
func (s *SpinnerData) BuildSpinnerVerbs() []string
```

---

## 3. Data Flow Diagram

### 3.1 Persona Resolution Flow

```
┌──────────────────────────────────────────────────────────┐
│                    CONFIGURATION                         │
│                                                          │
│  .claude/settings.local.json                             │
│    DO_PERSONA = "young-f"                                │
│    DO_USER_NAME = "승민"                                  │
│    DO_STYLE = "pair"   (optional, default in manifest)   │
│                                                          │
└──────────────────┬───────────────────────────────────────┘
                   │ env vars
                   ▼
┌──────────────────────────────────────────────────────────┐
│                     godo binary                          │
│                                                          │
│  1. Read DO_PERSONA env var → "young-f"                  │
│  2. Resolve persona dir:                                 │
│     {project}/.claude/personas/do/    (assembled output) │
│     OR fallback to hardcoded defaults                    │
│  3. Load characters/young-f.md                           │
│     → Parse YAML frontmatter                             │
│     → Extract: honorific_template, tone, character_summary│
│  4. Load spinners/young-f.yaml                           │
│     → Parse stems + suffix pattern                       │
│  5. Apply DO_USER_NAME to honorific_template             │
│     "{{name}}선배" → "승민선배"                             │
│                                                          │
└──────────────────┬───────────────────────────────────────┘
                   │
         ┌─────────┼─────────────┐
         │         │             │
         ▼         ▼             ▼
   SessionStart  UserPrompt   PostToolUse
                 Submit
```

### 3.2 Per-Hook Injection Detail

```
SessionStart:
  ┌──────────────────────────────────────────────┐
  │ 1. Load spinner YAML → build verb list       │
  │ 2. Apply verbs to ~/.claude/settings.json    │
  │ 3. Read mode state (.do/.current-mode)       │
  │ 4. Output: systemMessage with mode + project │
  │    (persona reminders deferred to other hooks)│
  └──────────────────────────────────────────────┘

UserPromptSubmit:
  ┌──────────────────────────────────────────────┐
  │ 1. Load character frontmatter (cached)       │
  │ 2. Build honorific: "승민선배"                 │
  │ 3. Build reminder line:                      │
  │    반드시 "승민선배"로 호칭할 것.               │
  │    말투: 반말+존댓말 혼합                       │
  │ 4. Append mode reminder                      │
  │ 5. Output: additionalContext                 │
  └──────────────────────────────────────────────┘

PostToolUse (.*):
  ┌──────────────────────────────────────────────┐
  │ 1. Load character frontmatter (cached)       │
  │ 2. Build same reminder line as above         │
  │ 3. Append checklist/plan pipeline results    │
  │ 4. Append lint results (if code file)        │
  │ 5. Output: additionalContext                 │
  └──────────────────────────────────────────────┘
```

### 3.3 Caching Strategy

godo is a short-lived process (one invocation per hook event). File parsing on every invocation is acceptable because:

- Character frontmatter is ~20 lines of YAML -- parse time < 1ms
- Spinner YAML is ~60 entries -- parse time < 1ms
- godo already reads `.do/.current-mode`, `.do/.latest-version`, and walks `.do/jobs/` on every PostToolUse
- No in-memory cache needed; filesystem IS the cache

If profiling shows measurable overhead, a future optimization is to serialize parsed persona data to `.do/cache/persona.json` and use file mtime for invalidation.

---

## 4. Approach Comparison

### Approach A: File-Based Externalization

Move all hardcoded persona data from Go source into markdown/YAML files under `personas/do/`. godo reads these files at runtime.

| Criterion | Evaluation |
|-----------|------------|
| Alignment with existing patterns | High -- matches persona.yaml manifest pattern already used by convert tool |
| Complexity | Low -- simple file read + YAML/frontmatter parse |
| Maintainability | High -- edit markdown files, no recompilation |
| Performance | Negligible impact -- files are tiny, godo already does file I/O |
| Testing strategy | Unit test file parsing; integration test hook output |
| Migration impact | Moderate -- extract hardcoded data to files, update godo to read from files with fallback |
| Extensibility | New persona = new .md file, no code change |

**Pros:**
- Consistent with convert tool's extraction/assembly model
- Persona content is version-controlled alongside other project assets
- Non-developers can modify personas by editing markdown
- The assembler can copy persona files to output `.claude/` directory
- Fallback to hardcoded defaults preserves backward compatibility

**Cons:**
- Adds file I/O to every hook invocation (but files are tiny)
- Need to handle "file not found" gracefully
- Two sources of truth during migration (files + hardcoded fallback)

### Approach B: Embedded Go Templates with Config

Keep persona data in Go source but use Go `embed` directive to bundle markdown files at compile time. Persona selection still reads `DO_PERSONA` env var, but content is baked into the binary.

| Criterion | Evaluation |
|-----------|------------|
| Alignment with existing patterns | Medium -- Go embed is standard but doesn't match convert's file-based model |
| Complexity | Low -- `//go:embed` is straightforward |
| Maintainability | Medium -- must recompile to change personas |
| Performance | Best -- no file I/O at runtime |
| Testing strategy | Same as Approach A |
| Migration impact | Low -- restructure Go code, embed files |
| Extensibility | Low -- adding persona requires code change + rebuild |

**Pros:**
- Zero runtime file dependency -- single binary, always works
- No "file not found" edge cases
- Fastest possible execution

**Cons:**
- Defeats the purpose: cannot modify personas without recompiling godo
- Breaks the convert tool's extraction/assembly workflow (assembler outputs files, but godo ignores them)
- Cannot add custom personas in a project without forking godo

### Approach C: Hybrid (File-First with Embedded Fallback) -- RECOMMENDED

Read from `personas/do/` files at runtime. If the file is not found, fall back to Go-embedded defaults. This is Approach A + B combined.

| Criterion | Evaluation |
|-----------|------------|
| Alignment with existing patterns | Highest -- file-first matches convert, fallback ensures reliability |
| Complexity | Medium -- two code paths (file load + embedded fallback) |
| Maintainability | High -- edit files when available, safe fallback |
| Performance | Good -- file read with fast fallback |
| Migration impact | Lowest -- gradual migration, no breaking changes |
| Extensibility | High -- file-based for customization, embedded for defaults |

**Pros:**
- Best of both worlds: flexibility + reliability
- Smooth migration path -- hardcoded defaults work until files are created
- Projects without persona files still get default behavior
- Aligns with convert tool's output model

**Cons:**
- More code to maintain (two paths)
- Risk of embedded defaults going stale if not synced with files

### Decision: Approach C (Hybrid)

**Rationale:** The hybrid approach provides the smoothest migration path. The current hardcoded `buildPersona()` becomes the fallback, and file-based loading is the primary path. This means:

1. Existing godo installations continue working without any persona files
2. The convert tool's assembler outputs persona files to `personas/do/`
3. godo prefers files when present, falls back to embedded defaults
4. Adding a 5th persona = adding a new .md file (no code change)

---

## 5. Integration with Convert Tool

### 5.1 Extraction Phase

The existing `ExtractorOrchestrator` already handles styles as fully persona-specific. Character and spinner files are new asset types that need extraction support.

Changes to convert tool:

| File | Change | Description |
|------|--------|-------------|
| `internal/model/persona_manifest.go` | Add fields | Add `Characters []string` and `Spinners []string` to `PersonaManifest` |
| `internal/extractor/character.go` | New file | `CharacterExtractor` -- classifies character .md files as persona |
| `internal/extractor/spinner.go` | New file | `SpinnerExtractor` -- classifies spinner .yaml files as persona |
| `internal/extractor/orchestrator.go` | Update routing | Route `characters/*.md` and `spinners/*.yaml` to new extractors |
| `internal/assembler/orchestrator.go` | Update copy | Include characters and spinners in persona file copy |

### 5.2 Assembly Phase

During assembly, the assembler copies character/spinner files from persona source to output:

```
personas/do/characters/young-f.md  →  output/.claude/personas/do/characters/young-f.md
personas/do/spinners/young-f.yaml  →  output/.claude/personas/do/spinners/young-f.yaml
```

### 5.3 godo Persona File Resolution

godo resolves persona files in this order:

```
1. {CLAUDE_PROJECT_DIR}/.claude/personas/do/characters/{type}.md    (assembled output)
2. {CLAUDE_PROJECT_DIR}/personas/do/characters/{type}.md            (source, for dev)
3. Fallback to hardcoded buildPersona() defaults                    (backward compat)
```

---

## 6. Go Code Changes in godo

### 6.1 New File: `cmd/godo/persona_loader.go`

Contains `PersonaData`, `SpinnerData` structs and `LoadCharacter()`, `LoadSpinner()` functions. Parses YAML frontmatter from character .md files and full YAML from spinner files.

### 6.2 Modified: `cmd/godo/hook_session_start.go`

```
Current:  buildPersona(personaType, userName)  → hardcoded Persona struct
          getPersonaSpinnerVerbs(personaType)   → hardcoded spinner arrays

New:      LoadCharacter(personaDir, personaType) → file-based PersonaData
          LoadSpinner(personaDir, personaType)   → file-based SpinnerData
          Fallback to buildPersona() if file not found
```

### 6.3 Modified: `cmd/godo/hook_post_tool_use.go`

```
Current:  buildPersona(persona, userName)
          p.Honorific + p.Tone → additionalContext

New:      LoadCharacter(personaDir, persona)
          pd.BuildReminder(userName) → additionalContext
          Fallback to buildPersona() if file not found
```

### 6.4 Modified: `cmd/godo/hook_user_prompt_submit.go`

Same pattern as PostToolUse -- replace `buildPersona()` call with `LoadCharacter()` + fallback.

### 6.5 Modified: `cmd/godo/spinner.go`

```
Current:  spinnerStemsYoungF (hardcoded array)
          spinnerStemsDefault (hardcoded array)
          getPersonaSpinnerVerbs() → switch on persona type

New:      LoadSpinner(personaDir, personaType) → file-based SpinnerData
          sd.BuildSpinnerVerbs() → generated verb list
          Fallback to hardcoded arrays if file not found
```

### 6.6 Unchanged Files

- `cmd/godo/hook_pre_tool.go` -- security checks, no persona involvement
- `cmd/godo/hook_stop.go` -- checklist checks, no persona involvement
- `cmd/godo/hook_compact.go` -- context compression, no persona involvement
- `cmd/godo/hook_subagent_stop.go` -- progress tracking, no persona involvement
- `cmd/godo/hook_session_end.go` -- session cleanup, no persona involvement
- `cmd/godo/mode.go` -- mode management, no persona involvement
- `cmd/godo/statusline.go` -- status display, no persona involvement

---

## 7. Testing Strategy

### Unit Tests

| Test Target | File | Method |
|-------------|------|--------|
| Character frontmatter parsing | `persona_loader_test.go` | Parse sample .md, verify all fields extracted |
| Spinner YAML parsing | `persona_loader_test.go` | Parse sample .yaml, verify stems + suffixes |
| Honorific template rendering | `persona_loader_test.go` | Test `{{name}}` replacement with various inputs |
| Reminder string building | `persona_loader_test.go` | Test `BuildReminder()` output format |
| Spinner verb generation | `persona_loader_test.go` | Test `BuildSpinnerVerbs()` matches expected output |
| Fallback behavior | `persona_loader_test.go` | Test file-not-found falls back to hardcoded |
| File resolution order | `persona_loader_test.go` | Test priority: .claude/personas > personas > fallback |

### Integration Tests

| Test Target | File | Method |
|-------------|------|--------|
| Hook output with file-based persona | `hook_session_start_test.go` | Verify systemMessage content |
| Hook output with fallback persona | `hook_post_tool_use_test.go` | Verify additionalContext with no persona files |
| Spinner application | `spinner_test.go` | Verify settings.json gets correct verbs from YAML |
| End-to-end persona switch | `e2e_test.go` | Change DO_PERSONA, verify all hooks output correct character |

### Test Matrix

| Layer | Method | Infrastructure |
|-------|--------|----------------|
| PersonaLoader parsing | Unit test | None (testdata fixtures) |
| Hook persona injection | Unit test | None (mock stdin/stdout) |
| Spinner verb generation | Unit test | None |
| File resolution fallback | Unit test | Temp directories |
| Full hook pipeline | Integration test | Docker (godo binary) |

---

## 8. Implementation Order

### Phase 1: Create Persona Files (no code changes)

1. Create `personas/do/characters/young-f.md` -- extract from `buildPersona()` case "young-f"
2. Create `personas/do/characters/young-m.md` -- extract from `buildPersona()` case "young-m"
3. Create `personas/do/characters/senior-f.md` -- extract from `buildPersona()` case "senior-f"
4. Create `personas/do/characters/senior-m.md` -- extract from `buildPersona()` case "senior-m"
5. Create `personas/do/spinners/young-f.yaml` -- extract from `spinnerStemsYoungF`
6. Create `personas/do/spinners/young-m.yaml` -- extract from `spinnerStemsDefault` + young-m suffix
7. Create `personas/do/spinners/senior-f.yaml` -- extract from `spinnerStemsDefault` + senior-f suffix
8. Create `personas/do/spinners/senior-m.yaml` -- extract from `spinnerStemsDefault` + senior-m suffix
9. Create `personas/do/styles/sprint.md` -- extract from existing do-focus styles
10. Create `personas/do/styles/pair.md` -- extract from existing do-focus styles
11. Create `personas/do/styles/direct.md` -- extract from existing do-focus styles
12. Update `personas/do/persona.yaml` -- add characters, spinners, styles references

### Phase 2: Implement PersonaLoader in godo

13. Create `cmd/godo/persona_loader.go` -- PersonaData, SpinnerData, LoadCharacter(), LoadSpinner()
14. Create `cmd/godo/persona_loader_test.go` -- unit tests for parsing + rendering
15. Add test fixtures in `cmd/godo/testdata/personas/`

### Phase 3: Wire PersonaLoader into Hooks

16. Modify `cmd/godo/hook_session_start.go` -- use LoadCharacter + LoadSpinner with fallback
17. Modify `cmd/godo/hook_post_tool_use.go` -- use LoadCharacter with fallback
18. Modify `cmd/godo/hook_user_prompt_submit.go` -- use LoadCharacter with fallback
19. Modify `cmd/godo/spinner.go` -- use LoadSpinner with fallback

### Phase 4: Update Convert Tool

20. Add `Characters` and `Spinners` fields to `PersonaManifest` in convert tool
21. Create `CharacterExtractor` and `SpinnerExtractor` in convert tool
22. Update `ExtractorOrchestrator` routing for new file types
23. Update `Assembler` to copy character/spinner files

### Phase 5: Verification

24. Run all existing tests -- ensure no regression
25. Integration test: full hook pipeline with file-based personas
26. Manual verification: switch DO_PERSONA, confirm all 4 characters work
27. Manual verification: switch DO_STYLE, confirm all 3 styles work

---

## 9. Risk Mitigation

| Risk | Impact | Likelihood | Mitigation |
|------|--------|------------|------------|
| Persona file not found at runtime | Medium -- persona reverts to neutral | Low (fallback exists) | Hardcoded `buildPersona()` remains as fallback; warn on stderr |
| YAML frontmatter parse error | Medium -- persona injection fails | Low (simple format) | Validate at build/test time; fallback on parse error |
| File I/O adds latency to hook execution | Low -- hooks must be fast | Very Low (files < 1KB) | Measured at < 1ms for frontmatter parse; acceptable |
| Persona files diverge from hardcoded defaults | Low -- confusing behavior | Medium (during migration) | Phase 5 verification; deprecation warning for hardcoded path |
| Convert tool doesn't handle new file types | Medium -- assembly fails | Low (clear implementation path) | Phase 4 adds support before files are required |
| Spinner YAML format breaks existing behavior | High -- user-visible regression | Low (test coverage) | Unit test verifies generated verbs match current hardcoded output exactly |

---

## 10. Migration Strategy

### Backward Compatibility

The hybrid approach (Approach C) ensures zero breaking changes:

1. **No persona files** -- godo falls back to `buildPersona()` and hardcoded spinners (current behavior)
2. **Partial persona files** -- godo uses available files, falls back for missing ones
3. **Full persona files** -- godo uses all files, hardcoded defaults are dormant

### Deprecation Path

After Phase 5 verification confirms file-based loading works correctly:

1. Add `// Deprecated: use LoadCharacter() instead` to `buildPersona()`
2. Add stderr warning when fallback is used: `"[godo] persona file not found, using built-in defaults"`
3. In a future major version, remove hardcoded defaults entirely

### Data Extraction

The hardcoded data in `hook_session_start.go` and `spinner.go` is the canonical source for creating the initial persona files. Phase 1 extracts this data verbatim -- no modifications, no enhancements. The goal is exact behavioral parity between file-based and hardcoded paths.

---

**Author**: team-architect
**Date**: 2026-02-16
**Status**: Design complete
**Next Step**: Plan creation based on this architecture
