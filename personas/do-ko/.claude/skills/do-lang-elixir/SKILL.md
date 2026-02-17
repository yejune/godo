---
name: do-lang-elixir
description: >
  Elixir 1.17+ 개발 전문가 - Phoenix 1.7, LiveView, Ecto, OTP 패턴涵盖.
  실시간 애플리케이션, 분산 시스템, Phoenix 프로젝트 구축 시 사용하세요.
license: Apache-2.0
compatibility: Designed for Claude Code
allowed-tools: Read Grep Glob Bash(mix:*) Bash(elixir:*) Bash(iex:*) Bash(erl:*) mcp__context7__resolve-library-id mcp__context7__get-library-docs
user-invocable: false
metadata:
  version: "1.1.0"
  category: "language"
  status: "active"
  updated: "2026-01-11"
  modularized: "true"
  tags: "language, elixir, phoenix, liveview, ecto, otp, genserver"

# MoAI Extension: Progressive Disclosure
progressive_disclosure:
  enabled: true
  level1_tokens: 100
  level2_tokens: 5000

# MoAI Extension: Triggers
triggers:
  keywords: ["Elixir", "Phoenix", "LiveView", "Ecto", "OTP", "GenServer", ".ex", ".exs", "mix.exs"]
  languages: ["elixir"]
---

## Quick Reference (30초 요약)

Elixir 1.17+ 개발 전문가 - Phoenix 1.7, LiveView, Ecto, OTP 패턴, 함수형 프로그래밍.

자동 트리거: .ex, .exs 파일, mix.exs, config/, Phoenix/LiveView 논의

핵심 기능:

- Elixir 1.17: 패턴 매칭, 파이프, 프로토콜, 행동, 매크로
- Phoenix 1.7: Controller, LiveView, Channels, PubSub, Verified Routes
- Ecto: Schema, Changeset, Query, 마이그레이션, Multi
- OTP: GenServer, Supervisor, Agent, Task, Registry
- ExUnit: setup, describe, async와 함께 테스트
- Mix: 빌드 도구, 작업, 릴리스
- Oban: 백그라운드 작업 처리

### Quick Patterns

Phoenix Controller: MyAppWeb에 :controller를 사용하여 모듈을 정의하세요. MyApp.Accounts와 같이 context 모듈에 대한 alias를 생성하세요. conn과 params map을 받고 destructuring한 id를 가진 show와 같은 action 함수를 정의하세요. context 함수와 bang(!)을 사용하여 get_user!로 데이터를 가져오고 데이터를 assigns로 전달하는 템플릿을 렌더링하세요.

create 액션의 경우 context 결과 튜플에 패턴 매칭하세요. ok 튜플의 경우 파이프 연산자를 사용하여 put_flash로 성공 메시지를 보내고 ~p sigil로 verified routes를 사용하여 redirect 하세요. Ecto.Changeset이 있는 error 튜플의 경우 changeset을 전달하는 폼 템플릿을 렌더링하세요.

Ecto Schema와 Changeset: Ecto.Schema와 Ecto.Changeset을 임포트하는 모듈을 정의하세요. :string과 같은 타입과 virtual 필드를 포함하는 필드 선언이 있는 schema 블록을 정의하세요. struct와 attrs를 받는 changeset 함수를 생성하고 파이프 연산자로 cast 대상 필드 목록과 체이닝, validate_required, regex를 사용하는 validate_format, min 옵션이 있는 validate_length, unique_constraint를 체인하세요.

GenServer 패턴: GenServer를 사용하여 모듈을 정의하세요. initial_value를 받고 GenServer.start_link를 __MODULE__, initial_value, name 옵션과 함께 호출하는 start_link를 생성하세요. GenServer.call과 __MODULE__ 및 message atom을 호출하는 클라이언트 API 함수를 생성하세요. 상태를 반환하는 ok 튜플로 init 콜백을 구현하세요. 각 message에 대한 handle_call 콜백을 구현하여 응답과 새 상태가 포함된 reply 튜플을 반환하세요.

---

## Implementation Guide (5분 가이드)

### Elixir 1.17 기능

고급 패턴 매칭: map 키와 타입에 패턴 매칭하는 함수 헤드를 정의하세요. is_binary 또는 byte_size 확인 같은 제약 조건을 추가하기 위해 when 가드를 사용하세요. 잘못된 입력을 위한 error 튜플을 반환하는 catch-all 절을 추가하세요.

에러 처리와 함께 파이프 연산자: 실패할 수 있는 작업을 연결하기 위해 with 특수 형식을 사용하세요. 왼쪽 화살표로 각 단계를 패턴 매치하고 모든 단계가 성공적으로 완료되면 최종 ok 튜플을 반환하세요. else 블록에서 error 튜플을 그대로 반환하여 처리하세요.

다형성을 위한 프로토콜: defprotocol과 @doc를 사용하여 프로토콜과 함수 사양을 정의하세요. for: 옵션과 함께 defimpl을 사용하여 특정 타입에 대해 프로토콜을 구현하세요. 각 구현은 해당 타입에 대한 특정 동작을 제공합니다.

### Phoenix 1.7 패턴

LiveView 컴포넌트: :live_view와 함께 MyAppWeb 모듈을 정의하세요. params, session, socket을 받는 mount 콜백을 구현하여 상태를 assigned한 ok 튜플을 반환하세요. 사용자 상호작용을 위한 handle_event 콜백을 구현하여 update helper로 socket을 업데이트한 noreply 튜플을 반환하세요. @assign를 접두사로 사용하여 assigns에 접근하는 ~H sigil이 있는 HEEx 템플릿을 사용하는 render 콜백을 구현하세요.

Changeset이 있는 LiveView 폼: mount에서 to_form helper로 초기 changeset을 생성하고 form을 할당하세요. :validate 액션으로 changeset를 생성하고 form을 재할당하는 validate event handler를 구현하세요. context create 함수를 호출하는 save event handler를 구현하고 성공 시 put_flash와 push_navigate를 사용하고 error changeset으로 form을 재할당하세요.

Phoenix Channels: :channel과 함께 MyAppWeb 모듈을 정의하세요. 대괄호 각진 세그먼트가 있는 topic 패턴과 일치하도록 join 콜백을 구현하세요. after_join 메시지를 보내기 위해 self()와 함께 send를 사용하세요. broadcast!를 사용하여 모든 구독자에게 보내기 위해 handle_info를 구현하세요. broadcast!를 사용하여 handle_in에서 클라이언트 메시지를 처리하세요.

Verified Routes: router.ex의 scope 블록에서 라이브 경로를 위해 live 매크로 경로를 정의하세요. 템플릿과 컨트롤러에서 동적 세그먼트를 위한 보간 구문이 있는 ~p sigil을 사용하여 verified routes를 사용하세요.

### Ecto 패턴

트랜잭션을 위한 Multi: Ecto.Multi.new()를 사용하고 name atom과 changeset 함수로 Ecto.Multi.update를 파이프하여 작업을 체인하세요. 이전 단계에서 결과가 필요한 경우 function 콜백과 함께 Ecto.Multi.insert를 사용하세요. 최종 Multi를 Repo.transaction()에 파이프하여 명명된 결과가 포함된 ok 또는 error 튜플을 반환합니다.

쿼리 구성: composable 쿼리 함수를 가진 query 모듈을 생성하세요. from 표현식을 반환하는 base 함수를 정의하세요. where 절이 포함된 from 표현식을 반환하는 filter 함수를 query, default 파라미터를 가진 채우기를 정의하세요. Repo.all에 전달하기 전에 파이프 연산자로 함수를 체인하세요.

---

## Advanced Implementation (10분 이상)

포괄적인涵盖内容包括:

- 릴리스와 함께 프로덕션 배포
- libcluster를 사용한 분산 시스템
- 고급 LiveView 패턴 (streams, components)
- OTP supervision 트리 및 동적 supervisor
- 원격 측정 및 관찰 가능성
- 보안 모범 사례
- CI/CD 통합 패턴

다음을 참조하세요:

- [Advanced Patterns](modules/advanced-patterns.md) - 완전한 고급 패턴 가이드

---

## Context7 라이브러리 매핑

- /elixir-lang/elixir Elixir 언어 문서
- /phoenixframework/phoenix Phoenix 웹 프레임워크
- /phoenixframework/phoenix_live_view LiveView 실시간 UI
- /elixir-ecto/ecto 데이터베이스 래퍼 및 쿼리 언어
- /sorentwo/oban 백그라운드 작업 처리

---

## Works Well With

- `do-domain-backend` - REST API 및 마이크로서비스 아키텍처
- `do-domain-database` - SQL 패턴 및 쿼리 최적화
- `do-workflow-testing` - DDD 및 테스트 전략
- `do-essentials-debug` - AI 기반 디버깅
- `do-platform-deploy` - 배포 및 인프라

---

## Troubleshooting

일반적인 문제:

Elixir 버전 확인: 1.17+을 확인하려면 elixir --version을 실행하세요. Mix 빌드 도구 버전을 확인하려면 mix --version을 실행하세요.

의존성 문제: 의존성을 가져오려면 mix deps.get을 실행하세요. 컴파일하려면 mix deps.compile을 실행하세요. 빌드 아티팩트를 제거하려면 mix clean을 실행하세요.

데이터베이스 마이그레이션: mix ecto.create로 데이터베이스를 생성하세요. 마이그레이션을 실행하려면 mix ecto.migrate를 실행하세요. 마지막 마이그레이션을 롤백하려면 mix ecto.rollback을 실행하세요.

Phoenix 서버: mix phx.server로 서버를 시작하세요. IEx 셸로 시작하려면 iex -S mix phx.server를 실행하세요. 프로덕션 릴리스를 빌드하려면 MIX_ENV=prod mix release를 실행하세요.

LiveView 로딩 안 됨:

- 브라우저 개발자 콘솔에서 웹소켓 연결을 확인하세요
- 엔드포인트 구성에 websocket transport가 포함되어 있는지 확인하세요
- Phoenix.LiveView가 mix.exs 의존성에 나열되어 있는지 확인하세요

---

Last Updated: 2026-01-11
Status: Active (v1.1.0)
