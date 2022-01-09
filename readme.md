# Data Mover #
![](./icon_128.png)

Data Mover is an Open Source Data Migration tool.

Data Mover features are:
- **Multi database**: SQLite, SQL Server, MySQL/MariaDB, Postgres
- **Auto Migrate Schema**: can migrate/update database schema
- **Scheduler**: jobs can be scheduled and start every x time or just only once a day at the set time
- **Move Data from Source to Target database**: this is Data Mover main feature 🐼

## How does it work ##

![](./_docs/slide_001.png)

Data Mover hosts "jobs".

A "Job" is basically a JSON file describing what Data Move should do at a predefined time.

Each "job" contains:
- Schedule: optional data to define WHEN the job must start
- Transaction: an array of "Action" to define WHAT the job must do

So, jobs define WHEN and WHAT about Data Mover.

Let's start analyze the anatomy of a Data Mover job.

### Schedule ###

Schedule collects some settings to define a timed task.
Data Mover has an internal task manager and a scheduler working on a thread safe environment.

```json
{
  "schedule": {
    "start_at": "",
    "timeline": "second:3"
  }
}
```

Is quite simple figure how Schedule works:
- start_at: optional value representing hour and minute. ex: "10:20", "18:30", etc..
- timeline: optional value representing a key-pair "unit:value". ex: "millisecond:100", "second:3", "minute:10", "hour:24".

Schedule is optional at all. If you do not specify any value, the job will not be scheduled, but remain a valid job that 
can be invoked from another job (see below about "Job Chains").

### Job Chains ###

```json
{
  "schedule": {
    "start_at": "",
    "timeline": "second:10"
  },
  "next_run": "job-sqlite-users"
}
```

Not all jobs must be scheduled.
Sometimes you should prefer schedule a master job and define different jobs for some other tasks to invoke after the master job.

That's a chain.

"next_run" is the field that tell a job what to do next.

### Transactions ###
```json
{
  "transaction": [
    {
      "uid": "source",
      "description": "select data from source",
      "connection": {
        "driver": "sqlite",
        "dsn": "./source.db"
      },
      "command": "SELECT * FROM users WHERE exported=false",
      "script": ""
    },
    {
      "uid": "target",
      "description": "save source to target",
      "connection": {
        "driver": "sqlite",
        "dsn": "./target.db",
        "schema": {}
      },
      "command": "INSERT INTO users (...) returning id",
      "script": ""
    },
    {
      "uid": "source-update",
      "description": "write flag in source",
      "connection": {
        "driver": "sqlite",
        "dsn": "./source.db"
      },
      "command": "UPDATE users SET exported=true, exported_time='<var>date|iso</var>' WHERE id=@id",
      "script": ""
    }
  ]
}
```

A transaction is an array of actions.

This is a sample action:
```json
{
      "uid": "source",
      "description": "select data from source",
      "connection": {
        "driver": "sqlite",
        "dsn": "./source.db"
      },
      "command": "SELECT * FROM users WHERE exported=false",
      "script": ""
    }
```
This action works on a SQLite db (./source.db) and select all not exported yet data from user table.

That's all. But transactions works using actions that interact together creating an execution context, a transaction context.

The context enable actions to share data during transaction execution. And here comes our first action: this action select data and keep them in context for next action.

Next action will do something new with data in context:

```json
{
  "uid": "target",
  "description": "save source to target",
  "connection": {
    "driver": "sqlite",
    "dsn": "./target.db",
    "schema": {}
  },
  "command": "INSERT INTO users (...) returning id",
  "script": ""
}
```

This action, using context data created from first action, executes an SQL Formula on a target database.

```
INSERT INTO users (...) returning id
```
This SQL Formula is a custom SQL like command with a special statement: `(...)`

The three-dots-statement 🌿 tells the panda 🐼 of Data Mover to extract all fields
and values from the datasource into context and execute an INSERT for each row
in the context.

So Data Mover's panda will start moving row after row into the context to the target database 
executing an INSERT statement for each source row.

Reassuming:
- First we selected some rows from a datasource
- then we started a loop on each row and executed and INSERT into a target database. The insert command was auto-completed from the panda 🐼 of Data Mover that is able to understand a three-dot-statement 🌿.

Now, to complete the transaction, we should need to mark all source data as exported.

And here comes our third action in Transaction:

```json
{
  "uid": "source-update",
  "description": "write flag in source",
  "connection": {
    "driver": "sqlite",
    "dsn": "./source.db"
  },
  "command": "UPDATE users SET exported=true, exported_time='<var>date|iso</var>' WHERE id=@id",
  "script": ""
}
```

The third action, using same context, now try to execute a new UPDATE command into source database setting all fields as "exported".

The "SQL Formula" (a Data Mover special SQL commands) is:
```
UPDATE users SET exported=true, exported_time='<var>date|iso</var>' WHERE id=@id
```

In this statement we have two strange things:
- 🌿 `<var>date|iso</var>`: a Special Expression
- 🌿 `@id`: a named parameter

Nothing magic, just panda 🐼 style.

Data Mover's panda is also able to interpret some Special Expressions like the one above (that return an ISO-8601 date time).

Otherwise, the `@id` named parameter uses the context to get a value for each loop.

That's all.
We just wrote three simple ACTIONs and the panda 🐼 did all the job.

### Schema Migration ###

TODO: add schema migration specifications

### Special Formulas ###

TODO: add Special Formulas documentation

## Binaries ##

Download binaries from this repository in [_build](./_build) directory.

Supported OS:
- Windows and Windows64
- Linux and Linux Embedded Systems
- OSX
- Raspbian

## MIT License NON-COMMERCIAL USE ##

Data Mover is distributed under MIT license fo non-commercial use.
If you use as a tool for your own projects, you can use Data Mover under MIT license.

NON-COMMERCIAL: non-commercial is that no money should be exchanged as part of the transaction of using of the materials – regardless of whether the money represents a break-even of marginal cost, reimbursement or profit.

## Commercial Use License ##

If you are a company that sell projects to its customers and need Data Mover, 
you should ask for a Commercial License.

For Commercial License, please write to [angelo.geminiani@ggtechnologies.sm](mailto:angelo.geminiani@ggtechnologies.sm)

 