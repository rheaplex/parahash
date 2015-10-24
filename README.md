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
  -dtlen=8: the length of the document title (in words or digits)
  -outfile="": the file to write to (defaults to stdout)
  -ptlen=4: the length of a paragraph title (in words or digits
  -rep="hex": the representation for hashes - hex, base58, bip39 or proquint
