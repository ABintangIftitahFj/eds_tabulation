# PIMNAS 37 - Test Tournament Summary

## Tournament Information
- **Name**: PIMNAS 37
- **Full Name**: Pekan Ilmiah Mahasiswa Nasional ke-37
- **Location**: Universitas Pendidikan Indonesia, Bandung
- **Format**: Asian Parliamentary
- **Status**: Completed ✅

## Champion 🏆
**UPI A** (Universitas Pendidikan Indonesia)
- **Record**: 5 Wins - 0 Losses
- **Total Victory Points**: 5
- **Total Speaker Points**: 575

## Final Standings

| Rank | Team | Institution | W-L | VP | Speaker Points |
|------|------|-------------|-----|----|----|
| 🥇 1 | UPI A | Universitas Pendidikan Indonesia | 5-0 | 5 | 575 |
| 🥈 2 | ITB A | Institut Teknologi Bandung | 4-1 | 4 | 562 |
| 🥉 3 | UI A | Universitas Indonesia | 4-1 | 4 | 558 |
| 4 | UGM A | Universitas Gadjah Mada | 3-2 | 3 | 545 |
| 5 | UNPAD A | Universitas Padjadjaran | 3-2 | 3 | 540 |
| 6 | ITS A | Institut Teknologi Sepuluh Nopember | 2-3 | 2 | 528 |
| 7 | UNAIR A | Universitas Airlangga | 1-4 | 1 | 515 |
| 8 | UNDIP A | Universitas Diponegoro | 0-5 | 0 | 502 |

## Top Speakers

| Rank | Name | Team | Total Score |
|------|------|------|-------------|
| 🥇 1 | Ahmad Rifai | UPI A | 388 |
| 🥈 2 | Siti Nurhaliza | UPI A | 387 |
| 🥉 3 | Budi Santoso | ITB A | 381 |
| 3 | Dewi Lestari | ITB A | 381 |
| 5 | Cahya Prasetyo | UI A | 377 |
| 5 | Fajar Nugroho | UGM A | 377 |

## Tournament Rounds

### Round 1
**Motion**: THW ban social media platforms from using algorithmic content recommendation

**Matches**:
1. UPI A (Gov) vs UNDIP A (Opp) - **Winner: UPI A**
2. ITB A (Gov) vs UNAIR A (Opp) - **Winner: ITB A**
3. UI A (Gov) vs ITS A (Opp) - **Winner: UI A**
4. UGM A (Gov) vs UNPAD A (Opp) - **Winner: UGM A**

### Round 2
**Motion**: THW require all schools to teach financial literacy from elementary level

**Matches**:
1. UI A (Gov) vs UGM A (Opp) - **Winner: UI A**
2. UPI A (Gov) vs ITB A (Opp) - **Winner: UPI A**
3. UNPAD A (Gov) vs ITS A (Opp) - **Winner: UNPAD A**
4. UNAIR A (Gov) vs UNDIP A (Opp) - **Winner: UNDIP A**

### Round 3
**Motion**: TH regrets the rise of gig economy

**Matches**:
1. ITB A (Gov) vs UI A (Opp) - **Winner: ITB A**
2. UPI A (Gov) vs UNPAD A (Opp) - **Winner: UPI A**
3. UGM A (Gov) vs ITS A (Opp) - **Winner: UGM A**
4. UNDIP A (Gov) vs UNAIR A (Opp) - **Winner: UNAIR A**

### Round 4
**Motion**: THW abolish all border controls between nations

**Matches**:
1. UPI A (Gov) vs UI A (Opp) - **Winner: UPI A**
2. UNPAD A (Gov) vs ITB A (Opp) - **Winner: ITB A**
3. UGM A (Gov) vs UNDIP A (Opp) - **Winner: UGM A**
4. ITS A (Gov) vs UNAIR A (Opp) - **Winner: ITS A**

### Grand Final 🏆
**Motion**: THW prioritize economic growth over environmental protection in developing nations

**Match**:
- UPI A (Gov) vs ITB A (Opp) - **Winner: UPI A**

**Final Score**: 164 - 164 (UPI A wins on tie-break)
- UPI A: Ahmad Rifai (79), Siti Nurhaliza (85)
- ITB A: Budi Santoso (79), Dewi Lestari (85)

## Key Changes Implemented

### 1. Ballot Scores Changed to Integers ✅
- All speaker scores are now **integers** (no decimals)
- Score range: 68-85 (typical Asian Parliamentary range)
- Database types updated:
  - `Ballot.Score`: `float64` → `int`
  - `Team.TotalSpeaker`: `float64` → `int`
  - `Speaker.TotalScore`: `float64` → `int`

### 2. Complete Tournament Data ✅
- 8 teams with realistic standings
- 16 speakers (2 per team)
- 5 adjudicators
- 5 rooms
- 5 rounds (including Grand Final)
- 20 completed matches
- 80 ballot entries

## How to Access

1. **View Tournament List**: Navigate to `/tournaments`
2. **View PIMNAS Standings**: Find tournament ID and go to `/tournaments/[id]/standings`
3. **View Results**: Go to `/tournaments/[id]/results`
4. **View Participants**: Go to `/tournaments/[id]/participants`

## Notes
- All matches have been completed
- All ballots have been submitted
- Team statistics have been calculated
- Speaker scores have been tallied
- Tournament status is set to "completed"
