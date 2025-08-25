import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';

enum UserColor { blue, green, red, yellow }

class CupertinoRadioTest extends StatefulWidget {
  const CupertinoRadioTest({super.key});

  @override
  State<CupertinoRadioTest> createState() => _CupertinoRadioTestState();
}

class _CupertinoRadioTestState extends State<CupertinoRadioTest> {
  UserColor? _color;

  void _updateColor(UserColor? value) {
    setState(() {
      _color = value;
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: Text('CupertinoRadio')),
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Text('Seba'),
            CupertinoListTile(
              title: const Text('Blue'),
              leading: CupertinoRadio<UserColor>(
                value: UserColor.blue,
                groupValue: _color,
                onChanged: (UserColor? value) => _updateColor(value),
                activeColor: Colors.blue,
              ),
            ),
            CupertinoListTile(
              title: const Text('Green'),
              leading: CupertinoRadio<UserColor>(
                value: UserColor.green,
                groupValue: _color,
                onChanged: (UserColor? value) => _updateColor(value),
                activeColor: Colors.green,
              ),
            ),
            CupertinoListTile(
              title: const Text('Red'),
              leading: CupertinoRadio<UserColor>(
                value: UserColor.red,
                groupValue: _color,
                onChanged: (UserColor? value) => _updateColor(value),
                activeColor: Colors.red,
              ),
            ),
            CupertinoListTile(
              title: const Text('Yellow'),
              leading: CupertinoRadio<UserColor>(
                value: UserColor.yellow,
                groupValue: _color,
                onChanged: (UserColor? value) => _updateColor(value),
                activeColor: Colors.yellow,
              ),
            ),
          ],
        ),
      ),
    );
  }
}
