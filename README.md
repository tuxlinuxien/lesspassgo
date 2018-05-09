LessPassGo [![Build Status](https://travis-ci.org/tuxlinuxien/lesspassgo.svg?branch=master)](https://travis-ci.org/tuxlinuxien/lesspassgo) [![GoDoc](https://godoc.org/github.com/tuxlinuxien/lesspassgo?status.svg)](https://godoc.org/github.com/tuxlinuxien/lesspassgo) [![Go Report Card](https://goreportcard.com/badge/tuxlinuxien/lesspassgo)](https://goreportcard.com/report/github.com/tuxlinuxien/lesspassgo)
==========

Password generator implemented in Go based on [LessPass](//github.com/lesspass/cli).

Now support **v1** and **v2**.

## Install

```
$ go get -u github.com/tuxlinuxien/lesspassgo/...
```

## CLI

### Usage

```
NAME:
   lesspassgo - LessPass password generator CLI.

USAGE:
   lesspassgo [options]

AUTHOR:
   Yoann Cerda <tuxlinuxien@gmail.com>

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --site value               
   --login value              
   --password value           
   --counter value, -c value  (default: 1)
   --length value, -L value   (default: 16)
   --version1, --v1           
   --upper, -u                
   --lower, -l                
   --numbers, -n              
   --symbols, -s              
   --help, -h                 show help
```

## Server

API compatible with LessPass web-browser plug-in.
Users and passwords are stored into sqlite3 database.

```
NAME:
    lesspassgo-server - LessPass password server.

USAGE:
    lesspassgo-server [options]

AUTHOR:
    Yoann Cerda <tuxlinuxien@gmail.com>

COMMANDS:
    help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
    --db value              (default: "./user.db") [$DB_PATH]
    --host value            (default: "0.0.0.0") [$HOST]
    --port value, -p value  (default: 1314) [$PORT]
    --disable-registration  [$DISABLE_REGISTRATION]
```

## License

Licensed under the MIT License.
