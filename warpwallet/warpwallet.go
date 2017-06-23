//ng.org/x/crypto/scrypt Copyright (c) 2013 Charles M. Ellison III
// All rights reserved.

// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met: 

// 1. Redistributions of source code must retain the above copyright notice, this
//    list of conditions and the following disclaimer. 
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution. 

// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
// ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
// (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package main

import "fmt"
import "bufio"
import "io"
import "os"
import "crypto/sha256"
import "strings"

import "github.com/btcsuite/btcutil/base58"
import "github.com/sour-is/koblitz/kelliptic"

import "golang.org/x/crypto/scrypt"
import "golang.org/x/crypto/pbkdf2"
import "golang.org/x/crypto/ripemd160"

func main() {
	passphrase := ""
	salt := ""

	if len(os.Args) == 3 {
		passphrase = os.Args[1]
		salt = os.Args[2]
	} else {
		passphrase, salt = getInputFromUser(os.Stdout, os.Stdin)
	}

	private, address := generate(passphrase, salt)
	fmt.Printf("Private key: %s\n", private)
	fmt.Printf("Public address: %s\n", address)
}

// to put a nice wrapper around dealing with input/output streams (helps with testing)
func getInputFromUser(ostream io.Writer, istream io.Reader) (string, string) {
	scanner := bufio.NewScanner(istream)

	fmt.Fprintf(ostream, "Please enter your passphrase: ")
	passphrase := readln(scanner);
	if strings.HasSuffix(passphrase, "\r") {
		fmt.Fprintf(ostream, "CAUTION!  Your passphrase ends with a \\r character.  Proceeding anyway...\n");
	}
	
	fmt.Fprintf(ostream, "Please enter your salt: ")
	salt := readln(scanner);
	if strings.HasSuffix(salt, "\r") {
		fmt.Fprintf(ostream, "CAUTION!  Your salt ends with a \\r character.  Proceeding anyway...\n");
	}

	return passphrase, salt
}

// reads a single line from the scanner.  does not return separator
func readln(scanner *bufio.Scanner) (string) {
	
	if scanner.Scan() {
	    return scanner.Text();
	}

	if err := scanner.Err(); err != nil {
	    panic(fmt.Sprintf("Trouble reading input: %v", err))
	}
	panic("Trouble reading input.")
}

// actually turns a passphrase + salt into a public + private
func generate(passphrase string, salt string) (string, string) {
	secret := secret([]byte(passphrase), []byte(salt))

	private := getPrivate(secret)
	public := getPublic(secret)
	address := getPublicAddress(public)

	return private, address
}

func getPrivate(pri []byte) string {
	bytes := []byte{0x80}
	bytes = append(bytes, pri...)
	
	sh := ShaTwice(bytes)
	checksum := make([]byte, 4)
	copy(checksum, sh[:4])
	
	bytes = append(bytes, checksum...)
	privWif := string(base58.Encode(bytes))
	return privWif
}

func getPublic(priv_key []byte) []byte {
	x, y := kelliptic.S256().ScalarBaseMult(priv_key)
	xbytes := x.Bytes()
	ybytes := y.Bytes()

	ret := make([]byte, 65)
	ret[0] = 4
	copy(ret[1 + 32 - len(xbytes):33], xbytes)
	copy(ret[33 + 32 - len(ybytes):65], ybytes)

	return ret
}

func getPublicAddress(pub []byte) string {
	if len(pub) != 65 {
		fmt.Printf("expected 65 long pub key")
	}
	// steps based on https://en.bitcoin.it/wiki/Technical_background_of_version_1_Bitcoin_addresses
	step3 := ShaRipemd(pub)
	step4 := append([]byte{0}, step3[:]...)
	step6 := ShaTwice(step4)
	checksum := step6[0:4]
	bbta := append(step4, checksum...)
	pubAddy := string(base58.Encode(bbta))

	return pubAddy
}

func secret(passphrase []byte, salt []byte) []byte {
	s1 := s1(passphrase, salt)
	s2 := s2(passphrase, salt)
	pri := xorBytes(s1, s2)
	return pri
}

func s1(passphrase []byte, salt []byte) []byte {
	salt = append(salt, 1)
	passphrase = append(passphrase, 1)
	s1, err := scrypt.Key(passphrase, salt, 262144, 8, 1, 32)
	if err != nil {
		panic(fmt.Sprintf("err: %v", err))
	}
	return s1
}

func s2(passphrase []byte, salt []byte) []byte {
	salt = append(salt, 2)
	passphrase = append(passphrase, 2)
	s2 := pbkdf2.Key(passphrase, salt, 65536, 32, sha256.New)
	return s2
}



func xorBytes(a []byte, b []byte) []byte {
	if len(a) != len(b) {
		panic("lengths not the same")
	}
	out := make([]byte, len(a))
	for i, x := range a {
		out[i] = x ^ b[i]
	}

	return out
}

func ShaTwice(a []byte) []byte {
	hasher := sha256.New()
	hasher.Write(a)
	a = hasher.Sum(nil)
	hasher.Reset()
	hasher.Write(a)
	return hasher.Sum(nil)
}

func ShaRipemd(a []byte) []byte {
	shaHasher := sha256.New()
	shaHasher.Write(a)
	ripemdHasher := ripemd160.New()
	ripemdHasher.Write(shaHasher.Sum(nil))
	return ripemdHasher.Sum(nil)
}
