---
name: do-lang-go
description: >
  Go 1.23+ 개발 전문가 - Fiber, Gin, GORM, 동시성 프로그래밍 패턴涵盖.
  고성능 마이크로서비스, CLI 도구, 클라우드 네이티브 애플리케이션 구축 시 사용하세요.
license: Apache-2.0
compatibility: Designed for Claude Code
user-invocable: false
metadata:
  version: "1.1.0"
  category: "language"
  status: "active"
  updated: "2026-01-11"
  modularized: "false"
  tags: "go, golang, fiber, gin, concurrency, microservices"
  context7-libraries: "/gofiber/fiber, /gin-gonic/gin, /go-gorm/gorm"
  related-skills: "do-lang-rust, do-domain-backend"

# MoAI Extension: Progressive Disclosure
progressive_disclosure:
  enabled: true
  level1_tokens: 100
  level2_tokens: 5000

# MoAI Extension: Triggers
triggers:
  keywords: ["Go", "Golang", "Fiber", "Gin", "GORM", "Echo", "Chi", ".go", "go.mod", "goroutine", "channel", "generics", "concurrent", "testing", "benchmark", "fuzzing", "microservices", "gRPC"]
  languages: ["go", "golang"]
---

## Quick Reference (30초 요약)

고성능 백엔드 시스템 및 CLI 애플리케이션을 위한 Go 1.23+ 개발 전문가.

자동 트리거: .go 확장자 파일, go.mod, go.sum, goroutines, channels, Fiber, Gin, GORM, Echo, Chi

핵심 사용 사례:

- 고성능 REST API 및 마이크로서비스
- 동시 및 병렬 처리 시스템
- CLI 도구 및 시스템 유틸리티
- 클라우드 네이티브 컨테이너화된 서비스

빠른 패턴:

Fiber API 패턴:

fiber.New() 함수를 호출하여 app를 생성하세요. fiber.Ctx를 받고 error를 반환하는 핸들러 함수로 api/users/:id에 get 경로를 정의하세요. 핸들러에서 c.Params에서 id를 가져와 fiber.Map과 함께 c.JSON을 호출하세요. 포트 3000으로 app.Listen을 호출하세요.

Gin API 패턴:

gin.Default() 함수를 호출하여 r을 생성하세요. gin.Context 포인터를 받는 핸들러 함수로 api/users/:id에 GET 경로를 정의하세요. 핸들러에서 c.Param에서 id를 가져와 상태 200과 gin.H와 함께 c.JSON을 호출하세요. 포트 3000으로 r.Run을 호출하세요.

에러 처리와 함께 Goroutine:

context.Background()와 함께 errgroup.WithContext를 호출하여 g와 ctx를 생성하세요. processUsers(ctx)를 반환하는 함수로 g.Go를 호출하세요. processOrders(ctx)를 반환하는 함수로 g.Go를 호출하세요. g.Wait()에서 에러가 nil이 아니면 log.Fatal과 함께 에러를 호출하세요.

---

## Implementation Guide (5분 가이드)

### Go 1.23 언어 기능

새로운 기능:

- for i range 10 구문으로 정수 범위에서 i를 반복하고 i를 출력하세요
- Profile-Guided Optimization PGO 2.0
- 더 나은 타입 추론으로 향상된 제네릭

제네릭 패턴:

T와 U를 any로 하는 타입 파라미터로 제네릭 Map 함수를 생성하세요. T 슬라이스와 T에서 U로의 함수를 받으세요. 동일한 길이의 U 슬라이스로 result를 생성하세요. values에 함수를 적용하여 result 요소를 설정하며 슬라이스를 반복하세요. result를 반환하세요.

### 웹 프레임워크 Fiber v3

ErrorHandler와 Prefork를 true로 설정한 fiber.Config와 함께 fiber.New로 app를 생성하세요. recover.New, logger.New, cors.New 미들웨어를 사용하세요. api/v1 경로에 api 그룹을 생성하세요. listUsers, id 파라미터로 getUser, id로 updateUser, id로 deleteUser 경로를 정의하세요. 포트 3000으로 app.Listen을 호출하세요.

### 웹 프레임워크 Gin

gin.Default로 r을 생성하세요. cors.Default 미들웨어를 사용하세요. api/v1 경로에 api 그룹을 생성하세요. listUsers를 호출하는 GET users, getUser를 호출하는 GET users/:id, createUser를 호출하는 POST users를 정의하세요. 포트 3000으로 r.Run을 호출하세요.

요청 바인딩 패턴:

Name과 Email 필드를 가진 CreateUserRequest 구조체를 정의하세요. json 태그와 required, min 길이 2, required email 검증을 위한 binding 태그를 추가하세요. createUser 핸들러에서 req 변수를 선언하고 포인터와 함께 c.ShouldBindJSON을 호출하세요. 에러가 있으면 상태 400과 에러와 함께 c.JSON을 호출하세요. 그렇지 않으면 201과 응답 데이터와 함께 c.JSON을 호출하세요.

### 웹 프레임워크 Echo

echo.New로 e를 생성하세요. middleware.Logger, middleware.Recover, middleware.CORS를 사용하세요. api/v1 경로에 api 그룹을 생성하세요. GET users와 POST users를 정의하세요. 포트 3000으로 e.Start와 함께 e.Logger.Fatal을 호출하세요.

### 웹 프레임워크 Chi

chi.NewRouter로 r을 생성하세요. middleware.Logger와 middleware.Recoverer를 사용하세요. api/v1 경로와 함수로 r.Route를 호출하세요. 내부에서 users 경로로 r.Route를 호출하세요. 목록을 위한 Get, 생성을 위한 Post, 단일 사용자를 위한 id 파라미터가 있는 Get을 정의하세요. 포트 3000으로 r과 함께 http.ListenAndServe를 호출하세요.

### ORM GORM 1.25

모델 정의:

gorm.Model을 임베딩하는 User 구조체를 정의하세요. uniqueIndex와 not null 태그가 있는 Name, uniqueIndex와 not null이 있는 Email, foreignKey AuthorID 태그가 있는 Posts 슬라이스를 추가하세요.

쿼리 패턴:

created_at desc로 정렬하고 10으로 제한하는 Posts로 db.Preload를 호출하고 id 1로 First를 호출하세요. 트랜잭션의 경우 tx 포인터를 받는 함수로 db.Transaction을 호출하세요. 내부에서 user와 profile을 생성하고 모든 에러를 반환하세요.

### sqlc와 함께 타입 안전 SQL

sqlc 2 버전, postgresql 엔진이 있는 sql 섹션, 쿼리 및 스키마 경로, 패키지 이름, 출력 디렉토리, pgx v5 sql_package가 있는 go 생성 설정으로 sqlc.yaml을 생성하세요.

query.sql 파일에서 id가 파라미터와 일치하는 곳에서 모든 열을 반환하는 GetUser를 이름으로 추가하세요. name과 email 값을 삽입하고 모든 열을 반환하는 CreateUser를 이름으로 추가하세요.

### 동시성 패턴

Errgroup 패턴:

errgroup.WithContext로 g와 ctx를 생성하세요. users 변수에 할당하는 fetchUsers를 위해 g.Go를 호출하세요. orders 변수에 할당하는 fetchOrders를 위해 g.Go를 호출하세요. g.Wait가 에러를 반환하면 nil과 에러를 반환하세요.

Worker Pool 패턴:

jobs receive-only 채널, results send-only 채널, n worker count를 받는 workerPool 함수를 정의하세요. WaitGroup을 생성하세요. n번 반복하며 WaitGroup을 증가하고 Done을 지연시키며 jobs를 반복하고 processJob 결과를 results로 보내는 goroutine을 생성하세요. Wait한 다음 results를 닫으세요.

Timeout과 함께 Context:

5초 동안 context.WithTimeout으로 ctx와 cancel을 생성하세요. defer로 cancel 호출하세요. fetchData와 함께 ctx를 호출하세요. 에러가 context.DeadlineExceeded이면 timeout과 StatusGatewayTimeout으로 응답하세요.

### 테스트 패턴

테이블 기반 테스트:

name 문자열, input CreateUserInput, wantErr bool을 포함하는 구조체의 tests 슬라이스를 정의하세요. 유효한 입력과 빈 이름에 대한 테스트 케이스를 추가하세요. name과 테스트 함수로 t.Run을 호출하며 tests를 반복하세요. service Create를 호출하고 wantErr가 true인지 require.Error를 확인하세요.

HTTP 테스트:

fiber.New로 app를 생성하세요. users/:id에서 getUser를 호출하는 GET 경로를 추가하세요. GET at users/1으로 httptest.NewRequest로 요청을 생성하세요. app.Test와 요청을 호출하여 응답을 가져오세요. 200 상태 코드를 단언하세요.

### Cobra와 Viper와 함께 CLI

Use와 Short 필드가 있는 cobra.Command 포인터로 rootCmd를 정의하세요. init 함수에서 cfgFile을 위한 PersistentFlags StringVar를 추가하세요. config와 lookup으로 viper.BindPFlag를 호출하세요. MYAPP으로 viper.SetEnvPrefix를 설정하고 viper.AutomaticEnv를 호출하세요.

---

## Advanced Patterns

포괄적인涵盖内容包括:

- 고급 동시성 패턴 (worker pools, rate limiting, errgroup)
- 제네릭 및 타입 제약 조건
- 인터페이스 설계 및 컴포지션
- 포괄적인 테스트 패턴 (TDD, 테이블 기반, 벤치마크, 퍼징)
- 성능 최적화 및 프로파일링

다음을 참조하세요: 고급 패턴은 [reference/advanced.md](reference/advanced.md), 테스트 패턴은 [reference/testing.md](reference/testing.md)

### 성능 최적화

PGO 빌드:

GODEBUG=pgo를 활성화하고 cpuprofile 출력으로 실행하세요. pgo 플래그를 프로필 파일을 가리키며 go build로 빌드하세요.

객체 풀링:

4096 바이트 슬라이스를 반환하는 New 함수를 가진 sync.Pool로 bufferPool을 생성하세요. 타입 어설션으로 버퍼를 가져오고, 풀로 반환하기 위해 Put을 지연하세요.

### 컨테이너 배포 10-20MB

멀티 스테이지 Dockerfile: 첫 번째 스테이지는 golang:1.23-alpine을 builder로 사용하고, WORKDIR을 설정하고, go.mod와 go.sum을 복사하고, go mod download를 실행하고, 소스를 복사하고, CGO_ENABLED 0와 스트립된 바이너리를 위한 ldflags로 빌드합니다. 두 번째 스테이지는 scratch를 사용하고, 빌더에서 바이너리를 복사하고, ENTRYPOINT를 설정합니다.

### 우아한 종료

app.Listen을 호출하는 goroutine을 생성하세요. 버퍼 1로 os.Signal용 quit 채널을 생성하세요. SIGINT 및 SIGTERM에 대해 signal.Notify를 호출하세요. quit에서 수신한 다음 app.Shutdown을 호출하세요.

---

## Context7 라이브러리

- golang/go Go 언어 및 stdlib
- gofiber/fiber Fiber 웹 프레임워크
- gin-gonic/gin Gin 웹 프레임워크
- labstack/echo Echo 웹 프레임워크
- go-chi/chi Chi 라우터
- go-gorm/gorm GORM ORM
- sqlc-dev/sqlc 타입 안전 SQL
- jackc/pgx PostgreSQL 드라이버
- spf13/cobra CLI 프레임워크
- spf13/viper 구성
- stretchr/testify 테스팅 툴킷

---

## Works Well With

- do-domain-backend REST API 아키텍처 및 마이크로서비스
- do-lang-rust 시스템 프로그래밍 동반자
- do-quality-security 보안 강화
- do-essentials-debug 성능 프로파일링
- do-workflow-ddd 도메인 주도 개발

---

## Troubleshooting

일반적인 문제:

- 모듈 에러: go mod tidy와 go mod verify를 실행하세요
- 버전 확인: go version과 go env GOVERSION을 실행하세요
- 빌드 문제: go clean -cache와 go build -v를 실행하세요

성능 진단:

- CPU 프로파일링: go test -cpuprofile cpu.prof -bench .를 실행하세요
- 메모리 프로파일링: go test -memprofile mem.prof -bench .를 실행하세요
- 경합 감지: go test -race ./...를 실행하세요

---

## Additional Resources

고급 동시성 패턴, 제네릭, 인터페이스 설계는 reference/advanced.md를 참조하세요.

TDD, 벤치마크, 퍼징을 포함한 포괄적인 테스트 패턴은 reference/testing.md를 참조하세요.

---

Last Updated: 2026-01-11
Version: 1.1.0
