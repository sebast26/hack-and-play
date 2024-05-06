import 'package:fetching/fetching.dart' as fetching;

void main(List<String> arguments) async {
  final post = await fetching.fetchPost();
  print(post.title);
  print(post.userId);
}
