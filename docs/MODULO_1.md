# Módulo 1 — Fundamentos de Arquitectura de Software

## Descripción General

Este módulo explora el papel estratégico del arquitecto de software, que va más allá de lo técnico para alinear soluciones tecnológicas con objetivos de negocio. Se analizan las responsabilidades del arquitecto y las tres dimensiones clave que definen su importancia:

1. **Estructuración de la complejidad**
2. **Análisis temprano de atributos de calidad**
3. **Toma de decisiones técnicas informadas**

**Ejercicio práctico:** Desarrollar un sistema de respuesta ultrarrápida que procese estímulos con **latencia menor a 1 milisegundo**, aplicando los principios de diseño arquitectónico para optimizar el rendimiento.

---

## Estándar ISO/IEC 25010 — Calidad del Software

Estándar internacional que define un modelo de calidad para el software. Forma parte de la familia de normas **SQuaRE** (System and Software Quality Requirements and Evaluation).

### ¿Por qué existe?

Antes de este estándar, cada empresa o equipo definía "calidad" a su manera. La ISO/IEC 25010 **unifica ese criterio** y da un lenguaje común para evaluar, comparar y mejorar software de forma objetiva.

---

### Modelo 1 — Calidad del Producto

Define **8 características principales** con sus subcaracterísticas:

#### 1. 🔧 Adecuación Funcional
*¿El software hace lo que debe hacer?*

- **Completitud funcional** — ¿Cubre todas las funciones requeridas?
- **Corrección funcional** — ¿Los resultados son correctos y precisos?
- **Pertinencia funcional** — ¿Las funciones facilitan los objetivos del usuario?

> Ejemplo: Una calculadora que suma incorrectamente falla en corrección funcional.

#### 2. ⚡ Eficiencia en el Desempeño
*¿El software responde bien bajo ciertas condiciones de recursos?*

- **Comportamiento temporal** — Tiempos de respuesta y procesamiento
- **Utilización de recursos** — CPU, memoria, red que consume
- **Capacidad** — Límites máximos que puede manejar

> Ejemplo: Una app que tarda 10 segundos en cargar una página tiene mal comportamiento temporal.

#### 3. 🔄 Compatibilidad
*¿Puede coexistir e intercambiar información con otros sistemas?*

- **Coexistencia** — Funciona bien junto a otros software sin conflictos
- **Interoperabilidad** — Puede intercambiar datos con otros sistemas

> Ejemplo: Un sistema que no puede exportar datos a Excel o PDF tiene baja interoperabilidad.

#### 4. 👤 Usabilidad
*¿Qué tan fácil y agradable es usarlo?*

- **Reconocimiento de idoneidad** — El usuario entiende si el sistema sirve para lo que necesita
- **Capacidad de aprendizaje** — ¿Qué tan fácil es aprender a usarlo?
- **Operabilidad** — Facilidad para operar y controlar el sistema
- **Protección contra errores** — Evita que el usuario cometa errores
- **Estética de la interfaz** — Interfaz agradable visualmente
- **Accesibilidad** — Usable por personas con distintas capacidades

> Ejemplo: Un formulario que no indica qué campos son obligatorios tiene baja operabilidad.

#### 5. 🛡️ Fiabilidad
*¿El software funciona correctamente durante un tiempo determinado?*

- **Madurez** — ¿Con qué frecuencia falla en condiciones normales?
- **Disponibilidad** — ¿Está accesible cuando se necesita?
- **Tolerancia a fallos** — ¿Sigue funcionando ante fallos internos?
- **Recuperabilidad** — ¿Puede recuperar datos tras una falla?

> Ejemplo: Un sistema bancario que cae cada semana tiene baja madurez.

#### 6. 🔐 Seguridad
*¿Protege la información y los datos adecuadamente?*

- **Confidencialidad** — Solo acceden quienes deben acceder
- **Integridad** — Los datos no son alterados sin autorización
- **No repudio** — Se puede demostrar que una acción ocurrió
- **Responsabilidad** — Se puede rastrear quién hizo qué
- **Autenticidad** — Verifica que usuarios y recursos son legítimos

> Ejemplo: Una app que guarda contraseñas en texto plano falla en confidencialidad.

#### 7. 🔨 Mantenibilidad
*¿Qué tan fácil es modificar y corregir el software?*

- **Modularidad** — Está dividido en partes independientes
- **Reusabilidad** — Sus componentes pueden usarse en otros sistemas
- **Analizabilidad** — Es fácil identificar dónde está un problema
- **Modificabilidad** — Se puede cambiar sin introducir nuevos errores
- **Capacidad de prueba** — Se pueden crear pruebas fácilmente

> Ejemplo: Un código con miles de líneas en un solo archivo tiene baja modularidad.

#### 8. 🌐 Portabilidad
*¿Puede trasladarse a otros entornos?*

- **Adaptabilidad** — Puede adaptarse a distintos sistemas operativos o hardware
- **Instalabilidad** — Se puede instalar/desinstalar sin problemas
- **Reemplazabilidad** — Puede reemplazar a otro software con el mismo propósito

> Ejemplo: Un software que solo corre en Windows XP tiene baja adaptabilidad.

---

### Modelo 2 — Calidad en Uso

Evalúa la calidad desde la perspectiva del **usuario final** usando el producto en un contexto real:

| Característica | Pregunta clave |
|---|---|
| Efectividad | ¿El usuario logra sus objetivos? |
| Eficiencia | ¿Lo logra con el esfuerzo correcto? |
| Satisfacción | ¿Está conforme con la experiencia? |
| Libertad de riesgo | ¿El uso genera riesgos económicos, humanos o ambientales? |
| Cobertura del contexto | ¿Funciona bien en todos los contextos previstos? |

### Resumen Visual del Modelo

```
ISO/IEC 25010
│
├── CALIDAD DEL PRODUCTO
│   ├── 1. Adecuación Funcional
│   ├── 2. Eficiencia en el Desempeño
│   ├── 3. Compatibilidad
│   ├── 4. Usabilidad
│   ├── 5. Fiabilidad
│   ├── 6. Seguridad
│   ├── 7. Mantenibilidad
│   └── 8. Portabilidad
│
└── CALIDAD EN USO
    ├── Efectividad
    ├── Eficiencia
    ├── Satisfacción
    ├── Libertad de riesgo
    └── Cobertura del contexto
```

### ¿Cómo se usa en la práctica?

Se aplica principalmente en tres momentos del ciclo de vida del software:

1. **Al definir requisitos** — Se usan las características como checklist para asegurarse de que se están considerando todos los aspectos de calidad desde el inicio.
2. **Durante el desarrollo** — Guía las decisiones de diseño y arquitectura. Por ejemplo, pensar en mantenibilidad desde el comienzo implica diseñar código modular.
3. **En la evaluación/auditoría** — Permite medir objetivamente qué tan bueno es un producto antes de entregarlo o certificarlo.

---

## Leyes de la Arquitectura de Software

### Primera Ley
> **"Todo en la arquitectura de software es un intercambio (trade-off)"**

Cada elección técnica conlleva compensaciones. Una decisión verdaderamente arquitectónica equilibra atributos de calidad contrapuestos (como integridad vs. costo), en lugar de buscar soluciones unilaterales.

### Segunda Ley
> **"El por qué es más importante que el cómo"**

---

## Alcance de la Arquitectura de Software

La arquitectura consiste en la estructura combinada con:

- **Características de la arquitectura** ("-ilidades") — lo que el sistema debe soportar
- **Decisiones de arquitectura** — reglas para construir sistemas
- **Principios de diseño** — pautas para construir sistemas

Las responsabilidades de un arquitecto trascienden la implementación directa. Su valor radica en tomar decisiones estructurales que aborden requisitos críticos, alineando tecnología con necesidades de negocio.

> **Ley de Conway:** El diseño de sistemas refleja la estructura de comunicación de las organizaciones. Esto resalta la necesidad de alinear equipos e infraestructura.

---

## Importancia de la Arquitectura en la Práctica

1. **Definir estructuras esenciales** — Establece elementos principales, sus relaciones y propiedades. Permite razonar sobre el sistema y sus atributos de calidad.
2. **Facilita el análisis temprano** — Permite predecir y gestionar atributos como rendimiento, seguridad y escalabilidad desde etapas iniciales, reduciendo riesgos y costos.
3. **Base para decisiones técnicas** — Al encapsular decisiones fundamentales, ayuda a los equipos a evitar soluciones contradictorias o costosas.

### Conexión con los objetivos estratégicos

1. **Traduce metas abstractas en soluciones concretas** — Por ejemplo, optimizar la escalabilidad para soportar un crecimiento imprevisto en usuarios.
2. **Influye en la estructura organizacional** — Según la Ley de Conway, el diseño refleja la estructura de comunicación de las organizaciones.
3. **Facilita la reutilización y líneas de productos** — Un modelo arquitectónico bien diseñado puede servir como base de múltiples productos.

### Contribución al crecimiento organizacional

1. **Habilitar la innovación y la agilidad** — Sistemas modulares bien documentados permiten iterar rápidamente y adaptarse a cambios del mercado.
2. **Mejorar la comunicación entre stakeholders** — Un diseño claro y compartido facilita la colaboración entre equipos técnicos, de negocio y otros interesados.
3. **Fomentar competencias** — Una arquitectura sólida fortalece la organización cuando enfrenta retos.

---

## 13 Puntos Clave de la Arquitectura de Software

1. Una arquitectura puede inhibir o habilitar los atributos de calidad de un sistema.
2. Las decisiones arquitectónicas permiten gestionar el cambio a medida que evoluciona el sistema.
3. El análisis de una arquitectura permite una predicción temprana de la calidad del sistema.
4. Una arquitectura documentada mejora la comunicación entre las partes.
5. Las decisiones arquitectónicas son más tempranas, fundamentales y difíciles de cambiar.
6. La arquitectura define un conjunto de restricciones en la implementación.
7. La arquitectura dicta la estructura de una organización.
8. La arquitectura puede proporcionar la base para un desarrollo incremental.
9. La arquitectura es el artefacto clave que permite razonar sobre costos y programas.
10. Se puede crear una arquitectura como modelo transferible y reutilizable (línea de productos).
11. El desarrollo basado en arquitectura centra la atención en el ensamblaje de componentes, no solo en su creación.
12. Al restringir las alternativas de diseño, la arquitectura canaliza la creatividad de los desarrolladores, reduciendo la complejidad.
13. Una arquitectura puede ser la base para la formación de un nuevo miembro del equipo.

> Estos 13 puntos son la brújula para crear soluciones técnicas alineadas con el negocio, garantizando calidad, adaptabilidad y crecimiento.

---

## Restricciones y Decisiones Arquitectónicas

En la arquitectura de software, las decisiones clave definen la estructura y el comportamiento del sistema, mientras que las **restricciones** actúan como límites que guían el diseño y la implementación.

### Restricciones de Arquitectura
Reglas o limitaciones que dictan ciertos aspectos del diseño y la implementación del sistema. Comprender las restricciones del negocio, las limitaciones tecnológicas y los recursos disponibles es esencial para tomar decisiones arquitectónicas efectivas.

### Decisiones de Arquitectura
Elecciones fundamentales que definen la estructura y el comportamiento de un sistema de software — las bases sobre las cuales será construido y evolucionará.

---

## Características Arquitectónicas

Los atributos de calidad y las características arquitectónicas son fundamentales para construir sistemas de software no solo funcionales, sino también sostenibles, eficientes y adaptables.

### Operacionales
*Cómo el sistema funciona y se comporta en producción. ¿Cómo vive el sistema en el mundo real?*

> **Ejemplo:** Un banco necesita que su sistema esté disponible 24/7, que responda en menos de 2 segundos y que pueda escalar en fechas de pago masivo.

| # | Característica | Descripción |
|---|---|---|
| 1 | **Disponibilidad** | Tiempo que el sistema debe estar operativo (ej. 24/7 requiere medidas de recuperación ante fallas) |
| 2 | **Continuidad** | Capacidad de recuperación ante desastres |
| 3 | **Rendimiento** | Pruebas de estrés, análisis de picos, frecuencia de uso, capacidad requerida y tiempos de respuesta |
| 4 | **Recuperabilidad** | Requisitos de continuidad del negocio (tiempo máximo para restaurar tras desastre). Impacta estrategia de backups y hardware redundante |
| 5 | **Confiabilidad / Seguridad** | Evalúa si el sistema debe ser a prueba de fallos o es crítico (afecta vidas humanas o genera pérdidas financieras) |
| 6 | **Robustez** | Capacidad para manejar errores y condiciones límite: caídas de conexión, cortes de energía, fallos de hardware |
| 7 | **Escalabilidad** | Capacidad para mantener rendimiento y operatividad ante aumento de usuarios o solicitudes |

### Estructurales
*Cómo está organizado y construido el sistema internamente. ¿Cómo está hecho por dentro?*

> **Ejemplo:** Decidir si una aplicación será monolítica con capas (presentación → negocio → datos) o si se dividirá en microservicios independientes que se comunican por API.

| # | Característica | Descripción |
|---|---|---|
| 1 | **Configurabilidad** | Capacidad de los usuarios finales para modificar aspectos de la configuración del software mediante interfaces intuitivas |
| 2 | **Extensibilidad** | Grado de importancia para incorporar nuevas funcionalidades de manera modular |
| 3 | **Capacidad de instalación** | Facilidad para desplegar el sistema en todas las plataformas requeridas |
| 4 | **Reutilización** | Habilidad para aprovechar componentes comunes en múltiples productos |
| 5 | **Localización** | Soporte para múltiples idiomas, caracteres multibyte, unidades de medida y monedas locales |
| 6 | **Mantenibilidad** | Facilidad para implementar cambios y mejoras a lo largo del tiempo |
| 7 | **Portabilidad** | ¿Requiere ejecutarse en múltiples plataformas? (ej. Oracle y SAP DB) |
| 8 | **Capacidad de actualización** | Facilidad para migrar desde versiones anteriores a versiones más recientes en servidores y clientes |

### Transversales (Cross-Cutting)
*Aspectos que afectan a todo el sistema por igual, sin importar qué parte del código se está mirando. ¿Qué reglas aplican en todas partes?*

> **Ejemplo:** La autenticación de usuarios no pertenece solo a una capa o servicio; debe aplicarse en toda la aplicación de forma consistente.

| # | Característica | Descripción |
|---|---|---|
| 1 | **Accesibilidad** | Garantizar el acceso a todos los usuarios sin importar si sufren discapacidad |
| 2 | **Capacidad de archivado** | ¿Los datos deberán archivarse o eliminarse después de un tiempo? (ej. cuentas obsoletas) |
| 3 | **Autenticación** | Requisitos de seguridad para verificar la identidad de los usuarios |
| 4 | **Autorización** | Requisitos de seguridad para controlar el acceso a funciones específicas |
| 5 | **Aspectos legales** | Restricciones legislativas: protección de datos, GDPR, Sarbanes-Oxley, derechos de reserva |
| 6 | **Privacidad** | Capacidad de proteger transacciones incluso del personal interno (encriptación que impida acceso a administradores de BD) |
| 7 | **Seguridad** | ¿Encriptación en BD? / ¿Encriptación en comunicaciones internas? / ¿Protocolos de autenticación para acceso remoto? |
| 8 | **Capacidad de soporte** | Complejidad de logs necesarios / Herramientas de diagnóstico / Recursos para solución de incidencias |
| 9 | **Usabilidad** | Curva de aprendizaje / Diseño intuitivo centrado en el usuario / Requisitos ergonómicos |

---

## Cierre del Módulo

En este primer módulo hemos profundizado en los fundamentos de la arquitectura de software, explorando su papel como **puente entre las estrategias de negocio y las soluciones técnicas**.

Desde las responsabilidades del arquitecto hasta las implicaciones de las restricciones tecnológicas, organizacionales y humanas, hemos construido una comprensión integral de cómo las decisiones arquitectónicas impactan la escalabilidad, seguridad y mantenibilidad de los sistemas.

Al analizar atributos de calidad, características operacionales y estándares como ISO/IEC 25010, hemos establecido las bases para diseñar sistemas **robustos y alineados con las necesidades del mundo real**.

> **Recurso de video:** [Arquitectura de Software — Papel Estratégico](https://www.youtube.com/watch?v=m8UmSPmw3jU&t=96s)
