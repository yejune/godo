---
description: 구현 계획 작성 - 코드베이스 분석 후 플랜 생성
allowed-tools: AskUserQuestion, Bash, Read, Write, Glob, Grep
argument-hint: [플랜 제목 또는 빈칸]
---

# /do:plan - 구현 계획 작성

## 실행 단계

### Step 1: 플랜 제목 결정

- `$ARGUMENTS`가 있으면 해당 제목 사용
- 없으면 AskUserQuestion으로 질문:
  - question: "플랜 제목을 입력하세요"
  - header: "Plan Title"

제목 형식화:
- 공백은 하이픈(-)으로 변환
- 소문자로 통일
- 특수문자 제거

### Step 2: 코드베이스 분석

1. **프로젝트 구조 파악**
```bash
# 프로젝트 루트 파일 확인
ls -la

# 주요 디렉토리 구조
find . -type d -maxdepth 3 -not -path '*/\.*' -not -path '*/node_modules/*' | head -30
```

2. **기술 스택 식별**
- package.json, pyproject.toml, go.mod, Cargo.toml 등 확인
- 사용 중인 프레임워크 및 라이브러리 파악

3. **기존 코드 패턴 분석**
- 디렉토리 구조 및 네이밍 컨벤션
- 테스트 구조 확인

4. **관련 파일 탐색**
- 제목과 관련된 기존 코드 검색
- 유사한 구현 사례 확인

### Step 3: 플랜 작성

아래 형식으로 구현 계획 작성:

```markdown
# {제목}

생성일: {YYYY-MM-DD}
상태: draft

## 목표
{사용자 요청 또는 플랜 제목에서 추론한 목표 요약}

## 분석

### 현재 상태
{코드베이스 분석 결과}
- 기술 스택: {식별된 스택}
- 관련 기존 코드: {있으면 경로 나열}

### 고려사항
{구현 시 고려해야 할 사항}

## 구현 계획

### Phase 1: {단계명}
- [ ] {작업 1}
- [ ] {작업 2}

### Phase 2: {단계명}
- [ ] {작업 3}
- [ ] {작업 4}

## 파일 변경 예상
- `{파일1}`: {변경 내용}
- `{파일2}`: {변경 내용}

## 테스트 계획
- {테스트 항목 1}
- {테스트 항목 2}

## 리스크 및 대안
- {리스크}: {대안}
```

### Step 4: 디렉토리 생성 및 저장

```bash
mkdir -p .do/jobs/{YY}/{MM}/{DD}/{제목-kebab-case}
```

파일 경로: `.do/jobs/{YY}/{MM}/{DD}/{제목-kebab-case}/plan.md`

예시: `.do/jobs/26/01/07/user-auth/plan.md`

### Step 5: 완료 메시지

```
플랜 생성 완료: .do/jobs/{YY}/{MM}/{DD}/{제목-kebab-case}/plan.md

다음 단계:
- 플랜 검토 후 수정 필요시 알려주세요
- 구현 시작: "플랜대로 구현해줘"
```

---

## 예시

```
/do:plan
→ "플랜 제목을 입력하세요" 질문
→ "user-auth" 입력
→ 분석 수행
→ .do/jobs/26/01/07/user-auth/plan.md 생성

/do:plan api-rate-limiting
→ 바로 분석 시작
→ .do/jobs/26/01/07/api-rate-limiting/plan.md 생성
```

---

## 주의사항

- 코드베이스 분석은 충분히 수행하되 과도하게 깊이 들어가지 않기
- 플랜은 실행 가능한 수준으로 구체적으로 작성
- 기존 코드 스타일과 패턴을 존중
- 테스트 계획 반드시 포함
