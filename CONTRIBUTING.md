# Contributing
Before making a pull request, please first discuss the change via an issue.

If you have questions, ask! Your level of experience is irrelevant. Questions are good. We're happy to help.

## Commits
We leverage [this rather amazing GitHub action](https://github.com/marvinpinto/action-automatic-releases/) for automated releases. One thing it does is gather commits between releases and create a changelog out of them. To keep changelogs neat, please label your commits. Here's an example:

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
For now, please submit your pull requests to the `develop` branch. Once a bunch of changes get tested, they'll get merged to `master` and an automated release should follow.

We might consider ditching `develop` at some point and just merging things directly to `master`.

## Code of Conduct
1. You DO NOT talk about Fight Club.