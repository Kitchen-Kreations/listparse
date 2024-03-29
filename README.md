# listparse
<img alt="GitHub release (latest by date)" src="https://img.shields.io/github/v/release/Kitchen-Kreations/listparse"> <img alt="GitHub" src="https://img.shields.io/github/license/Kitchen-Kreations/listparse"> <img alt="GitHub all releases" src="https://img.shields.io/github/downloads/Kitchen-Kreations/listparse/total">


<img src="https://github.com/Kitchen-Kreations/listparse/blob/main/img/Listparse.PNG?raw=true" data-canonical-src="https://gyazo.com/eb5c5741b6a9a16c692170a41a49c858.png" width="250" height="250" />
Parse through a wordlist and get only the entries you want. Useful for limiting wordlists for brute-force & dictionary attacks when password policy is known

## Quick Start
To get started, Download the latest release or compile from source.

```
# listparse help
$ ./listparse -h          
[sub]Command required
usage: listparse <Command> [-h|--help]

                 Creates more customized wordlist

Commands:

  gui     Launch listparse's gui
  parser  User listparse as a command line tool

Arguments:

  -h  --help  Print help information
```
```
$ listparse gui -h
usage: listparse gui [-h|--help]

                 Launch listparse's gui

Arguments:

  -h  --help  Print help information
```
```
$ listparse parser -h
usage: listparse parser [-h|--help] -w|--wordlist "<value>" -o|--output
                 "<value>" [-m|--min-length "<value>"] [-x|--max-length
                 "<value>"] [-p|--phrase "<value>"]
                 [-s|--require-special-characters] [-n|--require-number]
                 [-v|--verbose]

                 User listparse as a command line tool

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
  -v  --verbose                     Verbose mode

```
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
