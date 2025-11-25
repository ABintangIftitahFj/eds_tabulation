# 🔒 Fitur Publikasi Draw & Motion

## 📝 Deskripsi
Fitur ini memungkinkan admin untuk **mengontrol visibilitas Draw dan Motion** kepada user. Ketika admin membuat draw, hasilnya tidak langsung terlihat oleh user sampai admin menekan tombol "Publish Draw".

## 🎯 Masalah yang Diselesaikan
Sebelumnya, ketika admin sudah membuat draw/matchmaking, hasilnya **langsung terlihat** oleh semua user meskipun mungkin masih ada kesalahan atau perlu revisi. Dengan fitur ini:

- ✅ Admin bisa membuat draw terlebih dahulu
- ✅ Admin bisa mengecek dan memastikan draw sudah benar
- ✅ Admin bisa menekan tombol **"Lock Draw"** untuk mempublikasikan
- ✅ User hanya melihat draw yang sudah dipublikasikan
- ✅ User melihat pesan "TBA (To Be Announced)" jika draw belum dipublikasikan

## 🏗️ Arsitektur

### Backend (Go)
**File**: `Backend/models/schema.go`
```go
type Round struct {
    // ... fields lainnya
    IsDrawPublished   bool    `json:"is_draw_published"`   // Draw visible to users
    IsMotionPublished bool    `json:"is_motion_published"` // Motion visible to users
}
```

**File**: `Backend/controllers/setup_controller.go`
- `PublishDraw(c *gin.Context)` - Endpoint untuk publish/unpublish draw
- `PublishMotion(c *gin.Context)` - Endpoint untuk publish/unpublish motion

**Routes**: `Backend/main.go`
```go
api.PUT("/rounds/:id/publish-draw", controllers.PublishDraw)
api.PUT("/rounds/:id/publish-motion", controllers.PublishMotion)
```

### Frontend - Admin Panel
**File**: `frontend/app/admin/tournament/[id]/page.tsx`

Admin dapat mengontrol publikasi draw dan motion dengan toggle button:
- **🔒 Unlock Draw** - Tombol untuk mempublikasikan draw
- **✅ Draw Published** - Status ketika draw sudah dipublikasikan
- **🔒 Unlock Motion** - Tombol untuk mempublikasikan motion
- **✅ Motion Published** - Status ketika motion sudah dipublikasikan

### Frontend - User View
**File**: `frontend/app/tournaments/[id]/rounds/[roundId]/page.tsx`

User akan melihat:
- Jika `is_draw_published == false`: Pesan "🔒 Draw Belum Tersedia" dengan badge TBA
- Jika `is_draw_published == true`: Daftar matchmaking lengkap
- Jika `is_motion_published == false`: "🔒 TBA (To Be Announced)" untuk motion
- Jika `is_motion_published == true`: Motion ditampilkan

**File**: `frontend/app/tournaments/[id]/page.tsx`

Di halaman daftar rounds, user akan melihat:
- Badge "✅ Draw Published" jika draw sudah dipublikasikan dengan tombol "Lihat Draw"
- Badge "🔒 Draw Belum Diumumkan" jika draw belum dipublikasikan

## 📊 Alur Kerja (Workflow)

### 1. Admin Membuat Draw
```
Admin Panel > Tournament Detail > Create Matches > Isi pairing
```
- Draw dibuat dengan `is_draw_published = false` (default)
- User belum bisa melihat

### 2. Admin Mempublikasikan Draw
```
Admin Panel > Tournament Detail > Rounds Timeline > "🔒 Unlock Draw"
```
- Admin klik tombol "Unlock Draw"
- Backend update `is_draw_published = true`
- Tombol berubah menjadi "✅ Draw Published"

### 3. User Melihat Draw
```
User View > Tournament Detail > Rounds Tab > "Lihat Draw"
```
- Jika `is_draw_published = true`: User melihat matchmaking
- Jika `is_draw_published = false`: User melihat pesan TBA

## 🔧 API Endpoints

### Publish Draw
```http
PUT /api/rounds/:id/publish-draw
Content-Type: application/json

{
  "is_draw_published": true
}
```

**Response**:
```json
{
  "message": "Draw published to users",
  "data": {
    "ID": 1,
    "name": "Round 1",
    "is_draw_published": true,
    ...
  }
}
```

### Unpublish Draw
```http
PUT /api/rounds/:id/publish-draw
Content-Type: application/json

{
  "is_draw_published": false
}
```

**Response**:
```json
{
  "message": "Draw hidden from users",
  "data": {
    "ID": 1,
    "name": "Round 1",
    "is_draw_published": false,
    ...
  }
}
```

### Publish Motion
```http
PUT /api/rounds/:id/publish-motion
Content-Type: application/json

{
  "is_motion_published": true
}
```

## 🎨 UI/UX Details

### Admin View - Publication Controls
```
📢 Publication Controls
┌─────────────────────┬─────────────────────┐
│  🔒 Unlock Motion   │   🔒 Unlock Draw    │
└─────────────────────┴─────────────────────┘

Setelah dipublikasikan:
┌─────────────────────┬─────────────────────┐
│ ✅ Motion Published │ ✅ Draw Published   │
└─────────────────────┴─────────────────────┘
```

### User View - Draw Locked
```
┌──────────────────────────────────┐
│            🔒                     │
│     Draw Belum Tersedia          │
│                                  │
│  Draw untuk ronde ini belum      │
│  dipublikasikan oleh admin.      │
│  Silakan tunggu pengumuman       │
│  lebih lanjut.                   │
│                                  │
│  🕐 TBA (To Be Announced)        │
└──────────────────────────────────┘
```

### User View - Draw Published
```
Draw Pertandingan
┌──────────────────────────────────┐
│  1  📍 Ruangan A1       ✅ Selesai│
│  ┌────────────┬────────────┐     │
│  │ GOVERNMENT │ OPPOSITION │     │
│  │   UGM A    │   UI A     │     │
│  └────────────┴────────────┘     │
│  👨‍⚖️ JURI: Dr. Ahmad             │
└──────────────────────────────────┘
```

## 🧪 Testing Checklist

### Backend
- [x] Endpoint `PUT /rounds/:id/publish-draw` berfungsi
- [x] Endpoint `PUT /rounds/:id/publish-motion` berfungsi
- [x] Default value `is_draw_published` adalah `false`
- [x] Default value `is_motion_published` adalah `false`

### Frontend - Admin
- [x] Tombol "Unlock Draw" muncul ketika `is_draw_published = false`
- [x] Tombol berubah menjadi "Draw Published" ketika `is_draw_published = true`
- [x] Tombol "Unlock Motion" muncul ketika `is_motion_published = false`
- [x] Tombol berubah menjadi "Motion Published" ketika `is_motion_published = true`
- [x] Alert muncul setelah berhasil publish/unpublish

### Frontend - User
- [x] Draw tidak terlihat ketika `is_draw_published = false`
- [x] Draw terlihat ketika `is_draw_published = true`
- [x] Motion tidak terlihat (TBA) ketika `is_motion_published = false`
- [x] Motion terlihat ketika `is_motion_published = true`
- [x] Badge status draw muncul di daftar rounds

## 🚀 Deployment Notes

1. Database migration sudah otomatis (GORM AutoMigrate)
2. Field `is_draw_published` dan `is_motion_published` akan ditambahkan ke tabel `rounds`
3. Default value adalah `false` untuk semua round yang sudah ada
4. Tidak perlu manual migration

## 🔑 Key Features

1. **Independent Control**: Draw dan Motion bisa dipublikasikan secara terpisah
2. **Real-time Update**: Perubahan status langsung terlihat di user view
3. **Clear Feedback**: User tahu jelas kapan draw belum tersedia vs belum ada match
4. **Admin Friendly**: Toggle button yang mudah digunakan
5. **Reversible**: Admin bisa unpublish jika ada kesalahan

## 📌 Notes

- Fitur ini **tidak mempengaruhi** data match yang sudah ada
- Admin tetap bisa edit/delete match meskipun sudah dipublish
- Hanya visibilitas di user view yang terpengaruh
- Motion dan Draw adalah control yang independen
