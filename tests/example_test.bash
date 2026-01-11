#!/bin/bash
set -euo pipefail

BASE_DIR="$(cd "$(dirname "$0")/.." && pwd)"

has() {
    local command="$1"
    which "$command" >/dev/null 2>&1
}

prepare() {
    has go || throw "go is not installed"
    (cd "$BASE_DIR"
        go build -o ./tests/bin_fake/stdinsub ./tools/stdinsub/main.go
        go build -o ./tests/bin_fake/stdinsubst ./tools/stdinsubst/main.go
        git clean -fdx ./tests
        mkdir -p ./tests/prompt ./tests/output
    )
}

clean() {
    git clean -xdf .
    mkdir -p ./prompt ./output
}

main() {
    has parallel || throw "parallel is not installed"

    prepare

    export PATH="$BASE_DIR/tests/bin_fake:$PATH"

    # You can query multiple prompts to the Agentic AI by the following steps.
    (cd "$BASE_DIR/tests"
        set -x
        echo "====== TEST 1 ======"

        # 1. Prepare the prompt_template.md file.
        cat ./prompt_template.md

        # 2. Collect input TSV files.
        find ./input -name '*.tsv' -print0 > ./input_files

        # 3. Put a prompt generator script.
        cat ./prompt_generator

        # 4. Generate prompt.md for each input TSV file and collect them.
        parallel -0 <./input_files ./prompt_generator "{}" >./prompt_files

        # 5. Process the TSV files in parallel using 3 processes by Agentic AI.
        parallel -j3 -0 <./prompt_files 'claude -dangerously-skip-permissions -p < "{}"'

        set +x
        clean
    )

    # You can combine the above steps into fewer commands.
    (cd "$BASE_DIR/tests"
        set -x
        echo "====== TEST 2 ======"

        find ./input -name '*.tsv' -print0 \
            | parallel -0 ./prompt_generator "{}" \
            | parallel -j3 -0 'claude -dangerously-skip-permissions -p < "{}"'

        set +x
        clean
    )

    # If Agentic AI failed, you can resume by the following steps:
    (cd "$BASE_DIR/tests"
        set -x
        echo "====== TEST 3 ======"

        # 1. Collect input TSV files.
        find ./input -name '*.tsv' -print0 > ./input_files

        # 2. Collect processed files.
        find ./output -name '*.json' -print0 > ./output_files

        # 3. Convert output paths to input paths.
        parallel -0 <./output_files bash -c 'printf "{}" | sed -e "s|^\./output/|./input/|" -e "s|\.json$|.tsv|"' > ./input_files.processed

        # 4. Subtract processed files from input files.
        stdinsub -0 <./input_files ./input_files.processed >./input_files.unprocessed

        # 5. Remove previous prompt.md files.
        find ./prompt -name '*.md' -exec rm "{}" \;

        # 6. Re-generate prompt.md for unprocessed files.
        parallel -0 <./input_files.unprocessed ./prompt_generator "{}" >./prompt_files

        # 7. Resume the process by the following command.
        parallel -j3 -0 <./prompt_files 'claude -dangerously-skip-permissions -p < "{}"'

        set +x
        clean
    )

    (cd "$BASE_DIR/tests"
        set -x
        echo "====== TEST 4 ======"

        # You can combine the above steps into a single command.
        find ./input -name '*.tsv' -print0 \
            | stdinsub -0 <(find ./output -name '*.json' -print0 | sed -z -e "s|^\./output/|./input/|" -e "s|\.json$|.tsv|") \
            | parallel -0 ./prompt_generator "{}" \
            | parallel -j3 -0 'claude -dangerously-skip-permissions -p < "{}"'

        set +x
        clean
    )
}

main "$@"