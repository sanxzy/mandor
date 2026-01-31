# Agent Instructions - Mandor Project

## Output and Communication Rules

**Avoid Creating Summary Files**
- Do NOT create executive summary files unless explicitly asked
- Do NOT create summary.md or recap files
- Do NOT create files ending with "_SUMMARY" unless requested

**Focus on Action Over Documentation**
1. Make the fixes/changes
2. Test to verify they work
3. Provide concise output about what was done
4. Only document if user asks

**Response Style**
- Keep responses short and to the point
- Provide direct answers without lengthy preamble
- Avoid unnecessary explanations or summaries
- One-to-three sentences when possible
- Focus on what was accomplished, not how it was done

## Development Guidelines

This is the Mandor CLI project. When making changes:

1. Update the relevant source files in `internal/`
2. Build and test with `go build -o mandor-cli ./cmd/mandor`
3. Run CLI commands to verify fixes work
4. Document critical bugs in BUG_REPORT.md if needed
5. Create test scenarios to verify behavior

See root AGENTS.md for session completion and git push requirements.
