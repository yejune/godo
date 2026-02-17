---
name: do-lang-scala
description: >
  Scala 3.4+ 개발 전문가 - Akka, Cats Effect, ZIO, Spark 패턴涵盖.
  분산 시스템, 빅 데이터 파이프라인, 함수형 프로그래밍 애플리케이션 구축 시
  사용하세요.
license: Apache-2.0
compatibility: Designed for Claude Code
allowed-tools: Read Grep Glob mcp__context7__resolve-library-id mcp__context7__get-library-docs
user-invocable: false
metadata:
  version: "2.1.0"
  category: "language"
  status: "active"
  updated: "2026-01-11"
  modularized: "true"
  tags: "scala, akka, cats-effect, zio, spark, sbt"

# MoAI Extension: Progressive Disclosure
progressive_disclosure:
  enabled: true
  level1_tokens: 100
  level2_tokens: 5000

# MoAI Extension: Triggers
triggers:
  keywords: ["Scala", "Akka", "Cats Effect", "ZIO", "Spark", ".scala", ".sc", "build.sbt", "sbt"]
  languages: ["scala"]
---

# Scala 3.4+ 개발 전문가

JVM 애플리케이션을 위한 함수형 프로그래밍, effect systems, 빅 데이터 처리.

## Quick Reference

자동 트리거: Scala 파일 (.scala, .sc), 빌드 파일 (build.sbt, project/build.properties)

핵심 기능:

- Scala 3.4: Given/using, extension 메서드, enums, opaque types, match types
- Akka 2.9: Typed actors, streams, clustering, persistence
- Cats Effect 3.5: Pure FP 런타임, fibers, concurrent 구조
- ZIO 2.1: Effect system, layers, streaming, error handling
- Apache Spark 3.5: DataFrame API, SQL, structured streaming

주요 에코시스템 라이브러리:

- HTTP: Http4s 0.24, Tapir 1.10
- JSON: Circe 0.15, ZIO JSON 0.6
- Database: Doobie 1.0, Slick 3.5, Quill 4.8
- Streaming: FS2 3.10, ZIO Streams 2.1
- Testing: ScalaTest, Specs2, MUnit, Weaver

---

## Module Index

이 스킬은 전문涵盖을 위해 점진적 공개를 사용하는 특수화된 모듈을 사용합니다:

### Core Language

- [functional-programming.md](modules/functional-programming.md) - Scala 3.4 기능: Given/Using, Type Classes, Enums, Opaque Types, Extension Methods

### Effect Systems

- [cats-effect.md](modules/cats-effect.md) - Cats Effect 3.5: IO monad, Resources, Fibers, FS2 Streaming
- [zio-patterns.md](modules/zio-patterns.md) - ZIO 2.1: Effects, Layers, ZIO Streams, Error handling

### Frameworks

- [akka-actors.md](modules/akka-actors.md) - Akka Typed Actors 2.9: Actors, Streams, Clustering patterns
- [spark-data.md](modules/spark-data.md) - Apache Spark 3.5: DataFrame API, SQL, Structured Streaming

---

## Implementation Guide

### Project Setup (SBT 1.10)

build.sbt에서 ThisBuild / scalaVersion을 "3.4.2"로 설정하고 organization을 정의하세요. name과 libraryDependencies를 포함하는 settings를 포함하는 lazy val root 프로젝트를 정의하세요. cats-effect, zio, akka-actor-typed, http4s-ember-server, circe-generic 의존성과 test scope를 위한 scalatest를 추가하세요. deprecation, feature warnings, Xfatal-warnings를 위한 scalacOptions를 포함하세요.

### Quick Examples

Extension Methods: 괄호 안 파라미터와 함께 extension 키워드를 사용하세요. 공백으로 분할하고 길이를 확인한 후 문자를 가져오고 말줄표를 추가하는 truncate 같은 메서드를 정의하세요.

Given과 Using: 추상 메서드 서명을 가진 trait를 정의하세요. with 키워드와 함께 given 인스턴스를 생성하고 메서드를 구현하세요. using 절을 위한 implicit 해결을 위한 함수를 생성하세요.

Enum 타입: 제네릭 타입 파라미터와 plus 분산 주석을 가진 enum을 정의하세요. 파라미터를 가진 case 엔트리를 생성하세요. match 식을 사용하여 각 case를 처리하는 enum에 메서드를 정의하고 적절한 결과를 반환하세요.

---

## Context7 통합

최신 문서를 위한 라이브러리 매핑:

Core Scala:

- /scala/scala3 - Scala 3.4 언어 참조
- /scala/scala-library - 표준 라이브러리

Effect Systems:

- /typelevel/cats-effect - Cats Effect 3.5 문서
- /typelevel/cats - Cats 2.10 함수형 추상화
- /zio/zio - ZIO 2.1 문서
- /zio/zio-streams - ZIO Streams 2.1

Akka Ecosystem:

- /akka/akka - Akka 2.9 typed actors 및 streams
- /akka/akka-http - Akka HTTP REST API
- /akka/alpakka - Akka 커넥터

HTTP and Web:

- /http4s/http4s - 함수형 HTTP 서버/클라이언트
- /softwaremill/tapir - API 우선 설계

Big Data:

- /apache/spark - Spark 3.5 DataFrame 및 SQL
- /apache/flink - Flink 1.19 streaming
- /apache/kafka - Kafka clients 3.7

---

## Testing Quick Reference

ScalaTest: Matchers와 함께 AnyFlatSpec를 확장하세요. should in을 사용하는 문자열 설명으로 동작을 정의하세요. 같음 비교를 위해 shouldBe를 사용하세요.

Cats Effect와 함께 MUnit: CatsEffectSuite를 확장하세요. IO가 포함된 assertEquals assertion을 반환하는 string name으로 테스트를 정의하세요.

ZIO Test: ZIOSpecDefault를 확장하세요. spec을 suite로 하고 test 엔트리를 정의하세요. for-comprehension을 사용하여 effect를 실행하고 assertTrue assertion을 yield하세요.

---

## Troubleshooting

일반적인 문제:

- Implicit 해결: 상세 에러 메시지를 위해 scalac -explain 사용
- 타입 추론: 추론 실패 시 명시적 타입 어노테이션 추가
- SBT 컴파일 느림: build.sbt에서 Global / concurrentRestrictions 활성화

Effect System 문제:

- Cats Effect: cats.effect._ 또는 cats.syntax.all._ import가 누락되었는지 확인
- ZIO: ZIO.serviceWith과 ZIO.serviceWithZIO로 layer 구성 확인
- Akka: actor 계층 구조 및 supervision 전략 검토

---

## Works Well With

- do-lang-java JVM 상호 운용성, Spring Boot 통합
- do-domain-backend REST API, GraphQL, 마이크로서비스 패턴
- do-domain-database Doobie, Slick, 데이터베이스 패턴
- do-workflow-testing ScalaTest, MUnit, property-based testing

---

## Additional Resources

포괄적인 참조 자료:

- [reference.md](reference.md) - 완전한 Scala 3.4涵盖, Context7 매핑, 성능
- [examples.md](examples.md) - 프로덕션 준비 코드: Http4s, Akka, Spark 패턴

---

Last Updated: 2026-01-11
Status: Production Ready (v2.1.0)
