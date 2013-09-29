gor
===

## Execute previous go run command
When learning go or running trials with small programs, one usually has a few different .go files in the same directory.   

This program creates a .gorrc in your user directory that contains the last run/build/test file for your present working directory.  Thereafter, just run gor (or with -b/-t option for build and test respectively).  An entry exists per directory and per tool (run/build/test/doc) that you use.

## Installation
```
go get github.com/sathishvj/gor
```

> This will pull down the source and install the gor command in $GOROOT/bin

## Running gor
Make sure that your PATH contains $GOPATH/bin.  Then use gor from any directory.

```
gor hello.go 
```
> this will run hello.go for the first time and add it to .gorrc in your user folder.

```
gor
```

> this will re-run hello.go

```
gor -b hello.go
```

> this will build hello.go for the first time.

```
gor -b 
```

> this will re-build hello.go

Tip: After using "gor -b" for the first time, you might want to alias gob so that you can use just gob each time.  Since alias cannot accept parameters though, for the first time you build a file, you will have to use "gor -b"
```
alias gob="gor -b"
```

