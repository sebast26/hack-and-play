# gdoc-writer

Program creates new Google Document with the content from STDIN.

## Usage

```text
NAME:
   gdoc-writer - creates Google Document with the content from stdin

USAGE:
   gdoc-writer [global options] command [command options] [arguments...]

DESCRIPTION:
   creates Google Document with the content from stdin

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --prefix value, -p value  prefix for the document title
   --help, -h                show help (default: false)
```

## How to build and run

Before running the program please make sure that you have `credentials.json` file in program working directory. You may
get the credentials file from your `Google Console`.

Program expects being run inside pipe. To run the program you execute:

```bash
cat file1 file2 | go run gdoc-writer.go
```
