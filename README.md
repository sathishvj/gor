gor
===

## Execute previous go command
When learning go or running trials with small programs, one usually has a few different .go files in the same directory.  With this utility, just type gor and avoid retyping 'go run filename.go' each time.  

This program creates a .gorrc in your user directory that contains the last run/build/test/doc file for your present working directory.  After you do "gor filename.go" once, just run gor (or with -b/-t/-d option for build/test/doc respectively).  A persistent entry exists per directory and per tool (run/build/test/doc) that you use.  That is helpful when you are working in many directories.

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

Tip: You might want to alias gob="gor -b"  so that you can use just use gob to do builds. Similarly for doc and test, if you so wish.

```
alias gob="gor -b"
```

## Help
```
gor -h
```

## Testing
Checked this only on Mac OS X and with single files. If you find any bugs, please raise an issue for this project.  Thank you.
