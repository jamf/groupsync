# Contributing
Before making a pull request, please first discuss the change via an issue.

If you have questions, ask! Your level of experience is irrelevant. Questions
are good. We're happy to help.

## Commits
We leverage
[this rather amazing GitHub action](https://github.com/marvinpinto/action-automatic-releases/)
for automated releases. One thing it does is gather commits between releases
and create a changelog out of them. To keep changelogs neat, please label your
commits. Here's an example:

```
fix: groupsync now responds appropriately to HTTP 418
```

Available labels:

```
    ConventionalCommitTypes["feat"] = "Features";
    ConventionalCommitTypes["fix"] = "Bug Fixes";
    ConventionalCommitTypes["docs"] = "Documentation";
    ConventionalCommitTypes["style"] = "Styles";
    ConventionalCommitTypes["refactor"] = "Code Refactoring";
    ConventionalCommitTypes["perf"] = "Performance Improvements";
    ConventionalCommitTypes["test"] = "Tests";
    ConventionalCommitTypes["build"] = "Builds";
    ConventionalCommitTypes["ci"] = "Continuous Integration";
    ConventionalCommitTypes["chore"] = "Chores";
    ConventionalCommitTypes["revert"] = "Reverts";
```

## Pull requests
Pull requests should target the `master` branch. Every once in a while
a release will be kicked off from it (by pushing a new git tag), following
basic [semver v2 spec](https://semver.org/).

## Code of Conduct
1. You DO NOT talk about Fight Club.
