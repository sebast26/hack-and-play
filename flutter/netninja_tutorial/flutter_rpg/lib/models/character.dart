import 'package:flutter_rpg/models/stats.dart';

class Character with Stats{

  Character({
    required this.name, required this.slogan, required this.id
  });

  final String name;
  final String slogan;
  final String id;
  bool _isFav = false;

  // getters
  bool get isFav => _isFav;

  void toggleIsFav() {
    _isFav = !_isFav;
  }
}