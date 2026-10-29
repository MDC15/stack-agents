# Custom agent

> English: [docs/en/custom-agents.md](../../docs/en/custom-agents.md)

Các file rules sinh ra phủ các tool đọc một đường dẫn cố định. Nếu bạn tự chạy
vòng lặp riêng — LangChain, CrewAI, hoặc thứ bạn viết — `scripts/load_skill.py`
cung cấp đúng hai mảnh đó: rules dưới dạng system prompt, và skills dưới dạng
tool mà model gọi khi cần.

## `system_prompt()`

```python
from scripts.load_skill import system_prompt

prompt = system_prompt()   # nội dung của AGENTS.md
```

Trả về nội dung `AGENTS.md`. Dùng làm system prompt cho agent.

## `list_modules()`

```python
from scripts.load_skill import list_modules

list_modules()
# {'skills': ['architect', 'arena', ...], 'playbooks': ['debug', 'refactor']}
```

Trả về các tên khả dụng, đã sort, nhóm theo kind. Hữu ích khi dựng mô tả tool
hoặc một lệnh help.

## `load_skill(name, kind="skills")`

```python
from scripts.load_skill import load_skill

load_skill("tdd")                     # text của .ai/skills/tdd.md
load_skill("debug", kind="playbooks") # text của .ai/playbooks/debug.md
```

`kind` là `"skills"` (mặc định) hoặc `"playbooks"`.

## Name check là tính chất an toàn

`load_skill` chỉ nhận tên khớp `^[a-z0-9][a-z0-9-]*$`, và chỉ hai kind phía
trên. Một tên kiểu `../../etc/passwd` sẽ raise `ValueError`, nên model phát
ra tham số tool hostile hoặc bị hỏng không thể dùng hàm này để đọc file ngoài
`.ai/`.

Đây là lý do ràng buộc tồn tại chứ không phải một regex trang trí. Nếu bạn viết
lại loader, hãy giữ nó.

## Kiểm tra từ shell

```bash
python3 -I scripts/load_skill.py list
python3 -I scripts/load_skill.py skills tdd
python3 -I scripts/load_skill.py playbooks debug
```

`-I` cô lập interpreter khỏi môi trường, điều bạn muốn cho một sanity check.
Truyền kind lạ hoặc tên không hợp lệ thoát non-zero.

## Nối làm tool

Expose một tool với schema `{name: string, kind: "skills" | "playbooks"}` và
kêu body gọi `load_skill`, trả về text. Phần mô tả bảo model chỉ gọi khi một
lệnh được gọi, để hành vi khớp vòng vận hành trong `AGENTS.md`.

Bộ test của repo phủ name check — giờ là
[`repo_test.go`](../../repo_test.go), chạy bằng `go test ./...`; và
[cli.md](cli.md) ghi CLI tương tự mà các test điều khiển.

## Kế tiếp

- [how-it-works.md](how-it-works.md) — mô hình sync đằng sau các tool file.
- [skills.md](skills.md) — các file được load thực chứa gì.