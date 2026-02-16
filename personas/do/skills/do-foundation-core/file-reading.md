# File Reading Optimization

## Progressive Loading by Size

- Under 200 lines: Read full file
- 200-500 lines: Grep to find line numbers, then Read with offset/limit
- 500-1000 lines: Never read full file. Grep first, read 50-100 line chunks
- Over 1000 lines: Use Grep with context (-C) instead of Read. Delegate to Explore agent if full understanding needed

## Practical Rules

- Before reading: Use Grep to find exact line numbers of interest
- When modifying: Read only the section being modified + 20 lines context
- When exploring: Start with Grep for entry points, read only relevant sections

## Token Budget

Each line ≈ 10-20 tokens. A 1000-line file costs 10K-20K tokens.
Targeted 50-line reads from 10 files = 5-10K tokens (vs 50-100K for full reads).
