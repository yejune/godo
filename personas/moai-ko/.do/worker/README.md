# Do Worker Service

백그라운드 DB 작업 처리를 위한 Go 기반 Worker Service. Claude 토큰 소모 없이 메모리 관리를 수행합니다.

## 특징

- **토큰 효율**: Claude 세션 외부에서 DB 작업 처리
- **팀 컨텍스트**: 다른 팀원의 작업 내용 조회
- **에이전트 추적**: observations에 agent_name 필드로 에이전트별 추적
- **다중 DB 지원**: SQLite (기본) 및 MySQL 지원

## 빠른 시작

```bash
# 의존성 설치
make deps

# 빌드
make build

# 실행
./bin/do-worker
# 또는
make run
```

## API 엔드포인트

| Method | Endpoint | 설명 |
|--------|----------|------|
| GET | `/health` | 헬스체크 |
| GET | `/api/context/inject` | 세션 시작용 컨텍스트 주입 |
| POST | `/api/sessions` | 세션 생성 |
| PUT | `/api/sessions/:id/end` | 세션 종료 |
| POST | `/api/observations` | 관찰 저장 |
| POST | `/api/summaries` | 요약 저장 |
| POST | `/api/plans` | 플랜 저장 |
| GET | `/api/team/context` | 팀 컨텍스트 조회 |

## 환경 변수

```bash
# 서버 포트 (기본: 3778)
DO_WORKER_PORT=3778

# 데이터베이스 설정
DO_DB_TYPE=sqlite          # sqlite 또는 mysql
DO_DB_PATH=.do/memory.db   # SQLite 경로

# MySQL 설정 (DO_DB_TYPE=mysql인 경우)
DO_DB_HOST=localhost
DO_DB_PORT=3306
DO_DB_USER=doworker
DO_DB_PASSWORD=secret
DO_DB_DATABASE=do_memory

# 사용자 설정
DO_USER_NAME=max
```

## 사용 예시

### 세션 생성

```bash
curl -X POST http://localhost:3778/api/sessions \
  -H "Content-Type: application/json" \
  -d '{"id": "session-123", "user_name": "max"}'
```

### 관찰 저장

```bash
curl -X POST http://localhost:3778/api/observations \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "session-123",
    "agent_name": "expert-backend",
    "type": "decision",
    "content": "PostgreSQL 선택: 트랜잭션 지원 필요",
    "importance": 4,
    "tags": ["database", "architecture"]
  }'
```

### 컨텍스트 조회

```bash
curl "http://localhost:3778/api/context/inject?user=max"
```

### 팀 컨텍스트 조회

```bash
curl "http://localhost:3778/api/team/context?exclude_user=max"
```

## 아키텍처

```
.do/worker/
├── main.go              # 진입점
├── cmd/worker/main.go   # CLI
├── internal/
│   ├── server/          # HTTP 서버
│   ├── db/              # 데이터베이스 어댑터
│   ├── memory/          # 메모리 관리
│   └── context/         # 컨텍스트 빌더
└── pkg/models/          # 공유 타입
```

## 개발

```bash
# 테스트 실행
make test

# 커버리지 리포트
make test-cover

# 린트
make lint

# 포맷
make fmt
```

## 라이선스

MIT
