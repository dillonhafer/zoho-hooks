#!/bin/sh

set -e

go build
exec ./zoho_webhooks
