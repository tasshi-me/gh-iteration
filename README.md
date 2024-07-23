# gh iteration

## Usage

- gh iteration field-list
- gh iteration field-view
- gh iteration list
- gh iteration item-edit (in development)

See more details in https://tasshi-me.github.io/gh-iteration/

### Global options

```
--verbose
--format plain|json
--log-level debug|info|warn|error|none
--log-format plain|json
```

### gh iteration field-list

List the iteration fields in a project

```shell
gh iteration field-list
```
#### Options

```
--owner
--project (number)
```

### gh iteration field-view

View an iteration field

```shell
gh iteration field-view
```

#### Options

```
--owner
--project (number)
--field (name)
```

### gh iteration list

List the iterations for an iteration field

```shell
gh iteration list
```

#### Options

```
--owner
--project (number)
--field (name)
--completed

--field & --project & --owner
```

### gh iteration item-edit

Edit iteration of a project item

```shell
gh iteration item-edit
```

#### Options

```
--owner
--project (number)
--field (name)
--iteration (title)
--iteration-start-date
--iteration-current
--clear
```
