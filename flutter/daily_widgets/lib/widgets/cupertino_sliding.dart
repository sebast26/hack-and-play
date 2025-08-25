import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';

enum Emotion { happy, sad, shocked }

class CupertinoSlidingSegmentedControlTest extends StatefulWidget {
  const CupertinoSlidingSegmentedControlTest({super.key});

  @override
  State<CupertinoSlidingSegmentedControlTest> createState() =>
      _CupertinoSlidingSegmentedControlTestState();
}

class _CupertinoSlidingSegmentedControlTestState
    extends State<CupertinoSlidingSegmentedControlTest> {
  Emotion? _emotion;

  void _updateEmotion(Emotion? value) {
    setState(() {
      _emotion = value;
    });
  }

  var segmentsMap = <Emotion, Widget>{
    Emotion.happy: Icon(Icons.tag_faces_outlined),
    Emotion.sad: Icon(Icons.face_retouching_natural),
    Emotion.shocked: Icon(Icons.face),
  };

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: Text('CupertinoRadio')),
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const Text('Seba'),
            CupertinoSlidingSegmentedControl(
              children: segmentsMap,
              groupValue: _emotion,
              padding: EdgeInsets.fromLTRB(20, 10, 20, 10),
              thumbColor: Colors.greenAccent,
              backgroundColor: Colors.pinkAccent,
              onValueChanged: _updateEmotion,
            ),
          ],
        ),
      ),
    );
  }
}
