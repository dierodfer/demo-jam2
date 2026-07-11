package com.empresa.portal.web.dto;

/** Cuerpo de PUT /api/me: campos editables del perfil (sin id ni username). */
public record EmpleadoUpdate(
        String nombre,
        String email,
        String telefono,
        String puesto,
        String departamento,
        String direccion,
        String foto) {
}
