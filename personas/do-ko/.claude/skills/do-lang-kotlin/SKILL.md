---
name: do-lang-kotlin
description: >
  Kotlin 2.0+ 개발 전문가 - Ktor, coroutines, Compose Multiplatform,
  Kotlin-idiomatic 패턴涵盖. Kotlin 서버 앱, Android 앱, 멀티플랫폼
  프로젝트 구축 시 사용하세요.
license: Apache-2.0
compatibility: Designed for Claude Code
user-invocable: false
metadata:
  version: "1.1.0"
  category: "language"
  status: "active"
  updated: "2026-01-11"
  modularized: "false"
  tags: "kotlin, ktor, coroutines, compose, android, multiplatform"
  context7-libraries: "/ktorio/ktor, /jetbrains/compose-multiplatform, /jetbrains/exposed"
  related-skills: "do-lang-java, do-lang-swift"

# MoAI Extension: Progressive Disclosure
progressive_disclosure:
  enabled: true
  level1_tokens: 100
  level2_tokens: 5000

# MoAI Extension: Triggers
triggers:
  keywords: ["Kotlin", "Ktor", "coroutine", "Flow", "Compose", "Android", ".kt", ".kts", "build.gradle.kts"]
  languages: ["kotlin"]
---

## Quick Reference (30초 요약)

K2 컴파일러, coroutines, Ktor, Compose Multiplatform이 포함된 Kotlin 2.0+ 전문가.

자동 트리거: Kotlin 파일 (.kt, .kts), Gradle Kotlin DSL (build.gradle.kts, settings.gradle.kts)

핵심 기능:

- Kotlin 2.0: K2 컴파일러, coroutines, Flow, sealed classes, value classes
- Ktor 3.0: Async HTTP server/client, WebSocket, JWT 인증
- Exposed 0.55: coroutines 지원이 있는 Kotlin SQL 프레임워크
- Spring Boot (Kotlin): Kotlin-idiomatic Spring과 WebFlux
- Compose Multiplatform: Desktop, iOS, Web, Android UI
- 테스트: JUnit 5, MockK, Kotest, Flow 테스트용 Turbine

---

## Implementation Guide (5분 가이드)

### Kotlin 2.0 기능

Coroutines 및 Flow:

parallel 작업을 위해 coroutineScope와 async를 사용하세요. async로 deferred 값을 생성한 다음 각각을 await하여 결과를 가져오세요. data class로 결과를 결합하세요. reactive 스트림의 경우 while 루프 내에서 emit 호출이 포함된 flow 블록을 생성하세요. 간격은 delay를 사용하고 flowOn으로 디스패처를 지정하세요.

Sealed Classes 및 Value Classes:

generic 타입 파라미터가 있는 sealed interface를 정의하세요. success를 위한 data class 구현, 상태 없는 케이스(예: Loading)를 위한 data object를 생성하세요. @JvmInline과 원시 타입을 래핑하는 value class를 사용하세요. init 블록에서 검증을 위해 require를 사용하세요.

### Ktor 3.0 Server

Application Setup:

Netty, port, host 파라미터와 함께 embeddedServer를 호출하세요. 람다 내에서 Koin, security, routing, content negotiation 구성 함수를 호출하세요. wait이 true인 start를 호출하세요.

Koin 구성의 경우 Koin 플러그인을 설치하고 single 선언으로 싱글톤을 위한 모듈을 정의하세요. security의 경우 Authentication 플러그인을 설치하고 realm, verifier, validate 콜백으로 JWT를 구성하세요. content negotiation의 경우 json 구성과 함께 ContentNegotiation 플러그인을 설치하세요.

인증과 함께 Routing:

Application에서 routing 함수를 정의하세요. routing 블록 내에서 route를 경로 접두사로 사용하세요. create 엔드포인트를 정의하기 위해 call.receive로 요청 본문을 받는 post를 사용하세요. verifier 이름과 함께 authenticate 블록으로 보호된 경로를 정의하세요. route 블록 내에서 get 엔드포인트를 정의하고 call.parameters로 경로/쿼리 파라미터를 가져옵니다. call.respond와 함께 상태 코드 및 응답 본문을 반환하세요.

### Exposed SQL Framework

Table 및 Entity:

LongIdTable을 확장하고 테이블 이름을 가진 object를 정의하세요. varchar, enumerationByName, timestamp 함수로 열을 선언하세요. uniqueIndex()와 defaultExpression을 기본값으로 사용하세요.

LongEntity와 LongEntityClass를 확장하는 entity class를 생성하세요. 테이블 열 참조를 사용하는 by 구문으로 속성을 선언하세요. entity를 도메인 모델로 매핑하는 toModel 함수를 생성하세요.

Coroutines와 함께 Repository:

Database 파라미터를 받는 repository 구현을 생성하세요. dbQuery 헬퍼로 래핑된 Exposed 작업을 래핑하는 suspend 함수를 구현하세요. findById는 단일 엔티티 조회를 사용하세요. Entity.new는 삽입을 위해 사용하세요. IO 디스패처를 사용하는 newSuspendedTransaction으로 private dbQuery 함수를 정의하세요.

### Spring Boot와 함께 Kotlin

WebFlux Controller:

@RestController와 @RequestMapping로 클래스에 어노테이션을 적용하세요. @GetMapping과 @PostMapping이 있는 suspend 함수로 엔드포인트를 정의하세요. map으로 엔티티를 변환하는 Flow를 반환하기 위해 collect을 사용하세요. 상태 코드와 함께 ResponseEntity를 반환하세요. 요청 검증을 위해 @Valid를 사용하세요.

---

## Advanced Patterns

### Compose Multiplatform

공유 UI 컴포넌트:

ViewModel과 콜백 파라미터를 받는 @Composable 함수를 생성하세요. collectAsState로 uiState를 상태로 수집하세요. sealed 상태에서 when 식을 사용하여 Loading, Success, Error에 대해 다른 composable를 표시합니다.

목록 항목의 경우 Modifier.fillMaxWidth와 clickable이 있는 Card composable를 생성하세요. 패딩이 있는 Row, CircleShape로 클립된 AsyncImage for avatar, MaterialTheme.typography로 텍스트 내용을 위한 Column을 사용하세요.

### MockK와 함께 테스트

mockk로 의존성을 위한 테스트 클래스를 생성하세요. 선언에서 mock으로 mock을 초기화하고 서비스를 생성하세요. coroutine 테스트를 위해 runTest와 @Test를 사용하세요. async mocking을 위해 coEvery와 coAnswers를 delay와 함께 사용하세요. assertions를 위해 assertThat를 사용하세요. Flow 테스트의 경우 toList로 방출을 수집하고 크기와 내용을 assert하세요.

### Gradle Build Configuration

kotlin("jvm")과 kotlin("plugin.serialization") 버전 문자열과 함께 plugins 블록을 사용하세요. ktor.plugin id를 추가하세요. kotlin 블록에 jvmToolchain을 구성하세요. dependencies 블록에 ktor server 모듈, kotlinx coroutines, exposed 모듈, postgresql driver를 위한 implementation dependencies를 추가하세요. mockk, coroutines-test, turbine을 위한 test dependencies를 추가하세요.

---

## Context7 통합

최신 문서를 위한 라이브러리 매핑:

- `/ktorio/ktor` - Ktor 3.0 server/client 문서
- `/jetbrains/exposed` - Exposed SQL 프레임워크
- `/JetBrains/kotlin` - Kotlin 2.0 언어 참조
- `/Kotlin/kotlinx.coroutines` - Coroutines 라이브러리
- `/jetbrains/compose-multiplatform` - Compose Multiplatform
- `/arrow-kt/arrow` - Arrow 함수형 프로그래밍

사용법: context7CompatibleLibraryID, 특정 영역을 위한 topic 문자열, 응답 크기를 파라미터로 mcp__context7__get_library_docs를 호출하세요.

---

## Kotlin을 언제 사용할까요

다음의 경우 Kotlin을 사용하세요:

- Android 애플리케이션 개발 (공식 언어)
- Ktor로 모던 서버 애플리케이션 구축
- 간결하고 표현적인 구문을 선호
- coroutines와 Flow로 반응형 서비스 구축
- iOS, Desktop, Web을 위한 멀티플랫폼 애플리케이션 생성
- 완전한 Java 상호 운용성 필요

대안을 고려해야 할 때:

- 최소 변경을 필요로 하는 레거시 Java 코드베이스
- 빅 데이터 파이프라인 (Scala with Spark를 선호)

---

## Works Well With

- `do-lang-java` - Java 상호 운용성 및 Spring Boot 패턴
- `do-domain-backend` - REST API, GraphQL, 마이크로서비스 아키텍처
- `do-domain-database` - JPA, Exposed, R2DBC 패턴
- `do-quality-testing` - JUnit 5, MockK, TestContainers 통합
- `do-infra-docker` - JVM 컨테이너 최적화

---

## Troubleshooting

K2 컴파일러: gradle.properties에 kotlin.experimental.tryK2=true를 추가하세요. 전체 재빌드를 위해 .gradle 디렉토리를 지우세요.

Coroutines: suspend 컨텍스트에서 runBlocking을 피하세요. blocking 작업에는 Dispatchers.IO를 사용하세요.

Ktor: ContentNegotiation이 설치되어 있는지 확인하세요. JWT verifier 구성을 확인하세요. 라우팅� 계층 구조를 확인하세요.

Exposed: 모든 DB 작업이 트랜잭션 컨텍스트 내에서 실행되는지 확인하세요. 트랜잭션 외부에서 lazy entity 로딩을 주의하세요.

---

## Advanced Documentation

포괄적인 참조 자료:

- [reference.md](reference.md) - 완전한 에코시스템, Context7 매핑, 테스트 패턴, 성능
- [examples.md](examples.md) - 프로덕션 준비 코드 예제, Ktor, Compose, Android 패턴

---

Last Updated: 2026-01-11
Status: Production Ready (v1.1.0)
