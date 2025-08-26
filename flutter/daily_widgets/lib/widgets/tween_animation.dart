import 'package:flutter/material.dart';

class TextSizeAnimationTest extends StatefulWidget {
  const TextSizeAnimationTest({super.key});

  @override
  State<TextSizeAnimationTest> createState() => _TextSizeAnimationTestState();
}

class _TextSizeAnimationTestState extends State<TextSizeAnimationTest>
    with SingleTickerProviderStateMixin {
  late Animation<double> _animation;
  late AnimationController _controller;

  @override
  void initState() {
    super.initState();
    _controller = AnimationController(
      duration: const Duration(seconds: 2),
      vsync: this,
    );
    _animation = Tween<double>(
      begin: 16.0,
      end: 52.0,
    ).animate(_controller)
    ..addListener(() {
      setState(() {
        
      });
    })
    ..addStatusListener((status) => debugPrint('$status'));
    _controller.forward();
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: Text('CupertinoRadio')),
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Text('Seba', style: TextStyle(fontSize: _animation.value)),
          ],
        ),
      ),
    );
  }
}
