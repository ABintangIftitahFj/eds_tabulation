-- PIMNAS Test Tournament Data
-- This creates a complete tournament with finished rounds and a clear winner

-- 1. Create Tournament
INSERT INTO tournaments (name, slug, format, start_date, end_date, location, description, status, is_public, created_at, updated_at)
VALUES ('PIMNAS 37', 'pimnas-37', 'asian', '2025-01-15', '2025-01-18', 'Universitas Pendidikan Indonesia, Bandung', 'Pekan Ilmiah Mahasiswa Nasional ke-37', 'completed', true, NOW(), NOW());

SET @tournament_id = LAST_INSERT_ID();

-- 2. Create Teams (8 teams)
INSERT INTO teams (tournament_id, name, institution, total_vp, total_speaker, rank, wins, losses, created_at, updated_at)
VALUES 
(@tournament_id, 'UPI A', 'Universitas Pendidikan Indonesia', 5, 575, 1, 5, 0, NOW(), NOW()),
(@tournament_id, 'ITB A', 'Institut Teknologi Bandung', 4, 562, 2, 4, 1, NOW(), NOW()),
(@tournament_id, 'UI A', 'Universitas Indonesia', 4, 558, 3, 4, 1, NOW(), NOW()),
(@tournament_id, 'UGM A', 'Universitas Gadjah Mada', 3, 545, 4, 3, 2, NOW(), NOW()),
(@tournament_id, 'UNPAD A', 'Universitas Padjadjaran', 3, 540, 5, 3, 2, NOW(), NOW()),
(@tournament_id, 'ITS A', 'Institut Teknologi Sepuluh Nopember', 2, 528, 6, 2, 3, NOW(), NOW()),
(@tournament_id, 'UNAIR A', 'Universitas Airlangga', 1, 515, 7, 1, 4, NOW(), NOW()),
(@tournament_id, 'UNDIP A', 'Universitas Diponegoro', 0, 502, 8, 0, 5, NOW(), NOW());

-- Get Team IDs
SET @team_upi = (SELECT id FROM teams WHERE name = 'UPI A' AND tournament_id = @tournament_id);
SET @team_itb = (SELECT id FROM teams WHERE name = 'ITB A' AND tournament_id = @tournament_id);
SET @team_ui = (SELECT id FROM teams WHERE name = 'UI A' AND tournament_id = @tournament_id);
SET @team_ugm = (SELECT id FROM teams WHERE name = 'UGM A' AND tournament_id = @tournament_id);
SET @team_unpad = (SELECT id FROM teams WHERE name = 'UNPAD A' AND tournament_id = @tournament_id);
SET @team_its = (SELECT id FROM teams WHERE name = 'ITS A' AND tournament_id = @tournament_id);
SET @team_unair = (SELECT id FROM teams WHERE name = 'UNAIR A' AND tournament_id = @tournament_id);
SET @team_undip = (SELECT id FROM teams WHERE name = 'UNDIP A' AND tournament_id = @tournament_id);

-- 3. Create Speakers for each team (2 speakers per team)
INSERT INTO speakers (team_id, name, total_score, speaker_rank, created_at, updated_at)
VALUES 
-- UPI A
(@team_upi, 'Ahmad Rifai', 288, 1, NOW(), NOW()),
(@team_upi, 'Siti Nurhaliza', 287, 2, NOW(), NOW()),
-- ITB A
(@team_itb, 'Budi Santoso', 282, 3, NOW(), NOW()),
(@team_itb, 'Dewi Lestari', 280, 5, NOW(), NOW()),
-- UI A
(@team_ui, 'Cahya Prasetyo', 281, 4, NOW(), NOW()),
(@team_ui, 'Eka Putri', 277, 7, NOW(), NOW()),
-- UGM A
(@team_ugm, 'Fajar Nugroho', 279, 6, NOW(), NOW()),
(@team_ugm, 'Gita Savitri', 266, 12, NOW(), NOW()),
-- UNPAD A
(@team_unpad, 'Hendra Wijaya', 275, 8, NOW(), NOW()),
(@team_unpad, 'Indah Permata', 265, 13, NOW(), NOW()),
-- ITS A
(@team_its, 'Joko Widodo', 270, 10, NOW(), NOW()),
(@team_its, 'Kartika Sari', 258, 15, NOW(), NOW()),
-- UNAIR A
(@team_unair, 'Luthfi Hakim', 268, 11, NOW(), NOW()),
(@team_unair, 'Maya Anggraini', 247, 16, NOW(), NOW()),
-- UNDIP A
(@team_undip, 'Nanda Prakoso', 273, 9, NOW(), NOW()),
(@team_undip, 'Olivia Situmorang', 229, 23, NOW(), NOW());

-- 4. Create Adjudicators
INSERT INTO adjudicators (tournament_id, name, institution, level, is_available, created_at, updated_at)
VALUES 
(@tournament_id, 'Dr. Andi Setiawan', 'Universitas Indonesia', 'Chief', true, NOW(), NOW()),
(@tournament_id, 'Prof. Benny Kurniawan', 'Institut Teknologi Bandung', 'Chief', true, NOW(), NOW()),
(@tournament_id, 'Citra Maharani, M.A.', 'Universitas Gadjah Mada', 'Wing', true, NOW(), NOW()),
(@tournament_id, 'Doni Prasetya, S.S.', 'Universitas Pendidikan Indonesia', 'Wing', true, NOW(), NOW()),
(@tournament_id, 'Eka Wulandari', 'Universitas Padjadjaran', 'Panelist', true, NOW(), NOW());

SET @adj1 = (SELECT id FROM adjudicators WHERE name = 'Dr. Andi Setiawan' AND tournament_id = @tournament_id);
SET @adj2 = (SELECT id FROM adjudicators WHERE name = 'Prof. Benny Kurniawan' AND tournament_id = @tournament_id);
SET @adj3 = (SELECT id FROM adjudicators WHERE name = 'Citra Maharani, M.A.' AND tournament_id = @tournament_id);
SET @adj4 = (SELECT id FROM adjudicators WHERE name = 'Doni Prasetya, S.S.' AND tournament_id = @tournament_id);
SET @adj5 = (SELECT id FROM adjudicators WHERE name = 'Eka Wulandari' AND tournament_id = @tournament_id);

-- 5. Create Rooms
INSERT INTO rooms (tournament_id, name, location, capacity, is_available, created_at, updated_at)
VALUES 
(@tournament_id, 'R1', 'Gedung A Lantai 1', 50, true, NOW(), NOW()),
(@tournament_id, 'R2', 'Gedung A Lantai 2', 50, true, NOW(), NOW()),
(@tournament_id, 'R3', 'Gedung B Lantai 1', 50, true, NOW(), NOW()),
(@tournament_id, 'R4', 'Gedung B Lantai 2', 50, true, NOW(), NOW()),
(@tournament_id, 'Final Room', 'Auditorium Utama', 200, true, NOW(), NOW());

SET @room1 = (SELECT id FROM rooms WHERE name = 'R1' AND tournament_id = @tournament_id);
SET @room2 = (SELECT id FROM rooms WHERE name = 'R2' AND tournament_id = @tournament_id);
SET @room3 = (SELECT id FROM rooms WHERE name = 'R3' AND tournament_id = @tournament_id);
SET @room4 = (SELECT id FROM rooms WHERE name = 'R4' AND tournament_id = @tournament_id);
SET @room_final = (SELECT id FROM rooms WHERE name = 'Final Room' AND tournament_id = @tournament_id);

-- 6. Create Rounds with Matches
-- Round 1
INSERT INTO rounds (tournament_id, name, motion, info_slide, is_published, is_draw_published, is_motion_published, created_at, updated_at)
VALUES (@tournament_id, 'Round 1', 'THW ban social media platforms from using algorithmic content recommendation', 'Context: Social media algorithms prioritize engagement over wellbeing', true, true, true, NOW(), NOW());

SET @round1 = LAST_INSERT_ID();

-- Round 1 Matches
INSERT INTO matches (round_id, room_id, adjudicator_id, gov_team_id, opp_team_id, winner_id, is_completed, created_at, updated_at)
VALUES 
(@round1, @room1, @adj1, @team_upi, @team_undip, @team_upi, true, NOW(), NOW()),
(@round1, @room2, @adj2, @team_itb, @team_unair, @team_itb, true, NOW(), NOW()),
(@round1, @room3, @adj3, @team_ui, @team_its, @team_ui, true, NOW(), NOW()),
(@round1, @room4, @adj4, @team_ugm, @team_unpad, @team_ugm, true, NOW(), NOW());

-- Get speaker IDs for Round 1
SET @upi_s1 = (SELECT id FROM speakers WHERE name = 'Ahmad Rifai');
SET @upi_s2 = (SELECT id FROM speakers WHERE name = 'Siti Nurhaliza');
SET @undip_s1 = (SELECT id FROM speakers WHERE name = 'Nanda Prakoso');
SET @undip_s2 = (SELECT id FROM speakers WHERE name = 'Olivia Situmorang');
SET @itb_s1 = (SELECT id FROM speakers WHERE name = 'Budi Santoso');
SET @itb_s2 = (SELECT id FROM speakers WHERE name = 'Dewi Lestari');
SET @unair_s1 = (SELECT id FROM speakers WHERE name = 'Luthfi Hakim');
SET @unair_s2 = (SELECT id FROM speakers WHERE name = 'Maya Anggraini');
SET @ui_s1 = (SELECT id FROM speakers WHERE name = 'Cahya Prasetyo');
SET @ui_s2 = (SELECT id FROM speakers WHERE name = 'Eka Putri');
SET @its_s1 = (SELECT id FROM speakers WHERE name = 'Joko Widodo');
SET @its_s2 = (SELECT id FROM speakers WHERE name = 'Kartika Sari');
SET @ugm_s1 = (SELECT id FROM speakers WHERE name = 'Fajar Nugroho');
SET @ugm_s2 = (SELECT id FROM speakers WHERE name = 'Gita Savitri');
SET @unpad_s1 = (SELECT id FROM speakers WHERE name = 'Hendra Wijaya');
SET @unpad_s2 = (SELECT id FROM speakers WHERE name = 'Indah Permata');

-- Round 1 Match 1 Ballots (UPI vs UNDIP)
SET @match1 = (SELECT id FROM matches WHERE round_id = @round1 AND gov_team_id = @team_upi);
INSERT INTO ballots (match_id, speaker_id, score, position, is_reply, team_role, created_at, updated_at)
VALUES 
(@match1, @upi_s1, 76, 'PM', false, 'gov', NOW(), NOW()),
(@match1, @upi_s2, 74, 'DPM', false, 'gov', NOW(), NOW()),
(@match1, @undip_s1, 72, 'LO', false, 'opp', NOW(), NOW()),
(@match1, @undip_s2, 68, 'DLO', false, 'opp', NOW(), NOW());

-- Round 1 Match 2 Ballots (ITB vs UNAIR)
SET @match2 = (SELECT id FROM matches WHERE round_id = @round1 AND gov_team_id = @team_itb);
INSERT INTO ballots (match_id, speaker_id, score, position, is_reply, team_role, created_at, updated_at)
VALUES 
(@match2, @itb_s1, 75, 'PM', false, 'gov', NOW(), NOW()),
(@match2, @itb_s2, 73, 'DPM', false, 'gov', NOW(), NOW()),
(@match2, @unair_s1, 71, 'LO', false, 'opp', NOW(), NOW()),
(@match2, @unair_s2, 69, 'DLO', false, 'opp', NOW(), NOW());

-- Round 1 Match 3 Ballots (UI vs ITS)
SET @match3 = (SELECT id FROM matches WHERE round_id = @round1 AND gov_team_id = @team_ui);
INSERT INTO ballots (match_id, speaker_id, score, position, is_reply, team_role, created_at, updated_at)
VALUES 
(@match3, @ui_s1, 74, 'PM', false, 'gov', NOW(), NOW()),
(@match3, @ui_s2, 72, 'DPM', false, 'gov', NOW(), NOW()),
(@match3, @its_s1, 70, 'LO', false, 'opp', NOW(), NOW()),
(@match3, @its_s2, 68, 'DLO', false, 'opp', NOW(), NOW());

-- Round 1 Match 4 Ballots (UGM vs UNPAD)
SET @match4 = (SELECT id FROM matches WHERE round_id = @round1 AND gov_team_id = @team_ugm);
INSERT INTO ballots (match_id, speaker_id, score, position, is_reply, team_role, created_at, updated_at)
VALUES 
(@match4, @ugm_s1, 73, 'PM', false, 'gov', NOW(), NOW()),
(@match4, @ugm_s2, 71, 'DPM', false, 'gov', NOW(), NOW()),
(@match4, @unpad_s1, 72, 'LO', false, 'opp', NOW(), NOW()),
(@match4, @unpad_s2, 70, 'DLO', false, 'opp', NOW(), NOW());

-- Round 2
INSERT INTO rounds (tournament_id, name, motion, info_slide, is_published, is_draw_published, is_motion_published, created_at, updated_at)
VALUES (@tournament_id, 'Round 2', 'THW require all schools to teach financial literacy from elementary level', 'Context: Many adults struggle with basic financial management', true, true, true, NOW(), NOW());

SET @round2 = LAST_INSERT_ID();

-- Round 2 Matches
INSERT INTO matches (round_id, room_id, adjudicator_id, gov_team_id, opp_team_id, winner_id, is_completed, created_at, updated_at)
VALUES 
(@round2, @room1, @adj2, @team_ui, @team_ugm, @team_ui, true, NOW(), NOW()),
(@round2, @room2, @adj3, @team_upi, @team_itb, @team_upi, true, NOW(), NOW()),
(@round2, @room3, @adj4, @team_unpad, @team_its, @team_unpad, true, NOW(), NOW()),
(@round2, @room4, @adj5, @team_unair, @team_undip, @team_undip, true, NOW(), NOW());

-- Round 2 Ballots
SET @match5 = (SELECT id FROM matches WHERE round_id = @round2 AND gov_team_id = @team_ui);
INSERT INTO ballots (match_id, speaker_id, score, position, is_reply, team_role, created_at, updated_at)
VALUES 
(@match5, @ui_s1, 77, 'PM', false, 'gov', NOW(), NOW()),
(@match5, @ui_s2, 75, 'DPM', false, 'gov', NOW(), NOW()),
(@match5, @ugm_s1, 74, 'LO', false, 'opp', NOW(), NOW()),
(@match5, @ugm_s2, 70, 'DLO', false, 'opp', NOW(), NOW());

SET @match6 = (SELECT id FROM matches WHERE round_id = @round2 AND gov_team_id = @team_upi);
INSERT INTO ballots (match_id, speaker_id, score, position, is_reply, team_role, created_at, updated_at)
VALUES 
(@match6, @upi_s1, 78, 'PM', false, 'gov', NOW(), NOW()),
(@match6, @upi_s2, 76, 'DPM', false, 'gov', NOW(), NOW()),
(@match6, @itb_s1, 75, 'LO', false, 'opp', NOW(), NOW()),
(@match6, @itb_s2, 73, 'DLO', false, 'opp', NOW(), NOW());

SET @match7 = (SELECT id FROM matches WHERE round_id = @round2 AND gov_team_id = @team_unpad);
INSERT INTO ballots (match_id, speaker_id, score, position, is_reply, team_role, created_at, updated_at)
VALUES 
(@match7, @unpad_s1, 76, 'PM', false, 'gov', NOW(), NOW()),
(@match7, @unpad_s2, 74, 'DPM', false, 'gov', NOW(), NOW()),
(@match7, @its_s1, 73, 'LO', false, 'opp', NOW(), NOW()),
(@match7, @its_s2, 71, 'DLO', false, 'opp', NOW(), NOW());

SET @match8 = (SELECT id FROM matches WHERE round_id = @round2 AND gov_team_id = @team_unair);
INSERT INTO ballots (match_id, speaker_id, score, position, is_reply, team_role, created_at, updated_at)
VALUES 
(@match8, @unair_s1, 70, 'PM', false, 'gov', NOW(), NOW()),
(@match8, @unair_s2, 68, 'DPM', false, 'gov', NOW(), NOW()),
(@match8, @undip_s1, 72, 'LO', false, 'opp', NOW(), NOW()),
(@match8, @undip_s2, 70, 'DLO', false, 'opp', NOW(), NOW());

-- Round 3
INSERT INTO rounds (tournament_id, name, motion, info_slide, is_published, is_draw_published, is_motion_published, created_at, updated_at)
VALUES (@tournament_id, 'Round 3', 'TH regrets the rise of gig economy', 'Context: Freelance and contract work has become increasingly common', true, true, true, NOW(), NOW());

SET @round3 = LAST_INSERT_ID();

INSERT INTO matches (round_id, room_id, adjudicator_id, gov_team_id, opp_team_id, winner_id, is_completed, created_at, updated_at)
VALUES 
(@round3, @room1, @adj1, @team_itb, @team_ui, @team_itb, true, NOW(), NOW()),
(@round3, @room2, @adj2, @team_upi, @team_unpad, @team_upi, true, NOW(), NOW()),
(@round3, @room3, @adj3, @team_ugm, @team_its, @team_ugm, true, NOW(), NOW()),
(@round3, @room4, @adj4, @team_undip, @team_unair, @team_unair, true, NOW(), NOW());

-- Round 3 Ballots
SET @match9 = (SELECT id FROM matches WHERE round_id = @round3 AND gov_team_id = @team_itb);
INSERT INTO ballots (match_id, speaker_id, score, position, is_reply, team_role, created_at, updated_at)
VALUES 
(@match9, @itb_s1, 76, 'PM', false, 'gov', NOW(), NOW()),
(@match9, @itb_s2, 75, 'DPM', false, 'gov', NOW(), NOW()),
(@match9, @ui_s1, 74, 'LO', false, 'opp', NOW(), NOW()),
(@match9, @ui_s2, 72, 'DLO', false, 'opp', NOW(), NOW());

SET @match10 = (SELECT id FROM matches WHERE round_id = @round3 AND gov_team_id = @team_upi);
INSERT INTO ballots (match_id, speaker_id, score, position, is_reply, team_role, created_at, updated_at)
VALUES 
(@match10, @upi_s1, 77, 'PM', false, 'gov', NOW(), NOW()),
(@match10, @upi_s2, 75, 'DPM', false, 'gov', NOW(), NOW()),
(@match10, @unpad_s1, 73, 'LO', false, 'opp', NOW(), NOW()),
(@match10, @unpad_s2, 71, 'DLO', false, 'opp', NOW(), NOW());

SET @match11 = (SELECT id FROM matches WHERE round_id = @round3 AND gov_team_id = @team_ugm);
INSERT INTO ballots (match_id, speaker_id, score, position, is_reply, team_role, created_at, updated_at)
VALUES 
(@match11, @ugm_s1, 74, 'PM', false, 'gov', NOW(), NOW()),
(@match11, @ugm_s2, 72, 'DPM', false, 'gov', NOW(), NOW()),
(@match11, @its_s1, 71, 'LO', false, 'opp', NOW(), NOW()),
(@match11, @its_s2, 69, 'DLO', false, 'opp', NOW(), NOW());

SET @match12 = (SELECT id FROM matches WHERE round_id = @round3 AND gov_team_id = @team_undip);
INSERT INTO ballots (match_id, speaker_id, score, position, is_reply, team_role, created_at, updated_at)
VALUES 
(@match12, @undip_s1, 72, 'PM', false, 'gov', NOW(), NOW()),
(@match12, @undip_s2, 68, 'DPM', false, 'gov', NOW(), NOW()),
(@match12, @unair_s1, 73, 'LO', false, 'opp', NOW(), NOW()),
(@match12, @unair_s2, 71, 'DLO', false, 'opp', NOW(), NOW());

-- Round 4
INSERT INTO rounds (tournament_id, name, motion, info_slide, is_published, is_draw_published, is_motion_published, created_at, updated_at)
VALUES (@tournament_id, 'Round 4', 'THW abolish all border controls between nations', 'Context: Global migration and refugee crises', true, true, true, NOW(), NOW());

SET @round4 = LAST_INSERT_ID();

INSERT INTO matches (round_id, room_id, adjudicator_id, gov_team_id, opp_team_id, winner_id, is_completed, created_at, updated_at)
VALUES 
(@round4, @room1, @adj5, @team_upi, @team_ui, @team_upi, true, NOW(), NOW()),
(@round4, @room2, @adj1, @team_unpad, @team_itb, @team_itb, true, NOW(), NOW()),
(@round4, @room3, @adj2, @team_ugm, @team_undip, @team_ugm, true, NOW(), NOW()),
(@round4, @room4, @adj3, @team_its, @team_unair, @team_its, true, NOW(), NOW());

-- Round 4 Ballots
SET @match13 = (SELECT id FROM matches WHERE round_id = @round4 AND gov_team_id = @team_upi);
INSERT INTO ballots (match_id, speaker_id, score, position, is_reply, team_role, created_at, updated_at)
VALUES 
(@match13, @upi_s1, 78, 'PM', false, 'gov', NOW(), NOW()),
(@match13, @upi_s2, 77, 'DPM', false, 'gov', NOW(), NOW()),
(@match13, @ui_s1, 76, 'LO', false, 'opp', NOW(), NOW()),
(@match13, @ui_s2, 74, 'DLO', false, 'opp', NOW(), NOW());

SET @match14 = (SELECT id FROM matches WHERE round_id = @round4 AND gov_team_id = @team_unpad);
INSERT INTO ballots (match_id, speaker_id, score, position, is_reply, team_role, created_at, updated_at)
VALUES 
(@match14, @unpad_s1, 73, 'PM', false, 'gov', NOW(), NOW()),
(@match14, @unpad_s2, 71, 'DPM', false, 'gov', NOW(), NOW()),
(@match14, @itb_s1, 76, 'LO', false, 'opp', NOW(), NOW()),
(@match14, @itb_s2, 75, 'DLO', false, 'opp', NOW(), NOW());

SET @match15 = (SELECT id FROM matches WHERE round_id = @round4 AND gov_team_id = @team_ugm);
INSERT INTO ballots (match_id, speaker_id, score, position, is_reply, team_role, created_at, updated_at)
VALUES 
(@match15, @ugm_s1, 77, 'PM', false, 'gov', NOW(), NOW()),
(@match15, @ugm_s2, 75, 'DPM', false, 'gov', NOW(), NOW()),
(@match15, @undip_s1, 74, 'LO', false, 'opp', NOW(), NOW()),
(@match15, @undip_s2, 72, 'DLO', false, 'opp', NOW(), NOW());

SET @match16 = (SELECT id FROM matches WHERE round_id = @round4 AND gov_team_id = @team_its);
INSERT INTO ballots (match_id, speaker_id, score, position, is_reply, team_role, created_at, updated_at)
VALUES 
(@match16, @its_s1, 75, 'PM', false, 'gov', NOW(), NOW()),
(@match16, @its_s2, 73, 'DPM', false, 'gov', NOW(), NOW()),
(@match16, @unair_s1, 72, 'LO', false, 'opp', NOW(), NOW()),
(@match16, @unair_s2, 70, 'DLO', false, 'opp', NOW(), NOW());

-- Grand Final
INSERT INTO rounds (tournament_id, name, motion, info_slide, is_published, is_draw_published, is_motion_published, created_at, updated_at)
VALUES (@tournament_id, 'Grand Final', 'THW prioritize economic growth over environmental protection in developing nations', 'Context: Developing nations face trade-offs between development and sustainability', true, true, true, NOW(), NOW());

SET @round5 = LAST_INSERT_ID();

INSERT INTO matches (round_id, room_id, adjudicator_id, gov_team_id, opp_team_id, winner_id, is_completed, created_at, updated_at)
VALUES 
(@round5, @room_final, @adj1, @team_upi, @team_itb, @team_upi, true, NOW(), NOW());

-- Grand Final Ballots (UPI vs ITB)
SET @final_match = (SELECT id FROM matches WHERE round_id = @round5);
INSERT INTO ballots (match_id, speaker_id, score, position, is_reply, team_role, created_at, updated_at)
VALUES 
(@final_match, @upi_s1, 79, 'PM', false, 'gov', NOW(), NOW()),
(@final_match, @upi_s2, 85, 'DPM', false, 'gov', NOW(), NOW()),
(@final_match, @itb_s1, 80, 'LO', false, 'opp', NOW(), NOW()),
(@final_match, @itb_s2, 84, 'DLO', false, 'opp', NOW(), NOW());

-- Success message
SELECT 'PIMNAS 37 test tournament created successfully!' as message,
       'Winner: UPI A with 5 wins' as result,
       'Total Teams: 8' as teams,
       'Total Rounds: 5 (including Grand Final)' as rounds;
