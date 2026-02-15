# Persona System Rebuild Checklist

**Plan**: plan.md
**Architecture**: architecture-persona-system.md
**Analysis**: analysis-requirements.md

---

## 작업 목록

- [ ] #1 Phase 1a: 캐릭터 파일 YAML frontmatter 재작성 (4파일)
  - 담당: expert-backend
  - 서브: checklists/01_persona-files.md
  - 파일: characters/young-f.md, young-m.md, senior-f.md, senior-m.md

- [ ] #2 Phase 1b: 스피너 YAML 파일 신규 생성 (4파일)
  - 담당: expert-backend
  - 서브: checklists/01_persona-files.md
  - 파일: spinners/young-f.yaml, young-m.yaml, senior-f.yaml, senior-m.yaml
  - depends on: #1 (캐릭터 형식 확정 후 스피너 작성)

- [ ] #3 Phase 1c: 스타일 파일 새 위치 생성 + persona.yaml 업데이트 (4파일)
  - 담당: expert-backend
  - 서브: checklists/01_persona-files.md
  - 파일: styles/sprint.md, pair.md, direct.md, persona.yaml

- [ ] #4 Phase 2a: PersonaLoader 구조체 + LoadCharacter 구현 (1파일)
  - 담당: expert-backend
  - 서브: checklists/02_persona-loader.md
  - 파일: cmd/godo/persona_loader.go
  - depends on: #1 (파일 형식 확정)

- [ ] #5 Phase 2b: LoadSpinner + BuildSpinnerVerbs 구현 (1파일)
  - 담당: expert-backend
  - 서브: checklists/02_persona-loader.md
  - 파일: cmd/godo/persona_loader.go (같은 파일, 스피너 부분)
  - depends on: #2, #4

- [ ] #6 Phase 2c: PersonaLoader 단위 테스트 + 테스트 픽스처 (2파일)
  - 담당: expert-testing
  - 서브: checklists/02_persona-loader.md
  - 파일: cmd/godo/persona_loader_test.go, cmd/godo/testdata/
  - depends on: #4, #5

- [ ] #7 Phase 3a: hook_session_start.go에 PersonaLoader 연결 (1파일)
  - 담당: expert-backend
  - 서브: checklists/03_hook-integration.md
  - 파일: cmd/godo/hook_session_start.go
  - depends on: #6 (테스트 통과 후)

- [ ] #8 Phase 3b: hook_post_tool_use.go에 PersonaLoader 연결 (1파일)
  - 담당: expert-backend
  - 서브: checklists/03_hook-integration.md
  - 파일: cmd/godo/hook_post_tool_use.go
  - depends on: #6

- [ ] #9 Phase 3c: hook_user_prompt_submit.go에 PersonaLoader 연결 (1파일)
  - 담당: expert-backend
  - 서브: checklists/03_hook-integration.md
  - 파일: cmd/godo/hook_user_prompt_submit.go
  - depends on: #6

- [ ] #10 Phase 3d: spinner.go에 LoadSpinner 연결 (1파일)
  - 담당: expert-backend
  - 서브: checklists/03_hook-integration.md
  - 파일: cmd/godo/spinner.go
  - depends on: #6

- [ ] #11 Phase 4a: PersonaManifest 모델에 Characters/Spinners 필드 추가 (1파일)
  - 담당: expert-backend
  - 서브: checklists/04_convert-tool.md
  - 파일: internal/model/persona_manifest.go
  - depends on: #1 (파일 형식 확정)

- [ ] #12 Phase 4b: CharacterExtractor + SpinnerExtractor 생성 (2파일)
  - 담당: expert-backend
  - 서브: checklists/04_convert-tool.md
  - 파일: internal/extractor/character.go, internal/extractor/spinner.go
  - depends on: #11

- [ ] #13 Phase 4c: ExtractorOrchestrator 라우팅 + Assembler 복사 로직 (2파일)
  - 담당: expert-backend
  - 서브: checklists/04_convert-tool.md
  - 파일: internal/extractor/orchestrator.go, internal/assembler/orchestrator.go
  - depends on: #12

- [ ] #14 Phase 5a: sprint.md MoAI 잔재 제거 (1파일)
  - 담당: expert-backend
  - 서브: checklists/05_style-cleanup.md
  - 파일: personas/do/styles/sprint.md
  - depends on: #3 (스타일 파일이 새 위치에 존재)

- [ ] #15 Phase 5b: pair.md MoAI 잔재 제거 + 크기 축소 (1파일)
  - 담당: expert-backend
  - 서브: checklists/05_style-cleanup.md
  - 파일: personas/do/styles/pair.md
  - depends on: #3

- [ ] #16 Phase 5c: direct.md MoAI 잔재 제거 (1파일)
  - 담당: expert-backend
  - 서브: checklists/05_style-cleanup.md
  - 파일: personas/do/styles/direct.md
  - depends on: #3

- [ ] #17 전체 검증: 12 조합 동작 확인 + MoAI 잔재 최종 grep
  - 담당: expert-testing
  - 서브: checklists/05_style-cleanup.md
  - depends on: #7, #8, #9, #10, #14, #15, #16

- [ ] #18 Phase 6a: manager-ddd.md, manager-tdd.md 에이전트 정의 Do 철학 반영 (2파일)
  - 담당: expert-backend
  - 서브: checklists/06_agents-philosophy.md
  - 파일: agents/do/manager-ddd.md, agents/do/manager-tdd.md

- [ ] #19 Phase 6b: manager-quality.md, manager-project.md, team-quality.md Do 철학 반영 (3파일)
  - 담당: expert-backend
  - 서브: checklists/06_agents-philosophy.md
  - 파일: agents/do/manager-quality.md, agents/do/manager-project.md, agents/do/team-quality.md
  - depends on: #18 (일관된 패턴 확립 후)

- [ ] #20 Phase 7: rules/ 개발 규칙 갱신 (2파일)
  - 담당: expert-backend
  - 서브: checklists/07_rules-update.md
  - 파일: rules/do/workflow/spec-workflow.md, rules/do/workflow/workflow-modes.md

- [ ] #21 Phase 8a: SKILL.md + reference.md 스킬 철학 반영 (2파일)
  - 담당: expert-backend
  - 서브: checklists/08_skills-philosophy.md
  - 파일: skills/do/SKILL.md, skills/do/references/reference.md

- [ ] #22 Phase 8b: workflows/ 스킬 파일 Do 철학 반영 (6파일)
  - 담당: expert-backend
  - 서브: checklists/08_skills-philosophy.md
  - 파일: skills/do/workflows/do.md, team-do.md, report.md, run.md, plan.md, test.md
  - depends on: #21 (SKILL.md 패턴 확립 후)

- [ ] #23 Phase 9: commands/ 확인 및 갱신 (5파일)
  - 담당: expert-backend
  - 서브: checklists/09_commands-update.md
  - 파일: commands/do/setup.md, style.md, mode.md, checklist.md, plan.md

- [ ] #24 Phase 6-9 전체 검증: MoAI 잔재 최종 grep
  - 담당: expert-testing
  - 서브: checklists/09_commands-update.md
  - depends on: #18, #19, #20, #21, #22, #23

---

## 병렬 실행 가능 그룹

### Group A: Phase 1 (순차)
#1 → #2 → #3

### Group B: Phase 2 + Phase 3 (Phase 1 완료 후)
#4 → #5 → #6 → (#7, #8, #9, #10 병렬)

### Group C: Phase 4 (Phase 1 완료 후, Group B와 병렬)
#11 → #12 → #13

### Group D: Phase 5 (#3 완료 후, Group B/C와 병렬)
(#14, #15, #16 병렬) → #17

### Group E: Phase 6 (Group A-D와 독립, 병렬 가능)
#18 → #19

### Group F: Phase 7 (Group A-D와 독립, 병렬 가능)
#20

### Group G: Phase 8 (Group A-D와 독립, 병렬 가능)
#21 → #22

### Group H: Phase 9 (Group A-D와 독립, 병렬 가능)
#23

### Group I: Phase 6-9 최종 검증
#24 (depends on: #18-#23 전체 완료)

---

## 상태 요약

| Phase | 항목 수 | 상태 |
|-------|--------|------|
| Phase 1 (파일 생성) | #1, #2, #3 | [ ] |
| Phase 2 (PersonaLoader) | #4, #5, #6 | [ ] |
| Phase 3 (Hook 연결) | #7, #8, #9, #10 | [ ] |
| Phase 4 (convert 도구) | #11, #12, #13 | [ ] |
| Phase 5 (스타일 정리) | #14, #15, #16, #17 | [ ] |
| Phase 6 (agents 철학 반영) | #18, #19 | [ ] |
| Phase 7 (rules 갱신) | #20 | [ ] |
| Phase 8 (skills 철학 반영) | #21, #22 | [ ] |
| Phase 9 (commands 갱신) | #23, #24 | [ ] |
