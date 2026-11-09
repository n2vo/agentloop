### Agentloop

A minimal agent runtime loop in Go, built without an agent framework. It calls the Gemini Application Programming Interface (API) through the Google Gen AI Software Development Kit (SDK) for Go (`google.golang.org/genai`). The agent answers questions about an NVIDIA DGX Spark by calling a `gpu_status` tool, which runs `nvidia-smi` on the Spark over Secure Shell (SSH).

#### How the Loop Works

1. The question is sent to the model along with the declared tools.
2. If the response contains no tool calls, the answer is printed and the loop ends.
3. Otherwise, the model turn is appended to the conversation, each requested tool is executed, and the results are appended as a user turn.
4. The loop repeats, up to a cap of five steps.

#### Requirements

- Go 1.27.3 or later.
- A Gemini API key in the `GEMINI_API_KEY` environment variable.
- Key-based SSH access to the host that runs `nvidia-smi`.

#### Usage

```bash
go run .
go run . "Which GPU is installed?"
go run . -host <ssh-host> "Which GPU is installed?"
```

| Flag | Default | Description |
|---|---|---|
| `-host` | `spark-e2fe` | SSH host of the DGX Spark |

With no question argument, the agent asks for the installed GPU's utilization and temperature.

#### Step Log

Each iteration logs to standard error:

- **Model Call:** Step number, prompt token count, model latency, and number of tool calls requested.
- **Tool Call:** Tool name, arguments, result size in bytes, and tool latency.

#### Files

| File | Contents |
|---|---|
| `main.go` | Flags, the agent loop, and the step log |
| `tools.go` | The `gpu_status` tool declaration and implementation |
| `tools_test.go` | Test for `gpu_status` |

#### Testing

```bash
go test ./...
```

The test requires SSH access to the default host.
