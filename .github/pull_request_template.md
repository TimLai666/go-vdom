# Pull Request

## Description

<!-- Provide a brief description of your changes -->

## Type of Change

<!-- Mark the relevant option with an "x" -->

- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Documentation update
- [ ] Code refactoring
- [ ] Performance improvement
- [ ] Test addition or modification

## Related Issues

<!-- Link to related issues using #issue_number -->

Closes #
Related to #

## Changes Made

<!-- Describe what you changed and why -->

-
-
-

## Testing

<!-- Describe the tests you ran and their results -->

- [ ] All existing tests pass (`go test ./...`)
- [ ] Template linter passes (`./tools/template-linter/template-linter -fix ./`)
- [ ] Code is formatted (`go fmt ./...`)
- [ ] Code passes vet (`go vet ./...`)
- [ ] New tests added for new functionality
- [ ] Manual testing performed

### Test Results

```bash
# Paste test output here
```

## Checklist

<!-- Mark completed items with an "x" -->

### Code Quality

- [ ] Code follows the project's style guidelines
- [ ] Self-review of my own code completed
- [ ] Code is well-commented, particularly in hard-to-understand areas
- [ ] No unnecessary console logs or debug code

### Templates (if applicable)

- [ ] Template variables use correct syntax (`{{var}}` not `'{{var}}'`)
- [ ] Boolean comparisons use `{{flag}} === true` not `'{{flag}}' === 'true'`
- [ ] String concatenation handled in Go code, not in `${...}` expressions
- [ ] Template linter reports no issues

### Documentation

- [ ] Updated relevant documentation (README, docs/, etc.)
- [ ] Added/updated code comments where necessary
- [ ] Updated CHANGELOG.md (if applicable)
- [ ] Added examples or usage instructions (if new feature)

### Testing

- [ ] Added tests that prove the fix/feature works
- [ ] All tests pass locally
- [ ] Tests cover edge cases
- [ ] No test coverage regression

## Screenshots/Examples (if applicable)

<!-- Add screenshots or code examples demonstrating your changes -->

```go
// Example code here
```

## Breaking Changes

<!-- If this PR introduces breaking changes, describe them here -->

- None

OR

- **What breaks**:
- **Migration path**:
- **Deprecation notice**:

## Additional Context

<!-- Add any other context about the PR here -->

## Reviewer Notes

<!-- Any specific areas you'd like reviewers to focus on? -->

---

## Pre-Submission Checklist

Before submitting this PR, ensure you have:

1. **Run all checks locally**:
   ```bash
   # Linux/Mac
   make check

   # Windows PowerShell
   .\dev.ps1 check

   # Windows Batch
   dev.bat check
   ```

2. **Verified template quality**:
   ```bash
   cd tools/template-linter
   ./template-linter -fix ../../
   ```

3. **Reviewed your own code** carefully

4. **Tested manually** (if applicable)

5. **Updated documentation** as needed

---

**By submitting this PR, I confirm that:**
- [ ] My code follows the project's contribution guidelines
- [ ] I have performed a self-review of my own code
- [ ] I have tested my changes thoroughly
- [ ] I am willing to address any feedback from reviewers
