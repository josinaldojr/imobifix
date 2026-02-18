import 'package:flutter/widgets.dart';

class AppStrings {
  AppStrings(this.locale);

  final Locale locale;

  static const supportedLocales = [Locale('pt'), Locale('en')];

  static AppStrings of(BuildContext context) {
    final locale = Localizations.localeOf(context);
    return AppStrings(locale);
  }

  bool get _isPt => locale.languageCode.toLowerCase() == 'pt';

  String get appTitle => _isPt ? 'ImobiFX' : 'ImobiFX';
  String get adsTab => _isPt ? 'Anuncios' : 'Ads';
  String get createAdTab => _isPt ? 'Novo anuncio' : 'New ad';
  String get quoteTab => _isPt ? 'Cotacao' : 'Quote';
  String get loginTitle => _isPt ? 'Autenticacao' : 'Login';
  String get username => _isPt ? 'Usuario' : 'Username';
  String get password => _isPt ? 'Senha' : 'Password';
  String get login => _isPt ? 'Entrar' : 'Sign in';

  String get type => _isPt ? 'Tipo' : 'Type';
  String get sale => _isPt ? 'Venda' : 'Sale';
  String get rent => _isPt ? 'Aluguel' : 'Rent';
  String get priceBrl => _isPt ? 'Valor BRL' : 'Price BRL';
  String get priceUsd => _isPt ? 'Valor USD' : 'Price USD';
  String get cep => 'CEP';
  String get street => _isPt ? 'Rua' : 'Street';
  String get number => _isPt ? 'Numero' : 'Number';
  String get complement => _isPt ? 'Complemento' : 'Complement';
  String get neighborhood => _isPt ? 'Bairro' : 'Neighborhood';
  String get city => _isPt ? 'Cidade' : 'City';
  String get state => _isPt ? 'UF' : 'State';
  String get image => _isPt ? 'Imagem' : 'Image';
  String get optional => _isPt ? 'Opcional' : 'Optional';
  String get save => _isPt ? 'Salvar' : 'Save';
  String get loading => _isPt ? 'Carregando...' : 'Loading...';
  String get searchCep => _isPt ? 'Buscar CEP' : 'Lookup CEP';
  String get clearFilters => _isPt ? 'Limpar filtros' : 'Clear filters';
  String get filters => _isPt ? 'Filtros' : 'Filters';
  String get minPrice => _isPt ? 'Preco min' : 'Min price';
  String get maxPrice => _isPt ? 'Preco max' : 'Max price';
  String get page => _isPt ? 'Pagina' : 'Page';
  String get previous => _isPt ? 'Anterior' : 'Previous';
  String get next => _isPt ? 'Proxima' : 'Next';
  String get total => _isPt ? 'Total' : 'Total';
  String get noItems => _isPt ? 'Nenhum anuncio encontrado.' : 'No ads found.';
  String get pickImage => _isPt ? 'Selecionar imagem' : 'Pick image';
  String get quoteRate => _isPt ? 'Cotacao BRL -> USD' : 'BRL -> USD quote';
  String get effectiveAt => _isPt ? 'Vigencia (RFC3339)' : 'Effective at (RFC3339)';
  String get addressLookupFailed =>
      _isPt ? 'Nao foi possivel obter CEP. Preencha manualmente.' : 'Could not load CEP. Fill address manually.';
  String get requiredField => _isPt ? 'Campo obrigatorio' : 'Required field';
  String get invalidNumber => _isPt ? 'Numero invalido' : 'Invalid number';
  String get savedSuccess => _isPt ? 'Salvo com sucesso.' : 'Saved successfully.';
  String get requestFailed => _isPt ? 'Falha na requisicao.' : 'Request failed.';

  String errorForCode(String? code) {
    switch (code) {
      case 'IMAGE_TOO_LARGE':
        return _isPt ? 'Imagem excede o tamanho maximo.' : 'Image exceeds maximum size.';
      case 'UNSUPPORTED_IMAGE_TYPE':
        return _isPt ? 'Tipo de imagem nao suportado.' : 'Unsupported image type.';
      case 'VALIDATION_ERROR':
        return _isPt ? 'Dados invalidos.' : 'Invalid data.';
      case 'CEP_INVALID':
        return _isPt ? 'CEP invalido.' : 'Invalid CEP.';
      case 'CEP_NOT_FOUND':
        return _isPt ? 'CEP nao encontrado.' : 'CEP not found.';
      case 'VIA_CEP_UNAVAILABLE':
        return _isPt ? 'ViaCEP indisponivel.' : 'ViaCEP unavailable.';
      case 'INTERNAL_ERROR':
        return _isPt ? 'Erro interno.' : 'Internal error.';
      case 'AUTH_INVALID':
        return _isPt ? 'Credenciais invalidas.' : 'Invalid credentials.';
      default:
        return requestFailed;
    }
  }
}
