package models

import (
	"time"

	"gorm.io/gorm"
)

// ==========================================
// 🔐 CORE: USER & ADMIN
// ==========================================
type User struct {
	gorm.Model
	Username string `gorm:"unique;not null" json:"username"`
	Password string `gorm:"not null" json:"-"`
	Role     string `gorm:"default:'admin'" json:"role"` // "admin", "tabulator"
}

// ==========================================
// 🏆 FITUR 1: COMPANY PROFILE (PORTFOLIO)
// ==========================================

// Member: Data Anggota EDS UPI
type Member struct {
	gorm.Model
	Name        string `json:"name"`
	Role        string `json:"role"`  // e.g. "President"
	Batch       string `json:"batch"` // e.g. "2023"
	Quote       string `json:"quote"` // e.g. "Vini Vidi Vici"
	PhotoURL    string `json:"photo_url"`
	IsStar      bool   `json:"is_star"` // Member of the Month
	Email       string `json:"email"`
	Achievement string `json:"achievement"` // Prestasi highlight
	IsActive    bool   `gorm:"default:true" json:"is_active"`
}

// Article: Berita & Blog
type Article struct {
	gorm.Model
	Title       string    `json:"title"`
	Slug        string    `gorm:"unique" json:"slug"`
	Content     string    `gorm:"type:text" json:"content"`
	Author      string    `json:"author"`
	ImageURL    string    `json:"image_url"`
	Category    string    `json:"category"` // "News", "Event", "Op-Ed"
	PublishDate time.Time `json:"publish_date"`
	IsPublished bool      `gorm:"default:false" json:"is_published"`
	ViewCount   int       `gorm:"default:0" json:"view_count"`
}

// CompetitionHistory: (Portfolio) Lomba luar yang diikuti EDS UPI
// "Kita menang apa aja sih di luar?"
type CompetitionHistory struct {
	gorm.Model
	Name        string `json:"name"`   // e.g. "NUDC 2024"
	Format      string `json:"format"` // "AP", "BP"
	Level       string `json:"level"`  // "National", "International"
	Year        int    `json:"year"`
	Location    string `json:"location"`
	Result      string `json:"result"` // "Champion", "Semi Finalist"
	Description string `json:"description"`
	PhotoURL    string `json:"photo_url"`
	TeamMembers string `json:"team_members"` // Nama anggota tim (text biasa)
	IsHighlight bool   `json:"is_highlight"` // Tampil di Home
}

// Achievement: (Portfolio) Prestasi spesifik untuk Gallery
type Achievement struct {
	gorm.Model
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	PhotoURL    string    `json:"photo_url"`
	Category    string    `json:"category"`
}

// ==========================================
// 🗣️ FITUR 2: SISTEM TABULASI (APP UTAMA)
// ==========================================

// Tournament: Event yang sedang kita host (e.g., EDS CUP)
type Tournament struct {
	gorm.Model
	Name        string    `json:"name"`
	Slug        string    `gorm:"unique" json:"slug"` // e.g. "eds-cup-2025"
	Format      string    `json:"format"`             // "asian", "british"
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Location    string    `json:"location"`
	Description string    `json:"description"`
	Status      string    `gorm:"default:'upcoming'" json:"status"` // upcoming, ongoing, completed
	IsPublic    bool      `gorm:"default:true" json:"is_public"`
}

// Team: Peserta Turnamen
type Team struct {
	gorm.Model
	TournamentID uint       `json:"tournament_id"`
	Tournament   Tournament `json:"tournament"`
	Name         string     `json:"name"`        // "UGM A"
	Institution  string     `json:"institution"` // "Universitas Gadjah Mada"
	Speakers     []Speaker  `json:"speakers"`

	// Statistik Tabulasi (Diupdate tiap ronde)
	TotalVP      int     `gorm:"default:0" json:"total_vp"`      // Victory Points
	TotalSpeaker float64 `gorm:"default:0" json:"total_speaker"` // Total Speaker Score
	Rank         int     `gorm:"default:0" json:"rank"`
	Wins         int     `gorm:"default:0" json:"wins"`
	Losses       int     `gorm:"default:0" json:"losses"`
}

type Speaker struct {
	gorm.Model
	TeamID      uint    `json:"team_id"`
	Name        string  `json:"name"`
	TotalScore  float64 `json:"total_score"`
	SpeakerRank int     `json:"speaker_rank"`
}

type Round struct {
	gorm.Model
	TournamentID uint    `json:"tournament_id"`
	Name         string  `json:"name"`   // "Round 1"
	Motion       string  `json:"motion"` // "THW Ban TikTok"
	InfoSlide    string  `json:"info_slide"`
	IsPublished  bool    `json:"is_published"`
	Matches      []Match `json:"matches"`
}

// Match: Struktur Hybrid (Bisa AP, Bisa BP)
// Kita pakai teknik "Nullable Foreign Keys"
type Match struct {
	gorm.Model
	RoundID     uint   `json:"round_id"`
	Room        string `json:"room"`
	Adjudicator string `json:"adjudicator"`
	PanelJudges string `json:"panel_judges"`

	// --- KOLOM ASIAN PARLIAMENTARY (2 Teams) ---
	GovTeamID *uint `json:"gov_team_id"` // Pake Pointer (*) biar bisa NULL
	GovTeam   *Team `json:"gov_team"`
	OppTeamID *uint `json:"opp_team_id"`
	OppTeam   *Team `json:"opp_team"`
	// Hasil AP
	WinnerID *uint `json:"winner_id"` // Siapa yang menang (Gov/Opp)

	// --- KOLOM BRITISH PARLIAMENTARY (4 Teams) ---
	OGTeamID *uint `json:"og_team_id"` // Opening Gov
	OGTeam   *Team `json:"og_team"`
	OOTeamID *uint `json:"oo_team_id"` // Opening Opp
	OOTeam   *Team `json:"oo_team"`
	CGTeamID *uint `json:"cg_team_id"` // Closing Gov
	CGTeam   *Team `json:"cg_team"`
	COTeamID *uint `json:"co_team_id"` // Closing Opp
	COTeam   *Team `json:"co_team"`
	// Hasil BP (Ranking 1-4)
	Rank1TeamID *uint `json:"rank1_team_id"`
	Rank2TeamID *uint `json:"rank2_team_id"`
	Rank3TeamID *uint `json:"rank3_team_id"`
	Rank4TeamID *uint `json:"rank4_team_id"`

	IsCompleted bool `json:"is_completed"`
}

// Ballot: Lembar Skor Individu
type Ballot struct {
	gorm.Model
	MatchID   uint    `json:"match_id"`
	SpeakerID uint    `json:"speaker_id"`
	Speaker   Speaker `json:"speaker"`
	Score     float64 `json:"score"` // AP (68-82), BP (60-80)

	// Identitas Peran (Penting buat BP)
	Position string `json:"position"` // "PM", "LO", "Member", "Whip"
	IsReply  bool   `json:"is_reply"`
	TeamRole string `json:"team_role"` // "gov" or "opp"
}
