import 'dart:convert';

import 'package:file_picker/file_picker.dart';
import 'package:http/http.dart' as http;
import 'package:http_parser/http_parser.dart';

class ApiException implements Exception {
  ApiException(this.statusCode, this.message, {this.code, this.details});

  final int statusCode;
  final String message;
  final String? code;
  final dynamic details;

  @override
  String toString() => 'ApiException($statusCode): $message';
}

class ApiClient {
  ApiClient({required this.baseUrl});

  final String baseUrl;
  String? _authToken;

  void setAuthToken(String? token) {
    _authToken = token;
  }

  Future<Map<String, dynamic>> lookupAddress(String cep) async {
    final uri = Uri.parse('$baseUrl/api/addresses/$cep');
    final response = await http.get(uri, headers: _headers());
    return _decodeObject(response);
  }

  Future<Map<String, dynamic>> createQuote({
    required double brlToUsd,
    String? effectiveAt,
  }) async {
    final uri = Uri.parse('$baseUrl/api/quotes');
    final body = <String, dynamic>{'brl_to_usd': brlToUsd};
    if (effectiveAt != null && effectiveAt.trim().isNotEmpty) {
      body['effective_at'] = effectiveAt.trim();
    }
    final response = await http.post(
      uri,
      headers: _headers(contentType: 'application/json'),
      body: jsonEncode(body),
    );
    return _decodeObject(response);
  }

  Future<Map<String, dynamic>> createAd({
    required Map<String, String> fields,
    PlatformFile? image,
  }) async {
    final uri = Uri.parse('$baseUrl/api/ads');
    final req = http.MultipartRequest('POST', uri);
    req.fields.addAll(fields);
    final headers = _headers();
    if (headers.isNotEmpty) {
      req.headers.addAll(headers);
    }

    if (image != null && image.bytes != null) {
      final mime = _inferMimeType(image);
      req.files.add(
        http.MultipartFile.fromBytes(
          'image',
          image.bytes!,
          filename: image.name,
          contentType: mime == null ? null : MediaType.parse(mime),
        ),
      );
    }

    final streamed = await req.send();
    final response = await http.Response.fromStream(streamed);
    return _decodeObject(response);
  }

  Future<Map<String, dynamic>> listAds({
    required int page,
    required int pageSize,
    String? type,
    String? city,
    String? state,
    String? minPrice,
    String? maxPrice,
  }) async {
    final params = <String, String>{
      'page': '$page',
      'page_size': '$pageSize',
    };

    void putIfNotBlank(String key, String? value) {
      if (value != null && value.trim().isNotEmpty) {
        params[key] = value.trim();
      }
    }

    putIfNotBlank('type', type);
    putIfNotBlank('city', city);
    putIfNotBlank('state', state);
    putIfNotBlank('min_price', minPrice);
    putIfNotBlank('max_price', maxPrice);

    final uri = Uri.parse('$baseUrl/api/ads').replace(queryParameters: params);
    final response = await http.get(uri, headers: _headers());
    return _decodeObject(response);
  }

  Map<String, dynamic> _decodeObject(http.Response response) {
    final payload = response.body.isEmpty ? <String, dynamic>{} : jsonDecode(response.body);
    if (response.statusCode >= 200 && response.statusCode < 300) {
      if (payload is Map<String, dynamic>) {
        return payload;
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

  String? _inferMimeType(PlatformFile file) {
    final lower = file.name.toLowerCase();
    if (lower.endsWith('.png')) return 'image/png';
    if (lower.endsWith('.jpg') || lower.endsWith('.jpeg')) return 'image/jpeg';
    if (lower.endsWith('.webp')) return 'image/webp';
    return null;
  }

  Map<String, String> _headers({String? contentType}) {
    final out = <String, String>{};
    if (contentType != null) {
      out['Content-Type'] = contentType;
    }
    if (_authToken != null && _authToken!.trim().isNotEmpty) {
      out['Authorization'] = 'Bearer ${_authToken!}';
    }
    return out;
  }
}
