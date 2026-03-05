#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TEST_DIR="$ROOT_DIR/test"
WORK_DIR="$TEST_DIR/.work"

need_cmd() {
    local name="$1"
    if ! command -v "$name" >/dev/null 2>&1; then
        echo "missing command: $name" >&2
        return 1
    fi
}

need_cmd go
need_cmd bun

guard_no_invalid_nilable_defaults() {
    local go_sb_dir="$1"
    local bad_lines
    bad_lines="$(
        awk '
        BEGIN { in_func = 0; nilable_ret = 0; has_bad = 0 }
        /^func / {
            in_func = 1
            nilable_ret = 0
            if ($0 ~ /\(result (\*|\[\]).*, errCode RpcErrCode\)/) {
                nilable_ret = 1
            }
        }
        in_func && nilable_ret && /return 0, RpcRespErr/ {
            printf "%s:%d:%s\n", FILENAME, FNR, $0
            has_bad = 1
        }
        in_func && /^}/ {
            in_func = 0
            nilable_ret = 0
        }
        END {
            if (has_bad == 1) {
                exit 1
            }
        }
        ' "$go_sb_dir"/api.*.go 2>/dev/null || true
    )"

    if [ -z "$bad_lines" ]; then
        return 0
    fi
    echo "[guard] invalid nilable default return detected:"
    echo "$bad_lines"
    return 1
}

rm -rf "$WORK_DIR"
mkdir -p "$WORK_DIR"

run_core_case() {
    local case_name="robust_core"
    local schema="$TEST_DIR/${case_name}.sb"
    local case_dir="$WORK_DIR/${case_name}"
    local go_dir="$case_dir/go"
    local ts_dir="$case_dir/ts"

    mkdir -p "$go_dir" "$ts_dir"

    echo "[core] generate from $schema"
    (cd "$ROOT_DIR" && go run . -go "$go_dir" -ts "$ts_dir" "$schema")
    guard_no_invalid_nilable_defaults "$go_dir/sb"

    cat > "$go_dir/go.mod" <<'EOF'
module sbtest

go 1.22
EOF

    cp "$TEST_DIR/harness/robust_core_go_test.go.txt" "$go_dir/sb/robust_core_integration_test.go"

    echo "[core] go test"
    (cd "$go_dir" && go test ./... -count=1)

    echo "[core] bun test"
    (cd "$ts_dir/sb" && bun test rpc_smoke.test.ts)
}

run_smoke_case() {
    local case_name="$1"
    local schema="$TEST_DIR/${case_name}.sb"
    local case_dir="$WORK_DIR/${case_name}"
    local go_dir="$case_dir/go"
    local ts_dir="$case_dir/ts"

    mkdir -p "$go_dir" "$ts_dir"

    echo "[smoke] generate from $schema"
    (cd "$ROOT_DIR" && go run . -go "$go_dir" -ts "$ts_dir" "$schema")
    guard_no_invalid_nilable_defaults "$go_dir/sb"

    cat > "$go_dir/go.mod" <<'EOF'
module sbtest

go 1.22
EOF

    if [ "$case_name" = "robust_size_limits" ]; then
        cp "$TEST_DIR/harness/robust_size_limits_go_test.go.txt" "$go_dir/sb/robust_size_limits_integration_test.go"

        echo "[smoke] robust_size_limits go test"
        (cd "$go_dir" && go test ./... -count=1)

        echo "[smoke] robust_size_limits bun test"
        (cd "$ts_dir/sb" && bun test rpc_smoke.test.ts)
        return
    fi

    if [ "$case_name" = "robust_optional_matrix" ]; then
        cp "$TEST_DIR/harness/robust_optional_matrix_go_test.go.txt" "$go_dir/sb/robust_optional_matrix_integration_test.go"

        echo "[smoke] robust_optional_matrix go test"
        (cd "$go_dir" && go test ./... -count=1)

        echo "[smoke] robust_optional_matrix bun test"
        (cd "$ts_dir/sb" && bun test rpc_smoke.test.ts)
        return
    fi

    echo "[smoke] unsupported case: $case_name" >&2
    return 1
}

run_core_case
run_smoke_case robust_size_limits
run_smoke_case robust_optional_matrix

echo "all go+bun tests passed"
