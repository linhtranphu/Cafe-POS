1. Vấn đề của state machine hiện tại (BA critique)

Ở bản trước:

PAID → IN_PROGRESS : send_to_bar
IN_PROGRESS → SERVED : serve_drink


❌ Vấn đề:

Không phân biệt:

ai gửi

ai xác nhận

IN_PROGRESS vừa là:

trạng thái nghiệp vụ

vừa là hành động mơ hồ

👉 Dễ dẫn tới:

Waiter “tự bấm đã pha”

Không audit được lỗi pha chế

Barista bị “vô hình” trong hệ thống

2. Chốt lại vai trò chuẩn (BA definition)
Vai trò trong Order lifecycle
Role	Trách nhiệm
Waiter	Tạo order, thu tiền, gửi pha chế
Barista	Xác nhận pha, hoàn tất đồ uống
Cashier	Kiểm soát, khóa dữ liệu
System	Enforce rule

👉 Barista PHẢI là người làm thay đổi state vận hành

3. Làm rõ lại các state (rất quan trọng)
Định nghĩa lại state cho đúng nghiệp vụ
State	Ý nghĩa BA
CREATED	Order chưa thanh toán
PAID	Đã thu tiền, chưa giao cho pha chế
IN_PROGRESS	Barista đã nhận và đang pha
READY	Pha xong, chờ giao
SERVED	Đã giao cho khách
LOCKED	Đã chốt ca

👉 IN_PROGRESS không phải do waiter set, mà do Barista nhận việc

4. Phân quyền trigger state (KEY BA POINT)
Transition	Ai được làm
CREATED → PAID	Waiter
PAID → (send to bar)	Waiter
PAID → IN_PROGRESS	❌ (không ai trực tiếp)
(Barista nhận) → IN_PROGRESS	Barista
IN_PROGRESS → READY	Barista
READY → SERVED	Waiter
SERVED → LOCKED	Cashier

👉 Đây là điểm rất nhiều hệ thống POS làm sai.

5. STATE MACHINE – HOÀN CHỈNH (có Barista)
stateDiagram-v2
    [*] --> CREATED

    CREATED --> PAID : collect_payment\n(Waiter)

    PAID --> PAID : edit / refund / collect_more\n(Waiter/Cashier)

    PAID --> QUEUED : send_to_bar\n(Waiter)

    QUEUED --> IN_PROGRESS : accept_order\n(Barista)

    IN_PROGRESS --> READY : finish_preparing\n(Barista)

    READY --> SERVED : deliver_drink\n(Waiter)

    SERVED --> LOCKED : close_shift\n(Cashier)

    PAID --> CANCELLED : refund_full\n(Waiter/Cashier)

    CANCELLED --> LOCKED : close_shift\n(Cashier)

6. Vì sao cần state QUEUED? (BA justification)

Đây là state rất quan trọng, không phải thừa.

Nếu không có QUEUED:

Không biết:

Order đã gửi chưa

Barista có nhận chưa

Không đo được:

Thời gian chờ pha chế

Không quản lý được:

Hàng đợi pha

👉 QUEUED = handoff point giữa phục vụ & pha chế

7. Business Rules bổ sung (bắt buộc ghi)
BR-06: Only Barista can move order to IN_PROGRESS.
BR-07: Once order enters IN_PROGRESS, no modification or refund is allowed.
BR-08: Payment adjustments are allowed only before QUEUED.
BR-09: READY indicates drink is completed but not yet delivered.

8. BA insight chốt hạ 🔒☕️

State machine chuẩn không chỉ kiểm soát dữ liệu,
mà kiểm soát TRÁCH NHIỆM CON NGƯỜI

Waiter → tiền & giao tiếp

Barista → nguyên liệu & chất lượng

Cashier → tiền & sổ sách

Nếu thiếu Barista trong state machine → thiếu 1 trụ cột vận hành