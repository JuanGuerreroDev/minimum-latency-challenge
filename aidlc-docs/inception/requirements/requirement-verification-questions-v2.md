# Preguntas de Verificación — Nuevo Intent: Stimulus Client Mode

Necesito clarificar algunos detalles para implementar correctamente este fix.

## Question 1
¿Cómo prefieres la interfaz del cliente de estímulo standalone?

A) Un nuevo binario separado (`cmd/stimulus/main.go`) que envía un único estímulo y muestra la latencia
B) Un flag/subcomando en el binario de benchmark existente (`cmd/benchmark --mode=single`) que envía un solo estímulo
C) Un nuevo binario separado que quede en modo interactivo (puedes enviar múltiples estímulos manualmente, ej. presionando Enter)
D) Other (please describe after [Answer]: tag below)

[Answer]: C

## Question 2
¿Debe el cliente standalone generar log de trazabilidad (.log) al igual que el benchmark, o solo mostrar el resultado en consola?

A) Solo mostrar latencia en consola (stdout)
B) Mostrar en consola Y además escribir en un archivo .log (configurable)
C) Other (please describe after [Answer]: tag below)

[Answer]: B

## Question 3
¿Las extensiones (Security Baseline y Property-Based Testing) mantienen la misma configuración del intent anterior?

A) Sí — mantener Security Baseline (pragmático) y PBT (parcial)
B) No — deshabilitarlas para este fix
C) Other (please describe after [Answer]: tag below)

[Answer]: A
