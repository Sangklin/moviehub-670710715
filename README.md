# 🎬 MovieHub Starter

โปรเจกต์เริ่มต้นสำหรับสัปดาห์ Fetch / JSON / API ของวิชา 520 341 Client Side Web Programming
หน้าตาเสร็จแล้วทั้งหมด (React 19 + React Router 6 + Tailwind CSS 3) เหลือเติมเฉพาะส่วนที่คุยกับ API ทีละขั้นในคาบ

## สถานะ Week 12: Admin Backoffice

- เว็บ React: https://moviehub-670710715-ianr.vercel.app/
- Go API: https://moviehub-51xl.onrender.com
- Repository React และ Go อยู่ด้วยกัน: https://github.com/Sangklin/moviehub-670710715 (Go อยู่ใน `moviehub/`)
- GitHub Actions: https://github.com/Sangklin/moviehub-670710715/actions
- บทเรียน: https://plearnjai.com/fe2026/week12-admin-backoffice.html

หลังบ้านอยู่ที่ `/admin` มี Dashboard, จัดการหนัง และรายชื่อสมาชิก สมาชิกทั่วไปเห็น 403; admin เข้าจัดการได้ การตรวจสิทธิ์จริงอยู่ที่ Go และต้อง login ใหม่หลังเปลี่ยนบทบาท

ตรวจด่านเพิ่ม/แก้ไข/ลบหนังผ่านหน้าเว็บแล้ว และเก็บหนังสั้นสมมติ “คืนฉายหนังที่ทับแก้ว” ไว้เป็นข้อมูลฝึก การเลื่อนเพื่อนเป็น admin ใน Lab ข้อ 7 ข้ามไว้ตามคำสั่งเจ้าของโปรเจกต์

รัน React จากโฟลเดอร์หลักด้วย `npm start` โดยตั้ง `REACT_APP_API_URL=https://moviehub-51xl.onrender.com` ใน `.env.local` และใช้ `REACT_APP_TMDB_KEY` ของตนเองใน `.env` สำหรับส่วน TMDB ส่วน Go ใช้ `cd moviehub` แล้ว `go run .` โดยอ่าน `DATABASE_URL`, `JWT_SECRET` และ `ALLOWED_ORIGINS` จาก `moviehub/.env` ในเครื่อง หรือ Environment Variables บน Render

ตรวจงานด้วย `CI=true npm test -- --watchAll=false --runInBand` และ `npm run build` เทสต์ทั้งหมด 30 ข้อ รวม RequireRole และ AdminMembers ตามบทเรียน 6 ข้อ

ข้อความ Starter ด้านล่างเป็นเอกสารเดิมของช่วงเริ่มต่อ TMDB; หน้า Movies ปัจจุบันใช้ Go API แล้ว

## เริ่มใช้งานใน Codespaces

```bash
npm install               # Codespaces รันให้อัตโนมัติตอนสร้างครั้งแรก รันซ้ำได้ไม่เสียหาย
cp .env.example .env      # แล้วเปิด .env ใส่ key ของตัวเองหลังเครื่องหมาย =
npm start
```

ถ้ายังไม่ได้ใส่ key (หรือใส่ผิดรูปแบบ) จะเห็นแถบสีเหลืองเตือนอยู่ใต้เมนู ใส่แล้วต้อง **หยุดแล้วรัน npm start ใหม่** แถบถึงจะหาย

## โครงสร้างไฟล์

```
src/
  api/tmdb.js              ชั้นกลางคุยกับ TMDB (มี TODO ขั้นที่ 1 และ 2)
  api/cache.js             onceADay จำผลไว้ใน localStorage วันละครั้ง (มีให้แล้ว)
  components/
    Navbar.jsx  Footer.jsx  KeyBanner.jsx     โครงหน้า
    MovieCard.jsx  MovieGrid.jsx              การ์ดหนัง + การวาดตามสถานะ loading / error / success
    FeaturedCarousel.jsx                      แถบหนังแนะนำเลื่อนซ้ายขวา
    ReviewForm.jsx                            ฟอร์มรีวิว (state 3 ตัว)
  pages/
    Home.jsx                หน้าแรก หนังแนะนำสุ่ม (มี TODO ขั้นที่ 5)
    Movies.jsx              หน้ารายการหนัง ค้นหา + เลือกแนว กรองในเครื่อง (มี TODO ขั้นที่ 3)
    MovieDetail.jsx         หน้ารายละเอียด (มี TODO ขั้นที่ 4)
    ApiLab.jsx              สนามทดลองยิง API ดู JSON ดิบ และเทียบความเร็วเมื่อมี cache
    About.jsx  NotFound.jsx
  data/data.js              หนัง 10 เรื่องสำหรับใช้ก่อนต่อ API (id เป็นรหัสเดียวกับ TMDB)
docs/api-wishlist.md      แม่แบบรายการ API ที่จะขอจากทีม Backend (กรอกในขั้นที่ 6)
.env.example                แม่แบบไฟล์ .env
.devcontainer/              ตั้งค่า Codespaces (Node 20 + npm install อัตโนมัติ + เปิดพอร์ต 3000)
vercel.json                 ตั้งค่าให้ refresh หน้าใน ๆ แล้วไม่ 404
```

## สิ่งที่จะทำในคาบ (ค้นหาคำว่า TODO ในโค้ด)

| ขั้น | ไฟล์ | ทำอะไร | ดูผลได้ที่ |
|---|---|---|---|
| 0 | `.env` | ใส่ REACT_APP_TMDB_KEY | แถบเตือนหาย |
| 1 | `src/api/tmdb.js` | เขียนส่วน fetch ใน `getJSON` | หน้า API Lab ปุ่ม "ดู JSON ดิบ" ได้ JSON แทน error |
| 2 | `src/api/tmdb.js` | เขียน `toMovie` แปลง JSON ของ TMDB เป็นรูปร่างเดียวกับ data.js | หน้า API Lab ปุ่ม "ผ่าน toMovie แล้ว" ได้ field ชื่อเดียวกับ data.js |
| 3 | `src/pages/Movies.jsx` | useEffect โหลดหนังกำลังฉายผ่าน `getMovies()` (ยิง API จริงวันละครั้ง) ค้นหาและเลือกแนวกรองในเครื่อง | หน้า "หนัง" ขึ้นหนังจริงประมาณ 60 เรื่อง พิมพ์ค้นหาแล้วไม่มี request |
| 4 | `src/pages/MovieDetail.jsx` | โหลดรายละเอียดตาม id ด้วย `getMovie` | คลิกการ์ดแล้วได้เรื่องย่อภาษาไทย |
| 5 | `src/pages/Home.jsx` | หนังแนะนำสุ่มจาก `getMovies()` ชุดเดียวกับหน้าหนัง | หน้าแรกเลื่อนดูหนังกำลังฉาย ไม่มี request ซ้ำ |
| 6 | `docs/api-wishlist.md` | กรอกแม่แบบรายการ API ที่ทีมต้องการจาก Backend | ส่งให้ทีม Go |

## Deploy ขึ้น Vercel

1. push โค้ดขึ้น GitHub (ตรวจว่า `git status` ไม่มีไฟล์ `.env`)
2. Vercel: Add New Project แล้วเลือก repo ถ้าโปรเจกต์อยู่ในโฟลเดอร์ย่อย ให้ตั้ง Root Directory เป็นโฟลเดอร์นั้น
3. Settings แล้ว Environment Variables เพิ่ม `REACT_APP_TMDB_KEY` ค่าเดียวกับใน .env
4. Deploy (ถ้าเพิ่มตัวแปรทีหลังต้องกด Redeploy)

## เครดิต

ข้อมูลหนังและโปสเตอร์จาก [The Movie Database (TMDB)](https://www.themoviedb.org)
This product uses the TMDB API but is not endorsed or certified by TMDB.
โลโก้ TMDB สำหรับหน้า About ดาวน์โหลดได้จาก https://www.themoviedb.org/about/logos-attribution
