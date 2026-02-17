# Testing Mandatory Rules [HARD]

## 적용 범위
- [HARD] 이 규칙은 **테스트 가능한 코드**(비즈니스 로직, API, 데이터 계층)에 적용
- [HARD] 테스트 불가능한 변경(CSS, 설정 파일, 문서, 훅 스크립트 등)은 대안 검증으로 대체:
  - 빌드 확인: `go build`, `npm run build`, `docker compose config`
  - 수동 확인: 브라우저 확인, CLI 실행 확인
  - 구문 검사: lint, type check
- [HARD] 체크리스트 항목에 검증 방법 명시 필수 — 예: `(검증: 빌드 확인)`, `(테스트: 단위 테스트)`

## 테스트 철학
- [HARD] 테스트는 **행위(behavior)**를 검증 — 구현 방식이 아닌, 코드가 하는 일을 테스트
- [HARD] FIRST 원칙 준수: Fast(빠르게), Independent(독립적), Repeatable(반복 가능), Self-validating(자가 검증), Timely(적시 작성)
- [HARD] Test Pyramid 준수: Unit > Integration > E2E 비율 유지 — 단, 모든 계층에서 Real DB 사용
- [HARD] Flaky 테스트는 버그 — 발견 즉시 원인 조사 후 수정, 재실행으로 무시 금지
- [HARD] 변이 테스트 사고방식: "이 코드 한 줄을 바꾸면 테스트가 실패하는가?" — 실패하지 않으면 테스트 부족

## 테스트 통과 필수
- [HARD] 전체 테스트 반드시 통과 — skip, timeout bypass 금지
- [HARD] "시간 초과로 완료되지 않았다"는 보고 금지 — 반드시 해결
- [HARD] 목업하거나 생략 금지, 스킵된 테스트도 구현 필수
- [HARD] 테스트 실패 시 코드 수정 필수, 테스트를 수정하지 않음

## AI 안티패턴 금지 [HARD]
- [HARD] assertion 약화 금지 — `assertEqual`을 `assertContains`로, 정확한 값을 `any()`로 바꾸지 않음
- [HARD] try/catch로 에러 삼키기 금지 — 에러를 잡아서 테스트를 녹색으로 만들지 않음
- [HARD] 테스트 기대값을 잘못된 출력에 맞추기 금지 — 코드를 수정할 것
- [HARD] `time.sleep()` / 임의 지연 금지 — 타이밍 문제의 진짜 원인을 찾을 것
- [HARD] 실패하는 테스트 삭제/주석처리 금지
- [HARD] 정확한 값을 알 때 와일드카드 매처(`any()`, `mock.ANY`) 사용 금지
- [HARD] happy path만 테스트 금지 — 에러 경로, 엣지 케이스, 경계값 테스트 필수

## 테스트 데이터 관리 [HARD]
- [HARD] 각 테스트는 자체 데이터를 설정하고 정리 (Arrange-Act-Assert / Given-When-Then)
- [HARD] 정리 방식: 트랜잭션 롤백 또는 truncate — 공유 fixture에 의존 금지
- [HARD] 테스트 데이터는 최소한으로 — 해당 테스트에 필요한 것만, 그 이상 금지
- [HARD] 테스트 데이터 생성은 factory/builder 패턴 사용 — 매번 raw SQL 직접 작성 금지
- [HARD] 테스트별 고유 식별자 사용 (UUID/timestamp suffix) — 병렬 실행 시 충돌 방지

## 테스트 품질
- [HARD] 테스트 이름은 시나리오를 설명: `test_login_fails_with_expired_token` — `test_login_2` 금지
- [HARD] 하나의 테스트에 하나의 논리적 검증 (동일 검증을 위한 다중 assert는 허용)
- [HARD] 실패 메시지는 **무엇이** 왜 잘못되었는지 명확히 설명
- [HARD] 테스트 코드도 프로덕션 코드와 동일한 품질 기준 적용 — 중복 제거, 가독성 유지

## Real DB Only
- [HARD] 데이터베이스는 실제 쿼리만 — mock DB, in-memory DB, SQLite 대체 금지
- [HARD] 테스트는 Docker Compose 서비스의 실제 DB에 연결
- [HARD] 외부 API만 mock 허용 — 데이터 계층은 절대 mock 금지

## 병렬성 & 시간 제한
- [HARD] 테스트는 병렬 실행 안전하게 설계 (concurrency-safe)
- [HARD] 동시성/Lock/Race condition 테스트는 동일 파일 내 메소드로 분리하여 순차 실행
- [HARD] 순차 테스트는 서로의 트랜잭션을 침범하지 않아야 함
- [HARD] 병렬 테스트 시 에이전트에게 Docker Compose 로컬 환경/컨테이너 정보 전달 필수
- [HARD] 개별 테스트 < 3분, 전체 스위트 < 10분 — 초과 시 파일 분리/프로세스 개선
- [HARD] 테스트 간 리소스 격리: DB 스키마, 포트, 파일 경로 등 공유 자원 충돌 방지

## 실행 순서
- [HARD] 개별 테스트 먼저 실행, 확신이 생기면 전체 스위트 실행
- [HARD] 교차 검증 필수: view ↔ DB ↔ worker ↔ model ↔ business logic ↔ controller
- [HARD] 비즈니스 로직을 파악하고 올바른 테스트가 작성되었는지 검증
- [HARD] 테스트를 위한 테스트가 아닌, 문제를 찾아내고 개선에 기여하는 테스트 작성

## DB 트랜잭션
- [HARD] 논리적으로 분리된 DB가 동일 서버인지 반드시 사용자에게 확인 (AskUserQuestion)
- [HARD] 동일 서버 시: service가 X_DB 사용 → X_DB에 트랜잭션, Z_DB.table로 참조
