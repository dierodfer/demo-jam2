package com.empresa.portal.web;

import com.empresa.portal.model.Certificacion;
import com.empresa.portal.repo.CertificacionRepository;
import com.empresa.portal.web.dto.CertificacionDto;
import com.empresa.portal.web.dto.CertificacionInput;
import jakarta.servlet.http.HttpSession;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.DeleteMapping;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.PutMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.List;
import java.util.Map;

/** CRUD de conocimientos / certificaciones del empleado autenticado. */
@RestController
@RequestMapping("/api/certificaciones")
public class CertificacionController {

    private final CertificacionRepository repo;

    public CertificacionController(CertificacionRepository repo) {
        this.repo = repo;
    }

    @GetMapping
    public ResponseEntity<?> list(HttpSession session) {
        Long empleadoId = currentEmpleadoId(session);
        if (empleadoId == null) {
            return unauthorized();
        }
        List<CertificacionDto> dtos = repo.findByEmpleadoIdOrderByIdAsc(empleadoId)
                .stream().map(CertificacionDto::from).toList();
        return ResponseEntity.ok(dtos);
    }

    @PostMapping
    public ResponseEntity<?> create(@RequestBody CertificacionInput body, HttpSession session) {
        Long empleadoId = currentEmpleadoId(session);
        if (empleadoId == null) {
            return unauthorized();
        }
        Certificacion c = new Certificacion();
        c.setEmpleadoId(empleadoId);
        apply(c, body);
        repo.save(c);
        return ResponseEntity.status(201).body(CertificacionDto.from(c));
    }

    @PutMapping("/{id}")
    public ResponseEntity<?> update(@PathVariable Long id, @RequestBody CertificacionInput body, HttpSession session) {
        Long empleadoId = currentEmpleadoId(session);
        if (empleadoId == null) {
            return unauthorized();
        }
        Certificacion c = repo.findById(id).orElse(null);
        if (c == null || !empleadoId.equals(c.getEmpleadoId())) {
            return notFound();
        }
        apply(c, body);
        repo.save(c);
        return ResponseEntity.ok(CertificacionDto.from(c));
    }

    @DeleteMapping("/{id}")
    public ResponseEntity<?> delete(@PathVariable Long id, HttpSession session) {
        Long empleadoId = currentEmpleadoId(session);
        if (empleadoId == null) {
            return unauthorized();
        }
        Certificacion c = repo.findById(id).orElse(null);
        if (c == null || !empleadoId.equals(c.getEmpleadoId())) {
            return notFound();
        }
        repo.delete(c);
        return ResponseEntity.noContent().build();
    }

    private void apply(Certificacion c, CertificacionInput body) {
        c.setConocimiento(body.conocimiento());
        c.setEmpresaEmisora(body.empresaEmisora());
        c.setFecha(body.fecha());
    }

    private Long currentEmpleadoId(HttpSession session) {
        Object id = session.getAttribute(AuthController.SESSION_KEY);
        return (id instanceof Long empleadoId) ? empleadoId : null;
    }

    private ResponseEntity<?> unauthorized() {
        return ResponseEntity.status(401).body(Map.of("error", "No autenticado"));
    }

    private ResponseEntity<?> notFound() {
        return ResponseEntity.status(404).body(Map.of("error", "No encontrado"));
    }
}
