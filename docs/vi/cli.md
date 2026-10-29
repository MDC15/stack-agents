# CLI reference

> English: [docs/en/cli.md](../../docs/en/cli.md)

`stack` có sáu lệnh. Exit code thống nhất ở mọi nơi: 0 thành công, 1 drift
hoặc lỗi, 2 lỗi cú pháp (lệnh lạ, cờ sai).

```
usage: stack <command> [args]
commands: --version, check, init, status, sync, uninstall, version
```

## stack init

Biến một thư mục thành project stack-agents: dò agent, hiện kế hoạch dạng
bảng, xác nhận (trừ khi `--yes`), copy bộ ship (`AGENTS.md`, `VERSION`,
`.ai/**`, `scripts/load_skill.py`), render `.agents/*`, rồi link hoặc copy
đường dẫn cố định của từng agent được planned.

| Cờ | Lệnh | Tác dụng |
|---|---|---|
| `--path <dir>` | tất cả | Thư mục project, mặc định `.` |
| `--yes` | init, uninstall | Bỏ prompt xác nhận |
| `--dry-run` | init, sync, uninstall | In kế hoạch, không ghi gì |
| `--force` | init | Cập nhật install sẵn có, merge có giữ Project overrides |
| `--agents <csv>` | init, sync | Link các agent ID này (cộng thêm vào phần dò được) |
| `--all` | init, sync | Link mọi agent |

Detection là "present-as-opt-in": agent tìm thấy trên máy tự động nhận link,
các agent còn lại bị bỏ qua. Agent ID lạ là lỗi (exit 1).

## stack status

Đọc project và báo trạng thái từng agent: `ok`, `drift`, `missing link`, hoặc
`skipped (agent not found)`. Thoát 1 khi thư mục không phải project.

## stack sync

Render lại mọi file `.agents/*` từ `AGENTS.md` hiện tại, rồi tạo hoặc làm mới
entry đường dẫn cố định của từng agent được planned. Symlink do stack sở
hữu được trỏ lại; bản copy do stack sở hữu được ghi lại; entry không có nguồn
từ stack được để yên và được báo cáo. Link của một agent không nằm trong kế
hoạch được giữ nguyên — presence là opt-in. Sửa `AGENTS.md` làm các file sinh
ra lỗi thời cho tới lần sync kế tiếp; CI chạy [check](#stack-check) để bắt
điều đó.

## stack check

So sánh từng file `.agents/*` và entry đường dẫn cố định với thứ mà sync
sẽ ghi ngay lúc này. Mỗi entry lỗi thời in một dòng:

```
out of date: .agents/claude.md
```

`--diff` nối thêm diff kiểu unified cho từng file `.agents/*` lỗi thời: dòng
hiện tại với tiền tố `-`, dòng kỳ vọng với tiền tố `+`. Thoát 1 khi có gì đó
lỗi thời hoặc thiếu; cũng thoát 1 khi link của một agent được phát hiện biến
mất, hoặc khi thư mục không phải project.

## stack uninstall

Gỡ thư mục `.agents/`, các entry đường dẫn cố định là symlink vào `.agents/`
hoặc bản copy mang banner, `VERSION`, `scripts/load_skill.py`, và mọi file
`.ai/**` mà stack đã ship. Giữ lại `AGENTS.md`, `.ai/conventions.md`, và mọi
skill bạn tự thêm. Không có `--yes` thì hỏi trước; `--dry-run` chỉ in danh
sách thay vì xóa.

## stack version

In phiên bản nhúng trong binary, `stack 1.0.0` (trùng `--version`; alias được
giữ để tiện dùng).

## Migration từ ≤ 0.2

Phiên bản 0.2 shipped các bash script (`install.sh`, `sync.sh`, `test.sh`);
chúng không còn ở 1.0. Để chuyển một project đã có:

1. Vẫn chạy bộ ≤ 0.2? Chạy uninstall cũ trước:
   `bash /path/to/stack-agents/scripts/install.sh --uninstall <project>`.
   Bước này không bắt buộc: một lần `stack init` mới sẽ chuyển layout ngay
   tại chỗ.
2. Chạy `stack init --path <project>`. Install mới hoặc `--force` thay thế
   những gì còn lại, và các file rules vẫn còn banner cũ `scripts/sync.sh`
   được chuyển đổi tự động — lần `stack sync` đầu tiên nhận cả hai thế hệ
   banner là do stack sở hữu.
3. Chạy `stack check` và đưa nó vào CI để các file sinh ra không thể lệch nhau.

## Kế tiếp

- [getting-started.md](getting-started.md) — cài đặt và kiểm chứng.
- [how-it-works.md](how-it-works.md) — vì sao là `.agents/` và symlink.
- [extending.md](extending.md) — thêm agent, skill, playbook.