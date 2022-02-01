# basic-math-helper

## DESCRIPTION

This is sample program that generates simple mathematics exercises ready to be printed 
and give as a test to young kid learning math.

## Usage

To use program you can check what options do you have:

```bash
./basic-math-helper
```

This will print usage options

```
NAME:
   creates basic math exercises - A new cli application

USAGE:
   besic-math-helper [global options] command [command options] [arguments...]

COMMANDS:
   subtract_result_always_9  generates exercises with subtract operations that will always result in 9
   help, h                   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --num_exercises value, -n value  number of exercises to generate (default: 30)
   --help, -h                       show help (default: false)
```

To generate simple exercises about `subtract_result_always_9` subject you can invoke following command:

```bash
./basic-math-helper subtract_result_always_9
```
You can customise how many exercises will be generated with help of `--num_exercises` flag, ie:

```bash
./basic-math-helper -n 50 subtract_result_always_9
```

The default number of exercises is 30.

## Features

Program offers following options for generating exercises:

- `subtract_result_always_9` - generates exercises with subtract operations that will always result in 9
- `random_subtract` - generates random set of subtract exercises without any scheme