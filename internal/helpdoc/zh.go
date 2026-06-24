package helpdoc

func chineseReviewDocs() map[string]CommandDoc {
	return map[string]CommandDoc{
		"seekmoon": {
			Key: "seekmoon",
			ReviewZH: `SeekMoon 是 MoonBit 包发现工作台。它帮助依赖消费者在引入包之前完成候选发现、证据下钻、本地验证、采纳记录和报告输出。

首次使用 SeekMoon 时，从顶层 help 进入。首次使用某个命令时，从该命令 help 进入；命令 help 说明动作、输入、证据边界、输出模式和编号候选输入。

常用工作流：先运行 doctor 检查本地环境，再运行 sync 固定数据口径，然后用 search 生成候选。候选可以继续进入 view、api、source、compare 和 probe。判断结果用 record 保存，用 report 输出。

search 和 skill search 会写入编号候选。后续命令可以用 1、2 这样的编号引用当前项目默认 session 中的候选；编号不可用时，重新运行 search 或传完整坐标。

默认输出服务终端阅读。自动化使用 --json 或 --jq。字段学习使用 --shape。严格校验使用 --schema。`,
		},
		"seekmoon doctor": {
			Key: "seekmoon doctor",
			ReviewZH: `doctor 检查 SeekMoon 当前运行所需的本地环境。它读取 MoonBit 工具链、registry 路径、网络可达性和当前项目上下文。

doctor 不创建 snapshot，不更新 registry，不写 record。它的输出用于确认后续命令是否具备本地运行条件。

首次在一个项目中使用 SeekMoon 时，先运行 doctor。环境状态异常时，先根据错误表面处理本地工具链、路径或网络问题。`,
		},
		"seekmoon sync": {
			Key: "seekmoon sync",
			ReviewZH: `sync 创建带时间戳的数据 snapshot。Snapshot 固定当前 Mooncakes API、统计信息、本地 index 和工具链信息的数据口径。

后续 search、record 和 report 可以引用 snapshot，使调查结果回到同一读取时间和来源状态。

sync 可以执行本地 registry 更新，并把每个来源的状态写入 snapshot。部分来源失败时，失败状态保留在 snapshot 或错误表面中。`,
		},
		"seekmoon search": {
			Key: "seekmoon search",
			ReviewZH: `search 从 query 生成 library module 候选。它读取 snapshot 中的 Modules API 数据，并按 module、description、keywords 和 repository declaration 进行本地匹配。没有可用 snapshot 时，命令可以读取当前 API 数据形成临时口径。

search 会把可见候选写入当前项目默认 session。输出表中的编号可以被 view、api、source、compare、probe 和 record 继续使用。

--target 将 target 上下文加入候选列表。target 证据只来自已读取来源或后续命令的派生结果。

下一步通常是 view 1 查看单候选证据，compare 1 2 比较多个候选，或 probe 1 --target js 做本地验证。`,
			Flags: map[string]FlagDoc{
				"target": {ReviewZH: "将 target 上下文加入候选列表。"},
			},
		},
		"seekmoon view": {
			Key: "seekmoon view",
			ReviewZH: `view 展示单个 library module 的证据画像。输入可以是完整 module coordinate，也可以是 search 产生的编号。

view 读取 Manifest API，并按 manifest version 读取 module index asset。输出包含 description、license、repository declaration、downloads、build status、docs URL、package index 状态和 package 摘要。

view 不展开完整 API symbol。API 详情由 api 命令读取。需要源码材料时进入 source。`,
		},
		"seekmoon api": {
			Key: "seekmoon api",
			ReviewZH: `api 展示某个 package 的 API profile。输入是 module coordinate 或候选编号，并且必须提供 --package <path>。

api 先读取 module index，确认 package path 并派生 package relpath，然后读取同版本的 package_data.json。输出包含类型、函数、trait、docstring、signature 和 source location。

package path 不存在时，错误表面使用 module index 中的已知 package paths 帮助恢复。`,
			Flags: map[string]FlagDoc{
				"package": {ReviewZH: "所选 module 内部的 package path。"},
			},
		},
		"seekmoon source": {
			Key: "seekmoon source",
			ReviewZH: `source 定位 registry 发布版本对应的源码材料。输入可以是 module coordinate、带版本的 coordinate，或 search 产生的编号。

source 记录每次来源解析尝试，包括 moon fetch、source zip、本地 cache、core 本地源码和 repository signal。选中的源码来源来自成功的解析尝试。

源码定位会产生本地文件系统结果。默认结果属于 SeekMoon 控制的 source 或 cache 边界。`,
		},
		"seekmoon skill": {
			Key: "seekmoon skill",
			ReviewZH: `skill 处理 Mooncakes Skills API 中的 executable skill entry。Skill entry 是可执行 Wasm 或 runwasm 对象。

skill search 生成 skill 候选，并写入当前项目默认 session。skill view 展示 skill profile、asset 状态和 pinned runwasm coordinate。

记录 skill 判断时使用 record --kind skill。`,
		},
		"seekmoon skill search": {
			Key: "seekmoon skill search",
			ReviewZH: `skill search 从 query 生成 executable skill 候选。它读取 Skills API，并按 skill name、module、package 和 metadata description 匹配。

输出编号写入当前项目默认 session。后续可以用 skill view 1 读取 skill profile，也可以用 record 1 --kind skill 保存判断。`,
		},
		"seekmoon skill view": {
			Key: "seekmoon skill view",
			ReviewZH: `skill view 展示一个 executable skill 的证据画像。输入可以是 skill entry coordinate 或 skill search 产生的编号。

skill view 读取 skill detail、SKILL.md、Wasm asset、checksum asset，并派生 pinned runwasm coordinate。`,
		},
		"seekmoon compare": {
			Key: "seekmoon compare",
			ReviewZH: `compare 把多个候选放在同一个证据表面中。输入可以是多个编号，也可以是多个 module coordinate。

compare 对齐 manifest、package index、source、probe 和已加载 repository signal 等证据字段。它展示证据差异，供消费者继续下钻、验证或记录判断。

下一步通常是对差异明显的候选运行 view、api、source 或 probe。`,
		},
		"seekmoon probe": {
			Key: "seekmoon probe",
			ReviewZH: `probe 记录一个候选在当前工具链、版本、target 与命令序列下的本地验证证据。输入可以是候选编号或 module coordinate。

默认 probe 在隔离目录中创建验证项目，执行 moon add、moon check、moon test 和 target check/build。每个步骤记录 command、cwd、exit code、状态和 log path。

probe 结果属于 local derived evidence。上游来源字段保持各自的来源证据身份。`,
			Flags: map[string]FlagDoc{
				"target": {ReviewZH: "本地验证步骤使用的 target backend。"},
			},
		},
		"seekmoon record": {
			Key: "seekmoon record",
			ReviewZH: `record 保存一次采纳判断。输入是候选编号或 coordinate，并且必须提供 --conclusion <value>。

record 写入候选、版本、项目上下文、snapshot、证据引用、结论和 note。结论使用稳定枚举，JSON 输出保持英文枚举值。

library 候选默认使用 --kind library。skill 判断使用 --kind skill。`,
			Flags: map[string]FlagDoc{
				"conclusion": {ReviewZH: "采纳判断枚举值。"},
				"note":       {ReviewZH: "随 record 保存的人类备注。"},
				"kind":       {ReviewZH: "本次判断记录的候选类型。"},
			},
		},
		"seekmoon report": {
			Key: "seekmoon report",
			ReviewZH: `report 从已有 records、snapshot、项目上下文和证据引用生成调查报告。它输出调查轨迹的文档投影。

Report 只列已经被记录或引用的来源。没有执行的验证动作不会出现在报告中。

--format 指定报告格式。`,
			Flags: map[string]FlagDoc{
				"format": {ReviewZH: "需要输出的报告格式。"},
			},
		},
		"seekmoon raw": {
			Key: "seekmoon raw",
			ReviewZH: `raw 读取指定来源的原始 payload。它保留上游字段名和原始 shape，并附带 source status 和 metadata。

raw 服务来源审计、字段复查和失败复现。普通 discovery 路径优先使用 search、view、api、source 和 skill。`,
		},
	}
}
