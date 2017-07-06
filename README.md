# WarpWallet
This is an implementation of [WarpWallet](https://keybase.io/warp) in Go.  WarpWallet is a brain wallet generator (for Bitcoin), originally written by [Max Krohn](https://github.com/maxtaco) and [Chris Coyne](https://github.com/malgorithms).  You can use it to turn passphrases into Bitcoin wallets, so your money is as safe as your memory :)

Except for referencing some hash functions (described below), this package is entirely self contained.  It is released under the BSD 2-clause license, and includes some BSD-style code from [ThePiachu](https://github.com/thepiachu).

This program has been tested under Linux and Windows.

## Install
To install, you'll need to run these from the command line:

```
go get github.com/sour-is/koblitz/kelliptic
go get github.com/btcsuite/btcutil/base58
go get golang.org/x/crypto/ripemd160
go get golang.org/x/crypto/pbkdf2
go get golang.org/x/crypto/scrypt
```

Finally, you should be able to do a:

```
go get github.com/nearwood/warpwallet/warpwallet
```
If you add your `$GOPATH/bin` to your path, you should now be able to run `warpwallet`.

If the above instructions don't work, download this repo and run `go build` in the `warpwallet` directory.  This will create a `warpwallet` executable that you can then run.

## Test
To run the test suite, just run `go test` inside the `warpwallet` directory.

## Love
If you found this useful, please send me some love at `1GGCFrshLz46tdas9ZtKqX59n5UAFzR6sD` :)
