package com.imobifx.auth.api;

public class LoginResponse {
  private String token;
  private String expiresAt;

  public LoginResponse(String token, String expiresAt) {
    this.token = token;
    this.expiresAt = expiresAt;
  }

  public String getToken() {
    return token;
  }

  public String getExpiresAt() {
    return expiresAt;
  }
}
