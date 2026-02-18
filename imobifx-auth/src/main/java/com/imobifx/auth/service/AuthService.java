package com.imobifx.auth.service;

import java.time.Instant;
import java.time.temporal.ChronoUnit;
import java.util.Objects;
import java.util.UUID;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;

@Service
public class AuthService {
  private final String user;
  private final String pass;

  public AuthService(
      @Value("${auth.user:admin}") String user,
      @Value("${auth.pass:admin123}") String pass) {
    this.user = user;
    this.pass = pass;
  }

  public boolean validate(String username, String password) {
    return Objects.equals(user, username) && Objects.equals(pass, password);
  }

  public String newToken() {
    return UUID.randomUUID().toString();
  }

  public String expiresAtIso() {
    return Instant.now().plus(2, ChronoUnit.HOURS).toString();
  }
}
