import 'package:flutter/material.dart';
import 'package:flutter_localizations/flutter_localizations.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:imobifx_frontend/l10n/app_strings.dart';
import 'package:imobifx_frontend/pages/ads_list_page.dart';

import 'fakes.dart';

void main() {
  testWidgets('ads list renders items and prices', (tester) async {
    final api = FakeApiClient()
      ..listAdsResult = {
        'page': 1,
        'page_size': 10,
        'total': 1,
        'items': [
          {
            'id': '1',
            'type': 'SALE',
            'price_brl': 100.0,
            'price_usd': 20.0,
            'image_url': null,
            'address': {
              'cep': '58000-000',
              'street': 'Rua A',
              'number': null,
              'complement': null,
              'neighborhood': 'Centro',
              'city': 'Joao Pessoa',
              'state': 'PB',
            },
            'created_at': '2026-02-16T10:00:00Z',
          }
        ],
      };

    await tester.pumpWidget(
      MaterialApp(
        locale: const Locale('pt'),
        supportedLocales: AppStrings.supportedLocales,
        localizationsDelegates: const [
          GlobalMaterialLocalizations.delegate,
          GlobalWidgetsLocalizations.delegate,
          GlobalCupertinoLocalizations.delegate,
        ],
        home: Scaffold(body: AdsListPage(api: api, baseUrl: 'http://localhost:8080')),
      ),
    );

    await tester.pumpAndSettle();

    expect(find.text('Tipo: SALE'), findsOneWidget);
    expect(find.textContaining('CEP 58000-000'), findsOneWidget);
    expect(find.textContaining('Valor BRL:'), findsOneWidget);
    expect(find.textContaining('Valor USD:'), findsOneWidget);
  });

  testWidgets('pagination controls enable and disable correctly', (tester) async {
    final api = FakeApiClient()
      ..listAdsResult = {
        'page': 1,
        'page_size': 10,
        'total': 15,
        'items': [],
      };

    await tester.pumpWidget(
      MaterialApp(
        locale: const Locale('pt'),
        supportedLocales: AppStrings.supportedLocales,
        localizationsDelegates: const [
          GlobalMaterialLocalizations.delegate,
          GlobalWidgetsLocalizations.delegate,
          GlobalCupertinoLocalizations.delegate,
        ],
        home: Scaffold(body: AdsListPage(api: api, baseUrl: 'http://localhost:8080')),
      ),
    );

    await tester.pumpAndSettle();

    final prev = find.widgetWithText(OutlinedButton, 'Anterior');
    final next = find.widgetWithText(OutlinedButton, 'Proxima');

    final prevButton = tester.widget<OutlinedButton>(prev);
    final nextButton = tester.widget<OutlinedButton>(next);

    expect(prevButton.onPressed, isNull);
    expect(nextButton.onPressed, isNotNull);
  });
}
