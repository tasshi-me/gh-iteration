# gh iteration

## Documentation

https://tasshi-me.github.io/gh-iteration/

## Installation

```shell
gh extension install tasshi-me/iteration
$ gh auth refresh -s project
```

## Example

```shell
# Assign current sprint to all in progress tasks

gh iteration items-edit \
  --owner "myOrg" \
  --project "123" \
  --field "Sprint" \
  --query "(Item.Type == \"ISSUE\") && (
     (Item.Fields.Status.Name endsWith \"In progress\")
     || (Item.Fields.Status.Name endsWith \"In review\")
     || (Item.Fields.Status.Name endsWith \"In testing\")
     || (Item.Fields.Status.Name endsWith \"In AC Check\"))" \
  --current
```

## License

- [MIT](./LICENSE)
