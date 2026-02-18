package com.imobifx.auth.api;

import jakarta.validation.Valid;

import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RestController;

import com.imobifx.auth.service.AuthService;

@RestController
public class AuthController {
  private final AuthService service;

  public AuthController(AuthService service) {
    this.service = service;
  }

  @PostMapping("/auth/login")
  public ResponseEntity<?> login(@Valid @RequestBody LoginRequest req) {
    if (!service.validate(req.getUsername(), req.getPassword())) {
      return ResponseEntity.status(HttpStatus.UNAUTHORIZED)
          .body(new ErrorResponse("AUTH_INVALID", "Credenciais invalidas."));
    }
    return ResponseEntity.ok(new LoginResponse(service.newToken(), service.expiresAtIso()));
  }
}
