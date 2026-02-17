---
name: do-lang-rust
description: >
  Rust 1.92+ 개발 전문가 - Axum, Tokio, SQLx, 메모리 안전 시스템 프로그래밍涵盖.
  고성능, 메모리 안전 애플리케이션 또는 WebAssembly 구축 시 사용하세요.
license: Apache-2.0
compatibility: Designed for Claude Code
allowed-tools: Read Grep Glob mcp__context7__resolve-library-id mcp__context7__get-library-docs
user-invocable: false
metadata:
  version: "1.2.0"
  category: "language"
  status: "active"
  updated: "2026-01-11"
  modularized: "false"
  tags: "language, rust, axum, tokio, sqlx, serde, wasm, cargo"

# MoAI Extension: Progressive Disclosure
progressive_disclosure:
  enabled: true
  level1_tokens: 100
  level2_tokens: 5000

# MoAI Extension: Triggers
triggers:
  keywords: ["Rust", "Axum", "Tokio", "SQLx", "serde", ".rs", "Cargo.toml", "async", "await", "lifetime", "trait", "ownership", "borrowing", "performance", "optimization", "clippy", "memory safety"]
  languages: ["rust"]
---

## Quick Reference (30초 요약)

고성능, 메모리 안전 애플리케이션을 위한 Rust 1.92+ 개발 전문가입니다.

자동 트리거: .rs, Cargo.toml, async/await, Tokio, Axum, SQLx, serde, 라이프타임, 트레이트

핵심 사용 사례:

- 고성능 REST API 및 마이크로서비스
- 메모리 안전 동시 시스템
- CLI 도구 및 시스템 유틸리티
- WebAssembly 애플리케이션
- 저지연시 네트워킹 서비스

빠른 패턴:

Axum REST API: 경로와 핸들러를 연결하는 route 매크로 체이닝으로 Router를 생성하세요. 공유 상태를 위해 with_state를 추가하세요. tokio::net와 함께 TcpListener를 바인딩하고 axum::serve로 서비스하세요.

Async Handler와 SQLx: AppState를 위한 State extractor와 id를 위한 Path extractor를 받는 async handler 함수를 정의하세요. SQL 문자열과 바인드 파라미터로 sqlx::query_as! 매크로를 사용하세요. pool에서 fetch_optional을 호출하고 await한 다음 에러 변환을 위해 ok_or를 사용하세요. Json로 래핑된 결과를 반환하세요.

---

## Implementation Guide (5분 가이드)

### Rust 1.92 기능

모던 Rust 기능:

- Rust 2024 Edition 사용 가능 (Rust 1.85와 함께 출시)
- 안정적인 async 트레이트 (async-trait 크레이트 더 이상 필요 없음)
- 컴파일 시간 배열 크기를 위한 Const 제네릭
- 초기 반환을 위한 let-else 패턴 매칭
- polonius로 향상된 borrow checker

안정적인 Async 트레이트: async fn 서명을 가진 트레이트를 정의하세요. async fn 구현으로 구체적 타입에 트레이트를 구현하세요. 트레이트 메서드에서 직접 sqlx 매크로를 호출하세요.

Let-Else 패턴: return을 사용하여 초기 종료를 위해 let Some(value) = option else를 사용하세요. 순차적 검증을 위해 여러 let-else 문을 체이닝하세요. else 블록에서 에러 타입을 반환하세요.

### 웹 프레임워크 Axum 0.8

설치: Cargo.toml dependencies 섹션에 axum 버전 0.8, full features가 있는 tokio 버전 1.48, cors와 trace features가 있는 tower-http 버전 0.6을 추가하세요.

완전한 API 설정: axum::extract에서 extractor를, routing macros를 임포트하세요. PgPool을 보유하는 Clone-derive AppState 구조체를 정의하세요. tokio::main async main에서 DATABASE_URL에서 env와 함께 max_connections를 설정하여 PgPoolOptions로 pool을 생성하세요. 경로와 핸들러에 대한 route 체인, CorsLayer를 가진 Router를 빌드하고 with_state를 호출하세요. TcpListener를 바인딩하고 axum::serve를 호출하세요.

핸들러 패턴: State, Path, Query extractor를 적절한 타입과 함께 받는 async 핸들러를 정의하세요. 위치 바인드를 위한 타입 안전 쿼리로 sqlx::query_as!를 사용하세요. Json 성공과 AppError 실패를 가진 Result를 반환하세요.

### Async 런타임 Tokio 1.48

태스크 생성 및 채널: 용량으로 mpsc 채널을 생성하세요. 루프에서 채널에서 수신하는 worker 태스크를 생성하기 위해 tokio::spawn을 호출하세요. timeout의 경우 작업 분기와 sleep 분기가 있는 tokio::select! 매크로를 사용하고 timeout 시 에러를 반환하세요.

### 데이터베이스 SQLx 0.8

타입 안전 쿼리: 자동 매핑을 위해 구조체에 sqlx::FromRow를 파생하세요. 컴파일 시간 체크 쿼리로 query_as! 매크로를 사용하세요. pool에서 fetch_one 또는 fetch_optional을 호출하세요. 트랜잭션의 경우 pool.begin을 호출하고 트랜잭션 참조에서 쿼리를 실행한 다음 tx.commit을 호출하세요.

### 직렬화 Serde 1.0

구조체에 Serialize와 Deserialize를 파생하세요. 대소문자 변환을 위해 rename_all이 있는 serde 속성을 사용하세요. 필드별 명명을 위해 rename 속성을 사용하세요. Option::is_none과 함께 skip_serializing_if를 사용하세요. 기본값을 위해 default 속성을 사용하세요.

### 에러 처리

thiserror: display 메시지를 위한 error 속성이 있는 enum에 Error를 파생하세요. 소스 에러에서 자동 변환을 위해 from 속성을 사용하세요. 변형에서 일치하여 상태 코드와 Json 본문이 포함된 에러 메시지를 반환하여 IntoResponse를 구현하세요.

### CLI 개발 clap

command 속성으로 name, version, about를 가진 main Cli 구조체에 Parser를 파생하세요. 전역 플래그를 위해 arg 속성을 사용하세요. 명령을 위해 enum에 Subcommand를 파생하세요. main에서 명령과 일치하여 논리를 디스패치하세요.

### 테스트 패턴

cfg(test) 속성으로 테스트 모듈을 생성하세요. tokio::test async 함수를 정의하세요. 설정 도우미를 호출하고 테스트할 함수를 호출한 다음 검증을 위해 assert! 매크로를 사용하세요.

---

## Advanced Patterns

포괄적인涵盖内容包括:

- 고급 소유권 패턴, 라이프타임, 스마트 포인터
- 트레이트 설계 및 제네릭 프로그래밍
- 성능 최적화 전략 및 프로파일링
- 엔지니어링 모범 사례 및 코딩 가이드라인
- Async 패턴 및 동시성

다음을 참조하세요: 소유권 및 트레이트는 [reference/engineering.md](reference/engineering.md), 최적화는 [reference/performance.md](reference/performance.md), 코딩 표준은 [reference/guidelines.md](reference/guidelines.md)

### 성능 최적화

릴리스 빌드: Cargo.toml profile.release 섹션에서 lto를 활성화하고, codegen-units를 1로 설정하고, panic을 abort로 설정하고, strip을 활성화하세요.

### 배포

최소 컨테이너: 멀티 스테이지 Dockerfile을 사용하세요. 첫 번째 스테이지는 rust alpine 이미지를 사용하고 Cargo 파일을 복사하고 의존성 캐싱을 위해 더미 main을 생성하고 릴리스를 빌드하고 소스를 복사하고 재빌드를 위해 main.rs를 touch하고 최종 릴리스를 빌드합니다. 두 번째 스테이지는 alpine을 사용하고 빌더에서 바이너리를 복사하고 포트를 노출하고 CMD를 설정합니다.

### 동시성

속도 제한 작업: 최대 허가를 가진 Semaphore로 래핑된 Arc를 생성하세요. 항목을 매핑하여 permit을 획득하고 처리하고 결과를 반환하는 태스크를 생성합니다. futures::future::join_all을 사용하여 결과를 수집하세요.

---

## Context7 통합

라이브러리 문서 액세스:

- `/rust-lang/rust` - Rust 언어 및 stdlib
- `/tokio-rs/tokio` - Tokio async 런타임
- `/tokio-rs/axum` - Axum 웹 프레임워크
- `/launchbadge/sqlx` - SQLx async SQL
- `/serde-rs/serde` - 직렬화 프레임워크
- `/dtolnay/thiserror` - 에러 derive
- `/clap-rs/clap` - CLI parser

---

## Works Well With

- `do-lang-go` - Go 시스템 프로그래밍 패턴
- `do-domain-backend` - REST API 아키텍처 및 마이크로서비스 패턴
- `do-foundation-quality` - Rust 애플리케이션 보안 강화
- `do-workflow-testing` - 테스트 주도 개발 워크플로우

---

## Troubleshooting

일반적인 문제:

- Cargo 에러: cargo clean 다음에 cargo build를 실행하세요
- 버전 확인: rustc --version과 cargo --version을 실행하세요
- 의존성 문제: cargo update와 cargo tree를 실행하세요
- 컴파일 시간 SQL 체크: cargo sqlx prepare를 실행하세요

성능 특성:

- 시작 시간: 50-100ms
- 메모리 사용량: 5-20MB 기본
- 처리량: 100k-200k req/s
- 지연시간: p99 5ms 미만
- 컨테이너 크기: 5-15MB (alpine)

---

## Additional Resources

고급 소유권 패턴과 트레이트 설계는 reference/engineering.md를 참조하세요.

최적화 전략 및 프로파일링 기법은 reference/performance.md를 참조하세요.

Rust 코딩 표준 및 모범 사례는 reference/guidelines.md를 참조하세요.

---

Last Updated: 2026-01-11
Version: 1.2.0
