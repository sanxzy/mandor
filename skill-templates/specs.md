---
name: "Spec Writer"
description: Generates comprehensive spec files from brief.md Specs Index. Creates independently testable, traceable specs. Comprehensive IAE scenarios.
category: Planning
tags: [specs, documentation, requirements]
---

**Input**: Optional project ID (kebab-case). Ask if missing.

**PREREQUISITE CHECK**

1. Get project ID from input or ask user.
2. Verify `.mandor/projects/<project_id>/brief.md` exists.
3. **ERROR if missing**: "brief.md required. Run `Project Planner` first."
4. Load Specs Index from brief.md. Process specs in priority order (P0→P1→P2→P3→P4→P5).

**Steps** (repeat for EACH spec in Specs Index)

1. **Feature Info** - Ask: capability, problem solved, module, users, business value.

2. **Happy Path** - Ask: user goal, successful outcome, main steps, preconditions, input, trigger data, expected output, timeout, success confirmation.

3. **Negative Cases** - Ask: invalid inputs, error conditions, boundary conditions, external service unavailability, missing/malformed data, insufficient permissions, rate limits, timeouts, validation rules, retry policies, error codes/messages.

4. **Input Validation** - Ask: required inputs, validation rules, data types/formats, min/max values, character limits, required vs optional, regex/patterns, custom validation, error messages, validation timing (client/server/both).

5. **Output/API** - Ask: output format (JSON/XML/HTML/binary), HTTP status codes, success/error response structures, fields, pagination, content types, versioning, rate limiting headers, caching headers.

6. **Security** - Ask: authentication methods, authorization roles/permissions, data access per role, operations per role, sensitive data protection, encryption, audit logging, CSRF/XSS/injection prevention.

7. **Database** - Ask: data entities, tables/collections, relationships, indexes, CRUD operations, transactions, migrations, seeding data.

8. **Integration** - Ask: external APIs/services, integration patterns (sync/async/webhook), timeouts, retries, circuit breakers, fallbacks, mock/test endpoints.

9. **UI/UX** - Ask: UI components, user flows, accessibility (a11y), i18n needs, responsive design, loading/error/success states.

10. **Performance** - Ask: latency requirements, throughput requirements, caching strategies, monitoring metrics, alerting thresholds, logging, performance tests.

11. **Testing** - Ask: unit/integration/E2E tests, test data, mock dependencies, coverage targets, edge cases.

12. **Assumptions & Dependencies** - Ask: system assumptions, spec dependencies, library versions, infrastructure assumptions, security/performance assumptions.

13. **Validation** - Summarize ALL info. Ask user to confirm/correct EACH assumption. List all IAE scenarios, negative cases, assumptions, dependencies. Refuse to proceed without full confirmation.

14. **Generate spec file** - Create `.mandor/projects/<project_id>/specs/<spec-id>.md` with human-readable kebab-case ID. Include ALL IAE scenarios, assumptions, dependencies. Ensure traceability and independent testability.

**Output**

Generate ONE spec file per Specs Index entry in priority order:

```markdown
# Spec: <Title>

**Spec ID:** <human-readable-kebab-case-id> (e.g., "invitation-builder")
**Module:** <module>
**Priority:** <P0-P5>
**Status:** PENDING
**Created:** <date>
**Brief Reference:** [brief.md](../brief.md)

## 1. Intent
- User Goal: ...
- Business Value: ...
- Success Criteria: ...

## 2. Scope
- In-Scope: ...
- Out-of-Scope: ...

## 3. IAE Scenarios
### Happy Path
| Intent | Action | Expect |
|--------|--------|--------|
| User goal | Steps | Expected outcome |

### Negative Cases
| Intent | Action | Expect |
|--------|--------|--------|
| Invalid input | Validation | Error response |
| Service down | Error handling | Fallback |

## 4. Input Validation
Field-by-field rules, error messages, timing.

## 5. Output Specifications
Success/error format, HTTP status, headers.

## 6. Security
Auth, authorization, sensitive data, encryption, audit.

## 7. Database
Models, CRUD, migrations.

## 8. Integration
External services, API contracts, error handling.

## 9. UI Requirements
Components, flows, accessibility.

## 10. Testing
Unit, integration, E2E, coverage, edge cases.

## 11. Performance
Latency, throughput, caching.

## 12. Monitoring
Metrics, alerts, logging.

## 13. Assumptions
List all assumptions.

## 14. Dependencies
List all dependencies.
```

After ALL specs: Update brief.md status from "Ready" to "Specs Complete".

**Guardrails**

- **PREREQUISITE: brief.md MUST exist** - ERROR if missing.
- Generate ALL specs from Specs Index in priority order. Don't skip any.
- Don't generate if info incomplete. Don't skip ANY aspect.
- Don't make silent assumptions - ALWAYS ask.
- Don't create partial specs. Don't skip IAE scenarios, security, performance, monitoring.
- Don't proceed without explicit confirmation.
- Be exhaustive.

**Tool Reference**

- Claude Code: `AskUserQuestions`
- OpenCode: `questions`
