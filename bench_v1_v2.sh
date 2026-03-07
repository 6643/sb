#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")" && pwd)"
cd "$ROOT_DIR"

GO_CMD="${GO_CMD:-go}"
BUN_CMD="${BUN_CMD:-bun}"
GOCACHE_DIR="${GOCACHE:-/tmp/go-build}"

SB_SCHEMA="${SB_SCHEMA:-aaa.sb}"
SB_GO_DIR="${SB_GO_DIR:-./go}"
SB_TS_DIR="${SB_TS_DIR:-./ts}"
SB_GO_TAG="${SB_GO_TAG:-bson,json}"
SB_GO_BENCHTIME="${SB_GO_BENCHTIME:-100ms}"
SB_BENCH_ITERS="${SB_BENCH_ITERS:-1000}"
SB_BENCH_WARMUP="${SB_BENCH_WARMUP:-100}"
SB_BENCH_LIST_SIZE="${SB_BENCH_LIST_SIZE:-32}"
SB_SKIP_GEN="${SB_SKIP_GEN:-0}"

if ! command -v "$GO_CMD" >/dev/null 2>&1; then
    echo "missing go command: $GO_CMD" >&2
    exit 1
fi

if ! command -v "$BUN_CMD" >/dev/null 2>&1; then
    echo "missing bun command: $BUN_CMD" >&2
    exit 1
fi

if [[ "$SB_SKIP_GEN" != "1" ]]; then
    echo "==> generate v1/v2 code"
    GOCACHE="$GOCACHE_DIR" "$GO_CMD" run . -go "$SB_GO_DIR" -ts "$SB_TS_DIR" -tag "$SB_GO_TAG" -go-v2 -ts-v2 "$SB_SCHEMA"
fi

echo
echo "==> go v1 benchmark"
GOCACHE="$GOCACHE_DIR" "$GO_CMD" test ./go/sb -run '^$' -bench 'Sim' -benchtime "$SB_GO_BENCHTIME"

echo
echo "==> go v2 benchmark"
GOCACHE="$GOCACHE_DIR" "$GO_CMD" test ./go/sbv2 -run '^$' -bench 'V2' -benchtime "$SB_GO_BENCHTIME"

echo
echo "==> ts v1 benchmark"
SB_BENCH_ITERS="$SB_BENCH_ITERS" \
SB_BENCH_WARMUP="$SB_BENCH_WARMUP" \
SB_BENCH_LIST_SIZE="$SB_BENCH_LIST_SIZE" \
"$BUN_CMD" ts/sb/runtime_bench.ts

echo
echo "==> ts v2 benchmark"
SB_BENCH_ITERS="$SB_BENCH_ITERS" \
SB_BENCH_WARMUP="$SB_BENCH_WARMUP" \
SB_BENCH_LIST_SIZE="$SB_BENCH_LIST_SIZE" \
"$BUN_CMD" ts/sbv2/runtime_bench.ts
