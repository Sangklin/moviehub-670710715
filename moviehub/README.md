# MovieHub Go API

งานตาม [Week 12: The Ticket That Cannot Be Faked](https://plearnjai.com/be2026/Week12-The-Ticket-That-Cannot-Be-Faked.html) และ [Lab: Ticket Booth and Gatekeepers](https://plearnjai.com/be2026/Lab-Week12-Ticket-Booth-and-Gatekeepers.html) เชื่อมกับ React ตาม [Week 9](https://plearnjai.com/fe2026/week9-backend-mock.html) และทดสอบตาม [Week 10](https://plearnjai.com/fe2026/week10-testing.html)

## รันในเครื่อง

ใช้ Go ตามเวอร์ชันใน `go.mod`, Node.js และ Docker Compose คำสั่งต่อไปนี้ใช้ Bash และเริ่มที่โฟลเดอร์หลักของ repository

```bash
docker compose up -d
```

รอให้ PostgreSQL พร้อม แล้วสร้างและโหลดฐานข้อมูล **ครั้งแรกเท่านั้น** (ถ้ามีฐานข้อมูลและข้อมูลอยู่แล้วให้ข้าม):

```bash
docker compose exec db psql -U postgres -d kitchen -c 'CREATE DATABASE moviehub;'
docker compose exec -T db psql -U postgres -d moviehub -v ON_ERROR_STOP=1 < moviehub/schema.sql
docker compose exec -T db psql -U postgres -d moviehub -v ON_ERROR_STOP=1 < moviehub/seed.sql
```

เปิด API:

```bash
cd moviehub
test -s .jwt_secret || (umask 077; openssl rand -hex 32 > .jwt_secret)
export DATABASE_URL='postgres://postgres:secret@localhost:5432/moviehub?sslmode=disable'
export JWT_SECRET="$(cat .jwt_secret)"
export ALLOWED_ORIGINS='http://localhost:3000'
go run .
```

เปิด terminal อีกอันที่โฟลเดอร์หลัก แล้วรัน React ให้เรียก Go:

```bash
npm ci
REACT_APP_API_URL=http://localhost:8080 npm start
```

บัญชี `demo@moviehub.test / 1234` ในหน้า Login ใช้กับ mock เท่านั้น เมื่อใช้ Go ให้สมัครสมาชิกใหม่ รหัสผ่านอย่างน้อย 8 ตัวอักษรและไม่เกิน 72 bytes เก็บ `.jwt_secret` ไว้ใช้ต่อ และอย่าเพิ่มไฟล์นี้ลง Git

## ตรวจงาน

จากโฟลเดอร์หลัก:

```bash
CI=true npm test -- --watchAll=false
CI=true npm run build
cd moviehub
go test ./...
go vet ./...
go build ./...
```

ผลตรวจวันที่ 21 ตุลาคม 2026:

- React: 9 suites, 27 tests ผ่าน และ production build ผ่าน
- Go: ทดสอบ JWT ปลอม, login/role middleware และ CORS ผ่าน; `go vet` และ build ผ่าน
- ทดสอบ API จริง 81 คำขอ: สมัคร/เข้าสู่ระบบ, สิทธิ์เจ้าของรีวิวและ admin, Wishlist, คะแนน, validation, CORS และเวลาตอบกลับกรณีรหัสผ่านผิด/ไม่มีอีเมล
- ทดสอบเบราว์เซอร์จริง 11 กรณี: สมัครสมาชิก, ยืนยันรหัสผ่าน, คะแนน, Wishlist, รีวิว, logout และ CORS จาก origin ที่อนุญาต/ไม่อนุญาต
- ตรวจ API และเบราว์เซอร์ด้วยฐานข้อมูลทดสอบ PostgreSQL 17.11 ชั่วคราว โดยโหลด `schema.sql` และ `seed.sql` เดิม; Docker Compose ของโครงการยังใช้ PostgreSQL 16 ตามบทเรียน

การตรวจในเครื่องนี้ยังไม่ใช่การ deploy Week 11 ไป Supabase/Render
