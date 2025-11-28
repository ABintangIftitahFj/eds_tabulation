# 🗺️ EDS UPI Tabulation System - Development Roadmap

## ✅ FASE 1: COMPLETED (Sudah Jalan)
- [x] Authentication (Login/Register)
- [x] Tournament CRUD (dengan tanggal)
- [x] Team Registration
- [x] Round Management
- [x] Match/Pairing Creation
- [x] Ballot Submission (Input Skor)
- [x] Standings/Leaderboard (dengan Win/Loss indicator & Power Horse detection)
- [x] Article Management

---

## 🚧 FASE 2: PRIORITY FEATURES (Harus Segera)

### 1. **Tournament Detail Page** (`/admin/tournament/[id]`)
**Status:** Belum dibuat
**Fitur:**
- Overview statistics (Total Teams, Total Rounds, Total Matches)
- Quick actions (Add Round, Add Team, View Standings)
- Tournament settings (Edit name, dates, status)
- Timeline view (Rounds progress)

### 2. **Round Management Enhancement**
**Status:** Partial (Create ada, tapi UI kurang)
**Yang Kurang:**
- [ ] Edit/Delete round
- [ ] Publish/Unpublish round
- [ ] Motion display per round
- [ ] Info slide upload/input
- [ ] Round status indicator

### 3. **Match Results Display**
**Status:** Belum ada
**Yang Dibutuhkan:**
- [ ] List semua match per round dengan hasil
- [ ] Motion yang digunakan (dari round)
- [ ] Winner indicator (Gov/Opp)
- [ ] Link ke detailed ballot/scoresheet
- [ ] Filter by round

### 4. **Speaker Tab** (Individual Rankings)
**Status:** Belum ada
**Yang Dibutuhkan:**
- [ ] Ranking speaker berdasarkan total score
- [ ] Average score per speaker
- [ ] Number of speeches
- [ ] Best speaker badge

---

## 🎯 FASE 3: ADVANCED FEATURES (Nice to Have)

### 5. **Auto Draw/Pairing System**
**Status:** Belum ada (sekarang manual)
**Fitur:**
- [ ] Random pairing untuk preliminary rounds
- [ ] Power-paired draw (berdasarkan standings)
- [ ] Avoid same institution matchup
- [ ] Room allocation otomatis

### 6. **Adjudicator Management**
**Status:** Belum ada
**Fitur:**
- [ ] Register adjudicators
- [ ] Assign adjudicators to rooms
- [ ] Adjudicator feedback/rating
- [ ] Conflict of interest detection

### 7. **Public Pages** (Frontend untuk Peserta)
**Status:** Belum ada
**Fitur:**
- [ ] Public draw display
- [ ] Public standings
- [ ] Motion announcement
- [ ] Schedule/timeline

### 8. **Export & Reports**
**Status:** Belum ada
**Fitur:**
- [ ] Export standings to PDF/Excel
- [ ] Generate certificates
- [ ] Tournament summary report
- [ ] Speaker awards list

---

## 🐛 FASE 4: BUG FIXES & POLISH

### Known Issues:
1. **Backend:**
   - [ ] Lint error: "main redeclared" di seed.go (minor, tidak mengganggu)
   - [ ] Validasi input skor (min/max range)
   - [ ] Handle duplicate team names
   
2. **Frontend:**
   - [ ] Loading states kurang smooth
   - [ ] Error handling belum konsisten
   - [ ] Mobile responsiveness perlu diperbaiki
   - [ ] Toast notifications (ganti alert())

3. **Database:**
   - [ ] Add indexes untuk query performance
   - [ ] Soft delete untuk data penting
   - [ ] Backup strategy

---

## 📝 NOTES & RECOMMENDATIONS

### Prioritas Implementasi (Urutan):
1. **Tournament Detail Page** - Jadi hub utama management
2. **Match Results Display** - Biar bisa lihat hasil pertandingan
3. **Speaker Tab** - Untuk individual awards
4. **Round Management UI** - Biar lebih user-friendly
5. **Auto Pairing** - Menghemat waktu setup

### Tech Debt:
- Perlu refactor ballot_controller.go (terlalu panjang)
- Frontend components bisa di-extract jadi reusable
- Perlu tambah unit tests

### Inspiration dari Tabbycat:
- Draw display dengan room layout
- Motion info slide integration
- Break announcement system
- Feedback forms untuk adjudicators

---

## 🎨 UI/UX Improvements Needed:
- [ ] Dashboard dengan statistics cards
- [ ] Better navigation (sidebar?)
- [ ] Breadcrumbs
- [ ] Confirmation modals (sebelum delete)
- [ ] Success/error toast notifications
- [ ] Dark mode support
- [ ] Print-friendly views

---

**Last Updated:** 2025-11-24
**Version:** 0.2.0-alpha
