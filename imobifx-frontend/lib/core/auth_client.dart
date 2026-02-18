import 'dart:convert';

import 'package:http/http.dart' as http;

import 'api_client.dart';

class AuthClient {
  AuthClient({required this.baseUrl});

  final String baseUrl;

  Future<String> login({required String username, required String password}) async {
    final uri = Uri.parse('$baseUrl/auth/login');
    final response = await http.post(
      uri,
      headers: {'Content-Type': 'application/json'},
      body: jsonEncode({
        'username': username,
        'password': password,
      }),
    );

    final payload = response.body.isEmpty ? <String, dynamic>{} : jsonDecode(response.body);
    if (response.statusCode >= 200 && response.statusCode < 300) {
      if (payload is Map<String, dynamic> && payload['token'] != null) {
        return payload['token'].toString();
      }
      throw ApiException(response.statusCode, 'Unexpected payload type');
    }

    if (payload is Map<String, dynamic>) {
      final message = (payload['message'] ?? payload['code'] ?? 'Request failed').toString();
      final code = payload['code']?.toString();
      throw ApiException(response.statusCode, message, code: code, details: payload['details']);
    }
    throw ApiException(response.statusCode, 'Request failed');
  }
}
