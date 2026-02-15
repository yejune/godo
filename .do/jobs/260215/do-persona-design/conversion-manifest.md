# Do Persona Conversion Manifest

## File Origin Tracking

모든 `personas/do/` 파일의 출처, 변환 유형, 주의사항을 추적한다.

### Persona Root Files

| Do 파일 | 출처 (moai) | 변환 유형 | 수동 검토 필요 | 비고 |
|---------|-----------|----------|-------------|------|
| `CLAUDE.md` | `personas/moai/CLAUDE.md` | structural | YES | 삼원 구조, 체크리스트 워크플로우, 한국어 혼합으로 재작성 |
| `manifest.yaml` | `personas/moai/manifest.yaml` | structural | YES | hook_scripts 빈 배열, slot_content/agent_patches 비어있음 |
| `settings.json` | `personas/moai/settings.json` | structural | YES | godo 직접 호출, outputStyle:pair, SubagentStop/UserPromptSubmit 추가 |

### Agents (5개 persona agents)

| Do 파일 | 출처 (moai) | 변환 유형 | 수동 검토 필요 | 비고 |
|---------|-----------|----------|-------------|------|
| `agents/do/manager-ddd.md` | `agents/moai/manager-ddd.md` | mechanical+structural | YES | hooks->godo, skills: moai-->do-, memory: .moai/-->.do/ |
| `agents/do/manager-tdd.md` | `agents/moai/manager-tdd.md` | mechanical+structural | YES | 동일 변환 패턴 |
| `agents/do/manager-quality.md` | `agents/moai/manager-quality.md` | mechanical+structural | YES | 동일 변환 패턴 |
| `agents/do/manager-project.md` | `agents/moai/manager-project.md` | mechanical+structural | YES | 동일 변환 패턴 |
| `agents/do/team-quality.md` | `agents/moai/team-quality.md` | mechanical+structural | YES | 동일 변환 패턴 |
| (없음) | `agents/moai/manager-spec.md` | removed | N/A | Do는 SPEC 미사용, plan workflow로 대체 |

### Skills (Orchestrator)

| Do 파일 | 출처 (moai) | 변환 유형 | 수동 검토 필요 | 비고 |
|---------|-----------|----------|-------------|------|
| `skills/do/SKILL.md` | `skills/moai/SKILL.md` | structural | YES | Mode Router 추가, Intent Router 재설계, Execution Directive 8단계 |
| `skills/do/workflows/plan.md` | `skills/moai/workflows/plan.md` | structural | YES | SPEC -> Analysis/Architecture/Plan 3단계 |
| `skills/do/workflows/run.md` | `skills/moai/workflows/run.md` | structural | YES | SPEC 기반 -> Checklist 기반 |
| `skills/do/workflows/do.md` | `skills/moai/workflows/moai.md` | structural | YES | 자동 파이프라인. 이름: moai.md -> do.md |
| `skills/do/workflows/team-do.md` | `skills/moai/workflows/team-plan.md` + `team-run.md` | structural | YES | Team Plan+Run 통합 |
| `skills/do/workflows/report.md` | (없음 -- Do 고유) | original | YES | MoAI sync에 해당하지만 완전히 다른 구현 |
| `skills/do/workflows/test.md` | (없음 -- Do 고유) | original | YES | 테스트 워크플로우 |
| `skills/do/references/reference.md` | `skills/moai/references/reference.md` | mechanical+structural | YES | 브랜드 치환 + Do 설정 반영 |
| (없음) | `skills/moai/workflows/sync.md` | removed | N/A | Do는 report.md로 대체 |
| (없음) | `skills/moai/workflows/fix.md` | removed | N/A | Do에 해당 없음 |
| (없음) | `skills/moai/workflows/loop.md` | removed | N/A | Do에 해당 없음 |
| (없음) | `skills/moai/workflows/feedback.md` | removed | N/A | Do에 해당 없음 |
| (없음) | `skills/moai/workflows/project.md` | removed | N/A | godo setup으로 대체 |
| (없음) | `skills/moai/workflows/team-sync.md` | removed | N/A | Team sync 미지원 |
| (없음) | `skills/moai/workflows/team-debug.md` | removed | N/A | Do에 해당 없음 |
| (없음) | `skills/moai/workflows/team-review.md` | removed | N/A | Do에 해당 없음 |
| (없음) | `skills/moai/moai-workflow-team/SKILL.md` | removed | N/A | workflows/에 통합 |

### Override Skills (moai -> Do 대응)

| moai Override Skill | Do 대응 | 변환 유형 | 비고 |
|---------------------|---------|----------|------|
| `skills/moai-foundation-core/SKILL.md` | (없음) | absorbed | CLAUDE.md + rules/*.md에 흡수 |
| `skills/moai-foundation-quality/SKILL.md` | (없음) | absorbed | dev-testing.md + dev-workflow.md에 흡수 |
| `skills/moai-workflow-ddd/SKILL.md` | (없음) | absorbed | rules/do/workflow/workflow-modes.md에 흡수 |
| `skills/moai-workflow-tdd/SKILL.md` | (없음) | absorbed | dev-workflow.md TDD 섹션에 흡수 |
| `skills/moai-workflow-spec/SKILL.md` | (없음) | absorbed | .do/jobs/ + dev-checklist.md에 흡수 |
| `skills/moai-workflow-project/SKILL.md` | (없음) | absorbed | godo init/setup으로 대체 |
| `skills/moai-workflow-testing/` (modules) | (없음) | absorbed | dev-testing.md에 흡수 |

### Output Styles

| Do 파일 | 출처 (moai) | 변환 유형 | 수동 검토 필요 | 비고 |
|---------|-----------|----------|-------------|------|
| `output-styles/do/pair.md` | `output-styles/moai/moai.md` | structural | YES | 오케스트레이터 -> 협업 스타일로 재작성 |
| `output-styles/do/sprint.md` | `output-styles/moai/r2d2.md` | structural | YES | 페어 프로그래밍 -> 빠른 실행 스타일로 재작성 |
| `output-styles/do/direct.md` | `output-styles/moai/yoda.md` | structural | YES | 교육 마스터 -> 직설적 전문가로 재작성 |

**주의**: moai 3 스타일과 Do 3 스타일은 1:1 rename이 아니다. 성격과 목적이 다르므로 재작성 필요.

### Commands

| Do 파일 | 출처 (moai) | 변환 유형 | 수동 검토 필요 | 비고 |
|---------|-----------|----------|-------------|------|
| `commands/do/check.md` | (없음 -- Do 고유) | original | NO | Do 고유 커맨드 |
| `commands/do/checklist.md` | (없음 -- Do 고유) | original | NO | Do 고유 커맨드 |
| `commands/do/mode.md` | (없음 -- Do 고유) | original | NO | 3모드 전환 |
| `commands/do/plan.md` | (없음 -- Do 고유) | original | NO | Do 고유 커맨드 |
| `commands/do/setup.md` | (없음 -- Do 고유) | original | NO | Do 고유 커맨드 |
| `commands/do/style.md` | (없음 -- Do 고유) | original | NO | Do 고유 커맨드 |
| (없음) | `commands/moai/github.md` | removed | N/A | moai 프로젝트 특화 |
| (없음) | `commands/moai/99-release.md` | removed | N/A | moai-adk 전용 |

### Rules (Persona-specific)

| Do 파일 | 출처 (moai) | 변환 유형 | 수동 검토 필요 | 비고 |
|---------|-----------|----------|-------------|------|
| `rules/do/workflow/spec-workflow.md` | `rules/moai/workflow/spec-workflow.md` | mechanical+structural | YES | Plan/Run/Sync -> Plan/Run/Report |
| `rules/do/workflow/workflow-modes.md` | `rules/moai/workflow/workflow-modes.md` | mechanical | NO | DDD/TDD/Hybrid 동일, 브랜드명만 치환 |

### Hooks

| moai Hook 파일 | Do 대응 | 변환 유형 | 비고 |
|---------------|---------|----------|------|
| `hooks/moai/handle-session-start.sh` | (없음) | removed | `godo hook session-start` 직접 호출 |
| `hooks/moai/handle-session-end.sh` | (없음) | removed | `godo hook session-end` 직접 호출 |
| `hooks/moai/handle-pre-tool.sh` | (없음) | removed | `godo hook pre-tool` 직접 호출 |
| `hooks/moai/handle-post-tool.sh` | (없음) | removed | `godo hook post-tool-use` 직접 호출 |
| `hooks/moai/handle-compact.sh` | (없음) | removed | `godo hook compact` 직접 호출 |
| `hooks/moai/handle-stop.sh` | (없음) | removed | `godo hook stop` 직접 호출 |
| `hooks/moai/handle-agent-hook.sh` | (없음) | removed | `godo hook subagent-stop` 직접 호출 |

### Characters (Do 고유)

| Do 파일 | 출처 (moai) | 변환 유형 | 수동 검토 필요 | 비고 |
|---------|-----------|----------|-------------|------|
| `characters/young-f.md` | (없음 -- Do 고유) | original | NO | MoAI에 캐릭터 시스템 없음 |
| `characters/young-m.md` | (없음 -- Do 고유) | original | NO | |
| `characters/senior-f.md` | (없음 -- Do 고유) | original | NO | |
| `characters/senior-m.md` | (없음 -- Do 고유) | original | NO | |

---

## Conversion Types

| 유형 | 코드 | 설명 | 자동화 가능 |
|------|------|------|-----------|
| 기계적 치환 | `mechanical` | find/replace만으로 충분 | YES |
| 구조 변환 | `structural` | 섹션 구조/개념/워크플로우 다름, 재작성 필요 | NO |
| 새로 작성 | `original` | moai에 없는 Do 고유 파일 | NO |
| 보존 | `preserve` | 변경 불필요 | N/A |
| 삭제 | `removed` | moai에 있지만 Do에 불필요 | YES |
| 흡수 | `absorbed` | moai skill이 Do rules에 이미 통합 | N/A |

### 유형별 파일 수

| 유형 | 파일 수 | 예시 |
|------|---------|------|
| mechanical | 3 | workflow-modes.md, reference.md (부분), agent 경로 치환 |
| structural | 12 | CLAUDE.md, SKILL.md, manifest.yaml, settings.json, workflows, styles |
| original | 12 | 6 commands, 4 characters, report.md, test.md |
| removed | 16 | 7 hooks, 5 workflows, 2 moai commands, team-sync/debug |
| absorbed | 7 | 6 override skills + moai-workflow-testing modules |

---

## Manual Review Required Items

### CRITICAL (moai 업데이트 시 반드시 검토)

1. **CLAUDE.md**: 새 HARD 규칙, 새 섹션, 구조 변경
2. **SKILL.md**: Intent Router, Execution Directive, 에이전트 카탈로그
3. **manifest.yaml**: 새 필드, 스키마 변경

### HIGH (moai 업데이트 시 주로 검토)

4. **Agent frontmatter**: hooks, skills, memory 필드 구조 변경
5. **Workflow files**: 새 워크플로우 추가, 기존 구조 변경
6. **Slot 변수**: registry.yaml에 새 `{{slot:*}}` 패턴 추가

### MEDIUM (주기적 검토)

7. **Output styles**: 스타일 구조 변경
8. **Settings.json**: 새 hooks, permissions, env 변수

---

## Automated Verification Commands

### 변환 후 자동 실행 체크리스트

```bash
# 변수 설정
DO_PERSONA_DIR=./personas/do

# === 1. moai 참조 잔존 확인 ===
echo "=== 1. moai 참조 잔존 확인 ==="
MOAI_REFS=0
for f in ${DO_PERSONA_DIR}/CLAUDE.md \
         ${DO_PERSONA_DIR}/skills/do/SKILL.md \
         ${DO_PERSONA_DIR}/skills/do/workflows/*.md \
         ${DO_PERSONA_DIR}/skills/do/references/*.md \
         ${DO_PERSONA_DIR}/commands/do/*.md \
         ${DO_PERSONA_DIR}/output-styles/do/*.md; do
  if [ -f "$f" ]; then
    COUNT=$(grep -ci "moai" "$f" 2>/dev/null || echo 0)
    if [ "$COUNT" -gt 0 ]; then
      echo "  WARNING: $f has $COUNT moai references"
      MOAI_REFS=$((MOAI_REFS + COUNT))
    fi
  fi
done
echo "  Total moai references: ${MOAI_REFS} (should be 0)"

# === 2. Do 정체성 요소 확인 ===
echo "=== 2. Do 정체성 요소 확인 ==="
grep -q "나는 Do다" ${DO_PERSONA_DIR}/CLAUDE.md && echo "  [OK] Do 선언문" || echo "  [FAIL] Do 선언문 누락"
grep -q "나는 Focus다" ${DO_PERSONA_DIR}/CLAUDE.md && echo "  [OK] Focus 선언문" || echo "  [FAIL] Focus 선언문 누락"
grep -q "나는 Team이다" ${DO_PERSONA_DIR}/CLAUDE.md && echo "  [OK] Team 선언문" || echo "  [FAIL] Team 선언문 누락"
grep -q "삼원 실행 구조" ${DO_PERSONA_DIR}/CLAUDE.md && echo "  [OK] 삼원 구조" || echo "  [FAIL] 삼원 구조 누락"

# === 3. godo Hook 패턴 확인 ===
echo "=== 3. godo Hook 패턴 확인 ==="
grep -q "godo hook session-start" ${DO_PERSONA_DIR}/settings.json && echo "  [OK] session-start" || echo "  [FAIL] session-start"
grep -q "godo hook stop" ${DO_PERSONA_DIR}/settings.json && echo "  [OK] stop" || echo "  [FAIL] stop"
grep -q "godo hook subagent-stop" ${DO_PERSONA_DIR}/settings.json && echo "  [OK] subagent-stop" || echo "  [FAIL] subagent-stop"
grep -q "godo hook user-prompt-submit" ${DO_PERSONA_DIR}/settings.json && echo "  [OK] user-prompt-submit" || echo "  [FAIL] user-prompt-submit"

# === 4. Shell wrapper 부재 확인 ===
echo "=== 4. Shell wrapper 부재 확인 ==="
if [ -d "${DO_PERSONA_DIR}/hooks" ]; then
  SH_COUNT=$(find ${DO_PERSONA_DIR}/hooks -name "*.sh" 2>/dev/null | wc -l)
  echo "  Shell scripts: ${SH_COUNT} (should be 0)"
else
  echo "  [OK] No hooks directory"
fi

# === 5. 페르소나/스타일/커맨드 파일 확인 ===
echo "=== 5. 파일 존재 확인 ==="
for p in young-f young-m senior-f senior-m; do
  [ -f "${DO_PERSONA_DIR}/characters/${p}.md" ] && echo "  [OK] character/${p}" || echo "  [FAIL] character/${p}"
done
for s in sprint pair direct; do
  [ -f "${DO_PERSONA_DIR}/output-styles/do/${s}.md" ] && echo "  [OK] style/${s}" || echo "  [FAIL] style/${s}"
done
for c in check checklist mode plan setup style; do
  [ -f "${DO_PERSONA_DIR}/commands/do/${c}.md" ] && echo "  [OK] cmd/${c}" || echo "  [FAIL] cmd/${c}"
done

# === 6. manifest.yaml 필수 필드 ===
echo "=== 6. manifest.yaml ==="
grep -q "name: do" ${DO_PERSONA_DIR}/manifest.yaml 2>/dev/null && echo "  [OK] name: do" || echo "  [FAIL] name"
grep -q "brand: do" ${DO_PERSONA_DIR}/manifest.yaml 2>/dev/null && echo "  [OK] brand: do" || echo "  [FAIL] brand"

echo ""
echo "=== Verification Complete ==="
```

### 개별 확인 명령어

```bash
# moai 문자열 잔존 (persona 파일)
grep -ri "moai" personas/do/CLAUDE.md personas/do/skills/ personas/do/commands/ personas/do/output-styles/ 2>/dev/null

# .moai/ 경로 잔존
grep -r "\.moai/" personas/do/ 2>/dev/null

# moai- 스킬 접두사 잔존
grep -r "moai-" personas/do/ 2>/dev/null | grep -v "moai-constitution"

# /moai 커맨드 잔존
grep -r "/moai" personas/do/ 2>/dev/null

# SPEC-XXX 패턴 잔존
grep -r "SPEC-" personas/do/ 2>/dev/null

# shell wrapper 참조 잔존
grep -r "handle-.*\.sh" personas/do/ 2>/dev/null
grep -r "hooks/moai/" personas/do/ 2>/dev/null

# 삼원 구조 건재 확인
grep "나는 Do다" personas/do/CLAUDE.md
grep "나는 Focus다" personas/do/CLAUDE.md
grep "나는 Team이다" personas/do/CLAUDE.md
```

---

## Conversion Decision Log

### 결정 1: Override Skills 생성하지 않음
- **일시**: 2026-02-15 (architecture.md)
- **근거**: Do는 "rules -> 항상 로드" 방식. MoAI의 "skill -> progressive disclosure"와 근본적으로 다름.
- **영향**: agent_patches 비어있음.

### 결정 2: manager-spec 제거
- **일시**: 2026-02-15 (architecture.md)
- **근거**: Do는 SPEC 문서 체계 미사용. plan workflow 내 에이전트 위임으로 대체.
- **영향**: agent 5개 (moai 6개에서 1개 감소)

### 결정 3: 스타일 1:1 매핑 아님
- **일시**: 2026-02-15 (redo-instructions.md)
- **근거**: moai의 moai/r2d2/yoda와 Do의 sprint/pair/direct은 성격이 다름.
- **영향**: output-styles 3개 모두 structural 변환.

### 결정 4: Hook shell wrapper 삭제
- **일시**: 2026-02-15 (research-do-philosophy.md)
- **근거**: godo binary 직접 호출이 28개 hook 이슈를 원천 제거.
- **영향**: manifest.yaml에 `hook_scripts: []`.

### 결정 5: slot_content/agent_patches 비어있음
- **일시**: 2026-02-15 (architecture.md)
- **근거**: TRUST 5 브랜딩 미사용, TAG Chain 미사용, override skills 없음.
- **주의**: assemble 후 `grep -r "{{slot:" .claude/` 실행하여 미치환 슬롯 확인 필요.

---

**작성자**: synthesizer agent
**작성일**: 2026-02-15
