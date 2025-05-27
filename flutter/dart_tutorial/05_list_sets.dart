void main() {
  final List<int> scores = [50, 60, 75, 89];

  scores[0] = 25;
  print(scores[0]);

  // lists are MUTABLE!
  scores.add(12);
  print(scores);
  
  scores.shuffle();
  print(scores);
}
