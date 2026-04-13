# Call2Me Go SDK

The official Go SDK for [Call2Me](https://call2me.app) — the AI voice agent platform.

Build, deploy, and manage AI voice agents that handle real phone calls, extract data, and take automated actions.

[![Go Reference](https://pkg.go.dev/badge/github.com/call2me-app/go-sdk.svg)](https://pkg.go.dev/github.com/call2me-app/go-sdk)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

## Features

- **13 API Resources** — Agents, Calls, KB, Wallet, Campaigns, Schedules, Phone Numbers, SIP Trunks, API Keys, Users, Widgets, Voices, Payments
- **Zero Dependencies** — Only Go standard library
- **Simple API** — Consistent patterns across all resources

## Installation

```bash
go get github.com/call2me-app/go-sdk
```

Requires Go 1.21+

## Getting Your API Key

1. Sign up at [dashboard.call2me.app](https://dashboard.call2me.app/signup) — you get **$10 free credits**
2. Go to **API Keys** in the dashboard
3. Click **Create API Key** and copy your `sk_call2me_...` key

## Quick Start

```go
package main

import (
    "fmt"
    call2me "github.com/call2me-app/go-sdk"
)

func main() {
    client := call2me.New("sk_call2me_...")

    // List agents
    agents, _ := client.Agents.List()
    for _, a := range agents {
        fmt.Printf("%s — %s\n", a["agent_id"], a["agent_name"])
    }

    // Check balance
    balance, _ := client.Wallet.Balance()
    fmt.Printf("Balance: $%v\n", balance["balance_usd"])
}
```

## Full API Reference

### Agents
```go
client.Agents.List()
client.Agents.Get("agent_id")
client.Agents.Create(call2me.M{"agent_name": "My Agent", "voice_id": "elevenlabs-selin"})
client.Agents.Update("agent_id", call2me.M{"agent_name": "New Name"})
client.Agents.Delete("agent_id")
client.Agents.Duplicate("agent_id")
client.Agents.Stats("agent_id")
```

### Calls
```go
client.Calls.List()
client.Calls.Get("call_id")
client.Calls.End("call_id")
```

### Knowledge Base
```go
client.KnowledgeBase.List()
client.KnowledgeBase.Get("kb_id")
client.KnowledgeBase.Create(call2me.M{"name": "FAQ"})
client.KnowledgeBase.Delete("kb_id")
client.KnowledgeBase.Query("kb_id", "question")
```

### Campaigns
```go
client.Campaigns.List()
client.Campaigns.Get("campaign_id")
client.Campaigns.Create(call2me.M{"name": "Sale", "agent_id": "...", "from_number": "+908501234567"})
client.Campaigns.Start("campaign_id")
client.Campaigns.Pause("campaign_id")
client.Campaigns.Cancel("campaign_id")
```

### Scheduled Calls
```go
client.Schedules.List()
client.Schedules.Create(call2me.M{"agent_id": "...", "phone_number": "+905551234567", "scheduled_at": "2026-04-15T10:00:00+03:00"})
client.Schedules.Delete("schedule_id")
```

### Phone Numbers
```go
client.PhoneNumbers.List()
client.PhoneNumbers.Create(call2me.M{"phone_number": "+908501234567", "trunk_id": "..."})
client.PhoneNumbers.Delete("+908501234567")
client.PhoneNumbers.BindAgent("+908501234567", "agent_id")
```

### SIP Trunks
```go
client.SipTrunks.List()
client.SipTrunks.Create(call2me.M{"name": "My Trunk", "sip_server": "sip.provider.com", "sip_username": "user", "sip_password": "pass"})
client.SipTrunks.Delete("trunk_id")
client.SipTrunks.Test("trunk_id")
```

### Wallet & Billing
```go
client.Wallet.Balance()
client.Wallet.Transactions()
client.Wallet.Analytics()
```

### Payments
```go
client.Payments.Checkout(50.0, "USD")
client.Payments.History()
client.Payments.SavedCards()
```

### API Keys
```go
client.ApiKeys.List()
client.ApiKeys.Create(call2me.M{"name": "Production Key"})
client.ApiKeys.Revoke("key_id")
client.ApiKeys.Delete("key_id")
```

### Users & Branding
```go
client.Users.Me()
client.Users.Update(call2me.M{"full_name": "John Doe"})
client.Users.Stats()
client.Users.Branding()
```

### Widgets
```go
client.Widgets.List()
client.Widgets.Create(call2me.M{"agent_id": "...", "name": "Support"})
client.Widgets.Delete("widget_id")
client.Widgets.Chat("widget_id", "Hello")
```

### Voices
```go
client.Voices.List()
```

## Links

- **Website**: [call2me.app](https://call2me.app)
- **Dashboard**: [dashboard.call2me.app](https://dashboard.call2me.app)
- **API Docs**: [call2me.app/docs](https://call2me.app/docs)
- **Guides**: [call2me.app/guides](https://call2me.app/guides)
- **GitHub**: [github.com/call2me-app/go-sdk](https://github.com/call2me-app/go-sdk)
- **Support**: [support@call2me.app](mailto:support@call2me.app)

## License

MIT
