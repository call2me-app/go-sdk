# Call2Me Go SDK

The official Go SDK for [Call2Me](https://call2me.app) — the AI voice agent platform.

[![Go Reference](https://pkg.go.dev/badge/github.com/call2me-app/go-sdk.svg)](https://pkg.go.dev/github.com/call2me-app/go-sdk)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

## Features

- **13 API Resources** — Agents, Calls, KB, Wallet, Campaigns, Schedules, Phone Numbers, SIP Trunks, API Keys, Users, Widgets, Voices, Payments
- **Zero Dependencies** — Only Go standard library
- **Type-safe** — Strongly typed service methods

## Installation

```bash
go get github.com/call2me-app/go-sdk
```

## Getting Your API Key

1. Sign up at [dashboard.call2me.app](https://dashboard.call2me.app/signup) — **$10 free credits**
2. Go to **API Keys** → **Create API Key**

## Quick Start

```go
package main

import (
    "fmt"
    call2me "github.com/call2me-app/go-sdk"
)

func main() {
    client := call2me.New("sk_call2me_...")
    agents, _ := client.Agents.List()
    for _, a := range agents {
        fmt.Printf("%s — %s\n", a["agent_id"], a["agent_name"])
    }
}
```

## API Reference

| Service | Methods |
|---------|---------|
| `Agents` | `List` `Get` `Create` `Update` `Delete` `Duplicate` `Stats` |
| `Calls` | `List` `Get` `End` |
| `KnowledgeBase` | `List` `Get` `Create` `Delete` `Query` |
| `Wallet` | `Balance` `Transactions` `Analytics` |
| `Campaigns` | `List` `Get` `Create` `Start` `Pause` `Cancel` |
| `Schedules` | `List` `Create` `Delete` |
| `PhoneNumbers` | `List` `Create` `Delete` `BindAgent` |
| `SipTrunks` | `List` `Create` `Delete` `Test` |
| `ApiKeys` | `List` `Create` `Revoke` `Delete` |
| `Users` | `Me` `Update` `Stats` `Branding` |
| `Widgets` | `List` `Create` `Delete` `Chat` |
| `Voices` | `List` |
| `Payments` | `Checkout` `History` `SavedCards` |

## Links

- [call2me.app](https://call2me.app) · [API Docs](https://call2me.app/docs) · [Guides](https://call2me.app/guides) · [GitHub](https://github.com/call2me-app/go-sdk) · [support@call2me.app](mailto:support@call2me.app)

## License

MIT
