---
name: "Blueprint Designer"
description: Generates architectural blueprint.md from brief.md. Documents system design, components, data flows, decisions. Runs in parallel with Spec Writer.
category: Planning
tags: [architecture, design, blueprint]
---

**Input**: Optional project ID (kebab-case). Ask if missing.

**PREREQUISITE CHECK**

1. Get project ID from input or ask user.
2. Verify `.mandor/projects/<project_id>/brief.md` exists.
3. **ERROR if missing**: "brief.md required. Run `Project Planner` first."
4. Load brief.md for project context, architecture requirements, tech stack.

**Steps**

1. **Architecture Pattern** - Ask based on brief.md: preferred patterns (monolith/microservices/serverless/layered), main components/responsibilities, inter-component communication (REST/gRPC/events/MQ), deployment topology.

2. **Module Structure** - Ask: core modules from brief.md, boundaries between modules, inter-module dependencies, module organization (by feature/layer).

3. **Data Flow** - Ask: flow patterns (sync/async/event-driven), data storage locations (DB/cache/external), data movement through system, ownership/consistency models.

4. **Technology Stack** - Reference brief.md. Confirm: frontend tech/role, backend tech/role, DB tech/schema patterns, auth approach, external services integration.

5. **External Integrations** - Ask: integrated external services (from brief.md), exposed APIs, integration patterns (sync/async/webhook), API gateway/service mesh needs.

6. **Quality Attributes** - Ask based on brief.md: scalability (users/requests/data), availability (uptime %), security (encryption/auth/compliance), performance (latency/throughput).

7. **Infrastructure/DevOps** - Ask: deployment platform (cloud/on-prem/hybrid), containerization (Docker/K8s), CI/CD requirements, environment strategy (dev/staging/prod).

8. **Architectural Decisions** - Document each decision with rationale. For rejected alternatives: document what was rejected and why. Document trade-offs. Ask if unclear.

9. **Validation** - Summarize ALL architectural info. Ask user to confirm/correct EACH assumption. List all decisions, trade-offs. Refuse to proceed without confirmation.

10. **Generate blueprint.md** - Create `.mandor/projects/<project_id>/blueprint.md`. Reference brief.md as source of truth. Document high-level architecture, components/boundaries, data flows, decisions with rationale. Trace to brief.md. Independent from specs.

**Output**

Generate `blueprint.md`:

```markdown
# Architecture Blueprint: <Project Title>

**Project ID:** <project-id>
**Created:** <date>
**Brief Reference:** [brief.md](./brief.md)
**Status:** ARCHITECTED

## 1. Architectural Overview
- **Pattern:** Monolith/Microservices/Serverless/Layered
- **Key Design Principles**
- **Architectural Drivers** (from brief.md)

## 2. System Architecture Diagram
```
[ASCII/text diagram]
```

## 3. Component Architecture
- Frontend Components
- Backend Components
- Data Layer Components
- Integration Components
- Responsibilities and Boundaries

## 4. Module Structure
```
src/
├── modules/       # Feature modules
├── components/    # Shared components
├── lib/          # Utilities/configs
├── hooks/        # Custom hooks
└── types/        # TypeScript types
```

## 5. Data Architecture
- Data Models Overview
- Data Flow Diagrams
- Storage Strategy (DB/cache)
- Data Ownership

## 6. Integration Architecture
- External Services (from brief.md)
- API Design Patterns
- Integration Patterns (sync/async/webhook)
- Third-Party Dependencies

## 7. Quality Attributes
- Scalability
- Availability
- Security
- Performance

## 8. Infrastructure
- Deployment Architecture
- CI/CD Pipeline
- Environment Strategy
- Containerization

## 9. Architectural Decisions
| Area | Decision | Rationale | Alternatives |
|------|----------|-----------|--------------|
| Pattern | Microservices | Scalability | Monolith rejected |
| Database | PostgreSQL | ACID compliance | NoSQL rejected |

## 10. Trade-offs & Rejected Alternatives
List rejected alternatives and why.

## 11. Traceability
Map each decision to brief.md requirements. Document how blueprint satisfies brief.md.
```

**Guardrails**

- **PREREQUISITE: brief.md MUST exist** - ERROR if missing.
- **INDEPENDENT from Spec Writer** - Can run in PARALLEL.
- **Based on brief.md only** - Don't read specs. Trace decisions to brief.md.
- Don't create partial blueprints. Don't skip decision rationale or rejected alternatives.
- Don't make silent assumptions - ALWAYS ask.
- Be exhaustive.

**Tool Reference**

- Claude Code: `AskUserQuestions`
- OpenCode: `questions`
