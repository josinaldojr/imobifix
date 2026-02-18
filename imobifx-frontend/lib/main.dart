import 'package:flutter/material.dart';
import 'package:flutter_localizations/flutter_localizations.dart';

import 'core/api_client.dart';
import 'l10n/app_strings.dart';
import 'pages/ads_list_page.dart';
import 'pages/create_ad_page.dart';
import 'pages/create_quote_page.dart';

const _apiBaseUrl = String.fromEnvironment(
  'API_BASE_URL',
  defaultValue: 'http://localhost:8080',
);
const _appLocale = String.fromEnvironment(
  'APP_LOCALE',
  defaultValue: 'pt',
);

void main() {
  final locale = _appLocale.toLowerCase().startsWith('en') ? const Locale('en') : const Locale('pt');
  runApp(ImobiFxFrontend(locale: locale));
}

class ImobiFxFrontend extends StatelessWidget {
  const ImobiFxFrontend({super.key, required this.locale});

  final Locale locale;

  @override
  Widget build(BuildContext context) {
    final strings = AppStrings(locale);
    final api = ApiClient(baseUrl: _apiBaseUrl);

    return MaterialApp(
      title: strings.appTitle,
      locale: locale,
      supportedLocales: AppStrings.supportedLocales,
      localizationsDelegates: const [
        GlobalMaterialLocalizations.delegate,
        GlobalWidgetsLocalizations.delegate,
        GlobalCupertinoLocalizations.delegate,
      ],
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(seedColor: const Color(0xFF1F8A70)),
        useMaterial3: true,
      ),
      home: RootPage(api: api, baseUrl: _apiBaseUrl),
    );
  }
}

class RootPage extends StatefulWidget {
  const RootPage({super.key, required this.api, required this.baseUrl});

  final ApiClient api;
  final String baseUrl;

  @override
  State<RootPage> createState() => _RootPageState();
}

class _RootPageState extends State<RootPage> {
  int _index = 0;

  @override
  Widget build(BuildContext context) {
    final s = AppStrings.of(context);
    final pages = [
      AdsListPage(api: widget.api, baseUrl: widget.baseUrl),
      CreateAdPage(api: widget.api),
      CreateQuotePage(api: widget.api),
    ];

    return Scaffold(
      appBar: AppBar(title: Text(s.appTitle)),
      body: pages[_index],
      bottomNavigationBar: NavigationBar(
        selectedIndex: _index,
        onDestinationSelected: (index) => setState(() => _index = index),
        destinations: [
          NavigationDestination(icon: const Icon(Icons.list_alt), label: s.adsTab),
          NavigationDestination(icon: const Icon(Icons.home_work), label: s.createAdTab),
          NavigationDestination(icon: const Icon(Icons.currency_exchange), label: s.quoteTab),
        ],
      ),
    );
  }
}
