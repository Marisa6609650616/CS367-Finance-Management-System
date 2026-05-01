# Personal Finance Management API System
**วิชา:** CS367 Web Service Development Concepts  

---

## รายละเอียดโปรเจกต์ (Overview)
ระบบ Backend API สำหรับจัดการข้อมูลการเงินส่วนบุคคลที่ช่วยแก้ปัญหาการจดบันทึกแบบเดิมๆ โดยระบบสามารถ:
- บันทึกรายรับ-รายจ่ายแยกตามหมวดหมู่
- สรุปยอดคงเหลือสุทธิและยอดรวมรายเดือน
- ใช้งานอย่างปลอดภัยผ่านระบบ Login (JWT)

## ฟีเจอร์ที่พัฒนา (Features)
1. **Authentication:** ระบบสมาชิกและเข้าสู่ระบบ
2. **Transaction Management:** การจัดการ เพิ่ม/แก้ไข/ลบ รายการทางการเงิน
3. **Financial Summary:** ระบบสรุปผลยอดเงินคงเหลือและรายงานรายเดือน

## สมาชิกและการแบ่งงาน (Responsibilities)
| ชื่อ-นามสกุล | หน้าที่ (API) | งานส่วนกลาง |
| :--- | :--- | :--- |
| นายเสฎฐวุฒิ วิจิตรศิลป์ | Register, Login | Security & JWT Middleware |
| นายจันทร์พงศ์ วิทยอรุณธานี | Get, Post Transactions | DB Schema & Connection |
| นางสาวนัทธ์ชนัน ชาคริตบุษบง | Put, Delete Transactions | Docker & Git Management |
| นายธนดล กล่อมใจ | Summary Monthly | Code Review & Test Coverage |
| นางสาวมาริษา จันทร์ทอง | Summary Balance | Postman & README Documentation |

## Tech Stack (เบื้องต้น)
- **Language:** Go / Java / [ภาษาที่กลุ่มเลือก]
- **Database:** PostgreSQL / MySQL
- **Tool:** Docker, Postman, Git

## การทดสอบ (Testing)
- Unit Test Coverage: ไม่น้อยกว่า 80%
- API Testing via Postman Collection
