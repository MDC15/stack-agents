# stack-agents

Bộ quy tắc và quy trình làm việc dạng Markdown, dùng được với mọi AI coding
agent: Claude Code, Cursor, Codex, Cline, Roo Code, Windsurf, GitHub Copilot,
hoặc agent tự viết bằng LangChain / CrewAI.

Nguyên lý: **một nguồn duy nhất (`AGENTS.md`)**, script sinh ra file quy tắc
riêng cho từng công cụ. Kỹ năng chỉ được đọc khi cần, nên không chiếm ngữ cảnh
khi chưa dùng tới.

**Ngôn ngữ:** [English](README.md) | Tiếng Việt

## Cách hoạt động

`AGENTS.md` chứa vòng làm việc, bảng kỹ năng và các luật cứng. Bạn chỉ sửa file
này. `scripts/sync.sh` chép nội dung của nó vào mọi file quy tắc mà từng agent
đọc, nên chúng không thể lệch nhau.

Kỹ năng và playbook nằm trong `.ai/` và được nạp khi cần: agent đọc
`.ai/skills/tdd.md` khi `/tdd` được gọi, không đọc trước.

Chi tiết: [docs/vi/how-it-works.md](docs/vi/how-it-works.md).

## Cấu trúc

```
AGENTS.md                  # Nguồn sự thật: vòng làm việc, bảng kỹ năng, luật cứng
VERSION                    # Phiên bản hiện tại; được chép vào mọi bản cài
.ai/
├── principles.md          # 16 nguyên tắc kỹ thuật (đọc khi làm việc không tầm thường)
├── skills/                # Gọi bằng lệnh /...
│   ├── architect.md       # Thiết kế trước khi viết code
│   ├── arena.md           # So sánh các hướng giải quyết
│   ├── red-team.md        # Phản biện đối kháng (bí danh /swarm)
│   ├── tdd.md             # Đỏ - xanh - refactor
│   ├── verify.md          # Bằng chứng trước khi báo "xong" (tự chạy)
│   ├── checkpoint.md      # So sánh trạng thái quanh lệnh rủi ro
│   ├── unslop.md          # Cắt chữ thừa
│   └── conventions.md     # Học quy ước của repo, lưu vào .ai/conventions.md
└── playbooks/             # Chọn bằng /auto <mục tiêu>
    ├── debug.md           # Gồm nhánh xử lý sự cố production
    └── refactor.md
scripts/
├── install.sh             # Cài, cập nhật hoặc gỡ bằng một lệnh
├── sync.sh                # Sinh file quy tắc cho từng agent từ AGENTS.md
├── test.sh                # Tự kiểm: install.sh, sync.sh, load_skill.py, liên kết
└── load_skill.py          # Bộ nạp cho agent tự viết
docs/
├── en/                    # Tài liệu gốc (tiếng Anh)
└── vi/                    # Bản dịch tiếng Việt của docs/en
CHANGELOG.md               # Thay đổi theo từng phiên bản
CONTRIBUTING.md            # Quy trình đóng góp và checklist pull request
Makefile                   # make check | sync | diff | test | install
```

## Cài đặt

Một lệnh, không copy tay:

```bash
bash /đường/dẫn/tới/stack-agents/scripts/install.sh /đường/dẫn/tới/dự-án
```

Script chép `AGENTS.md`, `.ai/`, `scripts/` và `VERSION` vào dự án rồi chạy
`sync.sh` tại đó. Sau đó điền lệnh build, test, lint vào mục **Project overrides**
cuối `AGENTS.md`.

| Tùy chọn | Tác dụng |
|---|---|
| `--force` | Gộp `AGENTS.md`, `.ai/`, `scripts/` và `VERSION` hiện tại đè lên bản trong dự án, rồi chạy lại sync. Đây là cách cập nhật; script sẽ in phiên bản cũ và mới |
| `--dry-run` | In ra việc sẽ làm, không thay đổi gì |
| `--no-sync` | Chỉ chép file, không sinh file quy tắc cho agent |
| `--uninstall` | Gỡ những gì stack-agents đã cài (file sinh ra, script, kỹ năng kèm theo), giữ lại file của dự án. Kết hợp được với `--dry-run` |
| `--version` | In phiên bản của bản checkout rồi thoát |

Không có `--force`, script từ chối đụng vào bản cài đã có, nên không thể phá
hỏng công việc. Khi gộp, file nào có trong nguồn sẽ được thay, phần còn lại giữ
nguyên; `.ai/conventions.md` và mọi kỹ năng bạn tự thêm không nằm trong nguồn
nên vẫn còn sau khi cập nhật.

> **Đừng sửa file kỹ năng kèm theo** trong dự án đã cài — bản cập nhật
> `--force` sẽ thay chúng. Hãy tự tạo `.ai/skills/<tên>.md`; file thừa không bao
> giờ bị cập nhật hay gỡ bỏ đụng tới.

Để gỡ stack-agents khỏi một dự án:

```bash
bash /đường/dẫn/tới/stack-agents/scripts/install.sh --uninstall /đường/dẫn/tới/dự-án
```

<details>
<summary>Cài bằng tay, nếu script không hợp với cách bạn làm</summary>

1. Chép `AGENTS.md`, `.ai/` và `scripts/` vào thư mục gốc dự án.
2. Điền lệnh build, test, lint vào mục **Project overrides** cuối `AGENTS.md`.
3. Chạy `bash scripts/sync.sh`.

</details>

Hướng dẫn đầy đủ: [docs/vi/getting-started.md](docs/vi/getting-started.md).

## File được sinh ra

| Agent | File agent đọc |
|---|---|
| Claude Code | `CLAUDE.md` |
| Cursor | `.cursor/rules/stack-agents.mdc` (luôn bật) |
| Cline | `.clinerules` |
| Roo Code | `.roorules` |
| Windsurf | `.windsurfrules` |
| GitHub Copilot | `.github/copilot-instructions.md` |
| Codex và mọi agent đọc `AGENTS.md` | dùng thẳng `AGENTS.md` |

Tên đường dẫn thay đổi theo phiên bản công cụ. Nếu agent của bạn tìm ở chỗ
khác, sửa danh sách `place` trong `scripts/sync.sh`.

Sau mỗi lần sửa `AGENTS.md`, chạy lại `sync.sh`. Muốn xem sync sẽ thay đổi gì
mà chưa ghi ra, chạy `bash scripts/sync.sh --diff` (hoặc `make diff`). Trong CI,
chạy `bash scripts/sync.sh --check` để báo lỗi khi file sinh ra đã lệch.

## Lệnh

| Lệnh | Tác dụng |
|---|---|
| `/architect` | Thiết kế (dữ liệu, luồng, giao diện) rồi dừng chờ duyệt |
| `/arena` | 2-3 phương án khác nhau về bản chất, bảng so sánh, khuyến nghị |
| `/red-team` | Soi lỗi theo 3 góc: bảo mật, vận hành, QA |
| `/tdd` | Viết test fail trước, code tối thiểu, rồi refactor |
| `/verify` | Chạy lint, test, build và chạy thử thay đổi. Tự động trước khi báo xong |
| `/checkpoint` | Ghi và đối chiếu trạng thái quanh lệnh rủi ro |
| `/unslop` | Rút gọn văn bản, giữ nguyên số liệu và lưu ý quan trọng |
| `/conventions` | Rút quy ước từ repo và từ cuộc trò chuyện, lưu sau khi bạn duyệt |
| `/auto <mục tiêu>` | Chọn playbook phù hợp (`debug`, `refactor`) và chạy |

Đầy đủ từng lệnh: [docs/vi/skills.md](docs/vi/skills.md).

## Agent tự viết

Với vòng lặp riêng của bạn (LangChain, CrewAI, bất kỳ thứ gì khác),
`scripts/load_skill.py` đưa quy tắc ra thành system prompt và các kỹ năng thành
một tool:

```python
from scripts.load_skill import system_prompt, load_skill, list_modules

prompt = system_prompt()   # dùng làm system prompt
# đăng ký load_skill(name, kind="skills") thành một tool để agent nạp kỹ năng khi cần
```

`load_skill` chỉ nhận tên dạng `[a-z0-9-]` và hai loại `skills`, `playbooks`,
nên agent không thể dùng nó để đọc file ngoài `.ai/`.

Chi tiết: [docs/vi/custom-agents.md](docs/vi/custom-agents.md).

## Mở rộng

Thêm kỹ năng: tạo `.ai/skills/<tên>.md`, thêm một dòng vào bảng trong
`AGENTS.md`, rồi chạy `sync.sh`. Playbook làm tương tự trong `.ai/playbooks/`.
Giữ `AGENTS.md` ngắn; quy tắc dài hơn nên đặt trong `.ai/`, nơi chúng chỉ được
nạp khi cần.

Chi tiết: [docs/vi/extending.md](docs/vi/extending.md).

## Kiểm chứng

```bash
make check     # kiểm tra file sinh ra có lệch không, rồi chạy bộ test
```

`make check` chạy `scripts/sync.sh --check` và `scripts/test.sh`. Bộ test phủ
hành vi của `install.sh`, khả năng phát hiện lệch của `sync.sh`, CLI của
`load_skill.py`, việc mọi liên kết tương đối trong tài liệu đều tồn tại, và
việc `docs/en` cùng `docs/vi` không bị lệch. CI chạy đúng lệnh này, cộng thêm
ShellCheck cho các script Bash và Ruff cho file Python.

## Đóng góp

Quy trình, versioning và checklist pull request:
[CONTRIBUTING.md](CONTRIBUTING.md). Thay đổi gần đây:
[CHANGELOG.md](CHANGELOG.md).

## Giới hạn cần biết

- Đây là chỉ dẫn bằng văn bản, không phải cơ chế cưỡng chế. Mức độ tuân thủ phụ
  thuộc vào mô hình và công cụ. Với thao tác nguy hiểm, vẫn cần quyền hạn và
  sandbox ở phía công cụ.
- Đường dẫn file quy tắc của từng agent ở trên chưa được kiểm chứng với mọi
  phiên bản hiện hành của mọi công cụ. Hãy thử một yêu cầu đơn giản (ví dụ
  `/verify`) để xác nhận agent đã đọc cấu hình.

## Giấy phép

[MIT](LICENSE).
