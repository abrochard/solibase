# solibase
Go tool to maintain SQL schema under source control. Similar to liquibase.

Keep a single source of truth for your database schema, using the [toml format](https://github.com/toml-lang/toml).

## Installation
```
go get github.com/abrochard/solibase
```

## Usage
You must specify a driver and changelog file.


To apply all changes in your changelog
```
solibase -driver mysql -changelog changelog.toml -mysql-db testDB
```

To rollback to before a specific change
```
solibase -driver mysql -changelog changelog.toml -mysql-db testDB -rollback badChange.toml
```


### Help page
```
Usage of solibase:
  -changelog string
    	location of the changelog
  -driver string
    	driver to use
  -mysql-db string
    	MySQL database
  -mysql-host string
    	MySQL host (defaults to localhost)
  -mysql-password string
    	MySQL password (defaults to no password)
  -mysql-port int
    	MySQL port (defaults to 3306)
  -mysql-user string
    	MySQL username (defaults to root)
  -rollback string
    	rollback to before the specified change
  -sqlite-db string
    	SQLite database
```

## Example changelog and changeset
Full examples under the `examples` folder.
### `changelog.toml`
Lists changeset files using relative path.
```
files = [
       "changeset.toml",
       "changeset2.toml"
]
```
### `changeset.toml`
Contains arbitrary metadata, an optional condition, a sql query to execute, and potentially a rollback.
```
[metadata]
description = "create user table"


[change]

condition = '''
SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME="users"
'''


sql = '''
CREATE TABLE `users` (
	`id` INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
	`username` TEXT(64) NOT NULL,
	`email` TEXT(128) NOT NULL,
	`password` TEXT(254) NOT NULL
);
'''

[rollback]

sql = '''
DROP TABLE users;
'''

```

## Supported Database
- SQLite
- MySQL

If you want to help and extend to more systems, feel free to PR with a package implementing the `solibase.Driver` interface.
