# NFR Design Plan — minimum-latency-system

## Plan Checklist

- [ ] Diseñar patrones de rendimiento (zero-allocation, pre-allocated buffers, connection reuse)
- [ ] Diseñar patrón del Reactor (gnet event loop configuration)
- [ ] Diseñar patrones de resiliencia (error handling, graceful shutdown)
- [ ] Documentar componentes lógicos del sistema
- [ ] Generar artifacts:
  - [ ] nfr-design-patterns.md
  - [ ] logical-components.md
- [ ] Validar completitud

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
