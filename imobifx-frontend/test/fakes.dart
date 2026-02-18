import 'package:imobifx_frontend/core/api_client.dart';

class FakeApiClient extends ApiClient {
  FakeApiClient() : super(baseUrl: 'http://localhost:8080');

  Map<String, dynamic>? addressResult;
  Exception? addressError;

  Map<String, dynamic>? quoteResult;
  Exception? quoteError;

  Map<String, dynamic>? createAdResult;
  Exception? createAdError;

  Map<String, dynamic>? listAdsResult;
  Exception? listAdsError;

  @override
  Future<Map<String, dynamic>> lookupAddress(String cep) async {
    if (addressError != null) throw addressError!;
    return addressResult ?? {};
  }

  @override
  Future<Map<String, dynamic>> createQuote({
    required double brlToUsd,
    String? effectiveAt,
  }) async {
    if (quoteError != null) throw quoteError!;
    return quoteResult ?? {'brl_to_usd': brlToUsd};
  }

  @override
  Future<Map<String, dynamic>> createAd({
    required Map<String, String> fields,
    dynamic image,
  }) async {
    if (createAdError != null) throw createAdError!;
    return createAdResult ?? {'id': '1'};
  }

  @override
  Future<Map<String, dynamic>> listAds({
    required int page,
    required int pageSize,
    String? type,
    String? city,
    String? state,
    String? minPrice,
    String? maxPrice,
  }) async {
    if (listAdsError != null) throw listAdsError!;
    return listAdsResult ??
        {
          'page': page,
          'page_size': pageSize,
          'total': 0,
          'items': [],
        };
  }
}
