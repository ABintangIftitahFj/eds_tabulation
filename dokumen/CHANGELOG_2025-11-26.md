# Changelog - Ballot Score System Update

## Date: 2025-11-26

### Changes Implemented

#### 1. ✅ Ballot Scores Changed to Integers (No Decimals)

**Backend Changes:**
- **File**: `Backend/models/schema.go`
  - Changed `Ballot.Score` from `float64` to `int`
  - Changed `Team.TotalSpeaker` from `float64` to `int`
  - Changed `Speaker.TotalScore` from `float64` to `int`

- **File**: `Backend/controllers/ballot_controller.go`
  - Updated `totalGov` and `totalOpp` variables from `float64` to `int`

**Frontend Changes:**
- **File**: `frontend/app/tournaments/[id]/standings/page.tsx`
  - Removed `.toFixed(1)` formatting for `team.total_speaker`
  - Removed `.toFixed(1)` formatting for `speaker.total_score`
  - Removed `.toFixed(1)` formatting for `inst.total_speaks`
  - Scores now display as integers: `75` instead of `75.0`

**Impact:**
- Speaker scores are now whole numbers (e.g., 72, 75, 80)
- No more decimal points in ballot scoring
- Cleaner, more professional appearance
- Matches traditional debate scoring practices

#### 2. ✅ Created PIMNAS 37 Test Tournament

**What was created:**
- Complete tournament with finished matches
- **Tournament Name**: PIMNAS 37 (Pekan Ilmiah Mahasiswa Nasional ke-37)
- **Format**: Asian Parliamentary
- **Status**: Completed

**Tournament Data:**
- **8 Teams** from top Indonesian universities
- **16 Speakers** (2 per team)
- **5 Adjudicators** with different levels (Chief, Wing, Panelist)
- **5 Rooms** for matches
- **5 Rounds** including Grand Final
- **20 Completed Matches** with full ballot data
- **80 Ballot Entries** with integer scores

**Champion**: 🏆 **UPI A** (Universitas Pendidikan Indonesia)
- Record: 5 Wins - 0 Losses
- Total VP: 5
- Total Speaker Points: 575

**Files Created:**
1. `Backend/scripts/create_pimnas_data.go` - Script to populate tournament data
2. `Backend/pimnas_test_data.sql` - SQL script (alternative method)
3. `PIMNAS_TOURNAMENT_INFO.md` - Complete tournament documentation

**How to Run the Script:**
```bash
cd Backend
go run scripts/create_pimnas_data.go
```

### Testing the Changes

1. **View PIMNAS Tournament:**
   - Navigate to `/tournaments` in your frontend
   - Find "PIMNAS 37" in the list
   - Click to view tournament details

2. **Check Standings:**
   - Go to `/tournaments/[id]/standings`
   - Verify scores display as integers (no decimals)
   - Check team standings show UPI A as champion

3. **Verify Ballot Scores:**
   - All speaker scores should be whole numbers
   - No `.0` decimal points
   - Cleaner number display throughout

### Database Migration Notes

⚠️ **Important**: These schema changes require database migration. The changes will be automatically applied when you restart your Go server (AutoMigrate is enabled in `models/setup.go`).

If you need to manually reset the database:
1. Drop existing tables with float scores
2. Restart the Go server to recreate tables with int scores
3. Run the PIMNAS data script to populate test data

### Related Files Modified

**Backend:**
- ✅ `Backend/models/schema.go`
- ✅ `Backend/controllers/ballot_controller.go`

**Frontend:**
- ✅ `frontend/app/tournaments/[id]/standings/page.tsx`

**New Files:**
- ✅ `Backend/scripts/create_pimnas_data.go`
- ✅ `Backend/pimnas_test_data.sql`
- ✅ `PIMNAS_TOURNAMENT_INFO.md`
- ✅ `CHANGELOG_2025-11-26.md` (this file)

### Summary

✨ **All requested changes completed:**
1. Ballot scores now use integers instead of decimals
2. Complete PIMNAS 37 tournament created with a clear winner
3. All matches completed with realistic debate data
4. Frontend updated to display integer scores properly

The system is now ready for testing with realistic tournament data! 🎉
