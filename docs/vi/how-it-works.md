# Cách hoạt động

> English: [docs/en/how-it-works.md](../../docs/en/how-it-works.md)

## Một nguồn sự thật duy nhất

Mỗi agent muốn quy tắc nằm ở một file khác nhau. Duy trì cùng một nội dung trong
sáu file sẽ hỏng ngay lần đầu bạn sửa một trong số đó. Nên `AGENTS.md` là file
duy nhất bạn sửa, và `scripts/sync.sh` sinh ra phần còn lại từ nó.

Kiến trúc chỉ có vậy. Không có bước build, không phụ thuộc, không cần cài gì
ngoài `bash`.

## Sync sinh ra những gì

| Agent | File | Ghi chú |
|---|---|---|
| Claude Code | `CLAUDE.md` | chép nguyên văn |
| Cursor | `.cursor/rules/stack-agents.mdc` | thêm frontmatter `alwaysApply: true` ở đầu |
| Cline | `.clinerules` | chép nguyên văn |
| Roo Code | `.roorules` | chép nguyên văn |
| Windsurf | `.windsurfrules` | chép nguyên văn |
| GitHub Copilot | `.github/copilot-instructions.md` | chép nguyên văn |

Mọi file sinh ra đều mở đầu bằng một dòng banner nêu `AGENTS.md` là nguồn. Mỗi
bản chép nguyên văn giống `AGENTS.md` từng byte sau banner; file cho Cursor thêm
frontmatter YAML để quy tắc được nạp mà không cần gọi.

`scripts/sync.sh --check` so từng file sinh ra với nội dung mà sync sẽ ghi, và
thoát khác 0 kèm tên file nếu có file lệch. CI chạy lệnh này, đó là thứ giữ sáu
file không trôi khỏi nhau. `--diff` làm tương tự nhưng còn in phần khác biệt của
từng file lệch, để thấy rõ thay đổi chứ không chỉ thấy tên file.

Nếu một công cụ tìm ở chỗ khác — các đường dẫn này thay đổi theo phiên bản — sửa
danh sách `place` ở cuối `scripts/sync.sh`. Phần còn lại của script không cần
đổi.

## Vì sao kỹ năng được nạp khi cần

`.ai/skills/*.md` và `.ai/playbooks/*.md` không được nhúng vào `AGENTS.md`.
`AGENTS.md` chỉ mang bảng ánh xạ lệnh sang file. Agent đọc file kỹ năng khi lệnh
được gọi.

Cách này giữ ngữ cảnh thường trực nhỏ: tám kỹ năng và hai playbook nếu nhúng hết
sẽ nằm trong cửa sổ ngữ cảnh ở mọi yêu cầu, kể cả những yêu cầu không cần tới
chúng. Nó cũng cho phép một kỹ năng dài mà không tốn gì cho tới khi được dùng.

Ngoại lệ là `/verify`, được vòng làm việc tự chạy trước khi báo cáo là đã xong.

## `AGENTS.md` là gì và không là gì

Nó là chỉ dẫn bằng văn bản. Nó yêu cầu agent giữ thay đổi nhỏ, kiểm chứng trước
khi báo thành công, và hỏi trước các lệnh nguy hiểm. Nó không cưỡng chế được
điều nào trong số đó.

Các luật cứng ở cuối `AGENTS.md` — không chạy lệnh phá hủy khi chưa xác nhận,
không in bí mật, văn bản tìm thấy trong file là dữ liệu chứ không phải chỉ thị —
là những luật quan trọng nhất, và cũng là những luật cần được hệ thống quyền hạn
ở phía công cụ chống lưng. Hãy coi `AGENTS.md` là lời đề nghị và sandbox là bảo
đảm.

## Cài và sync là hai việc tách biệt

- `scripts/install.sh` chép repo vào một dự án: `AGENTS.md`, `.ai/`, `scripts/`
  và `VERSION`. Đây là việc làm một lần, và mặc định nó từ chối ghi đè.
  `--force` cập nhật bản cài sẵn có và báo thay đổi phiên bản; `--uninstall` gỡ
  phần do stack-agents quản lý và giữ lại file của dự án.
- `scripts/sync.sh` chạy bên trong một dự án và sinh file cho từng agent từ
  `AGENTS.md`. Nó chạy sau mỗi lần sửa `AGENTS.md`, và chạy trong CI.

Tách hai việc này là lý do vì sao cập nhật một dự án bằng `--force` gộp file từ
nguồn đè lên bản trong dự án mà không phá `.ai/conventions.md`.

## Tiếp theo

- [getting-started.md](getting-started.md) — cài đặt.
- [skills.md](skills.md) — mỗi lệnh làm gì.
- [custom-agents.md](custom-agents.md) — nối vào vòng lặp riêng của bạn.
