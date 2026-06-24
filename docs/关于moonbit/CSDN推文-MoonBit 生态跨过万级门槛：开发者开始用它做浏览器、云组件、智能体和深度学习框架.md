# MoonBit 生态跨过万级门槛：开发者开始用它做浏览器、云组件、智能体和深度学习框架

CSDN

 _2026年6月18日 12:30_ _江苏_ 22人

![图片](https://mmbiz.qpic.cn/mmbiz_gif/Pn4Sm0RsAujqPrtpN5QqqOuSG0AdUUx8wRg3vvr7aQJicdC0cHDBqZrJBZajQUXx97BwNaAE6L6x7OGlsTIpWjw/640?wx_fmt=gif&wxfrom=13&tp=wxpic#imgIndex=0)

![图片](https://mmbiz.qpic.cn/mmbiz_png/Pn4Sm0RsAujX5kS5KQ6BaBUsy1RqR06QuwjkSP1G6wEJHaJCLTONqlcQexqRgJcIICxofIOJs6B6tWBfibb7now/640?wx_fmt=png&wxfrom=13&tp=wxpic#imgIndex=1)

MoonBit 生态跨过万级门槛

MoonBit 生态跨过万级门槛，意义不只是“包数量多了”。

从全球编程语言发展史看，这其实是一个很严苛的标准。HOPL（Online Historical Encyclopaedia of Programming Languages）收录了 8,945 种 从 18 世纪至今出现过的编程语言，但其中绝大多数并没有形成可持续的现代包生态，更不用说拥有万级规模的 Packages。

对一门新语言来说，语法设计、编译器、benchmark 只能证明“这门语言能跑”。但包生态过万，意味着它已经不只是“能写代码”的新语言，而是开始具备承载真实工程、社区协作和 AI 时代开发范式的生态基础。

这也是 MoonBit 这个节点值得关注的原因。MoonBit 于 2023 年正式对外发布，Mooncakes 于 2024 年 2 月上线，并从一开始就和 moon 构建系统、依赖管理、文档生成流程结合在一起。到发稿前后，Mooncakes Packages 已跨过万级门槛，Modules 超过 1,300 个，Downloads 接近 400 万次。也就是说，MoonBit 大约用了两年多时间，就完成了从语言发布到包生态万级的第一阶段跨越。

![图片](https://mmbiz.qpic.cn/mmbiz_png/S1iaf4GgGjEyySWyggRmsIBNYNGRu6CyAMb19icjksXx0rSX7EFTqjSIEhKvd7SKEHicVubQKol6nzlFCB5ajNwiaDlenDMiaO9NnAgYyTnfA5GY/640?wx_fmt=png&from=appmsg&tp=wxpic&wxfrom=5&wx_lazy=1#imgIndex=2)

如果拿 Rust 作参照，这个速度会更直观。Rust 最早在 2006 年作为个人项目启动，2012 年发布 0.1 版本，到 crates.io 生态跨过万级规模，大致经历了一个较长的早期积累周期。按项目启动口径看，这个周期约十年左右；即便按 2012 年公开发布口径看，也用了数年时间。这个对比不是为了说明谁更强，而是为了说明：包生态过万对任何新语言都不是轻松跨过去的门槛。

当然，MoonBit 也赶上了新的时代窗口。AI coding 的普及、Wasm Component Model 的发展、云端运行时的变化，以及现代包管理和开源协作流程的成熟，都在加速新语言生态的形成。MoonBit 的优势在于，它不是在生态成熟后才补工具链，而是从早期就把编译器、构建系统、包管理、文档、AI 工具和多后端输出放在同一套开发体验里。

总的来说 「生态过万」 对一门新编程语言来说非常关键，因为开发者真正决定是否采用它时，看的不只是语法和性能，还会看有没有库可用、有没有稳定的包管理、有没有持续维护的社区项目，以及能不能把 demo 推进到真实工程。

值得一提的是，MoonBit 生态里出现了几个值得注意的工程样本：日本前端工程师 mizchi 正在用 MoonBit 构建轻量 headless browser environment；Golem Cloud 官方用 MoonBit 实现协作列表编辑器，并将其作为 Wasm Component 运行；一位清华大学计算机博士也开源了一个基于 MoonBit 的类 PyTorch 深度学习训练框架 MoonXi-net；社区开发者 brickfrog 则用 MoonBit 编写了本地 agent orchestrator Choir，尝试把多个 AI coding agents 放进任务拆分、独立分支、PR、CI 和 review feedback 的工程协作流程里。

这些案例覆盖了浏览器基础设施、云端 Wasm Component、深度学习训练框架和 AI agent 工具链。证明了一件事：MoonBit 正在从“语言特性展示”进入“真实项目验证”的阶段。

  

![图片](https://mmbiz.qpic.cn/mmbiz_png/Pn4Sm0RsAujX5kS5KQ6BaBUsy1RqR06QDyP5HTUOXsPlJWd79yygiasaqXicpN7ibIfqiak5WFpaxE1mGxxpfmMjiaA/640?wx_fmt=png&tp=wxpic&wxfrom=5&wx_lazy=1#imgIndex=3)

Crater：用 MoonBit 构建浏览器基础设施

第一个案例来自日本前端工程师 mizchi。他做了一个项目叫 Crater。

Crater 并不是一个完整浏览器，也不是要替代 Chromium 或 WebKit。它更准确的定位，是一个用 MoonBit 写的 headless browser environment，通过 WebDriver BiDi 协议接入 Playwright，让测试可以连接到一个更小、更可预测的浏览器目标，而不是每次都启动完整 Chromium。

这个项目的切入点，是浏览器里非常复杂的一层：布局、DOM、CSS、渲染和测试协议。  

![图片](https://mmbiz.qpic.cn/mmbiz_png/S1iaf4GgGjEwkhYcvSY5gnFubj2fdiaF7NQibv6iaZzsK6cfuflCTc2rWhIIUf0mRqaasbiaYianfYib8xjWPnbGeDQafsrQZ9XvpiczaPUvS6tC3dw/640?wx_fmt=png&from=appmsg&tp=wxpic&wxfrom=5&wx_lazy=1#imgIndex=4)

Crater Playground 展示了用 MoonBit 实现 HTML/CSS 布局计算，并在 Canvas 中预览渲染结果。

mizchi 在文章《Building My Own Browser Layout Engine from Scratch》中写到，自己一直想做浏览器，这次主要实现的是浏览器 GUI 层里的布局坐标计算，也就是 CSS box model、Flex 和 Grid 的坐标计算。他也实现了部分 HTML parser、CSS parser、CSS query engine 和 painting，但核心投入还是布局坐标逻辑。

这比普通 Web demo 更有分量。

浏览器布局引擎不是业务页面。它要处理标准、解析、层叠、布局树、渲染树、边界条件和兼容性测试。Crater 不是在展示 MoonBit 能跑一个网页，而是在展示 MoonBit 能否承载浏览器级复杂基础设施。

如果一门语言只能写小型 Wasm 函数，它的想象空间会受限。但像 Crater 这样的项目出现后，开发者看到的是另一种可能：MoonBit 可以进入解析器、布局引擎、协议层、测试运行时这些更底层的工程场景。

  

![图片](https://mmbiz.qpic.cn/mmbiz_png/Pn4Sm0RsAujX5kS5KQ6BaBUsy1RqR06Q0QtgZAx8xsoTReptvArfwbn9MvHGVfV98Qkl5PMRS6MCt2ljwLJwoQ/640?wx_fmt=png&tp=wxpic&wxfrom=5&wx_lazy=1#imgIndex=5)

Golem Cloud：把 MoonBit Wasm Component 跑成 durable worker

第二个案例来自美国云基础设施公司 Golem Cloud 官方博客《Using MoonBit with Golem Cloud》。

这篇文章是用 MoonBit 实现了一个协作列表编辑器。应用可以同时打开任意数量的列表，多个用户可以连接、编辑、查询变更、归档列表；当列表长时间没有更新时，还可以给编辑者发送邮件提醒。

这个案例拆成了三个 Golem components：list 负责某一个可编辑列表，archive 负责保存归档后的列表，email notifier 负责在列表长时间未更新时发送提醒。

Golem 的设计很适合这个场景。每个 editable list 可以映射到自己的 worker，通过列表 ID 标识。这样 list component 本身只需要处理一个列表，多列表扩展交给 Golem。archive component 则可以作为单例 worker。

文章里还提到，由于 Golem 提供 durable execution guarantee，简单场景下 archive 甚至可以先把状态放在内存里。所谓 durable execution，可以简单理解为：应用执行过程中的状态会被运行时托管和恢复，开发者不需要把所有状态都显式写入数据库后再恢复流程。

![图片](https://mmbiz.qpic.cn/mmbiz_png/S1iaf4GgGjEyZriaQ2B20RaMCRNoelndWSICxZh5aUtjuicdS5LaOydPicy6y3ef67eYtURQXib85iaqmmCgYXicGE6qHU7gT3q2AehN8Gk2MQ0WVA/640?wx_fmt=png&from=appmsg&tp=wxpic&wxfrom=5&wx_lazy=1#imgIndex=6)

这件事的重点不是“MoonBit 又能写一种云函数”。

更准确地说，MoonBit 负责把业务逻辑编译成 Wasm Component，Golem Cloud 负责让这些组件以 durable worker 的形式运行起来。一个负责生成小而清晰的组件，一个负责状态恢复、worker 隔离和组件通信。

这正好击中了 MoonBit 的核心叙事：Wasm-first。

对很多语言来说，Wasm 是后来补上的编译目标。对 MoonBit 来说，Wasm 是一开始就被放在核心位置的方向。Golem Cloud 这样的 runtime，给 MoonBit 提供了一个很自然的落地场景。

  

![图片](https://mmbiz.qpic.cn/mmbiz_png/Pn4Sm0RsAujX5kS5KQ6BaBUsy1RqR06Q3SAo8dibFicWvibnB5u4gdHLgb2AbA6UJ8VImKnyW92hibqZefwpDbPYAQ/640?wx_fmt=png&tp=wxpic&wxfrom=5&wx_lazy=1#imgIndex=7)

MoonXi-net：用 MoonBit 探索类 PyTorch 训练框架

第三个案例来自知乎文章《开源了一个基于 MoonBit 语言的深度学习训练框架：MoonXi-net》。

据作者披露，他基于 MoonBit、国内大模型和开源 harness，在不到一个月的业余时间里，实现了一个类 PyTorch 的深度学习训练框架 MoonXi-net。按作者文章中的实验口径，经过算子融合等优化，MoonXi-net 在其 GPU 训练实验中取得了约 2 倍于 PyTorch 的速度表现，但收敛精度略低于 PyTorch，具体原因可能来自随机种子、kernel 实现等差异，仍有待进一步研究。

MoonXi-net 的定位更接近一次工程实验。它的价值不在于马上替代 PyTorch。PyTorch 背后有成熟的算子体系、自动微分、GPU 生态、编译优化和大量社区经验，一个新项目不可能在短时间里覆盖这些东西。

真正值得注意的是，MoonXi-net 没有简单照搬 PyTorch 的动态风格，而是用 MoonBit 的 trait、泛型、FFI 和 Tagless Final 架构，探索一种更类型安全、更容易扩展到多后端的训练框架设计。

在架构上，MoonXi-net 先定义 Tensor trait，把加减乘、维度、zeros、scale、square、mean、scalar、size 等计算能力抽象出来，再分别实现 CPU 侧的 NpArray 和未来可接入的 GPUTensor。模型层通过泛型约束实现，例如 Linear[T] 的 forward 函数只写一次，只要 T 满足 Tensor 约束，就可以适配不同 Tensor 后端。

自动微分部分使用的是 tape-based autograd。简单说，就是在前向计算时，把反向传播需要执行的闭包记录到 tape 里；调用 backward 时，再反向播放这些记录。扩展算子时，则通过新的 trait 继续拆分能力，例如卷积、ReLU 等图像模型能力可以放到 ImageTensor trait 里。GPU 支持则通过 MoonBit 的 FFI 调用 CUDA kernel、cuBLAS、cuDNN 等算子接口。

![图片](https://mmbiz.qpic.cn/mmbiz_png/S1iaf4GgGjEzDN5qABicGBKl9BHNicjoPslhs7fuOYNrhGY0DqNe9QYQwN98zve8N4ZNLmCCsrK7GtaqUREwpSOLWN723dQwDQlQQbPhDoFia2M/640?wx_fmt=png&from=appmsg&tp=wxpic&wxfrom=5&wx_lazy=1#imgIndex=8)

作者也列出了项目当前的不足：全局 tape 需要手动重置，暂时不支持二阶梯度；Tensor 泛型参数还不能在编译期做精确的维度和形状检查；GPU 训练期间为了加速内存分配使用 pool，这会要求优化器状态和模型参数必须 in-place 更新。

需要说明的是，MoonXi-net 的性能结论来自作者文章中的实验结果，CSDN 未对相关 benchmark 进行独立复现。训练速度和收敛效果通常会受到硬件、数据集、batch size、算子实现、随机种子、CUDA kernel、框架配置等多重因素影响，相关结果应结合具体实验环境理解。

从 MoonBit 生态角度看，MoonXi-net 的意义不只是性能对比，而是它说明 MoonBit 已经被开发者拿来探索 AI 基础设施：Tensor 抽象、自动微分、多后端、GPU FFI、训练循环和性能优化。

  

![图片](https://mmbiz.qpic.cn/mmbiz_png/Pn4Sm0RsAujX5kS5KQ6BaBUsy1RqR06Q0jqAkXB2hEohZkOqYzXVsmJtnNRQncYh54ZQOpcC5ZGSicYFgtBZP4w/640?wx_fmt=png&tp=wxpic&wxfrom=5&wx_lazy=1#imgIndex=9)

Choir：外部开发者用 MoonBit 编排多个 AI coding agents

除了浏览器、云组件和深度学习训练框架，MoonBit 生态里也开始出现面向 AI agent 工具链的外部开发者案例。

比较典型的是 Choir。这是社区开发者 brickfrog 开源的一个本地 agent orchestrator，主体使用 MoonBit 编写。它不是简单的 LLM API wrapper，而是尝试把多个 AI coding agents 放进一个更接近真实软件工程的协作流程里。

Choir 的设计思路是：由一个能力更强的模型担任 team lead，例如 Claude；其他模型或工具作为 leaf agents，例如 Codex、Gemini、Moon Pilot、Cursor Agent 等。team lead 负责拆分任务，leaf agents 分别在独立的 git worktree 中完成实现，再通过 PR 回到 team lead 的分支。项目内置 poller 会继续处理 GitHub review、CI、Copilot review 等反馈，并把结果路由给对应的 agent。

也就是说，Choir 关注的不是“让一个模型直接写完整项目”，而是把任务拆分、并行开发、分支隔离、PR、CI、review feedback 和结果合并 串成一个多 agent 协作闭环。

![图片](https://mmbiz.qpic.cn/mmbiz_png/S1iaf4GgGjEyiazAkvficibMic8BGDjj4LicfwbyicoSec9CiclqmaFqkfibKUy4m51SHbiann6iaqr462ITSLALuRbqNYgTowbFpH7CXvA30diamGCtiatc/640?wx_fmt=png&from=appmsg&tp=wxpic&wxfrom=5&wx_lazy=1#imgIndex=10)

这个案例的意义在于，MoonBit 不只是 AI coding 的目标语言，也开始被外部开发者用于实现 AI coding 基础设施本身。前面 MoonXi-net 展示的是开发者用 AI 辅助完成 MoonBit 工程项目；Choir 则反过来，用 MoonBit 去编写 agent 编排工具。

当然，Choir 目前更适合被看作一个早期工程样本，而不是成熟的商业化 agent 平台。它的价值不在于替代 Claude Code、Cursor 或其他 agent 产品，而在于提供了一个有代表性的外部用户案例：MoonBit 正在被开发者用于构建多 agent 协作和 AI coding 基础设施。

放在整篇文章里，Choir 补上了一个新维度：Crater 展示浏览器基础设施，Golem Cloud 案例展示 Wasm Component 云运行时，MoonXi-net 展示深度学习框架探索，而 Choir 则说明 MoonBit 也开始进入 AI agent 工具链本身。

  

![图片](https://mmbiz.qpic.cn/mmbiz_png/Pn4Sm0RsAujX5kS5KQ6BaBUsy1RqR06QGHX4EaTS3AknUHvnayNiaLnANlsxF3g1JR5q9HIdppdNd8c8uyuC05w/640?wx_fmt=png&tp=wxpic&wxfrom=5&wx_lazy=1#imgIndex=11)

从 AI Coding 到复杂工程：MoonBit 的工具链闭环

把 Crater、Golem Cloud、MoonXi-net 和 Choir 放在一起看，会发现它们背后指向的是同一件事：MoonBit 不只是想做一门能编译到 Wasm 的语言，而是在尝试构建一套适合复杂工程和 AI coding 的工具链。

AI coding 真正需要的，不是模型一次性把代码写对。

现实情况往往是：模型先写一版，编译器打回来；再改一版，测试打回来；继续改，继续跑，直到结构稳定下来。也就是说，AI coding 的核心能力不只是“生成”，而是生成之后能不能快速验证、快速定位、快速修复。

MoonBit 比较有优势的地方，正在这里。

首先是反馈速度。AI Agent 高频试错时，编译和测试反馈越快，迭代效率越高。MoonBit 官方强调过编译速度和小体积 Wasm Component 输出，这类能力对人类开发者有用，对 agent 更有用。

其次是类型约束。AI 写代码很容易写出“看起来像那么回事”的代码，尤其是跨语言时会出现语法幻觉。MoonXi-net 作者也提到，AI 在写 MoonBit 时产生过不少 Rust 语法幻觉，比如 trait 泛型写法、高阶类型写法等，最后都被编译器拒绝。

这不是坏事。对 AI coding 来说，编译器越早拒绝错误，Agent 越容易修正方向。静态类型、trait 约束、错误处理、测试和文档，都是给 agent 的边界。

更值得注意的是，MoonBit 还把“证明”放进了原生工具链。MoonBit 0.9 中引入了 first-class formal verification，moon prove 与 moon build、moon test 一样，都是工具链里的内置命令。开发者可以在 .mbt 文件中编写可执行代码，在 .mbtp 文件中编写谓词、逻辑辅助函数和引理，再通过 moon prove 运行验证，并借助 SMT solver 处理证明义务。

这对 AI coding 很关键。测试能告诉开发者“这些样例有没有通过”，而形式化验证进一步尝试回答“某些性质在逻辑上是否成立”。当代码越来越多由 agent 生成时，单靠自然语言 prompt 和运行样例很难建立足够信任。编译、测试、诊断、证明如果能放在同一套工具链里，模型写出的代码就更容易被机器检查、约束和修正。

从更大的工程视角看，MoonBit 的几个特性也能解释为什么这些案例会出现。

Crater 需要 JS、Wasm、测试环境；Golem Cloud 需要 Wasm Component；MoonXi-net 涉及 CPU、GPU、FFI 和本地训练实验；Choir 则需要 CLI、本地工具、GitHub PR、CI 和多 agent 协作。它们都不是单一运行环境里的小 demo，而是需要语言具备多后端、清晰类型系统、测试能力、包管理、文档生成和工程自动化能力。

这也是 MoonBit 比较特殊的地方。它是在 Wasm、AI coding 和云端运行时已经成为明确趋势之后成长起来的新语言，不用背太多历史包袱，可以把编译器、构建系统、包管理、文档、测试、证明、AI 工具链和多后端输出放在同一套设计里。

所以，MoonBit 讨论 AI Coding 时，重点不是“让模型多写几行代码”，而是把模型放进编译、测试、诊断、证明和修复的工程闭环里。

资料来源：

1. Mooncakes 包管理平台首页

 https://mooncakes.io/

2. mizchi / Crater GitHub 仓库

 https://github.com/mizchi/crater

3. mizchi 文章：Building My Own Browser Layout Engine from Scratch

 https://zenn.dev/mizchi/articles/build-my-own-layout-engine?locale=en

4. Golem Cloud 官方博客：Using MoonBit with Golem Cloud

 https://golem.cloud/blog/using-moonbit-with-golem-cloud/

5. 知乎文章：《开源了一个基于 MoonBit 语言的深度学习训练框架：MoonXi-net》

 https://zhuanlan.zhihu.com/p/2039128564604342839

6. brickfrog / Choir GitHub 仓库

 https://github.com/brickfrog/choir

阅读 2.0万

MoonXi-net

​

**留言 7**

作者已设置关注后可以留言

- ![](https://wx.qlogo.cn/mmopen/PiajxSqBRaELsQc06ZiaiaqXPE3qVXkq6BKjoKMbQCHFMAOy4MFT8eW8Vwicib5B7RDQ1ATicgIPLnIibg1J9WzsNCas548qxkKPmuSGa3XYrB6qN2IicdVUGpc7TSqWGncRuicCv/64)
    
    下小雨
    
    陕西6天前
    
    赞2
    
    有什么优势
    
    ![](https://wx.qlogo.cn/mmopen/S1iaf4GgGjEwKS2Iic7mmHSqFf2Wicsau9BIUIL4bWA3LlfyJM4libhdXJg1OYHGkia6dzc6LHfu4vueibrYurwGoYqPRyoKRAy2AjCGKFbfEmOgRzh4ADFibtE5mgaWQ5dlXvg/64)
    
    Hongbo
    
    广东6天前
    
    赞4
    
    可靠性高，可以验证，性能好，编译快
    
- ![](https://wx.qlogo.cn/mmopen/ajNVdqHZLLC6blplAiaB7OZNpWOMuJ6JaWiagye4iaL0J5HO9nJMC9xfEL41bH4sicKd567KnjR4Ij3a1YjTWJJWkllFfmY2dGsNYtdlib4FfV3LY6Ca0P9M5zgxf0hAMp0a0/64)
    
    强降雨.
    
    山西6天前
    
    赞2
    
    期待的
    
- ![](https://wx.qlogo.cn/finderhead/0nn3FBrD9a1m1UWJhQ2B9jrBAhdlDnI1rHaxWbco9o9paibFUiby0BPzbz7XIoIc1gX91tmOx5Py8/64)
    
    刘国华YZ
    
    湖南6天前
    
    赞2
    
    国产崛起
    
- 
- ![](https://wx.qlogo.cn/mmhead/Q3auHgzwzM7ZOUPAFibGicN3MMDWOMaTiceRz8LvITCOcEt9qiaot73GYw/64)
    
    数字游民开发者
    
    广东6天前
    
    赞2
    
    👍👍👍
    
- ![](https://wx.qlogo.cn/mmopen/S1iaf4GgGjEyq0Csj58Xx65EOHW22OA8jmyibWdtPFYELicPhVUcpMoxAwwwHSicw0icJB1h0kaiamI07705Jp11gd8FdiapOViaLMCp/64)
    
    IQ 37
    
    北京5天前
    
    赞1
    
    18世纪的编程语言是啥？
    
- ![](https://wx.qlogo.cn/mmopen/S1iaf4GgGjEzM2lebptyjUjUosFaGVa3uK46ibPfZiaCdujcXaCgXV9jxNicpQHlv2SrgfHx2YCk1enWvCxESIiaHnicMpOtWlz2rvMd0OTpYiaCnKpDdeboKQnicyscGnRNGpza/64)
    
    老鹰球
    
    浙江6天前
    
    赞1
    
    国内企业呢，在生产里用的怎样
    

已无更多数据

[](javacript:;)

![](https://mmbiz.qpic.cn/sz_mmbiz_png/q7kM1wTg0ehdb83pia0aXf56jvuNIpxqhlD7GibqGAKnTtGibgV2IBuktPbPUhrcCrHERt1yeurQNGwX1fHibSHsAA/300?wx_fmt=png&wxfrom=18)

CSDN

2个朋友关注

关注

52

72

23

7

![](https://wx.qlogo.cn/mmopen/PiajxSqBRaEKPPEsdHy7cRGLwAibElBicw6eOr3UaGIgKzNhGLicUFgibSvyMxyS56tHP5VAdicMf5ZJYq5kBXjtZk5nw3aO415ribyLjbJr2UMlnCCAPh2dy3MHWicpcxKoODhR/96)

复制搜一搜

复制搜一搜

暂无评论