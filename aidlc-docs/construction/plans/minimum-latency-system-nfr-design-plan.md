# NFR Design Plan — minimum-latency-system

## Plan Checklist

- [x] Diseñar patrones de rendimiento (zero-allocation, pre-allocated buffers, connection reuse)
- [x] Diseñar patrón del Reactor (gnet event loop configuration)
- [x] Diseñar patrones de resiliencia (error handling, graceful shutdown)
- [x] Documentar componentes lógicos del sistema
- [x] Generar artifacts:
  - [x] nfr-design-patterns.md
  - [x] logical-components.md
- [x] Validar completitud

---

## Preguntas de NFR Design

### Performance Patterns

## Question 1

¿Quieres que el servidor use un solo event loop (single-loop) o múltiples event loops (multi-loop) de gnet?

A) Single event loop — 1 loop, mínimo overhead, ideal para benchmark local con 1 conexión
B) Multi event loop — NumCPU loops, mayor throughput para múltiples conexiones concurrentes
X) Other (please describe after [Answer]: tag below)

[Answer]: A

### Resilience Patterns

## Question 2

¿Quieres que el servidor implemente graceful shutdown (atrapar SIGINT/SIGTERM y cerrar conexiones limpiamente)?

A) Sí — atrapar señales y hacer shutdown limpio del event loop
B) No — simplemente terminar el proceso (suficiente para un benchmark local)
X) Other (please describe after [Answer]: tag below)

[Answer]: A
