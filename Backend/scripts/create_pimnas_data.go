package main

import (
	"fmt"
	"log"

	"github.com/star_fj/eds-backend/models"
)

// Script to create PIMNAS test tournament data
func main() {
	// Connect to database
	models.ConnectDatabase()

	fmt.Println("🚀 Creating PIMNAS 37 Test Tournament...")

	// Create Tournament
	tournament := models.Tournament{
		Name:        "PIMNAS 37",
		Slug:        "pimnas-37",
		Format:      "asian",
		Location:    "Universitas Pendidikan Indonesia, Bandung",
		Description: "Pekan Ilmiah Mahasiswa Nasional ke-37",
		Status:      "completed",
		IsPublic:    true,
	}
	if err := models.DB.Create(&tournament).Error; err != nil {
		log.Fatal("Failed to create tournament:", err)
	}
	fmt.Printf("✅ Tournament created: %s (ID: %d)\n", tournament.Name, tournament.ID)

	// Create Teams
	teams := []models.Team{
		{TournamentID: tournament.ID, Name: "UPI A", Institution: "Universitas Pendidikan Indonesia", TotalVP: 5, TotalSpeaker: 575, Wins: 5, Losses: 0},
		{TournamentID: tournament.ID, Name: "ITB A", Institution: "Institut Teknologi Bandung", TotalVP: 4, TotalSpeaker: 562, Wins: 4, Losses: 1},
		{TournamentID: tournament.ID, Name: "UI A", Institution: "Universitas Indonesia", TotalVP: 4, TotalSpeaker: 558, Wins: 4, Losses: 1},
		{TournamentID: tournament.ID, Name: "UGM A", Institution: "Universitas Gadjah Mada", TotalVP: 3, TotalSpeaker: 545, Wins: 3, Losses: 2},
		{TournamentID: tournament.ID, Name: "UNPAD A", Institution: "Universitas Padjadjaran", TotalVP: 3, TotalSpeaker: 540, Wins: 3, Losses: 2},
		{TournamentID: tournament.ID, Name: "ITS A", Institution: "Institut Teknologi Sepuluh Nopember", TotalVP: 2, TotalSpeaker: 528, Wins: 2, Losses: 3},
		{TournamentID: tournament.ID, Name: "UNAIR A", Institution: "Universitas Airlangga", TotalVP: 1, TotalSpeaker: 515, Wins: 1, Losses: 4},
		{TournamentID: tournament.ID, Name: "UNDIP A", Institution: "Universitas Diponegoro", TotalVP: 0, TotalSpeaker: 502, Wins: 0, Losses: 5},
	}

	for i := range teams {
		if err := models.DB.Create(&teams[i]).Error; err != nil {
			log.Fatal("Failed to create team:", err)
		}
	}
	fmt.Printf("✅ Created %d teams\n", len(teams))

	// Create Speakers (2 per team)
	speakers := []models.Speaker{
		{TeamID: teams[0].ID, Name: "Ahmad Rifai", TotalScore: 388},
		{TeamID: teams[0].ID, Name: "Siti Nurhaliza", TotalScore: 387},
		{TeamID: teams[1].ID, Name: "Budi Santoso", TotalScore: 381},
		{TeamID: teams[1].ID, Name: "Dewi Lestari", TotalScore: 381},
		{TeamID: teams[2].ID, Name: "Cahya Prasetyo", TotalScore: 377},
		{TeamID: teams[2].ID, Name: "Eka Putri", TotalScore: 371},
		{TeamID: teams[3].ID, Name: "Fajar Nugroho", TotalScore: 377},
		{TeamID: teams[3].ID, Name: "Gita Savitri", TotalScore: 348},
		{TeamID: teams[4].ID, Name: "Hendra Wijaya", TotalScore: 366},
		{TeamID: teams[4].ID, Name: "Indah Permata", TotalScore: 356},
		{TeamID: teams[5].ID, Name: "Joko Widodo", TotalScore: 359},
		{TeamID: teams[5].ID, Name: "Kartika Sari", TotalScore: 337},
		{TeamID: teams[6].ID, Name: "Luthfi Hakim", TotalScore: 354},
		{TeamID: teams[6].ID, Name: "Maya Anggraini", TotalScore: 338},
		{TeamID: teams[7].ID, Name: "Nanda Prakoso", TotalScore: 362},
		{TeamID: teams[7].ID, Name: "Olivia Situmorang", TotalScore: 282},
	}

	for i := range speakers {
		if err := models.DB.Create(&speakers[i]).Error; err != nil {
			log.Fatal("Failed to create speaker:", err)
		}
	}
	fmt.Printf("✅ Created %d speakers\n", len(speakers))

	// Create Adjudicators
	adjudicators := []models.Adjudicator{
		{TournamentID: tournament.ID, Name: "Dr. Andi Setiawan", Institution: "Universitas Indonesia", Level: "Chief"},
		{TournamentID: tournament.ID, Name: "Prof. Benny Kurniawan", Institution: "Institut Teknologi Bandung", Level: "Chief"},
		{TournamentID: tournament.ID, Name: "Citra Maharani, M.A.", Institution: "Universitas Gadjah Mada", Level: "Wing"},
		{TournamentID: tournament.ID, Name: "Doni Prasetya, S.S.", Institution: "Universitas Pendidikan Indonesia", Level: "Wing"},
		{TournamentID: tournament.ID, Name: "Eka Wulandari", Institution: "Universitas Padjadjaran", Level: "Panelist"},
	}

	for i := range adjudicators {
		if err := models.DB.Create(&adjudicators[i]).Error; err != nil {
			log.Fatal("Failed to create adjudicator:", err)
		}
	}
	fmt.Printf("✅ Created %d adjudicators\n", len(adjudicators))

	// Create Rooms
	rooms := []models.Room{
		{TournamentID: tournament.ID, Name: "R1", Location: "Gedung A Lantai 1", Capacity: 50},
		{TournamentID: tournament.ID, Name: "R2", Location: "Gedung A Lantai 2", Capacity: 50},
		{TournamentID: tournament.ID, Name: "R3", Location: "Gedung B Lantai 1", Capacity: 50},
		{TournamentID: tournament.ID, Name: "R4", Location: "Gedung B Lantai 2", Capacity: 50},
		{TournamentID: tournament.ID, Name: "Final Room", Location: "Auditorium Utama", Capacity: 200},
	}

	for i := range rooms {
		if err := models.DB.Create(&rooms[i]).Error; err != nil {
			log.Fatal("Failed to create room:", err)
		}
	}
	fmt.Printf("✅ Created %d rooms\n", len(rooms))

	// Create Rounds and Matches
	rounds := []struct {
		Name      string
		Motion    string
		InfoSlide string
		Matches   []struct {
			RoomIdx int
			AdjIdx  int
			GovIdx  int
			OppIdx  int
			WinIdx  int
			Ballots []struct {
				SpeakerIdx int
				Score      int
				Position   string
				TeamRole   string
			}
		}
	}{
		{
			Name:      "Round 1",
			Motion:    "THW ban social media platforms from using algorithmic content recommendation",
			InfoSlide: "Context: Social media algorithms prioritize engagement over wellbeing",
			Matches: []struct {
				RoomIdx int
				AdjIdx  int
				GovIdx  int
				OppIdx  int
				WinIdx  int
				Ballots []struct {
					SpeakerIdx int
					Score      int
					Position   string
					TeamRole   string
				}
			}{
				{RoomIdx: 0, AdjIdx: 0, GovIdx: 0, OppIdx: 7, WinIdx: 0, Ballots: []struct {
					SpeakerIdx int
					Score      int
					Position   string
					TeamRole   string
				}{
					{SpeakerIdx: 0, Score: 76, Position: "PM", TeamRole: "gov"},
					{SpeakerIdx: 1, Score: 74, Position: "DPM", TeamRole: "gov"},
					{SpeakerIdx: 14, Score: 72, Position: "LO", TeamRole: "opp"},
					{SpeakerIdx: 15, Score: 68, Position: "DLO", TeamRole: "opp"},
				}},
				{RoomIdx: 1, AdjIdx: 1, GovIdx: 1, OppIdx: 6, WinIdx: 1, Ballots: []struct {
					SpeakerIdx int
					Score      int
					Position   string
					TeamRole   string
				}{
					{SpeakerIdx: 2, Score: 75, Position: "PM", TeamRole: "gov"},
					{SpeakerIdx: 3, Score: 73, Position: "DPM", TeamRole: "gov"},
					{SpeakerIdx: 12, Score: 71, Position: "LO", TeamRole: "opp"},
					{SpeakerIdx: 13, Score: 69, Position: "DLO", TeamRole: "opp"},
				}},
				{RoomIdx: 2, AdjIdx: 2, GovIdx: 2, OppIdx: 5, WinIdx: 2, Ballots: []struct {
					SpeakerIdx int
					Score      int
					Position   string
					TeamRole   string
				}{
					{SpeakerIdx: 4, Score: 74, Position: "PM", TeamRole: "gov"},
					{SpeakerIdx: 5, Score: 72, Position: "DPM", TeamRole: "gov"},
					{SpeakerIdx: 10, Score: 70, Position: "LO", TeamRole: "opp"},
					{SpeakerIdx: 11, Score: 68, Position: "DLO", TeamRole: "opp"},
				}},
				{RoomIdx: 3, AdjIdx: 3, GovIdx: 3, OppIdx: 4, WinIdx: 3, Ballots: []struct {
					SpeakerIdx int
					Score      int
					Position   string
					TeamRole   string
				}{
					{SpeakerIdx: 6, Score: 73, Position: "PM", TeamRole: "gov"},
					{SpeakerIdx: 7, Score: 71, Position: "DPM", TeamRole: "gov"},
					{SpeakerIdx: 8, Score: 72, Position: "LO", TeamRole: "opp"},
					{SpeakerIdx: 9, Score: 70, Position: "DLO", TeamRole: "opp"},
				}},
			},
		},
		{
			Name:      "Round 2",
			Motion:    "THW require all schools to teach financial literacy from elementary level",
			InfoSlide: "Context: Many adults struggle with basic financial management",
			Matches: []struct {
				RoomIdx int
				AdjIdx  int
				GovIdx  int
				OppIdx  int
				WinIdx  int
				Ballots []struct {
					SpeakerIdx int
					Score      int
					Position   string
					TeamRole   string
				}
			}{
				{RoomIdx: 0, AdjIdx: 1, GovIdx: 2, OppIdx: 3, WinIdx: 2, Ballots: []struct {
					SpeakerIdx int
					Score      int
					Position   string
					TeamRole   string
				}{
					{SpeakerIdx: 4, Score: 77, Position: "PM", TeamRole: "gov"},
					{SpeakerIdx: 5, Score: 75, Position: "DPM", TeamRole: "gov"},
					{SpeakerIdx: 6, Score: 74, Position: "LO", TeamRole: "opp"},
					{SpeakerIdx: 7, Score: 70, Position: "DLO", TeamRole: "opp"},
				}},
				{RoomIdx: 1, AdjIdx: 2, GovIdx: 0, OppIdx: 1, WinIdx: 0, Ballots: []struct {
					SpeakerIdx int
					Score      int
					Position   string
					TeamRole   string
				}{
					{SpeakerIdx: 0, Score: 78, Position: "PM", TeamRole: "gov"},
					{SpeakerIdx: 1, Score: 76, Position: "DPM", TeamRole: "gov"},
					{SpeakerIdx: 2, Score: 75, Position: "LO", TeamRole: "opp"},
					{SpeakerIdx: 3, Score: 73, Position: "DLO", TeamRole: "opp"},
				}},
				{RoomIdx: 2, AdjIdx: 3, GovIdx: 4, OppIdx: 5, WinIdx: 4, Ballots: []struct {
					SpeakerIdx int
					Score      int
					Position   string
					TeamRole   string
				}{
					{SpeakerIdx: 8, Score: 76, Position: "PM", TeamRole: "gov"},
					{SpeakerIdx: 9, Score: 74, Position: "DPM", TeamRole: "gov"},
					{SpeakerIdx: 10, Score: 73, Position: "LO", TeamRole: "opp"},
					{SpeakerIdx: 11, Score: 71, Position: "DLO", TeamRole: "opp"},
				}},
				{RoomIdx: 3, AdjIdx: 4, GovIdx: 6, OppIdx: 7, WinIdx: 7, Ballots: []struct {
					SpeakerIdx int
					Score      int
					Position   string
					TeamRole   string
				}{
					{SpeakerIdx: 12, Score: 70, Position: "PM", TeamRole: "gov"},
					{SpeakerIdx: 13, Score: 68, Position: "DPM", TeamRole: "gov"},
					{SpeakerIdx: 14, Score: 72, Position: "LO", TeamRole: "opp"},
					{SpeakerIdx: 15, Score: 70, Position: "DLO", TeamRole: "opp"},
				}},
			},
		},
		{
			Name:      "Round 3",
			Motion:    "TH regrets the rise of gig economy",
			InfoSlide: "Context: Freelance and contract work has become increasingly common",
			Matches: []struct {
				RoomIdx int
				AdjIdx  int
				GovIdx  int
				OppIdx  int
				WinIdx  int
				Ballots []struct {
					SpeakerIdx int
					Score      int
					Position   string
					TeamRole   string
				}
			}{
				{RoomIdx: 0, AdjIdx: 0, GovIdx: 1, OppIdx: 2, WinIdx: 1, Ballots: []struct {
					SpeakerIdx int
					Score      int
					Position   string
					TeamRole   string
				}{
					{SpeakerIdx: 2, Score: 76, Position: "PM", TeamRole: "gov"},
					{SpeakerIdx: 3, Score: 75, Position: "DPM", TeamRole: "gov"},
					{SpeakerIdx: 4, Score: 74, Position: "LO", TeamRole: "opp"},
					{SpeakerIdx: 5, Score: 72, Position: "DLO", TeamRole: "opp"},
				}},
				{RoomIdx: 1, AdjIdx: 1, GovIdx: 0, OppIdx: 4, WinIdx: 0, Ballots: []struct {
					SpeakerIdx int
					Score      int
					Position   string
					TeamRole   string
				}{
					{SpeakerIdx: 0, Score: 77, Position: "PM", TeamRole: "gov"},
					{SpeakerIdx: 1, Score: 75, Position: "DPM", TeamRole: "gov"},
					{SpeakerIdx: 8, Score: 73, Position: "LO", TeamRole: "opp"},
					{SpeakerIdx: 9, Score: 71, Position: "DLO", TeamRole: "opp"},
				}},
				{RoomIdx: 2, AdjIdx: 2, GovIdx: 3, OppIdx: 5, WinIdx: 3, Ballots: []struct {
					SpeakerIdx int
					Score      int
					Position   string
					TeamRole   string
				}{
					{SpeakerIdx: 6, Score: 74, Position: "PM", TeamRole: "gov"},
					{SpeakerIdx: 7, Score: 72, Position: "DPM", TeamRole: "gov"},
					{SpeakerIdx: 10, Score: 71, Position: "LO", TeamRole: "opp"},
					{SpeakerIdx: 11, Score: 69, Position: "DLO", TeamRole: "opp"},
				}},
				{RoomIdx: 3, AdjIdx: 3, GovIdx: 7, OppIdx: 6, WinIdx: 6, Ballots: []struct {
					SpeakerIdx int
					Score      int
					Position   string
					TeamRole   string
				}{
					{SpeakerIdx: 14, Score: 72, Position: "PM", TeamRole: "gov"},
					{SpeakerIdx: 15, Score: 68, Position: "DPM", TeamRole: "gov"},
					{SpeakerIdx: 12, Score: 73, Position: "LO", TeamRole: "opp"},
					{SpeakerIdx: 13, Score: 71, Position: "DLO", TeamRole: "opp"},
				}},
			},
		},
		{
			Name:      "Round 4",
			Motion:    "THW abolish all border controls between nations",
			InfoSlide: "Context: Global migration and refugee crises",
			Matches: []struct {
				RoomIdx int
				AdjIdx  int
				GovIdx  int
				OppIdx  int
				WinIdx  int
				Ballots []struct {
					SpeakerIdx int
					Score      int
					Position   string
					TeamRole   string
				}
			}{
				{RoomIdx: 0, AdjIdx: 4, GovIdx: 0, OppIdx: 2, WinIdx: 0, Ballots: []struct {
					SpeakerIdx int
					Score      int
					Position   string
					TeamRole   string
				}{
					{SpeakerIdx: 0, Score: 78, Position: "PM", TeamRole: "gov"},
					{SpeakerIdx: 1, Score: 77, Position: "DPM", TeamRole: "gov"},
					{SpeakerIdx: 4, Score: 76, Position: "LO", TeamRole: "opp"},
					{SpeakerIdx: 5, Score: 74, Position: "DLO", TeamRole: "opp"},
				}},
				{RoomIdx: 1, AdjIdx: 0, GovIdx: 4, OppIdx: 1, WinIdx: 1, Ballots: []struct {
					SpeakerIdx int
					Score      int
					Position   string
					TeamRole   string
				}{
					{SpeakerIdx: 8, Score: 73, Position: "PM", TeamRole: "gov"},
					{SpeakerIdx: 9, Score: 71, Position: "DPM", TeamRole: "gov"},
					{SpeakerIdx: 2, Score: 76, Position: "LO", TeamRole: "opp"},
					{SpeakerIdx: 3, Score: 75, Position: "DLO", TeamRole: "opp"},
				}},
				{RoomIdx: 2, AdjIdx: 1, GovIdx: 3, OppIdx: 7, WinIdx: 3, Ballots: []struct {
					SpeakerIdx int
					Score      int
					Position   string
					TeamRole   string
				}{
					{SpeakerIdx: 6, Score: 77, Position: "PM", TeamRole: "gov"},
					{SpeakerIdx: 7, Score: 75, Position: "DPM", TeamRole: "gov"},
					{SpeakerIdx: 14, Score: 74, Position: "LO", TeamRole: "opp"},
					{SpeakerIdx: 15, Score: 72, Position: "DLO", TeamRole: "opp"},
				}},
				{RoomIdx: 3, AdjIdx: 2, GovIdx: 5, OppIdx: 6, WinIdx: 5, Ballots: []struct {
					SpeakerIdx int
					Score      int
					Position   string
					TeamRole   string
				}{
					{SpeakerIdx: 10, Score: 75, Position: "PM", TeamRole: "gov"},
					{SpeakerIdx: 11, Score: 73, Position: "DPM", TeamRole: "gov"},
					{SpeakerIdx: 12, Score: 72, Position: "LO", TeamRole: "opp"},
					{SpeakerIdx: 13, Score: 70, Position: "DLO", TeamRole: "opp"},
				}},
			},
		},
		{
			Name:      "Grand Final",
			Motion:    "THW prioritize economic growth over environmental protection in developing nations",
			InfoSlide: "Context: Developing nations face trade-offs between development and sustainability",
			Matches: []struct {
				RoomIdx int
				AdjIdx  int
				GovIdx  int
				OppIdx  int
				WinIdx  int
				Ballots []struct {
					SpeakerIdx int
					Score      int
					Position   string
					TeamRole   string
				}
			}{
				{RoomIdx: 4, AdjIdx: 0, GovIdx: 0, OppIdx: 1, WinIdx: 0, Ballots: []struct {
					SpeakerIdx int
					Score      int
					Position   string
					TeamRole   string
				}{
					{SpeakerIdx: 0, Score: 79, Position: "PM", TeamRole: "gov"},
					{SpeakerIdx: 1, Score: 85, Position: "DPM", TeamRole: "gov"},
					{SpeakerIdx: 2, Score: 79, Position: "LO", TeamRole: "opp"},
					{SpeakerIdx: 3, Score: 85, Position: "DLO", TeamRole: "opp"},
				}},
			},
		},
	}

	for roundIdx, roundData := range rounds {
		round := models.Round{
			TournamentID:      tournament.ID,
			Name:              roundData.Name,
			Motion:            roundData.Motion,
			InfoSlide:         roundData.InfoSlide,
			IsPublished:       true,
			IsDrawPublished:   true,
			IsMotionPublished: true,
		}

		if err := models.DB.Create(&round).Error; err != nil {
			log.Fatal("Failed to create round:", err)
		}

		for _, matchData := range roundData.Matches {
			match := models.Match{
				RoundID:       round.ID,
				RoomID:        &rooms[matchData.RoomIdx].ID,
				AdjudicatorID: &adjudicators[matchData.AdjIdx].ID,
				GovTeamID:     &teams[matchData.GovIdx].ID,
				OppTeamID:     &teams[matchData.OppIdx].ID,
				WinnerID:      &teams[matchData.WinIdx].ID,
				IsCompleted:   true,
			}

			if err := models.DB.Create(&match).Error; err != nil {
				log.Fatal("Failed to create match:", err)
			}

			for _, ballotData := range matchData.Ballots {
				ballot := models.Ballot{
					MatchID:   match.ID,
					SpeakerID: speakers[ballotData.SpeakerIdx].ID,
					Score:     ballotData.Score,
					Position:  ballotData.Position,
					TeamRole:  ballotData.TeamRole,
					IsReply:   false,
				}

				if err := models.DB.Create(&ballot).Error; err != nil {
					log.Fatal("Failed to create ballot:", err)
				}
			}
		}

		fmt.Printf("✅ Created Round %d: %s with %d matches\n", roundIdx+1, roundData.Name, len(roundData.Matches))
	}

	fmt.Println("\n🎉 SUCCESS! PIMNAS 37 tournament created!")
	fmt.Println("📊 Summary:")
	fmt.Println("   - Tournament: PIMNAS 37 (Completed)")
	fmt.Println("   - Winner: UPI A (5 wins, 0 losses)")
	fmt.Println("   - Total Teams: 8")
	fmt.Println("   - Total Rounds: 5 (including Grand Final)")
	fmt.Println("   - All matches completed with ballots")
	fmt.Printf("   - Tournament ID: %d\n", tournament.ID)
}
