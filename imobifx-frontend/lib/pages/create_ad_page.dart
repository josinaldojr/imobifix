import 'package:file_picker/file_picker.dart';
import 'package:flutter/material.dart';

import '../core/api_client.dart';
import '../l10n/app_strings.dart';

class CreateAdPage extends StatefulWidget {
  const CreateAdPage({super.key, required this.api});

  final ApiClient api;

  @override
  State<CreateAdPage> createState() => _CreateAdPageState();
}

class _CreateAdPageState extends State<CreateAdPage> {
  final _formKey = GlobalKey<FormState>();

  final _priceController = TextEditingController();
  final _cepController = TextEditingController();
  final _streetController = TextEditingController();
  final _numberController = TextEditingController();
  final _complementController = TextEditingController();
  final _neighborhoodController = TextEditingController();
  final _cityController = TextEditingController();
  final _stateController = TextEditingController();

  String _type = 'SALE';
  bool _saving = false;
  bool _lookingUpCep = false;
  PlatformFile? _pickedImage;

  @override
  void dispose() {
    _priceController.dispose();
    _cepController.dispose();
    _streetController.dispose();
    _numberController.dispose();
    _complementController.dispose();
    _neighborhoodController.dispose();
    _cityController.dispose();
    _stateController.dispose();
    super.dispose();
  }

  Future<void> _lookupCep(AppStrings s) async {
    final cep = _cepController.text.trim();
    if (cep.isEmpty) return;

    setState(() => _lookingUpCep = true);
    try {
      final addr = await widget.api.lookupAddress(cep);
      _streetController.text = (addr['street'] ?? '').toString();
      _neighborhoodController.text = (addr['neighborhood'] ?? '').toString();
      _cityController.text = (addr['city'] ?? '').toString();
      _stateController.text = (addr['state'] ?? '').toString();
    } on ApiException catch (e) {
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text(s.errorForCode(e.code))),
      );
    } catch (_) {
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text(s.addressLookupFailed)),
      );
    } finally {
      if (mounted) setState(() => _lookingUpCep = false);
    }
  }

  Future<void> _pickImage() async {
    final result = await FilePicker.platform.pickFiles(
      allowMultiple: false,
      type: FileType.image,
      withData: true,
    );
    if (result == null || result.files.isEmpty) return;
    setState(() => _pickedImage = result.files.first);
  }

  Future<void> _submit(AppStrings s) async {
    if (!(_formKey.currentState?.validate() ?? false)) return;

    setState(() => _saving = true);
    try {
      final fields = <String, String>{
        'type': _type,
        'price_brl': _priceController.text.trim().replaceAll(',', '.'),
        'cep': _cepController.text.trim(),
        'street': _streetController.text.trim(),
        'number': _numberController.text.trim(),
        'complement': _complementController.text.trim(),
        'neighborhood': _neighborhoodController.text.trim(),
        'city': _cityController.text.trim(),
        'state': _stateController.text.trim().toUpperCase(),
      };

      await widget.api.createAd(fields: fields, image: _pickedImage);
      if (!mounted) return;

      ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text(s.savedSuccess)));
      _formKey.currentState?.reset();
      _priceController.clear();
      _cepController.clear();
      _streetController.clear();
      _numberController.clear();
      _complementController.clear();
      _neighborhoodController.clear();
      _cityController.clear();
      _stateController.clear();
      setState(() {
        _type = 'SALE';
        _pickedImage = null;
      });
    } on ApiException catch (e) {
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text(s.errorForCode(e.code))),
      );
    } finally {
      if (mounted) setState(() => _saving = false);
    }
  }

  String? _requiredValidator(String? value, AppStrings s) {
    if ((value ?? '').trim().isEmpty) return s.requiredField;
    return null;
  }

  @override
  Widget build(BuildContext context) {
    final s = AppStrings.of(context);
    return SingleChildScrollView(
      padding: const EdgeInsets.all(16),
      child: Form(
        key: _formKey,
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            DropdownButtonFormField<String>(
              value: _type,
              decoration: InputDecoration(labelText: s.type),
              items: [
                DropdownMenuItem(value: 'SALE', child: Text(s.sale)),
                DropdownMenuItem(value: 'RENT', child: Text(s.rent)),
              ],
              onChanged: (value) => setState(() => _type = value ?? 'SALE'),
            ),
            const SizedBox(height: 12),
            TextFormField(
              controller: _priceController,
              decoration: InputDecoration(labelText: s.priceBrl),
              validator: (value) {
                final required = _requiredValidator(value, s);
                if (required != null) return required;
                final parsed = double.tryParse((value ?? '').trim().replaceAll(',', '.'));
                if (parsed == null || parsed < 0) return s.invalidNumber;
                return null;
              },
            ),
            const SizedBox(height: 12),
            Row(
              children: [
                Expanded(
                  child: TextFormField(
                    controller: _cepController,
                    decoration: InputDecoration(labelText: s.cep),
                    validator: (value) => _requiredValidator(value, s),
                  ),
                ),
                const SizedBox(width: 8),
                OutlinedButton(
                  onPressed: _lookingUpCep ? null : () => _lookupCep(s),
                  child: Text(_lookingUpCep ? s.loading : s.searchCep),
                ),
              ],
            ),
            const SizedBox(height: 12),
            TextFormField(
              controller: _streetController,
              decoration: InputDecoration(labelText: s.street),
              validator: (value) => _requiredValidator(value, s),
            ),
            const SizedBox(height: 12),
            Row(
              children: [
                Expanded(
                  child: TextFormField(
                    controller: _numberController,
                    decoration: InputDecoration(labelText: '${s.number} (${s.optional})'),
                  ),
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: TextFormField(
                    controller: _complementController,
                    decoration: InputDecoration(labelText: '${s.complement} (${s.optional})'),
                  ),
                ),
              ],
            ),
            const SizedBox(height: 12),
            TextFormField(
              controller: _neighborhoodController,
              decoration: InputDecoration(labelText: s.neighborhood),
              validator: (value) => _requiredValidator(value, s),
            ),
            const SizedBox(height: 12),
            Row(
              children: [
                Expanded(
                  child: TextFormField(
                    controller: _cityController,
                    decoration: InputDecoration(labelText: s.city),
                    validator: (value) => _requiredValidator(value, s),
                  ),
                ),
                const SizedBox(width: 12),
                SizedBox(
                  width: 120,
                  child: TextFormField(
                    controller: _stateController,
                    maxLength: 2,
                    decoration: InputDecoration(labelText: s.state, counterText: ''),
                    validator: (value) {
                      final raw = (value ?? '').trim();
                      if (raw.length != 2) return s.requiredField;
                      return null;
                    },
                  ),
                ),
              ],
            ),
            const SizedBox(height: 12),
            Row(
              children: [
                OutlinedButton(
                  onPressed: _pickImage,
                  child: Text(s.pickImage),
                ),
                const SizedBox(width: 8),
                Expanded(
                  child: Text(
                    _pickedImage?.name ?? '-',
                    overflow: TextOverflow.ellipsis,
                  ),
                ),
              ],
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
