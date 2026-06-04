# Dependencias de Componentes — Minimum Latency Challenge

## Matriz de Dependencias

| Componente | Depende de | Tipo de Dependencia |
|---|---|---|
| `cmd/server` | `internal/reactor`, `internal/logger` | Import directo |
| `cmd/benchmark` | `internal/protocol`, `internal/logger`, `internal/stats` | Import directo |
| `internal/reactor` | `internal/protocol`, `internal/logger` | Import directo |
| `internal/protocol` | Ninguno (standalone) | — |
| `internal/logger` | Ninguno (standalone) | — |
| `internal/stats` | Ninguno (standalone) | — |

## Dependencias Externas

| Componente | Dependencia Externa | Versión |
|---|---|---|
| `cmd/server` | `github.com/panjf2000/gnet/v2` | Latest stable |
| `internal/reactor` | `github.com/panjf2000/gnet/v2` | Latest stable |
| `cmd/benchmark` | Standard library only (`net`, `time`, `os`) | Go 1.22+ |

---

## Grafo de Dependencias

```mermaid
flowchart TD
    subgraph Executables["Ejecutables (cmd/)"]
        SERVER["cmd/server<br/><b>main</b>"]
        BENCH["cmd/benchmark<br/><b>main</b>"]
    end
    
    subgraph Internal["Paquetes Internos (internal/)"]
        REACTOR["internal/reactor<br/><b>ReactorHandler</b>"]
        PROTO["internal/protocol<br/><b>Encode/Decode</b>"]
        LOG["internal/logger<br/><b>Logger + BenchmarkLogger</b>"]
        STATS["internal/stats<br/><b>Calculate</b>"]
    end
    
    subgraph External["Dependencias Externas"]
        GNET["gnet/v2<br/><b>Event Loop</b>"]
        STDLIB["Go stdlib<br/><b>net, time, os</b>"]
    end
    
    SERVER --> REACTOR
    SERVER --> LOG
    SERVER --> GNET
    
    BENCH --> PROTO
    BENCH --> LOG
    BENCH --> STATS
    BENCH --> STDLIB
    
    REACTOR --> PROTO
    REACTOR --> LOG
    REACTOR --> GNET
    
    style SERVER fill:#2196F3,stroke:#0D47A1,stroke-width:2px,color:#fff
    style BENCH fill:#2196F3,stroke:#0D47A1,stroke-width:2px,color:#fff
    style REACTOR fill:#FF9800,stroke:#E65100,stroke-width:2px,color:#000
    style PROTO fill:#4CAF50,stroke:#1B5E20,stroke-width:2px,color:#fff
    style LOG fill:#4CAF50,stroke:#1B5E20,stroke-width:2px,color:#fff
    style STATS fill:#4CAF50,stroke:#1B5E20,stroke-width:2px,color:#fff
    style GNET fill:#9C27B0,stroke:#4A148C,stroke-width:2px,color:#fff
    style STDLIB fill:#607D8B,stroke:#263238,stroke-width:2px,color:#fff
    style Executables fill:#BBDEFB,stroke:#1565C0,stroke-width:2px,color:#000
    style Internal fill:#C8E6C9,stroke:#2E7D32,stroke-width:2px,color:#000
    style External fill:#F3E5F5,stroke:#7B1FA2,stroke-width:2px,color:#000
```

---

## Patrones de Comunicación

### Server ↔ Benchmark (TCP)
- **Protocolo**: Binario ultra-minimal (1 byte tipo + payload)
- **Patrón**: Request-Response síncrono sobre conexión persistente
- **Dirección**: Bidireccional (client envía stimulus, server responde response)
- **Formato**: `0x01 + "ping"` → `0x02 + "pong"`

### Internal Components (Go Imports)
- **Patrón**: Llamadas de función directas (in-process)
- **Acoplamiento**: Loose coupling via interfaces y packages separados
- **Nota**: `internal/protocol`, `internal/logger` y `internal/stats` son standalone sin dependencias entre sí

---

## Estructura del Proyecto

```
minimum-latency-challenge/
+-- cmd/
|   +-- server/
|   |   +-- main.go              # Server entry point
|   +-- benchmark/
|       +-- main.go              # Benchmark client entry point
+-- internal/
|   +-- reactor/
|   |   +-- handler.go           # gnet EventHandler implementation
|   +-- protocol/
|   |   +-- protocol.go          # Binary encode/decode
|   |   +-- protocol_test.go     # PBT roundtrip tests
|   +-- logger/
|   |   +-- logger.go            # Structured logger
|   |   +-- benchmark_logger.go  # Buffered benchmark logger
|   +-- stats/
|       +-- stats.go             # Latency statistics
|       +-- stats_test.go        # Stats calculation tests
+-- docs/
|   +-- MODULO_1.md              # Context document
|   +-- Estructura-del-conocimiento.md
|   +-- system-documentation.md  # Generated: architecture docs
|   +-- results-report.md        # Generated: benchmark results
+-- go.mod
+-- go.sum
+-- README.md
```
