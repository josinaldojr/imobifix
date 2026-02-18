import 'package:flutter/material.dart';

import '../core/auth_client.dart';
import '../core/api_client.dart';
import '../l10n/app_strings.dart';

class LoginPage extends StatefulWidget {
  const LoginPage({super.key, required this.auth, required this.onLoggedIn});

  final AuthClient auth;
  final void Function(String token) onLoggedIn;

  @override
  State<LoginPage> createState() => _LoginPageState();
}

class _LoginPageState extends State<LoginPage> {
  final _formKey = GlobalKey<FormState>();
  final _userController = TextEditingController();
  final _passController = TextEditingController();
  bool _loading = false;

  @override
  void dispose() {
    _userController.dispose();
    _passController.dispose();
    super.dispose();
  }

  Future<void> _submit(AppStrings s) async {
    if (!(_formKey.currentState?.validate() ?? false)) return;
    setState(() => _loading = true);
    try {
      final token = await widget.auth.login(
        username: _userController.text.trim(),
        password: _passController.text,
      );
      if (!mounted) return;
      widget.onLoggedIn(token);
    } on ApiException catch (e) {
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text(s.errorForCode(e.code))),
      );
    } finally {
      if (mounted) setState(() => _loading = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    final s = AppStrings.of(context);
    return Scaffold(
      appBar: AppBar(title: Text(s.loginTitle)),
      body: Center(
        child: ConstrainedBox(
          constraints: const BoxConstraints(maxWidth: 360),
          child: Padding(
            padding: const EdgeInsets.all(16),
            child: Form(
              key: _formKey,
              child: Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  TextFormField(
                    controller: _userController,
                    decoration: InputDecoration(labelText: s.username),
                    validator: (value) {
                      if ((value ?? '').trim().isEmpty) return s.requiredField;
                      return null;
                    },
                  ),
                  const SizedBox(height: 12),
                  TextFormField(
                    controller: _passController,
                    decoration: InputDecoration(labelText: s.password),
                    obscureText: true,
                    validator: (value) {
                      if ((value ?? '').isEmpty) return s.requiredField;
                      return null;
                    },
                  ),
                  const SizedBox(height: 20),
                  FilledButton(
                    onPressed: _loading ? null : () => _submit(s),
                    child: Text(_loading ? s.loading : s.login),
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }
}
