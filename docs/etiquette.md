# Git Etiquette

Git is a powerful tool allowing for several different ways of achieving the same version control goals. Over time sets of standards have been developed which make working with Git as an individual and a team easier. [Github flow](https://docs.github.com/en/get-started/quickstart/github-flow) is the Git etiquette standard used within this course and is encouraged to be used within all of your projects and exercises.

Github flow has several rules associated with proper usage, however your grade will only depend on a subset of these rules for each stage of the project. To understand the requirements for each release refer to the table below:

Categories are categories under which you will be marked and will appear on your marking transcript. Categories collect a series of rules that will be enabled or disabled for marking depending on which phase (Alpha, Beta, Final Project) of the project you are working in.

| Category       | Rule                        | Alpha | Beta | Final Project |
| -------------- | --------------------------- | ----- | ---- | ------------- |
| Commit         | Commit Atomicity            | ✓     | ✓    | ✓             |
| Commit         | No Binaries                 | ✓     | ✓    | ✓             |
| Commit         | Empty Commits               | ✓     | ✓    | ✓             |
| Commit Message | Commit Message Matches Diff |       | ✓    | ✓             |
| Commit Message | Short Commit Message        | ✓     | ✓    | ✓             |
| Branching      | Stale Branch                |       |      | ✓             |
| Branching      | Consistent Branch Name      |       | ✓    | ✓             |
| Pull Request   | Feature Branching           |       | ✓    | ✓             |
| Pull Request   | Linked Pull Request Issue   |       | ✓    | ✓             |
| Pull Request   | Approved Pull Request       |       | ✓    | ✓             |
| Pull Request   | Resolved Conversations      |       | ✓    | ✓             |
| Pull Request   | Constructive Feedback       |       | ✓    | ✓             |
| General        | Force Pushes                |       | ✓    | ✓             |
| General        | Unresolved Commit           | ✓     | ✓    | ✓             |

## Outline

The following section gives an outline and some guidance to each of the Git etiquette rules that you should be following for the duration of your projects.

### Commit Atomicity

Commits should focus on a single small ‘task’ (or part of a task) that you are trying to complete. Ensure that commits do not encompass multiple tasks by becoming too large and fragmented by properly breaking down the task into smaller commits. This best practice is encouraged fo the following reasons:

1. Ensures that changes are easy to revert or roll back at any point, such that part of a task can be rolled back without requiring that other parts of the task are also rolled back.
2. It makes it easier to merge work. Commits are applied individually when merging to a branch which means that if a commit conflicts, the type of changes conflicting are logically consistent and minimal in code size.

#### Guidance

- Roughly plan the stages of your task and commit when you feel that you have reached a ‘checkpoint’ where you would feel comfortable reverting to at any point in the future.
- Ensure that your commits do not leave the project in a broken state. If you are in a place where you are having to revert commits you’re not going to be happy when you have to fix additional bugs from a previous commit.

#### More Information

1. A guide on how to think about atomic commits: https://www.freshconsulting.com/insights/blog/atomic-commits/

### No Binaries

When compiling your program it may generate binary artifacts. Generated files should be excluded from your git tree and history as they are files that change between builds and are not necessary to track version history of via Git.

#### Guidelines

- Create a .gitignore file which allows Git to ignore certain files from your project directory. Using https://www.toptal.com/developers/gitignore may help you generate a gitignore.

### Empty Commits

Empty commits do not add code the the repository and serve little purpose in the Git history apart from possibly triggering a continuous integration (CI) action. For the purpose of your assignments excessive usage of empty commits is considered to be forging your git development history and should be avoided.

### Commit Message Matches Diff

When creating a commit message, ensure that the message is informative and relevant to the changes that the commit encompasses. Conventional commits will have verbiage surrounding a reference to changed code. We discourage very terse messages and irrelevant messages when committing as others seeing your commits will not have enough information to work off.

#### More Information

1. A good standard for conventional commit messages: https://www.conventionalcommits.org/en/v1.0.0/#summary

### Short Commit Message

Commit messages should not be less than **three** words. This is an arbitrary word count for a commit message, however it is helpful as it ensures you are providing value to your peers by creating succinct descriptions of your commit content.

### Stale Branches

Stale branches are branches that are considered to be unused and have not received any commit activity within **two weeks**. Stale branches clutter your Git project by making it difficult to search for active development branches. A good Git repository will only have branches in active development which additionally allows you to gather insight into what features are being worked on by team members. There are two options if you have a stale branch:

1. Abandon development of the branch by deleting the branch `git push –delete origin <branch-name>`.
2. Continue working on the branch. Ensure that the history of the ophan branch is updated by rebasing to the source branch. After development on the branch is completed, create a pull request, get your peers to review the code, merge it into your primary branch (main) and delete the branch.

#### Guidelines

Github has a helpful UI for deleting the source branch when a pull request has been merged. It is good practice to use this feature to keep your active branches tidy:

<p align="center">
  <img src="https://i.imgur.com/1PHi9oM.png" />
</p>

### Consistent Branch Names

Correctly name branches is important for ease of access and being able to see at a glance what a branch contains in code changes. An additional enhancement that you can make to branch names is using **group tokens**. Although not every branch requires group tokens, it gives the ability to group together families of branches. For example the following branch names:

- feature/context-change
- feature/cli-subcommand
- fix/markdown-generator
- fix/http-headers

Are branches which have the groups tokens **feature** and **fix** with a separator token _/_ (slash). Group tokens also allow you to quickly list and access all branches within a group by using tab expansion:

- `git checkout fix/<TAB>`

Will list all of the branches that exist in the fix group which is helpful to see what fixes there are on your version control tree.

#### Guidelines

- Organize with your group which branch naming contentions and naming groups you would like to use, including what group tokens and separator conventions you should use
- Ensure that your branch names and group tokens are short enough to allow people to quickly understand names and change branches. Short branch names will additionally ensure that the branch name does not pollute logs.

#### More Information

1. An overview of good branch naming conventions: https://stackoverflow.com/a/6065944

### Feature Branching

When adding code to the codebase a branch should be created, code committed to that branch and then a pull request should be made to your primary branch which is reviewed and merged. This forms the basic etiquette of the Github workflow. Directly committing features to the main branch is discouraged as it muddles development history and makes it difficult to organize the ordering of features to be put into your project.

#### More information

1. An overview of basic branch, pull request workflow: https://www.atlassian.com/git/tutorials/comparing-workflows/feature-branch-workflow

### Linked Pull Request Issue

When creating a pull request good practice in some cases (creating features, implementing fixes) to link to an existing issue in your description when creating a pull request. This gives context for and describes the nature of the pull request and allows you to link multiple pull requests to a single issue if it requires multiple atomic changes to the project.

#### Automatically linking:

<p align="center">
  <img src="https://i.imgur.com/boLVpLx.png" />
</p>

#### Manually Linking

<p align="center">
  <img src="https://i.imgur.com/F9dg0tT.png" />
</p>

#### Guidelines

1. Issues can be quickly added to a pull request by typing in the # symbol in the description dialog. This will display a list of the issues that exist on your git repository and allows for quick selection.
2. Issues can be manually added by selecting an issue from the `development` dropdown from the right sidebar of the pull request.

#### More Information

1. Pull request issue automatic link: https://docs.github.com/en/get-started/writing-on-github/working-with-advanced-formatting/autolinked-references-and-urls
2. Pull request Issue manually link: https://docs.github.com/en/issues/tracking-your-work-with-issues/linking-a-pull-request-to-an-issue

### Review Pull Requests

Pull requests should be reviewed by at least one other person who is not the author of the PR before merging. Reviewing code before merging can catch issues with code that were not apparent to the author, when resolved creates a higher quality codebase. Reviews should aim to constructively critique the state of the code and suggest improvements if any are required.

<p align="center">
  <img src="https://i.imgur.com/N1HQl50.png" />
</p>

### More Information

1. Pull request reviews: https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/reviewing-changes-in-pull-requests/about-pull-request-reviews

### Resolved Conversations

Reviews can also take place in the form of creating a series of comment threads (conversations) that relate to a part of the code contained within the pull request. This is a good practice as it allows more direct feedback to the code author. Before merging these conversations should be marked as resolved before the pull request is merged to indicate that all discussion has finished and that there are no remaining issues.

<p align="center">
  <img src="https://i.imgur.com/PNLx7ts.png" />
</p>

### Constructive Feedback

Try to make your review feedback as constructive as possible. Aiming to clearly state what is to be improved with the code without leaving the author of the code to feel bad about their work. Ideally your feedback should be **actionable**, letting your group partner know exactly what is wrong and a suggestion on how to fix it. Constructive feedback allows your team to maintain a positive working environment and aids your productivity, making it a better expedience for everyone involved.

#### Guidelines

1. A guide for creating good feedback: https://about.gitlab.com/handbook/people-group/guidance-on-feedback/

### Force Pushes

Force pushes are heavily discouraged as you risk losing commits as the branch history is overwritten. Instead of force, pushing consider using [git-revert](https://www.atlassian.com/git/tutorials/undoing-changes/git-revert) to revert the commits that have been made to the branch by creating a new commit. This way there is no risk of losing any commits that were made on a branch.

#### Guidelines

If a force push has overwritten somebodies work and you wish to retrieve that work there is still hope. Force pushes tell you which commit hash the force push updated from and to meaning it is easy to change the branch back to the old commit (see [here](https://www.jvt.me/posts/2021/10/23/undo-force-push/))

#### More Information

1. Basic git revert usage: https://www.atlassian.com/git/tutorials/undoing-changes/git-revert
2. Undoing a force push: https://www.jvt.me/posts/2021/10/23/undo-force-push/

### Unresolved Conflicts

When a merge conflict occurs it is easy to commit an unresolved conflict if you are not careful about resolving them before continuing a merge or commit. This is bad practice as it breaks the code by introducing symbols that will make the project fail to compile which will additionally need to be resolved at a later date. An example is below of a merge conflict that should needs resolved before continuing a merge or committing:

<p align="center">
  <img src="https://i.imgur.com/N814CZo.png" />
</p>

#### More Information

1. Resolving merge conflicts: https://www.atlassian.com/git/tutorials/using-branches/merge-conflicts
