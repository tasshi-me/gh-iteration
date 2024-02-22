# gh iteration

## Usage

- gh iteration field-list (in development)
- gh iteration field-view (in development)
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
gh iteration field list
```
#### Options

```
--owner
--project
--project-id
```

### gh iteration field-view

View a iteration field

```shell
gh iteration field view
```

#### Options

```
--owner
--project
--project-id
--name
--id
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
--project-id
--field (name)
--field-id
--completed

--field-id | (--field & (--project-id | (--project & --owner)))
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
--project-id
--field (name)
--field-id
--iteration (title)
--iteration-id
--iteration-start-date
--iteration-current
--clear
```
