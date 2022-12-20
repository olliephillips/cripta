# Cripta Messenger

## Peer to peer encrypted & signed messaging via MQTT pub/sub.

[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/olliephillips/cripta?style=flat-square)](https://go.dev/)
[![Go Report Card](https://goreportcard.com/badge/github.com/olliephillips/cripta?style=flat-square)](https://goreportcard.com/report/github.com/olliephillips/cripta)
[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/olliephillips/cripta/build.yml?branch=main&style=flat-square)](https://github.com/olliephillips/cripta/actions/workflows/build.yml)
[![GitHub Workflow Status (branch)](https://img.shields.io/github/actions/workflow/status/olliephillips/cripta/unit_test.yml?branch=main&label=tests&style=flat-square)](https://github.com/olliephillips/cripta/actions/workflows/unit_test.yml)

[![GitHub Release Date](https://img.shields.io/github/release-date/olliephillips/cripta?style=flat-square)](https://github.com/olliephillips/cripta/releases)
[![GitHub commits since latest release (by date)](https://img.shields.io/github/commits-since/olliephillips/cripta/latest?style=flat-square)](https://github.com/olliephillips/cripta/commits)
[![GitHub](https://img.shields.io/github/license/olliephillips/cripta?label=license&style=flat-square)](LICENSE)

With Cripta you can relay secure messages instantly using free (and mostly available) MQTT brokers. The broker is
configurable, to use a different broker edit, rebuild the app and distribute that.

## Features

- Asymmetric encryption using recipient public key
- Inbox management (list, read, delete & empty)
- Mail groups
- Message signing with private key (TODO)

## Releases

Download an executable for your OS [here](https://github.com/olliephillips/cripta/releases).

## Tests

No tests in this project. It started as an afternoon hack. Maybe there will never be any tests.

## License

MIT