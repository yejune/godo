# SPEC 워크플로우

토큰 예산 관리를 포함한 MoAI의 3단계 개발 워크플로우.

## 단계 개요

| 단계 | 명령어 | 에이전트 | 토큰 예산 | 목적 |
|-------|---------|-------|--------------|---------|
| Plan | /moai plan | manager-spec | 30K | SPEC 문서 생성 |
| Run | /moai run | manager-ddd/tdd (quality.yaml 설정에 따름) | 180K | DDD/TDD 구현 |
| Sync | /moai sync | manager-docs | 40K | 문서 동기화 |

## Plan 단계

EARS 형식을 사용하여 포괄적인 명세서를 작성합니다.

토큰 전략:
- 할당: 30,000 토큰
- 요구사항만 로드
- 완료 후 /clear 실행
- 구현을 위해 45-50K 토큰 절약

출력물:
- `.moai/specs/SPEC-XXX/spec.md`의 SPEC 문서
- EARS 형식 요구사항
- 인수 기준
- 기술적 접근 방식

## Run 단계

설정된 개발 방법론을 사용하여 명세서를 구현합니다.

토큰 전략:
- 할당: 180,000 토큰
- 선택적 파일 로드
- 70% 더 큰 구현 가능

개발 방법론:
- quality.yaml에서 설정 (development_mode: ddd, tdd, 또는 hybrid)
- 상세 방법론 사이클은 @workflow-modes.md 참조

성공 기준:
- 모든 SPEC 요구사항 구현 완료
- 방법론별 테스트 통과
- 85%+ 코드 커버리지
- TRUST 5 품질 게이트 통과

## Sync 단계

문서를 생성하고 배포를 준비합니다.

토큰 전략:
- 할당: 40,000 토큰
- 결과 캐싱
- 60% 더 적은 중복 파일 읽기

출력물:
- API 문서
- README 업데이트
- CHANGELOG 항목
- Pull request

## 완료 마커

AI는 작업 완료를 알리기 위해 마커를 사용합니다:
- `<moai>DONE</moai>` - 작업 완료
- `<moai>COMPLETE</moai>` - 전체 완료

## 컨텍스트 관리

/clear 전략:
- /moai plan 완료 후 (필수)
- 컨텍스트가 150K 토큰을 초과할 때
- 주요 단계 전환 전

점진적 공개:
- 레벨 1: 메타데이터만 (~100 토큰)
- 레벨 2: 트리거 시 스킬 본문 (~5000 토큰)
- 레벨 3: 번들 파일 온디맨드

## 단계 전환

Plan에서 Run으로:
- 트리거: SPEC 문서 승인
- 조치: /clear 실행 후 /moai run SPEC-XXX

Run에서 Sync로:
- 트리거: 구현 완료, 테스트 통과
- 조치: /moai sync SPEC-XXX 실행

## Agent Teams 변형

팀 모드가 활성화된 경우 (workflow.team.enabled 및 AGENT_TEAMS 환경변수), 단계들은 서브에이전트 대신 Agent Teams로 실행될 수 있습니다.

### 팀 모드 단계 개요

| 단계 | 서브에이전트 모드 | 팀 모드 | 조건 |
|-------|---------------|-----------|-----------|
| Plan | manager-spec (단일) | researcher + analyst + architect (병렬) | 복잡도 >= 임계값 |
| Run | manager-ddd/tdd (순차) | backend-dev + frontend-dev + tester (병렬) | 도메인 >= 3 또는 파일 >= 10 |
| Sync | manager-docs (단일) | manager-docs (항상 서브에이전트) | N/A |

### 팀 모드 Plan 단계
- 병렬 조사 팀을 위한 TeamCreate
- 팀원들이 코드베이스 탐색, 요구사항 분석, 접근 방식 설계
- MoAI가 SPEC 문서로 종합
- 팀 종료, Run 단계 전 /clear

### 팀 모드 Run 단계
- 구현 팀을 위한 TeamCreate
- 파일 소유권 경계가 있는 작업 분해
- 팀원들이 공유 목록에서 작업을 자체 할당
- 모든 구현 완료 후 품질 검증
- 팀 종료

### 모드 선택
- --team 플래그: 팀 모드 강제
- --solo 플래그: 서브에이전트 모드 강제
- 플래그 없음 (기본): 복잡도 기반 선택
- 임계값은 workflow.yaml의 team.auto_selection 참조

### 폴백
팀 모드가 실패하거나 사용 불가능한 경우:
- 서브에이전트 모드로 정상적인 폴백
- 마지막 완료 작업부터 계속
- 데이터 손실 또는 상태 손상 없음
