---
name: do-lang-typescript
description: >
  TypeScript 5.9+ 개발 전문가 - React 19, Next.js 16 App Router, tRPC로 타입 안전한
  API, Zod 검증, 모던 TypeScript 패턴涵盖. TypeScript 애플리케이션,
  React 컴포넌트, Next.js 페이지, 타입 안전한 API 개발 시 사용하세요.
license: Apache-2.0
compatibility: Designed for Claude Code
allowed-tools: Read Grep Glob mcp__context7__resolve-library-id mcp__context7__get-library-docs
user-invocable: false
metadata:
  version: "1.1.0"
  category: "language"
  status: "active"
  updated: "2026-01-11"
  modularized: "false"
  tags: "typescript, react, nextjs, frontend, fullstack"

# MoAI Extension: Progressive Disclosure
progressive_disclosure:
  enabled: true
  level1_tokens: 100
  level2_tokens: 5000

# MoAI Extension: Triggers
triggers:
  keywords: ["TypeScript", "React", "Next.js", "tRPC", "Zod", ".ts", ".tsx", "tsconfig.json"]
  languages: ["typescript", "tsx"]
---

## Quick Reference (30초 요약)

TypeScript 5.9+ 개발 전문가 - React 19, Next.js 16, 타입 안전한 API 패턴과 모던 TypeScript.

자동 트리거: .ts, .tsx, .mts, .cts 확장자 파일, TypeScript 구성, React 또는 Next.js 프로젝트

핵심 스택:

- TypeScript 5.9: 지연된 모듈 평가, 데코레이터, satisfies 연산자
- React 19: Server Components, use 훅, Actions, 동시성 기능
- Next.js 16: App Router, Server Actions, 미들웨어, ISR/SSG/SSR
- 타입 안전한 API: tRPC 11, Zod 3.23, tanstack-query
- 테스트: Vitest, React Testing Library, Playwright

빠른 명령어:

npx create-next-app으로 Next.js 16 프로젝트를 생성하세요. latest, typescript, tailwind, app 플래그를 사용하세요. npm install로 타입 안전한 API 스택을 설치하세요: trpc server, client, react-query, zod, tanstack react-query. npm install -D로 테스트 스택을 설치하세요: vitest, testing-library/react, playwright/test.

---

## Implementation Guide (5분 가이드)

### TypeScript 5.9 핵심 기능

타입 확장 없이 타입 검증하는 satisfies 연산자:

red, green, blue 문자열 리터럴의 합집합인 Colors 타입을 정의하세요. palette 객체를 red를 숫자 배열, green을 16진수 문자열, blue를 숫자 배열로 생성하세요. Record<Colors, string | number[]>와 함께 satisfies 연산자를 적용하세요. 이제 palette.red는 숫자 배열로 유추되어 map 메서드를 사용할 수 있고, palette.green은 문자열로 유추되어 toUpperCase를 사용할 수 있습니다.

지연된 모듈 평가:

heavy/module/path에서 import defer * as namespace를 사용하세요. 모듈이 필요한 함수에서 namespace 속성에 접근하면 첫 번째 사용 시 모듈이 로드됩니다.

모던 데코레이터 Stage 3:

target 함수와 ClassMethodDecoratorContext를 받는 logged 함수 데코레이터를 생성하세요. 메서드 이름을 로깅한 다음 apply로 target을 호출하는 함수를 반환하세요. 데이터를 가져오는 클래스 메서드에 logged 데코레이터를 적용하세요.

### React 19 패턴

App Router의 기본 Server Components:

app/users/[id]/page.tsx의 페이지 컴포넌트 경우, params를 Promise<{id: string}>으로 가진 PageProps 인터페이스를 정의하세요. params를 await하는 async default 함수를 생성하고, 사용자를 위해 DB를 쿼리하고, 찾지 못하면 notFound를 호출하고, 사용자 이름을 가진 main 요소를 반환하세요.

Promise와 Context를 언래핑하는 use 훅:

use client 지시어로 표시된 클라이언트 컴포넌트에서 react에서 use를 임포트하세요. userPromise prop을 Promise<User> 타입으로 받는 UserProfile 컴포넌트를 생성하세요. 프라미스에 use 훅을 호출하여 해결될 때까지 일시 중단하세요. 사용자 이름이 포함된 div를 반환하세요.

서버 함수와 함께 폼 처리하는 Actions:

use server 지시어로 표시된 server actions 파일에서 revalidatePath를 임포트하세요. name과 email 검증을 위해 zod로 CreateUserSchema를 정의하세요. FormData를 받는 async createUser 함수를 생성하세요. 스키마로 파싱하고, DB에 사용자를 생성하고, 경로를 재검증하세요.

폼 상태를 위한 useActionState:

클라이언트 컴포넌트에서 useActionState를 임포트하세요. createUser action과 함께 useActionState를 호출하여 state, action, isPending을 분해하는 form 컴포넌트를 생성하세요. action prop, pending일 때 비활성화된 input, 동적 텍스트가 있는 버튼, state에서의 에러 메시지가 있는 form을 반환하세요.

### Next.js 16 App Router

경로 구조:

app 디렉토리는 루트 레이아웃을 위한 layout.tsx, 홈 경로를 위한 page.tsx, 로딩 UI를 위한 loading.tsx, 에러 경계를 위한 error.tsx, API 경로를 위한 api/route.ts를 포함합니다. users와 같은 하위 디렉토리는 목록을 위한 page.tsx와 동적 경로를 위한 [id]/page.tsx를 포함합니다. 경로 그룹은 (marketing)/about/page.tsx와 같은 괄호를 사용합니다.

Metadata API:

Metadata 타입을 임포트하세요. default와 template을 포함하는 객체를 가진 title, description 문자열로 metadata 객체를 내보내세요. params를 받고 params를 await하고 사용자를 가져오며 title을 사용자 이름으로 설정하는 객체를 반환하는 async generateMetadata 함수를 내보내세요.

검증과 함께 Server Actions:

server 파일에서 zod, revalidatePath, redirect를 임포트하세요. id, name, email 검증으로 UpdateUserSchema를 정의하세요. prevState와 formData를 받는 async updateUser 함수를 생성하세요. safeParse로 파싱하고 실패하면 에러를 반환하고, DB를 업데이트하고, 경로를 재검증하고, 리디렉션하세요.

### tRPC와 타입 안전한 API

서버 설정:

trpc/server에서 initTRPC와 TRPCError를 임포트하세요. initTRPC.context<Context>().create()를 호출하여 t를 생성하세요. router, publicProcedure, protectedProcedure를 t에서 내보내세요. protectedProcedure는 세션 사용자를 확인하고 누락된 경우 UNAUTHORIZED 에러를 발생시키는 미들웨어를 사용합니다.

라우터 정의:

zod를 임포트하세요. router 함수로 userRouter를 생성하세요. publicProcedure와 input 스키마로 id를 uuid 문자열로 사용하여 getById 프로시저를 정의하고 id로 사용자를 찾는 쿼리를 정의하세요. protectedProcedure와 input 스키마로 name과 email을 사용하여 create 프로시저를 정의하고 사용자를 생성하는 뮤테이션을 정의하세요.

클라이언트 사용:

클라이언트 컴포넌트에서 page 파라미터로 trpc.user.list.useQuery를 호출하는 UserList 함수를 생성하세요. data와 isLoading을 분해하세요. trpc.user.create.useMutation으로 뮤테이션을 생성하세요. 로딩 상태 또는 사용자 목록을 반환하세요.

### Zod 스키마 패턴

복잡한 검증:

uuid 문자열 id, 최소/최대 길이 name, email 형식 email, admin/user/guest 열거 role, coerce.date createdAt을 가진 z.object로 UserSchema를 생성하세요. strict 메서드를 적용하세요. 스키마에서 User 타입을 유추하세요. CreateUserSchema를 생성할 때 id와 createdAt을 생략하고 password와 confirmPassword로 확장하며 커스텀 메시지와 경로로 password 일치 검증을 위한 refine을 추가하세요.

### 상태 관리

클라이언트 상태를 위한 Zustand:

zustand와 middleware에서 create를 임포트하세요. user를 User | null로, login 메서드, logout 메서드를 가진 AuthState 인터페이스를 정의하세요. devtools와 persist 미들웨어로 래핑된 create 함수로 useAuthStore를 생성하세요. 초기 user를 null로 설정하고 login은 user를 설정하고 logout은 user를 null로 설정하세요. persist는 auth-storage 이름을 사용합니다.

원자 상태를 위한 Jotai:

jotai에서 atom을, jotai/utils에서 atomWithStorage를 임포트하세요. 초기값 0으로 countAtom을 생성하세요. countAtom을 가져와 2를 곱하는 파생 atom으로 doubleCountAtom을 생성하세요. light/dark 테마를 저장에 지속하는 atomWithStorage로 themeAtom을 생성하세요.

---

## Advanced Patterns

고급 TypeScript 패턴, 성능 최적화, 테스트 전략, 배포 구성에 대한 포괄적인 문서는 다음을 참조하세요:

- reference.md 완전한 API 참조, Context7 라이브러리 매핑, 고급 타입 패턴
- examples.md 프로덕션 준비 코드 예제, 풀스택 패턴, 테스트 템플릿

### Context7 통합

TypeScript 문서의 경우 decorators satisfies 주제로 microsoft/TypeScript를 사용하세요. React 19의 경우 server-components use-hook 주제로 facebook/react를 사용하세요. Next.js 16의 경우 app-router server-actions 주제로 vercel/next.js를 사용하세요. tRPC의 경우 procedures middleware 주제로 trpc/trpc를 사용하세요. Zod의 경우 schema-validation 주제로 colinhacks/zod를 사용하세요.

---

## Works Well With

- do-domain-frontend UI 컴포넌트 및 스타일링 패턴
- do-domain-backend API 설계 및 데이터베이스 통합
- do-library-shadcn 컴포넌트 라이브러리 통합
- do-workflow-testing 테스트 전략 및 패턴
- do-foundation-quality 코드 품질 표준
- do-essentials-debug TypeScript 애플리케이션 디버깅

---

## Quick Troubleshooting

TypeScript 에러:

타입 검사만을 위해 npx tsc --noEmit을 실행하세요. 성능 추적을 위해 npx tsc --generateTrace --output-directory를 실행하세요.

React 및 Next.js 문제:

빌드 에러를 확인하기 위해 npm run build를 실행하세요. ESLint 검사를 위해 npx next lint를 실행하세요. 캐시를 지우기 위해 .next 디렉토리를 삭제하고 npm run dev를 실행하세요.

타입 안전 패턴:

예기치 않은 값에 대해 에러를 발생시키는 never 타입 파라미터를 가진 assertNever 함수를 생성하고, 포괄적인 switch 문에서 사용하세요. 값이 id 속성을 가진 객체인지 확인하고 타입 술어를 반환하는 isUser 타입 가드 함수를 생성하세요.

---

Last Updated: 2026-01-11
Status: Active (v1.1.0)
