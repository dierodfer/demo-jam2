package com.empresa.portal.web.dto;

/** Cuerpo de POST /api/login. La contraseña se ignora (login simulado). */
public record LoginRequest(String username, String password) {
}
