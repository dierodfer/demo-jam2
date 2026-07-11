package com.empresa.portal.config;

import com.empresa.portal.model.Empleado;
import com.empresa.portal.repo.EmpleadoRepository;
import com.empresa.portal.web.AuthController;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.CommandLineRunner;
import org.springframework.stereotype.Component;

/** Siembra un único empleado (id=1) al arrancar si la tabla está vacía. */
@Component
public class DataSeeder implements CommandLineRunner {

    private final EmpleadoRepository repo;
    private final String seedUsername;

    public DataSeeder(EmpleadoRepository repo, @Value("${app.seed.username}") String seedUsername) {
        this.repo = repo;
        this.seedUsername = seedUsername;
    }

    @Override
    public void run(String... args) {
        if (repo.count() > 0) {
            return;
        }
        Empleado e = new Empleado();
        e.setId(AuthController.EMPLEADO_ID);
        e.setUsername(seedUsername);
        e.setNombre("Ana García");
        e.setEmail("ana.garcia@empresa.com");
        e.setTelefono("+34 600 123 456");
        e.setPuesto("Desarrolladora Full Stack");
        e.setDepartamento("Tecnología");
        e.setDireccion("Calle Mayor 1, 28013 Madrid");
        e.setFoto("https://i.pravatar.cc/300?u=ana.garcia");
        repo.save(e);
    }
}
