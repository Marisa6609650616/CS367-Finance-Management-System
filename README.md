# Personal Finance Management API System

### 1. Authentication APIs
| :--- | :--- | :--- | :--- |
| `/api/auth/register` | `POST` | สมัครสมาชิกใหม่ | นายเสฎฐวุฒิ วิจิตรศิลป์ |
| `/api/auth/login` | `POST` | เข้าสู่ระบบ | นายเสฎฐวุฒิ วิจิตรศิลป์ |

### 2. Transaction APIs
| :--- | :--- | :--- | :--- |
| `/api/transactions` | `GET` | ดึงรายการรายรับ-รายจ่ายทั้งหมด | นายจันทร์พงศ์ วิทยอรุณธานี |
| `/api/transactions` | `POST` | เพิ่มรายการรายรับหรือรายจ่ายใหม่ | นายจันทร์พงศ์ วิทยอรุณธานี |
| `/api/transactions/{id}` | `PUT` | แก้ไขรายการที่มีอยู่ | นางสาวนัทธ์ชนัน ชาคริตบุษบง |
| `/api/transactions/{id}` | `DELETE` | ลบรายการรายรับรายจ่าย | นางสาวนัทธ์ชนัน ชาคริตบุษบง |

### 3. Summary APIs
| :--- | :--- | :--- | :--- |
| `/api/summary/monthly` | `GET` | สรุปยอดรายรับ-รายจ่ายรายเดือน | นายธนดล กล่อมใจ |
| `/api/summary/balance` | `GET` | แสดงยอดคงเหลือสุทธิทั้งหมด | นางสาวมาริษา จันทร์ทอง |

