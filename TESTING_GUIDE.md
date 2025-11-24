# 🧪 Complete Testing Guide - EDS UPI Tabulation System

## 📋 Pre-requisites

✅ Backend running on `http://localhost:8080`
✅ Frontend running on `http://localhost:3000`
✅ PostgreSQL database running
✅ Admin user created (username: `admin`, password: `admin123`)

---

## 🎯 TESTING FLOW (Step by Step)

### **STEP 1: Login** ✓

**URL:** `http://localhost:3000/login`

**Actions:**
1. Enter username: `admin`
2. Enter password: `admin123`
3. Click LOGIN
4. Should redirect to `/admin/dashboard` automatically (no alert)

**Expected Result:**
- ✅ Redirect to dashboard
- ✅ See 6 menu cards
- ✅ No errors in console

**If Failed:**
- Run: `cd Backend && go run cmd/seed.go` to create admin user
- Check backend is running on port 8080

---

### **STEP 2: Create Tournament** ✓

**URL:** `http://localhost:3000/admin/tournaments`

**Actions:**
1. Fill form:
   - Nama Turnamen: `EDS Cup 2025`
   - Format: `Asian Parliamentary`
   - Lokasi: `UPI Bandung`
   - Tanggal Mulai: `2025-01-15`
   - Tanggal Selesai: `2025-01-17`
2. Click `SIMPAN TURNAMEN`

**Expected Result:**
- ✅ Toast notification: "Turnamen berhasil dibuat!"
- ✅ Tournament appears in the list
- ✅ Status badge shows "UPCOMING"
- ✅ Quick action buttons visible

**Test Data Created:**
```json
{
  "name": "EDS Cup 2025",
  "format": "asian",
  "location": "UPI Bandung",
  "start_date": "2025-01-15",
  "end_date": "2025-01-17"
}
```

---

### **STEP 3: Add Teams** ✓

**URL:** `http://localhost:3000/admin/teams`

**Actions:**
Create 4 teams for testing:

**Team 1:**
- Tournament: `EDS Cup 2025`
- Nama Tim: `UPI A`
- Institusi: `Universitas Pendidikan Indonesia`
- Speaker 1: `Budi Santoso`
- Speaker 2: `Ani Wijaya`
- Speaker 3: `Citra Dewi`

**Team 2:**
- Nama Tim: `ITB A`
- Institusi: `Institut Teknologi Bandung`
- Speakers: `Dedi`, `Eka`, `Fani`

**Team 3:**
- Nama Tim: `UNPAD A`
- Institusi: `Universitas Padjadjaran`
- Speakers: `Gita`, `Hadi`, `Ika`

**Team 4:**
- Nama Tim: `UGM A`
- Institusi: `Universitas Gadjah Mada`
- Speakers: `Joko`, `Kiki`, `Lina`

**Expected Result:**
- ✅ Toast: "Tim berhasil didaftarkan!" (4 times)
- ✅ All 4 teams appear in table
- ✅ Speakers listed correctly
- ✅ Counter shows "4 Tim"

---

### **STEP 4: Create Rounds & Matches** ✓

**URL:** `http://localhost:3000/admin/matches`

**Actions:**

#### **4.1 Create Round 1:**
1. Select Tournament: `EDS Cup 2025`
2. Nama Ronde: `Round 1`
3. Mosi: `THW ban social media for teenagers`
4. Click `+ Tambah Ronde`

**Expected Result:**
- ✅ Toast: "Ronde berhasil dibuat!"
- ✅ Round 1 appears in left panel

#### **4.2 Create Match 1:**
1. Click on `Round 1` (should highlight in blue)
2. Fill pairing form:
   - Government: `UPI A`
   - Opposition: `ITB A`
   - Ruangan: `A1`
   - Juri: `Dr. Ahmad`
3. Click `✓ SIMPAN PAIRING`

**Expected Result:**
- ✅ Toast: "Pairing berhasil dibuat!"
- ✅ Match appears in table below
- ✅ Room shows "A1"
- ✅ Status shows "Pending"

#### **4.3 Create Match 2:**
- Gov: `UNPAD A` vs Opp: `UGM A`
- Room: `A2`
- Juri: `Prof. Siti`

**Expected Result:**
- ✅ 2 matches visible in table
- ✅ Different rooms (A1, A2)
- ✅ All teams assigned

---

### **STEP 5: Submit Ballot (Input Scores)** ✓

**URL:** `http://localhost:3000/admin/ballot`

**Actions:**

#### **5.1 Select Match:**
1. Pilih Turnamen: `EDS Cup 2025`
2. Pilih Ronde: `Round 1`
3. Pilih Match: `A1: UPI A vs ITB A`

**Expected Result:**
- ✅ Match details appear
- ✅ Room header shows "ROOM: A1"
- ✅ Teams displayed: "GOV: UPI A vs OPP: ITB A"
- ✅ Juri auto-filled: "Dr. Ahmad"

#### **5.2 Input Scores:**

**Government (UPI A):**
- Speaker 1 (Budi Santoso): 76
- Speaker 2 (Ani Wijaya): 78
- Speaker 3 (Citra Dewi): 75
- Reply: 38

**Opposition (ITB A):**
- Speaker 1 (Dedi): 74
- Speaker 2 (Eka): 73
- Speaker 3 (Fani): 72
- Reply: 36

**Total:**
- Gov: 267
- Opp: 255

3. Click `📥 SUBMIT BALLOT`

**Expected Result:**
- ✅ Toast: "Skor berhasil disimpan!"
- ✅ Redirect to dashboard
- ✅ No errors

---

### **STEP 6: View Standings** ✓

**URL:** `http://localhost:3000/admin/standings`

**Actions:**
1. Select Tournament: `EDS Cup 2025`

**Expected Result:**
- ✅ Table shows all teams
- ✅ UPI A ranked #1 (1 VP, 267 speaker points)
- ✅ ITB A ranked lower (0 VP, 255 speaker points)
- ✅ Win/Loss indicators: UPI A shows "1 ↑ - 0 ↓"
- ✅ ITB A shows "0 ↑ - 1 ↓"
- ✅ Medal icons for top 3

---

### **STEP 7: View Match Results** ✓

**URL:** `http://localhost:3000/admin/results`

**Actions:**
1. Select Tournament: `EDS Cup 2025`
2. Select Round: `Round 1`

**Expected Result:**
- ✅ Motion displayed: "THW ban social media for teenagers"
- ✅ Match A1 shows:
  - Gov: UPI A (with green ↑)
  - Opp: ITB A (no indicator)
  - Winner badge: "UPI A"
  - Status: "✓ COMPLETED"
- ✅ Match A2 shows:
  - Status: "PENDING"

---

## 📊 VERIFICATION CHECKLIST

After completing all steps, verify:

### **Database:**
- [ ] 1 Tournament created
- [ ] 4 Teams created
- [ ] 12 Speakers created (3 per team)
- [ ] 1 Round created
- [ ] 2 Matches created
- [ ] 8 Ballot entries (4 Gov + 4 Opp for Match A1)

### **Frontend:**
- [ ] All pages load without errors
- [ ] Toast notifications work
- [ ] No browser console errors
- [ ] Data displays correctly
- [ ] Navigation works

### **Backend:**
- [ ] Server running without crashes
- [ ] All API endpoints responding
- [ ] No 500 errors
- [ ] Data persisted in database

---

## 🐛 Common Issues & Solutions

### **Issue 1: Cannot create match**
**Solution:** Backend fixed! Restart server if needed.

### **Issue 2: Ballot submission fails**
**Possible causes:**
- Match not created
- Speaker names don't match
- Missing team_role field

**Solution:** Check match exists, verify speaker names exactly match team registration.

### **Issue 3: Standings show 0 VP**
**Possible causes:**
- Ballot not submitted
- Winner calculation failed

**Solution:** Re-submit ballot, check backend logs.

### **Issue 4: Toast not showing**
**Solution:** Check ToastContainer is in layout.tsx, refresh page.

---

## ✅ SUCCESS CRITERIA

All tests pass if:
1. ✅ Can login without errors
2. ✅ Can create tournament with dates
3. ✅ Can add multiple teams
4. ✅ Can create rounds with motions
5. ✅ **Can create matches (CRITICAL - FIXED!)**
6. ✅ Can submit ballots
7. ✅ Standings calculate correctly
8. ✅ Match results display properly

---

## 📝 Test Results Log

**Date:** _____________
**Tester:** _____________

| Step | Status | Notes |
|------|--------|-------|
| 1. Login | ⬜ Pass ⬜ Fail | |
| 2. Create Tournament | ⬜ Pass ⬜ Fail | |
| 3. Add Teams | ⬜ Pass ⬜ Fail | |
| 4. Create Rounds | ⬜ Pass ⬜ Fail | |
| 5. Create Matches | ⬜ Pass ⬜ Fail | |
| 6. Submit Ballot | ⬜ Pass ⬜ Fail | |
| 7. View Standings | ⬜ Pass ⬜ Fail | |
| 8. View Results | ⬜ Pass ⬜ Fail | |

**Overall Result:** ⬜ PASS ⬜ FAIL

**Issues Found:**
_______________________________________
_______________________________________

---

**Next:** After all tests pass, proceed to **OPSI 2** (Room & Adjudicator Management + Auto Pairing)
