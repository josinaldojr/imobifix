import 'package:flutter/material.dart';
import 'package:flutter_localizations/flutter_localizations.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:imobifx_frontend/core/api_client.dart';
import 'package:imobifx_frontend/l10n/app_strings.dart';
import 'package:imobifx_frontend/pages/create_ad_page.dart';

import 'fakes.dart';

void main() {
  testWidgets('create ad validates required fields', (tester) async {
    final api = FakeApiClient();

    await tester.pumpWidget(
      MaterialApp(
        locale: const Locale('pt'),
        supportedLocales: AppStrings.supportedLocales,
        localizationsDelegates: const [
          GlobalMaterialLocalizations.delegate,
          GlobalWidgetsLocalizations.delegate,
          GlobalCupertinoLocalizations.delegate,
        ],
        home: Scaffold(body: CreateAdPage(api: api)),
      ),
    );

    await tester.tap(find.byType(FilledButton));
    await tester.pump();

    expect(find.text('Campo obrigatorio'), findsWidgets);
  });

  testWidgets('lookup CEP shows localized error for backend code', (tester) async {
    final api = FakeApiClient()
      ..addressError = ApiException(503, 'CEP', code: 'VIA_CEP_UNAVAILABLE');

    await tester.pumpWidget(
      MaterialApp(
        locale: const Locale('en'),
        supportedLocales: AppStrings.supportedLocales,
        localizationsDelegates: const [
          GlobalMaterialLocalizations.delegate,
          GlobalWidgetsLocalizations.delegate,
          GlobalCupertinoLocalizations.delegate,
        ],
        home: Scaffold(body: CreateAdPage(api: api)),
      ),
    );

    await tester.enterText(find.widgetWithText(TextFormField, 'CEP'), '58000000');
    await tester.tap(find.text('Lookup CEP'));
    await tester.pump();

    expect(find.text('ViaCEP unavailable.'), findsOneWidget);
  });
}
