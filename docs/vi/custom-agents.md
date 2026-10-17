# Agent tự viết

> English: [docs/en/custom-agents.md](../../docs/en/custom-agents.md)

Các file quy tắc sinh ra phục vụ những công cụ đọc một đường dẫn cố định. Nếu
bạn chạy vòng lặp riêng — LangChain, CrewAI, hay thứ bạn tự viết —
`scripts/load_skill.py` cho bạn đúng hai thứ đó: quy tắc thành system prompt, và
kỹ năng thành một tool để mô hình gọi khi cần.

## `system_prompt()`

```python
from scripts.load_skill import system_prompt

prompt = system_prompt()   # nội dung của AGENTS.md
```

Trả về nội dung `AGENTS.md`. Dùng nó làm system prompt cho agent.

## `list_modules()`

```python
from scripts.load_skill import list_modules

list_modules()
# {'skills': ['architect', 'arena', ...], 'playbooks': ['debug', 'refactor']}
```

Trả về các tên đang có, đã sắp xếp, nhóm theo loại. Hữu ích khi dựng mô tả cho
tool hoặc một lệnh trợ giúp.

## `load_skill(name, kind="skills")`

```python
from scripts.load_skill import load_skill

load_skill("tdd")                     # nội dung .ai/skills/tdd.md
load_skill("debug", kind="playbooks") # nội dung .ai/playbooks/debug.md
```

`kind` là `"skills"` (mặc định) hoặc `"playbooks"`.

## Kiểm tra tên chính là tính an toàn

`load_skill` chỉ nhận tên khớp `^[a-z0-9][a-z0-9-]*$`, và chỉ nhận hai loại
trên. Một cái tên như `../../etc/passwd` sẽ ném `ValueError`, nên mô hình có
sinh ra tham số tool độc hại hay lộn xộn cũng không thể dùng hàm này để đọc file
ngoài `.ai/`.

Đây là lý do ràng buộc này tồn tại, chứ không phải một regex cho đẹp. Nếu bạn
viết lại bộ nạp, giữ nguyên nó.

## Kiểm tra từ shell

```bash
python3 -I scripts/load_skill.py list
python3 -I scripts/load_skill.py skills tdd
python3 -I scripts/load_skill.py playbooks debug
```

`-I` cô lập trình thông dịch khỏi môi trường, đúng thứ bạn cần cho một lần kiểm
tra nhanh. Truyền sai loại hoặc tên không hợp lệ sẽ thoát khác 0.

## Nối thành một tool

Mở một tool với schema `{name: string, kind: "skills" | "playbooks"}` và phần
thân gọi `load_skill`, trả về nội dung. Viết mô tả cho mô hình biết chỉ gọi khi
một lệnh được kích hoạt, để hành vi khớp với vòng làm việc trong `AGENTS.md`.

Bộ test của repo có phủ phần kiểm tra tên:

```bash
bash scripts/test.sh
```

## Tiếp theo

- [how-it-works.md](how-it-works.md) — mô hình sync đứng sau các file quy tắc.
- [skills.md](skills.md) — các file được nạp thực sự chứa gì.
