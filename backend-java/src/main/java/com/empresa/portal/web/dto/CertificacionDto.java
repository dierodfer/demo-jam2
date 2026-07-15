package com.empresa.portal.web.dto;

import com.empresa.portal.model.Certificacion;

/** Representación pública de una certificación (sin empleadoId). */
public record CertificacionDto(
        Long id,
        String conocimiento,
        String empresaEmisora,
        String fecha) {

    public static CertificacionDto from(Certificacion c) {
        return new CertificacionDto(c.getId(), c.getConocimiento(), c.getEmpresaEmisora(), c.getFecha());
    }
}
