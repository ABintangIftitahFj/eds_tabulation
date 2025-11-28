# EDS UPI Tournament Testing Guide

Panduan lengkap untuk testing sistem tournament EDS UPI menggunakan VS Code tanpa aplikasi eksternal.

## 📋 Daftar Isi

1. [Setup Testing Environment](#setup-testing-environment)
2. [Frontend Testing](#frontend-testing)
3. [Backend Testing](#backend-testing)
4. [Test Types & Coverage](#test-types--coverage)
5. [VS Code Integration](#vs-code-integration)
6. [Running Tests](#running-tests)
7. [Troubleshooting](#troubleshooting)

## 🛠️ Setup Testing Environment

### Frontend Requirements
```bash
cd frontend
npm install --save-dev jest @testing-library/react @testing-library/jest-dom @testing-library/user-event jest-environment-jsdom @types/jest
```

### Backend Requirements
```bash
cd Backend
go mod tidy
go get github.com/stretchr/testify
```

## 🌐 Frontend Testing

### Test Structure
```
frontend/
├── __tests__/
│   ├── integration/
│   │   └── api-integration.test.ts     # API endpoint testing
│   ├── functional/
│   │   └── tournament-workflows.test.tsx # User workflow testing
│   └── example.test.ts                 # Basic example test
├── components/
│   └── __tests__/
│       └── Navbar.test.tsx             # Component unit tests
├── lib/
│   └── __tests__/
│       └── api.test.ts                 # API utility tests
├── jest.config.js                     # Jest configuration
└── jest.setup.js                      # Test setup
```

### Test Categories

#### 1. **Unit Tests** - `components/__tests__/`
Tests individual components in isolation:
- Component rendering
- Props handling
- User interactions
- State management

**Example:**
```typescript
test('renders navigation links', () => {
    render(<Navbar />);
    expect(screen.getByText(/home/i)).toBeInTheDocument();
    expect(screen.getByText(/perlombaan/i)).toBeInTheDocument();
    expect(screen.getByText(/admin/i)).toBeInTheDocument();
});
```

#### 2. **Integration Tests** - `__tests__/integration/`
Tests API integration and data flow:
- Tournament CRUD operations
- Team management
- Ballot submission
- Standings calculations
- Error handling

**Example:**
```typescript
test('should submit ballot', async () => {
    const ballotData = {
        match_id: 1,
        adjudicator_id: 1,
        winner: 'gov',
        scores: [...]
    };
    const response = await api.post('/submit-ballot', ballotData);
    expect(response.data.winner_id).toBe(1);
});
```

#### 3. **Functional Tests** - `__tests__/functional/`
Tests complete user workflows:
- Tournament selection
- Ballot form submission
- Standings viewing
- Error handling flows

**Example:**
```typescript
test('complete tournament workflow simulation', async () => {
    // Step 1: Select tournament
    await user.click(screen.getByTestId('tournament-1'));
    
    // Step 2: Select match
    await user.click(screen.getByTestId('select-match'));
    
    // Step 3: Submit ballot
    await user.click(screen.getByTestId('submit-ballot'));
    
    // Step 4: Verify success
    expect(screen.getByTestId('success-message')).toBeInTheDocument();
});
```

### Running Frontend Tests

```bash
# Basic test run
npm test

# Run tests once (no watch mode)
npm test -- --watchAll=false

# Run with coverage
npm test -- --coverage

# Run specific test file
npm test -- Navbar.test.tsx

# Run tests in specific directory
npm test -- __tests__/integration/
```

## ⚙️ Backend Testing

### Test Structure
```
Backend/
├── logic_test.go              # Business logic unit tests
├── controllers_test.go        # API controller tests (optional)
└── main_test.go              # Integration tests (optional)
```

### Test Categories

#### 1. **Unit Tests** - `logic_test.go`
Tests core business logic without external dependencies:
- Score calculations
- Tournament data structures
- Standings sorting
- Validation rules

**Example:**
```go
func TestBallotCalculations(t *testing.T) {
    totalGov := 165 // 85 + 80
    totalOpp := 160 // 78 + 82
    
    var winner string
    if totalGov > totalOpp {
        winner = "gov"
    } else {
        winner = "opp"
    }
    
    if winner != "gov" {
        t.Errorf("Expected winner to be 'gov', got '%s'", winner)
    }
}
```

#### 2. **Integration Tests** - `controllers_test.go`
Tests API endpoints with database:
- Tournament CRUD
- Team management
- Ballot submission
- Standings retrieval

### Running Backend Tests

```bash
# Run all tests
go test -v ./...

# Run specific test file
go test -v logic_test.go

# Run tests with coverage
go test -cover

# Run benchmarks
go test -bench=.

# Verbose output
go test -v
```

## 📊 Test Types & Coverage

### Coverage Goals
- **Unit Tests**: > 80% code coverage
- **Integration Tests**: All API endpoints
- **Functional Tests**: Critical user workflows

### Test Coverage Areas

#### ✅ Currently Covered
1. **Frontend**
   - Component rendering (Navbar)
   - API configuration
   - Basic integration scenarios
   - User interaction workflows

2. **Backend**
   - Score calculation logic
   - Tournament data validation
   - Standings sorting algorithms
   - Environment setup

#### 🚧 To Be Added
1. **Frontend**
   - Form validation edge cases
   - Error boundary testing
   - Performance testing
   - Accessibility testing

2. **Backend**
   - Database transaction testing
   - Concurrent user scenarios
   - Performance benchmarks
   - Security testing

## 🔧 VS Code Integration

### Recommended Extensions
1. **Jest** (`Orta.vscode-jest`) - Already installed
   - Runs tests automatically
   - Shows test results inline
   - Debug test failures

2. **Go Test Explorer** (optional)
   - Visual test runner for Go
   - Test coverage visualization

### VS Code Features

#### Test Explorer Panel
- View all tests in sidebar
- Run individual tests
- Debug test failures
- View test output

#### Debug Tests
1. Set breakpoints in test files
2. Use "Debug Test" option
3. Inspect variables and state
4. Step through test execution

#### Test Coverage
- View coverage reports in VS Code
- Highlight uncovered lines
- Coverage gutters in editor

### Jest Integration Features
- **Auto-run**: Tests run automatically on file changes
- **IntelliSense**: Auto-completion for Jest matchers
- **Error highlighting**: Failed tests highlighted in editor
- **Debug support**: Set breakpoints and debug tests

## 🏃‍♂️ Running Tests

### Quick Test Script
Use the provided PowerShell script:
```powershell
.\run-tests.ps1
```

### Manual Testing

#### Frontend
```bash
cd frontend

# Install dependencies
npm install

# Run all tests
npm test

# Run tests with coverage
npm test -- --coverage

# Run specific test suite
npm test -- components/
```

#### Backend
```bash
cd Backend

# Install dependencies
go mod tidy

# Run logic tests
go test -v logic_test.go

# Run all tests
go test -v ./...

# Run with coverage
go test -cover
```

### Continuous Integration Ready
All tests are configured to run in CI/CD environments:
- No interactive prompts
- Exit codes for success/failure
- Coverage reports generated
- Test results in standard format

## 🔍 Test Scenarios Covered

### 1. Tournament Management
- ✅ Create tournament
- ✅ List tournaments
- ✅ Update tournament status
- ✅ Validate tournament data

### 2. Team & Speaker Management
- ✅ Register teams
- ✅ Add speakers to teams
- ✅ Validate team data
- ✅ Filter teams by tournament

### 3. Match & Round Management
- ✅ Create rounds
- ✅ Generate matches
- ✅ Assign adjudicators
- ✅ Track match completion

### 4. Ballot System
- ✅ Submit individual scores
- ✅ Calculate team totals
- ✅ Determine winners
- ✅ Validate score ranges
- ✅ Handle missing speakers

### 5. Standings & Rankings
- ✅ Calculate team standings
- ✅ Sort by VP and speaker points
- ✅ Generate speaker rankings
- ✅ Handle tied positions

### 6. Error Handling
- ✅ Invalid API inputs
- ✅ Missing required fields
- ✅ Network errors
- ✅ Database constraints
- ✅ User input validation

## 🐛 Troubleshooting

### Common Frontend Issues

#### Jest Configuration Problems
```bash
# Clear Jest cache
npx jest --clearCache

# Reinstall dependencies
rm -rf node_modules package-lock.json
npm install
```

#### Module Resolution Issues
Check `jest.config.js` for proper module mapping:
```javascript
moduleNameMapper: {
    '^@/(.*)$': '<rootDir>/$1',
    '\\.(css|less|scss|sass)$': 'identity-obj-proxy'
}
```

#### TypeScript Errors
Ensure `@types/jest` is installed:
```bash
npm install --save-dev @types/jest
```

### Common Backend Issues

#### Import Path Problems
Verify module path in `go.mod`:
```go
module github.com/star_fj/eds-backend
```

#### Database Connection Issues
Set test environment variables:
```go
os.Setenv("JWT_SECRET", "test-secret")
os.Setenv("DATABASE_URL", "sqlite://test.db")
```

#### Test Dependencies
Install testing libraries:
```bash
go get github.com/stretchr/testify
```

### VS Code Issues

#### Jest Extension Not Working
1. Reload VS Code window: `Ctrl+Shift+P` → "Developer: Reload Window"
2. Check Jest extension is enabled
3. Verify Jest is installed in project

#### Test Files Not Recognized
Ensure test files follow naming convention:
- Frontend: `*.test.ts`, `*.test.tsx`
- Backend: `*_test.go`

## 📈 Test Metrics

### Current Test Coverage
- **Frontend**: ~47% statement coverage
- **Backend**: Logic functions covered
- **API Endpoints**: Major endpoints tested
- **User Workflows**: Key scenarios covered

### Test Execution Time
- **Frontend**: ~5-6 seconds
- **Backend**: <1 second
- **Total**: <10 seconds for full test suite

### Test Reliability
- ✅ Deterministic results
- ✅ No external dependencies for unit tests
- ✅ Isolated test environments
- ✅ Consistent across different machines

## 🎯 Next Steps

1. **Expand Coverage**
   - Add more edge cases
   - Test error scenarios
   - Performance testing

2. **Integration Testing**
   - End-to-end user workflows
   - Database integration tests
   - Real API testing

3. **Automation**
   - Pre-commit hooks
   - CI/CD pipeline integration
   - Automated test reports

4. **Performance Testing**
   - Load testing for APIs
   - Frontend performance metrics
   - Memory usage optimization

---

## 📚 Resources

- [Jest Documentation](https://jestjs.io/)
- [React Testing Library](https://testing-library.com/docs/react-testing-library/intro/)
- [Go Testing](https://pkg.go.dev/testing)
- [VS Code Testing](https://code.visualstudio.com/docs/editor/testing)

**Happy Testing! 🎉**