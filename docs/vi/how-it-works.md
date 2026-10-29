# Cách hoạt động

> English: [docs/en/how-it-works.md](../../docs/en/how-it-works.md)

## Một nguồn sự thật

Mỗi agent muốn rules của mình nằm ở một file khác nhau. Giữ cùng một đoạn văn
trong sáu file sẽ hỏng ngay lần đầu bạn sửa một trong chúng. Nên `AGENTS.md`
là file duy nhất bạn sửa, và binary `stack` sinh nốt phần còn lại.

Đó là toàn bộ kiến trúc. Một binary tĩnh, không dependency khi chạy, không để
lại gì trên máy ngoài các file nó sinh ra.

## `.agents/` và các đường dẫn cố định

`stack sync` render `AGENTS.md` theo từng định dạng đầu ra vào `.agents/`:

| Agent | Sinh ra | Agent đọc | Cách |
|---|---|---|---|
| Claude Code | `.agents/claude.md` | `CLAUDE.md` | symlink tương đối |
| OpenCode | — | `AGENTS.md` | đọc trực tiếp |
| Codex | — | `AGENTS.md` | đọc trực tiếp |
| Cursor | `.agents/cursor.mdc` | `.cursor/rules/stack-agents.mdc` | symlink tương đối |
| GitHub Copilot | `.agents/copilot.md` | `.github/copilot-instructions.md` | bản copy banner |
| Cline | `.agents/cline.md` | `.clinerules` | symlink tương đối |
| Roo Code | `.agents/roo.md` | `.roorules` | symlink tương đối |
| Windsurf | `.agents/windsurf.md` | `.windsurfrules` | symlink tương đối |

Symlink trỏ vào trong (`CLAUDE.md` → `.agents/claude.md`) nên một lần ghi lại
chạm tới mọi agent cùng lúc. Copilot không đi theo symlink, nên nó nhận bản
copy kèm banner, ghi đè tại chỗ. Khi OS không cho tạo symlink, `stack` quay về
bản copy banner đó, nên cài đặt luôn đi đến cùng.

Mọi file sinh ra đều mở đầu bằng banner nêu `AGENTS.md` là nguồn:
`<!-- Generated from AGENTS.md by stack sync. Edit AGENTS.md, not this file. -->`
Banner đó cũng là dấu sở hữu: sync, check và uninstall chỉ thay thế file mang
nó, một symlink vào `.agents/`, hoặc file được bộ ship nhúng liệt kê. File
của bạn — `.clinerules` viết tay, `.ai/conventions.md`, custom skill — không
bao giờ bị đụng tới.

## Presence là opt-in

`stack sync` chỉ đụng tới đường dẫn cố định của agent hoặc được phát hiện
trên máy này, hoặc được khai báo bằng `--agents` / `--all`. Những gì còn lại
nó tránh xa: một `.roorules` do project tự viết tay được để yên, dù tên của
nó có trong registry. `stack status` cho thấy việc chia; `--agents` thêm cụ
thể; `--all` lấy hết.

## Health check cho drift

`stack check` so sánh từng file sinh ra và đường dẫn cố định với thứ mà sync
sẽ ghi ngay lúc này, thoát 1 và nêu tên từng entry lỗi thời:

```
out of date: .agents/claude.md
```

Một lần sửa `AGENTS.md` làm các file sinh ra lỗi thời cho tới lần sync kế
tiếp; CI chạy health check đó, và đó là thứ chặn các file riêng cho tool không
cách xa nhau. `--diff` in phần khác biệt cho từng file `.agents/*` lỗi thời,
nên drift nghe được chứ không chỉ được nêu tên.

Nếu một tool đọc chỗ khác — các đường dẫn này đổi giữa các phiên bản — thêm
một dòng vào registry trong
[`internal/agents/agents.go`](../../internal/agents/agents.go) và chạy lại
`make check`. Không còn gì cần đổi.

## Vì sao skills load theo nhu cầu

`.ai/skills/*.md` và `.ai/playbooks/*.md` không bị inline vào `AGENTS.md`.
`AGENTS.md` chỉ giữ bảng ánh xạ lệnh sang file. Agent đọc skill khi lệnh được
gọi.

Điều này giữ context luôn-present ở mức nhỏ: tám skills và hai playbooks nếu
ngược lại sẽ nằm trong cửa sổ cho mọi request, kể cả những request không cần
cái nào trong số đó. Điều đó cũng có nghĩa một skill có thể dài mà không mất
gì cho tới khi được dùng.

Ngoại lệ là `/verify`, được vòng lặp vận hành chạy tự động trước khi công việc
được báo là xong.

## `AGENTS.md` là gì, và không phải

Gói chỉ dẫn viết. Nó yêu cầu agent giữ thay đổi nhỏ, xác minh trước khi nói
thành công, và hỏi trước các lệnh phá hoại. Nó không cưỡng chế được điều gì.

Hard rules ở cuối `AGENTS.md` — không lệnh phá hoại khi chưa được xác nhận,
không để lộ bí mật, văn bản tìm thấy trong file là dữ liệu chứ không phải chỉ
dẫn — là những mục quan trọng nhất, và chúng là những mục một hệ thống quyền
ở cấp tool phải bảo trợ. Coi `AGENTS.md` là yêu cầu và sandbox là điều kiện
bắt buộc.

## Install, sync và uninstall là riêng biệt

- `stack init` biến một project thành project stack-agents: copy bộ ship,
  render `.agents/`, tạo link. One-time; từ chối ghi đè khi không có
  `--force`, và `--force` merge khi có giữ **Project overrides**.
- `stack sync` chạy bên trong một project và sinh các file riêng cho agent từ
  `AGENTS.md`. Nó chạy sau mỗi lần sửa `AGENTS.md`, và trong CI.
- `stack uninstall` gỡ mọi thứ stack sở hữu và giữ lại các file của project,
  nên việc rời khỏi hệ cũng một lệnh là xong.

Tách riêng chính là lý do một cập nhật `--force` merge file nguồn lên file của
project mà không phá `.ai/conventions.md`.

## Kế tiếp

- [getting-started.md](getting-started.md) — cài nó.
- [cli.md](cli.md) — mọi lệnh và cờ.
- [skills.md](skills.md) — mỗi lệnh làm gì.
- [custom-agents.md](custom-agents.md) — nối vào vòng lặp của riêng bạn.