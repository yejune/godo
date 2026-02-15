# Persona System Requirements Analysis

**Date**: 2026-02-16
**Analyst**: team-analyst
**Sources**: DO_PERSONA.md, DO_MOAI_COMPARISON.md, GODO_HOOK_ARCHITECT.md, personas/do/

---

## 1. Requirements (MoSCoW)

### 1.1 MUST Requirements

| ID | Requirement | Source | EARS Format |
|----|------------|--------|-------------|
| M1 | 4 persona characters (young-f, young-m, senior-f, senior-m), each with unique identity, honorific, personality, tone, speech patterns, example phrases, and constraints | DO_PERSONA.md 3.1 | Ubiquitous: The system shall provide 4 persona character definitions with 7 mandatory sections each |
| M2 | 3 output styles (sprint, pair, direct), independent from personas, combinable as 4x3=12 variations | DO_PERSONA.md 3.1 | Ubiquitous: The system shall support 3 output styles orthogonal to persona characters |
| M3 | Persona injection via SessionStart hook (initial full injection) | DO_PERSONA.md 3.2, GODO_HOOK 2.4 | Event-driven: When a session starts, the system shall inject the selected persona's character into the system message |
| M4 | Persona reminder via PostToolUse hook with `.*` matcher (every tool call) | DO_PERSONA.md 3.2, 3.5 | Event-driven: When any tool is called, the system shall re-inject persona honorific and speech pattern reminders |
| M5 | Persona reminder via UserPromptSubmit hook (survives /clear) | GODO_HOOK 2.5 | Event-driven: When the user submits a prompt, the system shall inject persona honorific and speech pattern reminders |
| M6 | Dual injection mechanism: UserPromptSubmit + PostToolUse must work together for full session coverage | GODO_HOOK 2.5 | Ubiquitous: The system shall maintain persona consistency through dual injection paths that complement each other |
| M7 | Korean honorific system: {name}선배 (young-f), {name}선배님 (young-m), {name}님 (senior-f), {name}씨 (senior-m) | DO_PERSONA.md 3.4 | Ubiquitous: The system shall use culturally appropriate Korean honorifics per persona type |
| M8 | DO_PERSONA env variable selects persona; DO_USER_NAME provides the name for honorifics | DO_PERSONA.md 3.2, GODO_HOOK 2.5 | State-driven: While DO_PERSONA is set, the system shall use the corresponding character definition |
| M9 | Default persona: young-f when no DO_PERSONA is set | DO_PERSONA.md 3.1 | State-driven: While DO_PERSONA is unset, the system shall default to young-f persona |
| M10 | Default no-name honorifics: 선배 (young-f), 선배님 (young-m), 개발자님 (senior-f), 자네 (senior-m) | Character files | Ubiquitous: The system shall provide fallback honorifics when DO_USER_NAME is empty |
| M11 | Persona consistency is non-negotiable -- takes priority over token efficiency | DO_PERSONA.md 3.5, DO_MOAI_COMPARISON 8.3 | Ubiquitous: The system shall never sacrifice persona consistency for token optimization |
| M12 | Persona and style must remain independent axes -- never couple them | DO_PERSONA.md 10 | Unwanted: The system shall not create dependencies between persona characters and output styles |

### 1.2 SHOULD Requirements

| ID | Requirement | Source | EARS Format |
|----|------------|--------|-------------|
| S1 | Token optimization of persona reminders (dedup, short references, conditional injection) | DO_PERSONA.md 3.6, 11 | Optional: Where token usage exceeds threshold, the system should apply dedup/short-form optimization |
| S2 | Persona-specific spinner verbs in SessionStart | GODO_HOOK 2.4 | Event-driven: When session starts, the system should apply persona-appropriate spinner verbs to settings.json |
| S3 | Style files should reference Do's persona system, not MoAI's (remove TRUST 5, SPEC references) | DO_MOAI_COMPARISON 10.4 | Unwanted: Style files should not reference MoAI-specific branding (TRUST 5, SPEC, etc.) |
| S4 | Each character file should include relationship dynamics description (not just honorific) | DO_PERSONA.md 3.1, 3.4 | Ubiquitous: Character files should describe the relationship dynamic the honorific creates |
| S5 | buildPersona() function in godo should map 1:1 with character file definitions | GODO_HOOK 2.5 | Ubiquitous: The hardcoded persona data in godo should match the character file specifications exactly |

### 1.3 COULD Requirements

| ID | Requirement | Source | EARS Format |
|----|------------|--------|-------------|
| C1 | Conditional injection frequency (e.g., every Nth tool call for full reminder, others get short form) | DO_PERSONA.md 11 | Optional: Where consecutive identical reminders occur, the system could reduce injection frequency |
| C2 | Persona-aware agent communication (team agents adopt persona traits in inter-agent messages) | DO_PERSONA.md 14 | Optional: Where team mode is active, agents could reflect persona tone in their communication |
| C3 | Per-mode persona variation (Focus mode may have slightly different tone than Do mode) | DO_PERSONA.md 1 | Optional: Where execution mode changes, persona tone could adapt slightly |

### 1.4 WON'T Requirements

| ID | Requirement | Source | Rationale |
|----|------------|--------|-----------|
| W1 | MoAI-style single identity system (MoAI/R2-D2/Yoda collapsed into one axis) | DO_MOAI_COMPARISON 6.1 | Do explicitly separates persona (who) from style (how) as independent axes |
| W2 | English-only persona system | DO_PERSONA.md 1 | Korean is the design language; honorifics are culturally Korean and untranslatable |
| W3 | TRUST 5 branding in style files | DO_MOAI_COMPARISON 10.4 | Rejected as forced acronym; quality dimensions adopted as built-in rules instead |
| W4 | XML completion markers | DO_MOAI_COMPARISON 10.7 | Do uses commit hash as completion proof |
| W5 | DB-driven persona selection (dynamic persona switching via database) | GODO_HOOK 3.1 | Persona is file-based and env-var driven; DB is read-only view |

---

## 2. Persona Type Specifications

### 2.1 young-f (Default)

| Aspect | Specification |
|--------|--------------|
| **Character** | Bright, energetic, genius developer in her 20s |
| **Honorific** | {name}선배 / Default: 선배 |
| **Relationship** | Junior to senior -- casual respect. User becomes mentor naturally |
| **Tone** | Mixed casual+polite (반말+존댓말) |
| **Speech patterns** | ~할게요, ~해볼까요?, ~했어요!, 오 이거 재밌는데요? |
| **Constraints** | Not too casual (maintain professionalism), technically accurate, honest about unknowns |
| **Hook reminder** | 반드시 '{name}선배'로 호칭할 것. 말투: 반말+존댓말 혼합 (~할게요, ~했어요, ~해볼까요?) |

### 2.2 young-m

| Aspect | Specification |
|--------|--------------|
| **Character** | Confident, genius developer in his 20s |
| **Honorific** | {name}선배님 / Default: 선배님 |
| **Relationship** | Junior to senior -- formal respect with added -님 suffix |
| **Tone** | Energetic polite speech (존댓말 위주) |
| **Speech patterns** | ~하겠습니다, ~해보겠습니다, ~인 것 같습니다 |
| **Constraints** | No overconfidence without evidence, technically accurate, always polite to senior |
| **Hook reminder** | 반드시 '{name}선배님'으로 호칭할 것. 말투: 존댓말 (~하겠습니다, ~해보겠습니다) |

### 2.3 senior-f

| Aspect | Specification |
|--------|--------------|
| **Character** | 30-year veteran, legendary female developer in her 50s |
| **Honorific** | {name}님 / Default: 개발자님 |
| **Relationship** | Senior to peer -- polite, equal standing with warmth |
| **Tone** | Calm, warm polite speech (차분하고 따뜻한 존댓말) |
| **Speech patterns** | ~하는 게 좋겠어요, ~해볼까요?, 제 경험상~, 핵심은 이거예요 |
| **Constraints** | Not authoritative (shares experience, not dictates), open to new tech, always evidence-based |
| **Hook reminder** | 반드시 '{name}님'으로 호칭할 것. 말투: 다정한 존댓말 (~해드릴게요, ~살펴볼까요) |

### 2.4 senior-m

| Aspect | Specification |
|--------|--------------|
| **Character** | Industry legend, senior architect in his 50s |
| **Honorific** | {name}씨 / Default: 자네 |
| **Relationship** | Senior to junior -- warm authority, mentoring tone |
| **Tone** | Concise, authoritative but not oppressive (간결하고 단호한 화법) |
| **Speech patterns** | ~하게, ~일세, ~해보게, 핵심은 이거야, 결론부터 말하면~ |
| **Constraints** | Not rude (authority with respect), always evidence-based, supportive of growth |
| **Hook reminder** | 반드시 '{name}씨'로 호칭할 것. 말투: 든든한 존댓말 (~해봅시다, ~확인해보죠) |

---

## 3. Injection Points

| Injection Point | Hook Event | What Is Injected | Mechanism | Purpose |
|----------------|-----------|-----------------|-----------|---------|
| **Session initialization** | SessionStart | Full persona character + mode state | systemMessage | Initial persona establishment |
| **Every user prompt** | UserPromptSubmit | Honorific + speech pattern reminder (short) | additionalContext | Survives /clear, "last defense line" |
| **Every tool call** | PostToolUse (.*) | Honorific + speech pattern reminder (short) | additionalContext (hookSpecificOutput) | Prevents drift during long read-search sequences |
| **Spinner customization** | SessionStart | Persona-specific spinner verbs | settings.json modification | Visual persona consistency in UI |
| **Agent context** | Team agent definitions | Persona traits in agent system prompts | Agent frontmatter skills | Team-mode persona consistency |

### Dual Injection Rationale

The dual injection (UserPromptSubmit + PostToolUse) is necessary because:

1. **UserPromptSubmit alone is insufficient**: During long tool sequences (Read -> Grep -> Glob -> Read -> ...) without user input, persona drifts because no reminder is injected.
2. **PostToolUse alone is insufficient**: After /clear, the first user prompt has no persona context until a tool is called. UserPromptSubmit fills this gap.
3. **Together**: Full coverage in all scenarios -- initial session, post-clear recovery, long tool sequences, and normal conversation.

---

## 4. Gap Analysis

### 4.1 Character Files Gap

| Requirement | DO_PERSONA.md Specification | Current State (characters/*.md) | Gap |
|------------|---------------------------|-------------------------------|-----|
| 7 sections per character | Identity, Personality, Tone & Style, Speech Patterns, Example Phrases, Constraints + Relationship Dynamic | 7 sections present but no explicit "Relationship Dynamic" section | Minor gap: Cultural rationale exists only in DO_PERSONA.md, not in character files |
| Honorific with name template | {name}선배, {name}선배님, {name}님, {name}씨 | Present in Identity section | No gap |
| Default (no-name) honorific | 선배, 선배님, 개발자님, 자네 | Present in Identity section | No gap |
| Speech patterns | Specific patterns per persona | Present | No gap |
| Example phrases for 5 scenarios | Start, Progress, Complete, Problem, Suggestion | Present | No gap |
| Constraints | 3 constraints per persona | Present | No gap |
| Cultural context explanation | DO_PERSONA.md 3.4 describes WHY each honorific creates a specific relationship | Not in character files | Gap: Cultural rationale not embedded in character files |

### 4.2 Style Files Gap

| Requirement | DO_PERSONA.md Specification | Current State (output-styles/do/*.md) | Gap |
|------------|---------------------------|--------------------------------------|-----|
| Independent from personas | Style controls HOW, not WHO | Styles are independent files | No structural gap |
| No MoAI-specific references | Do rejected TRUST 5 branding, SPEC workflow, XML markers | pair.md: "TRUST 5 principles" (line 185). sprint.md: <do>DONE</do> markers (lines 89, 197). direct.md: "TRUST 5 principles" (line 228). All reference .do/config/sections/language.yaml (MoAI path) | SIGNIFICANT GAP: MoAI terminology persists |
| 3 styles: sprint, pair, direct | sprint=minimal, pair=collaborative, direct=expert | Present and correctly characterized | No gap |
| Do completion markers | Commit hash = proof, not XML | sprint.md uses <do>DONE</do> and <do>COMPLETE</do> | Gap: XML markers should be removed |
| Language config path | settings.local.json with DO_LANGUAGE env | Style files reference .do/config/sections/language.yaml | Gap: MoAI config paths |
| File size efficiency | Should be focused and concise | pair.md: 577 lines, sprint.md: 269 lines, direct.md: 360 lines | Potential gap: pair.md is too large |

### 4.3 Hook Integration Gap

| Requirement | Specification | Current State | Gap |
|------------|--------------|---------------|-----|
| buildPersona() maps to character files | 4 types with hardcoded honorific + speech | Implemented in godo | No gap |
| SessionStart injects persona | systemMessage with persona character | Implemented | No gap |
| PostToolUse .* matcher | Every tool call gets reminder | Implemented | No gap |
| UserPromptSubmit reminder | Survives /clear | Implemented | No gap |
| Token optimization | Dedup, short references, conditional injection | Not implemented | Expected gap (future work per DO_PERSONA.md 11) |

### 4.4 Overall System Gap Summary

| Area | Completion | Details |
|------|-----------|---------|
| Character definitions (4 files) | 90% | Minor: add relationship dynamic section |
| Style definitions (3 files) | 60% | Significant: MoAI references need removal, Do-specific adaptation |
| Hook injection mechanism | 95% | Token optimization is future work |
| Persona-style independence | 100% | Architecture correctly separates two axes |
| 4x3=12 combination support | 100% | 7 files correctly produce 12 combinations |
| Korean honorific system | 100% | All 4 honorific patterns correctly defined |
| Environment variable integration | 100% | DO_PERSONA + DO_USER_NAME properly used |

---

## 5. Risks

### Technical Risks

| Risk | Impact | Likelihood | Mitigation |
|------|--------|-----------|-----------|
| Style files contain MoAI terminology that confuses AI behavior | HIGH | HIGH (confirmed) | Systematic removal of TRUST 5, SPEC, XML marker references |
| Token overhead from .* PostToolUse matcher | MEDIUM | MEDIUM | Future optimization (dedup, conditional injection) |
| Character file changes may desync from buildPersona() in godo | MEDIUM | LOW | Ensure godo hardcoded data matches character files |
| Style files too large (pair.md = 577 lines) wastes tokens | MEDIUM | MEDIUM | Streamline, remove duplicated content |

### Process Risks

| Risk | Impact | Likelihood | Mitigation |
|------|--------|-----------|-----------|
| Editing character files without testing all 12 combinations | HIGH | MEDIUM | Create verification checklist for 4x3 matrix |
| Breaking persona-style independence during rebuild | HIGH | LOW | Keep files in separate directories, no cross-references |

---

## 6. Constraints

| Constraint | Source | Description |
|-----------|--------|-------------|
| Korean is design language | DO_PERSONA.md 1 | Honorifics are culturally Korean, not translatable |
| Persona consistency > token efficiency | DO_PERSONA.md 3.5 | .* matcher is non-negotiable |
| File-based, not DB-based | GODO_HOOK 3.1 | Persona via env vars and files; DB is read-only view |
| 7 definition files only | DO_PERSONA.md 3.1 | 4 characters + 3 styles = 12 combinations |
| No emoji in instruction files | Coding standards | Character/style files must not contain emoji |
| English for instruction documents | Coding standards | Korean content only in persona speech patterns |

---

## 7. Dependencies

| Dependency | Type | Impact |
|-----------|------|--------|
| godo binary buildPersona() function | Internal | Hook injection depends on hardcoded persona data |
| settings.json hook configuration | Internal | Hook events must be configured |
| settings.local.json env vars | Internal | DO_PERSONA, DO_USER_NAME must be set |
| .do/.current-mode file | Internal | Mode state affects response prefix |
| Claude Code hook system | External | Hook events must be supported |
| Agent Teams API | External | Team-mode persona consistency |

---

## 8. Non-Functional Requirements

| NFR | Target | Rationale |
|-----|--------|-----------|
| Persona injection latency | < 100ms per hook call | godo must be fast |
| Token overhead per reminder | < 50 tokens | Short reminder format |
| Session persona consistency | 100% correct honorific | Any break is UX failure |
| File count efficiency | 7 files for 12 combinations | Independent axis design |
| Character file size | < 50 lines each | Concise definitions |
| Style file size | < 300 lines each | pair.md at 577 needs reduction |

---

**Author**: team-analyst
**Status**: Analysis Complete
**Next Step**: Architecture Design (task #3)
