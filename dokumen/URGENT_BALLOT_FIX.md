# URGENT: Ballot Page Corrupted - Manual Fix Required

## Problem
The file `frontend/app/admin/ballot/page.tsx` has been corrupted during automated editing attempts and cannot compile.

## Solution: Manual Restoration

### Step 1: Stop the Dev Server  
Press `Ctrl+C` in the terminal running `npm run dev`

### Step 2: Restore from Git (if tracked)
Run this command in PowerShell from the project root:
```
cd frontend
git checkout -- app/admin/ballot/page.tsx
```

OR if that doesn't work, use Git GUI:
1. Open your Git client (GitHub Desktop, VS Code Source Control, etc.)
2. Find `frontend/app/admin/ballot/page.tsx` in changed files
3. Right-click → Discard Changes

### Step 3: If File Wasn't in Git (New File)
You'll need to recreate it or use a backup. Let me know if this is the case.

### Step 4: After Restoration, Make These Manual Fixes

Once the file is restored, make ONLY these targeted changes:

#### Fix 1: Remove Adjudicator Rating Filter (Line ~510)
**Find this line:**
```tsx
{adjudicators.filter(a => a.rating >= 4).map(adj => (
```

**Replace with:**
```tsx
{adjudicators.map(adj => (
```

This fix will make ALL adjudicators show in the dropdown, not just high-rated ones.

#### Fix 2: Fix Solo Adjudicator Validation (Lines ~122-130)
**Find this code:**
```tsx
const createAdjudicatorPanel = async () => {
    if (!currentPanel.chief_adj_id) {
        return toast.warning("Select a chief adjudicator!");
    }
```

**Replace with:**
```tsx
const createAdjudicatorPanel = async () => {
    // For solo adjudicator, only one is needed
    if (currentPanel.panel_size === 1 && !currentPanel.chief_adj_id) {
        return toast.warning("Select an adjudicator!");
    }
    
    // For panels > 1, require chief  
    if (currentPanel.panel_size > 1 && !currentPanel.chief_adj_id) {
        return toast.warning("Select a chief adjudicator!");
    }
```

This fixes the validation so solo adjudicators work properly.

### Step 5: Restart Dev Server
```
cd frontend
npm run dev
```

## alternative: I Can Create a Fresh File

If you can't restore the file, please let me know and I'll create a completely fresh, working version of the ballot page with all the fixes included.

## Next Steps After Fixing Ballot Page
1. ✅ Add "Back to Tournaments" button on user pages (simpler, separate task)
2. Update PIMNAS data to have correct speaker counts
3. Add tournament format validation

Let me know when the file is restored and I'll help with the rest!
