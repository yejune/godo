---
name: do-lang-python
description: >
  Python 3.13+ 개발 전문가 - FastAPI, Django, async 패턴, 데이터 사이언스,
  pytest 테스트, 모던 Python 기능涵盖. Python API 개발, 웹 애플리케이션,
  데이터 파이프라인, 테스트 작성 시 사용하세요.
license: Apache-2.0
compatibility: Designed for Claude Code
allowed-tools: Read Grep Glob Bash(python:*) Bash(python3:*) Bash(pytest:*) Bash(ruff:*) Bash(pip:*) Bash(uv:*) Bash(mypy:*) Bash(pyright:*) Bash(black:*) Bash(poetry:*) mcp__context7__resolve-library-id mcp__context7__get-library-docs
user-invocable: false
metadata:
  version: "1.1.0"
  category: "language"
  status: "active"
  updated: "2026-01-11"
  modularized: "false"
  tags: "language, python, fastapi, django, pytest, async, data-science"

# MoAI Extension: Progressive Disclosure
progressive_disclosure:
  enabled: true
  level1_tokens: 100
  level2_tokens: 5000

# MoAI Extension: Triggers
triggers:
  keywords: ["Python", "Django", "FastAPI", "Flask", "asyncio", "pytest", "pyproject.toml", "requirements.txt", ".py"]
  languages: ["python"]
---

## Quick Reference (30초 요약)

Python 3.13+ 개발 전문가 - FastAPI, Django, async 패턴, pytest, 모던 Python 기능.

자동 트리거: .py 확장자 파일, pyproject.toml, requirements.txt, pytest.ini, FastAPI 또는 Django 관련 논의

핵심 기능:

- Python 3.13 기능: PEP 744 JIT 컴파일러, PEP 703 GIL-free 모드, match/case 문 패턴 매칭
- 웹 프레임워크: FastAPI 0.115 이상, Django 5.2 LTS
- 데이터 검증: Pydantic v2.9 model_validate 패턴
- ORM: SQLAlchemy 2.0 async 패턴
- 테스트: pytest fixtures, async 테스트, parametrize 데코레이터
- 패키지 관리: poetry, uv, pip, pyproject.toml
- 타입 힌트: Protocol, TypeVar, ParamSpec, 모던 typing 패턴
- Async: asyncio, async 제너레이터, 태스크 그룹
- 데이터 사이언스: numpy, pandas, polars 기초

### Quick Patterns

FastAPI 엔드포인트 패턴:

fastapi에서 FastAPI와 Depends를, pydantic에서 BaseModel을 임포트하세요. FastAPI 애플리케이션 인스턴스를 생성하세요. BaseModel을 상속받고 name과 email 문자열 필드를 가진 UserCreate 모델 클래스를 정의하세요. users 경로에 UserCreate 파라미터를 받고 await UserService.create를 호출하여 User를 반환하는 async post 엔드포인트를 생성하세요.

Pydantic v2.9 검증 패턴:

pydantic에서 BaseModel과 ConfigDict를 임포트하세요. BaseModel을 상속받는 User 클래스를 정의하세요. from_attributes를 True로, str_strip_whitespace를 True로 설정하여 model_config를 ConfigDict로 설정하세요. id를 정수, name을 문자열, email을 문자열 필드로 추가하세요. ORM 객체에서 model_validate를 사용하여 생성하고, JSON 데이터에서 model_validate_json을 사용하여 생성하세요.

pytest Async 테스트 패턴:

pytest와 pytest.mark.asyncio 데코레이터를 임포트하세요. async_client fixture 파라미터를 받는 async 테스트 함수를 생성하세요. name 필드가 포함된 JSON 본문으로 users 엔드포인트에 post 요청을 보내세요. response status_code가 201과 같은지 단언하세요.

---

## Implementation Guide (5분 가이드)

### Python 3.13 새 기능

PEP 744 JIT 컴파일러:

- 기본적으로 비활성화된 실험적 기능
- PYTHON_JIT 환경변수를 1로 설정하여 활성화
- enable-experimental-jit 플래그로 빌드 옵션 제공
- CPU 바운드 코드에 성능 향상 제공
- 전문 바이트코드를 머신 코드로 변환하는 copy-and-patch JIT 사용

PEP 703 GIL-Free 모드:

- python3.13t로 실험적 free-threaded 빌드 제공
- 실제 병렬 스레드 실행 허용
- 공식 Windows 및 macOS 설치 프로그램에서 사용 가능
- CPU 집약적 멀티스레드 애플리케이션에 가장 적합
- 아직 프로덕션 사용 권장 안 함

match/case 패턴 매칭:

response 디셔너리를 받아 문자열을 반환하는 process_response 함수를 생성하세요. response에 match 문을 사용하세요. status가 ok이고 data 필드가 있는 케이스의 경우 data와 함께 성공 메시지를 반환하세요. status가 error이고 message 필드가 있는 케이스의 경우 에러 메시지를 반환하세요. guard 조건을 사용하여 status가 pending 또는 processing과 일치하는 케이스의 경우 진행 중 메시지를 반환하세요. 언더스코어를 사용하는 기본 케이스의 경우 알 수 없는 응답을 반환하세요.

### FastAPI 0.115+ 패턴

Async 의존성 주입:

fastapi에서 FastAPI, Depends를, sqlalchemy.ext.asyncio에서 AsyncSession을, contextlib에서 asynccontextmanager를 임포트하세요. asynccontextmanager로 데코레이팅되고 FastAPI app을 받는 lifespan async 컨텍스트 매니저를 생성하세요. lifespan에서 시작 시 await init_db를 호출하고, yield한 후 종료 시 await cleanup을 호출하세요. lifespan 파라미터로 FastAPI app을 생성하세요. async_session에서 async with를 사용하여 AsyncSession의 AsyncGenerator를 반환하는 async get_db 함수를 정의하세요. get_db와 함께 Depends를 사용하여 DB 세션을 주입하는 user_id 경로 파라미터로 get 엔드포인트를 생성하세요. await get_user_by_id를 호출하고 UserResponse.model_validate와 user를 반환하세요.

클래스 기반 의존성:

page 기본값 1, size 기본값 20을 받는 init 메서드로 Paginator 클래스를 생성하세요. self.page를 1과 page 중 최대값으로, self.size를 100과 1과 size 중 최대값의 최소값으로, self.offset을 page에서 1을 뺀 값에 size를 곱한 값으로 설정하세요. Paginator에 Depends를 사용하는 list_items 엔드포인트를 생성하여 페이지네이션을 주입하고 offset과 size로 get_page를 사용하여 항목을 반환하세요.

### Django 5.2 LTS 기능

복합 기본 키:

CASCADE 삭제로 Order에 대한 ForeignKey, CASCADE 삭제로 Product에 대한 ForeignKey, 수량을 위한 IntegerField를 가진 OrderItem 모델을 생성하세요. Meta 클래스에서 pk를 order와 product 필드를 가진 models.CompositePrimaryKey로 설정하세요.

쿼리 파라미터가 있는 URL 역변환:

django.urls에서 reverse를 임포트하세요. search 뷰 이름, q를 django로, page를 1로 설정한 쿼리 디셔너리, results로 설정한 fragment와 함께 reverse를 호출하세요. 결과는 쿼리 문자열과 fragment가 있는 검색 경로입니다.

셸의 자동 모델 임포트:

python manage.py shell을 실행하면 설치된 모든 앱의 모델이 명시적 임포트 없이 자동으로 임포트됩니다.

### Pydantic v2.9 심화 패턴

Annotated와 함께 재사용 가능한 검증기:

typing에서 Annotated를, pydantic에서 AfterValidator와 BaseModel을 임포트하세요. 정수 v를 받아 정수를 반환하는 validate_positive 함수를 정의하세요. v가 0 이하이면 "must be positive" 메시지와 함께 ValueError를 발생시키세요. 그렇지 않으면 v를 반환하세요. AfterValidator와 validate_positive를 사용하는 Annotated로 PositiveInt를 생성하세요. 가격과 수량 필드에 PositiveInt를 사용하세요.

교차 필드 검증을 위한 모델 검증기:

pydantic에서 BaseModel과 model_validator를, typing에서 Self를 임포트하세요. start_date와 end_date를 date 필드로 가진 DateRange 모델을 생성하세요. mode를 after로 설정하여 model_validator 데코레이터를 추가하세요. Self를 반환하는 validate_dates 메서드에서 end_date가 start_date보다 앞서 있는지 확인하고 그렇다면 ValueError를 발생시키고, 그렇지 않으면 self를 반환하세요.

ConfigDict 모범 사례:

model_config를 ConfigDict로 설정한 BaseSchema 모델을 생성하세요. ORM 객체 지원을 위해 from_attributes를 True로, 별칭 허용을 위해 populate_by_name을 True로, 알 수 없는 필드 실패를 위해 extra를 forbid로, 문자열 정리를 위해 str_strip_whitespace를 True로 설정하세요.

### SQLAlchemy 2.0 Async 패턴

엔진 및 세션 설정:

sqlalchemy.ext.asyncio에서 create_async_engine, async_sessionmaker, AsyncSession을 임포트하세요. postgresql+asyncpg 연결 문자열, pool_pre_ping을 True로, echo를 True로 설정하여 create_async_engine으로 엔진을 생성하세요. engine, class_를 AsyncSession으로, expire_on_commit을 False로 설정하여 async_session을 생성하여 분리된 인스턴스 오류를 방지하세요.

리포지토리 패턴:

AsyncSession을 받는 init 메서드를 가진 UserRepository 클래스를 생성하세요. user_id에 대한 where 절이 있는 select 쿼리를 실행하여 scalar_one_or_none 결과를 반환하는 async get_by_id 메서드를 정의하세요. UserCreate model_dump에서 User를 생성하고, 세션에 추가하고, 커밋하고, 새로고침하고, user를 반환하는 async create 메서드를 정의하세요.

대량 결과 스트리밍:

AsyncSession을 받는 async stream_users 함수를 생성하세요. select User 쿼리로 await db.stream을 호출하세요. async for를 사용하여 result.scalars를 반복하고 각 user를 yield하세요.

### pytest 고급 패턴

pytest-asyncio와 함께 Async fixtures:

pytest, pytest_asyncio, httpx에서 AsyncClient를 임포트하세요. pytest_asyncio.fixture로 fixtures를 데코레이트하세요. app과 base_url로 가진 AsyncClient에 async with를 사용하여 async_client fixture를 생성하세요. async_session과 session.begin에 async with를 사용하여 db_session fixture를 생성하고 session을 yield한 후 await session.rollback을 호출하세요.

매개변수화된 테스트:

input_data와 expected_status 파라미터 이름과 함께 pytest.mark.parametrize 데코레이터를 사용하세요. 딕셔너리와 예상 상태 코드가 포함된 튜플로 테스트 케이스를 제공하세요. valid, empty_name, missing_name 케이스에 대한 id를 추가하세요. 테스트 함수는 async_client, input_data, expected_status를 받고 users 엔드포인트에 post하고 status_code가 expected와 일치하는지 단언하세요.

Fixture 팩토리:

async 함수를 반환하는 user_factory fixture를 생성하세요. 내부 함수는 db를 AsyncSession으로 받고 키워드 인자를 받습니다. name과 email으로 기본값 딕셔너리를 설정하세요. 파이프 연산자를 사용하여 kwargs와 병합된 기본값으로 User를 생성하고, db에 추가하고, 커밋하고, user를 반환하세요.

### 타입 힌트 모던 패턴

구조적 타이핑을 위한 Protocol:

typing에서 Protocol과 runtime_checkable을 임포트하세요. runtime_checkable 데코레이터를 적용하세요. 제네릭 타입 T로 Repository 프로토콜을 정의하세요. int id를 받아 T 또는 None을 반환하는 abstract async get 메서드, dict data를 받아 T를 반환하는 async create 메서드, int id를 받아 bool을 반환하는 async delete 메서드를 추가하세요.

데코레이터를 위한 ParamSpec:

typing에서 ParamSpec, TypeVar, Callable을, functools에서 wraps를 임포트하세요. P를 ParamSpec으로, R을 TypeVar로 정의하세요. times 기본값 3을 받고 callable 래퍼를 반환하는 retry 데코레이터 함수를 생성하세요. 내부 데코레이터는 함수를 래핑하고 래퍼는 지정된 횟수만큼 반복하며 함수를 await 시도하고 마지막 시도에서 재발생합니다.

### 패키지 관리

Poetry와 함께 pyproject.toml:

tool.poetry 섹션에서 name, version, python 버전 제약을 설정하세요. dependencies 아래에 fastapi, pydantic, asyncio extras가 포함된 sqlalchemy를 추가하세요. dev dependencies 아래에 pytest, pytest-asyncio, ruff를 추가하세요. line-length와 target-version으로 ruff를 구성하세요. ini_options에서 pytest asyncio_mode를 auto로 설정하세요.

uv 빠른 패키지 관리자:

curl을 사용하여 astral.sh의 설치 스크립트로 uv를 설치하세요. uv venv로 가상 환경을 생성하세요. uv pip install로 requirements.txt에서 의존성을 설치하세요. uv add 명령으로 의존성을 추가하세요.

---

## Advanced Implementation (10분 이상)

포괄적인涵盖内容包括:

- Docker 및 Kubernetes용 프로덕션 배포 패턴
- 태스크 그룹 및 세마포어를 포함한 고급 async 패턴
- numpy, pandas, polars와의 데이터 사이언스 통합
- 성능 최적화 기법
- OWASP 패턴을 따르는 보안 모범 사례
- CI/CD 통합 패턴

다음을 참조하세요:

- reference.md 완전한 참조 문서
- examples.md 프로덕션 준비 코드 예제

---

## Context7 라이브러리 매핑

- tiangolo/fastapi FastAPI async 웹 프레임워크
- django/django Django 웹 프레임워크
- pydantic/pydantic 타입 어노테이션과 데이터 검증
- sqlalchemy/sqlalchemy SQL 툴킷 및 ORM
- pytest-dev/pytest 테스팅 프레임워크
- numpy/numpy 수치 계산
- pandas-dev/pandas 데이터 분석 라이브러리
- pola-rs/polars 빠른 DataFrame 라이브러리

---

## Works Well With

- do-domain-backend REST API 및 마이크로서비스 아키텍처
- do-domain-database SQL 패턴 및 ORM 최적화
- do-workflow-testing DDD 및 테스트 전략
- do-essentials-debug AI 기반 디버깅
- do-foundation-quality TRUST 5 품질 원칙

---

## Troubleshooting

일반적인 문제:

Python 버전 확인:

python에 version 플래그를 사용하여 3.13 이상인지 확인하세요. 상세한 버전 정보를 위해 -c 플래그와 함께 sys.version_info를 출력하는 python을 사용하세요.

Async 세션 분리 오류:

세션 구성에서 expire_on_commit을 False로 설정하세요. 또는 커밋 후 객체와 함께 await session.refresh를 사용하세요.

pytest asyncio 모드 경고:

pyproject.toml의 tool.pytest.ini_options 아래에 asyncio_mode를 auto로, asyncio_default_fixture_loop_scope를 function으로 설정하세요.

Pydantic v2 마이그레이션:

parse_obj 메서드는 이제 model_validate입니다. parse_raw 메서드는 이제 model_validate_json입니다. from_orm 기능은 ConfigDict에서 from_attributes를 True로 설정해야 합니다.

---

Last Updated: 2026-01-11
Status: Active (v1.1.0)
