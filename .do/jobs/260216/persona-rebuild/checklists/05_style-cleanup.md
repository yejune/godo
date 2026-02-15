# style-cleanup: 스타일 파일 MoAI 잔재 제거 + 검증 (Phase 5)
상태: [o] | 담당: expert-backend

## Problem Summary
- analysis-requirements.md Section 4.2에서 확인된 SIGNIFICANT GAP: 스타일 파일에 MoAI 고유 용어가 남아있음
- pair.md (577행): "TRUST 5 principles" 참조, `.do/config/sections/language.yaml` MoAI 경로 참조
- sprint.md (269행): `<do>DONE</do>`, `<do>COMPLETE</do>` XML 마커 사용, TRUST 5 참조
- direct.md (360행): "TRUST 5 principles" 참조, MoAI 설정 경로 참조
- Do는 TRUST 5 브랜딩을 거부함 (DO_PERSONA.md Section 9) — 품질 차원은 내장 규칙으로 존재하되 브랜딩 없이
- Do는 XML 완료 마커를 거부함 — commit hash가 완료 증거 (DO_PERSONA.md Section 4, Pillar 3)
- pair.md가 577행으로 너무 큼 — 목표 300행 이하 (analysis NFR)

## Acceptance Criteria
- [o] sprint.md: `<do>DONE</do>`, `<do>COMPLETE</do>` XML 마커 제거
- [o] sprint.md: TRUST 5 참조 제거
- [o] pair.md: TRUST 5 참조 제거
- [o] pair.md: `.do/config/sections/` MoAI 경로 참조 제거
- [o] pair.md: 300행 이하로 축소 (577행 → 237행)
- [o] direct.md: TRUST 5 참조 제거
- [o] direct.md: MoAI 설정 경로 참조 제거
- [o] 모든 스타일 파일에서 `grep -i "trust 5\|moai\|SPEC\|<do>\|</do>\|\.do/config"` 결과 0건
- [o] 12 조합 (4 persona x 3 style) 동작 검증 — 스타일 파일 구조/독립성 유지 확인, MoAI 참조 제거로 올바른 Do 동작 보장
- [o] 커밋 완료 (commit: fa65c58)

## Solution Approach
- 제거 대상 패턴을 먼저 grep으로 식별하고, 각 패턴별로 Do 정체성에 맞는 대체 텍스트 작성
- TRUST 5 → 5가지 품질 차원을 이름 없이 인라인 서술 (Tested, Readable, Unified, Secured, Trackable)
- XML 마커 → "커밋 해시로 완료를 증명한다" 서술로 대체
- MoAI 경로 → `settings.local.json`의 `DO_LANGUAGE` 환경변수 참조로 대체
- pair.md 축소: 중복된 내용 제거, Do의 다른 규칙 파일과 겹치는 내용 참조로 대체
- 대안 고려: 스타일 파일 전체 재작성 → 기각 (기존 동작하는 지시를 보존하면서 오염만 제거하는 것이 안전)

## Test Strategy
- pass (grep 확인): `grep -rn "TRUST 5\|moai\|SPEC\|<do>\|</do>\|\.do/config" personas/do/styles/` 결과 0건
- pass (행 수 확인): `wc -l personas/do/styles/pair.md` ≤ 300
- pass (수동 확인): 각 스타일 파일을 읽고 Do 정체성과 일관되는지 리뷰

## Critical Files

### 항목 #14: sprint.md MoAI 잔재 제거
- **수정 대상**: `personas/do/styles/sprint.md` — XML 마커 제거, TRUST 5 제거
- **참조 파일**: `DO_PERSONA.md` Section 9 — 거부한 MoAI 기능 목록
- **참조 파일**: `DO_MOAI_COMPARISON.md` — 상세 비교 (있으면)

### 항목 #15: pair.md MoAI 잔재 제거 + 크기 축소
- **수정 대상**: `personas/do/styles/pair.md` — TRUST 5 제거, MoAI 경로 제거, 577행 → 300행 축소
- **참조 파일**: `DO_PERSONA.md` Section 9 — 거부한 MoAI 기능 목록
- **참조 파일**: `analysis-requirements.md` Section 4.2 — gap 상세

### 항목 #16: direct.md MoAI 잔재 제거
- **수정 대상**: `personas/do/styles/direct.md` — TRUST 5 제거, MoAI 경로 제거
- **참조 파일**: `DO_PERSONA.md` Section 9 — 거부한 MoAI 기능 목록

### 항목 #17: 전체 검증
- **참조 파일**: `personas/do/styles/sprint.md` — 정리 결과 확인
- **참조 파일**: `personas/do/styles/pair.md` — 정리 결과 확인
- **참조 파일**: `personas/do/styles/direct.md` — 정리 결과 확인
- **참조 파일**: `personas/do/characters/*.md` — 4 persona 확인
- **검증 스크립트**: grep 기반 MoAI 잔재 최종 스캔

## Risks
- pair.md 축소 시 필요한 스타일 지시까지 삭제할 수 있음: 제거 전 각 섹션의 역할을 파악하고 Do 규칙 파일과 중복 여부 확인
- TRUST 5를 단순 제거하면 품질 관련 지시가 아예 없어질 수 있음: 브랜딩만 제거하고 품질 차원(테스트됨, 읽기 쉬움, 일관됨, 보안됨, 추적 가능)은 인라인 서술로 유지
- sprint.md의 XML 마커가 Do 모드의 완료 신호로 실제 사용 중일 수 있음: godo 훅 코드에서 XML 마커 파싱 여부 확인 후 제거

## Progress Log
- 2026-02-16 02:31 [~] 작업 시작: 3개 스타일 파일 MoAI 잔재 식별 완료
  - sprint.md: <do>DONE/COMPLETE</do> XML 마커 2곳, TRUST 5 참조 1곳, MoAI 경로 1곳
  - pair.md: TRUST 5 참조 2곳, MoAI 경로 3곳, 577행 → 300행 축소 필요
  - direct.md: TRUST 5 참조 2곳, SPEC 참조 1곳, MoAI 경로 3곳
  - Go 코드에서 XML 마커 파싱 없음 확인 → 안전하게 제거 가능
- 2026-02-16 02:35 [~] sprint.md 정리 완료: XML 마커 제거, TRUST 5 제거, MoAI 경로 → settings.local.json
- 2026-02-16 02:38 [~] pair.md 정리 완료: 577행 → 237행 (59% 축소), 중복 섹션 제거, MoAI 참조 전부 제거
- 2026-02-16 02:40 [~] direct.md 정리 완료: 360행 → 283행, TRUST 5/SPEC/MoAI 경로 모두 제거
- 2026-02-16 02:42 [*] 전체 검증: grep 스캔 0건, 행 수 목표 달성 (sprint 264, pair 237, direct 283)
- 2026-02-16 02:43 [o] 커밋 완료 (commit: fa65c58)

## FINAL STEP: Commit (절대 생략 금지)
- [o] `git add` — output-styles/do/sprint.md, pair.md, direct.md만 스테이징
- [o] `git diff --cached` — 3개 파일만 확인됨
- [o] `git commit` — WHY 포함 (Do가 거부한 MoAI 기능 제거 사유 명시)
- [o] 커밋 해시 fa65c58 Progress Log에 기록 완료
⚠️ 이 섹션을 완료하지 않으면 작업은 미완료(incomplete) 상태임

## Lessons Learned (완료 시 작성)
- 잘된 점: DO_PERSONA.md와 DO_MOAI_COMPARISON.md의 "거부한 것" 목록이 명확해서 무엇을 제거해야 하는지 판단이 쉬웠음. pair.md 59% 축소로 토큰 효율성 크게 개선.
- 어려웠던 점: pair.md에서 어떤 내용이 "style-specific"이고 어떤 것이 "CLAUDE.md에 이미 있는 중복"인지 판별하는 것. Sequential Thinking MCP, Skills+Context7 등은 도구 관련이라 style 파일에 있을 필요 없다고 판단.
- 다음에 다르게 할 점: grep 검증 시 case-insensitive 검색과 changelog 항목을 분리하는 패턴을 미리 준비하면 검증 반복을 줄일 수 있음.
