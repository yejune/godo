# persona-files: 페르소나 파일 생성 (Phase 1)
상태: [ ] | 담당: expert-backend

## Problem Summary
- godo의 `buildPersona()` 함수에 하드코딩된 페르소나 데이터를 파일로 외부화해야 한다
- 기존 characters/*.md 파일은 YAML frontmatter가 없어 machine-readable하지 않다
- 스피너 데이터는 Go 소스의 배열에만 존재하고 파일이 아예 없다
- 스타일 파일은 `output-styles/do/`에 있으나 architecture에서 `styles/`로 이동이 필요하다
- persona.yaml 매니페스트가 새 파일 타입을 참조하지 않는다

## Acceptance Criteria
- [ ] 4개 캐릭터 파일에 YAML frontmatter 추가 (id, name, honorific_template, honorific_default, tone, character_summary, relationship)
- [ ] 4개 스피너 YAML 파일 생성 (persona, suffix_pattern, stems)
- [ ] 3개 스타일 파일이 `personas/do/styles/`에 존재
- [ ] persona.yaml에 characters, spinners, styles 섹션 추가
- [ ] 모든 YAML이 파싱 가능 (검증: `yq` 또는 Go yaml 파서로 확인)
- [ ] 커밋 완료

## Solution Approach
- Architecture 문서 Section 2 (Core Interfaces / Data Formats)의 형식을 정확히 따름
- 캐릭터 파일: 기존 markdown body는 보존하되 YAML frontmatter를 추가
- 스피너 파일: godo 소스의 `spinnerStemsYoungF`, `spinnerStemsDefault` 배열에서 1:1 추출
- 스타일 파일: `output-styles/do/`에서 `styles/`로 복사 (MoAI 정리는 Phase 5에서)
- 대안 고려: 캐릭터 파일을 전체 YAML로 작성 → 기각 (markdown body의 자유도를 보존하기 위해 frontmatter 방식 선택)

## Test Strategy
- pass (빌드 확인): YAML frontmatter를 `yq` 또는 스크립트로 파싱하여 필수 필드 존재 확인
- pass (구조 확인): persona.yaml의 참조 경로가 실제 파일과 일치하는지 확인

## Critical Files

### 항목 #1: 캐릭터 파일 YAML frontmatter 재작성
- **수정 대상**: `personas/do/characters/young-f.md` — frontmatter 추가
- **수정 대상**: `personas/do/characters/young-m.md` — frontmatter 추가
- **수정 대상**: `personas/do/characters/senior-f.md` — frontmatter 추가
- **수정 대상**: `personas/do/characters/senior-m.md` — frontmatter 추가
- **참조 파일**: godo `hook_session_start.go` — buildPersona() 데이터 원본
- **참조 파일**: `architecture-persona-system.md` Section 2.1 — frontmatter 형식

### 항목 #2: 스피너 YAML 파일 신규 생성
- **생성 대상**: `personas/do/spinners/young-f.yaml` — young-f 전용 stems + suffix
- **생성 대상**: `personas/do/spinners/young-m.yaml` — default stems + young-m suffix
- **생성 대상**: `personas/do/spinners/senior-f.yaml` — default stems + senior-f suffix
- **생성 대상**: `personas/do/spinners/senior-m.yaml` — default stems + senior-m suffix
- **참조 파일**: godo `spinner.go` — spinnerStemsYoungF, spinnerStemsDefault 배열

### 항목 #3: 스타일 파일 + persona.yaml 업데이트
- **생성 대상**: `personas/do/styles/sprint.md` — output-styles/do/sprint.md에서 복사
- **생성 대상**: `personas/do/styles/pair.md` — output-styles/do/pair.md에서 복사
- **생성 대상**: `personas/do/styles/direct.md` — output-styles/do/direct.md에서 복사
- **수정 대상**: `personas/do/persona.yaml` — characters, spinners, styles 참조 추가
- **참조 파일**: `architecture-persona-system.md` Section 2.4 — manifest 형식

## Risks
- 스피너 데이터 추출 시 stems 누락 가능: godo 소스와 1:1 비교 필수
- 캐릭터 frontmatter의 tone 필드가 기존 buildPersona()의 Tone과 정확히 일치해야 함 — 불일치 시 리마인더 메시지 변경됨
- output-styles/do/ 원본을 삭제하면 안 됨 (Phase 5에서 정리 후 판단)

## Progress Log
- (작업 시작 시 기록)

## FINAL STEP: Commit (절대 생략 금지)
- [ ] `git add` — 변경된 파일만 스테이징 (characters 4개 + spinners 4개 + styles 3개 + persona.yaml)
- [ ] `git diff --cached` — 의도한 변경만 포함되었는지 확인
- [ ] `git commit` — 커밋 메시지에 WHY 포함
- [ ] 커밋 해시를 Progress Log에 기록
⚠️ 이 섹션을 완료하지 않으면 작업은 미완료(incomplete) 상태임

## Lessons Learned (완료 시 작성)
- 잘된 점:
- 어려웠던 점:
- 다음에 다르게 할 점:
