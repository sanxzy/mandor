---
name: "Project Planner"
description: Rigorous project planning for NEW/EXISTING projects. Generates comprehensive brief.md through exhaustive questioning.
category: Planning
tags: [planning, requirements, analysis]
---

**Input**: Project identifier (kebab-case) or description. Ask if incomplete.

**Steps**

1. **Identity** - Ask: name (kebab), title, NEW/EXISTING, directory, version, standalone/part.
   - EXISTING: Ask current state, tech stack, issues, debt. Skip greenfield.

2. **Business** - Ask: problem solved, stakeholders, success criteria, build/buy/extend, ROI, timeline, budget, competitors, market, lifetime, metrics, compliance, strategic alignment.

3. **Technical** - Ask: deliverable type, required/forbidden tech, integrations, architecture, components, APIs, third-party services, data storage, performance (latency/throughput), availability (99.9%+), security (encryption/auth), languages, databases, message queues, caching, monitoring, testing.

4. **Operational** - Ask: deployment (cloud/on-prem), scaling needs, compliance (SOC2/HIPAA/GDPR), maintenance lifespan, disaster recovery, backups, logging, monitoring, cost targets, IaC, containers (Docker/K8s), CI/CD, environments, networking, SSL/TLS, rate limiting, backwards compatibility.

5. **Team** - Ask: team size/structure, experience level, existing CI/CD tools, roles (FE/BE/DevOps/QA), on-call rotation, code review process, documentation standards, testing requirements (unit/integration/e2e), code coverage, sprint length, release frequency, communication tools, knowledge management, onboarding, tech debt approach.

6. **Risks** - Ask: biggest technical risks, external dependencies, constraints, critical path items, single points of failure, security vulnerabilities, performance bottlenecks, integration risks, data migration, vendor lock-in, scalability, complexity, tech obsolescence, skills gaps, scope creep, library dependencies.

7. **Validation** - Summarize all info. Ask user to confirm/correct EACH assumption. Explicitly list: non-goals, OUT OF SCOPE, IN SCOPE, ASSUMPTIONS, RISKS, CONSTRAINTS. Refuse to proceed without full confirmation.

8. **Generate brief.md** - Create `.mandor/projects/<project_id>/brief.md` with ALL gathered info including Background & Motivation, Scope (in/out), Repository Context (NEW/EXISTING), and **Specs Index** with human-readable kebab-case IDs.

**Output**

Generate `brief.md` with: Background & Motivation, Scope Definition, Repository Context, Specs Index (REQUIRED), all requirements, user validation confirmation.

### Specs Index Template

```markdown
## Specs Index

| Spec ID | Title | Description | Priority |
|---------|-------|-------------|----------|
| invitation-builder | Invitation Builder | Template editor | P0 |
| guest-management | Guest Management | RSVP tracking | P0 |
| authentication | Authentication | User auth | P0 |
| xendit-integration | Xendit Integration | Payment processing | P1 |
| calendar-export | Calendar Export | Calendar integration | P2 |
| islamic-templates | Islamic Templates | Islamic-specific templates | P2 |
| gender-separation | Gender Separation | Gender session config | P2 |
| wali-management | Wali Management | Wali information | P3 |
| cash-gifting | Cash Gifting | Digital envelope system | P3 |
| media-gallery | Media Gallery | Photo/video gallery | P4 |
| multi-language | Multi-language | RTL and translations | P4 |
| pwa-support | PWA Support | Progressive web app | P5 |

Spec ID must be kebab-case. Each spec created by `Spec Writer` skill.
```

**Guardrails**

- Generate NEW or UPDATE EXISTING brief.md only. Don't load existing docs.
- For EXISTING: Focus on discovering current state.
- Don't generate if ANY info incomplete. Don't skip ANY question/domain.
- Don't generate partial artifacts. Don't proceed without explicit confirmation.
- Don't infer silently - ALWAYS ask. Don't proceed if conflicts unresolved.
- Be exhaustive - no aspect skipped.

**Tool Reference**

- Claude Code: `AskUserQuestions`
- OpenCode: `questions`
