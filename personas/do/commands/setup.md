# /do:setup

Do 환경 설정을 진행합니다.

## allowed-tools
Read, Write, Edit, AskUserQuestion

## Step 1: 현재 설정 확인

`.claude/settings.local.json` 파일을 읽어서 현재 설정 확인.
없으면 기본값 사용.

## Step 2-1: 기본 설정 (첫 번째 AskUserQuestion)

**먼저 시스템 사용자명을 가져옵니다:**
- `whoami` 명령 실행하여 시스템 사용자명 확인

AskUserQuestion으로 4개 질문:

1. **이름** (DO_USER_NAME)
   - 옵션 1: "{시스템사용자명} 사용 (추천)" - 설명: "현재 시스템 로그인 이름을 사용합니다"
   - 옵션 2: "설정 안 함" - 설명: "이름 없이 진행합니다"
   - 옵션 3 (Other): "직접 입력 (아래 입력칸에 원하는 이름 작성)"
   - 질문 텍스트: "사용자 이름을 선택하거나, 직접 입력하려면 맨 아래 입력칸을 사용하세요"

2. **대화 언어** (DO_LANGUAGE)
   - 한국어 (ko), English (en), 日本語 (ja), 中文 (zh)

3. **커밋 언어** (DO_COMMIT_LANGUAGE)
   - 한국어 (ko), English (en)

4. **페르소나** (DO_PERSONA)
   - young-f: 밝고 에너지 넘치는 20대 여성 천재 개발자
   - young-m: 자신감 넘치는 20대 남성 천재 개발자
   - senior-f: 30년 경력 레전드 50대 여성 천재 개발자
   - senior-m: 업계 전설 50대 남성 시니어 아키텍트

**처리 로직:**
- 옵션 1 선택: DO_USER_NAME = 시스템사용자명
- 옵션 2 선택: DO_USER_NAME = "" (빈 문자열)
- 옵션 3 (직접 입력): DO_USER_NAME = 입력된 값

## Step 2-2: 추가 설정 (두 번째 AskUserQuestion)

AskUserQuestion으로 4개 질문:

1. **응답 스타일** (DO_STYLE)
   - Sprint: 말 최소화, 바로 실행, 결과만
   - Pair (기본값): 협업적 톤, 의사결정 함께
   - Direct: 필요한 것만 직설적으로

2. **AI 푸터** (DO_AI_FOOTER)
   - 예 (true), 아니오 (false)

3. **실행 모드** (DO_MODE)
   - do: 모든 작업을 에이전트에게 위임 (대규모 작업)
   - focus: 코드를 직접 작성 (소규모 작업)
   - auto (기본값): 작업 규모에 따라 자동 선택

4. **Jobs 폴더 언어** (DO_JOBS_LANGUAGE)
   - English (en) (기본값), 한국어 (ko), 日本語 (ja), 中文 (zh)

## Step 3: 설정 저장

`.claude/settings.local.json` 파일 업데이트:

```json
{
  "env": {
    "DO_USER_NAME": "{이름}",
    "DO_LANGUAGE": "{언어코드}",
    "DO_COMMIT_LANGUAGE": "{커밋언어}",
    "DO_PERSONA": "{페르소나}",
    "DO_STYLE": "{스타일}",
    "DO_MODE": "{모드}",
    "DO_AI_FOOTER": "{true/false}",
    "DO_JOBS_LANGUAGE": "{jobs언어}"
  }
}
```

기존 settings.local.json 내용 유지하면서 env 필드만 업데이트.

## Step 3.5: 스피너 즉시 적용

Bash로 스피너를 바로 적용합니다:

```bash
godo spinner apply
```

## Step 4: 완료 메시지

설정 완료!
- 이름: {이름}
- 대화 언어: {언어}
- 커밋 언어: {커밋언어}
- 페르소나: {페르소나}
- 스타일: {스타일}
- 실행 모드: {모드}
- AI 푸터: {예/아니오}
- Jobs 폴더 언어: {jobs언어}
- 에이전트 확인: {예/아니오}

