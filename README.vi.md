# stack-agents

Một bộ Markdown rules và workflows dùng được với mọi AI coding agent:
Claude Code, Cursor, Codex, Cline, Roo Code, Windsurf, OpenCode, GitHub
Copilot, hoặc agent bạn tự xây bằng LangChain hay CrewAI.

Ý tưởng: **một nguồn sự thật duy nhất (`AGENTS.md`)**, được một binary tĩnh
(`stack`) render ra file rules riêng của từng tool. Skills load theo nhu cầu,
nên không bao giờ nằm trong context window vô ích.

**Ngôn ngữ:** [English](README.md) | Tiếng Việt

## Cách hoạt động

`AGENTS.md` giữ vòng vận hành, bảng skills và các hard rules. Bạn chỉ sửa file
đó. `stack sync` render nó ra thư mục `.agents/` rồi link hoặc copy tới đúng
file mà tool của bạn đọc — không chạy gì thêm, không có gì lệch nhau.

Skills và playbooks nằm trong `.ai/` và load theo nhu cầu: agent đọc
`.ai/skills/tdd.md` khi `/tdd` được gọi, không trước.

Chi tiết: [docs/vi/how-it-works.md](docs/vi/how-it-works.md).

## Cấu trúc

```
AGENTS.md                       # Nguồn sự thật duy nhất: vòng lặp, bảng skill, hard rules
VERSION                         # Phiên bản hiện tại; nhúng vào binary
.ai/
├── principles.md               # 16 nguyên tắc kỹ thuật (đọc cho việc không tầm thường)
├── skills/                     # Gọi dạng /...
│   ├── architect.md            # Thiết kế trước khi viết code
│   ├── arena.md                # So sánh các phương án trước khi commit
│   ├── red-team.md             # Review đối kháng (alias /swarm)
│   ├── tdd.md                  # Đỏ - xanh - refactor
│   ├── verify.md               # Bằng chứng trước khi "xong" (tự chạy)
│   ├── checkpoint.md           # So sánh trạng thái quanh lệnh rủi ro
│   ├── unslop.md               # Cắt chữ thừa, giữ chất
│   └── conventions.md          # Học thói quen của repo, lưu vào .ai/conventions.md
└── playbooks/                  # Được /auto <goal> chọn
    ├── debug.md                # Có nhánh sự cố production
    └── refactor.md
cmd/stack/                      # CLI: init, status, sync, check, uninstall, version
internal/
├── agents/                     # Registry: 8 agent và các đường dẫn chúng đọc
├── detect/                     # Tìm agent đã cài trên máy này
├── generate/                   # Render AGENTS.md ra từng định dạng file
├── project/                    # Install, sync, check, uninstall cho một project
└── ui/                         # Bảng agent và prompt y/n
embed.go                        # Nhúng AGENTS.md, VERSION, .ai/ để shipping
.agents/                        # Sinh ra; commit để CI kiểm tra drift
scripts/load_skill.py           # Loader cho custom agent
docs/en/                        # Tài liệu gốc (tiếng Anh)
docs/vi/                        # Bản dịch tiếng Việt của docs/en
CHANGELOG.md                    # Mục thay đổi, theo từng phiên bản
CONTRIBUTING.md                 # Workflow và checklist pull-request
Makefile                        # make check | sync | diff | test
```

## Cài đặt

`stack` là một binary đơn. Ba cách để có nó:

```bash
go install github.com/MDC15/stack-agents/cmd/stack@latest
```

hoặc tải bộ nén từ trang GitHub Releases (linux, darwin,
windows × amd64, arm64), hoặc build từ checkout:

```bash
git clone https://github.com/MDC15/stack-agents && cd stack-agents && go build -o stack ./cmd/stack
```

Sau đó, trong bất kỳ project nào:

```bash
stack init --path /path/to/your-project
```

Lệnh này dò các agent đang cài trên máy, hiện kế hoạch, hỏi một lần, rồi ghi
`AGENTS.md`, `VERSION`, `.ai/`, `scripts/load_skill.py` cùng các file sinh ra
trong `.agents/` và link cho từng agent. Không có `--force`, nó từ chối đụng
vào một install đã có sẵn, nên lần hai không thể phá gì. Cập nhật bằng
`--force` sẽ merge `AGENTS.md` mới lên của bạn nhưng giữ nguyên mục
**Project overrides**.

Bản ghi thật (OK), trên máy có OpenCode:

```
╭──────────────┬─────────────────┬─────────────────╮
│ AGENT        │ VERSION         │ STATUS          │
├──────────────┼─────────────────┼─────────────────┤
│ Claude Code  │ —               │ ✓ link          │
│ OpenCode     │ opencode v2.0.2 │ ✓ reads AGENTS. │
...
╰──────────────┴─────────────────┴─────────────────╯
created .agents/claude.md
linked CLAUDE.md -> .agents/claude.md
...
Next: edit AGENTS.md, then run `stack sync --path .`
```

| Tuỳ chọn | Tác dụng |
|---|---|
| `--agents <csv>` | Link đúng danh sách agent ID (claude, opencode, codex, cursor, copilot, cline, roo, windsurf) |
| `--all` | Link mọi agent |
| `--force` | Cập nhật install hiện có, giữ Project overrides |
| `--dry-run` | In kế hoạch, không ghi gì |
| `--yes` | Bỏ prompt xác nhận |
| `--path` | Thư mục project (mặc định `.`) |

> **Đừng sửa các skill file gốc** trong project đã cài — cập nhật `--force`
> sẽ thay thế chúng. Thay vào đó, thêm `.ai/skills/<name>.md` của riêng bạn;
> file thêm vào không bao giờ bị update hay uninstall đụng tới.

Hướng dẫn đầy đủ: [docs/vi/getting-started.md](docs/vi/getting-started.md).

## Các file sinh ra

| Agent | Sinh ra trong `.agents/` | Agent đọc | Cách link |
|---|---|---|---|
| Claude Code | `claude.md` | `CLAUDE.md` | symlink |
| OpenCode | — | `AGENTS.md` | trực tiếp, không file |
| Codex | — | `AGENTS.md` | trực tiếp, không file |
| Cursor | `cursor.mdc` | `.cursor/rules/stack-agents.mdc` | symlink |
| GitHub Copilot | `copilot.md` | `.github/copilot-instructions.md` | copy kèm banner |
| Cline | `cline.md` | `.clinerules` | symlink |
| Roo Code | `roo.md` | `.roorules` | symlink |
| Windsurf | `windsurf.md` | `.windsurfrules` | symlink |

Mặc định "có mặt" là opt-in: chỉ agent được phát hiện mới nhận link, dù thêm
các agent bạn khai báo bằng `--agents` hoặc `--all`. Sau mỗi lần sửa
`AGENTS.md`, chạy lại `stack sync`. `stack check` thoát lỗi và nêu tên file
lỗi thời — CI chạy chính lệnh đó, nên các file đích cụ thể không thể lệch nhau.

Thêm agent mới chỉ là một dòng trong
[`internal/agents/agents.go`](internal/agents/agents.go): xem
[docs/vi/extending.md](docs/vi/extending.md).

## Các lệnh

| Lệnh | Tác dụng |
|---|---|
| `stack init` | Dò agent, hiện kế hoạch, cài đặt và link |
| `stack status` | Trạng thái từng agent: ok / drift / missing / skipped |
| `stack sync` | Render lại `.agents/*` và làm mới các link đã planned |
| `stack check` | Thoát lỗi nêu tên từng file lỗi thời; `--diff` in phần khác biệt |
| `stack uninstall` | Gỡ những gì stack đã cài, giữ lại file của bạn |
| `stack version` | In phiên bản nhúng trong binary |

Exit codes: 0 thành công, 1 drift hoặc lỗi, 2 lỗi cú pháp. Mọi lệnh và cờ:
[docs/vi/cli.md](docs/vi/cli.md).

## Custom agent

Với vòng lặp riêng của bạn (LangChain, CrewAI, v.v.), `scripts/load_skill.py`
cung cấp rules dưới dạng system prompt và skills dưới dạng tool:

```python
from scripts.load_skill import system_prompt, load_skill, list_modules

prompt = system_prompt()   # dùng làm system prompt cho agent
# expose load_skill(name, kind="skills") như một tool để agent load skill theo nhu cầu
```

`load_skill` chỉ nhận tên dạng `[a-z0-9-]` và hai loại `skills`,
`playbooks`, nên agent không thể dùng nó để đọc file ngoài `.ai/`.

Chi tiết: [docs/vi/custom-agents.md](docs/vi/custom-agents.md).

## Mở rộng

Thêm một skill bằng cách tạo `.ai/skills/<name>.md`, thêm một dòng vào bảng
trong `AGENTS.md`, rồi chạy `stack sync`. Playbook tương tự trong
`.ai/playbooks/`. Giữ `AGENTS.md` ngắn; rule dài nên nằm trong `.ai/`, nơi chỉ
load khi cần.

Chi tiết: [docs/vi/extending.md](docs/vi/extending.md).

## Kiểm chứng

```bash
make check     # gofmt, go vet, test, rồi `stack check`
```

Bộ test phủ render và phát hiện drift, hành vi install và uninstall, phát hiện
agent, chính CLI (test subprocess trên binary vừa build), CLI của
`load_skill.py`, mọi link tương đối trong tài liệu, và `docs/en` với `docs/vi`
phải bước cùng nhau. CI chạy lại toàn bộ cùng health check drift, và Ruff trên
loader Python.

## Đóng góp

Workflow, versioning và checklist pull-request:
[CONTRIBUTING.md](CONTRIBUTING.md). Thay đổi gần đây:
[CHANGELOG.md](CHANGELOG.md).

## Giới hạn

- Phát hiện là best-effort: `stack` tìm binary của agent trên PATH,
  một thư mục đánh dấu trong `$HOME`, hoặc một VS Code extension đã cài.
  Agent cài theo cách khác sẽ được báo là không thấy; khai báo thực dụng bằng
  `--agents`.
- Symlink được thay bằng bản copy (kèm banner) khi OS không cho tạo, nên
  cài đặt vẫn đi đến cùng.
- File rules dạng symlink hiển thị kém khi checkout git trên Windows: chúng
  hoặc liền thành file văn bản thường, hoặc cần bật Developer Mode. Bản copy
  ở mục trước là cùng vấn đề.
- Bản văn này là chỉ dẫn viết, không phải cơ chế cưỡng chế. Agent nghe tới đâu
  phụ thuộc model và tool. Phá hoại vẫn cần quyền + sandbox ở phía tool.

## Giấy phép

[MIT](LICENSE).