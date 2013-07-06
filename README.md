toil [Working Title] 
====

![Toil](http://dl.dropbox.com/u/3155323/HostedShots/toil.gif)

Run multiple processes, defined by a simple json file, with a single log stream. Super useful for people who change from project to project.

## Quick Start
### Install
A binary for Mac OS x64 is included in the [releases section](https://github.com/snikch/toil/releases).

#### From Source
You’ll need [Go installed](http://golang.org/doc/install).

```
go get github.com/snikch/toil
```

### Running
Navigate to your project root.

```
toil init
toil add redis redis-server
toil add rails "bundle exec rails server"
toil
```

## Commands

### init
Initialises the `toil.json` file.

### add [name] [command]
Adds the given command with the given name. Use quotes for commands with spaces.

```
toil add redis redis-server
➜ ✔ Added redis

toil add rails "bundle exec rails server"
➜ ✔ Added rails
```

### rm [name]
Removes the command with the given name.

```
toil rm potato
➜ ✘ No process named potato

toil rm rails
➜ ✔ Removed rails
```

### list
Lists the current processes.

```
toil list
➜ 2 processes in toilfile
➜ Local:
➜  - rails: bundle exec rails s -p 3001
➜  - redis: redis-server
 ```

## TODO
* `toil register [project name]`
* `toil on [project name]`
*  Handle global processes when starting another instance of toil
