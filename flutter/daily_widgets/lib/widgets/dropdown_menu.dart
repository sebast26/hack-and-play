import 'package:flutter/material.dart';

class DropdownMenuTest extends StatefulWidget {
  const DropdownMenuTest({super.key});

  @override
  State<DropdownMenuTest> createState() => _DropdownMenuTestState();
}

class _DropdownMenuTestState extends State<DropdownMenuTest> {
  final Map<Color, String> colors = {
    Colors.red: 'Red',
    Colors.blue: 'Blue',
    Colors.green: 'Green',
    Colors.purple: 'Purple',
    Colors.yellow: 'Yellow',
  };

  void updateColor(color) {

  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Dropdown Menu')),
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const Text('Seba'),
            DropdownMenu(
              enableSearch: true,
              enableFilter: true,
              onSelected: updateColor,
              dropdownMenuEntries: [
                ...List.generate(5, (index) {
                  final element = colors.entries.elementAt(index);
                  return DropdownMenuEntry(value: element.key, label: element.value);
                }),
              ],
            ),
          ],
        ),
      ),
    );
  }
}
