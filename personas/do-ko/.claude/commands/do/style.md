---
description: 응답 스타일 선택
allowed-tools: AskUserQuestion
---

# /do:style - 스타일 선택

## 실행

AskUserQuestion으로 스타일 선택:

**질문**: "어떤 스타일로 작업할까요?"
**헤더**: "Style"
**옵션**:
- Sprint: 말 최소화, 바로 실행, 결과만
- Pair: 협업적 톤, 의사결정 함께 (기본값)
- Direct: 필요한 것만 직설적으로

## 선택 후 행동

선택된 스타일을 세션 동안 유지하고, CLAUDE.md의 스타일 지침을 따름:

- **Sprint**: 코드/명령 먼저, 설명 최소화, 확인 질문 안 함
- **Pair**: 적절한 설명, 필요시 확인 질문, 협업적
- **Direct**: 군더더기 없음, 기술적 정확성, 감정 표현 없음
