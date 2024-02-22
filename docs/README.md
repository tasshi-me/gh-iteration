## gh iteration

GitHub CLI extension to list/view/edit iterations in GitHub Projects.


### Index

- [Commands](#commands) 
- [Installation](#installation)
- [Examples](#examples)
- [Links](#links) 

### Commands

|Command|Description|
|-|-|
|[gh iteration list](gh_iteration_list.md)|List the iterations for an iteration field|
|[gh iteration update](gh_iteration_update.md)|A brief description of your command|

### Installation

#### 1. Install GitHub CLI (`gh`)

See the [document](https://github.com/cli/cli).

#### 2. Install gh-iteration

```shell
gh extension install tasshi-me/iteration
````
#### 3. Refresh token with `project` scope (if required)

To access to GitHub Projects, your token must have `project` scope.

You can check current your token scopes by:
```shell
$ gh auth status
```
If `project` is not in token scopes, add `project` scope by:
```shell
$ gh auth refresh -s project
```


### Examples

```shell
# Update the iteration of the project item to current iteration.
$ gh iteration item-edit \
    --owner <OWNER> \
    --project <PROJECT_NUM> \
    --field <FIELD_NAME> \
    --item-id <ITEM_ID> \
    --iteration-current

# If you have field ID, you can simplify.
$ gh iteration item-edit --field-id <FIELD_ID> --item-id <ITEM_ID> --iteration-current
```


### Links

- [Repository](https://github.com/tasshi-me/gh-iteration)
- [Release notes](https://github.com/tasshi-me/gh-iteration/releases)

