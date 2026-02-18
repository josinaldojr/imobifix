import 'package:intl/intl.dart';

String formatMoneyBrl(double value, String languageCode) {
  final locale = languageCode == 'pt' ? 'pt_BR' : 'en_US';
  return NumberFormat.currency(locale: locale, symbol: 'R\$').format(value);
}

String formatMoneyUsd(double value, String languageCode) {
  final locale = languageCode == 'pt' ? 'pt_BR' : 'en_US';
  return NumberFormat.currency(locale: locale, symbol: '\$').format(value);
}
