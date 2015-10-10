// hashparas
// Read a file of paragraphs and write them out with crypto hash titles
// Copyright 2015 Rob Myers <rob@robmyers.org>
// License: GNU GPLv3 or, at your option, any later version.

package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"encoding/hex"
	"github.com/tv42/base58"
	"github.com/tyler-smith/go-bip39"
	"hash"
	"io"
	"io/ioutil"
	"math/big"
	"os"
	"regexp"
	"strings"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

// Split the text into a list of paragraphs
// (a paragraph is text followed by \n\n)

func textParas(str string) []string {
	paras := strings.Split(str, "\n\n")
	nonempty := make([]string, 0)
    for _, para := range paras {
		stripped := strings.TrimSpace(para)
        if len(stripped) > 0 {
            nonempty = append(nonempty, stripped)
        }
    }
	return nonempty
}

// We don't want to treat formatting as significant for hashing
// So collapse whitespace to single spaces, and strip Markdown formatting

func stripParasMarkdown(paras []string) []string {
	whitespace, _ := regexp.Compile(`\s+`)
	//FIXME This doesn't handle escaped _ or *
	emphasis, _ := regexp.Compile(`[\*_]+`)
	url, _ := regexp.Compile(`\[([^\]]+)\]\([^)]+\)`)
	stripped := make([]string, 0)
	for _, para := range paras {
		para = emphasis.ReplaceAllString(para, "")
		para = whitespace.ReplaceAllString(para, " ")
		para = url.ReplaceAllString(para, "$1")
		stripped = append(stripped, para)
	}
	return stripped
}

// Allow the user to specify different hash representations

func hashString(hasher hash.Hash, hashRep string) string {
	str := ""
	if hashRep == "base58" {
		bigNum := new(big.Int)
		bigNum.SetBytes(hasher.Sum(nil)[:])
		rep := base58.EncodeBig(nil, bigNum)
		str = string(rep)
	} else if hashRep == "bip39" {
		str, _ = bip39.NewMnemonic(hasher.Sum(nil)[:])
	} else {
		// Convert [32]byte to []byte
		str = hex.EncodeToString(hasher.Sum(nil)[:])
	}
	return str
}

// Create a hash string representation for each string in the array

func parasHashes(paras []string) []hash.Hash {
	hashes := make([]hash.Hash, 0)
	for _, para := range paras {
		hasher := sha256.New()
		hasher.Write([]byte(para))
		hashes = append(hashes, hasher)
	}
	return hashes
}

// Combine all the hash strings into a single hash.
// A block hash, as it were

func hashOfHashes(hashes []hash.Hash, hashRep string) string {
	hasher := sha256.New()
	for _, h := range hashes {
		hasher.Write(h.Sum(nil))
	}
	return hashString(hasher, hashRep)
}

// Truncate the hash representation

func truncateHashRep (hashstr string, hashSliceLength int) string {
	if hashSliceLength > 0 {
		if (strings.Contains(hashstr, " ")) {
			words := strings.Split(hashstr, " " )
			hashstr = strings.Join(words[:hashSliceLength], " ")
		} else {
			hashstr = hashstr[:hashSliceLength]
		}
	}
	return hashstr
}

// Format the paragraphs to file with their hashes as their (h2) titles

func printParas (outfile io.Writer, paras []string, hashes []hash.Hash,
	hashRep string, hashSliceLength int) {
	for i, hash := range hashes {
		hashstr := truncateHashRep(hashString(hash, hashRep), hashSliceLength)
		fmt.Fprintln(outfile)
		fmt.Fprintln(outfile, "## " + hashstr)
		fmt.Fprintln(outfile)
		fmt.Fprintln(outfile, paras[i])
	}
}

func printDocumentTitle (outfile io.Writer, hashstr string,
	hashSliceLength int) {
	hashstr = truncateHashRep(hashstr, hashSliceLength)
	fmt.Fprintln(outfile, "# " + hashstr)
}

// Get configuration from command line flags
// Set up input and output files (or use standard input/output)
// Read the input, generate hashes for paragraph titles, write results to file.

func main() {
	infile := os.Stdin
	outfile := os.Stdout

	flag.CommandLine.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: parahash [OPTION] [FILE]")
		flag.PrintDefaults()
	}

	hashRep := flag.String("rep", "hex",
		"the representation for hashes - hex, base58 or bip39")
	paraTitleLength := flag.Int("ptlen", 4, "the length of a paragraph title")
	docTitleLength := flag.Int("dtlen", 8, "the length of the document title")
	outfilepath := flag.String("outfile", "",
		"the file to write to (defaults to stdout)")
	flag.Parse()

	if flag.NArg() > 1 {
		flag.CommandLine.Usage()
		os.Exit(1)
	} else if flag.NArg() == 1 {
		f, err := os.Open(flag.Arg(0))
		check(err)
		infile = f
	}

	if *outfilepath != "" {
		// Fail if exists
		f, err := os.OpenFile(*outfilepath, os.O_RDWR|os.O_CREATE|os.O_EXCL,
			0666)
		check(err)
		outfile = f
	}

	bytes, err := ioutil.ReadAll(infile)
	check(err)
	text := string(bytes)
	paras := textParas(text)
	stripped := stripParasMarkdown(paras)
	hashes := parasHashes(stripped)
	doctitlehash := hashOfHashes(hashes, *hashRep)
	printDocumentTitle(outfile, doctitlehash, *docTitleLength)
	printParas(outfile, paras, hashes, *hashRep, *paraTitleLength)
}
