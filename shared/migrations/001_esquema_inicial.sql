-- 001 — Esquema inicial: empleado y sus certificaciones (PostgreSQL).
--
-- Las tablas usan IF NOT EXISTS para adoptar bases de datos creadas antes de
-- existir el sistema de migraciones. Una vez aplicada, esta migración no debe
-- editarse: los cambios de esquema van en una migración nueva (002, 003, ...).

CREATE TABLE IF NOT EXISTS empleado (
    id           INTEGER PRIMARY KEY,
    username     TEXT,
    nombre       TEXT,
    email        TEXT,
    telefono     TEXT,
    puesto       TEXT,
    departamento TEXT,
    direccion    TEXT,
    foto         TEXT
);

CREATE TABLE IF NOT EXISTS certificacion (
    id              INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    empleado_id     INTEGER,
    conocimiento    TEXT,
    empresa_emisora TEXT,
    fecha           TEXT
);
