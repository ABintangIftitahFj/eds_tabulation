package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/star_fj/eds-backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 1. Connect to Database
	dsn := "host=localhost user=admin password=password123 dbname=eds_upi port=5433 sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}
	fmt.Println("✅ Connected to database")

	// 2. Create Tournament
	tournament := models.Tournament{
		Name:        "eds upi testing",
		Slug:        "eds-upi-testing",
		Format:      "asian",
		StartDate:   time.Now(),
		EndDate:     time.Now().AddDate(0, 0, 2),
		Location:    "UPI Bandung",
		Description: "Testing tournament seeded by AI",
		Status:      "completed", // We will simulate it to completion
		IsPublic:    true,
	}

	// Check if exists
	var existingTournament models.Tournament
	if err := db.Where("slug = ?", tournament.Slug).First(&existingTournament).Error; err == nil {
		fmt.Println("⚠️ Tournament 'eds upi testing' already exists. Cleaning up old data...")
		// Delete related data to start fresh
		db.Where("tournament_id = ?", existingTournament.ID).Delete(&models.Adjudicator{})
		db.Where("tournament_id = ?", existingTournament.ID).Delete(&models.Team{})
		// Note: Deleting rounds/matches/ballots cascades usually, but let's be safe or just delete the tournament
		db.Unscoped().Delete(&existingTournament)
	}

	if err := db.Create(&tournament).Error; err != nil {
		log.Fatal("❌ Failed to create tournament:", err)
	}
	fmt.Printf("✅ Created Tournament: %s (ID: %d)\n", tournament.Name, tournament.ID)

	// 3. Create Adjudicators (Simulating "eds 2025" style names)
	adjNames := []string{"Jane Doe", "John Smith", "Alice Johnson", "Bob Brown", "Charlie Davis", "Diana Evans", "Frank Green", "Grace Hall", "Henry Hill", "Ivy Ivy"}
	var adjudicators []models.Adjudicator
	for _, name := range adjNames {
		adj := models.Adjudicator{
			TournamentID: tournament.ID,
			Name:         name,
			Institution:  "Independent",
			Level:        "A",
		}
		db.Create(&adj)
		adjudicators = append(adjudicators, adj)
	}
	fmt.Printf("✅ Created %d Adjudicators\n", len(adjudicators))

	// 4. Create Teams and Speakers (10 teams, 3 speakers each)
	teamNames := []string{"Alpha", "Bravo", "Charlie", "Delta", "Echo", "Foxtrot", "Golf", "Hotel", "India", "Juliet"}
	var teams []models.Team
	for _, name := range teamNames {
		team := models.Team{
			TournamentID: tournament.ID,
			Name:         name,
			Institution:  "UPI",
		}
		db.Create(&team)
		teams = append(teams, team)

		// Create 3 speakers for each team
		for j := 1; j <= 3; j++ {
			speaker := models.Speaker{
				TeamID: team.ID,
				Name:   fmt.Sprintf("%s Speaker %d", name, j),
			}
			db.Create(&speaker)
		}
		fmt.Printf("  -> Created Team: %s\n", name)
	}

	// Helper function to create a match
	createMatch := func(roundID uint, govTeam, oppTeam models.Team, adj models.Adjudicator) models.Match {
		match := models.Match{
			RoundID:       roundID,
			GovTeamID:     &govTeam.ID,
			OppTeamID:     &oppTeam.ID,
			AdjudicatorID: &adj.ID,
			IsCompleted:   true,
		}

		// Determine winner (randomly for simulation, or fixed)
		// Let's make Gov win if index is even, Opp if odd, to mix it up
		winner := govTeam
		loser := oppTeam
		if rand.Intn(2) == 0 {
			winner = oppTeam
			loser = govTeam
		}
		match.WinnerID = &winner.ID

		db.Create(&match)

		// Create Ballots
		// Gov Speakers
		var govSpeakers []models.Speaker
		db.Where("team_id = ?", govTeam.ID).Find(&govSpeakers)
		for k, sp := range govSpeakers {
			score := 75.0 + rand.Float64()*5 // 75-80
			if winner.ID == govTeam.ID {
				score += 2 // Bonus for winning team
			}
			db.Create(&models.Ballot{
				MatchID:   match.ID,
				SpeakerID: sp.ID,
				Score:     score,
				TeamRole:  "gov",
				Position:  fmt.Sprintf("Speaker %d", k+1),
			})
		}

		// Opp Speakers
		var oppSpeakers []models.Speaker
		db.Where("team_id = ?", oppTeam.ID).Find(&oppSpeakers)
		for k, sp := range oppSpeakers {
			score := 75.0 + rand.Float64()*5
			if winner.ID == oppTeam.ID {
				score += 2
			}
			db.Create(&models.Ballot{
				MatchID:   match.ID,
				SpeakerID: sp.ID,
				Score:     score,
				TeamRole:  "opp",
				Position:  fmt.Sprintf("Speaker %d", k+1),
			})
		}

		// Update Team Stats
		db.Model(&winner).Update("wins", gorm.Expr("wins + ?", 1))
		db.Model(&winner).Update("total_vp", gorm.Expr("total_vp + ?", 1))
		db.Model(&loser).Update("losses", gorm.Expr("losses + ?", 1))

		return match
	}

	// 5. Round 1
	round1 := models.Round{
		TournamentID:      tournament.ID,
		Name:              "Round 1",
		Motion:            "THW Ban TikTok",
		IsPublished:       true,
		IsDrawPublished:   true,
		IsMotionPublished: true,
	}
	db.Create(&round1)
	fmt.Println("✅ Created Round 1")

	// Pairings R1 (Random 0v1, 2v3, ...)
	for i := 0; i < 10; i += 2 {
		createMatch(round1.ID, teams[i], teams[i+1], adjudicators[i/2])
	}

	// 6. Round 2
	round2 := models.Round{
		TournamentID:      tournament.ID,
		Name:              "Round 2",
		Motion:            "THW Support Universal Basic Income",
		IsPublished:       true,
		IsDrawPublished:   true,
		IsMotionPublished: true,
	}
	db.Create(&round2)
	fmt.Println("✅ Created Round 2")

	// Fetch updated teams to sort by wins/scores (simple mock sort)
	var updatedTeams []models.Team
	db.Order("wins desc").Order("total_vp desc").Find(&updatedTeams)

	// Pairings R2 (Power paired: 0v1, 2v3...)
	for i := 0; i < 10; i += 2 {
		createMatch(round2.ID, updatedTeams[i], updatedTeams[i+1], adjudicators[i/2])
	}

	// 7. Semi Final (Top 4)
	db.Order("wins desc").Order("total_vp desc").Find(&updatedTeams)
	top4 := updatedTeams[:4]

	semiFinal := models.Round{
		TournamentID:      tournament.ID,
		Name:              "Semi Final",
		Motion:            "THW Abolish the UN Security Council",
		IsPublished:       true,
		IsDrawPublished:   true,
		IsMotionPublished: true,
	}
	db.Create(&semiFinal)
	fmt.Println("✅ Created Semi Final")

	// SF 1: Rank 1 vs Rank 4
	matchSF1 := createMatch(semiFinal.ID, top4[0], top4[3], adjudicators[0])
	// SF 2: Rank 2 vs Rank 3
	matchSF2 := createMatch(semiFinal.ID, top4[1], top4[2], adjudicators[1])

	// 8. Final
	var finalist1, finalist2 models.Team
	// We need to fetch the winners of the SF matches
	// Since createMatch updates the DB, we can just check who won in the struct returned,
	// but the struct returned by createMatch might not have the full updated data if we didn't reload it.
	// However, createMatch sets WinnerID.

	db.First(&finalist1, *matchSF1.WinnerID)
	db.First(&finalist2, *matchSF2.WinnerID)

	finalRound := models.Round{
		TournamentID:      tournament.ID,
		Name:              "Grand Final",
		Motion:            "THBT AI is a threat to humanity",
		IsPublished:       true,
		IsDrawPublished:   true,
		IsMotionPublished: true,
	}
	db.Create(&finalRound)
	fmt.Println("✅ Created Grand Final")

	finalMatch := createMatch(finalRound.ID, finalist1, finalist2, adjudicators[2])

	var champion models.Team
	db.First(&champion, *finalMatch.WinnerID)

	fmt.Printf("\n🏆 TOURNAMENT COMPLETED! 🏆\n")
	fmt.Printf("🥇 Champion: %s\n", champion.Name)
}
