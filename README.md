# 💰 Personal Finance Management API System

> **วิชา:** CS367 Web Service Development Concepts

---

## 📋 รายละเอียดโปรเจกต์ (Overview)

ระบบ **Backend API** สำหรับจัดการข้อมูลการเงินส่วนบุคคล ออกแบบตามสถาปัตยกรรม **RESTful API** พร้อมระบบแบ่งแยกสิทธิ์ข้อมูลของผู้ใช้แต่ละคนอย่างชัดเจน

**ความสามารถหลักของระบบ:**

- 🔐 สมัครสมาชิกและเข้าสู่ระบบอย่างปลอดภัยด้วย JWT
- 📝 บันทึกรายรับ-รายจ่ายแยกตามหมวดหมู่ (Category)
- 📊 สรุปยอดคงเหลือสุทธิ (Balance) และยอดรวมรายเดือน

---

## ✨ ฟีเจอร์ที่พัฒนา (Features)

| #   | Feature                    | รายละเอียด                                                |
| --- | -------------------------- | --------------------------------------------------------- |
| 1   | **Authentication**         | ระบบสมาชิกและเข้าสู่ระบบ ตรวจสอบสิทธิ์ผ่าน JWT Middleware |
| 2   | **Transaction Management** | เพิ่ม / แก้ไข / ลบ รายการทางการเงินแยกตามรายบุคคล         |
| 3   | **Financial Summary**      | สรุปยอดเงินคงเหลือสะสม และรายงานภาพรวมรายเดือน            |

---

## 👥 สมาชิกและการแบ่งงาน (Responsibilities)

| ชื่อ-นามสกุล                | หน้าที่ (API)            | งานส่วนกลาง / ส่วนเสริม                   |
| :-------------------------- | :----------------------- | :---------------------------------------- |
| นายเสฎฐวุฒิ วิจิตรศิลป์     | Register, Login          | Security & JWT Middleware                 |
| นายจันทร์พงศ์ วิทยอรุณธานี  | Get, Post Transactions   | DB Schema & Connection                    |
| นางสาวนัทธ์ชนัน ชาคริตบุษบง | Put, Delete Transactions | Docker & Git Management                   |
| นายธนดล กล่อมใจ             | Summary Monthly          | Code Review & Test Coverage               |
| นางสาวมาริษา จันทร์ทอง      | Summary Balance          | Postman Collection & README Documentation |

---

## 🛠️ Tech Stack

| Layer                | Technology                    |
| -------------------- | ----------------------------- |
| **Language**         | Go (Golang)                   |
| **Web Framework**    | Gin Web Framework             |
| **Database**         | SQLite (Driver: `go-sqlite3`) |
| **Containerization** | Docker & Docker Compose       |
| **API Testing**      | Postman                       |

---

## 🚀 วิธีการติดตั้งและรันระบบ (How to Run)

### 1. ตั้งค่า Environment Variables

สร้างไฟล์ `.env` ไว้ที่ **Root Directory** ของโปรเจกต์ แล้วกำหนดค่าดังนี้:

```env
PORT=8080
DB_PATH=finance.db
JWT_SECRET=your_secret_key_here
```

### 2. รันแบบ Local (Go Command)

```bash
go run cmd/server/main.go
```

### 3. รันผ่าน Docker Compose

```bash
docker-compose up --build
```

---

## 📡 API Documentation & Endpoints

ระบบผ่านการทดสอบแบบ **End-to-End** เรียบร้อยแล้ว  
สามารถนำไฟล์ `finance-api.postman_collection.json` ไป Import ใน Postman เพื่อทดสอบได้ทันที

### 📊 รายการ Endpoints ทั้งหมด

| Endpoint            | Method | Path                   | Auth Required | Description                                   |
| :------------------ | :----: | :--------------------- | :-----------: | :-------------------------------------------- |
| **Register**        | `POST` | `/api/auth/register`   |      ❌       | สมัครสมาชิกใหม่ (`email`, `password`, `name`) |
| **Login**           | `POST` | `/api/auth/login`      |      ❌       | เข้าสู่ระบบเพื่อรับ JWT Bearer Token          |
| **Transactions**    | `POST` | `/api/transactions`    |      ✅       | เพิ่มรายการรายรับ / รายจ่าย                   |
| **Summary Balance** | `GET`  | `/api/summary/balance` |      ✅       | ดูสรุปยอดเงินคงเหลือสะสม                      |

> **หมายเหตุ:** สำหรับ Endpoint ที่ต้องใช้ Auth ให้แนบ Token ในรูปแบบ `Bearer Token` ทุกครั้ง  
> ระบบจะผูกการคำนวณกับสิทธิ์ `user_id` ของผู้ใช้แต่ละคน
