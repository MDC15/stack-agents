# Mở rộng

> English: [docs/en/extending.md](../../docs/en/extending.md)

## Thêm một skill

1. Tạo `.ai/skills/<name>.md`. Tên chỉ chứa chữ thường, chữ số và gạch nối —
   `load_skill` từ chối mọi thứ khác.
2. Thêm một dòng vào bảng skill trong `AGENTS.md`, ánh xạ lệnh sang file và
   nêu khi nào dùng nó.
3. Chạy `stack sync` để các file render và link nhặt được bảng mới.

Chỉ vậy thôi. Skill giờ khả dụng dạng `/<name>` trong mọi tool, và qua
`load_skill("<name>")` trong custom agent.

Viết chính file skill theo cách các file hiện tại được viết: tiêu đề ngắn, khi
nào dùng, các bước, và nó tạo ra cái gì. Giữ về quy trình, không về một project
cụ thể.

## Thêm một playbook

Cùng dáng, trong `.ai/playbooks/`. Thêm tên vào dòng `/auto` trong `AGENTS.md`
để model biết nó tồn tại, rồi chạy `stack sync`.

## Thêm một agent

Một dòng trong `internal/agents/agents.go` — ID, tên hiển thị, kind
(symlink, bản copy banner, hoặc direct như OpenCode đọc trực tiếp
`AGENTS.md`), file được sinh tới nơi nào trong `.agents/`, và đường dẫn cố
định tool đọc. Tín hiệu detection nằm cạnh registry trong
`internal/detect/detect.go`; thêm probe agent cần — binary trên PATH, một thư
mục đánh dấu dưới `$HOME`, hoặc một VS Code extension đã cài. Commit cùng
`stack sync --all` để file và link sinh ra của agent mới xuất hiện trong cùng
một commit; `make check` xác nhận cả hai.

## Giữ `AGENTS.md` ngắn

`AGENTS.md` nằm trong context window cho mỗi request. Bất kỳ thứ gì dài hơn
vài dòng nên nằm trong `.ai/` và chỉ load khi cần. Bảng skill là cỡ đúng cho
`AGENTS.md`; các bước của skill thì không.

Quy tắc ngón tay cái: nếu đoạn văn chỉ quan trọng khi một lệnh được gọi, nó
thuộc về file mà lệnh đó load.

## Xác minh thay đổi

```bash
make check
```

Lệnh này chạy gofmt, `go vet`, `go test ./...` (bao gồm link và parity
kiểm tra tài liệu, phát hiện drift, hành vi install và uninstall, test
subprocess của CLI, và CLI của `load_skill.py`) và cuối cùng là `stack check`.
Nếu bạn thêm skill, xác nhận nó xuất hiện:

```bash
python3 -I scripts/load_skill.py list
```

## Giữ hai bản dịch cùng bước

Tài liệu tồn tại hai lần — `docs/en/` và `docs/vi/`, đặt trùng tên file. Nếu
bạn đổi một bên, đổi bên còn lại trong cùng commit. Test parity song ngữ
thoát lỗi khi file chỉ tồn tại một bên, hoặc hai phiên bản của một file mang
số lượng heading khác nhau — hai dấu hiệu của bản dịch bị bỏ lại. Chỉ văn
bản lỗi thời trượt qua, nên coi hai file là một thay đổi.

`AGENTS.md` và `.ai/` cố ý chỉ dùng tiếng Anh: chúng được load vào context của
model mỗi phiên, và một bản copy thứ hai sẽ gấp đôi chi phí đó.

## Đóng góp

Checklist đầy đủ nằm ở [../../CONTRIBUTING.md](../../CONTRIBUTING.md).
Riêng cho việc mở rộng:

1. `make check` pass.
2. `AGENTS.md` và các file sinh ra in cùng bước — CI thất bại nếu ngược lại.
3. Skill mới được đăng ký trong bảng `AGENTS.md`.
4. Cả hai phiên bản ngôn ngữ của mọi thay đổi tài liệu đều có.

## Kế tiếp

- [skills.md](skills.md) — các skill hiện có, làm mẫu để bắt chước.
- [how-it-works.md](how-it-works.md) — vì sao ranh giới được vẽ thế này.
- [cli.md](cli.md) — các lệnh mà phần kiểm chứng xoay quanh.
- [../../README.md](../../README.md) — bản tóm tắt.