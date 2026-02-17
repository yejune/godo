---
name: do-library-shadcn
description: >
  React 애플리케이션과 Tailwind CSS를 위한 shadcn/ui 컴포넌트 라이브러리 전문가입니다.
  UI 컴포넌트, 디자인 시스템, 또는 shadcn/ui, Radix primitives,
  Tailwind 기반 컴포넌트 라이브러리를 구현할 때 사용하세요.
  React 프레임워크가 아닌 custom CSS-only 스타일링의 경우
  (대신 do-domain-frontend를 사용하세요).
license: Apache-2.0
compatibility: Designed for Claude Code
allowed-tools: Read Grep Glob mcp__context7__resolve-library-id mcp__context7__get-library-docs
user-invocable: false
metadata:
  version: "2.1.0"
  category: "library"
  modularized: "true"
  status: "active"
  updated: "2026-01-11"
  tags: "library, shadcn, enterprise, development, ui"
  aliases: "do-library-shadcn"

# MoAI Extension: Triggers
triggers:
  keywords: ["shadcn", "component library", "design system", "radix", "tailwind", "ui components"]
---

## Quick Reference

Enterprise shadcn/ui 컴포넌트 라이브러리 전문가

AI 기반 컴포넌트 아키텍처, Context7 통합, 지능형 컴포넌트 오케스트레이션이 있는 모던 React 애플리케이션을 위한 포괄한 shadcn/ui 전문가입니다.

핵심 기능:

- Context7 MCP를 사용한 AI 기반 컴포넌트 아키텍처
- 자동화된 테마 커스터마이제이션
- 접근성 및 성능이 포함된 고급 컴포넌트 오케스트레이션
- 제로 구성 가능한 엔터프라이즈 UI 프레임워크
- 사용량 통찰을 위한 예측 컴포넌트 분석

사용 시기:

- shadcn/ui 컴포넌트 라이브러리 논의
- React 컴포넌트 아키텍처 계획
- Tailwind CSS 통합 및 design token
- 접근성 구현
- 디자인 시스템 커스터마이제이션

Module Organization:

- Core Concepts: 이 파일은 shadcn/ui 개요, 아키텍처, 에코시스템을 다룹니다.
- Components: shadcn-components.md 모듈은 컴포넌트 라이브러리와 고급 패턴을 다룹니다.
- Theming: shadcn-theming.md 모듈은 테마 시스템과 커스터마이제이션 전략을 다룹니다.
- Advanced Patterns: advanced-patterns.md 모듈은 복잡한 구현을 다룹니다.
- Optimization: optimization.md 모듈은 성능 튜닝을 다룹니다.

---

## Implementation Guide

### shadcn/ui 개요

shadcn/ui는 Radix UI와 Tailwind CSS로 구축된 재사용 가능한 컴포넌트 모음입니다. 전통적인 컴포넌트 라이브러리와 달리 npm 패키지가 아니라 프로젝트로 복사하는 컴포넌트 모음입니다.

주요 이점:

- 컴포넌트에 대한 완전한 제어 및 소유권
- Radix UI primitives 이외의 zero 의존성
- Tailwind CSS를 통한 완전한 커스터마이제이션
- 훌률한 타입 안전성 (excellent type safety를 가진 TypeScript-first)
- WCAG 2.1 AA 준수를 위한 내장된 접근성

아키텍처 철학:

shadcn/ui 컴포넌트는 기본적으로 Radix UI Primitives(스타일링되지 않은 접근성이 가능한 primitives) 위에 구축됩니다. Tailwind CSS는 유틸리티 기반 스타링을 제공합니다. TypeScript는 전체적으로 타입 안전성을 보장합니다. 사용자의 커스터마이제이션 레이어 최종 구현을 제어합니다.

### Core Component Categories

폼 컴포넌트: Input, Select, Checkbox, Radio, Textarea를 포함합니다. react-hook-form과 Zod와 통합되어 폼 검증을 지원합니다. ARIA labels로 접근성이 보장됩니다.

디스플레이 컴포넌트: Card, Dialog, Sheet, Drawer, Popover를 포함합니다. responsive design pattern이 내장되어 있습니다. dark mode 지원이 포함되어 있습니다.

내비게이션 컴포넌트: Navigation Menu, Breadcrumb, Tabs, Pagination을 포함합니다. 키보드 navigation 지원이 내장되어 있습니다. focus management가 자동으로 처리됩니다.

데이터 컴포넌트: Table, Calendar, DatePicker, Charts를 포함합니다. large dataset을 위한 virtual scrolling이 지원됩니다. TanStack Table 통합이 지원됩니다.

피드백 컴포넌트: Alert, Toast, Progress, Badge, Avatar를 포함합니다. loading states와 skeletons이 사용 가능합니다. error boundaries가 지원됩니다.

### Installation and Setup

1단계: npx와 함께 latest 버전으로 shadcn-ui init 명령을 실행하여 shadcn/ui를 초기화하세요.

2단계: ui.shadcn.com/schema.json에 schema URL이 포함된 components.json을 구성하세요. 스타일을 default로 설정하고 RSC과 TSX를 활성화합니다. Tailwind 설정을 포함하여 config 경로, CSS 경로, 기본 색상, CSS variables 활성화, 선택적 prefix를 설정합니다. 컴포넌트, utils, ui 경로에 대한 alias를 설정합니다.

3단계: npx와 함께 shadcn-ui add 명령과 함께 컴포넌트를 개별적으로 추가하세요. button, form, dialog 같은 컴포넌트 이름을 지정하세요.

### Foundation Technologies

React 19 기능: Server Components 지원, concurrent rendering features, automatic batching 개선, streaming SSR 향상.

TypeScript 5.5: 전체에 걸쳐 타입 안전, 향상된 추론, 더 나은 에러 메시지, 개선된 개발자 경험.

Tailwind CSS 3.4: JIT compilation, CSS variable 지원, dark mode variants, container queries.

Radix UI: 스타링되지 않은 접근성 가능한 primitives, keyboard navigation, focus management, ARIA 속성.

통합 스택: React Hook Form for 폼 상태 관리, Zod for 스키마 검증, class-variance-authority for variant management, Framer Motion for 애니메이션 라이브러리, Lucide React for 아이콘 라이브러리.

### AI-Powered Architecture Design

ShadcnUIArchitectOptimizer 클래스는 Context7 MCP 통합을 최적의 shadcn/ui 아키텍처를 설계합니다. Context7 client, component analyzer, theme optimizer를 초기화합니다. design_optimal_shadcn_architecture 메서드는 design system 요구사항을 받아 최신 shadcn/ui 및 React 문서를 Context7를 통해 가져오고 UI 컴포넌트와 사용자 요구를 기반으로 컴포넌트 선택을 최적화하고, 브랜드 가이드라인과 접근성 요구사항을 기반으로 테마 구성을 반환합니다.

최종 ShadcnUIArchitecture는 component library, theme system, accessibility compliance, performance optimization, integration patterns, customization strategy를 포함합니다.

### Best Practices

요구사항: CSS 변수를 테마 커스터마이제이션을 전용으로 사용하세요. 이는 동적 테마를 가능하게 하고, dark mode 전환을 지원하며, 모든 컴포넌트에서 디자인 시스템 일관성을 유지합니다. CSS 변수가 없으면 테마 변경에 코드 수정이 필요하고, dark mode가 실패하고 브랜드 커스터마이징이 유지 불가능합니다.

[HARD] 모든 대화형 요소에 접근성 속성을 포함하세요. 이는 WCAG 2.1 AA 준수를 준수하고, screen reader 호환성을 보장하며, 장애인 사용자를 포함합니다. 접근성 속성이 누�되면 장애인 사용자를 배제하고 법적合规 요구를 위반합니다.

[HARD] 모든 대화형 컴포넌트에 키보드 내비게이션을 구현하세요. 키보드 사용자에게 필수적인 탐색 방법을 제공하고, assistive technologies를 지원하며, 사용자 효율성을 개선합니다. 키보드 내비게이션 없으면 power user가 효율적으로 애플리케이션을 사용할 수 없고 접근성 준수에 실패합니다.

[SOFT] 비동기 작업을 위한 loading states를 제공하세요. 작업 진행률을 사용자에게 전달하고, perceived latency를 줄이며, 애플리케이션 응답성에 대한 사용자 자신감을 높입니다.

[HARD] 컴포넌트 트리 주변 에러 방지를 위해 error boundaries를 구현하세요. 고립된 컴포넌트 실패로 인한 전체 애플리케이션 충돌하는 것을 방지하고, 우아한 에러 복구를 가능하게 하며, 애플리케이션 안정성을 유지합니다.

[HARD] inline styles 대신 Tailwind CSS 클래스를 적용하세요. 이는 design system과 일관성을 유지하고, JIT compilation 이점을 활용하며, responsive design variants를 지원하며, bundle size 최적화를 개선합니다.

[SOFT] 모든 컴포넌트에 dark mode 지원을 구현하세요. 사용자 선호를 존중하고, 저조도 환경에서 눈의 스트레인을 줄이고, 모던 UI 기대와 일치합니다.

### Performance Optimization

Bundle Size 최적화: tree-shaking으로 사용하지 않은 컴포넌트 제거, large components를 code splitting, lazy loading과 React.lazy 사용, heavy dependencies에 대한 dynamic imports.

Runtime Performance 최적화: expensive components에 React.memo, computations에 useMemo와 useCallback 사용, large lists에 virtual scrolling, user interactions을 debouncing.

접근성: 모든 대화형 요소에 ARIA 속성, keyboard navigation support, focus management, screen reader testing.

---

## Advanced Patterns

### Component Composition

composable 패턴: Card, CardHeader, CardTitle, CardContent를 ui/card에서 임포트합니다. DashboardCard 컴포넌트는 title과 children prop를 받아 Card를 래핑합니다. CardHeader 안에 CardTitle을, CardContent 안에 children을 배치하여 구조를 형성합니다.

### Form Validation

Zod 및 React Hook Form 통합 패턴: useForm를 react-hook-form에서, zodResolver를 hookform/resolvers/zod에서, z에서 ZodObject schema를 임포트합니다. z.object를 사용하여 FormSchema를 정의하고 email 필드에 z.string().email()을, password 필드에 z.string().min(8)을 지정합니다. FormSchema에서 User 타입을 유추합니다. form 컴포넌트에서 useForm을 zodResolver(formSchema)와 함께 사용하여 form을 초기화하고 form.handleSubmit with onSubmit handler를 사용합니다.

---

## Works Well With

- shadcn-components.md 모듈 - 고급 컴포넌트 패턴 및 구현
- shadcn-theming.md 모듈 - 테마 시스템 및 커스터마이제이션 전략
- do-domain-uiux - design system 아키텍처 및 원칙
- do-lang-typescript - TypeScript 모범 사례
- code-frontend - frontend 개발 패턴

---

## Context7 Integration

관련 라이브러리:

- shadcn/ui at /shadcn-ui/ui provides reusable components built with Radix UI and Tailwind
- Radix UI at /radix-ui/primitives provides unstyled accessible component primitives
- Tailwind CSS at /tailwindlabs/tailwindcss provides utility-first CSS framework

공식 문서:

- shadcn/ui Documentation at ui.shadcn.com/docs
- API Reference at ui.shadcn.com/docs/components
- Radix UI Documentation at radix-ui.com
- Tailwind CSS Documentation at tailwindcss.com

Latest Versions as of November 2025:

- React 19
- TypeScript 5.5
- Tailwind CSS 3.4
- Radix UI Latest

---

Last Updated: 2026-01-11
Status: Production Ready
