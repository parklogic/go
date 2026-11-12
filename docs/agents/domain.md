# Domain Docs

How the engineering skills should consume this repo's domain documentation when exploring the codebase.

## Before exploring, read these

- **`GLOSSARY.md`** at the repo root.
- **`docs/adr/`**: read the ADRs that touch the area you are about to work in.

When a file is missing, proceed. The domain-modeling skill creates it when a term or a decision is resolved.

When an ADR names a successor, obey the successor. Wording under a History heading is the old decision.

## File structure

```
/
├── GLOSSARY.md
├── docs/adr/
└── docs/agents/
```

## Use the glossary's vocabulary

When your output names a domain concept (in an issue title, a refactor proposal, a hypothesis, a test name), use the term as defined in `GLOSSARY.md`. Stay with the glossary when it names a term to avoid.

When the concept you need is not in the glossary yet, you are either inventing language the project does not use, or there is a real gap. Note a gap for the domain-modeling skill.

## Flag ADR conflicts

When your output contradicts an existing ADR, name that ADR and the rule, and say why reopening it is worth doing.
