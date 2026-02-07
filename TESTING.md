# Test Coverage Roadmap

## Current Status
- **Overall Coverage**: 27.0% (improved from 22.9%)
- **FS Layer Progress**: 74.2% (improved from 13.4%)
- **All Tests Passing**: ✓

### Package Coverage
| Package | Coverage | Status |
|---------|----------|--------|
| internal/domain | 85.7% | ✓ Strong |
| internal/util | 67.4% | ✓ Good |
| internal/cmd | 45.5% | ⚠ Partial |
| internal/ai | 91.7% | ✓ Excellent |
| internal/fs | 74.2% | ✓ Excellent (was 13.4%) |
| internal/service | 56.6% | ✓ Improving |
| cmd/mandor | 0.0% | ✗ Untested |
| internal/cmd/* | 0.0% | ✗ Untested (14 packages) |

## Path to 100% Coverage

### Phase 1: Foundation (Complete)
- ✓ Add tests for `internal/ai` package (92% coverage achieved)
- ✓ Verify all existing tests pass

### Phase 2: Service Layer (In Progress)
**Status**: 6 of 9 services complete | **Current Coverage**: 56.6%

#### Completed Services

**BriefService** (`internal/service/brief_service.go`)
- 28 test functions covering:
  - Constructor tests (NewBriefServiceWithPaths)
  - Validation tests (backlog not found, empty name, invalid name, why length, no capabilities, empty capability fields, duplicate brief)
  - Create tests (success, capability separation, impact preservation)
  - Read tests (success, not found, empty file)
  - Update tests (status, metadata, timestamp)
  - Delete tests (success, non-existent)
  - Serialization tests (roundtrip, markdown conversion)
  - Edge cases (minimal why, multiple capabilities)

**SpecService** (`internal/service/spec_service.go`)
- 27 test functions covering:
  - Constructor tests (NewSpecServiceWithPaths)
  - Validation tests (backlog not found, empty/invalid capability ID, capability not in brief, empty summary, no requirements, empty requirement summary, missing IAE scenarios, empty intent/action/expect)
  - Create tests (success, multiple requirements, multiple IAE scenarios)
  - Read tests (success, not found)
  - Update tests (status, summary, timestamp)
  - Delete tests (success, not found)
  - Serialization tests (markdown generation, parsing, error cases, roundtrip)
  - Helper function tests (capabilityExistsInBrief, specFileExists)
  - Roundtrip serialization test

**BlueprintService** (`internal/service/blueprint_service.go`)
- 26 test functions covering:
  - Constructor tests (NewBlueprintServiceWithPaths)
  - Workspace initialization test
  - Validation tests (no backlog, empty brief ID, brief not found, missing spec, invalid spec, empty problem statement, no architecture decisions, empty decision title/decision, short rationale)
  - Create tests (success, with risks, constraints, user types, goals, data models, implementation strategy, multiple decisions)
  - Read tests (success, not found)
  - Update tests (success, timestamp)
  - Delete tests (success, non-existent)
  - Serialization tests (markdown generation, parsing, error cases, roundtrip)
  - Helper function tests (verifyAllSpecsExist, verifyAllSpecsValid, loadBrief)
  - File exists tests

**TaskService** (`internal/service/task_service.go`)
- 15 test functions covering:
  - Constructor tests (NewTaskServiceWithPaths)
  - Workspace initialization test
  - ParseTaskID tests (valid, invalid format)
  - Extract backlog ID from feature ID tests (valid, invalid)
  - Priority comparison tests
  - Goal minimum length tests
  - Status transition validation tests
  - Dependency validation tests
  - Cycle detection tests
  - Duplicate name validation tests
  - IAE scenario validation tests
  - Gate checking tests (all gates met, unmet gates)
  - Dependent finding tests

**FeatureService** (`internal/service/feature_service.go`)
- 30 test functions covering:
  - Constructor tests (NewFeatureServiceWithPaths)
  - Workspace initialization test
  - Validation tests (backlog not found, capability ID required/invalid/not in brief, spec ID required/format mismatch, name required, goal required/too short, scope invalid, priority invalid, self-dependency, dependency not found, dependency cancelled)
  - Create tests (success, with blocked dependencies)
  - Status transition validation tests (all valid/invalid transitions)
  - Goal minimum length tests
  - Capability existence in brief tests
  - Duplicate name validation tests
  - Cycle detection tests
  - Dependency validation tests
  - Dependent finding tests
  - Update input validation tests (not found, cancelled feature, invalid status transition)
  - Update tests (name change, dry run, cancel)
  - Delete input validation tests (not found, already cancelled)
  - Delete tests (success)
  - List tests (empty, with features)
  - Detail tests (not found, success)

**WorkspaceService** (`internal/service/workspace_service.go`)
- 20 test functions covering:
  - Constructor tests (NewWorkspaceService)
  - Struct tests
  - InitWorkspace tests (already initialized, success)
  - GetWorkspace tests (not initialized, success)
  - UpdateWorkspaceConfig tests (valid/invalid priority, strict mode, goal lengths, unknown key)
  - GetConfigValue tests (valid, unknown key)
  - GetGoalLength tests (backlog, feature, task, issue, unknown)
  - SetGoalLength tests (success, invalid entity)
  - ResetGoalLength tests (success, invalid entity)
  - ResetAllGoalLengths tests (all entity types)

**IssueService** (`internal/service/issue_service.go`)
- 35+ test functions covering:
  - Constructor tests (NewIssueService, NewIssueServiceWithPaths)
  - Workspace initialization test
  - Validation tests (backlog not found, name required, goal required/too short, issue type required/invalid, affected files/tests/implementation steps required, priority valid/invalid, dependency not found/cancelled, cross-backlog dependency allowed/not allowed)
  - Create tests (success, blocked dependencies, resolved dependency)
  - List tests (empty, with issues, filter by type)
  - Detail tests (success, not found, cancelled not included)
  - Update validation tests (issue not found, cancelled issue, empty name, invalid issue type/priority/status)
  - Update tests (dry run, name change, cancel, start, resolve, wontfix, reopen)
  - Find dependents test
  - Unblock dependents test
  - Backlog exists test

#### Remaining Services
| Service | Functions | Status |
|---------|-----------|--------|
| backlog_service.go | 24 | Already Complete (from earlier) |

### Phase 3: File System Layer (COMPLETED ✓)
**Impact**: +60.8% coverage (13.4% → 74.2%)

**Tested Functions**:

**Reader Tests** (16 test functions):
- ReadWorkspace (success, not found, corrupted JSON)
- ReadBacklogMetadata (success, not found, corrupted JSON)
- ReadBacklogSchema (success, not found)
- CountLines (non-existent, with content)
- CountEntityLines (non-existent, NDJSON with entries)
- ReadNDJSON (success, empty file, not found)
- ReadFeature (success, not found)
- ReadTask (success, not found)
- ReadIssue (success, not found)

**Writer Tests** (22 test functions):
- CreateMandorDir (success)
- WriteWorkspace (success)
- WriteJSON (success, nested directories)
- AppendNDJSON (single entry, multiple entries)
- WriteBacklogMetadata (success)
- WriteBacklogSchema (success)
- CreateBacklogDir (success, verifies entity files)
- DeleteBacklogDir (success)
- WriteFeature/Task/Issue (success, roundtrip verification)
- ReplaceFeature/Task/Issue (success, verifies update)
- ReplaceTasks/Features/Issues (batch update success)
- WriteFile (success, nested directories)

### Phase 4: CLI Commands
**Effort**: 12-16 hours | **Impact**: +25% coverage

Requires E2E tests for:
- 14 cmd/* subpackages (backlog, blueprint, brief, change, feature, issue, populate, session, spec, task, track, workspace, ai)
- Each command's flags and validations
- Error handling and help text

**Approach**:
- Use cobra's built-in test utilities
- Create workspace fixtures for integration tests
- Mock external dependencies (Git, Mandor config)

### Phase 5: Edge Cases & Main
**Effort**: 4-8 hours | **Impact**: +10% coverage

- Test `cmd/mandor/main.go` entry point
- Edge cases and error paths
- Platform-specific behaviors

## Implementation Strategy

### Tools to Use
```bash
# Generate coverage report
go test -v ./... -coverprofile=coverage.out
go tool cover -func=coverage.out | grep "0.0%"

# View HTML coverage
go tool cover -html=coverage.out

# Run specific package tests
go test -v ./internal/service/...

# Run only IssueService tests
go test -v ./internal/service/... -run "Issue"
```

### Testing Patterns Established

1. **Test Setup Pattern**:
```go
func setupTestIssueService(t *testing.T) (*IssueService, string) {
    tmpDir := t.TempDir()
    paths := &fs.Paths{WorkspaceRoot: tmpDir}
    writer := fs.NewWriter(paths)
    reader := fs.NewReader(paths)

    // Create required directories
    os.MkdirAll(filepath.Join(tmpDir, ".mandor", "backlogs"), 0755)

    // Write workspace file
    workspace := &domain.Workspace{...}
    writer.WriteWorkspace(workspace)

    return &IssueService{reader, writer, paths}, tmpDir
}
```

2. **Helper for Backlog + Issue Setup with Cross-Backlog Support**:
```go
func setupTestIssueServiceWithBacklog(t *testing.T, backlogID, issueDepRule string) (*IssueService, *BacklogService) {
    // Creates workspace, backlog with specified dependency rule, returns services
    // issueDepRule: "same_backlog_only", "cross_backlog_allowed", or "disabled"
}
```

3. **Validation Testing**: Test both success and failure paths for all validation rules

4. **Cross-Backlog Dependency Testing**: Set up multiple backlogs with different dependency rules

5. **Roundtrip Testing**: Create → Read → Verify all fields preserved

### File Structure for Tests
```
internal/service/
  ├── backlog_service.go
  ├── backlog_service_test.go      (Complete - 24 tests)
  ├── brief_service.go
  ├── brief_service_test.go        (Complete - 28 tests)
  ├── spec_service.go
  ├── spec_service_test.go         (Complete - 27 tests)
  ├── blueprint_service.go
  ├── blueprint_service_test.go     (Complete - 26 tests)
  ├── task_service.go
  ├── task_service_test.go         (Complete - 15 tests)
  ├── feature_service.go
  ├── feature_service_test.go       (Complete - 30 tests)
  ├── issue_service.go
  ├── issue_service_test.go        (Complete - 35+ tests)
  └── workspace_service.go
    └── workspace_service_test.go   (Complete - 20 tests)
```

## Success Metrics
- [x] All tests pass: `go test ./...`
- [x] BriefService coverage: 80-100%
- [x] SpecService coverage: 80-100%
- [x] BlueprintService coverage: 80-100%
- [x] TaskService coverage: 80-100%
- [x] FeatureService coverage: 80-100%
- [x] IssueService coverage: 80-100% (NEW - just completed)
- [x] WorkspaceService coverage: 80-100%
- [ ] Overall coverage > 80%: `go tool cover -func=coverage.out | grep total`
- [ ] No test warnings or errors
- [ ] CI/CD integration ready

## Priority Ranking (Updated)
1. **Critical (In Progress)**: internal/service/* (56.6% coverage, target: 80%+)
2. **High (Second)**: internal/fs/* (13.4% coverage)
3. **Medium (Third)**: internal/cmd/* (45.5% coverage)
4. **Low (Last)**: cmd/mandor/main.go (0% coverage)

## Estimated Timeline (Updated)
- **Phase 1**: Complete
- **Phase 2**: 8-12 hours (6/9 services complete - 3 remaining)
- **Phase 3**: 4-6 hours
- **Phase 4**: 12-16 hours
- **Phase 5**: 4-8 hours
- **Remaining**: ~20-25 hours → **100% coverage**

## Next Steps (For Continue Session)
1. Move to Phase 3 (fs layer tests)
   - `internal/fs/io.go` (36 untested functions)
   - `internal/fs/paths.go` (17 untested functions)
2. Or continue with remaining services:
   - BacklogService tests (if not already complete)
