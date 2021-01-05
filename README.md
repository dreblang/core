# Status

[![Build Status](https://travis-ci.com/dreblang/core.svg?branch=master)](https://travis-ci.com/dreblang/core)

# Dreblang core

Compiler written in Golang to learn about the concepts. This is based on the book called **'Writing a Compiler in Go'** by **Thorsten Ball**. Starting from what the book offers, additional data types and features has been added.

> This an experimental language, not to be used in production.

## How to use?

### CLI

```
$ go get github.com/dreblang/core/cmd/drebli
$ drebli
```

### Compile and Execute

```
$ go get github.com/dreblang/core/cmd/dreblc
$ dreblc <file>.dreb
```

Use sample.dreb for code reference. No documentation is available as of now.

Contact me for any queries.


## Feature overview

### Data Types

- **Basic Types** - int, float, bool, string, null
- **Advanced Types** - Array, Hash
- **Callables** - Closure, Member & Built-in Functions

### Control Flow

- if (else), loop, scope, fn
- Supports recursion

Basic arithmatic and comparision operators.
