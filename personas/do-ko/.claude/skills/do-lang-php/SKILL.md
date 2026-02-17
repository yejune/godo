---
name: do-lang-php
description: >
  PHP 8.3+ 개발 전문가 - Laravel 11, Symfony 7, Eloquent ORM, 모던 PHP 패턴涵盖.
  PHP API, 웹 애플리케이션, Laravel/Symfony 프로젝트 개발 시 사용하세요.
license: Apache-2.0
compatibility: Designed for Claude Code
allowed-tools: Read Grep Glob Bash(php:*) Bash(composer:*) Bash(phpunit:*) Bash(phpstan:*) Bash(phpcs:*) Bash(artisan:*) Bash(laravel:*) mcp__context7__resolve-library-id mcp__context7__get-library-docs
user-invocable: false
metadata:
  version: "1.1.0"
  category: "language"
  status: "active"
  updated: "2026-01-11"
  modularized: "true"
  tags: "language, php, laravel, symfony, eloquent, doctrine, phpunit"

# MoAI Extension: Progressive Disclosure
progressive_disclosure:
  enabled: true
  level1_tokens: 100
  level2_tokens: 5000

# MoAI Extension: Triggers
triggers:
  keywords: ["PHP", "Laravel", "Symfony", "Eloquent", "Doctrine", "PHPUnit", "Pest", ".php", "composer.json", "artisan"]
  languages: ["php"]
---

## Quick Reference (30초 요약)

PHP 8.3+ 개발 전문가 - Laravel 11, Symfony 7, Eloquent, Doctrine, 모던 PHP 패턴.

자동 트리거: .php 확장자 파일, composer.json, artisan 명령, symfony.yaml, Laravel 또는 Symfony 논의

핵심 기능:

- PHP 8.3 기능: readonly 클래스, typed 속성, 속성, enums, named arguments
- Laravel 11: Controller, Model, 마이그레이션, Form Request, API Resource, Eloquent
- Symfony 7: 속성 기반 라우팅, Doctrine ORM, 서비스, 의존성 주입
- ORM: Laravel용 Eloquent, Symfony용 Doctrine
- 테스트: PHPUnit, Pest, feature 및 unit 테스트 패턴
- 패키지 관리: autoloading이 있는 Composer
- 코딩 표준: PSR-12, Laravel Pint, PHP CS Fixer
- Docker: PHP-FPM, nginx, 멀티 스테이지 빌드

### Quick Patterns

Laravel Controller 패턴:

App\Http\Controllers\Api 네임스페이스에 Controller를 확장하는 UserController를 생성하세요. StoreUserRequest, UserResource, User, JsonResponse를 임포트하세요. StoreUserRequest를 받는 store 메서드를 정의하고 validated 데이터로 User를 생성하며 User를 래핑하는 UserResource와 상태 201인 JsonResponse를 반환하세요.

Laravel Form Request 패턴:

App\Http\Requests 네임스페이스에 FormRequest를 확장하는 StoreUserRequest를 생성하세요. authorize 메서드는 true를 반환합니다. rules 메서드는 name에 required, string, max 255 검증, email에 required, email, users 테이블의 unique 검증, password에 required, min 8, confirmed 검증이 있는 배열을 반환합니다.

Symfony Controller 패턴:

App\Controller 네임스페이스에 AbstractController를 확장하는 UserController를 생성하세요. User, EntityManagerInterface, JsonResponse, Route 속성을 임포트하세요. api/users 경로를 위한 클래스 레벨에 Route 속성을 적용하세요. 빈 경로와 POST 메서드를 위한 create 메서드에 Route 속성을 적용하세요. EntityManagerInterface를 주입하고, 새 User를 생성하고 persist하고 flush한 다음 user와 상태 201로 json 응답을 반환하세요.

---

## Implementation Guide (5분 가이드)

### PHP 8.3 모던 기능

Readonly Classes:

int id, string name, string email 속성을 public으로 촉진하는 readonly class UserDTO를 선언하세요.

Enums with Methods:

string backed enum OrderStatus를 정의하세요. pending 값, processing 값, completed 값인 케이스를 정의하세요. $this에 대한 match 식을 사용하여 각 케이스에 적절한 디스플레이 레이블을 반환하는 label 메서드를 추가하세요.

Attributes:

Attribute 속성과 함께 Validate 속성 클래스를 생성하세요. 생성자는 string rule과 선택적 string message를 받습니다. email 속성이 required 및 email 규칙을 지정하는 Validate 속성과 함께 UserRequest 클래스를 생성하세요.

### Laravel 11 패턴

Relationships와 함께 Eloquent Model:

App\Models 네임스페이스에 Model을 확장하는 Post model을 생성하세요. title, content, user_id, status를 가진 protected fillable 배열을 설정하세요. status를 PostStatus enum으로, published_at을 datetime으로 캐스팅하는 protected casts 배열을 설정하세요. BelongsTo 관계를 반환하는 user 메서드를 정의하세요. HasMany 관계를 반환하는 comments 메서드를 정의하세요. published status로 필터링하는 scopePublished 메서드를 추가하세요.

API Resource 패턴:

App\Http\Resources 네임스페이스에 JsonResource를 확장하는 PostResource를 생성하세요. Request 파라미터를 받는 toArray 메서드를 정의하세요. id, title, whenLoaded로 user 관계를 위한 UserResource, whenCounted로 comments_count, ISO 8601 문자열로 formatted created_at을 포함하는 배열을 반환하세요.

Migration 패턴:

Migration을 확장하는 익명 클래스를 생성하세요. up 메서드는 Schema의 posts 테이블에서 create를 호출합니다. id, user_id에 constrained와 cascadeOnDelete가 있는 foreignId, title을 위한 string, content를 위한 text, status 기본값이 draft인 string, timestamps, softDeletes를 정의하세요.

Service Layer 패턴:

App\Services 네임스페이스에 UserService 클래스를 생성하세요. UserDTO를 받는 create 메서드를 정의하세요. DB transaction으로 User 생성 (DTO 속성에서), 기본 bio로 profile 생성, 로드된 profile 관계와 함께 user를 반환하는 래핑을 래핑하세요. ActiveRecord\RecordInvalid 예외를 catch하여 검증 실패를 처리하세요.

### Symfony 7 패턴

Doctrine Attributes와 함께 Entity:

App\Entity 네임스페이스에 User 클래스를 생성하세요. ORM\Entity(repositoryClass=UserRepository::class)와 ORM\Table(name="users") 속성을 적용하세요. private nullable int id에 ORM\Id, ORM\GeneratedValue, ORM\Column 속성을 추가하세요. private nullable string name에 ORM\Column(length=255)과 Assert\NotBlank를 추가하세요. private nullable string email에 ORM\Column(length=180, unique=true)와 Assert\Email를 추가하세요.

Dependency Injection과 함께 Service:

App\Service 네임스페이스에 UserService 클래스를 생성하세요. readonly EntityManagerInterface와 readonly UserPasswordHasherInterface를 받는 생성자를 사용하세요. email과 password 문자열을 받는 createUser 메서드를 정의하세요. 새 User를 생성하고, email을 설정하고, password hasher로 비밀번호를 해싱하고, entity manager로 persist하고, flush하고 user를 반환하세요.

### 테스트 패턴

Laravel용 PHPUnit Feature Test:

Tests\Feature 네임스페이스에 RefreshDatabase trait이 있는 TestCase를 확장하는 UserApiTest를 생성하세요. name, email, password, password_confirmation이 포함된 JSON으로 api/users에 post하는 test_can_create_user 메서드를 정의하세요. 상태 201과 data에 id, name, email이 포함된 JSON 구조를 assert합니다. email이 있는 users 테이블을 assert합니다.

Laravel용 Pest Test:

App\Models\User와 Post를 사용하세요. it 함수를 사용하여 "post를 생성할 수 있음" 테스트를 생성합니다. factory로 user를 생성하고 actingAs(user)로 api/posts에 title과 content가 포함된 JSON을 post합니다. 상태 201을 assert하고 Post count가 1인 것을 expect합니다. "인증이 필요함" 테스트를 생성하고 인증 없이 post하고 상태 401을 assert합니다.

---

## Advanced Implementation (10분 이상)

포괄적인涵盖内容包括:

- Docker 및 Kubernetes용 프로덕션 배포 패턴
- observers, accessors, mutators를 포함한 고급 Eloquent 패턴
- embeddables와 상속을 포함한 Doctrine 고급 매핑
- 큐 및 작업 처리
- 이벤트 기반 아키텍처
- Redis 및 Memcached와 함께 캐싱 전략
- OWASP 패턴을 따르는 보안 모범 사례
- CI/CD 통합 패턴

다음을 참조하세요:

- modules/advanced-patterns.md 완전한 고급 패턴 가이드

---

## Context7 라이브러리 매핑

- laravel/framework Laravel 웹 프레임워크
- symfony/symfony Symfony 구성 요소 및 프레임워크
- doctrine/orm PHP용 Doctrine ORM
- phpunit/phpunit PHP 테스팅 프레임워크
- pestphp/pest 우아한 PHP 테스팅 프레임워크
- laravel/sanctum Laravel API 인증
- laravel/horizon Laravel 큐 대시보드

---

## Works Well With

- do-domain-backend REST API 및 마이크로서비스 아키텍처
- do-domain-database SQL 패턴 및 ORM 최적화
- do-workflow-testing DDD 및 테스트 전략
- do-platform-deploy Docker 및 배포 패턴
- do-essentials-debug AI 기반 디버깅
- do-foundation-quality TRUST 5 품질 원칙

---

## Troubleshooting

일반적인 문제:

PHP 버전 확인:

8.3 이상을 위해 php --version을 실행하세요. pdo, mbstring, openssl 확장을 확인하기 위해 php -m | grep을 사용하세요.

Composer Autoload 문제:

최적화된 autoloader를 위해 composer dump-autoload -o를 실행하세요. 패키지 캐시를 지우려면 composer clear-cache를 실행하세요.

Laravel 캐시 문제:

구성 캐시를 지우려면 php artisan config:clear를 실행하세요. 애플리케이션 캐시를 지우려면 php artisan cache:clear를 실행하세요. 경로 캐시를 지우려면 php artisan route:clear를 실행하세요. 컴파일된 뷰를 지우려면 php artisan view:clear를 실행하세요.

Symfony 캐시 문제:

캐시를 지우려면 php bin/console cache:clear를 실행하세요. 캐시를 예열하려면 php bin/console cache:warmup을 실행하세요.

데이터베이스 연결:

DB::connection()->getPdo() 호출 주위에 try-catch 블록을 사용하세요. 연결 시 성공 메시지를 출력하고 실패 시 예외 메시지를 출력하세요.

Migration 롤백:

마지막 마이그레이션을 롤백하려면 php artisan migrate:rollback --step=1을 사용하세요. 개발 재설정을 위해 php artisan migrate:fresh --seed를 사용하세요. Symfony의 경우 php bin/console doctrine:migrations:migrate prev를 사용하여 롤백하세요.

---

Version: 1.1.0 | Updated: 2026-01-11 | Status: Active
