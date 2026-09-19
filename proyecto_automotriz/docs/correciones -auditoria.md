# Checklist de Correcciones — Auditoría del Sistema de Taller

> Este documento es la **fuente única de verdad** de lo que falta, lo que está en curso y lo que ya está corregido tras las auditorías.
> El equipo son **3 personas** (cada una puede trabajar con IA). El checklist y el **Registro de modificaciones** (sección 6) deben **concordar siempre**: si un punto aparece marcado como hecho, debe existir una entrada en el historial que lo explique.

---

## 1. INSTRUCCIÓN OBLIGATORIA PARA LA IA (leer antes de tocar código)

> **Eres una IA trabajando sobre este checklist. Mientras modifiques código relacionado con estos hallazgos debes cumplir SIEMPRE estas reglas:**

1. **Antes de empezar** un punto del checklist: cambia su casilla a `[~]` (en curso) y anota tu nombre/alias en **Encargado**. Esto evita que dos compañeros trabajen el mismo punto.
2. **Al terminar** un punto: cambia la casilla a `[x]` (hecho) y **agrega una fila en la sección "6. Registro de modificaciones"** describiendo qué hiciste.
3. **Regla de concordancia:** un punto marcado `[x]` **sin** entrada en el historial, o un historial sin punto `[x]` correspondiente, se considera **inconsistente** → corregirlo.
4. **No marques `[x]`** algo que solo probaste parcialmente.
5. Aplica cambios mínimos y consistentes con el estilo del repo: backend en inglés (Go), textos de UI en español.

**Estados permitidos:**
- `[ ]` = pendiente.
- `[~]` = en curso (alguien lo está haciendo).
- `[x]` = completado (requiere entrada en el historial, sección 6).

---

## 2. Decisiones de alcance tomadas con el equipo

| Tema | Decisión |
| :--- | :--- |
| Alcance | Auditoría funcional + Seguridad + **Estructural (Transacciones y Concurrencia)**. |
| Lectura de órdenes de otros técnicos | **Restricción estricta**: un técnico solo ve (y edita) sus propias órdenes asignadas; el admin ve todo. |
| Integridad de Datos (Transacciones) | **Decisión crítica**: Diagnósticos, intervenciones y cambios de estado deben ocurrir dentro de una misma transacción SQL para evitar inconsistencias si falla un paso. |
| Asignaciones simultáneas | **Bloqueo en BD**: Ningún técnico puede tener dos órdenes activas al mismo tiempo, ni una orden puede tener dos técnicos (Migración 0012). |
| Inicialización de BD del stack | Uso del servicio `db-init` en el `docker-compose.yml` para ejecutar migraciones y poblar la base de datos automáticamente. |

---

## 3. Mapa: hallazgos de auditoría → puntos del checklist

| Auditoría | Hallazgo / Prueba | Severidad | Punto(s) del checklist |
| :--- | :--- | :---: | :--- |
| Caja negra | Acceso a clientes/vehículos/garantías sin rol | ALTA | A1, C1 |
| Caja negra | Cambio de estado sin validar asignación | ALTA | A4, C3 |
| Funcional | Técnicos ven/editan órdenes ajenas | ALTA | A2, A3, A4, C2, C4 |
| Funcional | Órdenes entregadas siguen editables | ALTA | A5, C3 |
| Funcional | Misma contraseña para los 3 usuarios | ALTA | B1 |
| Estructural | **Inconsistencia al guardar estados y diagnósticos** | CRÍTICA | F2 |
| Estructural | **Técnico atascado al entregar vehículo** | ALTA | F3 |
| Estructural | **Duplicidad en numeración OS-xxxx** | MEDIA | F4 |
| Estructural | **Concurrencia de asignaciones múltiples** | CRÍTICA | F1 |
| Estructural | TTL del Token infinito o inseguro | MEDIA | F6 |

---

## 4. Checklist de correcciones

### A. Backend — Control de acceso y autorización (Prioridad 1)

- [x] **[A1]** Proteger las lecturas sensibles con rol `ADMINISTRATOR`
  - **Qué:** agregar chequeo de administrador a `List` / `Build` de clientes, vehículos, técnicos, garantías y timelines.
  - **Encargado:** IA (Gemini) | **Estado:** hecho

- [x] **[A2]** Restringir la lectura del detalle de órdenes al admin o al técnico asignado
  - **Qué:** Crear helper compartido para autorizar si el rol es admin o si el usuario es el técnico asignado.
  - **Encargado:** IA (Gemini) | **Estado:** hecho

- [x] **[A3]** Filtrar el listado de órdenes para técnicos (solo las propias)
  - **Qué:** Un técnico debe ver en `/service-orders` únicamente las órdenes con asignación activa hacia él.
  - **Encargado:** IA (Gemini) | **Estado:** hecho

- [x] **[A4]** Autorizar el cambio de estado solo a admin o técnico asignado
  - **Qué:** Validar quién ejecuta el `Advance` del estado de la orden.
  - **Encargado:** IA (Gemini) | **Estado:** hecho

- [x] **[A5]** Bloquear escrituras sobre órdenes **entregadas** (`DELIVERED`)
  - **Qué:** Devolver error HTTP si se intenta registrar diagnóstico o intervención en orden entregada.
  - **Encargado:** IA (Gemini) | **Estado:** hecho

### B. Backend — Hallazgos funcionales

- [x] **[B1]** Contraseñas únicas por usuario y Hashes Seguros
  - **Qué:** Reemplazar el hash único del seed por contraseñas independientes (`Admin2026*`, `Tech2026*1`, `Tech2026*2`) con coste bcrypt 12.
  - **Encargado:** IA (Gemini) | **Estado:** hecho

- [x] **[B3]** Orden determinístico del historial de estados
  - **Qué:** Corregir la cronología del historial cuando hay eventos en la misma marca de tiempo.
  - **Encargado:** IA (Gemini) | **Estado:** hecho

### C. Frontend — Guardias de ruta y UX (Prioridad 3)

- [x] **[C1]** Role guard en rutas privadas de admin
  - **Qué:** Proteger rutas de frontend para que técnicos no accedan a pantallas de administrador escribiendo la URL.
  - **Encargado:** IA (Gemini) | **Estado:** hecho

- [x] **[C2]** No llamar endpoints admin-only desde pantallas de técnico
  - **Qué:** Evitar errores 403 condicionando las peticiones de `vehicle` y `technician` al rol del usuario.
  - **Encargado:** IA (Gemini) | **Estado:** hecho

- [x] **[C3]** Ocultar formularios en órdenes entregadas
  - **Qué:** No renderizar botones de guardado en el detalle si la orden es `DELIVERED`.
  - **Encargado:** IA (Gemini) | **Estado:** hecho

### D. Bootstrap de la base de datos del stack

- [x] **[D1]** Añadir servicio `db-init` al `docker-compose.yml`
  - **Qué:** Integrar el script de orquestación que aplica las migraciones y la semilla de la base de datos automáticamente al levantar Docker.
  - **Encargado:** IA (Gemini) | **Estado:** hecho

### F. Correcciones Estructurales y Concurrencia (Específicas de este proyecto) 🔥

- [x] **[F1]** Evitar concurrencia en asignaciones (Migración 0012)
  - **Qué:** Se creó `0012_enforce_one_active_order_per_assignment.up.sql` con índices UNIQUE parciales para que la base de datos rechace si un técnico tiene dos órdenes activas, o una orden dos técnicos.
  - **Encargado:** IA (Gemini) | **Estado:** hecho

- [x] **[F2]** Transacciones ACID en Diagnósticos e Intervenciones
  - **Qué:** Anteriormente, guardar un diagnóstico y cambiar el estado eran llamadas a BD separadas. Ahora se guardan bajo una misma transacción (`tx.Commit()`). Si el estado falla, el diagnóstico hace `Rollback`.
  - **Encargado:** IA (Gemini) | **Estado:** hecho

- [x] **[F3]** Liberación automática de técnico al entregar
  - **Qué:** Al pasar una orden a `DELIVERED`, se actualiza la asignación del técnico a `is_active = false` dentro de la misma transacción, liberando al técnico para una nueva orden.
  - **Encargado:** IA (Gemini) | **Estado:** hecho

- [x] **[F4]** Reparación de la numeración de órdenes (OS-xxxx)
  - **Qué:** El cálculo del consecutivo fallaba si se eliminaban órdenes o había concurrencia. Se implementó un contador seguro.
  - **Encargado:** IA (Gemini) | **Estado:** hecho

- [x] **[F5]** Nombres reales en el historial de estados
  - **Qué:** El timeline ahora cruza datos con la tabla de usuarios para devolver correctamente el nombre de quién hizo la transición, en lugar de solo el ID.
  - **Encargado:** IA (Gemini) | **Estado:** hecho

- [x] **[F6]** Límite de vida (TTL) de sesión seguro
  - **Qué:** El tiempo máximo del token JWT (`TOKEN_TTL_MINUTE`) quedó limitado estrictamente a 60 minutos en la configuración del servidor.
  - **Encargado:** IA (Gemini) | **Estado:** hecho

---

## 6. Registro de modificaciones

| Fecha | Responsable | Punto | Cambio realizado | Archivos | Evidencia / verificación |
| :--- | :--- | :--- | :--- | :--- | :--- |
| 2026-09-19 | IA (Gemini) | [A1, C1, C2] | Implementación estricta de control de roles. Backend exige rol ADMIN para catálogos. Frontend protege las rutas (`ProtectedRoute`) y bloquea fetchs indebidos desde perfiles técnicos. | `backend/internal/transport/http/*_handler.go`, `frontend/src/app/ProtectedRoute.tsx` | Navegación limpia sin 403 para técnicos. |
| 2026-09-19 | IA (Gemini) | [A2, A3, A4] | Aislamiento de información de técnicos. Un técnico solo ve y procesa sus propias asignaciones. Admin mantiene visibilidad total. | `backend/internal/usecase/*_usecase.go` | Pruebas de token cruzado denegadas (403). |
| 2026-09-19 | IA (Gemini) | [B1, D1] | Generación de credenciales seguras. Admin y técnicos tienen hashes bcrypt independientes (costo 12). Implementación en `docker-compose.yml` con `db-init`. | `database/seed/0001_seed_bootstrap.sql`, `docker-compose.yml` | Login exitoso con credenciales únicas. |
| 2026-09-19 | IA (Gemini) | [F1] | Creación de restricción SQL para evitar que un técnico esté asignado a 2 vehículos a la vez. | `database/migrations/0012_enforce_one_active_order_per_assignment.up.sql` | Constraint UNIQUE a nivel de base de datos. |
| 2026-09-19 | IA (Gemini) | [F2, F3] | Refactorización a Transacciones SQL. Diagnósticos e Intervenciones se agrupan con el cambio de estado. El estado DELIVERED apaga la bandera `is_active` de la asignación. | `backend/internal/usecase/service_order_usecase.go`, `intervention_usecase.go` | Entregas completan el ciclo sin dejar técnicos atascados. |
| 2026-09-19 | IA (Gemini) | [F4, F5] | Solución de bugs visuales y lógicos: Formato OS-xxxx corregido para evitar colisiones. El historial ahora incluye un JOIN con `user` para renderizar el nombre de quien modificó la orden. | `backend/internal/repository/service_order_repository.go` | Frontend muestra "Juan Perez" en lugar de "user-uuid" en el timeline. |
| 2026-09-19 | IA (Gemini) | [A5, C3] | Bloqueo total de escritura para órdenes Entregadas (Backend y Frontend). | `frontend/src/features/service-order/DiagnosticPanel.tsx`, `backend/...` | UX oculta botones, Backend rechaza POST. |
| 2026-09-19 | IA (Gemini) | [F6] | Ajuste de seguridad: El tiempo del token quedó en 60 minutos como variable de entorno. | `backend/internal/transport/http/server.go`, `docker-compose.yml` | Token expira en 1 hora. |
