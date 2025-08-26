import 'package:flutter/material.dart';
import 'package:video_player/video_player.dart';

class VideoAppTest extends StatefulWidget {
  const VideoAppTest({super.key});

  @override
  State<VideoAppTest> createState() => _VideoAppTestState();
}

class _VideoAppTestState extends State<VideoAppTest> {
  late VideoPlayerController _controller;

  @override
  void initState() {
    super.initState();
    _controller =
        VideoPlayerController.networkUrl(
            Uri.parse(
              'https://flutter.github.io/assets-for-api-docs/assets/videos/bee.mp4',
            ),
          )
          ..initialize().then((_) {
            setState(() {});
          });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Center(
        child: _controller.value.isInitialized
            ? AspectRatio(
                aspectRatio: _controller.value.aspectRatio,
                child: VideoPlayer(_controller),
              )
            : Container(),
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () {
          setState(() {
            _controller.value.isPlaying
                ? _controller.pause()
                : _controller.play();
          });
        },
        child: Icon(
          _controller.value.isPlaying ? Icons.pause : Icons.play_arrow,
        ),
      ),
    );
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }
}
