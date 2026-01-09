CLI Tools for CLI AI Agents
===========================

Examples
--------

````console
$ # First, prepare the prompt.md file.
$ cat ./prompt.md
You process the TSV file and output the JSON result to the output location.

# Output Format
```json
{
  "category": string, /* category of the result */
  "reason": string /* reason of the result */
}
```

# Output Location
```
%%OUTPUT%%
```

# Category definition
```tsv
Category	Description
A Foo
B Bar
C Baz
```

# Input
```tsv
%%INPUT_TSV%%
```

$ # Next, collect input TSV files.
$ find ./input -name '*.tsv' -print0 > ./input_files

$ # Next, put a prompt generator script in the ./prompt directory.
$ cat ./prompt_generator
#!/bin/bash
set -euo pipefail
input_file="$1"
stdinsubst < ./prompt/template.md "%%OUTPUT%%" <(echo "$input_file" | sed -e "s|\./input/|./output/|" -e "s|\.tsv|.json|") "%%INPUT_TSV%%" "$input_file" >"./prompt/$(basename "$input_file").md"

$ # Next, generate prompt.md for each input TSV file.
$ stdinexec -0 <./input_files ./prompt_generator "{}"

$ # Next, collect prompt.md files.
$ find ./prompt -name '*.md' -print0 > ./prompt_files

$ # Then, process the TSV files in parallel using 3 processes by Claude Code or Gemini or Codex or so on.
$ stdinexec -0 <./prompt_files bash -c 'claude -dangerously-skip-permissions -p < "{}"'

$ # You can combine the above steps into a single command.
$ find ./input -name '*.tsv' -print0 \
    | stdinexec -0 bash -c 'stdinsubst < ./prompt/template.md "%%OUTPUT%%" <(echo "{}" | sed -e "s|\./input/|./output/|" -e "s|\.tsv|.json|") "%%INPUT_TSV%%" "{}" >"./prompt/$(basename "{}").md"' \
    | stdinexec -0 bash -c 'claude -dangerously-skip-permissions -p < "{}"'

$ # If Claude Code fails, you can resume by the following steps.
$ # 1. Collect processed files.
$ find ./output -name '*.json' -print0 > ./output_files

$ # 2. Convert output paths to input paths.
$ stdinexec -0 <./output_files bash -c 'printf "{}" | sed -e "s|^\./output/|./input/|" -e "s|\.json$|.tsv|"' > ./input_files.processed

$ # 3. Subtract processed files from input files.
$ stdinsub -0 <./input_files ./input_files.processed > ./input_files.unprocessed

$ # 4. Remove previous prompt.md files.
$ find ./prompt -name '*.md' -exec rm "{}" \;

$ # 5. Generate prompt.md for unprocessed files.
$ stdinexec -0 <./input_files.unprocessed ./prompt_generator "{}"

$ # 6. Collect prompt.md files.
$ find ./prompt -name '*.md' -print0 > ./prompt_files

$ # 7. Resume the process by the following command.
$ stdinexec -0 <./prompt_files bash -c 'claude -dangerously-skip-permissions -p < "{}"'

$ # You can also combine the steps 1-7 into a single command.
$ find ./input -name '*.tsv' -print0 \
    | stdinsub -0 <(find ./output -name '*.json' -print0 | stdinexec -0 bash -c 'printf "{}" | sed -e "s|^\./output/|./input/|" -e "s|\.json$|.tsv|"') \
    | stdinexec -0 bash -c 'stdinsubst < ./prompt/template.md "%%OUTPUT%%" <(echo "{}" | sed -e "s|\./input/|./output/|" -e "s|\.tsv$|.json|") "%%INPUT_TSV%%" "{}"' | claude -dangerously-skip-permissions -p'
````


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


### stdinsplit

```console
$ stdinsplit -h
Usage: stdinsplit [-0] (-l <line-count> | -n <total-count>) -o <out-dir> [-t <template>] < <input>

Split the input by the separator and write each part to a file in the output directory.
If <line-count> is specified, split the input into <line-count> lines.
If <total-count> is specified, split the input into <total-count> parts.

Options:
  -0	use null byte as the record separator
  -l int
    	number of lines per part
  -line-count int
    	number of lines per part
  -n int
    	number of parts
  -o string
    	output directory
  -out-dir string
    	output directory
  -t string
    	basename template (default: "%03d.txt")
  -template string
    	basename template (default: "%03d.txt")
  -total-count int
    	number of parts
  -v	print version and exit
  -version
    	print version and exit

Examples:
  $ # Split the input into 10 parts.
  $ echo "Hello\nWorld\n" | stdinsplit -o ./output -n 10 -t "part-%02d.txt"
  ./output/part-00.txt
  ./output/part-01.txt
  ...
  ./output/part-09.txt

  $ # Split the input into 10 lines.
  $ echo "Hello\nWorld\n" | stdinsplit -o ./output -l 1 -t "part-%02d.txt"
  ./output/part-00.txt
  ./output/part-01.txt

  $ # Use with stdinexec to process each part in parallel.
  $ echo "Hello\nWorld\n" | stdinsplit -0 -o ./output -l 1 | stdinexec -0 -p 2 bash -c 'claude -p < "{}"'
```

### stdinsub

```console
$ stdinsub -h
Usage: stdinsub <number1> <number2>

Subtract the second number from the first number.

Options:
  -0	use null byte as the record separator
  -v	print version and exit
  -version
    	print version and exit

Examples:
  $ cat ./minuend.txt
  line 1
  line 2
  line 3

  $ cat ./subtrahend1.txt
  line 2

  $ cat ./subtrahend2.txt
  line 2
  line 3

  $ stdinsub ./subtrahend1.txt ./subtrahend2.txt < ./minuend.txt
  line 1

  $ # It is useful to drop processed files from the input.
  $ find ./input -name '*.md' -print0 | stdinsub -0 <(find ./output -name '*.md' -print0 | sed -e 's|^\./input/|./output/|')
  ./input/file1.md
  ...

  $ # Process unprocessed ./input/*.md files in parallel using 3 processes by Claude Code.
  $ stdinsub -0 <(find ./input -name '*.md' -print0) <(find ./output -name '*.md' -print0 | sed -e 's|^\./input/|./output/|') | stdinexec -0 bash -c 'claude -p < "{}"'
```

License
-------

[MIT License](./LICENSE)
