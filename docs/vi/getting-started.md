# Getting started

> English: [docs/en/getting-started.md](../../docs/en/getting-started.md)

## Điều kiện cần

- `stack` — một binary tĩnh. Cài bằng
  `go install github.com/MDC15/stack-agents/cmd/stack@latest`, từ một gói
  release, hoặc build từ checkout: `go build -o stack ./cmd/stack`.
- `python3`, chỉ dùng cho custom agent (`scripts/load_skill.py`) và một test
  trong repo; bản thân CLI không cần gì thêm.

## Khởi tạo một project

```bash
stack init --path /path/to/your-project
```

Kết quả là bảng tám agent được tìm thấy hoặc bỏ qua trên máy này.

```
╭──────────────┬─────────────────┬─────────────────┬───────────────────────────────╮
│AGENT         │VERSION         │STATUS           │TARGET                         │
├──────────────┼─────────────────┼─────────────────┼───────────────────────────────┤
│Claude Code   │—               │✓ link           │CLAUDE.md                      │
│OpenCode      │opencode v2.0.20│✓ reads AGENTS.md│                               │
│Codex         │—               │✗ skip           │                               │
│Cursor        │—               │✗ skip           │.cursor/rules/stack-agents.mdc │
│GitHub Copilot│—               │✓ link           │.github/copilot-instructions.md│
│Cline         │—               │✗ skip           │.clinerules                    │
│Roo Code      │—               │✗ skip           │.roorules                      │
│Windsurf      │—               │✗ skip           │.windsurfrules                 │
╰──────────────┴─────────────────┴─────────────────┴───────────────────────────────╯
```

Sau đó một lần xác nhận (`--yes` bỏ qua) và quá trình cài đặt: bộ ship được
copy, `.agents/` được render, và mọi đường dẫn cố định của các agent được
planned sẽ hiện ra.

```
Next: edit AGENTS.md, then run `stack sync --path /path/to/your-project`
```

Ở lần chạy thứ hai, khi không có `--force`, nó từ chối với
"already installed" — không có gì bị ghi đè nhầm.

| Tuỳ chọn | Tác dụng |
|---|---|
| `--agents <csv>` | Link những agent này dù detection bỏ sót |
| `--all` | Link mọi agent |
| `--force` | Cập nhật install hiện có; merge `AGENTS.md` mới và giữ mục **Project overrides** của bạn |
| `--dry-run` | In kế hoạch, không ghi gì |
| `--yes` | Bỏ prompt |

> **Đừng sửa các skill file gốc** trong `.ai/skills/` và `.ai/playbooks/` —
> cập nhật `--force` sẽ thay thế chúng. Thay vào đó, thêm riêng
> `.ai/skills/<name>.md` của bạn; file thêm vào vượt qua update và uninstall
> nguyên vẹn.

## Điền Project overrides

`AGENTS.md` kết thúc bằng mục **Project overrides**. Mục này có quyền ưu tiên
so với mọi thứ phía trên, trừ hard rules. Đặt đúng các lệnh thật của project
vào đó:

```
- Build: `npm run build`
- Test: `npm test`
- Lint/typecheck: `npm run lint`
```

`/verify` đọc mục này để quyết định chạy gì. Nếu bỏ trống, agent phải tự đoán
từ `package.json`, `Makefile`, `pyproject.toml`, `Cargo.toml` hoặc cấu hình
CI.

## Sinh lại sau mỗi lần sửa

Sửa `AGENTS.md`, không bao giờ sửa file sinh ra. Sau đó:

```bash
stack sync --path /path/to/your-project
```

Chỉ kiểm tra, không ghi — ví dụ trong CI:

```bash
stack check --path /path/to/your-project          # exit 1, nêu tên file lỗi thời
stack check --diff --path /path/to/your-project   # tương tự, và in phần khác biệt
```

## Kiểm chứng cài đặt

```bash
stack status --path /path/to/your-project
stack check --path /path/to/your-project
```

`status` hiển thị một dòng cho mỗi agent; `check` sẽ in "no drift". Xác nhận
agent nhận config bằng cách yêu cầu một việc nhỏ, ví dụ `/verify`.

Chưa nhận? Xem caveat mới tại [cli.md](cli.md#migration-từ-≤-02) — trước
khi có `stack`, các project được cài bằng bash script; lần `stack sync` đầu
tiên nhận diện thời kỳ đó qua banner cũ và chuyển đổi file.

## Agent của bạn bây giờ đọc gì

Mỗi tool đọc file riêng của mình, sinh từ `.agents/`: `CLAUDE.md` cho Claude
Code, `.cursor/rules/stack-agents.mdc` cho Cursor, `.clinerules` cho Cline,
`.roorules` cho Roo Code, `.windsurfrules` cho Windsurf,
`.github/copilot-instructions.md` cho Copilot, và chính `AGENTS.md` cho Codex
và OpenCode. Bảng đầy đủ, kèm ghi chú, nằm ở
[how-it-works.md](how-it-works.md#cac-file-sinh-ra).

## Kế tiếp

- [how-it-works.md](how-it-works.md) — sync sinh ra loại gì.
- [skills.md](skills.md) — từng lệnh một.
- [extending.md](extending.md) — thêm skill hoặc playbook của bạn.
- [../../README.md](../../README.md) — bản tóm tắt.