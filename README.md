gor
===

# Execute previous go run command again
To gain a few seconds from repeatedly retyping "go run filename.go".
This creates a .gorfile in your user directory that contains the last run file.  Thereafter, just run gor.

# Installation
go get github.com/sathishvj/gor

This will pull down the source and install the gor command in $GOROOT/bin

# Running it
Make sure that your PATH contains $GOPATH/bin.  Then use gor from any directory.
