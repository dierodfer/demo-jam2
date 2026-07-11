package com.empresa.portal.web;

import com.empresa.portal.model.Empleado;
import com.empresa.portal.repo.EmpleadoRepository;
import com.empresa.portal.web.dto.EmpleadoDto;
import com.empresa.portal.web.dto.EmpleadoUpdate;
import jakarta.servlet.http.HttpSession;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PutMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.Map;

/** Perfil del empleado autenticado. Requiere sesión activa (401 si no hay). */
@RestController
@RequestMapping("/api")
public class MeController {

    private final EmpleadoRepository repo;

    public MeController(EmpleadoRepository repo) {
        this.repo = repo;
    }

    @GetMapping("/me")
    public ResponseEntity<?> me(HttpSession session) {
        Empleado empleado = current(session);
        if (empleado == null) {
            return unauthorized();
        }
        return ResponseEntity.ok(EmpleadoDto.from(empleado));
    }

    @PutMapping("/me")
    public ResponseEntity<?> update(@RequestBody EmpleadoUpdate body, HttpSession session) {
        Empleado empleado = current(session);
        if (empleado == null) {
            return unauthorized();
        }
        if (body.nombre() != null) empleado.setNombre(body.nombre());
        if (body.email() != null) empleado.setEmail(body.email());
        if (body.telefono() != null) empleado.setTelefono(body.telefono());
        if (body.puesto() != null) empleado.setPuesto(body.puesto());
        if (body.departamento() != null) empleado.setDepartamento(body.departamento());
        if (body.direccion() != null) empleado.setDireccion(body.direccion());
        if (body.foto() != null) empleado.setFoto(body.foto());
        repo.save(empleado);
        return ResponseEntity.ok(EmpleadoDto.from(empleado));
    }

    private Empleado current(HttpSession session) {
        Object id = session.getAttribute(AuthController.SESSION_KEY);
        if (!(id instanceof Long empleadoId)) {
            return null;
        }
        return repo.findById(empleadoId).orElse(null);
    }

    private ResponseEntity<?> unauthorized() {
        return ResponseEntity.status(401).body(Map.of("error", "No autenticado"));
    }
}
