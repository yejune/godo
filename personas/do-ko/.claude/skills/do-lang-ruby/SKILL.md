---
name: do-lang-ruby
description: >
  Ruby 3.3+ 개발 전문가 - Rails 7.2, ActiveRecord, Hotwire/Turbo, 모던 Ruby 패턴涵盖.
  Ruby API, 웹 애플리케이션, Rails 프로젝트 개발 시 사용하세요.
license: Apache-2.0
compatibility: Designed for Claude Code
allowed-tools: Read Grep Glob Bash(ruby:*) Bash(gem:*) Bash(bundle:*) Bash(rake:*) Bash(rspec:*) Bash(rubocop:*) Bash(rails:*) mcp__context7__resolve-library-id mcp__context7__get-library-docs
user-invocable: false
metadata:
  version: "1.1.0"
  category: "language"
  status: "active"
  updated: "2026-01-11"
  modularized: "true"
  tags: "language, ruby, rails, activerecord, hotwire, turbo, rspec"

# MoAI Extension: Progressive Disclosure
progressive_disclosure:
  enabled: true
  level1_tokens: 100
  level2_tokens: 5000

# MoAI Extension: Triggers
triggers:
  keywords: ["Ruby", "Rails", "ActiveRecord", "Hotwire", "Turbo", "RSpec", ".rb", "Gemfile", "Rakefile", "config.ru"]
  languages: ["ruby"]
---

## Quick Reference (30초 요약)

Ruby 3.3+ 개발 전문가 - Rails 7.2, ActiveRecord, Hotwire/Turbo, RSpec, 모던 Ruby 패턴.

자동 트리거: .rb 확장자 파일, Gemfile, Rakefile, config.ru, Rails 또는 Ruby 논의

핵심 기능:

- Ruby 3.3 기능: YJIT 프로덕션 준비, pattern matching, Data class, endless methods
- 웹 프레임워크: Turbo, Stimulus, ActiveRecord가 있는 Rails 7.2
- 프론트엔드: SPA 같은 경험을 위한 Hotwire (Turbo 및 Stimulus)
- 테스트: factories, request specs, system specs가 있는 RSpec
- 백그라운드 작업: ActiveJob과 함께 Sidekiq
- 패키지 관리: Gemfile과 함께 Bundler
- 코드 품질: Rails cops가 있는 RuboCop
- 데이터베이스: 마이그레이션, 연관, scopes가 있는 ActiveRecord

### Quick Patterns

Rails Controller 패턴:

ApplicationController를 상속받는 UsersController를 생성하세요. show, edit, update, destroy 액션에서만 호출되는 set_user를 위하여 before_action을 추가하세요. User.all을 인스턴스 변수에 할당하는 index 메서드를 정의하세요. user_params로 새 User를 생성하는 create 메서드를 정의하세요. 성공 시 redirect 또는 unprocessable_entity 상태로 new를 렌더링하는 format.html, Turbo 응답을 위한 format.turbo_stream을 사용하여 respond_to 블록을 사용하세요. params[:id]로 사용자를 찾는 private set_user 메서드를 추가하세요. user를 요구하고 name과 email을 허용하는 user_params 메서드를 추가하세요.

ActiveRecord Model 패턴:

ApplicationRecord를 상속받는 User model을 생성하세요. dependent: destroy로 has_many :posts, dependent: destroy로 has_one :profile를 정의하세요. presence, uniqueness, URI::MailTo::EMAIL_REGEXP 형식으로 email을 검증합니다. presence와 length(minimum: 2, maximum: 100)으로 name을 검증합니다. active가 true인 것을 필터링하는 active scope를 정의하세요. created_at DESC로 정렬하는 recent scope를 정의하세요. first_name과 last_name을 공백으로 결합하고 제거하는 full_name 메서드를 정의하세요.

RSpec Test 패턴:

User model 타입에 RSpec.describe를 생성하세요. describe validations 블록에서 email에 validate_presence_of와 validate_uniqueness_of에 대한 expectation을 추가하세요. describe full_name 블록에서 first_name이 "John", last_name이 "Doe"인 user를 빌드하는 let을 사용하세요. user.full_name이 "John Doe"와 같은지 expect하는 it 블록을 추가하세요.

---

## Implementation Guide (5분 가이드)

### Ruby 3.3 새 기능

YJIT Production-Ready:

Ruby 3.3에서는 YJIT가 기본적으로 활성화되어 Rails 애플리케이션에 15~20% 성능 향상을 제공합니다. ruby --yjit 플래그로 ruby를 실행하거나 RUBY_YJIT_ENABLE 환경변수를 1로 설정하여 활성화하세요. RubyVM::YJIT.enabled? 메서드를 호출하여 상태를 확인하세요.

case/in과 함께 Pattern Matching:

response 파라미터를 받는 process_response 메서드를 생성하세요. response에 case를 사용하고 패턴 매칭을 위해 in을 사용하세요. status가 ok이고 data를 추출하는 케이스에서 성공 메시지와 함께 data를 출력하세요. status가 error이고 msg를 추출하는 케이스에서 에러 메시지를 출력하세요. guard 조건으로 pending 또는 processing과 일치하는 status 케이스를 사용하세요. 알 수 없는 응답을 위한 else를 사용하세요.

Immutable Structs를 위한 Data Class:

Data.define을 사용하여 name과 email 심볼로 User를 생성하세요. greeting 메서드를 반환하는 블록을 추가하여 name이 포함된 hello 메시지를 반환하세요. 키워드 인자로 user 인스턴스를 생성하세요. name 속성에 접근하고 greeting 메서드를 호출하세요.

Endless Method Definition:

add, multiply, positive? 메서드에 대해 equals 신택자 구문을 사용하는 Calculator 클래스를 생성하세요.

### Rails 7.2 패턴

Gemfile의 Application Setup:

rubygems.org로 source를 설정하세요. 7.2 버전 제약 조건이 있는 rails, 1.5인 pg, 6.0 이상인 puma, turbo-rails, stimulus-rails, 7.0인 sidekiq를 추가하세요. development 및 test 그룹에 7.0인 rspec-rails, factory_bot_rails, faker, require: false가 있는 rubocop-rails를 추가하세요. test 그룹에 capybara와 shoulda-matchers를 추가하세요.

Concerns와 함께 Model:

ActiveSupport::Concern을 확장하는 Sluggable 모듈을 생성하세요. included 블록에서 create 시 generate_slug를 위한 before_validation과 slug에 presence와 uniqueness를 검증합니다. to_param을 반환하여 slug를 정의합니다. title이 존재하고 slug가 비어있으면 parameterized title에서 slug를 설정하는 private generate_slug 메서드를 추가하세요. belongs_to :user, dependent: destroy로 has_many :comments, has_many_attached :images를 포함하는 Post model에 Sluggable을 include하세요. 검증과 published scope를 추가하세요.

Service Objects 패턴:

user_params를 받는 initialize로 UserRegistrationService를 생성하세요. User를 생성하고 ActiveRecord::Base.transaction을 사용하여 사용자를 저장하고 프로필을 생성하고 환영 이메일을 보내는 call 메서드를 정의하세요. success가 true이고 user가 있는 Result를 반환합니다. RecordInvalid를 구조하여 success가 false이고 errors가 있는 Result를 반환합니다. create_profile과 send_welcome_email을 위한 private 메서드를 추가하세요. success, user, errors를 가진 Data.define으로 Result를 정의하고 success?와 failure? 술어 메서드를 추가하세요.

### Hotwire Turbo와 Stimulus

Turbo Frames 패턴:

index 뷰에서 posts id로 turbo_frame_tag를 사용하고 각 post를 렌더링합니다. post 부분 분할에서 post의 dom_id로 turbo_frame_tag를 사용하여 h2 링크와 잘린 content 단락을 포함하는 article를 포함하세요.

Turbo Streams 패턴:

controller create 액션에서 current_user.posts에서 post를 빌드합니다. save 성공 여부에 따라 redirect 또는 render를 위한 format.html, format.turbo_stream과 함께 respond_to를 사용하세요. create.turbo_stream.erb 뷰에서 post로 posts에 turbo_stream.prepend를, new_post 폼 부분으로 turbo_stream.update를 사용하세요.

Stimulus Controller 패턴:

JavaScript 컨트롤러 파일에서 hotwired/stimulus에서 Controller를 임포트하고 export하세요. static targets 배열이 있는 input과 submit을 가진 클래스를 default로 확장하세요. connect 메서드에서 validate를 호출합니다. validate 메서드에서 모든 input targets에 값이 있는지 확인하고 submit target을 그에 따라 비활성화합니다.

### RSpec 테스트 기본

Factory Bot 패턴:

factories 파일에서 email에 sequence, name에 Faker::Name.name, password에 password123을 가진 user factory를 정의하세요. role을 admin 심볼로 설정하는 admin trait를 추가하세요. posts_count가 3인 transient를 사용하고 after create 콜백으로 사용자를 위한 posts를 create_list하는 with_posts trait를 추가하세요.

---

## Advanced Implementation (10분 이상)

포괄적인涵盖内容包括:

- Docker 및 Kubernetes용 프로덕션 배포 패턴
- 다형성, STI, 쿼리 객체를 포함한 고급 ActiveRecord 패턴
- Action Cable 실시간 기능
- 성능 최적화 기법
- 보안 모범 사례
- CI/CD 통합 패턴
- 완전한 RSpec 테스트 패턴

다음을 참조하세요:

- modules/advanced-patterns.md 프로덕션 패턴 및 고급 기능
- modules/testing-patterns.md 완전한 RSpec 테스트 가이드

---

## Context7 라이브러리 매핑

- rails/rails Ruby on Rails 웹 프레임워크
- rspec/rspec RSpec 테스팅 프레임워크
- hotwired/turbo-rails Rails용 Turbo
- hotwired/stimulus-rails Rails용 Stimulus
- sidekiq/sidekiq 백그라운드 작업 처리
- rubocop/rubocop Ruby 스타일 가이드 강제
- thoughtbot/factory_bot 테스트 데이터 팩토리

---

## Works Well With

- do-domain-backend REST API 및 웹 애플리케이션 아키텍처
- do-domain-database SQL 패턴 및 ActiveRecord 최적화
- do-workflow-testing DDD 및 테스트 전략
- do-essentials-debug AI 기반 디버깅
- do-foundation-quality TRUST 5 품질 원칙

---

## Troubleshooting

일반적인 문제:

Ruby 버전 확인:

3.3 이상을 위해 ruby --version을 실행하세요. YJIT 상태를 확인하기 위해 ruby -e 'puts RubyVM::YJIT.enabled?'를 실행하세요.

Rails 버전 확인:

7.2 이상을 위해 rails --version을 실행하세요. 전체 환경 정보를 위해 bundle exec rails about를 실행하세요.

데이터베이스 연결 문제:

- config/database.yml 구성을 확인하세요
- PostgreSQL 또는 MySQL 서비스가 실행 중인지 확인하세요
- 데이터베이스가 존재하지 않으면 rails db:create를 실행하세요

Asset Pipeline 문제:

assets를 컴파일하려면 rails assets:precompile을 실행하세요. 컴파일된 assets를 지우려면 rails assets:clobber를 실행하세요.

RSpec 설정 문제:

초기 설정을 위해 rails generate rspec:install을 실행하세요. 단일 spec을 위해 bundle exec rspec spec 파일 경로를 실행하세요. 상세한 출력을 위해 bundle exec rspec --format documentation을 실행하세요.

Turbo 및 Stimulus 문제:

JavaScript 설정을 위해 rails javascript:install:esbuild를 실행하세요. Turbo 설치를 위해 rails turbo:install을 실행하세요.

---

Last Updated: 2026-01-11
Status: Active (v1.1.0)
