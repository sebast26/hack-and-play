import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';

class CupertinoSwitchTest extends StatefulWidget {
  const CupertinoSwitchTest({super.key});

  @override
  State<CupertinoSwitchTest> createState() => _CupertinoSwitchTestState();
}

class _CupertinoSwitchTestState extends State<CupertinoSwitchTest> {
  bool _darkModeEnabled = false;

  void updateDarkMode(bool value) {
    setState(() {
      _darkModeEnabled = value;
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
            const Text('Seba'),
            Row(
              children: [
                const Text('Dark mode enabled?'),
                CupertinoSwitch(
                  value: _darkModeEnabled,
                  onChanged: updateDarkMode,
                  activeTrackColor: CupertinoColors.systemPurple,
                  inactiveTrackColor: CupertinoColors.darkBackgroundGray,
                  thumbColor: CupertinoColors.activeOrange,
                  inactiveThumbColor: CupertinoColors.systemTeal,
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }
}
