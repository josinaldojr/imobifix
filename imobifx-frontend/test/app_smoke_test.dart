import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:imobifx_frontend/core/api_client.dart';
import 'package:imobifx_frontend/main.dart';

void main() {
  testWidgets('renders root app shell', (tester) async {
    await tester.pumpWidget(
      const ImobiFxFrontend(locale: Locale('pt')),
    );

    expect(find.text('ImobiFX'), findsOneWidget);
    expect(find.text('Anuncios'), findsOneWidget);
  });

  test('api client stores base url', () {
    final api = ApiClient(baseUrl: 'http://localhost:8080');
    expect(api.baseUrl, 'http://localhost:8080');
  });
}
