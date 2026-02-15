# Persona System Rebuild Plan

**Date**: 2026-02-16
**Based on**: architecture-persona-system.md, analysis-requirements.md, DO_PERSONA.md

---

## Problem Statement

Do의 페르소나 시스템에 세 가지 핵심 문제가 있다:

1. **characters/*.md는 dead code**: 4개 캐릭터 파일이 존재하지만, godo의 `buildPersona()` 함수가 하드코딩된 데이터를 사용하므로 파일은 참조되지 않음
2. **스타일 파일에 MoAI 오염**: `output-styles/do/` 아래 3개 스타일 파일에 TRUST 5, SPEC, XML 마커, MoAI 설정 경로 등 MoAI 잔재가 남아있음
3. **주입이 1줄**: SessionStart hook에서 페르소나 주입이 최소한의 하드코딩된 문자열 — 캐릭터 파일의 풍부한 정의가 활용되지 않음

### 목표

Architecture 문서의 Approach C (Hybrid: File-First with Embedded Fallback)를 구현하여:
- 캐릭터/스피너 파일을 machine-readable YAML frontmatter 형식으로 재구성
- PersonaLoader를 godo에 구현하여 파일 기반 로딩 + 하드코딩 폴백
- 스타일 파일에서 MoAI 잔재를 제거하고 Do 정체성으로 재작성
- convert 도구가 새 파일 타입(characters, spinners)을 인식하도록 업데이트

---

## 9-Phase Implementation

### Phase 1: 페르소나 파일 생성 (하드코딩에서 추출)

**목적**: godo의 `buildPersona()`와 `spinnerStemsYoungF/Default`에서 데이터를 추출하여 파일로 외부화

**작업 내용**:
- 4개 캐릭터 파일을 YAML frontmatter 형식으로 재작성 (id, honorific_template, tone 등)
- 4개 스피너 YAML 파일 신규 생성 (stems + suffix_pattern)
- `personas/do/styles/` 디렉토리에 3개 스타일 파일 생성 (output-styles/do/에서 이동)
- persona.yaml 매니페스트에 characters, spinners, styles 참조 추가

**변경 파일**: 12개 (characters 4 + spinners 4 + styles 3 + manifest 1)
**의존성**: 없음 (코드 변경 없이 파일만 생성)
**검증**: 각 파일의 YAML frontmatter 파싱 가능 여부 확인, persona.yaml 구조 검증

### Phase 2: PersonaLoader 구현

**목적**: godo에 파일 기반 페르소나 로딩 기능 추가

**작업 내용**:
- `cmd/godo/persona_loader.go` 신규 생성: PersonaData, SpinnerData 구조체 + LoadCharacter(), LoadSpinner() 함수
- `cmd/godo/persona_loader_test.go` 신규 생성: 파싱, 렌더링, 폴백 단위 테스트
- `cmd/godo/testdata/personas/` 테스트 픽스처 생성

**변경 파일**: 3개 (loader + test + testdata)
**의존성**: Phase 1 완료 (파일 형식이 확정되어야 파서 작성 가능)
**검증**: 단위 테스트 전체 통과

### Phase 3: 훅에 PersonaLoader 연결

**목적**: 기존 하드코딩된 `buildPersona()` 호출을 PersonaLoader로 교체 (폴백 유지)

**작업 내용**:
- `hook_session_start.go` 수정: LoadCharacter + LoadSpinner 사용, 폴백 유지
- `hook_post_tool_use.go` 수정: LoadCharacter → BuildReminder 사용, 폴백 유지
- `hook_user_prompt_submit.go` 수정: LoadCharacter → BuildReminder 사용, 폴백 유지
- `spinner.go` 수정: LoadSpinner → BuildSpinnerVerbs 사용, 폴백 유지

**변경 파일**: 4개 (hook 3개 + spinner 1개)
**의존성**: Phase 2 완료 (PersonaLoader가 있어야 연결 가능)
**검증**: 기존 테스트 통과 + 파일 기반 로딩 동작 확인

### Phase 4: convert 도구 업데이트

**목적**: convert 도구가 character/spinner 파일을 인식하고 추출/조립하도록 업데이트

**작업 내용**:
- `PersonaManifest`에 Characters, Spinners 필드 추가
- `CharacterExtractor`, `SpinnerExtractor` 신규 생성
- `ExtractorOrchestrator` 라우팅 업데이트
- `Assembler`에 character/spinner 복사 로직 추가
- 기존 테스트 업데이트

**변경 파일**: 5개 (model 1 + extractor 2 + orchestrator 1 + assembler 1)
**의존성**: Phase 1 완료 (파일 형식 확정), Phase 2와 병렬 가능
**검증**: convert 도구 테스트 통과, 추출/조립 파이프라인 동작 확인

### Phase 5: 스타일 파일 MoAI 잔재 제거 + 검증

**목적**: 스타일 파일에서 MoAI 고유 용어를 제거하고 Do 정체성으로 재작성

**작업 내용**:
- sprint.md: `<do>DONE</do>` XML 마커 제거, TRUST 5 참조 제거
- pair.md: TRUST 5 참조 제거, `.do/config/sections/` 경로 제거, 크기 축소 (577행 → 300행 이하)
- direct.md: TRUST 5 참조 제거, MoAI 설정 경로 제거
- 전체 12 조합 (4 persona x 3 style) 검증

**변경 파일**: 3개 (styles 3개)
**의존성**: Phase 1 완료 (styles가 새 위치에 생성된 후)
**검증**: MoAI 용어 grep으로 잔재 없음 확인, 파일 크기 목표 달성

---

## File Impact Summary

### 신규 생성 (Phase 1)
- `personas/do/characters/young-f.md` (재작성 — YAML frontmatter 추가)
- `personas/do/characters/young-m.md` (재작성)
- `personas/do/characters/senior-f.md` (재작성)
- `personas/do/characters/senior-m.md` (재작성)
- `personas/do/spinners/young-f.yaml` (신규)
- `personas/do/spinners/young-m.yaml` (신규)
- `personas/do/spinners/senior-f.yaml` (신규)
- `personas/do/spinners/senior-m.yaml` (신규)
- `personas/do/styles/sprint.md` (새 위치)
- `personas/do/styles/pair.md` (새 위치)
- `personas/do/styles/direct.md` (새 위치)

### 수정 (Phase 1)
- `personas/do/persona.yaml` (manifest 업데이트)

### 신규 생성 (Phase 2 — godo)
- `cmd/godo/persona_loader.go`
- `cmd/godo/persona_loader_test.go`
- `cmd/godo/testdata/personas/` (test fixtures)

### 수정 (Phase 3 — godo)
- `cmd/godo/hook_session_start.go`
- `cmd/godo/hook_post_tool_use.go`
- `cmd/godo/hook_user_prompt_submit.go`
- `cmd/godo/spinner.go`

### 신규/수정 (Phase 4 — convert)
- `internal/model/persona_manifest.go` (수정)
- `internal/extractor/character.go` (신규)
- `internal/extractor/spinner.go` (신규)
- `internal/extractor/orchestrator.go` (수정)
- `internal/assembler/orchestrator.go` (수정)

### 수정 (Phase 5 — styles cleanup)
- `personas/do/styles/sprint.md`
- `personas/do/styles/pair.md`
- `personas/do/styles/direct.md`

---

## Test Strategy

| Phase | 테스트 유형 | 방법 |
|-------|----------|------|
| Phase 1 | pass (빌드 확인) | YAML frontmatter 파싱 검증 스크립트, persona.yaml 구조 확인 |
| Phase 2 | unit test | `persona_loader_test.go` — 파싱, 렌더링, 폴백, 해상도 순서 |
| Phase 3 | unit + integration | 기존 hook 테스트 통과 + 파일 기반 로딩 통합 테스트 |
| Phase 4 | unit test | convert 도구 기존 테스트 + 새 extractor 테스트 |
| Phase 5 | pass (grep 확인) | MoAI 용어 잔재 grep, 파일 크기 확인 |

---

## Dependencies

```
Phase 1 (파일 생성)
  ├── Phase 2 (PersonaLoader) depends on Phase 1
  │     └── Phase 3 (Hook 연결) depends on Phase 2
  ├── Phase 4 (convert 도구) depends on Phase 1 (병렬 가능 with Phase 2/3)
  └── Phase 5 (스타일 정리) depends on Phase 1

Phase 6 (agents 철학 반영) — 독립 (Phase 1-5와 병렬 가능)
Phase 7 (rules 갱신) — 독립 (Phase 1-5와 병렬 가능)
Phase 8 (skills 철학 반영) — 독립 (Phase 1-5와 병렬 가능)
Phase 9 (commands 갱신) — 독립 (Phase 1-5와 병렬 가능)
```

Phase 2와 Phase 4는 Phase 1 완료 후 **병렬 실행 가능**.
Phase 5는 Phase 1 완료 후 **병렬 실행 가능**.
Phase 3은 Phase 2 완료 후에만 진행.
Phase 6-9는 Phase 1-5와 **완전 독립** — 다른 파일셋이므로 병렬 실행 가능.

### Phase 6: agents/do/ 에이전트 정의 철학 반영

**목적**: 에이전트 정의 파일을 Do 철학에 맞게 갱신 — MoAI 잔재 제거, Do 체크리스트 기반 워크플로우 반영

**작업 내용**:
- 5개 에이전트 파일을 Do 철학에 맞게 재작성
- plan-mode artifact writing rule 반영 (플랜 산출물 위치 규칙)
- MoAI 잔재 제거 (TRUST 5, SPEC, EARS, MoAI 설정 경로 등)
- Do 체크리스트 기반 워크플로우 반영 (살아있는 체크리스트, 커밋 기반 증명)

**변경 파일**: 5개 (manager-ddd.md, manager-tdd.md, team-quality.md, manager-quality.md, manager-project.md)
**의존성**: 없음 (Phase 1-5와 다른 파일셋이므로 병렬 가능)
**검증**: MoAI 용어 grep으로 잔재 없음 확인, 에이전트 정의 구문 검증

### Phase 7: rules/ 개발 규칙 갱신

**목적**: 개발 규칙 파일을 Do 철학에 맞게 갱신 — 커밋 기반 증명, 살아있는 체크리스트, 멱등성 규칙 반영

**작업 내용**:
- 2개 규칙 파일을 Do 철학에 맞게 재작성
- MoAI 특유 용어/개념 제거 (SPEC, EARS, MoAI config 경로 등)
- Do 워크플로우 규칙 반영 (Plan → Checklist → Develop → Test → Report)

**변경 파일**: 2개 (spec-workflow.md, workflow-modes.md)
**의존성**: 없음 (Phase 1-5와 다른 파일셋이므로 병렬 가능)
**검증**: MoAI 용어 grep으로 잔재 없음 확인

### Phase 8: skills/do/ 스킬 철학 반영

**목적**: 스킬 파일을 Do 철학에 맞게 갱신 — Progressive Disclosure 채택, MoAI 잔재 제거

**작업 내용**:
- 8개 스킬 파일을 Do 철학에 맞게 재작성
- Progressive Disclosure 패턴 반영
- MoAI 잔재 제거 (TRUST 5, SPEC, XML 마커 등)

**변경 파일**: 8개 (SKILL.md, references/reference.md, workflows/do.md, workflows/team-do.md, workflows/report.md, workflows/run.md, workflows/plan.md, workflows/test.md)
**의존성**: 없음 (Phase 1-5와 다른 파일셋이므로 병렬 가능)
**검증**: MoAI 용어 grep으로 잔재 없음 확인, 스킬 YAML frontmatter 파싱 가능 여부 확인

### Phase 9: commands/ 확인 및 갱신

**목적**: 명령 파일을 Do 워크플로우에 맞게 갱신

**작업 내용**:
- 5개 커맨드 파일을 확인하고 Do 워크플로우에 맞게 갱신
- MoAI 잔재 제거 (필요시)
- Do 워크플로우 반영 (Plan → Checklist → Develop 등)

**변경 파일**: 5개 (setup.md, style.md, mode.md, checklist.md, plan.md)
**의존성**: 없음 (Phase 1-5와 다른 파일셋이므로 병렬 가능)
**검증**: MoAI 용어 grep으로 잔재 없음 확인, 커맨드 구문 검증

---

## File Impact Summary (Phase 6-9)

### 수정 (Phase 6 — agents)
- `personas/do/agents/do/manager-ddd.md`
- `personas/do/agents/do/manager-tdd.md`
- `personas/do/agents/do/team-quality.md`
- `personas/do/agents/do/manager-quality.md`
- `personas/do/agents/do/manager-project.md`

### 수정 (Phase 7 — rules)
- `personas/do/rules/do/workflow/spec-workflow.md`
- `personas/do/rules/do/workflow/workflow-modes.md`

### 수정 (Phase 8 — skills)
- `personas/do/skills/do/SKILL.md`
- `personas/do/skills/do/references/reference.md`
- `personas/do/skills/do/workflows/do.md`
- `personas/do/skills/do/workflows/team-do.md`
- `personas/do/skills/do/workflows/report.md`
- `personas/do/skills/do/workflows/run.md`
- `personas/do/skills/do/workflows/plan.md`
- `personas/do/skills/do/workflows/test.md`

### 수정 (Phase 9 — commands)
- `personas/do/commands/do/setup.md`
- `personas/do/commands/do/style.md`
- `personas/do/commands/do/mode.md`
- `personas/do/commands/do/checklist.md`
- `personas/do/commands/do/plan.md`

---

## Risks

| 위험 | 영향 | 완화 |
|------|------|------|
| godo 소스가 다른 프로젝트에 있음 | Phase 2/3 작업 위치 확인 필요 | 작업 시작 전 godo 프로젝트 경로 확인 |
| 스피너 데이터 추출 시 누락 | 기존 동작과 불일치 | 하드코딩된 배열과 1:1 비교 검증 |
| 스타일 파일 축소 시 필요한 내용 삭제 | 스타일 동작 변경 | MoAI 잔재만 제거, Do 고유 내용 보존 |
| persona.yaml 형식 변경이 기존 convert 파이프라인에 영향 | 빌드 실패 | Phase 4에서 호환성 테스트 |

---

**Author**: plan agent
**Status**: Plan complete
**Next Step**: Checklist 작성

---

## 변경 이력

### 2026-02-16 Phase 10 추가
- Phase 10: dev-* 개발 규칙 파일을 페르소나 패키지에 추가

---

### Phase 10: dev-* 규칙 파일 페르소나 패키지 추가

**목적**: do-focus의 핵심 개발 규칙(dev-checklist, dev-workflow, dev-testing, dev-environment, file-reading)을 convert 페르소나 패키지에 포함 — 배포 시 사용자 프로젝트에 설치되도록

**작업 내용**:
- do-focus `.claude/rules/` 에서 5개 dev-* 규칙 파일을 `personas/do/rules/`에 동일 복사
- 소스와 동일성 검증 (diff)
- 파일: dev-checklist.md, dev-workflow.md, dev-testing.md, dev-environment.md, file-reading.md

**변경 파일**: 5개 (신규 — 페르소나 패키지에 추가)
**의존성**: 없음 (Phase 1-9와 독립)
**검증**: `diff` 명령으로 소스와 동일성 확인

### File Impact Summary (Phase 10)

### 신규 (Phase 10 — rules)
- `personas/do/rules/dev-checklist.md`
- `personas/do/rules/dev-workflow.md`
- `personas/do/rules/dev-testing.md`
- `personas/do/rules/dev-environment.md`
- `personas/do/rules/file-reading.md`
