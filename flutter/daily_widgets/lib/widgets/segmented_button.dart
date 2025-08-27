import 'package:flutter/material.dart';

class SegmentedButtonTest extends StatefulWidget {
  const SegmentedButtonTest({super.key});

  @override
  State<SegmentedButtonTest> createState() => _SegmentedButtonTestState();
}

class _SegmentedButtonTestState extends State<SegmentedButtonTest> {
  final Map<String, IconData> options = {
    'Inbox': Icons.inbox,
    'Primary': Icons.priority_high,
    'Everything else': Icons.web_stories,
  };
  Set<String> _selected = {'Inbox'};

  void updateSelected(Set<String> newSelection) {
    setState(() {
      _selected = newSelection;
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: Text('SegmentedButton')),
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const Text('Seba'),
            SegmentedButton(
              multiSelectionEnabled: false,
              selected: _selected,
              onSelectionChanged: updateSelected,
              segments: [
                ...List.generate(3, (index) {
                  var element = options.entries.elementAt(index);
                  return ButtonSegment<String>(
                    value: element.key, 
                    label: Text(element.key),
                    icon: Icon(element.value)
                  );
                }),
              ],
            ),
          ],
        ),
      ),
    );
  }
}
