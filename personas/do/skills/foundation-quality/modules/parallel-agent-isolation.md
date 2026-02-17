# Parallel Agent Commit Isolation

Rules ensuring safe concurrent commits when multiple agents work in parallel.

## Rules [HARD]

- [HARD] `git add` -> `git diff --cached` -> `git commit` must execute as a **single Bash call** -- prevent other agents from interleaving between staging and commit
  - Correct: `git add file1.go file2.go && git diff --cached --name-only && git commit -m "msg"`
  - Forbidden: `git add file1.go` (separate call) -> (time passes) -> `git commit` (separate call)
- [HARD] `git reset HEAD` is absolutely forbidden -- unstaging files staged by another agent causes their changes to be lost
- [HARD] If other files exist in the staging area beyond your own, **never touch them** -- do not unstage/reset/checkout other agents' files
- [HARD] On commit failure (conflict, staging contamination, etc.) **do not attempt self-recovery** -- report the error as-is and wait for the orchestrator to resolve it
  - Reason: self-recovery attempts by an agent can destroy other agents' work -- it is safer to surface the conflict
- [HARD] The orchestrator must instruct parallel agents: "Execute staging + commit in one call, and never touch other agents' files"

## Safe Commit Pattern

```bash
# Single atomic call -- safe for parallel agents
git add file1.go file2.go && \
  git diff --cached --name-only && \
  git commit -m "$(cat <<'EOF'
feat: implement login handler

Add JWT token generation and validation
EOF
)"
```

## Failure Scenarios

| Scenario | Correct Response | Forbidden Response |
|----------|-----------------|-------------------|
| Unknown files in staging area | Ignore, commit only your files | `git reset HEAD` to clean staging |
| Commit conflict | Report error to orchestrator | Attempt merge/rebase |
| Staging contaminated | Report error to orchestrator | `git checkout .` to discard |

## Related

- [Commit Discipline](commit-discipline.md)
- [Coding Discipline](coding-discipline.md)
