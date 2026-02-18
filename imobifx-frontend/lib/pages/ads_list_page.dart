import 'package:flutter/material.dart';
import 'package:flutter_svg/flutter_svg.dart';

import '../core/api_client.dart';
import '../core/formatters.dart';
import '../l10n/app_strings.dart';

class AdsListPage extends StatefulWidget {
  const AdsListPage({super.key, required this.api, required this.baseUrl});

  final ApiClient api;
  final String baseUrl;

  @override
  State<AdsListPage> createState() => _AdsListPageState();
}

class _AdsListPageState extends State<AdsListPage> {
  final _cityController = TextEditingController();
  final _stateController = TextEditingController();
  final _minController = TextEditingController();
  final _maxController = TextEditingController();

  bool _loading = false;
  int _page = 1;
  final int _pageSize = 10;
  int _total = 0;
  String? _type;
  List<dynamic> _items = const [];

  @override
  void initState() {
    super.initState();
    _load();
  }

  @override
  void dispose() {
    _cityController.dispose();
    _stateController.dispose();
    _minController.dispose();
    _maxController.dispose();
    super.dispose();
  }

  Future<void> _load() async {
    setState(() => _loading = true);
    try {
      final response = await widget.api.listAds(
        page: _page,
        pageSize: _pageSize,
        type: _type,
        city: _cityController.text,
        state: _stateController.text,
        minPrice: _minController.text,
        maxPrice: _maxController.text,
      );
      setState(() {
        _items = (response['items'] as List<dynamic>? ?? const []);
        _total = (response['total'] as num? ?? 0).toInt();
      });
    } catch (_) {
      setState(() {
        _items = const [];
        _total = 0;
      });
    } finally {
      if (mounted) setState(() => _loading = false);
    }
  }

  bool get _canGoPrev => _page > 1;
  bool get _canGoNext => _page * _pageSize < _total;

  @override
  Widget build(BuildContext context) {
    final s = AppStrings.of(context);
    final lang = Localizations.localeOf(context).languageCode;

    return Padding(
      padding: const EdgeInsets.all(16),
      child: Column(
        children: [
          Wrap(
            spacing: 8,
            runSpacing: 8,
            crossAxisAlignment: WrapCrossAlignment.end,
            children: [
              SizedBox(
                width: 160,
                child: DropdownButtonFormField<String?>(
                  value: _type,
                  decoration: InputDecoration(labelText: s.type),
                  items: [
                    const DropdownMenuItem(value: null, child: Text('-')),
                    DropdownMenuItem(value: 'SALE', child: Text(s.sale)),
                    DropdownMenuItem(value: 'RENT', child: Text(s.rent)),
                  ],
                  onChanged: (value) => setState(() => _type = value),
                ),
              ),
              SizedBox(
                width: 180,
                child: TextField(
                  controller: _cityController,
                  decoration: InputDecoration(labelText: s.city),
                ),
              ),
              SizedBox(
                width: 100,
                child: TextField(
                  controller: _stateController,
                  maxLength: 2,
                  decoration: InputDecoration(labelText: s.state, counterText: ''),
                ),
              ),
              SizedBox(
                width: 140,
                child: TextField(
                  controller: _minController,
                  decoration: InputDecoration(labelText: s.minPrice),
                ),
              ),
              SizedBox(
                width: 140,
                child: TextField(
                  controller: _maxController,
                  decoration: InputDecoration(labelText: s.maxPrice),
                ),
              ),
              FilledButton(
                onPressed: _loading
                    ? null
                    : () {
                        _page = 1;
                        _load();
                      },
                child: Text(s.filters),
              ),
              OutlinedButton(
                onPressed: _loading
                    ? null
                    : () {
                        _type = null;
                        _cityController.clear();
                        _stateController.clear();
                        _minController.clear();
                        _maxController.clear();
                        _page = 1;
                        _load();
                      },
                child: Text(s.clearFilters),
              ),
            ],
          ),
          const SizedBox(height: 12),
          Row(
            children: [
              Text('${s.total}: $_total'),
              const Spacer(),
              Text('${s.page} $_page'),
              const SizedBox(width: 8),
              OutlinedButton(
                onPressed: !_loading && _canGoPrev
                    ? () {
                        setState(() => _page -= 1);
                        _load();
                      }
                    : null,
                child: Text(s.previous),
              ),
              const SizedBox(width: 8),
              OutlinedButton(
                onPressed: !_loading && _canGoNext
                    ? () {
                        setState(() => _page += 1);
                        _load();
                      }
                    : null,
                child: Text(s.next),
              ),
            ],
          ),
          const SizedBox(height: 8),
          Expanded(
            child: _loading
                ? Center(child: Text(s.loading))
                : _items.isEmpty
                    ? Center(child: Text(s.noItems))
                    : ListView.builder(
                        itemCount: _items.length,
                        itemBuilder: (context, index) {
                          final item = _items[index] as Map<String, dynamic>;
                          final address = item['address'] as Map<String, dynamic>? ?? const {};
                          final brl = (item['price_brl'] as num?)?.toDouble() ?? 0;
                          final usd = (item['price_usd'] as num?)?.toDouble();
                          final imageUrl = item['image_url']?.toString();

                          return Card(
                            child: Padding(
                              padding: const EdgeInsets.all(12),
                              child: Row(
                                crossAxisAlignment: CrossAxisAlignment.start,
                                children: [
                                  SizedBox(
                                    width: 120,
                                    height: 90,
                                    child: imageUrl != null && imageUrl.isNotEmpty
                                        ? Image.network(
                                            '${widget.baseUrl}$imageUrl',
                                            fit: BoxFit.cover,
                                            errorBuilder: (_, __, ___) => SvgPicture.asset(
                                              'assets/images/placeholder.svg',
                                              fit: BoxFit.cover,
                                            ),
                                          )
                                        : SvgPicture.asset(
                                            'assets/images/placeholder.svg',
                                            fit: BoxFit.cover,
                                          ),
                                  ),
                                  const SizedBox(width: 12),
                                  Expanded(
                                    child: Column(
                                      crossAxisAlignment: CrossAxisAlignment.start,
                                      children: [
                                        Text('${s.type}: ${item['type']}'),
                                        Text(
                                          '${address['street'] ?? ''}, ${address['number'] ?? '-'} - ${address['neighborhood'] ?? ''}',
                                        ),
                                        Text(
                                          '${address['city'] ?? ''}/${address['state'] ?? ''} - CEP ${address['cep'] ?? ''}',
                                        ),
                                        const SizedBox(height: 4),
                                        Text('${s.priceBrl}: ${formatMoneyBrl(brl, lang)}'),
                                        Text(
                                          '${s.priceUsd}: ${usd == null ? '-' : formatMoneyUsd(usd, lang)}',
                                        ),
                                      ],
                                    ),
                                  ),
                                ],
                              ),
                            ),
                          );
                        },
                      ),
          ),
        ],
      ),
    );
  }
}
