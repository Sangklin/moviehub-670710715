# MovieHub Go API

งานตาม [Week 12: The Ticket That Cannot Be Faked](https://plearnjai.com/be2026/Week12-The-Ticket-That-Cannot-Be-Faked.html) และ [Lab: Ticket Booth and Gatekeepers](https://plearnjai.com/be2026/Lab-Week12-Ticket-Booth-and-Gatekeepers.html) เชื่อมกับ React ตาม [Week 9](https://plearnjai.com/fe2026/week9-backend-mock.html) และทดสอบตาม [Week 10](https://plearnjai.com/fe2026/week10-testing.html)

## สถานะตรวจเทียบกับบทเรียน

ผู้ใช้ยืนยันว่าเวอร์ชันเครื่องมือแตกต่างได้ แต่โค้ดและขั้นตอนต้องตรงกับบทเรียน จึงใช้ Go 1.27.1 และ PostgreSQL 17.11 เดิม โดยไม่เปลี่ยน dependency

วันที่ 21 ตุลาคม 2026 ตรวจกรณีของแต่ละด่านซ้ำตามลำดับในใบงาน ใช้ฐานข้อมูลทดสอบใหม่ชื่อ `moviehub` โหลด `schema.sql` และ `seed.sql` เดิม ผลอยู่ใน `course-reference/verification/lab-by-stage.log` (ไฟล์ตรวจใน workspace ไม่ส่งขึ้น Git)

### ส่วนที่มีโจทย์รองรับ

- Week 9 หัวข้อ 7: โหลด Wishlist, ปุ่มคะแนน/หัวใจ และต่อยอดช่องยืนยันรหัสผ่าน
- Week 10 หัวข้อ 9: เทสต์ในคาบ 17 ข้อ, ฝึก Wishlist 2 ข้อ, MovieActions 3 ข้อ และต่อยอด Login รวม 24 ข้อ
- Lab ด่าน 1–5: bcrypt, สมัคร/เข้าสู่ระบบ, JWT, ชุดทดสอบตั๋วปลอม, middleware, รีวิว, admin และโจทย์ Wishlist
- Lab ด่าน 6: CORS, อ่าน error สองชั้นใน `apiFetch` และงานต่อยอดคะแนน
- Lab ด่าน 7: dummy hash ในกรณีไม่พบอีเมล

ไฟล์ Go ที่บทเรียนให้มาทั้งไฟล์และตรวจเทียบตรงกันหลัง `gofmt`: `store/store.go`, `auth/password.go`, `cmd/hashlab/main.go`, `auth/token.go`, `cmd/peek/main.go`, `auth/token_test.go`, `mw/logging.go`, `mw/cors.go`, `handler/admin.go`

### ขั้นตอนที่ตรวจแล้ว

- เตรียม: สร้างฐานข้อมูล โหลด schema/seed หนัง 100 เรื่อง สมาชิกเริ่มต้น 0 คน, build, ค้น `evil` ได้ 2 เรื่อง, รีวิวว่าง และ id ผิดได้ 400
- ด่าน 1: `go run ./cmd/hashlab`, สมัคร 201, อีเมลซ้ำ 409, รหัสสั้น/เกิน 72 bytes ได้ 400, ฐานข้อมูลเก็บ bcrypt และ role เป็น member
- ด่าน 2: ไม่มีกุญแจแล้ว server ปฏิเสธการเริ่ม, รหัสผิด/ไม่มีอีเมลตอบเหมือนกัน, อีเมลตัวใหญ่เข้าได้ และ `cmd/peek` อ่าน JWT ได้
- ด่าน 3: `go test ./auth -v` ผ่าน; ถอด `WithValidMethods` แล้ว HS512 หลุด; คืนแล้วถอด `WithExpirationRequired` พบว่าไม่มี exp หลุด; คืนทั้งสองบรรทัดแล้วผ่าน
- ด่าน 4: `/api/me` ได้ 401/200/401/401, รีวิวได้ 400/201/409/404, `memberId` จาก body เปลี่ยนเจ้าของไม่ได้ และ log ไม่มี Bearer
- ด่าน 5: สมาชิกได้ 403, ตั๋ว admin ใบเก่ายังได้ 403/ใบใหม่ 200, ลบรีวิวข้ามเจ้าของ 404, admin/เจ้าของลบได้ 204 และ Wishlist แยกแต่ละคน
- ด่าน 6: สร้าง `.env.local` แล้วหยุด/เริ่ม `npm start`; curl OPTIONS ได้ headers ตามใบงาน; เบราว์เซอร์อ่านข้อมูลนิดาได้และส่ง OPTIONS ก่อน DELETE; ถอด Authorization แล้วเบราว์เซอร์ถูกบล็อกแต่ curl ผ่าน; ใส่คืนแล้ว origin 3001 ยังถูกบล็อก; Login แสดง error ภาษาไทย; ตรวจงานต่อยอดคะแนน 200/400/204
- ด่าน 7: curl ห้าครั้งต่อกรณีก่อนและหลัง dummy hash: ก่อนแก้ รหัสผิดประมาณ 67–70 ms/ไม่มีอีเมลต่ำกว่า 1 ms; หลังแก้ทั้งคู่ใช้เวลาประมาณ 67–98 ms และรหัสถูกยังได้ 200
- Week 9: เปิด `npm run mock` และ `npm start`, หนัง 100 เรื่อง, login ผิด/ถูก, protected route, รีวิว 400/201/409 และ Wishlist
- Week 10: ทดลองเปลี่ยนความยาวเป็น 30 และถอด Content-Type ให้เทสต์แดง แล้วคืนให้เขียว; ชุดในคาบ 5 suites / 17 tests ผ่าน; รวมฝึกและต่อยอด 8 suites / 24 tests ผ่าน

การทดลองถอดเงื่อนไข Go ทำในสำเนาที่ `course-reference/lab-runtime` เพื่อไม่ทิ้งโค้ดที่ถอดเงื่อนไขไว้ในโครงการ หลังจบด่าน 7 ตรวจว่าไฟล์ Go ในสำเนาตรงกับโครงการทุกไฟล์

### URL ที่ผู้ใช้ยืนยัน

ผู้ใช้เปิด React ที่ `http://127.0.0.1:3000/` จึงตั้ง `.env.local` ที่ราก React เป็น `REACT_APP_API_URL=http://127.0.0.1:8080` และรัน Go ด้วย `ALLOWED_ORIGINS=http://127.0.0.1:3000` แล้วหยุด/เริ่มทั้งสองฝั่งใหม่ตามด่าน 6

ตรวจด้วยเบราว์เซอร์ภายใน Codespace ที่ URL นี้แล้ว: หนัง 100 เรื่อง, Login รหัสผิดแสดงข้อความภาษาไทย, รหัสถูกเข้าสู่ระบบได้ และ `/api/me` ได้ 200 การเปิดจากเครื่องผู้ใช้ต้องเข้าถึงพอร์ต 3000 และ 8080 ของ Codespace ได้ด้วย ยังไม่ถือว่าผลจากเบราว์เซอร์ภายในเป็นหลักฐานหน้าจอของผู้ใช้

### หลักฐานที่ยังต้องจัดทำตามเว็บ

- Week 9: คลิปไม่เกิน 1 นาที เห็นหนัง 100 เรื่องและ log ของ **mock** ตอน login, ส่งรีวิว และเพิ่ม Wishlist
- Week 10: ภาพ terminal ของผู้ใช้ที่เห็น 5 suites / 17 tests สำหรับส่วนในคาบ ส่วนฝึกและต่อยอดรวมเป็น 8 suites / 24 tests
- Lab ด่าน 3: ภาพรวมผลทดสอบ JWT ทั้งสามรอบตามใบงาน (ปกติ, ถอดการล็อกวิธีเซ็น, ถอดการบังคับวันหมดอายุ) แล้วคืนโค้ดให้ครบ

ห้ามใช้ภาพหรือผลทดสอบที่สร้างขึ้นแทนหลักฐานจากการรันจริง และยังไม่ถือว่า deploy Week 11 ไป Supabase/Render แล้ว
