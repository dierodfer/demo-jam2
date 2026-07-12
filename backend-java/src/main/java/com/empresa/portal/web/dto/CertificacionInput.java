package com.empresa.portal.web.dto;

/** Cuerpo de POST/PUT de certificaciones (sin id ni empleadoId). */
public record CertificacionInput(
        String conocimiento,
        String empresaEmisora,
        String fecha) {
}
