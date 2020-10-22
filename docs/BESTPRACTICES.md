# Best Practices for Developing Bubbly

## Branching strategy
We are following the GitHubFlow approach to branching. For more reading, see: [GitHubFlow](https://guides.github.com/introduction/flow/)

### Best Practices for Branches
When working on a `feature` branch, it is helpful to keep the branch small, and minimize the amount of time that code is checked out. Frequent merges and code-reviews keep the PR's manageable, and improve code quality.

When working on a feature branch, rebase against development often
`$ rebase origin/master`

It is considered best practice to `squash` your commits into 1 commit before creating a PR
`git rev-list --count HEAD ^master` # Returns the count of commits past develop's HEAD
`git rebase -i HEAD~[N]` # Where N is the amount of commits past develop's HEAD
