import 'package:flutter/material.dart';

import '../core/api_client.dart';
import '../l10n/app_strings.dart';

class CreateQuotePage extends StatefulWidget {
  const CreateQuotePage({super.key, required this.api});

  final ApiClient api;

  @override
  State<CreateQuotePage> createState() => _CreateQuotePageState();
}

class _CreateQuotePageState extends State<CreateQuotePage> {
  final _formKey = GlobalKey<FormState>();
  final _rateController = TextEditingController();
  final _effectiveAtController = TextEditingController();

  bool _saving = false;

  @override
  void dispose() {
    _rateController.dispose();
    _effectiveAtController.dispose();
    super.dispose();
  }

  Future<void> _submit(AppStrings s) async {
    if (!(_formKey.currentState?.validate() ?? false)) return;

    setState(() => _saving = true);
    try {
      final rate = double.parse(_rateController.text.trim().replaceAll(',', '.'));
      await widget.api.createQuote(
        brlToUsd: rate,
        effectiveAt: _effectiveAtController.text.trim(),
      );

      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text(s.savedSuccess)));
      _formKey.currentState?.reset();
      _rateController.clear();
      _effectiveAtController.clear();
    } on ApiException catch (e) {
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text(s.errorForCode(e.code))),
      );
    } finally {
      if (mounted) setState(() => _saving = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    final s = AppStrings.of(context);
    return Padding(
      padding: const EdgeInsets.all(16),
      child: Form(
        key: _formKey,
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              s.quoteRate,
              style: Theme.of(context).textTheme.titleLarge,
            ),
            const SizedBox(height: 16),
            TextFormField(
              controller: _rateController,
              decoration: InputDecoration(labelText: s.quoteRate),
              validator: (value) {
                final raw = (value ?? '').trim();
                if (raw.isEmpty) return s.requiredField;
                final parsed = double.tryParse(raw.replaceAll(',', '.'));
                if (parsed == null || parsed <= 0) return s.invalidNumber;
                return null;
              },
            ),
            const SizedBox(height: 12),
            TextFormField(
              controller: _effectiveAtController,
              decoration: InputDecoration(
                labelText: s.effectiveAt,
                hintText: '2026-02-16T10:00:00Z',
              ),
            ),
            const SizedBox(height: 20),
            FilledButton(
              onPressed: _saving ? null : () => _submit(s),
              child: Text(_saving ? s.loading : s.save),
            ),
          ],
        ),
      ),
    );
  }
}
