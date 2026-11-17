---
id: 977cce49e1b7
kind: refactor
title: A1 Expose minimal execution handler adapter for owned Bash# engine
seq: 2
status: todo
priority: p0
created: 2026-11-17T23:01:05.589361Z
sprint: 575
sprint_id: 65ce85f5-87c0-5863-bf55-62de86f5df36
sprint_title: 'Layer separation: Classic <- Bash# <- Yoke import-graph delineation'
---

Sprint575 A1 follows accepted A0 source relocation. Implement ONLY the minimum public handler adapter boundary described in umbrella docs/s575-runtime-boundary.md, independently reviewed by Codex/Muse. Add interp.WithHandlerContext(ctx,hc), AdaptHandlerContext(hc,kind,hooks), public HandlerKind constants Exec/Call/Open/ReadDir, six HandlerHooks: DryRun func()bool, Async func()bool, Umask func()os.FileMode, SignalIgnored func(string)bool, Builtin func(context.Context,[]string)error, DescribeValue func(string)(ValueDescription,bool). Preserve existing public fields and original private Runner paths when unadapted; adapted views consult live hooks. Map public kind explicitly onto existing private enum (exec starts1), preserve exec-only Builtin restriction and nil-hook defaults. No parser/Runner API expansion or behavior changes. Do not modify defaults, stat semantics, private markers, engine code, CLI, module pins, tests unrelated to adapter. Meaningful focused tests must verify callback liveness, all6 hook forwarding, kind restriction, ordinary Runner fallback, context roundtrip incl execflags/stdio/env/invalidpos, no fake defaults using foreign Runner. No builds/tests on dev box; static formatting only. Official-image tests owned by manager. No native foreign suites, dependency upgrades, skip/assertion weakening, KB.md, pushes or root umbrella edits. Commit isolated workspace with proper Sprint/Story/Story-ID and actual agent. Blocking prerequisite A0 must be accepted before assignment.

Manager readiness: A0 is accepted and published at adb3b4e0c28d85586caefd3837ab05e81f108b5e. Start from clean synced main. Review the current runtime-boundary design, not historical seam notes. For adapted nil hooks, provide safe neutral results; preserve the ordinary unadapted behavior. This slice exposes only the six callbacks and context helpers; no engine copy or stripping. The manager runs all build/test gates in the locked image with the recorded CPU profile.
