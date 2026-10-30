#!/bin/sh
# Run the bot over a spread of seeds and both rule sets, and leave a summary
# behind. This is the overnight job: every line is one full play attempt.
#
#   tools/soak.sh [frames] [seeds]
set -u
cd "$(dirname "$0")/.." || exit 1
FRAMES=${1:-36000}
SEEDS=${2:-6}
OUT=build/soak
mkdir -p "$OUT"
LOG=$OUT/summary.txt
: >"$LOG"

go build -o "$OUT/helios-bot" . || exit 1

say() { printf '%s\n' "$*" | tee -a "$LOG"; }

say "soak: $FRAMES frames ($((FRAMES / 60)) s) per run, $SEEDS seeds"
for mode in advanced classic; do
	for i in $(seq 1 "$SEEDS"); do
		seed=$((i * 37 + 3))
		rep="$OUT/$mode-$seed.json"
		out=$("$OUT/helios-bot" -bot -mode "$mode" -seed "$seed" -frames "$FRAMES" -report "$rep" 2>&1)
		line=$(printf '%s\n' "$out" | grep -E "^cash" | head -1)
		death=$(printf '%s\n' "$out" | grep -E "^DIED" | head -1)
		viol=$(printf '%s\n' "$out" | grep -c "frame .*: ")
		say "$(printf '%-9s seed %-4s %s %s' "$mode" "$seed" "$line" "$death")"
		if [ "$viol" -gt 0 ]; then
			say "    violations: $viol"
			printf '%s\n' "$out" | grep -A6 "VIOLATIONS" | tee -a "$LOG"
		fi
	done
done
say "soak finished"
