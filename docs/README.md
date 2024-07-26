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
|[gh iteration field-list](gh_iteration_field-list.md)|List the iteration fields in a project|
|[gh iteration field-view](gh_iteration_field-view.md)|View an iteration field|
|[gh iteration item-edit](gh_iteration_item-edit.md)|Edit iteration of a project item|
|[gh iteration item-view](gh_iteration_item-view.md)|View a project item|
|[gh iteration items-edit](gh_iteration_items-edit.md)|Edit iteration of multiple project items|
|[gh iteration list](gh_iteration_list.md)|List the iterations for an iteration field|

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

