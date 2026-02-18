import 'package:flutter/material.dart';
import 'package:flutter_localizations/flutter_localizations.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:imobifx_frontend/l10n/app_strings.dart';
import 'package:imobifx_frontend/pages/create_quote_page.dart';

import 'fakes.dart';

void main() {
  testWidgets('create quote submits and shows success', (tester) async {
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
        home: Scaffold(body: CreateQuotePage(api: api)),
      ),
    );

    await tester.enterText(find.byType(TextFormField).first, '0.25');
    await tester.tap(find.byType(FilledButton));
    await tester.pumpAndSettle();

    expect(find.text('Salvo com sucesso.'), findsOneWidget);
  });

  testWidgets('create quote validates required field', (tester) async {
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
        home: Scaffold(body: CreateQuotePage(api: api)),
      ),
    );

    await tester.tap(find.byType(FilledButton));
    await tester.pump();

    expect(find.text('Campo obrigatorio'), findsOneWidget);
  });
}
