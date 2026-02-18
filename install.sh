#!/bin/bash
# install.sh — deploy orchestration commands and agent personas to ~/.claude/
#
# Usage: ./install.sh
#
# Copies slash commands and agent personas from this repo into ~/.claude/,
# making them globally available in all projects.
#
# Safe to run multiple times (idempotent). Backs up existing files before
# overwriting. Does NOT modify ~/.claude/CLAUDE.md or delete files not
# managed by this script.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TARGET_DIR="${HOME}/.claude"

COMMANDS_SRC="${SCRIPT_DIR}/.claude/commands"
AGENTS_SRC="${SCRIPT_DIR}/.manager/agents"

COMMANDS_DST="${TARGET_DIR}/commands"
AGENTS_DST="${TARGET_DIR}/agents"

# counters
installed=0
up_to_date=0
backed_up=0

# install_file copies a single file from src to dst.
# if dst exists and differs, it creates a .bak backup first.
# if dst exists and is identical, it skips the copy.
install_file() {
    local src="$1"
    local dst="$2"
    local name
    name="$(basename "$src")"

    if [[ -f "$dst" ]]; then
        if diff -q "$src" "$dst" >/dev/null 2>&1; then
            printf "  %-20s up to date\n" "$name"
            ((up_to_date++)) || true
            return
        fi
        cp "$dst" "${dst}.bak"
        printf "  %-20s updated (backup: %s.bak)\n" "$name" "$name"
        ((backed_up++)) || true
    else
        printf "  %-20s installed\n" "$name"
    fi

    cp "$src" "$dst"
    ((installed++)) || true
}

# validate source directories exist
if [[ ! -d "$COMMANDS_SRC" ]]; then
    echo "error: commands source not found: ${COMMANDS_SRC}" >&2
    exit 1
fi
if [[ ! -d "$AGENTS_SRC" ]]; then
    echo "error: agents source not found: ${AGENTS_SRC}" >&2
    exit 1
fi

echo "installing orchestration system to ${TARGET_DIR}/"
echo ""

# create target directories
mkdir -p "$COMMANDS_DST"
mkdir -p "$AGENTS_DST"

# install commands
echo "commands:"
for src in "${COMMANDS_SRC}"/*.md; do
    [[ -f "$src" ]] || continue
    install_file "$src" "${COMMANDS_DST}/$(basename "$src")"
done
echo ""

# install agents
echo "agents:"
for src in "${AGENTS_SRC}"/*.md; do
    [[ -f "$src" ]] || continue
    install_file "$src" "${AGENTS_DST}/$(basename "$src")"
done
echo ""

# summary
echo "done: ${installed} installed, ${up_to_date} up to date, ${backed_up} backed up"
