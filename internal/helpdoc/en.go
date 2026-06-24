package helpdoc

func englishDocs() map[string]CommandDoc {
	return map[string]CommandDoc{
		"seekmoon": {
			Key:     "seekmoon",
			ShortEN: "MoonBit package discovery workbench",
			LongEN: `SeekMoon is a MoonBit package discovery workbench. It helps dependency consumers discover candidate packages, inspect evidence, run local validation, record adoption judgments, and render investigation reports before adding a dependency.

Start here when using SeekMoon for the first time. Start with a command's help before using that command for the first time; command help explains the action, inputs, evidence boundary, output modes, and numbered candidate inputs.

Common path:
  doctor -> sync -> search -> view/api/source/compare -> probe -> record -> report

search and skill search write numbered candidates into the current project's default session. Later commands can use numbers such as 1 or 2. When a number is unavailable, run search again or pass a full coordinate.

Pretty output is for terminal reading. Use --json or --jq for automation, --shape for field paths, and --schema for strict validation.`,
			ExampleEN: `seekmoon doctor
seekmoon sync
seekmoon search markdown --target js
seekmoon view 1
seekmoon api 1 --package moonbitlang/core/argparse
seekmoon probe 1 --target js
seekmoon record 1 --conclusion continue-verification
seekmoon report --format markdown`,
		},
		"seekmoon doctor": {
			Key:     "seekmoon doctor",
			ShortEN: "Check the local environment",
			LongEN: `doctor checks the local environment that SeekMoon commands depend on. It reads MoonBit toolchain availability, registry paths, network reachability, and the current project context.

doctor does not create a snapshot, update the registry, or write adoption records. Its output tells later commands which local capabilities are available.

Run doctor first in a project when using SeekMoon there for the first time.`,
			ExampleEN: `seekmoon doctor
seekmoon doctor --json
seekmoon doctor --shape`,
		},
		"seekmoon sync": {
			Key:     "seekmoon sync",
			ShortEN: "Create a data snapshot",
			LongEN: `sync creates a dated evidence snapshot. The snapshot fixes the current Mooncakes API data, statistics, local registry index summary, and toolchain context used by later investigation steps.

search, record, and report can refer back to the snapshot so the investigation keeps a stable data basis.

sync records source states. When a source fails, the failure belongs to that source action and is reported through the snapshot or error surface.`,
			ExampleEN: `seekmoon sync
seekmoon sync --json
seekmoon sync --schema`,
		},
		"seekmoon search": {
			Key:     "seekmoon search",
			ShortEN: "Search library module candidates",
			LongEN: `search turns a query into library module candidates. It reads module summaries from the current snapshot and matches module name, description, keywords, and repository declaration. If no usable snapshot exists, search can read current API data for a transient search basis.

search writes visible candidates into the current project's default session. The numbers in the output table can be used by view, api, source, compare, probe, and record.

--target adds target context to the candidate list. Target evidence appears only when it is loaded from a source or derived by a later command.

After search, inspect a candidate with view, compare multiple candidates, or run a local probe.`,
			ExampleEN: `seekmoon search markdown
seekmoon search markdown --target js
seekmoon search markdown --json
seekmoon search markdown --jq '.results[].module'
seekmoon view 1
seekmoon compare 1 2`,
			Flags: map[string]FlagDoc{
				"target": {UsageEN: "add target context to the candidate list"},
			},
		},
		"seekmoon view": {
			Key:     "seekmoon view",
			ShortEN: "View a library module profile",
			LongEN: `view shows the evidence profile for one library module. The input can be a full module coordinate or a numbered candidate produced by search.

view reads the Manifest API and the module index asset for the manifest version. The profile includes description, license, repository declaration, downloads, build status, docs URL, package index state, and package summary.

view does not expand full API symbols. Use api for package API details and source for published source material.`,
			ExampleEN: `seekmoon view mizchi/markdown
seekmoon view 1
seekmoon view 1 --json
seekmoon api 1 --package mizchi/markdown/src/api
seekmoon source 1`,
		},
		"seekmoon api": {
			Key:     "seekmoon api",
			ShortEN: "View a package API profile",
			LongEN: `api shows the API profile for one package. The input is a module coordinate or numbered candidate, plus --package <path>.

api reads the module index, resolves the package path to a package relpath, and then reads package_data.json for the same module version. The output contains types, values, traits, docstrings, signatures, and source locations.

When the package path is absent, the error surface uses known package paths from the module index to support recovery.`,
			ExampleEN: `seekmoon api mizchi/markdown --package mizchi/markdown/src/api
seekmoon api 1 --package mizchi/markdown/src/api
seekmoon api 1 --package mizchi/markdown/src/api --jq '.values[].name'
seekmoon api --shape`,
			Flags: map[string]FlagDoc{
				"package": {UsageEN: "package path inside the selected module"},
			},
		},
		"seekmoon source": {
			Key:     "seekmoon source",
			ShortEN: "Locate published source",
			LongEN: `source locates source material for the registry-published module version. The input can be a module coordinate, a versioned coordinate, or a numbered candidate from search.

source records each resolution attempt, including moon fetch, source zip, local cache, core local source, and repository signal. The selected source is derived from successful attempts.

Source resolution can produce local filesystem material under SeekMoon-controlled source or cache paths.`,
			ExampleEN: `seekmoon source mizchi/markdown@0.6.2
seekmoon source 1
seekmoon source 1 --json
seekmoon source 1 --jq '.selected_source.path'`,
		},
		"seekmoon skill": {
			Key:     "seekmoon skill",
			ShortEN: "Search or view executable skill entries",
			LongEN: `skill works with executable skill entries from the Mooncakes Skills API. A skill entry is an executable Wasm or runwasm object.

skill search creates numbered skill candidates in the current project's default session. skill view shows the skill profile, asset states, and pinned runwasm coordinate.

Use record --kind skill when recording a skill investigation result.`,
			ExampleEN: `seekmoon skill search cowsay
seekmoon skill view 1
seekmoon record 1 --kind skill --conclusion continue-verification`,
		},
		"seekmoon skill search": {
			Key:     "seekmoon skill search",
			ShortEN: "Search skill entries",
			LongEN: `skill search turns a query into executable skill candidates. It reads the Skills API and matches skill name, module, package, and metadata description.

The output numbers are written into the current project's default session. Use skill view to inspect a candidate and record --kind skill to save an investigation result.`,
			ExampleEN: `seekmoon skill search cowsay
seekmoon skill search cowsay --json
seekmoon skill view 1`,
		},
		"seekmoon skill view": {
			Key:     "seekmoon skill view",
			ShortEN: "View a skill profile",
			LongEN: `skill view shows the evidence profile for one executable skill. The input can be a skill entry coordinate or a numbered candidate produced by skill search.

skill view reads skill detail, SKILL.md, Wasm asset state, checksum asset state, and derives the pinned runwasm coordinate.`,
			ExampleEN: `seekmoon skill view Yoorkin/cowsay
seekmoon skill view 1
seekmoon skill view 1 --jq '.runwasm_coordinate'`,
		},
		"seekmoon compare": {
			Key:     "seekmoon compare",
			ShortEN: "Compare candidate evidence",
			LongEN: `compare places multiple candidates on one evidence surface. Inputs can be numbered candidates or module coordinates.

compare aligns evidence fields such as manifest data, package index state, source state, probe state, and loaded repository signals. The output shows evidence differences for further inspection, validation, or recording.

After compare, inspect a specific difference with view, api, source, or probe.`,
			ExampleEN: `seekmoon compare 1 2
seekmoon compare 1 2 3 --json
seekmoon compare mizchi/markdown moonbit-community/cmark
seekmoon probe 1 --target js`,
		},
		"seekmoon probe": {
			Key:     "seekmoon probe",
			ShortEN: "Run local validation for a candidate",
			LongEN: `probe records local validation evidence for one candidate in the current toolchain, version, target, and command sequence. The input can be a numbered candidate or module coordinate.

The default probe creates an isolated validation project and runs moon add, moon check, moon test, and target check/build steps. Each step records command, cwd, exit code, status, and log path.

A probe result is local derived evidence. Upstream source fields keep their own source-evidence identity.`,
			ExampleEN: `seekmoon probe 1
seekmoon probe 1 --target js
seekmoon probe mizchi/markdown@0.6.2 --target wasm-gc
seekmoon probe 1 --json`,
			Flags: map[string]FlagDoc{
				"target": {UsageEN: "target backend for local validation steps"},
			},
		},
		"seekmoon record": {
			Key:     "seekmoon record",
			ShortEN: "Save an adoption judgment",
			LongEN: `record saves one adoption judgment. The input is a numbered candidate or coordinate, plus --conclusion <value>.

record writes candidate, version, project context, snapshot, evidence references, conclusion, and note. Conclusions use stable enum values, and JSON output keeps those enum values in English.

Library candidates use --kind library by default. Use --kind skill when recording a skill investigation.`,
			ExampleEN: `seekmoon record 1 --conclusion continue-verification
seekmoon record 1 --conclusion reject-for-now --note "API coverage does not match this project."
seekmoon record 1 --kind skill --conclusion continue-verification
seekmoon record 1 --json`,
			Flags: map[string]FlagDoc{
				"conclusion": {UsageEN: "adoption conclusion enum value"},
				"note":       {UsageEN: "human note stored with the record"},
				"kind":       {UsageEN: "candidate kind recorded by this judgment"},
			},
		},
		"seekmoon report": {
			Key:     "seekmoon report",
			ShortEN: "Render an investigation report",
			LongEN: `report renders an investigation report from existing records, snapshots, project context, and evidence references. It is the document projection of the recorded investigation path.

A report lists sources that were recorded or referenced. Validation steps that were not run do not appear as validation results.

--format selects the report format.`,
			ExampleEN: `seekmoon report --format markdown
seekmoon report --format json
seekmoon report --format markdown --json`,
			Flags: map[string]FlagDoc{
				"format": {UsageEN: "report format to render"},
			},
		},
		"seekmoon raw": {
			Key:     "seekmoon raw",
			ShortEN: "Read a raw source payload",
			LongEN: `raw reads the requested source payload without normalizing its upstream field names or shape. The output includes source status and metadata.

raw serves source audit, field inspection, and failure reproduction. Use search, view, api, source, and skill for the ordinary discovery path.`,
			ExampleEN: `seekmoon raw modules
seekmoon raw manifest mizchi/markdown
seekmoon raw module-index mizchi/markdown@0.6.2
seekmoon raw package-data mizchi/markdown@0.6.2 mizchi/markdown/src/api
seekmoon raw skills`,
		},
	}
}

// CommonFlagDocs returns shared output-mode flag documentation.
func CommonFlagDocs() map[string]FlagDoc {
	return map[string]FlagDoc{
		"json": {
			UsageEN:  "render the command JSON projection for scripts and automation",
			ReviewZH: "输出当前命令结果的 JSON 投影。JSON 服务脚本、CI 和报告生成，不服务终端排版。",
		},
		"jq": {
			UsageEN:  "evaluate a jq expression against the command JSON projection",
			ReviewZH: "先生成当前命令的 JSON 投影，再用内置 jq 求值表达式。",
		},
		"shape": {
			UsageEN:  "show the command JSON field tree without running the data action",
			ReviewZH: "输出当前命令 JSON 投影的字段树。该模式用于写 jq、理解字段路径和学习输出契约。",
		},
		"schema": {
			UsageEN:  "show the command JSON Schema without running the data action",
			ReviewZH: "输出当前命令 JSON 投影的 JSON Schema。该模式用于严格校验和代码生成。",
		},
	}
}
