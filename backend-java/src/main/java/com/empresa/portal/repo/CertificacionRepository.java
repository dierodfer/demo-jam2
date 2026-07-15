package com.empresa.portal.repo;

import com.empresa.portal.model.Certificacion;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.List;

public interface CertificacionRepository extends JpaRepository<Certificacion, Long> {
    List<Certificacion> findByEmpleadoIdOrderByIdAsc(Long empleadoId);
}
