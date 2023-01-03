# list-parse
Parse through a wordlist and get only the entries you want. Useful for limiting wordlists for brute-force attacks when password policy is known

## Quick Start
To get started, Download the latest release or compile from source

```
# listparse help
$ listparse -h
usage: list-parse [-h|--help] -w|--wordlist "<value>" -o|--output "<value>"
                  [-m|--min-length "<value>"] [-x|--max-length "<value>"]
                  [-p|--phrase "<value>"] [-s|--require-special-characters]
                  [-n|--require-number] [-r|--replace-letters]

                  Creates more customized wordlist

Arguments:

  -h  --help                        Print help information
  -w  --wordlist                    Wordlist to parse through; Full path to
                                    wordlist
  -o  --output                      File to output to
  -m  --min-length                  Minimum length of line. Default: 0
  -x  --max-length                  Maximum length of line. Default: 0
  -p  --phrase                      Phrase/word that is required to be in the
                                    line. Default:
  -s  --require-special-characters  Require special characters
  -n  --require-number              Require Number
  -r  --replace-letters             Replace Letters with special characters in
                                    phrase search
```
*Please note -r currently does not work*
## Examples
```
# Get all lines with only 8 characters and require special characters
$ listparse -w rockyou.txt -o outfile.txt -m 8 -x 8 -s
```

```
# Get all lines that contains numbers
$ listparse -w rockyou.txt -o outfile.txt -n
```

```
# Get all lines that contain the phrase "pass"
$ listparse -w rockyou.txt -o outfile.txt -p pass
```

```
# Get all lines between 8 and 12 characters that have a special character and number
$ listparse -w rockyou.txt -o outfile.txt -m 8 -x 12 -s -n
```
