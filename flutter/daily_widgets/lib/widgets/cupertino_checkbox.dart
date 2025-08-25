import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';

class CupertinoCheckboxTest extends StatefulWidget {
  const CupertinoCheckboxTest({super.key});

  @override
  State<CupertinoCheckboxTest> createState() => _CupertinoCheckboxTestState();
}

class _CupertinoCheckboxTestState extends State<CupertinoCheckboxTest> {
  bool? _firstChecked;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: Text('CupertinoRadio')),
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const Text('Seba'),
            CupertinoCheckbox(
              value: _firstChecked,
              tristate: true,
              onChanged: (value) => setState(() {
                _firstChecked = value;
              }),
            ),
          ],
        ),
      ),
    );
  }
}
