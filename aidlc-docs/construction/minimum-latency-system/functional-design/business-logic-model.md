# Business Logic Model — Minimum Latency System

## Visión General

La lógica de negocio del sistema es intencionalmente simple: recibir un estímulo y retornar una respuesta lo más rápido posible. La complejidad reside en **cómo** se ejecuta (Reactor Pattern, zero-allocation), no en **qué** se ejecuta.

---

## 1. Server — Reactor Handler Logic

### OnTraffic (HOT PATH)

Este es el método más crítico del sistema. Se ejecuta cada vez que hay datos disponibles en el event loop de gnet.

```
OnTraffic(conn):
  1. data = conn.ReadAll()           // Lee todos los bytes disponibles
  2. IF len(data) == 0:
       RETURN None                   // No hay datos, no hacer nada
  3. msgType = data[0]               // Primer byte = tipo de mensaje
  4. IF msgType == TypeStimulus (0x01):
       response[0] = TypeResponse (0x02)  // Buffer pre-allocated
       conn.Write(response[:1])           // Escribir 1 byte de respuesta
  5. ELSE:
       logger.Info("unknown message type", type=msgType)  // Ignorar, solo loggear
  6. RETURN None                     // Continuar el event loop
```

**Invariantes del HOT PATH**:
- Zero allocations: No se crea ningún objeto nuevo
- No branching complejo: Un solo `if/else`
- Buffer pre-allocated: `response` es un `[1]byte` fijo en el struct del handler
- I/O directo: `conn.Write()` escribe directo al socket vía event loop

### OnBoot

```
OnBoot(engine):
  1. Guardar referencia al engine
  2. logger.Info("server started", port=engine.Port)
  3. RETURN None
```

### OnOpen

```
OnOpen(conn):
  1. logger.Info("connection opened", remoteAddr=conn.RemoteAddr())
  2. RETURN (nil, None)  // Sin datos iniciales
```

### OnClose

```
OnClose(conn, err):
  1. IF err != nil:
       logger.Error("connection closed with error", err=err)
  2. ELSE:
       logger.Info("connection closed", remoteAddr=conn.RemoteAddr())
  3. RETURN None
```

---

## 2. Benchmark Client Logic

### runBenchmark

```
runBenchmark(conn, iterations):
  1. Pre-allocar:
     - stimulus = [1]byte{TypeStimulus}   // 1 byte fijo
     - readBuf = [1]byte{}               // Buffer de lectura fijo
     - durations = make([]Duration, 0, iterations)  // Slice pre-allocated
     - errors = 0
  
  2. FOR i = 0; i < iterations; i++:
     a. sendTime = time.Now()
     b. _, writeErr = conn.Write(stimulus[:])
     c. IF writeErr != nil:
          errors++
          CONTINUE                       // Skip and continue (Q3=A)
     d. _, readErr = conn.Read(readBuf[:])
     e. IF readErr != nil:
          errors++
          CONTINUE                       // Skip and continue (Q3=A)
     f. recvTime = time.Now()
     g. latency = recvTime - sendTime
     h. durations = append(durations, latency)
     i. benchmarkLogger.Record(sendTime, recvTime, latency)
  
  3. RETURN durations, errors
```

### main (Benchmark)

```
main():
  1. Parsear flags: host, port, iterations (default: 10000)
  2. addr = host:port
  3. conn = net.Dial("tcp", addr)
  4. defer conn.Close()
  5. benchLogger = NewBenchmarkLogger(iterations)
  6. durations, errors = runBenchmark(conn, iterations)
  7. stats = Calculate(durations)
  8. Print stats.Report() a stdout     // Reporte solo al final (Q4=A)
  9. benchLogger.FlushToFile("benchmark.log", stats)
  10. Print "Log written to benchmark.log"
  11. IF errors > 0:
        Print "WARNING: {errors} iterations failed"
```

---

## 3. Protocol Encode/Decode Logic

### Encode

```
Encode(msgType, buf):
  1. buf[0] = msgType
  2. RETURN 1                          // 1 byte escrito (payload vacío, Q2=B)
```

### Decode

```
Decode(data):
  1. IF len(data) < 1:
       RETURN 0, ErrEmptyMessage
  2. msgType = data[0]
  3. RETURN msgType, nil
```

---

## 4. Stats Calculation Logic

### Calculate

```
Calculate(durations):
  1. IF len(durations) == 0:
       RETURN Stats{} con todos los valores en 0
  2. sort(durations)                   // Sort para percentiles
  3. stats.Count = len(durations)
  4. stats.Min = durations[0]
  5. stats.Max = durations[len-1]
  6. stats.Total = sum(durations)
  7. stats.Avg = stats.Total / stats.Count
  8. stats.Median = durations[len/2]
  9. stats.P50 = durations[index(50)]  // index = (percentile/100) * len
  10. stats.P95 = durations[index(95)]
  11. stats.P99 = durations[index(99)]
  12. RETURN stats
```

---

## 5. BenchmarkLogger Logic

### Record (durante benchmark)

```
Record(sendTime, recvTime, latency):
  1. entry = buffer[currentIndex]      // Pre-allocated slot
  2. entry.SendTime = sendTime
  3. entry.RecvTime = recvTime
  4. entry.Latency = latency
  5. currentIndex++
```

### FlushToFile (post-benchmark)

```
FlushToFile(filename, stats):
  1. file = os.Create(filename)
  2. writer = bufio.NewWriter(file)
  3. Write header: "# Benchmark Trace Log"
  4. Write metadata: timestamp, iterations
  5. Write header row: "SendTime | RecvTime | Latency"
  6. FOR each entry in buffer[:currentIndex]:
       Write formatted: "{entry.SendTime} | {entry.RecvTime} | {entry.Latency}"
  7. Write separator
  8. Write stats summary section
  9. writer.Flush()
  10. file.Close()
```
