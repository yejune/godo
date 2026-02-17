# 워크플로우: Team Review - 다관점 코드 리뷰

목적: 여러 관점에서 동시에 코드 변경 사항을 리뷰한다. 각 리뷰어가 특정 품질 차원에 집중한다.

흐름: TeamCreate -> 관점 배정 -> 병렬 리뷰 -> 보고서 통합

## 전제 조건

- workflow.team.enabled: true
- CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1
- 트리거: /do review --team 또는 명시적 다관점 리뷰 요청

## Phase 0: 리뷰 설정

1. 리뷰할 코드 변경 사항 식별 (diff, PR, 또는 파일 목록)
2. 팀 생성:
   ```
   TeamCreate(team_name: "do-review-{target}")
   ```
3. 리뷰 태스크 생성:
   ```
   TaskCreate: "보안 리뷰: OWASP 준수, 입력 검증, 인증" (의존성 없음)
   TaskCreate: "성능 리뷰: 알고리즘 복잡도, 리소스 사용, 캐싱" (의존성 없음)
   TaskCreate: "품질 리뷰: TRUST 5, 패턴, 유지보수성, 테스트 커버리지" (의존성 없음)
   TaskCreate: "리뷰 결과 통합" (위 태스크들에 의해 차단됨)
   ```

## Phase 1: 리뷰 팀 소환

리뷰 팀 패턴 사용:

팀원 1 - security-reviewer (team-quality 에이전트, inherit 모델):
- 프롬프트: "보안 이슈에 대해 다음 변경 사항을 리뷰하라. OWASP Top 10 준수, 입력 검증, 인증/인가, 시크릿 노출, 인젝션 위험 확인. 변경 사항: {diff_summary}"

팀원 2 - perf-reviewer (team-quality 에이전트, inherit 모델):
- 프롬프트: "성능 이슈에 대해 다음 변경 사항을 리뷰하라. 알고리즘 복잡도, 데이터베이스 쿼리 효율성, 메모리 사용, 캐싱 기회, 번들 크기 영향 확인. 변경 사항: {diff_summary}"

팀원 3 - quality-reviewer (team-quality 에이전트, inherit 모델):
- 프롬프트: "코드 품질에 대해 다음 변경 사항을 리뷰하라. TRUST 5 준수, 네이밍 컨벤션, 에러 처리, 테스트 커버리지, 문서화, 프로젝트 패턴과의 일관성 확인. 변경 사항: {diff_summary}"

## Phase 2: 병렬 리뷰

리뷰어들이 독립적으로 작업 (모두 읽기 전용):
- 각자 배정된 품질 차원에 집중
- 해당 관점으로 변경된 모든 파일 리뷰
- 각 발견 사항을 심각도별로 평가 (critical, warning, suggestion)
- 팀 리더에게 결과 보고

## Phase 3: 보고서 통합

모든 리뷰 완료 후:
1. 모든 리뷰어의 결과 수집
2. 겹치는 이슈 중복 제거
3. 심각도 순 우선순위 지정 (critical 우선)
4. 다음 내용을 포함하여 사용자에게 통합 리뷰 보고서 제시:
   - 즉시 주의가 필요한 critical 이슈
   - 처리해야 할 warning
   - 개선 제안
   - TRUST 5 차원별 전반적인 품질 평가

## Phase 4: 정리

1. 모든 리뷰 팀원 종료
2. 리소스 정리를 위한 TeamDelete
3. 선택적으로 critical 이슈에 대한 수정 태스크 생성

## 폴백

팀 생성 실패 시:
- 단일 관점 리뷰를 위해 manager-quality 서브에이전트로 폴백
- 보안, 성능, 품질 순으로 순차 리뷰

---

Version: 1.0.0
