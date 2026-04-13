# Call2Me Go SDK

The official Go SDK for [Call2Me](https://call2me.app) — the AI voice agent platform.

Build, deploy, and manage AI voice agents that handle real phone calls, extract data, and take automated actions.

[![Go Reference](https://pkg.go.dev/badge/github.com/call2me-app/go-sdk.svg)](https://pkg.go.dev/github.com/call2me-app/go-sdk)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

## Features

- **Voice Agents** — Create and manage AI agents with custom voices and personalities
- **Phone & Web Calls** — Handle inbound/outbound calls via SIP or browser
- **30+ AI Models** — GPT, Claude, Gemini, DeepSeek, Llama, and more
- **Knowledge Base** — RAG-powered answers from your documents
- **Campaigns** — Bulk outbound calling at scale
- **Scheduled Calls** — Book follow-up calls automatically
- **Post-Call Intelligence** — Extract structured data from every conversation
- **Zero Dependencies** — Only standard library

## Installation

```bash
go get github.com/call2me-app/go-sdk
```

Requires Go 1.21+

## Getting Your API Key

1. Sign up at [dashboard.call2me.app](https://dashboard.call2me.app/signup) — you get **$10 free credits**
2. Go to **API Keys** in the dashboard
3. Click **Create API Key** and copy it

## Quick Start

```go
package main

import (
    "fmt"
    "log"

    call2me "github.com/call2me-app/go-sdk"
)

func main() {
    client := call2me.New("sk_call2me_...")

    // List agents
    agents, err := client.Agents.List()
    if err != nil {
        log.Fatal(err)
    }

    for _, agent := range agents {
        fmt.Printf("%s — %s\n", agent["agent_id"], agent["agent_name"])
    }
}
```

## Usage Examples

### Create an Agent

```go
agent, err := client.Agents.Create(map[string]interface{}{
    "agent_name": "Customer Support",
    "voice_id":   "elevenlabs-selin",
    "language":   "tr-TR",
    "response_engine": map[string]interface{}{
        "type":          "call2me-llm",
        "system_prompt": "You are a friendly support agent.",
    },
})
```

### Get Call History

```go
calls, err := client.Calls.List()
if err != nil {
    log.Fatal(err)
}

for _, call := range calls {
    fmt.Printf("%s | %s | %s\n",
        call["call_id"], call["direction"], call["call_status"])
}
```

### Get Agent Details

```go
agent, err := client.Agents.Get("agent_abc123")
fmt.Println(agent["agent_name"])
```

## API Reference

| Service | Methods |
|---------|---------|
| `client.Agents` | `List()`, `Get(id)`, `Create(data)` |
| `client.Calls` | `List()` |

More services coming soon: Knowledge Base, Wallet, Campaigns, Schedules.

## Documentation

- **Full API Docs**: [call2me.app/docs](https://call2me.app/docs)
- **Guides**: [call2me.app/guides](https://call2me.app/guides)
- **Pricing**: [call2me.app/pricing](https://call2me.app/pricing)

## Support

- Email: [support@call2me.app](mailto:support@call2me.app)
- Website: [call2me.app](https://call2me.app)

## License

MIT — see [LICENSE](LICENSE) for details.
