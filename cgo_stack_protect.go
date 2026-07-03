//go:build linux && cgo

/*
 * Copyright (C) 2026 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package main

/*
// _ensure_stack_protect guarantees the __stack_chk_guard symbol is present in
// the final binary, satisfying BinSkim rule BA3003 (EnableStackProtector).
//
// Background: the Go toolchain passes -gno-record-gcc-switches to GCC, which
// prevents BinSkim from reading compiler flags from the binary's DWARF info.
// BinSkim therefore falls back to symbol-based detection of __stack_chk_guard.
// __attribute__((stack_protect)) forces GCC to emit a stack canary for this
// function regardless of -fstack-protector-strong heuristics.
//
// Reference: https://github.com/microsoft/go/issues/2240
#ifdef __GNUC__
__attribute__((used, noinline, stack_protect))
static void _ensure_stack_protect(void) {
    volatile char buf[8];
    buf[0] = 0;
}
#endif
*/
import "C"
