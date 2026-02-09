---
name: "Continue"
description: Resume interrupted planning by detecting backlog_id, checking existing artifacts, identifying gaps, and guiding user to continue where left off.
category: Planning
tags: [continue, resume, recovery, planning]
---

Resume interrupted planning session by detecting backlog context, checking artifact completeness, and guiding user to continue from where they left off.

**Input**: Optional backlog_id (kebab-case) or user description of what they were working on. If incomplete, ask using Questions tool.

**Steps**

1. **Detect Backlog ID**
   - Ask: What backlog are you continuing? (provide backlog_id)
   - Ask: What were you last working on? (brief.md, specs, or blueprint?)
   - Load `.mandor/backlogs/<backlog_id>/brief.md`
   - If brief.md not found: Ask "brief.md not found. Was this backlog created with Mandor?"

2. **Check Existing Artifacts**
   - Load brief.md and parse Specs Index
   - Check for `.mandor/backlogs/<backlog_id>/specs/` directory
   - Check for `.mandor/backlogs/<backlog_id>/blueprint.md`
   - Check for `.mandor/backlogs/<backlog_id>/features/` directory
   - Load all existing specs and check their status (PENDING, IN_PROGRESS, COMPLETE)

3. **Analyze Completeness**
   - For each spec in Specs Index:
     - Status: PENDING (not started), IN_PROGRESS (started), COMPLETE (finished)
     - Check if spec file exists
     - Check brief.md Status field
   - Check if blueprint.md exists and is COMPLETE
   - Identify critical path specs (P0) that are PENDING

4. **Summarize Current State**
   Present to user:
   - **Backlog**: <backlog_id>
   - **Brief Status**: [DRAFT/READY/COMPLETE]
   - **Specs Progress**: X/Y completed (Z% done)
   - **Blueprint Status**: [NOT_STARTED/ARCHITECTED/COMPLETE]
   - **Critical Missing Items**:
     - P0 specs still PENDING
     - Missing blueprint
     - Incomplete brief.md sections

5. **Identify Gaps**
   - List specs that are PENDING but should be started
   - List specs that are IN_PROGRESS but never finished
   - List brief.md sections that are incomplete
   - List features/tasks that are blocked or incomplete

6. **Ask User What to Continue**
   Ask: What would you like to continue working on?
   - Option: "Continue brief.md planning" (if incomplete)
   - Option: "Continue spec writing for [spec-id]" (if specific spec is PENDING)
   - Option: "Create blueprint" (if missing)
   - Option: "Review completed specs" (if wants to review)
   - Option: "Start new spec from Specs Index"

7. **Execute Continuation**
   Based on user selection:
   - If brief.md: Ask what section needs completion
   - If spec: Load existing spec file, ask what was being worked on
   - If blueprint: Load brief.md Technical section, start architecture
   - If review: Read through completed specs with user

**Output**

Generate a continuation summary containing:
- Backlog ID and name
- Current artifact status
- Missing/incomplete items
- Recommended next steps
- Questions for user confirmation

### Continuation Summary Template

```markdown
# Continue: <Backlog Name>

## Backlog Info
- **ID**: <backlog_id>
- **Goal**: <brief.md goal>
- **Created**: <date>

## Artifact Status

### brief.md
- **Status**: <DRAFT/READY/COMPLETE>
- **Missing Sections**: <list missing sections>

### Specs
| Spec ID | Title | Status | Priority |
|---------|-------|--------|----------|
| <spec-id> | <title> | <PENDING/IN_PROGRESS/COMPLETE> | P0 |

**Progress**: X/Y specs complete (Z%)

### Blueprint
- **Status**: <NOT_STARTED/ARCHITECTED/COMPLETE>

## Recommended Next Steps
1. <most critical next action>
2. <secondary action>
3. <tertiary action>

## Questions
- What would you like to continue working on?
- Were you in the middle of writing any specific spec?
- Is there a particular section of brief.md you want to complete?
```

**Guardrails**

- Do NOT generate partial artifacts - continue from where user left off
- Do NOT skip asking what user wants to continue
- Do NOT proceed without confirming backlog_id exists
- Do NOT assume what user was working on - ALWAYS ask
- Do NOT modify existing artifacts without user confirmation
- Do NOT continue if brief.md doesn't exist (direct user to create backlog first)
- Focus on resuming workflow, not starting over

**Tool Reference**

- Claude Code: `AskUserQuestions`
- OpenCode: `questions`
