# NFR Requirements Plan — minimum-latency-system

## Plan Checklist

- [ ] Refinar requerimientos de rendimiento con métricas concretas
- [ ] Definir decisiones de tech stack pendientes (Go runtime tuning)
- [ ] Establecer criterios de aceptación de rendimiento
- [ ] Definir estrategia de testing (PBT framework selection — PBT-09)
- [ ] Generar artifacts:
  - [ ] nfr-requirements.md
  - [ ] tech-stack-decisions.md
- [ ] Validar completitud

---

## Preguntas de NFR Requirements

Por favor responde las siguientes preguntas llenando la etiqueta `[Answer]:` con la letra de tu elección.

### Performance Requirements

## Question 1

¿Cuál es tu criterio de éxito principal para la latencia?

A) p99 < 1ms — el 99% de las peticiones deben responder en menos de 1ms
B) p95 < 1ms — el 95% de las peticiones deben responder en menos de 1ms, permitiendo outliers
C) Average < 1ms — el promedio debe ser menor a 1ms, aunque haya picos
D) p50 < 100µs — la mediana debe estar en el rango de microsegundos
X) Other (please describe after [Answer]: tag below)

[Answer]: A

## Question 2

¿Quieres aplicar tuning del Go runtime para optimizar latencia?

A) Sí, agresivo — GOGC=off (desactivar GC), GOMAXPROCS=1 (single core), CPU affinity
B) Sí, moderado — GOGC=off, GOMAXPROCS=default, sin CPU affinity
C) Mínimo — solo GOMAXPROCS=1 para reducir context switching
D) Default — dejar Go runtime con configuración por defecto para comparar
X) Other (please describe after [Answer]: tag below)

[Answer]: D, pero documentemos las otras opciones para el futuro.

### Testing Requirements

## Question 3

¿Cuál versión mínima de Go quieres soportar?

A) Go 1.22+ — última versión estable, mejores optimizaciones
B) Go 1.21+ — versión anterior estable, ampliamente disponible
C) La última disponible — siempre usar la versión más reciente
X) Other (please describe after [Answer]: tag below)

[Answer]: A
