# Ask Discipline

Rules for when and how to ask questions to users. These are HARD rules.

## Self-Assessment Before Asking

Before using AskUserQuestion or asking ANY question to the user, the AI MUST complete this internal checklist:

1. **Restate the request**: What exactly did the user ask for?
2. **Restate the context**: What did the user explain or demonstrate?
3. **Restate the expectation**: What outcome does the user expect?
4. **Self-evaluate**: Can I decide this myself based on the above?

If all three restatements point to a clear action, execute it. Do NOT ask.

## When Questions Are Justified

Only ask when:
- The decision involves irreversible consequences (data loss, breaking changes)
- Multiple valid approaches exist with genuinely different trade-offs that depend on user preference
- Missing critical information that cannot be inferred from context
- The user explicitly requested to be consulted

## When Questions Are PROHIBITED

Do NOT ask when:
- The answer is obvious from context (user just explained what they want)
- You can make a reasonable default choice and explain it after
- The question is just confirming what the user already said
- You're asking "should I do X?" when the user literally just asked you to do X
- You're presenting options that have a clearly superior choice

## Violation Examples

VIOLATION: User says "rename godo sync to godo moai sync" -> AI asks "should I rename it?"
CORRECT: AI renames it and reports the change.

VIOLATION: User explains problem A, B, C in detail -> AI asks "what do you want me to do?"
CORRECT: AI synthesizes A+B+C into an action plan and executes.

VIOLATION: User says "yes/confirmed" -> AI asks for more confirmation
CORRECT: AI proceeds with the confirmed action.

## Decide-and-Explain Pattern

When the decision is not critical but has multiple options:
1. Choose the most reasonable option
2. Execute it
3. Explain WHY you chose that option
4. Mention alternatives briefly

This is preferred over asking the user to choose between options you already have an informed opinion about.

## Integration with AskUserQuestion Tool

AskUserQuestion constraints (from CLAUDE.md) still apply:
- Maximum 4 options per question
- No emoji in question text
- Questions in user's conversation language

Additional constraint:
- [HARD] Every AskUserQuestion call MUST be preceded by a visible self-assessment (internal reasoning about whether the question is truly necessary)
- If the self-assessment concludes the AI can decide, it MUST decide instead of asking
