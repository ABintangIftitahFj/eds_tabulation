# Ballot System Fix Requirements

## Issues to Fix:

### 1. Adjudicator Panel Creation Error
**Problem**: When selecting "Solo Adjudicator", the system still requires selecting a "Chief Adjudicator" which doesn't show up in the list.

**Solution**:
- File: `frontend/app/admin/ballot/page.tsx`
- Line 122-152: Update `createAdjudicatorPanel()` function
- Change validation logic:
  ```typescript
  // For solo adjudicator (panel_size === 1), only need one adjudicator
  if (currentPanel.panel_size === 1 && !currentPanel.chief_adj_id) {
      return toast.warning("Select an adjudicator!");
  }
  
  // For panels > 1, require chief
  if (currentPanel.panel_size > 1 && !currentPanel.chief_adj_id) {
      return toast.warning("Select a chief adjudicator!");
  }
  ```

- Line 510: REMOVE the filter `.filter(a => a.rating >= 4)` 
  Change from:
  ```typescript
  {adjudicators.filter(a => a.rating >= 4).map(adj => (
  ```
  To:
  ```typescript
  {adjudicators.map(adj => (
  ```

### 2. Team Speaker Requirements Based on Tournament Format

**Problem**: PIMNAS teams have only 2 members, but Asian Parliamentary requires 3 speakers + 1 optional reply

**Solution**:
- Need to update team registration validation
- Asian Parliamentary: Require exactly 3 speakers (PM, DPM, GW)  + 1 optional reply speaker
- British Parliamentary: Require exactly 2 speakers per team

**Files to Update**:
1. `Backend/controllers/team_controller.go` - Add validation when creating teams
2. `frontend/app/admin/teams/*` - Add speaker count validation based on tournament format
3. PIMNAS data script needs to be updated to have 3 speakers per team

### 3. Add "Back to Tournaments" Button on User Pages

**Problem**: User pages don't have a button to go back to the tournament list

**Files to Update**:
- `frontend/app/tournaments/[id]/standings/page.tsx`
- `frontend/app/tournaments/[id]/motions/page.tsx`
- `frontend/app/tournaments/[id]/results/page.tsx`
- `frontend/app/tournaments/[id]/participants/page.tsx`

Add this to the navigation tabs section (around line 186):
```tsx
<Link
    href="/tournaments"
    className="px-4 py-2 text-sm rounded-lg font-semibold transition-all text-gray-700 hover:bg-gray-100"
>
    ← All Tournaments
</Link>
```

## Summary of Changes Needed:

### Immediate Fixes (High Priority):
1. ✅ Fix adjudicator panel creation for solo adjudicators
2. ✅ Remove rating filter that prevents adjudicators from showing
3. ✅ Add "Back to tournaments" button on all user tournament pages

### Data Fixes (Medium Priority):
4. Update PIMNAS test data to have 3 speakers per team instead of 2
5. Add tournament format validation for speaker counts

### Backend Validation (Medium Priority):
6. Add speaker count validation in team creation API
7. Validate speaker requirements based on tournament format

## Steps to Apply:

1. First, restore the ballot page file from git if it's corrupted
2. Apply the adjudicator panel fixes
3. Add navigation buttons to user pages  
4. Update PIMNAS data script
5. Add backend validation for speaker counts

## Test After Fixes:
- [ ] Can create solo adjudicator panel
- [ ] Adjudicators show in the dropdown list
- [ ] Can submit ballots successfully
- [ ] User can navigate back to tournament list
- [ ] Teams have correct number of speakers for their format
