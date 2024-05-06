import 'dart:convert' as convert;

import 'package:http/http.dart' as http;

Future<Post> fetchPost() async {
  var url = Uri.https('jsonplaceholder.typicode.com', '/posts/2');
  final response = await http.get(url);
  if (response.statusCode == 200) {
    Map<String, dynamic> data = convert.jsonDecode(response.body);
    return Post(data["title"], data["userId"]);
  }
  return Post("", 0);
}

class Post {
  String title;
  int userId;

  Post(this.title, this.userId);
}