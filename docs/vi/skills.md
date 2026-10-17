# Lệnh

> English: [docs/en/skills.md](../../docs/en/skills.md)

Tám kỹ năng và một bộ chọn playbook, cộng `/auto` để chọn playbook. Mỗi lệnh
ứng với một file trong `.ai/`; trang này tóm tắt chúng, còn file vẫn là nguồn
chuẩn.

Chỉ `/verify` chạy mà không cần gọi. Các lệnh còn lại chạy khi bạn gọi, hoặc khi
agent hỏi trước và bạn đồng ý.

## /architect

**Dùng khi** công việc cần thiết kế trước khi có code: tính năng mới, module
mới, thay đổi schema hoặc API.

Nhắc lại mục tiêu và ràng buộc trong tối đa ba dòng, liệt kê giả định, và chỉ
hỏi những câu chặn thiết kế. Đọc code mà thiết kế đụng tới, một lần. Rồi chỉ
trình bày những phần phù hợp: thành phần và trách nhiệm, mô hình dữ liệu (dùng
Mermaid `erDiagram` khi có hơn hai thực thể), luồng dữ liệu, hợp đồng giao diện,
các kiểu hỏng hóc, migration và rollout, chiến lược test. Với mỗi lựa chọn quan
trọng, nó nêu phương án đã bị loại.

**Kết quả:** một trang. Nếu cần hơn, phạm vi quá lớn và nó đề xuất tách ra. Nó
không viết code và dừng chờ bạn duyệt hoặc yêu cầu sửa.

## /arena

**Dùng khi** có từ hai phương án khả thi trở lên và việc chọn rất khó đảo
ngược. Bỏ qua khi một phương án rõ ràng tốt hơn.

Chọn tiêu chí phù hợp với việc đang làm — độ trễ, bộ nhớ, sự đơn giản, khả năng
test, khả năng vận hành, chi phí migration — chứ không dùng một danh sách cố
định. Đưa ra hai phương án, tối đa ba, gọi tên theo bản chất ("cache trong
tiến trình", "Redis") và khác nhau ở một đánh đổi thật. Phác qua cách mỗi phương
án hoạt động, rủi ro chính và chi phí đảo ngược về sau.

**Kết quả:** một bảng so sánh, một khuyến nghị, và điều gì sẽ khiến nó đổi ý. Nó
không triển khai gì cho tới khi bạn chọn. Nếu một phương án bị lấn át rõ ràng,
nó nói thẳng thay vì kéo dài bảng cho đủ.

## /red-team (bí danh /swarm)

**Dùng khi** một thiết kế hoặc một diff cần bị phản biện trước khi ra mắt.

Tấn công qua ba góc:

- **Bảo mật** — kiểm tra đầu vào, xác thực/phân quyền, injection, bí mật, rò rỉ
  dữ liệu, rủi ro từ phụ thuộc.
- **Vận hành** — hành vi khi lỗi và thử lại, timeout, tải, khả năng quan sát,
  rollback, migration, cấu hình trôi lệch.
- **QA** — nhánh chưa được test, trường hợp biên, đồng thời, đầu vào xấu, tương
  thích ngược, yêu cầu mơ hồ.

Nếu có subagent, nó chạy mỗi góc trong một subagent riêng để các góc không neo
vào nhau; nếu không, nó chạy tuần tự và nói rõ mức độ độc lập bị hạn chế.

**Kết quả:** tối đa năm phát hiện mỗi góc, mỗi phát hiện có mức độ nghiêm trọng,
bằng chứng (file và dòng, hoặc một tình huống cụ thể) và gợi ý sửa, rồi một kết
luận: **ship**, **sửa trước**, hoặc **thiết kế lại**. Nó không tự sửa gì cho tới
khi bạn chọn.

## /tdd

**Dùng khi** hành vi có thể test được và repo đã có bộ test.

Mỗi vòng một hành vi: **đỏ** (viết một test fail và xác nhận nó fail vì đúng lý
do), **xanh** (viết ít code nhất để pass), **refactor** (cải thiện cấu trúc khi
test còn xanh). Sửa lỗi thì bắt đầu bằng một test tái hiện lỗi.

Không bao giờ viết code production khi chưa có test fail đòi hỏi nó. Test hành
vi qua giao diện công khai, không qua chi tiết nội bộ. Chỉ mock ở ranh giới hệ
thống — mạng, đồng hồ, hệ thống file — và không mock thứ mình sở hữu trừ khi
không còn cách nào rẻ hơn. Không sửa test để nó pass, trừ khi chính yêu cầu đã
đổi, và phải nói rõ khi điều đó xảy ra. Nếu repo chưa có bộ test, nó nói ra và
đề xuất bộ nhỏ nhất phù hợp với công nghệ đang dùng.

## /verify

**Dùng khi** sắp báo cáo là đã xong. Lệnh này tự chạy.

Viết các điểm nghiệm thu từ yêu cầu ban đầu thành một danh sách. Chạy các lệnh
trong **Project overrides**, theo thứ tự: format và lint, typecheck, test trọng
tâm, toàn bộ bộ test, rồi build. Sau đó chạy thử thay đổi trực tiếp — mở CLI, gọi
endpoint, mở trang — vì test pass không phủ được hành vi UI, CLI hay tích hợp.
Cuối cùng nó tự soát diff của mình để tìm output debug sót lại, thay đổi không
liên quan, bí mật lọt vào, và TODO còn dang dở.

**Kết quả:** hai danh sách — đã kiểm chứng (mỗi lệnh và kết quả) và chưa kiểm
chứng (cái gì và vì sao). Nó không bao giờ nói "xong", "đã sửa" hay "chạy được"
mà không có kết quả lệnh đứng sau.

## /checkpoint

**Dùng khi** một lệnh thay đổi nhiều file hoặc khó hoàn tác: migration, sinh
code, viết lại hàng loạt, nâng cấp phụ thuộc, chạy formatter trên toàn repo.
Không dùng cho lệnh chỉ đọc.

Trước: ghi lại `git status --short` và `git rev-parse HEAD`, hoặc liệt kê thư mục
liên quan kèm thời gian sửa đổi nếu ngoài git, và cảnh báo nếu có thay đổi chưa
commit mà lệnh có thể ghi đè. Sau: so lại status và diff, xác nhận chỉ những file
mong đợi đã thay đổi.

**Nếu có gì bất ngờ thay đổi** nó dừng, báo rõ file nào và khác thế nào, rồi đưa
lựa chọn. Nó không chạy thêm lệnh phá hủy nào để dọn dẹp.

## /unslop

**Dùng khi** cần rút gọn văn bản: câu trả lời, tài liệu, commit message, mô tả
PR, comment.

Cắt phần mở đầu và kết thúc, phần nhắc lại câu hỏi, các từ rào đón không mang
thông tin, phần giải thích điều mà code đã nói rõ, và tính từ không có gì đo
đếm được. Giữ nguyên tên, số liệu, phiên bản, câu lệnh và đường dẫn, cùng với
các lưu ý làm thay đổi việc người đọc sẽ làm, các quyết định và lý do, và các
bước theo đúng thứ tự.

**Cách làm:** viết nháp, xóa mọi thứ không làm thay đổi điều người đọc biết hoặc
làm, ưu tiên một ý mỗi câu và câu chủ động, rồi kiểm tra không mất dữ kiện hay
bước nào. Với code, nó xóa code chết, comment nhắc lại code, và lớp bọc chỉ
chuyển tiếp lời gọi — không đụng tới hành vi.

## /conventions

**Dùng khi** cần ghi lại quy ước của một repo để các phiên sau theo mà không phải
nhắc lại.

Thu thập quy ước ứng viên từ cấu hình (linter, formatter, `tsconfig`,
`.editorconfig`, workflow CI), từ code (cách đặt tên, bố cục file, xử lý lỗi,
phong cách test ở vài file tiêu biểu), từ `git log` để biết cách viết commit, và
từ những chỗ bạn chỉnh nó trong cuộc trò chuyện. Chỉ giữ quy ước xuất hiện ở ít
nhất hai nơi hoặc do bạn nói ra. Viết mỗi quy ước thành một dòng mệnh lệnh kèm
bằng chứng trong ngoặc, nhóm theo chủ đề, và giữ file dưới 40 dòng.

**Kết quả:** đề xuất file `.ai/conventions.md` và chỉ lưu sau khi bạn duyệt.
Không bao giờ chứa bí mật hay dữ liệu cá nhân, đánh dấu `(inferred)` cho quy ước
chưa được bạn xác nhận, và hợp nhất vào file đang có thay vì viết lại.

## /auto <mục tiêu>

**Dùng khi** bạn có mục tiêu nhưng không có lệnh cụ thể.

Chọn playbook phù hợp trong `.ai/playbooks/`, nói rõ nó chọn cái nào trong một
dòng, rồi chạy. Hiện có hai:

- **`debug`** — tái hiện, cô lập, nêu nguyên nhân gốc trong một câu, sửa ít nhất,
  giữ lại test tái hiện lỗi, dọn các thí nghiệm thất bại. Nó có nhánh riêng cho
  sự cố production, bắt đầu ở chế độ chỉ đọc và đề xuất giảm thiểu (rollback,
  feature flag, tăng tài nguyên) trước khi sửa gốc.
- **`refactor`** — nêu rõ cái gì được cải thiện và hành vi nào không được đổi,
  thêm test đặc tả ở chỗ thiếu độ phủ, chia bước sao cho mỗi bước vẫn build và
  test xanh, checkpoint giữa các bước, và hoàn tác bước nào chuyển đỏ mà chưa rõ
  nguyên nhân trong một lần thử.

Nếu không playbook nào phù hợp, nó nói ra và tiếp tục theo vòng làm việc.

## Tiếp theo

- [getting-started.md](getting-started.md) — cài đặt.
- [extending.md](extending.md) — thêm của riêng bạn.
