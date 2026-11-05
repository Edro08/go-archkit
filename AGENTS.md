# Agent Instructions

Read the applicable specifications in `agents/specs/` before modifying code:
- `agents/specs/architecture.md`: package boundaries, dependency direction, and composition rules.

## Context and Progress Tracking

Agents may use the `agents/context/` directory to deposit and maintain context `.md` files for task tracking, reporting, and knowledge retention:
- `history.md`: chronicle of project milestones, phases, and snapshots. Must start with `# Historial de Archkit` and list phases/milestones in reverse chronological order (most recent first / cortes en orden cronológico descendente).
- `phases.md` / `fases.md`: implementation phases, roadmaps, and step-by-step checklists.
- `reports.md` / `reportes.md`: analysis reports, technical audits, and benchmark summaries.
- Additional markdown files as needed to retain execution context across sessions.
