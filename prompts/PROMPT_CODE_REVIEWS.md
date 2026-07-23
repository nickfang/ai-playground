## Code Review Instructions
### Always read the actual files
- Never rely on cached or stale context. Re-read files from disk before making claims
about their state.
- If the user says "take another look", read the files again — they may have changed.
- Never assume a file exists, a library is used, or a configuration is present without
verifying.
- Always verify the actual state of the codebase before making statements, suggestions,
or claiming an issue exists. Do not present a hypothetical problem as a fact.
### Commenting & Documentation Philosophy
Document the *Why*, not the *What*.
- **Avoid Redundancy:** Do not explain *what* the code is doing (e.g., `i++ // increment
i`). The code itself should be readable enough to convey the "what".
- **Provide Context:** Use comments to explain the reasoning, trade-offs, and
architectural decisions behind a specific implementation.
- **Explain the "Why":** Why was this specific library chosen? Why is there a pointer
here instead of a value? Why is this constant set to 45 minutes?
- **Future-Proofing:** Write comments that help a future developer (including your
future self) understand the intent and constraints that led to the current design.
### Architectural Integrity
- **Respect the Architectural Vision:** Adhere strictly to the decisions documented in
the `docs/` folder. Do not accept "standard" or "idiomatic" patterns that contradict
established project decisions without a compelling reason.
- **Separation of Concerns:** Verify that business logic, persistence, and API concerns
are in the correct layers. Do not leak implementation details across boundaries.
- **DRY (Don't Repeat Yourself):** Actively look for duplicated logic or patterns. Solve
problems at the highest appropriate level (e.g., shared library vs. copy-paste) rather
than patching locally.
- **Holistic Approach:** When reviewing a fix, do not just look at the symptom. Analyze
the system-wide implications. Ask: "Is this the right place for this logic? Does this
create debt elsewhere?"
### Compare against established patterns
- Identify the reference implementation (e.g., a sibling service that's already been
refactored).
- Check that the code under review follows the same conventions: directory structure,
naming, function signatures, file organization.
- Call out misalignments between similar services or modules.
### Check the full picture, not just the code
- **Build configs:** Do build paths (Dockerfiles, CI pipelines, task runners) match the
actual source layout?
- **Dev tooling:** Are dev/build/test targets consistent with the current structure?
- **Ignore files:** Are build artifacts and generated files excluded from version
control?
- **Planning docs:** Does the code match what was planned? Is anything missing?
### Check for stale artifacts
- Flag TODO comments that reference completed work.
- Identify dead code left over from refactoring — duplicate functions, unused imports,
orphaned helpers.
- Check for outdated references in documentation, comments, or config files.
### Verify error propagation
- When code is reorganized, check that errors still flow correctly through the layers.
- Does a failure in one layer still prevent dependent operations from running?
- Do callers receive meaningful error messages with enough context to diagnose the
issue?
### Check for dependency cycles
- Verify that the dependency graph between modules and packages stays clean.
- Flag any pattern where two packages or modules would need to import each other.
- When logic moves between layers, confirm no circular references are introduced.
### Verify the build, not just the tests
- Passing tests do not guarantee the application builds or deploys correctly.
- Check that build pipelines, entry points, and deployment configs reflect any
structural changes.
- Verify that implementations satisfy their interfaces or contracts.
### Organize the review as an outline
- When there's a lot to track, present the status as a structured outline or table.
- Show what's done, what's in progress, and what remains — layer by layer or file by
file.
- Map old code to new locations so nothing gets lost in a refactor.
### Do a line-by-line audit for refactors
- Compare every function, type, and constant from the old code against the new
structure.
- Verify that behavior is preserved (or that changes are intentional improvements).
- Check that tests cover the same cases as before, plus any new orchestration tests.
