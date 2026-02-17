---
name: do-library-nextra
description: >
  Next.js와 함께 Enterprise Nextra 문서화 프레임워크입니다.
  문서 사이트, 지식 베이스, API 참조 문서 구축 시 사용하세요.
license: Apache-2.0
compatibility: Designed for Claude Code
allowed-tools: Read Write Edit Grep Glob mcp__context7__resolve-library-id mcp__context7__get-library-docs
user-invocable: false
metadata:
  version: "2.2.0"
  category: "library"
  modularized: "true"
  status: "active"
  updated: "2026-01-11"
  tags: "library, nextra, nextjs, documentation, mdx, static-site"
  aliases: "do-library-nextra"

# MoAI Extension: Triggers
triggers:
  keywords: ["documentation", "nextra", "docs site", "knowledge base", "api reference", "mdx", "static site"]
---

## Quick Reference

목적: Nextra와 Next.js로 전문적인 문서 사이트를 구축하세요.

Nextra 장점:

- Zero config MDX - Markdown과 JSX의 무리미Seam 통합
- File-system routing - 자동 route 생성
- 코드 분할 및 프리패칭으로 성능 최적화
- 플러그인 가능하고 커스터마이짜 가능한 테마
- 내장된 국제화 지원

핵심 파일:

- pages 디렉토리는 MDX 형식의 문서 페이지를 포함합니다
- theme.config.tsx 파일은 사이트 구성을 포함합니다
- _meta.js 파일들이 네비게이션 구조를 제어합니다

## Implementation Guide

### Features

이 스킬은 Nextra 3.x 및 4.x 문서화 프레임워크 아키텍처 패턴, Next.js 14 및 15와의 최적 구성, theme.config.tsx 또는 Layout props를 통한 테마 커스터마이제이션, FlexSearch 통합한 고급 검색, internationalization 지원, React 컴포넌트가 포함된 MDX 콘텐츠, App Router 지원(Nextra 4.x), Turbopack 호환성을 포함합니다.

### When to Use

다음의 경우 이 스킬을 사용하세요:

- 모던 React 기능으로 문서 사이트 구축
- 고급 검색 기능이 있는 지식 베이스 생성
- multi-language 문서 포털 구축
- custom 문서 테마 구현
- 기술 docs에서 interactive 예제 통합

### Project Setup

Nextra 문서 사이트를 초기화하려면 create-nextra-app 명령어를 npx으로 사용하여 docs 템플릿을 지정하세요. 결과 프로젝트 구조는 pages 디렉토리에 App component 파일, index MDX 파일(홈 페이지), 섹션을 위한 하위 디렉토리를 포함합니다. 각 섹션에는 MDX 콘텐츠와 navigation 구성을 위한 _meta.json 파일이 포함됩니다.

### Theme Configuration

theme.config.tsx 파일에서 export하는 구성 객체의 여러 주요 속성: logo는 사이트 브랜딩 엘티티입니다. project는 프로젝트 저장소 링크입니다. docsRepositoryBase는 edit link 기능을 위한 기본 URL입니다. useNextSeoProps 함수는 title 표본을 포함한 SEO 구성을 반환합니다.

필수 구성 옵션: branding settings을 위한 logo와 logoLink, navigation settings을 위한 project links와 repository base URLs, sidebar settings을 위한 default collapse level과 toggle button visibility, table of contents settings을 위한 back-to-top 기능, footer settings를 위한 custom footer text를 포함합니다.

### Navigation Structure

_meta.js 파일들이 사이드바 메뉴 순서와 display 이름을 제어합니다. 각 파일은 키가 파일/디렉토리 이름이고 값이 display label인 객체를 default로 export합니다. triple dash를 키로 하고 빈 문자열 값을 사용하여 separator line을 추가합니다. external link는 title, href, newWindow 속성이 포함된 중첩 객체로 구성할 수 있습니다.

### MDX Content and JSX Integration

Nextra는 MDX 파일에서 Markdown과 React 컴포넌트를 직접 혼합할 수 있습니다. 파일 상단에서 컴포넌트를 import하고 Markdown 콘텐츠와 함께 인라인으로 사용할 수 있습니다. 컴포넌트를 정의하고 내보내서 export하거나 MDX 파일 자체 내에서 export하여 커스텀츠를 생성할 수 있습니다. nextra/components의 Callout 컴포넌트는 notes, warnings, tips를 위한 스타일링 박스를 제공합니다.

### Search and SEO Optimization

테마 구성에는 사용자 지정 placeholder가 있는 내장된 검색이 포함됩니다. SEO 메타데이터는 head 속성에 JSX를 사용하여 meta tags를 포함할 수 있습니다. Open Graph title, description, image를 포함합니다. useNextSeoProps 함수는 동적 title template 구성을 제공합니다.

---

## Advanced Documentation

이 스킬은 Progressive Disclosure를 사용합니다. 상세 패턴은 modules 디렉토리에서 확인하세요:

- modules/configuration.md 완전한 theme.config 참조
- modules/mdx-components.md MDX 컴포넌트 라이브러리
- modules/i18n-setup.ini 국제화 가이드
- modules/deployment.md 호스팅 및 배포

---

## Theme Options

내장된 테마:

- nextra-theme-docs - 문서 사이트용으로 권장
- nextra-theme-blog - 블로그 구현용

Customization options:

- 색상을 위한 CSS variables
- custom sidebar 컴포넌트
- footer customization
- navigation layout modifications

---

## Deployment

인기 배포 플폼:

- Vercel - zero-config 권장 설정
- GitHub Pages - 무료 self-hosted 옵션
- Netlify - flexible CI/CD 통합
- custom servers - full control

Vercel 배포의 경우: npm install -g로 Vercel CLI를 전역으로 설치하고 vercel 명령으로 프로젝트를 선택하고 배포하세요.

---

## Integration with Other Skills

상호 보완적인 스킬:

- do-docs_generation 코드에서 자동 문서 생성
- do-workflow-docs 문서 품질 검증
- do-cc-claude-md Markdown 형식

---

## Version History

버전 2.2.0 2026-01-11 릴리스: CLAUDE.md 문서 표준 준수를 준수하도록 재구성 - 코드 블록 제거, examples를 내러러리 텍스트 형식으로 변환

버전 2.1.0 2025-12-30: config.md에 완전한 Nextra별 theme.config.tsx 패턴 추가, Nextra 4.x App Router 구성 패턴 추가, Next.js 14 및 15 버전 호환 업데이트, Turbopack 지원 문서 추가

버전 2.0.0 2025-11-23: Progressive Disclosure로 리팩토, highlighted configuration patterns, MDX 통합 가이드 추가

버전 1.0.0 2025-11-12: 초기 Nextra 아키텍처 가이드, theme configuration, i18n 지원

---

Maintained by: MoAI-ADK Team
Domain: Documentation Architecture
Generated with: MoAI-ADK Skill Factory

---

## Works Well With

Agents:

- workflow-docs documentation generation
- code-frontend Nextra implementation
- workflow-spec architecture documentation

Skills:

- do-docs-generation content generation
- do-workflow-docs documentation validation
- do-library-mermaid diagram integration

Commands:

- moai:3-sync documentation deployment
- moai:0-project Nextra project initialization
