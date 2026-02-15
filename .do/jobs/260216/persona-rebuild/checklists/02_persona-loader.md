# persona-loader: PersonaLoader 구현 (Phase 2)
상태: [o] | 담당: expert-backend

## Problem Summary
- godo가 페르소나 데이터를 `buildPersona()` 하드코딩으로만 읽음 — 파일 기반 로딩 없음
- Phase 1에서 생성한 캐릭터/스피너 파일을 런타임에 파싱하는 코드가 필요
- Architecture의 Approach C (Hybrid)에 따라 파일 우선 로딩 + 하드코딩 폴백 구현 필요
- Honorific 템플릿 렌더링(`{{name}}` → 실제 이름), 리마인더 생성, 스피너 동사 생성 기능 필요

## Acceptance Criteria
- [o] `PersonaData` 구조체: ID, Name, HonorificTemplate, HonorificDefault, Tone, CharacterSummary, Relationship, FullContent
- [o] `SpinnerData` 구조체: Persona, SuffixCycle, Suffixes, Stems
- [o] `LoadCharacter(personaDir, personaType)` — .md 파일의 YAML frontmatter 파싱
- [o] `LoadSpinner(personaDir, personaType)` — .yaml 파일 파싱
- [o] `BuildHonorific(userName)` — `{{name}}` 치환, 빈 이름 시 default 반환
- [o] `BuildReminder(userName)` — `반드시 "{honorific}"로 호칭할 것. 말투: {tone}` 형식
- [o] `BuildSpinnerVerbs()` — stems + suffixes 조합으로 동사 리스트 생성
- [o] 파일 미존재 시 에러 반환 (호출부에서 폴백 처리)
- [o] 단위 테스트 전체 통과 (파싱, 렌더링, 폴백, 해상도 순서)
- [o] 커밋 완료

## Solution Approach
- Go 표준 라이브러리 `gopkg.in/yaml.v3` 또는 프로젝트 기존 YAML 라이브러리 사용
- frontmatter 파싱: `---` 구분자 사이의 YAML을 추출하여 파싱
- 파일 해상도 순서: `.claude/personas/do/` → `personas/do/` → 에러 반환
- 대안 고려: `go:embed`로 파일 임베딩 → 기각 (파일 수정 시 재컴파일 필요, Approach C의 유연성 상실)

## Test Strategy
- unit test: `persona_loader_test.go`
  - 캐릭터 frontmatter 파싱: 모든 필드 정확히 추출
  - 스피너 YAML 파싱: stems + suffix 정확히 추출
  - Honorific 렌더링: 이름 있을 때/없을 때
  - Reminder 생성: 형식 일치
  - SpinnerVerbs 생성: 기존 하드코딩 출력과 1:1 일치
  - 파일 미존재 시 에러 반환
  - 파일 해상도 순서: .claude/personas > personas > 에러

## Critical Files

### 항목 #4: PersonaData 구조체 + LoadCharacter 구현
- **생성 대상**: `cmd/godo/persona_loader.go` — PersonaData, LoadCharacter, BuildHonorific, BuildReminder
- **참조 파일**: `architecture-persona-system.md` Section 2.5 — Go 인터페이스 정의
- **참조 파일**: `cmd/godo/hook_session_start.go` — 기존 buildPersona() 구현

### 항목 #5: SpinnerData + LoadSpinner + BuildSpinnerVerbs 구현
- **수정 대상**: `cmd/godo/persona_loader.go` — SpinnerData, SpinnerStem, LoadSpinner, BuildSpinnerVerbs 추가
- **참조 파일**: `architecture-persona-system.md` Section 2.3 — 스피너 형식
- **참조 파일**: `cmd/godo/spinner.go` — 기존 getPersonaSpinnerVerbs() 구현

### 항목 #6: 단위 테스트 + 테스트 픽스처
- **생성 대상**: `cmd/godo/persona_loader_test.go` — 모든 공개 함수 테스트
- **생성 대상**: `cmd/godo/testdata/personas/characters/young-f.md` — 테스트 픽스처
- **생성 대상**: `cmd/godo/testdata/personas/spinners/young-f.yaml` — 테스트 픽스처
- **생성 대상**: `cmd/godo/testdata/personas/characters/invalid.md` — 잘못된 형식 픽스처

## Risks
- YAML frontmatter 파서가 `---` 구분자를 잘못 인식할 수 있음: markdown body에 `---`가 포함된 경우 → 첫 번째 `---` 쌍만 파싱하도록 구현
- 기존 YAML 라이브러리 의존성 확인 필요: 프로젝트에 이미 사용 중인 라이브러리 우선 사용
- SpinnerVerbs 출력이 기존 하드코딩과 정확히 일치해야 함 — suffix 조합 로직에 off-by-one 주의

## Progress Log
- 2026-02-16 02:40:00 [~] 작업 시작: PersonaData/SpinnerData 구조체 설계
- 2026-02-16 02:43:22 [~] persona_loader.go 구현 완료: LoadCharacter, LoadSpinner, BuildHonorific, BuildReminder, BuildSpinnerVerbs, resolvePersonaDir
- 2026-02-16 02:47:15 [~] persona_loader_test.go 작성 완료: 전체 공개 함수 테스트 + 테스트 픽스처
- 2026-02-16 02:49:30 [*] 단위 테스트 실행 — 파일 기반 출력이 기존 hardcoded buildPersona() 출력과 일치 확인
- 2026-02-16 02:51:14 [o] 완료 (commit: edca1c6)

## FINAL STEP: Commit (절대 생략 금지)
- [o] `git add` — persona_loader.go, persona_loader_test.go, testdata/ 파일만 스테이징
- [o] `git diff --cached` — 의도한 변경만 포함되었는지 확인
- [o] `git commit` — 커밋 메시지에 WHY 포함 (commit: edca1c6)
- [o] 커밋 해시를 Progress Log에 기록 (edca1c6)
⚠️ 이 섹션을 완료하지 않으면 작업은 미완료(incomplete) 상태임

## Lessons Learned (완료 시 작성)
- 잘된 점: Architecture의 Approach C (Hybrid) 설계대로 파일 우선 + 하드코딩 폴백 구현 완료. 테스트가 기존 출력과 1:1 일치를 검증하여 회귀 방지.
- 어려웠던 점: frontmatter 파싱에서 --- 구분자 처리 시 markdown body에도 ---가 있을 수 있어 첫 번째 쌍만 추출하는 로직 필요.
- 다음에 다르게 할 점: 서브 체크리스트를 구현 완료 즉시 갱신해야 작업 추적이 끊기지 않음.
