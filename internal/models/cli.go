/*
 * Copyright (C) 2023 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 *
 *
 */

package models

type ResponderHeaderFields struct {
	RequestId string
	TraceId   string
}

var RespHeaderFields ResponderHeaderFields
