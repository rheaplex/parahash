Parahash
========

Parahash reads a text file containing a limited subset of Markdown, prints
each paragraph with its cryptographic has as its tile, and prints the hash of
those hashes as the resulting document's title.

Recognised Markdown formatting and extraneous whitespace does not affect the
hashes.

Usage
-----

Usage: parahash [OPTION] [FILE]
  -dtlen int
    	the length of the document title (default 8)
  -outfile string
    	the file to write to (defaults to stdout)
  -ptlen int
    	the length of a paragraph title (default 4)
  -rep string
    	the representation for hashes - hex, base58 or bip39 (default "hex")
