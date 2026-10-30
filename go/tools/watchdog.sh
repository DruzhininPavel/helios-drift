#!/bin/sh
# Overnight watchdog: keeps running the suite and the bot soak, and leaves a
# log behind so the morning starts with results rather than guesses.
#
#   tools/watchdog.sh [rounds]
set -u
cd "$(dirname "$0")/.." || exit 1
mkdir -p build
LOG=build/watchdog.log
ROUNDS=${1:-0}
n=0

say() { printf '%s %s\n' "$(date '+%H:%M:%S')" "$*" >>"$LOG"; }

say "watchdog start (rounds=${ROUNDS:-forever})"
while :; do
	n=$((n + 1))
	say "--- round $n ---"

	if ! go vet ./... >>"$LOG" 2>&1; then
		say "FAIL vet"
	fi
	if ! go test ./... >>"$LOG" 2>&1; then
		say "FAIL tests"
	fi

	# one play attempt per mode, different seed each round
	for mode in advanced classic; do
		seed=$((n * 101 + 7))
		out=$(go run . -bot -mode "$mode" -seed "$seed" -frames 18000 2>&1)
		say "$mode seed $seed $(printf '%s\n' "$out" | grep -E '^cash|^DIED|^STUCK' | tr '\n' ' ')"
		printf '%s\n' "$out" | grep -A6 "VIOLATIONS" | while read -r line; do
			say "  $line"
		done
	done

	if [ "$ROUNDS" -gt 0 ] && [ "$n" -ge "$ROUNDS" ]; then
		say "watchdog done after $n rounds"
		exit 0
	fi
	sleep 120
done
