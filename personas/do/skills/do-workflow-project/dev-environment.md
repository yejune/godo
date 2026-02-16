# Development Environment Rules [HARD]

## Docker 필수

### 기본 원칙
- [HARD] 모든 프로젝트는 반드시 Dockerized (docker-compose.yml 필수)
- [HARD] Dockerfile 생성 지양 — 공식 이미지(`node:20`, `python:3.13`, `golang:1.23` 등) + docker-compose `command:`로 구성
- [HARD] Dockerfile이 필요한 경우: 시스템 패키지 설치, 멀티스테이지 빌드 등 공식 이미지만으로 불가능할 때만
- [HARD] 모든 코드는 Docker 컨테이너 내에서 실행 — 명령어는 컨테이너 대상
- [HARD] Docker Compose가 로컬 개발 환경의 Single Source of Truth
- [HARD] 호스트와 컨테이너는 볼륨 마운트로 연결 — 코드 변경은 즉시 컨테이너에 반영

### 명령 실행 구분
- [HARD] 실행 중인 컨테이너에 명령: `docker compose exec <service> <command>`
- [HARD] 일회성 명령 (마이그레이션, 시드 등): `docker compose run --rm <service> <command>`
- [HARD] 명령 실행 전 `docker compose ps`로 컨테이너 상태 확인 — 중지 상태면 먼저 `docker compose up -d`
- [HARD] healthcheck 정의된 서비스는 healthy 상태 확인 후 테스트/마이그레이션 실행

### 네트워크 — bootapp 도메인 기반
- [HARD] 포트 외부 노출 불필요 — bootapp이 도메인 기반 라우팅 제공
- [HARD] docker-compose.yml에 `ports:` 섹션 생략 가능 — 도메인으로 접근
- [HARD] 서비스 접근은 DOMAIN 환경변수로 설정 (예: `DOMAIN: app.test`, `DOMAIN: db.test`)
- [HARD] SSL_DOMAINS 환경변수 필수 — 포트 없이 HTTPS 도메인으로 통신, 인증서 자동 생성/신뢰
- [HARD] 호스트에서 접근: `https://app.test` (NOT `localhost:8080`)
- [HARD] 컨테이너 간 통신: Docker 서비스명 또는 DOMAIN 값 사용 (예: `db.test`, `redis.test`)
- [HARD] 컨테이너 내부 코드에서 `localhost`로 다른 서비스 접근 금지
- [HARD] TLD는 `.test` 권장 (RFC 2606 예약) — `.local` 금지 (macOS mDNS 충돌)

### 빌드 & 재시작
- [HARD] `docker bootapp up/down`은 `/etc/hosts`를 수정하므로 에이전트가 직접 실행 금지 — 사용자에게 실행할 명령어를 알려주고 실행을 요청할 것 (AskUserQuestion)
- [HARD] 프로젝트 시작: 사용자에게 `docker bootapp up` 실행 요청 — 서브넷 할당, 도메인 등록, SSL 자동 처리
- [HARD] 프로젝트 종료: 사용자에게 `docker bootapp down` 실행 요청 — 도메인/hosts 정리
- [HARD] Dockerfile 변경 시: `docker compose build` 후 사용자에게 `docker bootapp up` 실행 요청
- [HARD] docker-compose.yml 변경 시: 사용자에게 `docker bootapp up` 실행 요청 (자동 반영)
- [HARD] bootapp up 완료 후 반드시 서비스 상태 확인: `docker compose ps`로 컨테이너 기동 확인 + healthcheck 통과 대기
- [HARD] 의존성 추가 시 (package.json, go.mod 등): 컨테이너 내에서 설치 또는 이미지 리빌드

## 환경변수 관리

### 12-Factor 원칙
- [HARD] 설정은 환경변수로 주입 — 코드에 하드코딩 절대 금지
- [HARD] 커넥션 스트링, API URL 등은 반드시 환경변수에서 읽기
- [HARD] 기본값은 컨테이너 환경 기준 (예: `DB_HOST=postgres`, NOT `DB_HOST=localhost`)

### docker-compose.yml 환경변수
- [HARD] 일반 환경변수는 docker-compose.yml의 `environment` 섹션에 직접 정의
- [HARD] 서비스 관심사별 그룹핑 — DB 관련, 캐시 관련, 앱 설정 등 주석으로 구분
- [HARD] 시크릿(AWS 키, DB 비밀번호 등)만 별도 분리 — Docker secrets 또는 외부 주입(AWS SSM, Vault 등)

### 금지 사항
- [HARD] `.env` 자동 로드 파일 생성 절대 금지 (docker-compose가 암묵적으로 읽는 `.env`)
- [HARD] `.env.local`, `.env.development`, `.env.production` 파일 생성 금지

### 시크릿 주입 (env_file 허용)
- env_file: 디렉티브는 **시크릿 전용**으로 허용 — 로컬 개발 시 시크릿 대체용
- 로컬: `env_file:` → 프로덕션: AWS SSM / Vault 등 외부 주입 — 앱 코드는 동일 (멱등)
- env_file 대상 파일은 반드시 `.gitignore`에 등록 — 커밋 금지
- 일반 환경변수는 여전히 `environment:` 섹션에 직접 정의

## 코드 품질

### Read before Write 원칙
- [HARD] 코드 작성 전 기존 패턴 파악 필수 — 파일 구조, 네이밍, 에러 핸들링 스타일 확인
- [HARD] 새 코드는 프로젝트 기존 컨벤션을 따름 — 근거 없이 새로운 패턴 도입 금지
- [HARD] 유사한 유틸리티/헬퍼가 이미 있으면 재사용 — 중복 파일 생성 금지

### 구문 검사 필수
- [HARD] 코드 작성/수정 후 반드시 언어별 구문 검사 실행:
  - Go: `go build ./...` 또는 `go vet ./...`
  - TypeScript/JS: `npm run lint` 또는 `npx tsc --noEmit`
  - Rust: `cargo check`
  - Python: `ruff check` 또는 `flake8`
- [HARD] 구문 검사는 컨테이너 우선 — 프로덕션 이미지 등 도구 미설치 시 호스트에서 실행 허용 (볼륨 마운트로 소스 공유 전제)

### 의존성 관리
- [HARD] 새 의존성 추가 전 기존 의존성으로 해결 가능한지 확인
- [HARD] 기존 지식과 경험을 먼저 활용 — 검색하고 참고, 새로운 발견은 문서화

## AI 에이전트 금지 패턴

- [HARD] 의존성 설치는 호스트에서 실행 — 볼륨 마운트로 컨테이너에 자동 반영 (`npm install`, `pip install`, `go mod download` 등)
- [HARD] 구문 검사/개발 도구도 호스트에 플랫폼별로 설치 허용 (lint, formatter, type checker 등)
- [HARD] 컨테이너 간 통신에 `localhost` 사용 금지 — Docker 서비스명 또는 도메인 사용
- [HARD] `ports:` 매핑으로 포트 노출하지 않음 — bootapp 도메인 사용
- [HARD] Docker Compose healthcheck 상태 무시 금지 — healthy 확인 후 후속 작업
- [HARD] `.env` 계열 파일 생성 금지 — 어떤 이유로도 허용하지 않음
- [HARD] 컨테이너 밖에서 테스트 실행 금지 — `docker compose exec`로 실행
- [HARD] 컨테이너 내부 셸 진입 금지 (`docker exec -it ... bash/sh`) — 반드시 `docker compose exec <service> <command>` 단발 명령으로 실행
- [HARD] `docker cp` 금지 — 컨테이너에 파일 복사하지 않음, 볼륨 마운트로 공유
- [HARD] 컨테이너 내부에 임시 파일/스크립트 생성 금지 — 모든 변경은 호스트 소스에서, 볼륨으로 반영
- [HARD] 멱등성 필수 — 같은 명령을 여러 번 실행해도 동일한 결과, 임시 상태에 의존하지 않음
- [HARD] Dockerfile에서 소스 코드 `COPY` 금지 — 소스와 의존성 모두 볼륨 마운트로 공유
