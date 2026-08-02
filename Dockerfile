# Copyright 2021 The terraform-docs Authors.
# Copyright 2026 Northwood Labs, LLC <license@northwood-labs.com>.
#
# Licensed under the MIT license (the "License"); you may not
# use this file except in compliance with the License.
#
# You may obtain a copy of the License at the LICENSE file in
# the root directory of this source tree.

#-------------------------------------------------------------------------------
FROM golang:1.26.4-alpine@sha256:3ad57304ad93bbec8548a0437ad9e06a455660655d9af011d58b993f6f615648 AS builder

# hadolint ignore=DL3018
RUN apk add --update --no-cache make

WORKDIR /go/src/terraform-docs

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
# RUN make build

#-------------------------------------------------------------------------------
FROM alpine:3.23.4@sha256:5b10f432ef3da1b8d4c7eb6c487f2f5a8f096bc91145e68878dd4a5019afde11

# Mitigate CVE-2023-5363
RUN apk add --no-cache --upgrade "openssl>=3.1.4-r1"

# COPY --from=builder /go/src/terraform-docs/bin/linux-*/terraform-docs /usr/local/bin/

# ENTRYPOINT ["terraform-docs"]
