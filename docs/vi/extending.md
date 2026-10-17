# Mở rộng

> English: [docs/en/extending.md](../../docs/en/extending.md)

## Thêm một kỹ năng

1. Tạo `.ai/skills/<tên>.md`. Tên chỉ được chứa chữ thường, chữ số và dấu gạch
   ngang — `load_skill` từ chối mọi thứ khác.
2. Thêm một dòng vào bảng kỹ năng trong `AGENTS.md`, ánh xạ lệnh sang file và
   nói khi nào dùng nó.
3. Chạy `bash scripts/sync.sh` để các file quy tắc đã sinh nhận bảng mới.

Hết. Kỹ năng giờ có mặt dưới dạng `/<tên>` ở mọi công cụ, và qua
`load_skill("<tên>")` trong agent tự viết.

Viết file kỹ năng theo đúng cách các file đang có được viết: tiêu đề ngắn, khi
nào dùng, các bước, và kết quả. Giữ nó nói về quy trình, không nói về một dự án
cụ thể.

## Thêm một playbook

Cùng hình dạng, đặt trong `.ai/playbooks/`. Thêm tên đó vào dòng `/auto` trong
`AGENTS.md` để mô hình biết nó tồn tại, rồi chạy `sync.sh`.

## Giữ `AGENTS.md` ngắn

`AGENTS.md` nằm trong cửa sổ ngữ cảnh ở mọi yêu cầu. Thứ gì dài hơn vài dòng nên
đặt trong `.ai/` và chỉ được nạp khi cần. Bảng kỹ năng là cỡ phù hợp cho
`AGENTS.md`; các bước của kỹ năng thì không.

Quy tắc ngón tay: nếu nội dung chỉ quan trọng khi một lệnh được gọi, nó thuộc về
file mà lệnh đó nạp.

## Kiểm chứng thay đổi của bạn

```bash
make check
```

Lệnh này chạy `scripts/sync.sh --check` và `scripts/test.sh`. Bộ test kiểm tra
`install.sh`, `sync.sh`, CLI của `load_skill.py`, việc mọi liên kết tương đối
trong tài liệu đều tồn tại, và việc `docs/en` cùng `docs/vi` không bị lệch. Nếu
bạn thêm kỹ năng, xác nhận nó xuất hiện:

```bash
python3 -I scripts/load_skill.py list
```

## Giữ hai bản dịch khớp nhau

Tài liệu tồn tại hai lần — `docs/en/` và `docs/vi/`, với tên file giống hệt nhau.
Nếu bạn sửa một bên, sửa luôn bên kia trong cùng một lần commit. `scripts/test.sh`
sẽ fail khi một file chỉ có ở một bên, hoặc khi hai bản của cùng một file có số
tiêu đề khác nhau — cả hai đều là dấu hiệu bản dịch bị bỏ quên. Chỉ văn bản đã
cũ mà cấu trúc không đổi mới lọt lưới, nên hãy coi hai file là một thay đổi.

`AGENTS.md` và `.ai/` cố ý chỉ có tiếng Anh: chúng được nạp vào ngữ cảnh của mô
hình ở mọi phiên, và một bản thứ hai sẽ nhân đôi chi phí đó.

## Đóng góp

Checklist đầy đủ nằm trong [../../CONTRIBUTING.md](../../CONTRIBUTING.md).
Riêng phần mở rộng:

1. `make check` phải pass.
2. `AGENTS.md` và các file sinh ra phải khớp — nếu không, CI sẽ fail.
3. Kỹ năng mới phải được đăng ký trong bảng ở `AGENTS.md`.
4. Mọi thay đổi tài liệu phải có đủ cả hai bản ngôn ngữ.

## Tiếp theo

- [skills.md](skills.md) — các kỹ năng đang có, làm mẫu để theo.
- [how-it-works.md](how-it-works.md) — vì sao ranh giới được vẽ như vậy.
- [../../README.vi.md](../../README.vi.md) — bản ngắn.
