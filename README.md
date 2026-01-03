CLI Tools for CLI AI Agents
===========================

Examples
--------


Usage
-----

### stdinexec

```console
$ stdinexec -h
Usage: stdinexec [-0] [-p <parallel>] <command> [<args>...]

Execute a command for each line of the input, similar to "find -exec". "{}" in arguments is replaced with the line of the stdin.

Options:
  -0	use null byte as the record separator
  -p int
    	number of parallel executions
  -parallel int
    	number of parallel executions
  -v	print version and exit
  -version
    	print version and exit

Examples:
  $ # Process ./input/*.md files in parallel using 3 processes by Claude Code.
  $ find ./input -name '*.md' -print0 | stdinexec -0 -p 3 bash -c 'claude -p < "{}"'
```

### stdinsubst

```console
$ stdinsubst -h
Usage: stdinsubst <before-string> <after-file-path> [<before-string> <after-file-path> ...] < <template>

<before-string> and <after-file-path> are the strings to replace and the file path to read the replacement from.

Options:
  -v	print version and exit
  -version
    	print version and exit

Examples:
  $ # Substitute the "%INPUT_FILE%" in the template with the contents of the "input-1.txt" file.
  $ stdinsubst "%INPUT_FILE%" ./input-1.txt < ./prompt.md

  $ # Process ./input/*.md files in parallel using 3 processes by Claude Code.
  $ find ./input -name '*.md' -print0 | stdinexec -0 bash -c 'stdinsubst "%INPUT_FILE%" "{}" < ./prompt.md | claude -p'

  $ # Substitute the "%GREETING%" and "%NOUN%" in the template with the contents of the "hello.txt" and "world.txt" files.
  $ echo '%GREETING%, %NOUN%!' | stdinsubst '%GREETING%' ./hello.txt '%NOUN%' ./world.txt

  $ # Also, you can substitute the "%GREETING%" and "%NOUN%" in the template with the process substitution.
  $ echo '%GREETING%, %NOUN%!' | stdinsubst '%GREETING%' <(printf Hello) '%NOUN%' <(printf World)
  Hello, World!
```


### splitsec

```console
$ splitsec -h
Usage: splitsec -o <output_directory> < <markdown>

Split a markdown file by sections into files based on the number of seconds in each section.

Options:
  -o string
    	output directory
  -out-dir string
    	output directory
  -t string
    	basename template
  -tmpl string
    	basename template
  -v	print version and exit
  -version
    	print version and exit
```