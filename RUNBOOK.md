# Convert Runbook: moai 분해 및 페르소나 관리 가이드

## 1. 개요

### 1.1 이 문서의 목적

이 문서는 moai-adk(upstream 참조 구현)를 core + persona로 분해하고, do 페르소나를 관리하며, moai 버전업 시 변경사항을 do에 반영하는 전체 운영 프로세스를 정의한다.

### 1.2 전체 아키텍처

```
moai-adk (.claude/)                    do-focus (.claude/)
  |                                      ^
  v                                      |
[Extract]                             [Assemble]
  |                                      |
  v                                      |
core/  ─────────────────────────────>  core/ (공유)
  + registry.yaml                        +
  + agents/ (22개 core)                personas/do/ (do 고유)
  + rules/ (공유 규칙)                   + CLAUDE.md
  + skills/ (공유 스킬 40+)              + manifest.yaml
                                         + settings.json
personas/moai/ ──(참고)──>               + agents/do/ (6개)
  + CLAUDE.md                            + skills/do/ (SKILL.md + workflows)
  + manifest.yaml                        + commands/do/ (6개)
  + agents/moai/ (6개)                   + rules/do/ (2개)
  + skills/moai/ (SKILL.md + 13)         + output-styles/do/ (3개)
  + commands/moai/ (2개)
  + hooks/moai/ (7개 .sh)
  + output-styles/moai/ (3개)
  + rules/moai/workflow/ (2개)
```

### 1.3 핵심 개념

**Core (공유 계층)**
- 브랜드에 무관한 범용 파일: 언어별 rules, 도메인 skills, core agents
- `{{slot:BRAND}}` 등의 슬롯 마커로 브랜드 참조를 추상화
- moai든 do든 동일한 core를 사용

**Persona (브랜드별 계층)**
- 특정 방법론/워크플로우를 정의하는 파일
- CLAUDE.md, 오케스트레이터 스킬, 커맨드, 훅, 스타일 등
- manifest.yaml이 persona의 모든 구성을 선언

**Slot 변수**
| 변수 | moai 값 | do 값 | 용도 |
|------|---------|-------|------|
| `{{slot:BRAND}}` | moai | do | 브랜드명 |
| `{{slot:BRAND_DIR}}` | moai | do | 디렉토리 접두사 |
| `{{slot:BRAND_CMD}}` | moai | do | 커맨드 접두사 |

**Slot 유형**
| 유형 | 문법 | 용도 |
|------|------|------|
| 인라인 슬롯 | `{{slot:SLOT_ID}}` | 텍스트 내 브랜드 참조 치환 |
| 섹션 슬롯 | `<!-- BEGIN_SLOT:ID -->...<!-- END_SLOT:ID -->` | 블록 단위 페르소나 콘텐츠 |
| 프로젝트 변수 | `{{VAR_NAME}}` | converter가 처리하지 않는 프로젝트 변수 |

---

## 2. 버전 관리 전략

### 2.1 moai 버전 핀닝

moai-adk의 특정 git 버전(태그 또는 커밋 해시)을 기준으로 분해한다. 어떤 버전에서 분해했는지를 `versions.yaml`로 추적한다.

**파일 위치**: `~/Work/new/convert/versions.yaml`

```yaml
moai:
  repo: ~/Work/moai-adk
  version: v3.0.0          # git tag 또는 commit hash
  extracted_at: 2026-02-15  # 분해 실행일
  core_files: 395           # core 파일 수 (registry.yaml 기준)
  persona_files: 61         # persona 파일 수 (manifest.yaml 기준)
  notes: "초기 분해"

do:
  based_on_moai: v3.0.0     # 어떤 moai 버전 기반인지
  version: "3.0.0"          # do persona 자체 버전
  updated_at: 2026-02-15
  notes: "초기 생성"

history:
  - date: 2026-02-15
    moai_version: v3.0.0
    action: initial_extract
    changes: "초기 분해 및 do persona 생성"
```

### 2.2 버전업 감지

moai-adk 저장소의 변경을 주기적으로 확인한다.

```bash
# 현재 핀닝된 버전 확인
cd ~/Work/moai-adk
git log --oneline -1

# 새 태그/릴리즈 확인
git fetch --tags
git tag --sort=-v:refname | head -5

# 핀닝 버전과 최신 버전 간 diff
git diff v3.0.0..v3.1.0 --stat -- .claude/
git diff v3.0.0..v3.1.0 --name-only -- .claude/
```

### 2.3 변경 영향도 평가 기준

| 영향도 | 기준 | 예시 |
|--------|------|------|
| **HIGH** | CLAUDE.md 구조 변경, SKILL.md 변경, manifest 스키마 변경, 새 슬롯 패턴 | 새 Intent Router 추가, 새 필드 도입 |
| **MEDIUM** | persona 에이전트/스킬 추가/삭제, workflow 추가, settings 구조 변경 | 새 workflow 파일, hook 이벤트 추가 |
| **LOW** | core 파일만 변경, 문서 수정, 기존 파일 내용 보강 | 언어 rules 업데이트, skill 모듈 추가 |

---

## 3. 초기 분해 프로세스 (최초 1회)

### 3.1 사전 준비

```bash
# converter 빌드
cd ~/Work/new/convert
go build -o convert ./cmd/convert

# moai-adk 버전 확인 및 기록
cd ~/Work/moai-adk
git log --oneline -1
# 출력 예: a1b2c3d (tag: v3.0.0) feat: release v3.0.0
```

### 3.2 moai-adk 분해 (Extract)

```bash
cd ~/Work/new/convert

# Extract 실행
./convert extract \
  --src ~/Work/moai-adk/.claude \
  --out ./extracted-moai-v3.0.0

# 또는 GitHub에서 직접 (로컬 클론이 없을 때)
./convert extract \
  --repo org/moai-adk \
  --branch v3.0.0 \
  --out ./extracted-moai-v3.0.0
```

### 3.3 분해 결과 확인

```bash
# core 구조 확인
ls ./extracted-moai-v3.0.0/core/
# 기대: agents/ rules/ skills/ registry.yaml

# core 파일 수 확인
find ./extracted-moai-v3.0.0/core/ -type f | wc -l

# persona 구조 확인
ls ./extracted-moai-v3.0.0/personas/moai/
# 기대: manifest.yaml CLAUDE.md agents/ skills/ commands/ hooks/ output-styles/ rules/

# manifest 검토
cat ./extracted-moai-v3.0.0/personas/moai/manifest.yaml

# 슬롯 레지스트리 확인
cat ./extracted-moai-v3.0.0/core/registry.yaml
```

### 3.4 do 페르소나 생성

moai persona가 클린 템플릿이다. moai persona를 통째로 복제한 뒤 do 브랜드로 치환하고 do 방법론으로 내용을 수정한다. do-focus에서 복사하는 것이 아니다.

**3단계 프로세스**:

```
1. Extract moai-adk → core/ + personas/moai/    (3.2에서 완료)
2. Copy personas/moai/ → personas/do/            (moai persona를 통째로 복제)
3. Modify personas/do/                           (moai→do 브랜드 치환 + do 방법론으로 내용 수정)
```

**Step 2: moai persona 복제**

```bash
# moai persona를 do persona로 통째로 복사
cp -r ./extracted-moai-v3.0.0/personas/moai ./personas/do

# 디렉토리 구조 내 moai/ 서브폴더를 do/로 리네임
mv ./personas/do/agents/moai ./personas/do/agents/do
mv ./personas/do/skills/moai ./personas/do/skills/do
mv ./personas/do/commands/moai ./personas/do/commands/do
mv ./personas/do/hooks/moai ./personas/do/hooks/do         # 이후 삭제 대상
mv ./personas/do/output-styles/moai ./personas/do/output-styles/do
mv ./personas/do/rules/moai ./personas/do/rules/do
```

**Step 3: moai→do 변환**

변환 치환표에 따라 파일 내용을 수정한다.

| moai 패턴 | do 치환 | 비고 |
|-----------|--------|------|
| `agents/moai/` | `agents/do/` | 경로 |
| `.claude/hooks/moai/*.sh` | godo hook 직접 호출 | 구조 변경 (shell wrapper 제거) |
| `.moai/` | `.do/` | 디렉토리 |
| `/moai` | `/do` | 커맨드 접두사 |
| `moai-` (skill prefix) | `do-` | 스킬명 |
| `MoAI` | `Do` | 브랜드명 |
| `SPEC-XXX` | checklist 기반 | 워크플로우 |
| `YYMMDD` | `YY/MM/DD` | 날짜 형식 |

파일별 변환 내용:

1. `manifest.yaml` -- 브랜드, 경로 치환 + hook_scripts를 빈 배열로
2. `CLAUDE.md` -- 브랜드 치환 + do 3모드(Do/Focus/Team) 방법론 반영
3. `settings.json` -- godo hook 직접 호출로 변환, outputStyle 변경
4. `skills/do/SKILL.md` -- Intent Router/Mode Router를 do 방식으로 재작성
5. `skills/do/workflows/*.md` -- moai SPEC 워크플로우를 do 체크리스트 워크플로우로 변환
6. `skills/do/references/reference.md` -- 브랜드 치환
7. `output-styles/do/*.md` -- 스타일명 매핑 (moai→pair, r2d2→sprint, yoda→direct)
8. `commands/do/*.md` -- do 고유 커맨드로 교체 (moai 2개 → do 6개)
9. `rules/do/workflow/*.md` -- 브랜드 치환 + 워크플로우 단계명 변경
10. `agents/do/*.md` -- hooks 경로를 godo hook으로, memory/skills 참조 치환
11. `hooks/do/` -- 삭제 (godo binary가 직접 처리하므로 shell wrapper 불필요)

### 3.5 Assemble 검증

```bash
# core + do persona로 assemble
./convert assemble \
  --core ./extracted-moai-v3.0.0/core \
  --persona ./personas/do/manifest.yaml \
  --out /tmp/assembled-do

# 기대하는 .claude/ 구조와 비교
diff -rq /tmp/assembled-do ~/Work/do-focus.workspace/do-focus/.claude

# 미치환 슬롯 확인 (0개여야 함)
grep -r "{{slot:" /tmp/assembled-do || echo "OK: 미치환 슬롯 없음"

# HARD 규칙 무결성 확인
grep -r "\[HARD\]" /tmp/assembled-do | wc -l

# moai 참조 잔존 확인 (do persona에 moai 문자열이 없어야 함)
grep -ri "moai" /tmp/assembled-do/CLAUDE.md && echo "WARNING: moai 참조 잔존" || echo "OK"
```

### 3.6 Roundtrip 검증 (moai 측)

moai persona로도 assemble이 올바른지 확인한다.

```bash
# core + moai persona로 assemble
./convert assemble \
  --core ./extracted-moai-v3.0.0/core \
  --persona ./extracted-moai-v3.0.0/personas/moai/manifest.yaml \
  --out /tmp/assembled-moai

# 원본과 비교 (roundtrip 일치)
diff -rq /tmp/assembled-moai ~/Work/moai-adk/.claude
```

---

## 4. 버전업 프로세스 (반복)

이 섹션이 가장 핵심이다. moai-adk가 버전업될 때마다 이 프로세스를 따른다.

### 4.1 사전 분석: 무엇이 변했는지 파악

```bash
# 변수 설정
OLD_VER=v3.0.0
NEW_VER=v3.1.0

# moai-adk에서 변경사항 확인
cd ~/Work/moai-adk
git fetch --tags

# 변경 파일 통계
git diff ${OLD_VER}..${NEW_VER} --stat -- .claude/

# 변경 파일 목록 (핵심)
git diff ${OLD_VER}..${NEW_VER} --name-only -- .claude/

# 커밋 히스토리 (변경 의도 파악)
git log ${OLD_VER}..${NEW_VER} --oneline -- .claude/
```

### 4.2 변경 영향 분류

변경된 각 파일을 아래 4가지 카테고리로 분류한다. **분류의 정확성이 전체 프로세스의 품질을 결정한다.**

| 카테고리 | 설명 | do 페르소나 영향 | 조치 |
|---------|------|----------------|------|
| **A. Core 변경** | 브랜드 무관 파일 변경 (core agents, 공유 rules, 공유 skills) | 자동 반영 (re-extract하면 됨) | Extract만 재실행 |
| **B. Persona 구조 변경** | manifest 필드 추가, 새 파일 유형, 슬롯 스키마 변경 | do manifest + converter 코드 수정 필요 | manifest.yaml 업데이트 + converter 코드 수정 |
| **C. Persona 내용 변경** | moai CLAUDE.md, SKILL.md, workflow 파일 등의 내용 수정 | do 대응 파일 검토 필요 | 변경 의도 파악 후 do 버전 반영 |
| **D. 신규 기능 추가** | 새 workflow, 새 에이전트 유형, 새 커맨드 | do에 대응 파일 추가 여부 판단 | 분석 후 필요하면 추가 |

**분류 판단 기준**:

```
파일이 변경됐을 때:
├── .claude/agents/moai/에 없는 agent? → Core 변경 (A)
├── .claude/skills/moai*가 아닌 skill? → Core 변경 (A)
├── .claude/rules/moai/가 아닌 rule? → Core 변경 (A)
├── manifest 스키마 자체가 변경? → Persona 구조 변경 (B)
├── 새 slot 패턴 추가? → Persona 구조 변경 (B)
├── .claude/skills/moai/SKILL.md 내용 변경? → Persona 내용 변경 (C)
├── .claude/CLAUDE.md 변경? → Persona 내용 변경 (C)
├── .claude/skills/moai/workflows/* 변경? → Persona 내용 변경 (C)
├── 새 workflow 파일 추가? → 신규 기능 추가 (D)
└── 새 persona agent 추가? → 신규 기능 추가 (D)
```

### 4.3 분해 실행

```bash
cd ~/Work/new/convert

# Step 1: 새 버전 체크아웃 (또는 태그 확인)
cd ~/Work/moai-adk
git checkout ${NEW_VER}

# Step 2: 새 버전으로 Extract
cd ~/Work/new/convert
./convert extract \
  --src ~/Work/moai-adk/.claude \
  --out ./extracted-moai-${NEW_VER}

# Step 3: core 변경 비교
diff -rq ./extracted-moai-${OLD_VER}/core ./extracted-moai-${NEW_VER}/core

# Step 4: persona 변경 비교
diff -rq \
  ./extracted-moai-${OLD_VER}/personas/moai \
  ./extracted-moai-${NEW_VER}/personas/moai

# Step 5: 슬롯 레지스트리 변경 비교
diff ./extracted-moai-${OLD_VER}/core/registry.yaml \
     ./extracted-moai-${NEW_VER}/core/registry.yaml
```

### 4.4 카테고리별 do 페르소나 보완

#### A. Core 변경 (자동 반영)

core 파일만 변경된 경우 가장 간단하다.

```bash
# 새 core를 채택
cp -r ./extracted-moai-${NEW_VER}/core ./core-${NEW_VER}

# assemble로 검증
./convert assemble \
  --core ./core-${NEW_VER} \
  --persona ./personas/do/manifest.yaml \
  --out /tmp/verify-core-upgrade

# 미치환 슬롯 확인
grep -r "{{slot:" /tmp/verify-core-upgrade || echo "OK"
```

**주의**: 새 core 파일이 새 슬롯을 도입했을 수 있다. `registry.yaml` diff를 반드시 확인할 것.

#### B. Persona 구조 변경 (converter + manifest 수정)

manifest 스키마 또는 슬롯 패턴이 변경된 경우.

**절차**:
1. moai manifest의 새 필드/구조 확인
2. converter Go 코드가 새 구조를 지원하는지 확인
3. converter 코드 수정 (필요시)
4. do manifest.yaml에 새 필드 추가
5. assemble 재검증

```bash
# moai manifest 변경 확인
diff ./extracted-moai-${OLD_VER}/personas/moai/manifest.yaml \
     ./extracted-moai-${NEW_VER}/personas/moai/manifest.yaml

# 새 필드가 있다면 do manifest에도 추가
# (에디터로 personas/do/manifest.yaml 수정)

# converter 코드 수정 후 리빌드
go build -o convert ./cmd/convert

# 검증
./convert assemble \
  --core ./extracted-moai-${NEW_VER}/core \
  --persona ./personas/do/manifest.yaml \
  --out /tmp/verify-schema-upgrade
```

#### C. Persona 내용 변경 (수동 검토)

moai의 CLAUDE.md, SKILL.md, workflow 등의 내용이 변경된 경우. **가장 주의가 필요한 카테고리**.

**절차**:
1. 변경된 파일의 diff 확인
2. 커밋 메시지/PR에서 변경 의도 파악
3. Section 5 파일 매핑표에서 do 대응 파일 찾기
4. 변경 의도를 do 방법론에 맞게 반영
5. 검증

```bash
# moai SKILL.md 변경 확인 (예시)
diff ./extracted-moai-${OLD_VER}/personas/moai/skills/moai/SKILL.md \
     ./extracted-moai-${NEW_VER}/personas/moai/skills/moai/SKILL.md

# 커밋 메시지에서 의도 확인
cd ~/Work/moai-adk
git log ${OLD_VER}..${NEW_VER} --oneline -- .claude/skills/moai/SKILL.md
git show <commit-hash>  # 상세 확인

# do 대응 파일 수정 (예: skills/do/SKILL.md)
# → 매핑표 참조하여 어떤 do 파일을 수정할지 결정
```

**반영 원칙**:
- moai의 구조적 변경 (섹션 추가/삭제, 순서 변경) → do에도 동일 구조 반영
- moai의 HARD 규칙 추가/삭제 → do에 **무조건** 반영
- moai의 방법론 변경 (SPEC 관련) → do의 대응 방법론(체크리스트 기반)으로 번역
- moai의 브랜딩 변경 (TRUST 5, TAG Chain 등) → do에서는 무시 (rules에 이미 흡수)

#### D. 신규 기능 추가 (판단 필요)

새 파일이 추가된 경우.

**판단 매트릭스**:
| 추가된 파일 유형 | core에 해당? | do에 필요? | 조치 |
|----------------|-------------|-----------|------|
| 새 core agent | 예 | 자동 포함 | 조치 없음 |
| 새 core skill | 예 | 자동 포함 | 조치 없음 |
| 새 core rule | 예 | 자동 포함 | 조치 없음 |
| 새 persona agent | 아니오 | 판단 필요 | do 에이전트 작성 or 무시 |
| 새 persona workflow | 아니오 | 판단 필요 | do 워크플로우 작성 or 무시 |
| 새 persona command | 아니오 | 판단 필요 | do 커맨드 작성 or 무시 |
| 새 override skill | 아니오 | 보통 불필요 | do는 rules로 대체 |

**판단 기준**:
- do의 3모드(Do/Focus/Team) 체계에서 해당 기능이 필요한가?
- 기존 do 기능으로 대체 가능한가?
- 구현 비용 대비 가치가 있는가?

### 4.5 검증 체크리스트

모든 변경 작업 후 아래 항목을 **반드시** 확인한다.

```bash
OUT=/tmp/verify-upgrade

# 1. assemble 성공
./convert assemble \
  --core ./extracted-moai-${NEW_VER}/core \
  --persona ./personas/do/manifest.yaml \
  --out ${OUT}
echo "=== Assemble: $? (0이면 성공) ==="

# 2. 미치환 슬롯 0개
SLOTS=$(grep -r "{{slot:" ${OUT} 2>/dev/null | wc -l)
echo "=== 미치환 슬롯: ${SLOTS}개 (0이어야 함) ==="

# 3. HARD 규칙 무결성
HARD_COUNT=$(grep -r "\[HARD\]" ${OUT} 2>/dev/null | wc -l)
echo "=== HARD 규칙: ${HARD_COUNT}개 ==="

# 4. moai 참조 잔존 확인 (do persona 영역)
MOAI_REFS=$(grep -ri "moai" ${OUT}/CLAUDE.md ${OUT}/skills/do/ ${OUT}/commands/do/ 2>/dev/null | grep -v "moai-constitution" | wc -l)
echo "=== moai 잔존 참조: ${MOAI_REFS}개 (0이어야 함) ==="

# 5. 브랜드 슬롯이 올바르게 치환됐는지
grep -r "{{slot:BRAND" ${OUT} && echo "FAIL: 브랜드 슬롯 미치환" || echo "OK: 브랜드 슬롯 치환 완료"

# 6. 파일 구조 비교 (이전 assemble 결과와)
diff -rq ${OUT} ~/Work/do-focus.workspace/do-focus/.claude 2>/dev/null | head -20

# 7. roundtrip 검증 (moai 측)
./convert assemble \
  --core ./extracted-moai-${NEW_VER}/core \
  --persona ./extracted-moai-${NEW_VER}/personas/moai/manifest.yaml \
  --out /tmp/verify-moai-roundtrip
diff -rq /tmp/verify-moai-roundtrip ~/Work/moai-adk/.claude | head -20
```

### 4.6 버전 기록 업데이트

모든 검증 통과 후 `versions.yaml`을 업데이트한다.

```yaml
# versions.yaml 업데이트
moai:
  repo: ~/Work/moai-adk
  version: v3.1.0              # 새 버전
  extracted_at: 2026-03-01     # 오늘 날짜
  core_files: 401              # 새 core 파일 수
  persona_files: 63            # 새 persona 파일 수
  notes: "workflow 2개 추가, agent 1개 추가"

do:
  based_on_moai: v3.1.0        # 업데이트된 기반 버전
  version: "3.1.0"
  updated_at: 2026-03-01
  notes: "core 업데이트 반영, 새 workflow 미채택"

history:
  - date: 2026-03-01
    moai_version: v3.1.0
    action: upgrade
    changes: |
      - core: agent 1개 추가 (expert-ml), skill 모듈 5개 업데이트
      - persona 내용: SKILL.md Intent Router에 ml 키워드 추가
      - do 반영: core 자동 반영, SKILL.md 변경은 do에 해당 없음
  - date: 2026-02-15
    moai_version: v3.0.0
    action: initial_extract
    changes: "초기 분해 및 do persona 생성"
```

---

## 5. 파일 매핑표 (moai <-> do)

### 5.1 Persona 에이전트 분류

moai-adk의 28개 에이전트 중 **core**(22개)와 **persona**(6개)로 분류된다.

**Core 에이전트 (22개)** -- 양쪽 모두 자동 포함:

| 카테고리 | 에이전트 |
|---------|---------|
| Builder (3) | builder-agent, builder-plugin, builder-skill |
| Expert (9) | expert-backend, expert-chrome-extension, expert-debug, expert-devops, expert-frontend, expert-performance, expert-refactoring, expert-security, expert-testing |
| Manager (3) | manager-docs, manager-git, manager-strategy |
| Team (7) | team-analyst, team-architect, team-backend-dev, team-designer, team-frontend-dev, team-researcher, team-tester |

**Persona 에이전트**:

| moai persona 에이전트 | do persona 에이전트 | 비고 |
|---------------------|-------------------|------|
| manager-ddd.md | manager-ddd.md | DDD 구현 매니저 (양쪽 동일 역할) |
| manager-project.md | manager-project.md | 프로젝트 설정 |
| manager-quality.md | manager-quality.md | 품질 검증 |
| manager-spec.md | manager-spec.md (제거 후보) | moai: SPEC 문서, do: plan 워크플로우로 대체 가능 |
| manager-tdd.md | manager-tdd.md | TDD 구현 매니저 |
| team-quality.md | team-quality.md | Team 모드 품질 검증 |

### 5.2 Persona 스킬 매핑

#### 오케스트레이터 스킬

| moai persona 파일 | do persona 파일 | 관계 |
|------------------|----------------|------|
| `skills/moai/SKILL.md` | `skills/do/SKILL.md` | 구조 동일, Intent Router/Mode Router 내용 다름 |
| `skills/moai/workflows/plan.md` | `skills/do/workflows/plan.md` | moai: SPEC 생성, do: Analysis->Architecture->Plan |
| `skills/moai/workflows/run.md` | `skills/do/workflows/run.md` | moai: SPEC 기반, do: Checklist 기반 |
| `skills/moai/workflows/moai.md` | `skills/do/workflows/do.md` | 자동 파이프라인 (이름만 다름) |
| `skills/moai/workflows/sync.md` | (없음) | moai: 문서 동기화, do: report.md로 대체 |
| `skills/moai/workflows/feedback.md` | (없음) | moai 고유 |
| `skills/moai/workflows/fix.md` | (없음) | moai 고유 |
| `skills/moai/workflows/loop.md` | (없음) | moai 고유 |
| `skills/moai/workflows/project.md` | (없음) | godo setup으로 대체 |
| `skills/moai/workflows/team-plan.md` | `skills/do/workflows/team-plan.md` | 구조 유사 |
| `skills/moai/workflows/team-run.md` | `skills/do/workflows/team-run.md` | 구조 유사 |
| `skills/moai/workflows/team-sync.md` | (없음) | moai 고유 |
| `skills/moai/workflows/team-debug.md` | (없음) | moai 고유 |
| `skills/moai/workflows/team-review.md` | (없음) | moai 고유 |
| `skills/moai/references/reference.md` | `skills/do/references/reference.md` | 공통 패턴 참조 |
| `skills/moai/moai-workflow-team/SKILL.md` | (없음) | do는 team workflow를 workflows/에 통합 |

#### Override 스킬 (moai -> do 대응)

| moai override 스킬 | do 대응 | 비고 |
|-------------------|--------|------|
| `moai-foundation-core/SKILL.md` | 없음 | `rules/*.md`에 이미 흡수 |
| `moai-foundation-quality/SKILL.md` | 없음 | `dev-testing.md` + `dev-workflow.md` |
| `moai-workflow-ddd/SKILL.md` | 없음 | `rules/do/workflow/workflow-modes.md` |
| `moai-workflow-tdd/SKILL.md` | 없음 | `dev-workflow.md` TDD 섹션 |
| `moai-workflow-spec/SKILL.md` | 없음 | `.do/jobs/` + `dev-checklist.md` |
| `moai-workflow-project/SKILL.md` | 없음 | `godo init/setup` 대체 |

**핵심 인사이트**: moai는 "skill -> progressive disclosure"로 지식을 로드하지만, do는 "rules -> 항상 로드"로 직접 주입한다. 따라서 override skill 변환이 아니라 rules 유지가 올바른 전략이다. `agent_patches`도 비워둔다.

### 5.3 기타 Persona 파일 매핑

#### CLAUDE.md

| moai | do | 비고 |
|------|-----|------|
| `CLAUDE.md` (~358줄, lean) | `CLAUDE.md` (~200줄, lean) | 구조 유사, 내용(방법론) 다름 |

#### Output Styles

| moai 스타일 | do 스타일 | 대응 |
|------------|----------|------|
| `output-styles/moai/moai.md` | `output-styles/do/pair.md` | 기본 스타일 (협업적 톤) |
| `output-styles/moai/r2d2.md` | `output-styles/do/sprint.md` | 최소 대화 (빠른 실행) |
| `output-styles/moai/yoda.md` | `output-styles/do/direct.md` | 간결한 전문성 |

#### Commands

| moai 커맨드 | do 커맨드 | 비고 |
|------------|----------|------|
| `commands/moai/github.md` | (없음) | moai 프로젝트 특화 |
| `commands/moai/99-release.md` | (없음) | moai 프로젝트 특화 |
| (없음) | `commands/do/check.md` | do 고유 |
| (없음) | `commands/do/checklist.md` | do 고유 |
| (없음) | `commands/do/mode.md` | do 고유 (3모드 전환) |
| (없음) | `commands/do/plan.md` | do 고유 |
| (없음) | `commands/do/setup.md` | do 고유 |
| (없음) | `commands/do/style.md` | do 고유 |

#### Hooks

| moai 훅 | do 대응 | 비고 |
|--------|--------|------|
| `hooks/moai/handle-session-start.sh` | (없음) | `godo hook session-start` 직접 호출 |
| `hooks/moai/handle-session-end.sh` | (없음) | `godo hook session-end` 직접 호출 |
| `hooks/moai/handle-pre-tool.sh` | (없음) | `godo hook pre-tool` 직접 호출 |
| `hooks/moai/handle-post-tool.sh` | (없음) | `godo hook post-tool-use` 직접 호출 |
| `hooks/moai/handle-compact.sh` | (없음) | `godo hook compact` 직접 호출 |
| `hooks/moai/handle-stop.sh` | (없음) | `godo hook stop` 직접 호출 |
| `hooks/moai/handle-agent-hook.sh` | (없음) | `godo hook subagent-stop` 직접 호출 |

**구조적 차이**: moai는 shell wrapper 스크립트를 통해 moai CLI를 호출하지만, do는 godo binary를 settings.json에서 직접 호출한다. 따라서 `hook_scripts: []` (빈 배열).

#### Rules

| moai persona 규칙 | do persona 규칙 | 비고 |
|------------------|----------------|------|
| `rules/moai/workflow/spec-workflow.md` | `rules/do/workflow/spec-workflow.md` | 구조 동일, 단계명 다름 (Plan/Run/Sync vs Plan/Run/Report) |
| `rules/moai/workflow/workflow-modes.md` | `rules/do/workflow/workflow-modes.md` | DDD/TDD/Hybrid 방법론 (동일) |

#### Slot Content

| moai slot | do 대응 | 비고 |
|-----------|--------|------|
| `QUALITY_FRAMEWORK` (TRUST 5) | `slot_content`에 유지 또는 제거 | do는 `dev-testing.md` + `dev-workflow.md`로 대체 |
| `QUALITY_GATE_TEXT` | `slot_content`에 유지 또는 제거 | 텍스트 참조 |
| `TRACEABILITY_SYSTEM` (TAG Chain) | 제거 | do는 TAG Chain 미사용 |

### 5.4 Settings.json 구조 비교

| 필드 | moai | do | 비고 |
|------|------|-----|------|
| `outputStyle` | moai | pair | 기본 스타일 |
| `plansDirectory` | (moai 기본) | `.do/jobs` | do 고유 경로 |
| `statusLine` | moai CLI | `godo statusline` | CLI 도구 차이 |
| `hooks.SessionStart` | shell wrapper | `godo hook session-start` | 직접 호출 |
| `hooks.PostToolUse` | shell wrapper | `godo hook post-tool-use` | matcher: `.*` |
| `hooks.PreToolUse` | shell wrapper | `godo hook pre-tool` | matcher: `Write\|Edit\|Bash` |
| `hooks.UserPromptSubmit` | (없음) | `godo hook user-prompt-submit` | do 고유 |
| `hooks.SubagentStop` | (없음) | `godo hook subagent-stop` | do 고유 |

---

## 6. 중점 검토 항목

버전업 시 특히 주의해야 할 변경 유형별 가이드.

### 6.1 CLAUDE.md 변경

**영향도: HIGH**

moai CLAUDE.md가 변경되면 do CLAUDE.md도 검토해야 한다.

| 변경 유형 | do 대응 | 우선순위 |
|----------|--------|---------|
| 새 섹션 추가 | do에도 대응 섹션 추가 여부 판단 | HIGH |
| 섹션 구조 변경 | do에도 동일 구조로 반영 | HIGH |
| HARD 규칙 추가/삭제 | **무조건 반영** | CRITICAL |
| 참조 파일 경로 변경 | do의 참조 경로도 업데이트 | MEDIUM |
| Request Processing Pipeline 변경 | do의 파이프라인도 검토 | HIGH |

### 6.2 SKILL.md (오케스트레이터) 변경

**영향도: HIGH**

| 변경 유형 | do 대응 | 우선순위 |
|----------|--------|---------|
| Intent Router 변경 | do SKILL.md에도 동일하게 반영 | CRITICAL |
| 새 workflow 참조 추가 | do에 대응 workflow 추가 여부 판단 | HIGH |
| Execution Directive 변경 | do 8단계에 반영 | HIGH |
| Agent Catalog 변경 | do Agent Catalog도 업데이트 | MEDIUM |
| Core Rules 변경 | do Core Rules 검토 | HIGH |
| 새 Priority 레벨 추가 | do Intent Router에 반영 | HIGH |

### 6.3 Manifest 스키마 변경

**영향도: HIGH** (converter 코드 수정 필요할 수 있음)

| 변경 유형 | do 대응 | 우선순위 |
|----------|--------|---------|
| 새 필드 추가 | converter Go 코드 + do manifest 모두 업데이트 | CRITICAL |
| 기존 필드 변경 | 호환성 확인 + converter + manifest 업데이트 | CRITICAL |
| 필드 삭제 | converter에서 삭제 필드 처리 확인 | HIGH |
| 기본값 변경 | do manifest에서 명시적 값 설정 여부 확인 | MEDIUM |

### 6.4 Slot 변수 추가

**영향도: MEDIUM~HIGH**

```bash
# 새 슬롯 패턴 확인
diff ./extracted-moai-${OLD_VER}/core/registry.yaml \
     ./extracted-moai-${NEW_VER}/core/registry.yaml
```

| 변경 유형 | do 대응 | 우선순위 |
|----------|--------|---------|
| 새 `{{slot:*}}` 패턴 | converter의 slotifier/deslotifier에 추가 | HIGH |
| 새 `slot_content` 항목 | do manifest의 `slot_content`에 대응 값 추가 | HIGH |
| 기존 슬롯 삭제 | do manifest에서도 삭제 | MEDIUM |

### 6.5 에이전트/스킬 추가

**영향도: LOW~MEDIUM**

| 추가 유형 | do 대응 | 우선순위 |
|----------|--------|---------|
| 새 core 에이전트 | 자동 반영 (조치 없음) | LOW |
| 새 core 스킬 | 자동 반영 (조치 없음) | LOW |
| 새 persona 에이전트 | do에 대응 에이전트 추가 필요 여부 판단 | MEDIUM |
| 새 persona 스킬 | do에 대응 필요 여부 판단 | MEDIUM |
| 새 override 스킬 | 보통 불필요 (do는 rules로 대체) | LOW |

### 6.6 Workflow 추가/변경

**영향도: MEDIUM**

```bash
# 새 workflow 파일 확인
diff <(ls ./extracted-moai-${OLD_VER}/personas/moai/skills/moai/workflows/) \
     <(ls ./extracted-moai-${NEW_VER}/personas/moai/skills/moai/workflows/)
```

**판단 기준**: 새 workflow가 do의 3모드(Do/Focus/Team) 또는 체크리스트 체계에서 가치가 있는가?

### 6.7 Hooks 이벤트 추가

**영향도: MEDIUM**

moai가 새 hook 이벤트를 추가하면 do의 settings.json에도 대응 godo hook을 추가해야 할 수 있다.

```bash
# moai settings.json에서 hooks 변경 확인
diff <(grep -A2 "hooks" ./extracted-moai-${OLD_VER}/personas/moai/settings.json) \
     <(grep -A2 "hooks" ./extracted-moai-${NEW_VER}/personas/moai/settings.json)
```

---

## 7. 트러블슈팅

### 7.1 Assemble 실패

| 증상 | 원인 | 해결 |
|------|------|------|
| "manifest.yaml not found" | `--persona` 경로 오류 | 절대 경로 사용, 파일 존재 확인 |
| "core directory not found" | `--core` 경로 오류 | `ls <core-path>/registry.yaml` 확인 |
| "unknown slot" | registry에 없는 슬롯 사용 | `registry.yaml`에 슬롯 등록 확인 |
| 파일 누락 | manifest에 선언됐지만 실제 파일 없음 | manifest의 파일 경로와 실제 파일 매칭 |

### 7.2 미치환 슬롯 잔존

```bash
# 미치환 슬롯 찾기
grep -rn "{{slot:" /tmp/assembled-output/

# 원인: slot_content에 해당 값이 없음
# 해결: manifest.yaml의 slot_content에 값 추가
```

### 7.3 moai 참조 잔존

```bash
# do persona 영역에서 moai 문자열 검색
grep -ri "moai" /tmp/assembled-output/CLAUDE.md
grep -ri "moai" /tmp/assembled-output/skills/do/
grep -ri "moai" /tmp/assembled-output/commands/do/

# 예외: moai-constitution.md는 core 파일이므로 moai 참조 허용
# 예외: core agents/skills 내 moai 참조 (슬롯 치환 대상)
```

### 7.4 Extract 결과 불일치

roundtrip(extract -> assemble)이 원본과 다를 때:

```bash
# 상세 diff
diff -r /tmp/assembled-moai ~/Work/moai-adk/.claude | head -50

# 원인 1: classifier가 파일을 잘못 분류
# → converter의 detector 코드 확인

# 원인 2: 슬롯 치환이 원본과 다르게 동작
# → registry.yaml과 slot_content 교차 검증
```

---

## 8. 자동화 로드맵

### Phase 1: 수동 프로세스 (현재)

이 runbook을 따라 수동으로 분해/비교/보완한다.

- Extract/Assemble CLI는 이미 구현됨
- 비교/분류/보완은 수동 diff + 판단

### Phase 2: 반자동화

**`convert diff` 명령 추가**:
```bash
# 두 extracted 버전 간 변경사항을 자동 분류
convert diff \
  --old ./extracted-moai-v3.0.0 \
  --new ./extracted-moai-v3.1.0

# 출력 예:
# Core Changes (auto-reflect):
#   modified: rules/do/languages/python.md
#   added:    skills/do-lang-zig/SKILL.md
#
# Persona Content Changes (manual review):
#   modified: personas/moai/skills/moai/SKILL.md
#   modified: personas/moai/CLAUDE.md
#
# Persona Structure Changes (manifest update):
#   (none)
#
# New Features (decision needed):
#   added: personas/moai/skills/moai/workflows/review.md
```

**`convert check` 명령 추가**:
```bash
# do persona 무결성 검증
convert check \
  --core ./core \
  --persona ./personas/do/manifest.yaml

# 출력 예:
# [PASS] Assemble succeeds
# [PASS] No unresolved slots
# [PASS] No moai references in do persona files
# [PASS] All HARD rules preserved
# [WARN] slot_content.TRACEABILITY_SYSTEM is empty
```

### Phase 3: 완전 자동화

**`convert upgrade` 명령**:
```bash
# moai 버전업 → 분석 → 제안 → 적용
convert upgrade \
  --from v3.0.0 \
  --to v3.1.0 \
  --persona do

# 자동으로:
# 1. 두 버전 extract
# 2. diff 분석 및 분류
# 3. core 변경 자동 적용
# 4. persona 변경 보고서 생성
# 5. 검증 체크리스트 실행
```

**CI/CD 통합**:
```yaml
# .github/workflows/moai-upgrade.yml
on:
  push:
    tags: ['v*']
    # moai-adk 태그 push 시

jobs:
  upgrade:
    steps:
      - run: convert upgrade --from $OLD --to $NEW --persona do
      - run: convert check --core ./core --persona ./personas/do/manifest.yaml
      - run: gh pr create --title "moai ${NEW} upgrade" --body "..."
```

---

## 부록 A: 디렉토리 구조 전체 참조

### converter 프로젝트 구조

```
~/Work/new/convert/
├── cmd/convert/main.go              # CLI 진입점
├── internal/
│   ├── cli/                         # Cobra 커맨드 (extract, assemble)
│   ├── detector/                    # persona 감지 (패턴 매칭, 분류)
│   ├── extractor/                   # Extract 파이프라인
│   ├── assembler/                   # Assemble 파이프라인
│   ├── parser/                      # Markdown 파싱 (frontmatter, 섹션)
│   ├── template/                    # 슬롯 레지스트리, 슬롯 연산
│   └── model/                       # 공유 타입 (Document, PersonaManifest, Slot)
├── extracted-do/                    # do-focus에서 추출한 결과
│   ├── core/                        # core 파일들
│   └── personas/do/                 # do persona 파일들
├── personas/do/                     # do persona (작업 중)
├── versions.yaml                    # 버전 추적 파일
├── RUNBOOK.md                       # 이 문서
└── README.md                        # 프로젝트 설명
```

### Assemble 출력 구조 (.claude/)

```
.claude/
├── CLAUDE.md                        # persona CLAUDE.md
├── settings.json                    # persona settings (core base + persona overrides)
├── agents/do/                       # core agents (22개) + persona agents (5~6개)
├── rules/
│   ├── dev-*.md                     # core rules (공유)
│   ├── file-reading.md              # core
│   ├── go.md                        # core
│   └── do/                          # persona rules + core do/ rules
│       ├── core/                    # core (moai-constitution 등)
│       ├── development/             # core (agent-authoring 등)
│       ├── languages/               # core (언어별 규칙)
│       └── workflow/                # persona (spec-workflow, workflow-modes)
├── skills/
│   ├── do/                          # persona 오케스트레이터 스킬
│   │   ├── SKILL.md
│   │   ├── workflows/
│   │   └── references/
│   └── do-*/                        # core 스킬 (40+ 스킬 패키지)
├── commands/do/                     # persona 커맨드 (6개)
└── output-styles/do/                # persona 스타일 (3개)
```

---

## 부록 B: 빠른 참조 커맨드 모음

```bash
# === 버전 확인 ===
cd ~/Work/moai-adk && git describe --tags

# === Extract ===
cd ~/Work/new/convert
./convert extract --src ~/Work/moai-adk/.claude --out ./extracted-moai-vX.Y.Z

# === Assemble (do) ===
./convert assemble \
  --core ./extracted-moai-vX.Y.Z/core \
  --persona ./personas/do/manifest.yaml \
  --out /tmp/assembled-do

# === Assemble (moai roundtrip) ===
./convert assemble \
  --core ./extracted-moai-vX.Y.Z/core \
  --persona ./extracted-moai-vX.Y.Z/personas/moai/manifest.yaml \
  --out /tmp/assembled-moai

# === 검증 ===
grep -r "{{slot:" /tmp/assembled-do          # 미치환 슬롯 (0이어야 함)
grep -ri "moai" /tmp/assembled-do/CLAUDE.md  # moai 잔존 참조 (0이어야 함)
grep -r "\[HARD\]" /tmp/assembled-do | wc -l # HARD 규칙 수
diff -rq /tmp/assembled-do ~/Work/do-focus.workspace/do-focus/.claude  # 구조 비교

# === 버전 간 비교 ===
diff -rq ./extracted-moai-vOLD/core ./extracted-moai-vNEW/core
diff -rq ./extracted-moai-vOLD/personas/moai ./extracted-moai-vNEW/personas/moai
diff ./extracted-moai-vOLD/core/registry.yaml ./extracted-moai-vNEW/core/registry.yaml
```

---

## 9. Do 페르소나 정체성 검증

이 섹션은 moai 코드가 변경되었을 때 Do 페르소나의 고유 정체성이 보존되는지 검증하는 가이드이다.
상세 정체성 정의는 `.do/jobs/260215/do-persona-design/do-identity.md`를, 파일별 출처 추적은 `.do/jobs/260215/do-persona-design/conversion-manifest.md`를 참조한다.

### 9.1 moai 코드 변경 시 Do 정체성 체크리스트

moai-adk가 업데이트될 때마다 아래 항목을 순서대로 확인한다.

#### Step 1: 정체성 경계(Identity Boundaries) 침범 여부

| 정체성 요소 | 보호 수준 | 확인 방법 | 침범 시 조치 |
|------------|----------|----------|-------------|
| 삼원 실행 구조 (Do/Focus/Team) | CRITICAL | `grep "나는 Do다" personas/do/CLAUDE.md` | 절대 변경 금지. moai 변경 무시. |
| 체크리스트 시스템 (`[o][~][*][!][x]`) | CRITICAL | `grep "\[o\]" personas/do/skills/do/workflows/run.md` | 절대 변경 금지. |
| 페르소나 4종 캐릭터 | HIGH | `ls personas/do/characters/` (4개 파일) | moai에 캐릭터 추가돼도 Do 캐릭터는 독립. |
| 스타일 3종 (sprint/pair/direct) | HIGH | `ls personas/do/output-styles/do/` (3개 파일) | moai 스타일 변경과 무관. |
| godo 직접 호출 패턴 | HIGH | `grep "godo hook" personas/do/settings.json` | shell wrapper 도입 금지. |
| 한국어 선언문 | HIGH | `grep "말하면 한다" personas/do/CLAUDE.md` | 영어로 대체 금지. |
| 6개 개별 커맨드 | MEDIUM | `ls personas/do/commands/do/` (6개 파일) | `/moai` 통합 진입점으로 전환 금지. |
| Jobs 디렉토리 경로 | MEDIUM | `grep "plansDirectory" personas/do/settings.json` | `.moai/specs/` 패턴 도입 금지. |
| DO_* 환경변수 체계 | MEDIUM | `grep "DO_MODE" personas/do/` | `.moai/config/*.yaml` 방식으로 전환 금지. |

#### Step 2: moai 변경이 Do에 미치는 영향 분류

변경된 각 파일을 아래 기준으로 분류한다:

```
moai 파일이 변경됐을 때:
├── Do 정체성 경계를 침범하는가?
│   ├── YES → 변경 거부 (Do 정체성 보존)
│   └── NO → 계속
├── core 파일인가? (agents/moai/ 밖, skills/moai* 밖)
│   ├── YES → 자동 반영 (re-extract)
│   └── NO → persona 파일 → 아래 분류 계속
├── HARD 규칙 추가/삭제인가?
│   ├── YES → Do에 무조건 반영 (방법론에 맞게 번역)
│   └── NO → 계속
├── 구조적 변경인가? (섹션 추가/삭제/재구성)
│   ├── YES → Do 대응 파일 수동 검토
│   └── NO → 기계적 치환 가능 여부 확인
```

#### Step 3: 변환 후 자동 검증 실행

```bash
# 정체성 보존 자동 검증 (9.3 스크립트 참조)
cd ~/Work/new/convert
bash .do/jobs/260215/do-persona-design/verify-identity.sh
# 또는 conversion-manifest.md의 Automated Verification Commands 참조
```

### 9.2 파일별 변환 유형

| 변환 유형 | 설명 | 파일 수 | 자동화 | 예시 |
|----------|------|---------|--------|------|
| **mechanical** | find/replace만으로 충분 (브랜드명, 경로) | 3 | YES | workflow-modes.md |
| **structural** | 개념/구조가 다름, 재작성 필요 | 12 | NO | CLAUDE.md, SKILL.md, workflows |
| **original** | moai에 없는 Do 고유 파일 | 12 | NO | commands, characters |
| **removed** | moai에 있지만 Do에 불필요 | 16 | YES | hooks, sync/fix/loop workflows |
| **absorbed** | moai skill이 Do rules에 통합됨 | 7 | N/A | 6 override skills |

**기계적 치환 가능한 패턴**:

| 패턴 | 치환 |
|------|------|
| `agents/moai/` | `agents/do/` |
| `.moai/` | `.do/` |
| `moai-` (skill prefix) | `do-` |
| `skills/moai/` | `skills/do/` |
| `moai hook` | `godo hook` |
| `moai statusline` | `godo statusline` |

**수동 검토 필요한 파일** (structural 변환):

| 파일 | 검토 포인트 |
|------|-----------|
| `CLAUDE.md` | 삼원 구조, 한국어 혼합, 체크리스트 워크플로우 보존 |
| `SKILL.md` | Mode Router, Intent Router, Execution Directive |
| `manifest.yaml` | hook_scripts 빈 배열, slot_content/agent_patches 비어있음 |
| `settings.json` | godo 직접 호출, SubagentStop, UserPromptSubmit |
| `workflows/plan.md` | Analysis->Architecture->Plan (SPEC 아님) |
| `workflows/run.md` | Checklist 기반 (SPEC 아님) |
| `output-styles/do/*.md` | sprint/pair/direct (moai/r2d2/yoda 아님) |

### 9.3 Identity Boundary 위반 감지

#### 자동 grep 패턴

다음 명령어로 정체성 위반을 감지한다:

```bash
PERSONA=./personas/do

# 1. moai 브랜드 잔존 (persona 파일에 moai 문자열)
echo "=== moai 브랜드 잔존 ==="
grep -ri "moai" ${PERSONA}/CLAUDE.md ${PERSONA}/skills/do/ \
  ${PERSONA}/commands/do/ ${PERSONA}/output-styles/do/ 2>/dev/null | \
  grep -v "moai-constitution" | grep -v "^Binary"
# 결과: 0줄이어야 함

# 2. SPEC 패턴 잔존 (Do는 SPEC 미사용)
echo "=== SPEC 패턴 잔존 ==="
grep -ri "SPEC-[0-9]" ${PERSONA}/ 2>/dev/null
grep -ri "EARS" ${PERSONA}/CLAUDE.md ${PERSONA}/skills/ 2>/dev/null
# 결과: 0줄이어야 함

# 3. shell wrapper 잔존
echo "=== Shell wrapper 잔존 ==="
grep -r "handle-.*\.sh" ${PERSONA}/ 2>/dev/null
grep -r "hooks/moai/" ${PERSONA}/ 2>/dev/null
# 결과: 0줄이어야 함

# 4. .moai/ 경로 잔존
echo "=== .moai/ 경로 잔존 ==="
grep -r "\.moai/" ${PERSONA}/ 2>/dev/null
# 결과: 0줄이어야 함

# 5. moai CLI 참조 (godo가 아닌)
echo "=== moai CLI 잔존 ==="
grep -r '"moai ' ${PERSONA}/ 2>/dev/null
grep -r "moai hook" ${PERSONA}/ 2>/dev/null
# 결과: 0줄이어야 함

# 6. 삼원 구조 건재 확인
echo "=== 삼원 구조 건재 ==="
grep -c "나는 Do다" ${PERSONA}/CLAUDE.md    # 1이어야 함
grep -c "나는 Focus다" ${PERSONA}/CLAUDE.md # 1이어야 함
grep -c "나는 Team이다" ${PERSONA}/CLAUDE.md # 1이어야 함

# 7. 페르소나 호칭 건재
echo "=== 페르소나 호칭 건재 ==="
grep -c "선배" ${PERSONA}/characters/young-f.md 2>/dev/null  # 1+ 이어야 함
grep -c "선배님" ${PERSONA}/characters/young-m.md 2>/dev/null # 1+ 이어야 함

# 8. Completion marker 교체 확인 (Do는 XML 마커 미사용)
echo "=== XML 마커 잔존 ==="
grep -r "<moai>" ${PERSONA}/ 2>/dev/null
# 결과: 0줄이어야 함
```

#### 수동 체크 항목

| # | 항목 | 확인 방법 | 기대값 |
|---|------|----------|--------|
| 1 | CLAUDE.md에 삼원 선언문 3개 | 육안 확인 | 3개 모두 존재 |
| 2 | CLAUDE.md가 한국어+영어 혼합 | 육안 확인 | 한국어 섹션 다수 |
| 3 | settings.json에 godo 직접 호출 | `cat settings.json \| grep godo` | 7개 hook 모두 godo |
| 4 | manifest.yaml에 hook_scripts 빈 배열 | `grep hook_scripts manifest.yaml` | `hook_scripts: []` |
| 5 | 스타일이 sprint/pair/direct | `ls output-styles/do/` | 3개 파일 |
| 6 | 커맨드가 /do:* 개별 방식 | `ls commands/do/` | 6개 파일 |
| 7 | Jobs 경로가 `.do/jobs/` | `grep plansDirectory settings.json` | `.do/jobs` |
| 8 | 페르소나 캐릭터 4종 | `ls characters/` | 4개 파일 |

### 9.4 변환 워크플로우

moai 업데이트 시 Do 페르소나를 동기화하는 전체 워크플로우.

```
Step 1: 변경 감지
├── moai-adk git diff 확인 (Section 4.1)
├── 변경 파일을 4가지 카테고리로 분류 (Section 4.2)
└── 정체성 경계 침범 여부 확인 (Section 9.1 Step 1)

Step 2: Extract
├── 새 버전으로 extract 실행 (Section 3.2)
├── core diff 확인 (자동 반영 대상)
└── persona diff 확인 (수동 검토 대상)

Step 3: Convert
├── core 변경 → 자동 반영 (re-extract만으로 충분)
├── persona mechanical 변환 → 치환표 적용 (Section 9.2)
├── persona structural 변환 → 수동 검토 + 재작성
│   ├── do-identity.md 참조하여 정체성 보존 확인
│   ├── conversion-manifest.md 참조하여 출처/유형 확인
│   └── 변경 의도를 Do 방법론으로 "번역"
├── 신규 기능 → 판단 (Section 4.4 D)
└── HARD 규칙 변경 → Do에 무조건 반영

Step 4: Verify
├── 자동 검증 스크립트 실행 (Section 9.3)
├── assemble 검증 (Section 3.5)
│   ├── 미치환 슬롯 0개
│   ├── moai 참조 잔존 0개
│   ├── HARD 규칙 수 확인
│   └── 구조 비교 (이전 결과와)
├── 정체성 경계 건재 확인 (삼원 구조, 체크리스트, 페르소나, godo)
└── roundtrip 검증 (moai 측)

Step 5: Assemble & Deploy
├── assemble 실행
├── do-focus 프로젝트에 결과 복사/비교
├── versions.yaml 업데이트 (Section 4.6)
└── 변경 이력 기록
```

#### 핵심 원칙

1. **정체성 보존 우선**: moai의 기능 추가보다 Do의 정체성 보존이 우선이다.
2. **기계적 치환 최대화**: 자동화 가능한 부분은 자동화하여 실수를 줄인다.
3. **구조적 변환은 신중하게**: structural 파일은 do-identity.md의 Architecture Decisions를 참조하여 Do 철학에 맞게 변환한다.
4. **검증은 자동으로**: 모든 변환 후 자동 검증 스크립트를 실행한다.
5. **HARD 규칙은 무조건 반영**: moai의 HARD 규칙 변경은 Do에도 반영하되, Do 방법론으로 번역한다.
