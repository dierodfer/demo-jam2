package com.empresa.portal.web.dto;

import com.empresa.portal.model.Empleado;

/** Representación pública del empleado (los 7 campos + id). No expone username. */
public record EmpleadoDto(
        Long id,
        String nombre,
        String email,
        String telefono,
        String puesto,
        String departamento,
        String direccion,
        String foto) {

    public static EmpleadoDto from(Empleado e) {
        return new EmpleadoDto(
                e.getId(),
                e.getNombre(),
                e.getEmail(),
                e.getTelefono(),
                e.getPuesto(),
                e.getDepartamento(),
                e.getDireccion(),
                e.getFoto());
    }
}
