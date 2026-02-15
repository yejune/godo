# MoAI SPEC Philosophy: Deep Analysis

**Date**: 2026-02-15
**Analyst**: spec-analyst (team agent)
**Scope**: MoAI의 SPEC-first 접근법이 존재하는 이유, 그것이 대표하는 세계관, 그리고 그 한계

---

## 1. SPEC의 근본 철학: 왜 "사양서 먼저"인가?

### 핵심 주장

SPEC-first는 "AI 에이전트는 인간과 달리 암묵적 맥락을 공유하지 못한다"는 인식에서 출발한다. 사양서는 AI 에이전트 간의 공유 계약서이자, 토큰이라는 유한한 자원 안에서 정보 손실을 최소화하는 압축 매체다.

### 근거

MoAI의 SPEC은 EARS(Easy Approach to Requirements Syntax) 형식을 채택한다. EARS는 Rolls-Royce의 Alistair Mavin이 2009년에 개발한 요구사항 명세 문법으로, 2025년에 AWS Kiro IDE와 GitHub Spec-Kit이 채택하면서 산업 표준이 되었다.

> "EARS provides unambiguous, testable requirement syntax that eliminates interpretation errors."
> -- `/tmp/e2e4-extract/personas/moai/agents/moai/manager-spec.md:236`

SPEC 문서는 세 개의 파일로 구성된다:
- `spec.md`: EARS 형식의 5가지 요구사항 유형 (Ubiquitous, Event-driven, State-driven, Unwanted, Optional)
- `plan.md`: 구현 계획과 기술 접근법
- `acceptance.md`: Given/When/Then 형식의 인수 조건

> "MoAI is the Strategic Orchestrator for Claude Code. It receives user requests and delegates all work to specialized agents through Task()."
> -- `/tmp/e2e4-extract/personas/moai/skills/moai/SKILL.md:31`

### 철학적 해석

전통 소프트웨어 공학에서 "사양서 먼저"는 Waterfall의 유산이자, 과도한 선행 설계(Big Design Up Front)의 상징으로 비판받아 왔다. MoAI가 이것을 부활시킨 이유는 AI 에이전트 오케스트레이션이라는 새로운 맥락 때문이다.

인간 개발자 팀은 회의실에서 대화하고, 슬랙에서 질문하고, 코드 리뷰에서 의도를 읽는다. AI 에이전트에게는 이 모든 "암묵지"가 없다. 각 에이전트는 `Task()` 호출 시 전달받은 프롬프트가 세계의 전부다. SPEC은 이 "세계"를 표준화된 형식으로 명시하는 것이다.

EARS의 5가지 패턴은 자연어의 모호성을 제거한다:
- "The system **shall** [response]" (항상 참인 요구사항)
- "**When** [event], the system **shall** [response]" (트리거 기반)
- "**If** [undesired], **then** the system **shall** [response]" (금지 행동)

이것은 요구사항 언어를 프로그래밍 언어처럼 구조화하려는 시도다. 모호한 자연어를 "파싱 가능한" 형식으로 변환함으로써, 에이전트가 요구사항을 기계적으로 검증할 수 있게 만든다.

### 트레이드오프

SPEC 생성에 30K 토큰을 소비한다. 간단한 CSS 수정에도 EARS 형식의 사양서를 만들어야 한다면 이것은 과잉이다. MoAI는 이를 인식하고 있으며, `/moai fix`와 `/moai loop`이라는 SPEC-free 경로를 제공한다. 하지만 기본 워크플로우(`/moai`)는 항상 Plan -> Run -> Sync 파이프라인을 거친다. "언제 SPEC을 생략할 수 있는가?"의 판단이 오케스트레이터에게 위임되어 있다는 점에서, 과잉 설계의 위험은 오케스트레이터의 판단력에 의존한다.

---

## 2. 토큰 예산 전략: 200K를 30K+180K+40K로 나눈 근거

### 핵심 주장

Phase별 토큰 예산과 `/clear` 강제는 AI의 가장 근본적인 제약 -- 유한한 컨텍스트 윈도우 -- 에 대한 구조적 대응이다. 정보를 "버리는" 것이 아니라, 각 Phase의 산출물(SPEC, 코드, 문서)이라는 형태로 "물질화"시켜서 다음 Phase에 전달하는 것이다.

### 근거

> "Token Strategy: Allocation: 30,000 tokens. Load requirements only. Execute /clear after completion. Saves 45-50K tokens for implementation."
> -- `/tmp/e2e4-extract/personas/moai/rules/moai/workflow/spec-workflow.md:15-17`

> "Run Phase Token Strategy: Allocation: 180,000 tokens. Selective file loading. Enables 70% larger implementations."
> -- `/tmp/e2e4-extract/personas/moai/rules/moai/workflow/spec-workflow.md:33-35`

> "Sync Phase Token Strategy: Allocation: 40,000 tokens. Result caching. 60% fewer redundant file reads."
> -- `/tmp/e2e4-extract/personas/moai/rules/moai/workflow/spec-workflow.md:51-53`

Research report는 총 예산을 250K로 기술한다:
> "Total budget: 250K tokens across all phases: SPEC Phase: 30K, DDD Phase: 180K, Docs Phase: 40K"
> -- `research-moai-philosophy.md:300-303`

### 철학적 해석

이 분배는 수학적이라기보다 실용적이다.

**30K (Plan)**: SPEC 문서 생성은 주로 "읽기"(기존 프로젝트 문서, 코드 탐색)와 "쓰기"(EARS 문서 3개)로 구성된다. manager-spec은 product.md, structure.md, tech.md를 읽고 SPEC을 생성하면 된다. 30K는 이를 위한 최소한의 예산이다.

**180K (Run)**: 구현 Phase에 전체의 72%를 할당한 것은 의미심장하다. 코드를 읽고, 테스트를 작성하고, 구현하고, 검증하는 DDD/TDD 사이클은 가장 많은 토큰을 소비한다. `/clear`로 Plan Phase의 탐색 맥락을 모두 제거하고, SPEC 문서라는 "물질화된 맥락"만 남김으로써 구현에 최대한의 공간을 확보한다.

**40K (Sync)**: 문서화는 기존 산출물(코드, SPEC, git diff)을 참조하는 작업이다. 새로운 것을 만들기보다 기존 것을 정리하는 Phase이므로 상대적으로 적은 토큰이 필요하다.

`/clear`의 철학은 "기억의 물질화"다. 인간은 메모를 하고 책상을 정리한 후 다음 작업을 시작한다. `/clear`는 AI 에이전트의 "책상 정리"다. 중요한 것은 SPEC 문서와 코드 파일로 이미 디스크에 저장되어 있고, 나머지 탐색 과정의 맥락은 버려도 된다.

### 트레이드오프

Phase 간 정보 손실은 실제로 발생한다. Plan Phase에서 발견한 "이 코드는 이런 이유로 이렇게 짜여져 있다"는 통찰이 `/clear` 후 사라진다. Run Phase의 에이전트는 SPEC 문서에 적힌 것만 알 수 있다. SPEC에 적히지 않은 미묘한 맥락 -- 예를 들어 "이 레거시 코드는 건드리면 안 된다"는 구전 지식 -- 은 유실된다.

MoAI는 이를 "Implementation Divergence Tracking"으로 보완한다. Run Phase의 에이전트가 SPEC과 다르게 구현한 부분을 추적하고, Sync Phase에서 SPEC을 업데이트한다. 이것은 정보 손실을 인정하고, 사후에 교정하는 전략이다.

---

## 3. Progressive Disclosure 3단계: 토큰 경제학

### 핵심 주장

Progressive Disclosure는 "모든 지식을 항상 로드하면 토큰이 부족하고, 필요할 때만 로드하면 누락의 위험이 있다"는 딜레마에 대한 3단계 타협이다.

### 근거

> "Progressive Disclosure: Level 1: Metadata only (~100 tokens). Level 2: Skill body when triggered (~5000 tokens). Level 3: Bundled files on-demand."
> -- `/tmp/e2e4-extract/personas/moai/rules/moai/workflow/spec-workflow.md:77-79`

SKILL.md의 frontmatter에서 Progressive Disclosure가 어떻게 설정되는지:

```yaml
progressive_disclosure:
  enabled: true
  level1_tokens: 100
  level2_tokens: 5000
```
-- `/tmp/e2e4-extract/personas/moai/skills/moai/workflows/plan.md:19-22`

MoAI-ADK CLAUDE.md에서:
> "Level 1 (Metadata): ~100 tokens per skill, always loaded. Level 2 (Body): ~5K tokens, loaded when triggers match. Level 3 (Bundled): On-demand, Claude decides when to access."
> -- `/Users/max/Work/moai-adk/CLAUDE.md:301-303`

그리고 Research report에서 좀 더 구체적인 수치:
> "Level 1 (Quick): 30 seconds, ~1,000 tokens. Level 2 (Implementation): 5 minutes, ~3,000 tokens. Level 3 (Advanced): 10+ minutes, ~5,000 tokens."
> -- `research-moai-philosophy.md:315-318`

### 철학적 해석

50개 이상의 skill이 있다. 모든 skill의 Level 2를 로드하면 50 x 5,000 = 250,000 토큰 -- 컨텍스트 윈도우 전체를 skill만으로 채우게 된다. Level 1(메타데이터)만 로드하면 50 x 100 = 5,000 토큰으로 전체 skill 카탈로그를 "인덱스"할 수 있다.

이것은 운영체제의 **가상 메모리**와 같은 원리다. 모든 프로그램을 물리 메모리에 올리지 않고, 필요할 때 페이지를 로드하듯이, MoAI는 skill의 메타데이터만 "상주"시키고 본문은 트리거 조건이 일치할 때 로드한다.

트리거 시스템이 이 "페이지 폴트"의 역할을 한다:

```yaml
triggers:
  keywords: ["plan", "spec", "design", "architect"]
  agents: ["manager-spec", "Explore"]
  phases: ["plan"]
```

사용자의 요청에 "plan"이라는 키워드가 있으면 Level 2가 로드된다. 에이전트가 `manager-spec`이면 해당 skill이 로드된다. 이것은 정적 분석이 아니라 런타임 매칭이다.

### 트레이드오프

트리거가 누락되면 필요한 skill이 로드되지 않는다. "데이터베이스 최적화"를 요청했는데 키워드에 "optimization"만 있고 "database"가 없으면, 잘못된 skill이 로드될 수 있다. 이는 "거짓 부정"(필요한 지식의 누락)의 위험이다.

반대로, 트리거가 너무 넓으면 불필요한 skill이 로드되어 토큰을 낭비한다. "거짓 긍정"(불필요한 지식의 로드)의 비용은 토큰 낭비이지만, "거짓 부정"의 비용은 잘못된 구현이다. 비용 비대칭이 있으므로 트리거는 넓게 설정하는 것이 합리적이다 -- 그러나 MoAI는 토큰 예산도 타이트하기 때문에 이 균형이 미묘하다.

---

## 4. Agent Patches & Override Skills: 전문화의 가치

### 핵심 주장

"모든 에이전트에게 같은 규칙을 적용하면 각 에이전트의 전문성이 희석된다." Agent별 skill 주입(frontmatter의 `skills` 필드)은 각 에이전트에게 "전공 과목"을 부여하는 것이다.

### 근거

manager-spec의 frontmatter:
```
skills: moai-foundation-claude, moai-foundation-core, moai-foundation-philosopher, moai-workflow-spec, moai-workflow-project, moai-workflow-thinking, moai-lang-python, moai-lang-typescript
```
-- `/tmp/e2e4-extract/personas/moai/agents/moai/manager-spec.md:14`

manager-ddd의 frontmatter:
```
skills: moai-foundation-claude, moai-foundation-core, moai-foundation-quality, moai-workflow-ddd, moai-workflow-tdd, moai-workflow-testing, moai-tool-ast-grep
```
-- `/tmp/e2e4-extract/personas/moai/agents/moai/manager-ddd.md:17`

manager-quality의 frontmatter:
```
skills: moai-foundation-claude, moai-foundation-core, moai-foundation-quality, moai-workflow-testing, moai-tool-ast-grep, moai-workflow-loop
```
-- `/tmp/e2e4-extract/personas/moai/agents/moai/manager-quality.md:15`

### 철학적 해석

세 에이전트가 공유하는 skill은 `moai-foundation-claude`와 `moai-foundation-core` 뿐이다. 이것은 "헌법"에 해당한다 -- 모든 에이전트가 따라야 하는 기본 규칙이다.

그 위에:
- manager-spec은 `moai-foundation-philosopher`(전략적 사고), `moai-workflow-spec`(SPEC 작성법), `moai-lang-python/typescript`(기술 스택 판단)을 가진다. 이것은 "기획자의 전공 과목"이다.
- manager-ddd는 `moai-workflow-ddd`(리팩토링 방법론), `moai-tool-ast-grep`(AST 분석), `moai-workflow-testing`(테스트 작성법)을 가진다. 이것은 "구현자의 전공 과목"이다.
- manager-quality는 `moai-foundation-quality`(품질 기준), `moai-workflow-loop`(반복 수정)을 가진다. 이것은 "검수관의 전공 과목"이다.

이 설계는 인간 조직의 "직무 기술서"와 동일한 원리를 따른다. CEO에게 코딩 스킬을 주입할 필요가 없고, QA 엔지니어에게 요구사항 분석 스킬을 주입할 필요가 없다. 각 역할에 필요한 지식만 로드함으로써:

1. **토큰 효율**: 불필요한 skill 로드를 피한다
2. **전문성 강화**: 에이전트가 자신의 영역에 집중한다
3. **격리**: 구현 에이전트가 기획 에이전트의 규칙에 의해 혼란스러워지지 않는다

Hook 시스템도 에이전트별로 스코핑된다. manager-ddd에는 PreToolUse/PostToolUse hook이 있어서 매 코드 수정 전후에 검증이 실행되지만, manager-spec에는 SubagentStop hook만 있다. 이것은 "구현 에이전트에게는 안전벨트를, 기획 에이전트에게는 완료 확인을"이라는 차별화된 감독이다.

### 트레이드오프

에이전트 수 x skill 수의 조합이 폭발적으로 증가한다. 27개 에이전트와 50+ skill의 매핑을 관리하는 것은 운영 복잡성을 높인다. 새 skill을 추가할 때 어떤 에이전트에 주입해야 하는지 판단해야 하고, skill의 변경이 특정 에이전트에만 영향을 미치는지 확인해야 한다. 이것은 마이크로서비스 아키텍처의 "서비스 메쉬" 관리 문제와 유사하다.

---

## 5. SPEC 워크플로우 (Plan -> Run -> Sync): Phase 경계의 철학

### 핵심 주장

세 Phase는 소프트웨어 개발의 세 가지 인지 모드 -- "무엇을 만들 것인가"(Plan), "어떻게 만들 것인가"(Run), "만든 것을 어떻게 전달할 것인가"(Sync) -- 를 분리하여 각 모드에 최적화된 에이전트와 토큰 예산을 배정하는 것이다.

### 근거

전체 파이프라인:
> "Flow: Explore -> Plan -> Run -> Sync -> Done"
> -- `/tmp/e2e4-extract/personas/moai/skills/moai/workflows/moai.md:33`

Phase 전환의 규칙:
> "Plan to Run: Trigger: SPEC document approved. Action: Execute /clear, then /moai run SPEC-XXX."
> "Run to Sync: Trigger: Implementation complete, tests passing. Action: Execute /moai sync SPEC-XXX."
> -- `/tmp/e2e4-extract/personas/moai/rules/moai/workflow/spec-workflow.md:83-89`

Sync Phase가 항상 단일 에이전트인 이유:
> "Documentation generation is sequential by nature (README depends on code analysis, CHANGELOG depends on git history, PR depends on both). File outputs are few (3-5 files) with heavy interdependency. Token budget is small (40K) making team overhead wasteful. Single coherent voice produces better documentation quality."
> -- `/tmp/e2e4-extract/personas/moai/skills/moai/workflows/team-sync.md:9-14`

### 철학적 해석

**Plan Phase가 존재하는 이유**: "코드를 쓰기 전에 무엇을 쓸지 합의하라." 이것은 인간 팀에서도 유효한 원칙이지만, AI 에이전트에게는 더 절실하다. 에이전트는 "대충 시작해서 방향을 잡아가는" 것이 불가능하다. 토큰이 유한하기 때문이다. 잘못된 방향으로 50K 토큰을 소비하면 복구할 수 없다.

**Run Phase가 존재하는 이유**: 구현은 가장 많은 토큰을 소비한다. Plan Phase의 탐색 맥락을 `/clear`로 제거하고 SPEC이라는 "계약서"만 남김으로써, 구현 에이전트는 "무엇을 만들 것인가"에 대한 고민 없이 "어떻게 만들 것인가"에만 집중할 수 있다.

**Sync Phase가 별도인 이유**: 문서화를 구현과 분리한 것은 의도적이다. 구현 중에 문서를 쓰면 구현 에이전트의 토큰 예산을 잠식한다. 구현이 완료된 후, 새로운 컨텍스트에서 "만든 것을 정리"하는 것이 효율적이다. 또한 Sync Phase는 "SPEC-Implementation Divergence Analysis"를 수행한다. 계획과 실제 구현의 차이를 분석하여 SPEC을 업데이트하고, 프로젝트 문서를 동기화한다. 이것은 "사후 정산"이다.

**Phase 경계에서 `/clear`하는 이유**: 이것이 가장 급진적인 설계 결정이다. Phase 간의 모든 맥락을 버린다. SPEC 문서와 코드 파일만이 Phase 간의 "브릿지"다. 이것은 마치 릴레이 경주에서 바톤만 전달하는 것과 같다. 선수(에이전트)의 피로(토큰 소비)는 전달되지 않고, 바톤(SPEC/코드)만 전달된다.

### 트레이드오프

Phase 간 전환은 "콜드 스타트" 문제를 야기한다. Run Phase의 에이전트는 Plan Phase의 탐색 과정을 모른다. "왜 이 기술을 선택했는지"는 SPEC에 적혀 있지만, "왜 다른 기술을 기각했는지"는 적혀 있지 않을 수 있다. 이 "기각된 대안의 근거"가 사라지면, Run Phase의 에이전트가 이미 기각된 대안을 다시 고려할 위험이 있다.

---

## 6. 에러 유형별 에이전트 매핑 5종

### 핵심 주장

모든 에러를 같은 에이전트에게 재시도시키는 것은 "만능 의사에게 모든 환자를 보내는 것"과 같다. MoAI는 에러를 진단하고, 전문 에이전트에게 라우팅함으로써 해결 확률을 높인다.

### 근거

CLAUDE.md의 에러 핸들링 섹션:
> "Agent execution errors: Use expert-debug subagent. Token limit errors: Execute /clear, then guide user to resume. Permission errors: Review settings.json manually. Integration errors: Use expert-devops subagent. MoAI-ADK errors: Suggest /moai feedback."
> -- `/Users/max/Work/moai-adk/CLAUDE.md:271-275`

Fix 워크플로우의 에이전트 선택:
> "Level 1 (import, formatting): expert-backend or expert-frontend subagent. Level 2 (rename, type): expert-refactoring subagent. Level 3 (logic, API): expert-debug or expert-backend subagent (after user approval)."
> -- `/tmp/e2e4-extract/personas/moai/skills/moai/workflows/fix.md:96-98`

Loop 워크플로우의 에이전트 선택:
> "Type errors, logic bugs: expert-debug subagent. Import/module issues: expert-backend or expert-frontend subagent. Test failures: expert-testing subagent. Security issues: expert-security subagent. Performance issues: expert-performance subagent."
> -- `/tmp/e2e4-extract/personas/moai/skills/moai/workflows/loop.md:88-93`

### 철학적 해석

이것은 5종의 에러 매핑이 아니라, 사실상 두 차원의 분류 체계다:

**차원 1: 에러의 심각도** (fix.md의 Level 1-4)
- Level 1 (Immediate): 포맷팅, import 정렬 -- 자동 수정
- Level 2 (Safe): 변수명 변경, 타입 추가 -- 로그만 남기고 자동 수정
- Level 3 (Review): 로직 변경, API 수정 -- 사용자 승인 후 수정
- Level 4 (Manual): 보안 취약점, 아키텍처 변경 -- 자동 수정 불가

**차원 2: 에러의 도메인** (loop.md의 에이전트 매핑)
- 타입/로직 에러 -> expert-debug (디버깅 전문)
- 모듈/import 에러 -> expert-backend/frontend (도메인 전문)
- 테스트 실패 -> expert-testing (테스트 전문)
- 보안 이슈 -> expert-security (보안 전문)
- 성능 이슈 -> expert-performance (성능 전문)

일률적 재시도보다 나은 이유는 명확하다. "import error"를 security expert에게 보내면 토큰만 낭비하고 해결되지 않는다. 각 에이전트는 자신의 도메인에 특화된 skill을 로드하고, 해당 도메인의 패턴을 알고 있다. 이것은 병원의 진료 과목 분류와 같다.

Team Debug 모드에서는 더 흥미로운 패턴이 등장한다:
> "Formulate 2-3 competing hypotheses. Each teammate explores a different theory independently."
> -- `/tmp/e2e4-extract/personas/moai/skills/moai/workflows/team-debug.md:16-17`

복잡한 버그에 대해 3개의 가설을 동시에 탐색하고, 증거를 종합하여 가장 가능성 높은 원인을 찾는다. 이것은 의학의 "감별 진단"과 동일한 원리다.

### 트레이드오프

에이전트 라우팅 자체가 토큰을 소비한다. 에러를 분류하고, 적절한 에이전트를 선택하고, 컨텍스트를 전달하는 과정에서 오버헤드가 발생한다. 단순한 에러(예: 누락된 세미콜론)에 대해 "분류 -> 라우팅 -> 에이전트 호출" 파이프라인을 거치는 것은 직접 수정하는 것보다 비효율적이다.

---

## 7. MCP 통합 (Context7 등): 내부 지식과 외부 지식의 역할 분담

### 핵심 주장

MoAI의 skill은 "검증된 패턴"(안정적이지만 오래된), MCP/Context7은 "최신 정보"(신선하지만 미검증)를 제공한다. 이 둘을 결합하여 "할루시네이션 없는 코드 생성"을 추구한다.

### 근거

> "Hallucination-Free Code Generation Process: 1. Load Relevant Skills (proven patterns). 2. Query Context7 (latest API versions). 3. Combine Both (merge stability with freshness). 4. Cite Sources (every pattern has attribution). 5. Include Tests (follow Skill test patterns)."
> -- `research-moai-philosophy.md:648-652`

MCP 서버 목록:
> "Sequential Thinking: Complex problem analysis, architecture decisions. Context7: Up-to-date library documentation. Pencil: UI/UX design editing. claude-in-chrome: Browser automation."
> -- `/Users/max/Work/moai-adk/CLAUDE.md:288-291`

### 철학적 해석

이것은 "두 가지 지식 소스의 장단점 보완"이라는 단순한 이야기가 아니다. 더 깊은 철학이 있다.

AI 에이전트의 훈련 데이터에는 "특정 시점까지의 세계"가 들어 있다. 라이브러리 버전, API 시그니처, 베스트 프랙티스는 모두 과거 시점의 것이다. 사용자가 "FastAPI 최신 버전으로 API를 만들어줘"라고 하면, 에이전트는 훈련 데이터의 버전을 참조한다 -- 이것이 할루시네이션의 원인이다.

MoAI의 해결책은 이중 참조 체계다:
1. **Skills** (내부, 정적): 프로젝트 내에서 검증된 패턴. "우리 프로젝트에서는 이렇게 한다."
2. **Context7 MCP** (외부, 동적): 최신 라이브러리 문서. "현재 FastAPI 버전은 0.118이다."

이 둘을 결합하면: "우리 프로젝트의 패턴을 따르되, API는 최신 버전을 사용한다."

Sequential Thinking MCP(`--ultrathink`)는 또 다른 축이다. 이것은 외부 지식이 아니라 "깊은 사고"를 위한 도구다. 복잡한 아키텍처 결정에서 chain-of-thought를 확장하기 위해 사용된다.

### 트레이드오프

MCP 의존성은 외부 서비스 가용성에 의존한다. Context7 서버가 다운되면 최신 정보를 얻을 수 없다. MoAI는 이를 "graceful fallback"으로 처리하지만, 폴백 시 생성되는 코드의 품질은 떨어질 수 있다.

또한 MCP 호출은 지연 시간을 추가한다. 매 skill 로드 시 Context7을 조회하면 워크플로우가 느려진다. 이 때문에 MCP 통합은 skill의 `allowed-tools` 필드로 제어되며, 필요한 skill만 MCP에 접근한다.

---

## 8. TAG Chain: Task-Assigned Groups와 의존성 그래프

### 핵심 주장

TAG Chain은 "구현 순서를 코드에 임베드"하여, 에이전트가 올바른 순서로 작업을 실행하고, 완료 여부를 기계적으로 검증할 수 있게 하는 메커니즘이다.

### 근거

> "TAG: Task-Assigned Group -- implementation unit in a chain. TAG Chain: Sequence of TAGs with dependencies forming implementation order."
> -- `research-moai-philosophy.md:447-448`

manager-quality의 TAG chain 검증:
> "3.3 TAG chain verification: 1. Explore TAG comments (extract TAG list by file). 2. TAG order verification (compare with TAG order in implementation-plan, check missing TAG, check wrong order). 3. Check feature completion conditions."
> -- `/tmp/e2e4-extract/personas/moai/agents/moai/manager-quality.md:239-253`

Run workflow의 Task Decomposition:
> "Tasks for manager-strategy: Decompose plan into atomic implementation tasks. Each task must be completable in a single DDD/TDD cycle. Assign priority and dependencies for each task. Generate task tracking entries."
> -- `/tmp/e2e4-extract/personas/moai/skills/moai/workflows/run.md:109-113`

Team Run의 Task 의존성:
> "TaskCreate: 'Implement data models and schema' (no deps). TaskCreate: 'Implement API endpoints' (blocked by data models). TaskCreate: 'Implement UI components' (blocked by API endpoints). TaskCreate: 'Write unit and integration tests' (blocked by API + UI). TaskCreate: 'Quality validation - TRUST 5' (blocked by all above)."
> -- `/tmp/e2e4-extract/personas/moai/skills/moai/workflows/team-run.md:36-40`

### 철학적 해석

단순 리스트(1, 2, 3, ...)는 "순서"만 표현한다. TAG Chain은 "의존성 그래프"를 표현한다. 차이는 근본적이다.

단순 리스트에서 task 3이 실패하면, task 4와 5도 순차적으로 실패한다 -- task 4가 task 3에 의존하는지 여부와 무관하게. TAG Chain에서는 task 4가 task 3에 의존하지 않는다면, task 3이 실패해도 task 4는 실행 가능하다.

이것이 **병렬 실행의 기반**이다. Team Run 워크플로우에서 backend-dev와 frontend-dev가 동시에 작업할 수 있는 이유는, 두 에이전트의 task가 서로 독립적이거나, 의존성이 명확히 정의되어 있기 때문이다.

TAG Chain은 또한 **완료 검증의 기반**이다. manager-quality가 TAG order를 검증하는 것은 "모든 조각이 올바른 순서로 구현되었는가?"를 기계적으로 확인하는 것이다. 누락된 TAG는 누락된 기능이다. 순서가 뒤바뀐 TAG는 의존성 위반이다.

### 트레이드오프

TAG Chain의 오버헤드는 의존성 분석 자체에 있다. manager-strategy가 SPEC을 분석하여 task를 분해하고 의존성을 결정하는 Phase 1.5("Task Decomposition")는 추가 토큰을 소비한다. 작은 기능(1-2개 파일 변경)에 대해 의존성 그래프를 만드는 것은 과잉이다.

또한 의존성 추정이 잘못될 수 있다. "API가 먼저 구현되어야 UI를 만들 수 있다"는 것은 대체로 맞지만, mock API로 UI를 먼저 만들 수도 있다. 의존성 그래프가 지나치게 직렬적이면 병렬 실행의 이점을 잃는다.

---

## 9. 완료 마커 (`<moai>DONE</moai>`): 명시적 완료 신호의 철학

### 핵심 주장

AI 에이전트에게 "완료"는 자명하지 않다. 인간은 직감으로 "끝났다"를 느끼지만, AI는 명시적인 완료 조건과 신호가 필요하다. 완료 마커는 이 "완료의 모호성"을 해결하는 프로토콜이다.

### 근거

> "AI uses markers to signal task completion: `<moai>DONE</moai>` - Task complete. `<moai>COMPLETE</moai>` - Full completion."
> -- `/tmp/e2e4-extract/personas/moai/rules/moai/workflow/spec-workflow.md:65-66`

> "These markers enable automation detection of workflow state."
> -- `/tmp/e2e4-extract/personas/moai/skills/moai/SKILL.md:222`

Loop 워크플로우에서의 활용:
> "Step 1 - Completion Check: Check for completion marker in previous iteration response. Marker types: `<moai>DONE</moai>`, `<moai>COMPLETE</moai>`. If marker found: Exit loop with success."
> -- `/tmp/e2e4-extract/personas/moai/skills/moai/workflows/loop.md:51-53`

Hook 시스템과의 연동:
- manager-spec의 SubagentStop hook: `handle-agent-hook.sh spec-completion`
- manager-ddd의 SubagentStop hook: `handle-agent-hook.sh ddd-completion`
- manager-quality의 SubagentStop hook: `handle-agent-hook.sh quality-completion`

### 철학적 해석

완료 마커가 해결하는 문제는 세 가지다:

**1. 루프 종료 조건**: `/moai loop`은 최대 100회 반복하며 에러를 수정한다. "언제 멈추는가?" 에러 카운트가 0이 되면? 테스트가 모두 통과하면? 커버리지가 85%를 넘으면? 이 모든 조건을 `<moai>DONE</moai>`이라는 단일 신호로 통합한다.

**2. Hook 연동**: SubagentStop hook은 에이전트가 종료될 때 실행된다. Go 바이너리(`moai hook spec-completion`)가 완료 마커를 감지하고, 후속 처리(상태 업데이트, 다음 Phase 트리거)를 수행한다. XML 태그는 사용자에게 노출되지 않지만(HARD rule), 기계 간 통신에서는 파싱이 용이한 구분자 역할을 한다.

**3. 인간 감독 포인트**: 완료 마커는 "여기서 멈추고 인간에게 보고하라"는 신호이기도 하다. SPEC 완료 후 사용자에게 "구현 진행할까요?"를 묻고, 구현 완료 후 "문서화 진행할까요?"를 묻는 것은 모두 완료 마커가 트리거하는 AskUserQuestion이다.

`DONE`과 `COMPLETE`의 구분은 미묘하다. `DONE`은 개별 task 완료, `COMPLETE`는 전체 워크플로우 완료다. 이것은 Unix의 exit code(0 = 성공)처럼, 단순하지만 필수적인 "완료의 문법"이다.

### 트레이드오프

완료 마커가 잘못 생성되면 워크플로우가 조기 종료된다. 에이전트가 "아직 더 할 게 있지만 토큰이 부족해서" `DONE`을 출력하면, 미완성된 작업이 완료로 기록된다. 이것은 인간 개발자가 금요일 퇴근 전에 "Done"이라고 Jira 티켓을 닫는 것과 같은 문제다 -- 프로토콜이 있어도 신뢰성은 실행자에 의존한다.

---

## 10. SPEC의 한계: 과도한 상황과 오버헤드의 비용

### 핵심 주장

SPEC-first는 "중규모 이상의 기능 개발"에 최적화되어 있으며, 작은 버그 수정, 설정 변경, CSS 조정에는 과도한 오버헤드를 발생시킨다. MoAI는 이를 인식하고 있으나, 기본 경로가 SPEC을 거치도록 설계되어 있어 "언제 생략하는가?"의 판단이 중요하다.

### 근거

단일 도메인 작업의 SPEC 우회:
> "Single-domain routing: If task is single-domain (e.g., 'SQL optimization'): Delegate directly to expert agent, skip SPEC generation."
> -- `/tmp/e2e4-extract/personas/moai/skills/moai/workflows/moai.md:111-112`

SPEC-free 워크플로우의 존재:
- `/moai fix`: SPEC 없이 에러를 스캔하고 수정한다. Level 1-2는 자동, Level 3은 승인 후.
- `/moai loop`: SPEC 없이 반복적으로 에러를 수정한다. 최대 100회 반복.

이 두 워크플로우는 SPEC을 생략한다. 이것은 MoAI가 "모든 것에 SPEC이 필요한 것은 아니다"를 인정하는 것이다.

Run Phase에서 methodology 선택의 분기:
> "Mode | Workflow Cycle | Best For: DDD | ANALYZE-PRESERVE-IMPROVE | Existing projects, < 10% coverage. TDD | RED-GREEN-REFACTOR | New projects, 50%+ coverage. Hybrid | Mixed per change type | Partial coverage (10-49%)."
> -- `/tmp/e2e4-extract/personas/moai/rules/moai/workflow/workflow-modes.md:9-15`

### 철학적 해석

SPEC의 과잉이 발생하는 상황:

**1. 한 줄 버그 수정**: `if (x > 0)` -> `if (x >= 0)`. 이것에 EARS 형식의 요구사항, 인수 조건, 구현 계획이 필요한가? 아니다. `/moai fix`가 이 경우를 처리한다.

**2. CSS 조정**: 버튼 색상을 `#3498db`에서 `#2ecc71`로 변경. 이것에 5가지 EARS 패턴이 필요한가? 아니다.

**3. 설정 변경**: `docker-compose.yml`의 포트 변경. SPEC이 아니라 그냥 수정하면 된다.

**4. 문서 업데이트**: README의 오타 수정. SPEC은 물론이고 어떤 워크플로우도 필요 없다.

MoAI가 이 한계를 어떻게 처리하는가:

| 작업 유형 | MoAI 경로 | SPEC 필요 |
|----------|----------|----------|
| 신규 기능 (5+ 파일) | /moai plan -> run -> sync | 필수 |
| 중규모 기능 (2-4 파일) | /moai (자동 파이프라인) | 자동 생성 |
| 단일 도메인 수정 | Direct expert 위임 | 생략 |
| 에러 수정 | /moai fix | 생략 |
| 반복 수정 | /moai loop | 생략 |

핵심 문제는 **기본 경로의 편향**이다. 사용자가 `/moai "로그인 버튼 색상 변경"`이라고 하면, Intent Router의 Priority 3(자연어 분류)이 이것을 "implementation"으로 분류하고, 기본 파이프라인(Plan -> Run -> Sync)을 실행할 수 있다. 사용자가 `/moai fix`를 명시적으로 호출해야 SPEC을 우회할 수 있다.

이 편향은 의도적인 것일 수 있다. "과잉 문서화"가 "과소 문서화"보다 안전하다는 판단이다. 하지만 비용은 토큰이다. 불필요한 SPEC 생성에 30K 토큰을 쓰면, 실제 구현에 사용할 토큰이 줄어든다.

### SPEC 과잉의 실질적 비용

1. **토큰 낭비**: Plan Phase 30K + /clear 오버헤드. 단순 작업에는 이 30K가 구현에 쓰이는 편이 낫다.
2. **사용자 인터랙션 피로**: SPEC 승인, 구현 계획 승인, 다음 단계 선택 등 최소 3번의 AskUserQuestion. 단순 작업에는 "그냥 해줘"가 더 나은 UX다.
3. **시간 지연**: SPEC 생성 -> 승인 -> /clear -> 구현 -> 검증 -> 문서화의 파이프라인은 최소 3개의 에이전트 호출을 필요로 한다. 단순 작업에는 1개의 expert 에이전트 호출로 충분하다.

### 더 근본적인 한계

SPEC-first의 가장 근본적인 한계는 **요구사항이 불확실한 상황**에서의 적용이다. 스타트업 초기 단계에서 "사용자가 뭘 원하는지 아직 모르는" 상태에서 EARS 형식의 사양서를 작성하는 것은 Waterfall의 함정에 빠지는 것이다. "빠르게 프로토타입하고, 피드백을 받고, 방향을 수정하는" 애자일 접근법과 SPEC-first는 근본적으로 긴장 관계에 있다.

MoAI는 이것을 `/moai fix`와 `/moai loop`이라는 "반복 수정 경로"로 부분적으로 해결하지만, 이 경로들은 "이미 존재하는 코드의 에러 수정"에 특화되어 있고, "요구사항 탐색을 위한 빠른 프로토타이핑"에는 적합하지 않다.

---

## 총평: SPEC의 세계관

MoAI의 SPEC-first는 단순한 방법론이 아니라 하나의 세계관이다. 그 세계관은 다음과 같이 요약할 수 있다:

**"AI 에이전트 오케스트레이션에서, 구조화된 사양서는 에이전트 간의 유일한 신뢰할 수 있는 소통 매체다."**

이 세계관의 전제:
1. AI 에이전트는 암묵지를 공유하지 못한다
2. 토큰은 유한하며, Phase 간에 맥락을 "물질화"해야 한다
3. 에이전트는 전문화될수록 효과적이다
4. 품질은 자동으로 검증 가능한 기준으로 측정해야 한다
5. 완료는 명시적으로 선언되어야 한다

이 세계관이 참인 영역 -- 중규모 이상의 기능 개발, 멀티 도메인 작업, 팀 기반 병렬 개발 -- 에서 SPEC-first는 강력하다. 이 세계관이 과잉인 영역 -- 단순 수정, 프로토타이핑, 요구사항 탐색 -- 에서는 `/moai fix`와 `/moai loop`이 탈출구를 제공한다.

MoAI의 진짜 기여는 SPEC 형식 자체가 아니라, **"언제 SPEC을 쓰고, 언제 생략하는가?"의 라우팅 시스템** -- Intent Router와 Phase별 에이전트 매핑 -- 에 있다. SPEC은 도구이고, 그 도구를 언제 꺼내는지 아는 것이 오케스트레이터의 지혜다.

---

**Analyst**: spec-analyst
**Files Analyzed**: 20+
**Sources Cited**: 30+ direct quotations from source files
