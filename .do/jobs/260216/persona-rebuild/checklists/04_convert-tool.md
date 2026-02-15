# convert-tool: convert 도구 업데이트 (Phase 4)
상태: [ ] | 담당: expert-backend

## Problem Summary
- convert 도구의 `PersonaManifest` 모델이 characters, spinners 필드를 지원하지 않음
- 추출(extraction) 단계에서 `characters/*.md`와 `spinners/*.yaml`을 페르소나 에셋으로 분류하는 로직이 없음
- 조립(assembly) 단계에서 캐릭터/스피너 파일을 출력 디렉토리로 복사하지 않음
- Phase 1에서 persona.yaml에 추가한 characters/spinners 참조를 convert 파이프라인이 처리해야 함

## Acceptance Criteria
- [ ] PersonaManifest에 `Characters []string`과 `Spinners []string` 필드 추가
- [ ] CharacterExtractor: `characters/*.md` 파일을 페르소나 에셋으로 분류
- [ ] SpinnerExtractor: `spinners/*.yaml` 파일을 페르소나 에셋으로 분류
- [ ] ExtractorOrchestrator가 character/spinner 파일을 새 extractor로 라우팅
- [ ] Assembler가 character/spinner 파일을 출력 디렉토리로 복사
- [ ] 기존 convert 테스트 전체 통과 (회귀 없음)
- [ ] 새 extractor에 대한 단위 테스트 추가
- [ ] 커밋 완료

## Solution Approach
- 기존 StyleExtractor 구현 패턴을 참고하여 CharacterExtractor, SpinnerExtractor 생성
- ExtractorOrchestrator의 파일 라우팅 로직에 `characters/` 경로 매칭 추가
- Assembler의 복사 로직에 characters/spinners 디렉토리 처리 추가
- persona.yaml 파싱 시 새 필드를 자동으로 읽도록 Go struct tag 추가
- 대안 고려: 단일 GenericExtractor로 모든 파일 처리 → 기각 (파일 타입별 검증 로직이 다름 — .md vs .yaml 형식 검증)

## Test Strategy
- unit test: CharacterExtractor — .md 파일 분류 정확성
- unit test: SpinnerExtractor — .yaml 파일 분류 정확성
- unit test: PersonaManifest 파싱 — Characters/Spinners 필드 포함
- integration: 전체 extract → assemble 파이프라인에서 character/spinner 파일이 출력에 포함되는지 확인

## Critical Files

### 항목 #11: PersonaManifest 모델 업데이트
- **수정 대상**: `internal/model/persona_manifest.go` — Characters, Spinners 필드 + yaml struct tag 추가
- **참조 파일**: `personas/do/persona.yaml` — 매니페스트 실제 형식
- **참조 파일**: `architecture-persona-system.md` Section 2.4 — 매니페스트 명세

### 항목 #12: CharacterExtractor + SpinnerExtractor 생성
- **생성 대상**: `internal/extractor/character.go` — CharacterExtractor 구현
- **생성 대상**: `internal/extractor/spinner.go` — SpinnerExtractor 구현
- **참조 파일**: `internal/extractor/style.go` — 기존 StyleExtractor 패턴 참고 (있으면)
- **참조 파일**: `internal/extractor/orchestrator.go` — 기존 라우팅 패턴 참고

### 항목 #13: Orchestrator 라우팅 + Assembler 복사
- **수정 대상**: `internal/extractor/orchestrator.go` — `characters/` 및 `spinners/` 경로 라우팅 추가
- **수정 대상**: `internal/assembler/orchestrator.go` — character/spinner 파일 복사 로직 추가
- **참조 파일**: `architecture-persona-system.md` Section 5 — convert 도구 통합 명세

## Risks
- 기존 ExtractorOrchestrator의 라우팅 로직 구조를 먼저 파악해야 함 — Read Before Write 필수
- persona.yaml의 Characters/Spinners 필드가 optional이어야 함 — 기존 매니페스트(이 필드 없는)와 호환 유지
- Assembler가 디렉토리 구조를 보존하면서 복사해야 함 — `characters/young-f.md` → `output/characters/young-f.md`
- convert 도구의 기존 테스트가 깨지지 않도록 새 필드를 추가할 때 omitempty 사용

## Progress Log
- (작업 시작 시 기록)

## FINAL STEP: Commit (절대 생략 금지)
- [ ] `git add` — persona_manifest.go, character.go, spinner.go, orchestrator.go (extractor + assembler) 만 스테이징
- [ ] `git diff --cached` — 의도한 변경만 포함되었는지 확인
- [ ] `git commit` — 커밋 메시지에 WHY 포함
- [ ] 커밋 해시를 Progress Log에 기록
⚠️ 이 섹션을 완료하지 않으면 작업은 미완료(incomplete) 상태임

## Lessons Learned (완료 시 작성)
- 잘된 점:
- 어려웠던 점:
- 다음에 다르게 할 점:
