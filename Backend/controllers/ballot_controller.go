package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/star_fj/eds-backend/models"
)

type BallotInput struct {
	MatchID       uint            `json:"match_id"`
	AdjudicatorID uint            `json:"adjudicator_id"`
	Adjudicator   string          `json:"adjudicator"`
	Scores        []models.Ballot `json:"scores"`
	Winner        string          `json:"winner"`    // "gov" or "opp" - explicit winner selection
	GovReply      *int            `json:"gov_reply"` // Optional reply score
	OppReply      *int            `json:"opp_reply"` // Optional reply score
}

func SubmitBallot(c *gin.Context) {
	var input BallotInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Mulai Transaksi Database (Biar Aman)
	tx := models.DB.Begin()

	// 0. Ambil data match terlebih dahulu untuk mendapat teamID
	var match models.Match
	if err := tx.Preload("GovTeam").Preload("OppTeam").First(&match, input.MatchID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Match tidak ditemukan"})
		return
	}

	// 0.5. Jika match sudah pernah di-ballot, revert stats lama dulu
	if match.IsCompleted {
		// Ambil semua ballot lama untuk match ini
		var oldBallots []models.Ballot
		tx.Where("match_id = ?", input.MatchID).Find(&oldBallots)

		// Hitung total skor lama per tim
		var oldGovScore, oldOppScore int
		for _, ballot := range oldBallots {
			if ballot.TeamRole == "gov" {
				oldGovScore += ballot.Score
			} else if ballot.TeamRole == "opp" {
				oldOppScore += ballot.Score
			}
		}

		// Revert speaker scores
		for _, ballot := range oldBallots {
			if ballot.SpeakerID != 0 {
				var speaker models.Speaker
				if err := tx.First(&speaker, ballot.SpeakerID).Error; err == nil {
					speaker.TotalScore -= ballot.Score
					tx.Save(&speaker)
				}
			}
		}

		// Revert team stats
		if match.GovTeam.ID != 0 {
			match.GovTeam.TotalSpeaker -= oldGovScore
			if match.WinnerID != nil && *match.WinnerID == match.GovTeam.ID {
				match.GovTeam.TotalVP -= 1
				match.GovTeam.Wins -= 1
			} else {
				match.GovTeam.Losses -= 1
			}
			tx.Save(&match.GovTeam)
		}
		if match.OppTeam.ID != 0 {
			match.OppTeam.TotalSpeaker -= oldOppScore
			if match.WinnerID != nil && *match.WinnerID == match.OppTeam.ID {
				match.OppTeam.TotalVP -= 1
				match.OppTeam.Wins -= 1
			} else {
				match.OppTeam.Losses -= 1
			}
			tx.Save(&match.OppTeam)
		}

		// Hapus ballot lama
		tx.Where("match_id = ?", input.MatchID).Delete(&models.Ballot{})
	}

	// 1. Simpan Skor Individu
	var totalGov int = 0
	var totalOpp int = 0

	// Track speaker IDs and their scores for later update
	speakerScores := make(map[uint]int)

	for _, ballot := range input.Scores {
		ballot.MatchID = input.MatchID

		// Validasi data pembicara (Jaga-jaga)
		if ballot.Speaker.Name == "" {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Nama speaker tidak boleh kosong"})
			return
		}

		// Tentukan TeamID berdasarkan TeamRole dan MatchID
		var teamID uint
		if ballot.TeamRole == "gov" {
			if match.GovTeamID == nil {
				tx.Rollback()
				c.JSON(http.StatusBadRequest, gin.H{"error": "Match tidak memiliki Government Team"})
				return
			}
			teamID = *match.GovTeamID
		} else if ballot.TeamRole == "opp" {
			if match.OppTeamID == nil {
				tx.Rollback()
				c.JSON(http.StatusBadRequest, gin.H{"error": "Match tidak memiliki Opposition Team"})
				return
			}
			teamID = *match.OppTeamID
		} else {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "TeamRole harus 'gov' atau 'opp'"})
			return
		}

		// Debug log
		fmt.Printf("Debug: TeamRole=%s, TeamID=%d, MatchID=%d, SpeakerID=%d, Score=%d\n", ballot.TeamRole, teamID, ballot.MatchID, ballot.SpeakerID, ballot.Score)

		// Validasi TeamID tidak boleh 0
		if teamID == 0 {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "TeamID tidak valid (0)"})
			return
		}

		// Cari atau buat speaker baru jika belum ada
		var speaker models.Speaker
		var speakerID uint

		if ballot.SpeakerID != 0 {
			// Jika SpeakerID sudah ada, gunakan yang ada
			if err := tx.First(&speaker, ballot.SpeakerID).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusNotFound, gin.H{"error": "Speaker tidak ditemukan"})
				return
			}
			speakerID = speaker.ID
		} else if ballot.Speaker.Name != "" {
			// Cari speaker berdasarkan nama dan tim
			err := tx.Where("name = ? AND team_id = ?", ballot.Speaker.Name, teamID).First(&speaker).Error
			if err != nil {
				// Jika tidak ditemukan, buat speaker baru
				speaker = models.Speaker{
					Name:   ballot.Speaker.Name,
					TeamID: teamID,
				}
				if err := tx.Create(&speaker).Error; err != nil {
					tx.Rollback()
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat speaker: " + err.Error()})
					return
				}
			}
			speakerID = speaker.ID
		} else {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "SpeakerID atau Speaker.Name harus diisi"})
			return
		}

		// Track speaker score for later update
		speakerScores[speakerID] += ballot.Score
		fmt.Printf("Debug: Tracked speaker %d with score %d (total: %d)\n", speakerID, ballot.Score, speakerScores[speakerID])

		// Simpan ballot dengan SpeakerID yang valid
		ballotToSave := models.Ballot{
			MatchID:       ballot.MatchID,
			AdjudicatorID: input.AdjudicatorID,
			SpeakerID:     speakerID,
			Score:         ballot.Score,
			Position:      ballot.Position,
			IsReply:       ballot.IsReply,
			TeamRole:      ballot.TeamRole,
			Winner:        input.Winner,
		}

		if err := tx.Create(&ballotToSave).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal simpan skor: " + err.Error()})
			return
		}

		// Hitung Total Skor
		if ballot.TeamRole == "gov" {
			totalGov += ballot.Score
		} else if ballot.TeamRole == "opp" {
			totalOpp += ballot.Score
		}
	}

	// 2. Tentukan Pemenang (Logic Asian Parliamentary)
	// Match sudah diambil di atas, tidak perlu diambil lagi

	// Update Status Match
	match.IsCompleted = true
	// match.Adjudicator = input.Adjudicator // Removed because Adjudicator is now a relation

	// Logika Penentuan Pemenang
	// Prioritas: Gunakan winner yang dipilih manual oleh adjudicator
	// Kalau tidak ada, gunakan skor tertinggi
	if input.Winner == "gov" {
		match.WinnerID = match.GovTeamID
	} else if input.Winner == "opp" {
		match.WinnerID = match.OppTeamID
	} else if totalGov > totalOpp {
		// Fallback ke skor jika winner tidak diset
		match.WinnerID = match.GovTeamID
	} else {
		match.WinnerID = match.OppTeamID
	}

	if err := tx.Save(&match).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update hasil match"})
		return
	}

	// 3. Update Klasemen Tim (Standings)
	// Tambahkan VP (Victory Point) dan Speaker Score ke masing-masing tim

	// Update Gov Team
	match.GovTeam.TotalSpeaker += totalGov
	if match.WinnerID == match.GovTeamID {
		match.GovTeam.TotalVP += 1 // Menang dapat 1 poin
		match.GovTeam.Wins += 1
	} else {
		match.GovTeam.Losses += 1
	}
	if err := tx.Save(&match.GovTeam).Error; err != nil {
		tx.Rollback()
		return
	}

	// Update Opp Team
	match.OppTeam.TotalSpeaker += totalOpp
	if match.WinnerID == match.OppTeamID {
		match.OppTeam.TotalVP += 1
		match.OppTeam.Wins += 1
	} else {
		match.OppTeam.Losses += 1
	}
	if err := tx.Save(&match.OppTeam).Error; err != nil {
		tx.Rollback()
		return
	}

	// 4. Update Speaker Individual Scores (untuk Speaker Standings)
	fmt.Printf("Debug: Updating %d speakers scores\n", len(speakerScores))
	for speakerID, score := range speakerScores {
		if speakerID != 0 {
			var speaker models.Speaker
			if err := tx.First(&speaker, speakerID).Error; err == nil {
				oldScore := speaker.TotalScore
				speaker.TotalScore += score
				if err := tx.Save(&speaker).Error; err != nil {
					fmt.Printf("Warning: Failed to update speaker %d score: %v\n", speakerID, err)
				} else {
					fmt.Printf("Debug: Updated speaker %d score from %d to %d\n", speakerID, oldScore, speaker.TotalScore)
				}
			} else {
				fmt.Printf("Warning: Speaker %d not found: %v\n", speakerID, err)
			}
		}
	}

	// Selesai!
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{
		"message":   "Skor disimpan & Pemenang ditentukan!",
		"winner_id": match.WinnerID,
		"total_gov": totalGov,
		"total_opp": totalOpp,
	})
}

func GetBallots(c *gin.Context) {
	matchID := c.Query("match_id")
	roundID := c.Query("round_id")

	var ballots []models.Ballot

	if matchID != "" {
		// Query by specific match
		if err := models.DB.Preload("Speaker").Preload("Adjudicator").Where("match_id = ?", matchID).Find(&ballots).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else if roundID != "" {
		// Query by round - get all match IDs first
		var matchIDs []uint
		if err := models.DB.Table("matches").Where("round_id = ?", roundID).Pluck("id", &matchIDs).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get matches: " + err.Error()})
			return
		}

		if len(matchIDs) == 0 {
			// No matches found, return empty array
			c.JSON(http.StatusOK, gin.H{"data": []models.Ballot{}})
			return
		}

		if err := models.DB.Preload("Speaker").Preload("Adjudicator").Where("match_id IN ?", matchIDs).Find(&ballots).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "match_id or round_id is required"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": ballots})
}
