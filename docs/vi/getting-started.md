# Bắt đầu

> English: [docs/en/getting-started.md](../../docs/en/getting-started.md)

## Cần có trước

- `bash` phiên bản 4 trở lên.
- `python3`, cho agent tự viết (`scripts/load_skill.py`) và cho bộ test.
- `make`, không bắt buộc — chỉ để dùng các lệnh tắt `make`.

Không cần gì thêm. Repo không có phụ thuộc lúc chạy.

## Cài đặt

Một lệnh:

```bash
bash /đường/dẫn/tới/pstack-universal/scripts/install.sh /đường/dẫn/tới/dự-án
```

Script chép `AGENTS.md`, `.ai/` và `scripts/` vào dự án rồi chạy
`scripts/sync.sh` tại đó, nên file quy tắc cho từng agent được sinh trong cùng
một bước.

| Tùy chọn | Tác dụng |
|---|---|
| `--force` | Gộp `AGENTS.md`, `.ai/` và `scripts/` từ nguồn đè lên bản trong dự án, rồi chạy lại sync |
| `--dry-run` | In ra việc sẽ làm, không thay đổi gì |
| `--no-sync` | Chép file nhưng không sinh file quy tắc cho agent |

Không có `--force`, script từ chối đụng vào bản cài đã có, nên chạy lần hai vào
cùng dự án là an toàn và không làm gì cả. Khi có `--force`, file trong nguồn
thay thế bản tương ứng và mọi thứ khác giữ nguyên; `.ai/conventions.md` không có
trong nguồn nên vẫn còn sau khi cập nhật.

<details>
<summary>Cài bằng tay</summary>

1. Chép `AGENTS.md`, `.ai/` và `scripts/` vào thư mục gốc dự án.
2. Chạy `bash scripts/sync.sh`.

</details>

## Điền mục Project overrides

`AGENTS.md` kết thúc bằng mục **Project overrides**. Mục này ưu tiên hơn mọi thứ
ở trên, trừ các luật cứng. Điền lệnh thật của dự án vào đó:

```
- Build: `npm run build`
- Test: `npm test`
- Lint/typecheck: `npm run lint`
```

`/verify` đọc mục này để biết cần chạy gì. Nếu để trống, agent phải đoán từ
`package.json`, `Makefile`, `pyproject.toml`, `Cargo.toml` hoặc cấu hình CI.

## Sinh lại sau mỗi lần sửa

Sửa `AGENTS.md`, đừng bao giờ sửa file đã sinh. Rồi chạy:

```bash
bash scripts/sync.sh
```

Kiểm tra mà không ghi — ví dụ trong CI:

```bash
bash scripts/sync.sh --check     # thoát mã 1 và nêu tên file đã lệch
```

## Kiểm chứng bản cài

```bash
bash scripts/test.sh
```

Lệnh này phủ hành vi của `install.sh` và CLI của `load_skill.py`. Bản cài trong
dự án không có `docs/`, nên phần kiểm tra liên kết được bỏ qua thay vì báo lỗi.

Khi làm việc trong chính repo pstack-universal, `make check` chạy cùng bộ test đó
cộng thêm phần kiểm tra file sinh ra có lệch không.

## Agent của bạn giờ đọc gì

Mỗi công cụ đọc file riêng do sync sinh ra: `CLAUDE.md` cho Claude Code,
`.cursor/rules/pstack.mdc` cho Cursor, `.clinerules` cho Cline, `.roorules` cho
Roo Code, `.windsurfrules` cho Windsurf, `.github/copilot-instructions.md` cho
Copilot, và chính `AGENTS.md` cho Codex cùng mọi agent đọc file đó. Bảng đầy đủ,
kèm ghi chú, nằm trong
[how-it-works.md](how-it-works.md).

Xác nhận đã ăn bằng cách yêu cầu một việc nhỏ, ví dụ `/verify`.

## Tiếp theo

- [how-it-works.md](how-it-works.md) — sync thực sự sinh ra những gì.
- [skills.md](skills.md) — từng lệnh một.
- [extending.md](extending.md) — thêm kỹ năng hoặc playbook của bạn.
- [../../README.vi.md](../../README.vi.md) — bản ngắn.
