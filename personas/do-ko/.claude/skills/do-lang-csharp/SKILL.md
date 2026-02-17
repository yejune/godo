---
name: do-lang-csharp
description: >
  C# 12 / .NET 8 개발 전문가 - ASP.NET Core, Entity Framework, Blazor, 모던 C# 패턴涵盖.
  .NET API, 웹 애플리케이션, 엔터프라이즈 솔루션 개발 시 사용하세요.
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
  tags: "language, csharp, dotnet, aspnet-core, entity-framework, blazor"
  context7-libraries: "/dotnet/aspnetcore, /dotnet/efcore, /dotnet/runtime"

# MoAI Extension: Progressive Disclosure
progressive_disclosure:
  enabled: true
  level1_tokens: 100
  level2_tokens: 5000

# MoAI Extension: Triggers
triggers:
  keywords: ["C#", "Csharp", ".NET", "ASP.NET", "Entity Framework", "Blazor", ".cs", ".csproj", ".sln", "dotnet"]
  languages: ["csharp", "c#"]
---

# C# 12 / .NET 8 개발 전문가

ASP.NET Core, Entity Framework Core, Blazor, 엔터프라이즈 패턴과 모던 C# 개발.

## Quick Reference

자동 트리거: .cs, .csproj, .sln 파일, C# 프로젝트, .NET 솔루션, ASP.NET Core 애플리케이션

핵심 스택:

- C# 12: 기본 생성자, 컬렉션 식, 별칭 any 타입, 기본 lambda 파라미터
- .NET 8: Minimal API, 네이티브 AOT, 향상된 성능, WebSocket
- ASP.NET Core 8: Controller, 엔드포인트, 미들웨어, 인증
- Entity Framework Core 8: DbContext, 마이그레이션, LINQ, 쿼리 최적화
- Blazor: Server/WASM 컴포넌트, InteractiveServer, InteractiveWebAssembly
- 테스트: xUnit, NUnit, FluentAssertions, Moq

빠른 명령어:

새 .NET 8 웹 API 프로젝트를 생성하려면 dotnet new webapi -n 프로젝트이름 --framework net8.0을 실행하세요.

Blazor 웹 앱을 생성하려면 dotnet new blazor -n 프로젝트이름 --interactivity Auto를 실행하세요.

Entity Framework Core를 추가하려면 dotnet add package Microsoft.EntityFrameworkCore.SqlServer 다음에 Microsoft.EntityFrameworkCore.Design을 실행하세요.

FluentValidation과 MediatR을 추가하려면 dotnet add package FluentValidation.AspNetCore 다음에 dotnet add package MediatR을 실행하세요.

---

## Module Index

이 스킬은 전문涵盖을 위해 점진적 공개를 사용하는 특수화된 모듈을 사용합니다.

### Language Features

- [C# 12 Features](modules/csharp12-features.md) - 기본 생성자, 컬렉션 식, 타입 별칭, 기본 lambda

### Web Development

- [ASP.NET Core 8](modules/aspnet-core.md) - Minimal API, Controller, 미들웨어, 인증
- [Blazor Components](modules/blazor-components.md) - Server, WASM, InteractiveServer, 컴포넌트

### Data Access

- [Entity Framework Core 8](modules/efcore-patterns.md) - DbContext, 리포지토리 패턴, 마이그레이션, 쿼리 최적화

### Architecture Patterns

- [CQRS and Validation](modules/cqrs-validation.md) - MediatR CQRS, FluentValidation, 핸들러 패턴

### Reference Materials

- [API Reference](reference.md) - 완전한 API 참조, Context7 라이브러리 매핑
- [Code Examples](examples.md) - 프로덕션 준비 예제, 테스트 템플릿

---

## Implementation Quick Start

### Project Structure (Clean Architecture)

src 폴더에 4개의 주요 프로젝트로 프로젝트를 구성하세요. MyApp.Api는 ASP.NET Core 웹 API 레이어를 포함합니다: API Controller용 Controllers 폴더, Minimal API 엔드포인트용 Endpoints 폴더, 애플리케이션 진입점인 Program.cs. MyApp.Application은 비즈니스 로직을 포함합니다: CQRS Commands용 Commands 폴더, CQRS Queries용 Queries 폴더, FluentValidation용 Validators 폴더. MyApp.Domain은 도메인 엔티티를 포함합니다: 도메인 모델용 Entities 폴더, 리포지토리 인터페이스용 Interfaces 폴더. MyApp.Infrastructure는 데이터 액세스를 포함합니다: DbContext용 Data 폴더, 리포지토리 구현용 Repositories 폴더.

### Essential Patterns

DI와 함께 기본 생성자: IUserRepository와 ILogger<UserService>를 위한 생성자 파라미터로 public class UserService를 정의하세요. GetByIdAsync(Guid id)와 같은 async 메서드를 생성하고 UserId로 구조화된 로깅과 함께 logger를 사용하여 정보를 기록하며 repository.FindByIdAsync에서 결과를 반환하세요.

Minimal API 엔드포인트: "/api/users/{id:guid}"와 같은 경로 패턴과 Guid id 및 IUserService를 받는 async lambda로 app.MapGet을 사용하세요. 서비스 메서드를 호출하고 null 결과인 경우 Results.Ok를, 찾은 엔티티인 경우 Results.NotFound를 반환하세요. 경로 명명을 위해 WithName을 체이닝하고 OpenAPI 문서를 위해 WithOpenApi를 체이닝하세요.

엔티티 구성: 엔티티 타입을 위한 IEntityTypeConfiguration<T>를 구현하는 클래스를 생성하세요. EntityTypeBuilder<T>를 받는 Configure 메서드에서 기본 키를 설정하기 위해 HasKey를 호출하고, HasMaxLength와 IsRequired로 필드를 구성하기 위해 Property를 사용하며, 고유 제약 조건을 위해 HasIndex와 IsUnique를 사용하세요.

---

## Context7 Integration

최신 문서의 경우 Context7 MCP 도구를 사용하세요.

ASP.NET Core 문서의 경우 "aspnetcore"로 mcp__context7__resolve-library-id를 사용하여 라이브러리 ID를 확인한 다음 "minimal-apis middleware"와 같은 주제와 확인된 라이브러리 ID로 mcp__context7__get-library-docs를 호출하세요.

Entity Framework Core 문서의 경우 "efcore"로 확인하고 "dbcontext migrations"와 같은 주제로 가져오세요.

.NET Runtime 문서의 경우 "dotnet runtime"으로 확인하고 "collections threading"과 같은 주제로 가져오세요.

---

## Quick Troubleshooting

빌드 및 런타임: 상세한 출력을 위해 dotnet build --verbosity detailed를 실행하세요. HTTPS 프로필로 실행하려면 dotnet run --launch-profile https를 실행하세요. EF 마이그레이션을 적용하려면 dotnet ef database update를 실행하세요. 새 마이그레이션을 생성하려면 dotnet ef migrations add 마이그레이션이름을 실행하세요.

일반적인 패턴:

null 참조 처리의 경우 컨텍스트에서 가져온 후 ArgumentNullException.ThrowIfNull을 변수와 nameof 표현식과 함께 사용하세요.

async enumerable 스트리밍의 경우 IAsyncEnumerable<T>를 반환하는 async 메서드를 생성하세요. CancellationToken 파라미터에 EnumeratorCancellation 속성을 추가하세요. AsAsyncEnumerable와 WithCancellation이 있는 await foreach를 사용하여 반복하고 각 항목을 yield하세요.

---

## Works Well With

- `do-domain-backend` - API 설계, 데이터베이스 통합 패턴
- `do-platform-deploy` - Azure, Docker, Kubernetes 배포
- `do-workflow-testing` - 테스트 전략 및 패턴
- `do-foundation-quality` - 코드 품질 표준
- `do-essentials-debug` - .NET 애플리케이션 디버깅
