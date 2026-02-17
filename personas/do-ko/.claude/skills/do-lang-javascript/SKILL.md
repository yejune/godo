---
name: do-lang-javascript
description: >
  JavaScript ES2024+ 개발 전문가 - Node.js 22 LTS, Bun 1.x (serve, SQLite, S3,
  shell, test), Deno 2.x, 테스트 (Vitest, Jest), 린팅 (ESLint 9, Biome), 백엔드
  프레임워크 (Express, Fastify, Hono)涵盖. JavaScript API, 웹 애플리케이션,
  Node.js 프로젝트 개발 시 사용하세요.
license: Apache-2.0
compatibility: Designed for Claude Code
allowed-tools: Read Grep Glob Bash(node:*) Bash(npm:*) Bash(npx:*) Bash(yarn:*) Bash(pnpm:*) Bash(bun:*) Bash(deno:*) Bash(jest:*) Bash(vitest:*) Bash(eslint:*) Bash(prettier:*) Bash(biome:*) mcp__context7__resolve-library-id mcp__context7__get-library-docs
user-invocable: false
metadata:
  version: "1.2.0"
  category: "language"
  status: "active"
  updated: "2026-01-11"
  modularized: "false"
  tags: "language, javascript, nodejs, bun, deno, vitest, eslint, express"

# MoAI Extension: Progressive Disclosure
progressive_disclosure:
  enabled: true
  level1_tokens: 100
  level2_tokens: 5000

# MoAI Extension: Triggers
triggers:
  keywords: ["JavaScript", "Node.js", "Bun", "Deno", "Express", "Fastify", "Hono", "Vitest", "Jest", "ESLint", ".js", "package.json"]
  languages: ["javascript", "js"]
---

## Quick Reference (30초 요약)

JavaScript ES2024+ 개발 전문가 - Node.js 22 LTS, 여러 런타임, 현대 도구와 모던 JavaScript.

자동 트리거: .js, .mjs, .cjs 확장자 파일, package.json, Node.js 프로젝트, JavaScript 논의

핵심 스택:

- ES2024+: Set 메서드, Promise.withResolvers, 불변 배열, import 속성
- Node.js 22 LTS: 네이티브 TypeScript, 내장 WebSocket, 안정적인 감시 모드
- 런타임: Node.js 20 및 22 LTS, Deno 2.x, Bun 1.x
- 테스트: Vitest, Jest, Node.js 테스트 러너
- 린팅: ESLint 9 플랫 config, Biome
- 번들러: Vite, esbuild, Rollup
- 프레임워크: Express, Fastify, Hono, Koa

빠른 명령어:

npm create vite@latest 프로젝트이름 --template vanilla로 Vite 프로젝트를 생성하세요. npm init와 npm install -D vitest eslint eslint/js로 모던 도구로 초기화하세요. node --watch로 Node.js 감시 모드로 실행하세요. node --experimental-strip-types로 Node.js 22+에서 TypeScript를 직접 실행하세요.

---

## Implementation Guide (5분 가이드)

### ES2024 핵심 기능

Set 연산:

값 1, 2, 3, 4로 setA를, 값 3, 4, 5, 6으로 setB를 생성하세요. setA.intersection(setB)를 호출하여 3과 4를 포함하는 Set을 가져오세요. setA.union(setB)를 호출하여 1부터 6까지를 포함하는 Set을 가져오세요. setA.difference(setB)를 호출하여 1과 2를 포함하는 Set을 가져오세요. setA.symmetricDifference(setB)를 호출하여 1, 2, 5, 6을 포함하는 Set을 가져오세요. set 비교를 위해 isSubsetOf, isSupersetOf, isDisjointFrom 메서드를 호출하세요.

Promise.withResolvers:

Promise.withResolvers() 호출에서 promise, resolve, reject를 분해하는 createDeferred 함수를 생성하세요. 이 세 가지 속성을 가진 객체를 반환하세요. deferred 인스턴스를 생성하고 1000밀리초 후에 done으로 확인하도록 timeout을 설정하고 결과를 위해 프라미스를 await하세요.

불변 배열 메서드:

3, 1, 4, 1, 5 값으로 original 배열을 생성하세요. 원본을 수정하지 않고 정렬된 새 배열을 가져오려면 toSorted를 호출하세요. 새로운 역순 배열을 가져오려면 toReversed를 호출하세요. 인덱스 1, 삭제 수 2, 삽입 값 9로 toSpliced를 호출하세요. 인덱스 2에서 값 99와 함께 with를 호출하여 대체된 요소를 가진 새 배열을 가져오세요. 원본 배열은 변경되지 않습니다.

Object.groupBy 및 Map.groupBy:

type과 name 속성을 가진 객체를 포함하는 items 배열을 생성하세요. items와 item.type을 반환하는 함수로 Object.groupBy를 호출하여 type으로 그룹화된 배열을 가진 객체를 가져오세요. 동일한 인자로 Map.groupBy를 호출하여 type 키와 배열 값을 가진 Map을 가져오세요.

### ES2025 기능

JSON 모듈을 위한 Import 속성:

type 속성을 json으로 설정하여 config.json에서 config를 임포트하세요. type 속성을 css로 설정하여 styles.css에서 styles를 임포트하세요. config.apiUrl 속성에 접근하세요.

RegExp.escape:

괄호와 같은 특수 문자를 포함하는 userInput 문자열을 생성하세요. userInput과 함께 RegExp.escape를 호출하여 이스케이프된 패턴 문자열을 가져오세요. 안전한 패턴으로 새 RegExp를 생성하세요.

### Node.js 22 LTS 기능

내장 WebSocket 클라이언트:

wss URL로 새 WebSocket을 생성하세요. JSON 문자열화된 메시지를 보내는 open 이벤트용 이벤트 리스너를 추가하세요. event.data를 JSON으로 파싱하고 받은 데이터를 로깅하는 message 이벤트용 이벤트 리스너를 추가하세요.

네이티브 TypeScript 지원 실험적:

node --experimental-strip-types 플래그와 함께 Node.js 22.6+에서 .ts 파일을 직접 실행하세요. Node.js 22.18+에서는 타입 제거가 기본적으로 활성화되어 파일을 직접 실행할 수 있습니다.

안정적인 감시 모드:

파일 변경 시 자동 재시작을 위해 node --watch를 사용하세요. src 및 config와 같은 특정 디렉토리를 감시하려면 --watch-path 플래그를 여러 번 사용하세요.

권한 모델:

파일 시스템 액세스를 제한하려면 node --permission --allow-fs-read를 특정 경로로 설정하여 사용하세요. 도메인 이름과 함께 --allow-net 플래그를 사용하여 네트워크 액세스를 제한하세요.

### 백엔드 프레임워크

Express 전통적 패턴:

express를 임포트하세요. express() 함수를 호출하여 app을 생성하세요. express.json 미들웨어를 사용하세요. api/users에 DB 쿼리를 await하고 json으로 응답하는 get 엔드포인트를 생성하세요. 사용자를 생성하고 상태 201과 json으로 응답하는 post 엔드포인트를 생성하세요. 포트 3000에서 서버 실행을 로깅하는 콜백과 함께 listen을 호출하세요.

Fastify 고성능 패턴:

Fastify를 임포트하세요. logger를 true로 설정하여 fastify 인스턴스를 생성하세요. 바디에 type object, required에 name과 email 배열, validation 제약 조건이 있는 속성을 가진 userSchema를 정의하세요. schema 옵션과 사용자를 생성하고 code 201을 반환하는 async 핸들러로 post 엔드포인트를 생성하세요. 포트 3000으로 listen을 호출하세요.

Hono 엣지 퍼스트 패턴:

Hono와 미들웨어 함수를 임포트하세요. app 인스턴스를 생성하세요. 모든 경로에 logger 미들웨어를 사용하세요. api 경로에 cors 미들웨어를 사용하세요. DB 쿼리를 await하고 c.json을 반환하는 api/users에 get 엔드포인트를 생성하세요. 필수 필드를 확인하는 validator 미들웨어가 있는 post 엔드포인트를 생성한 다음 사용자를 생성하고 상태 201로 c.json을 반환하세요. 기본값으로 app을 내보내세요.

### Vitest와 함께 테스트

구성:

defineConfig로 vitest.config.js를 생성하세요. globals true, environment node, v8 제공자와 text, json, html 리포터가 있는 coverage를 설정한 test 객체를 설정하세요.

테스트 예제:

테스트 파일에서 describe, it, expect, vi, beforeEach를 vitest에서 임포트하세요. 테스트할 함수를 임포트하세요. User Service용 describe 블록을 생성하세요. beforeEach에서 vi.clearAllMocks를 호출하세요. createUser를 await하고 결과가 name과 email이 있는 객체와 일치하고 id가 정의되어 있는지 expect하는 "사용자를 생성해야 함" it 블록을 생성하세요. createUser가 Invalid email 에러로 reject할 것을 expect하는 "잘못된 이메일에서 throw해야 함" it 블록을 생성하세요.

### ESLint 9 플랫 Config

eslint.config.js를 생성하세요. eslint/js에서 js와 globals를 임포트하세요. js.configs.recommended 뒤에 ecmaVersion 2025, sourceType module, globals.node와 globals.es2025를 병합한 languageOptions를 포함하는 객체 배열을 내보내세요. args ignore 패턴이 있는 no-unused-vars, 허용된 메서드가 있는 no-console, error인 prefer-const, error인 no-var에 대한 규칙을 설정하세요.

### Biome 올인원

schema URL로 biome.json을 생성하세요. organizeImports를 활성화하세요. 권장 규칙으로 linter를 활성화하세요. indentStyle space, indentWidth 2로 formatter를 활성화하세요. javascript.formatter에서 quoteStyle을 single로, semicolons를 always로 설정하세요.

---

## Advanced Patterns

고급 async 패턴, 모듈 시스템 세부 정보, 성능 최적화, 프로덕션 배포 구성에 대한 포괄적인 문서는 다음을 참조하세요:

- reference.md 완전한 API 참조, Context7 라이브러리 매핑, 패키지 관리자 비교
- examples.md 프로덕션 준비 코드 예제, 풀스택 패턴, 테스트 템플릿

### Context7 통합

Node.js 문서의 경우 esm modules async 주제로 nodejs/node와 context7 get library docs를 사용하세요. Express의 경우 middleware routing 주제로 expressjs/express를 사용하세요. Fastify의 경우 plugins hooks 주제로 fastify/fastify를 사용하세요. Hono의 경우 middleware validators 주제로 honojs/hono를 사용하세요. Vitest의 경우 mocking coverage 주제로 vitest-dev/vitest를 사용하세요.

---

## Works Well With

- do-lang-typescript JSDoc와 함께 TypeScript 통합 및 타입 검사
- do-domain-backend API 설계 및 마이크로서비스 아키텍처
- do-domain-database 데이터베이스 통합 및 ORM 패턴
- do-workflow-testing DDD 워크플로우 및 테스트 전략
- do-foundation-quality 코드 품질 표준
- do-essentials-debug JavaScript 애플리케이션 디버깅

---

## Quick Troubleshooting

모듈 시스템 문제:

type 필드를 위해 package.json을 확인하세요. ESM은 type module과 import/export를 사용합니다. CommonJS는 type commonjs 또는 필드 생략과 require/module.exports를 사용합니다.

Node.js 버전 확인:

20.x 또는 22.x LTS를 위해 node --version을 실행하세요. 10.x 이상을 위해 npm --version을 실행하세요.

일반적인 수정사항:

npm cache clean --force로 npm 캐시를 지우세요. node_modules와 package-lock.json을 삭제한 다음 npm install을 실행하세요. npm config prefix를 홈 디렉토리 npm-global 폴더로 설정하여 권한 문제를 수정하세요.

ESM 및 CommonJS 상호 운용:

ESM에서 CommonJS를 임포트하려면 기본값을 임포트한 다음 거기서 named exports를 분해하세요. CommonJS에서 동적 임포트의 경우 await import를 사용하고 default 속성을 분해하세요.

---

Last Updated: 2026-01-11
Status: Active (v1.2.0)
