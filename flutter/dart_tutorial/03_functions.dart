void main() {
  final greeting = greet("Sebastian", 55);
  print(greeting);

  final greeting2 = greet2(age: 18);
  print(greeting2);
}

String greet(String name, int age) {
  return "Hi, my name is $name and I am $age";
}

String greet2({String? name, required int age}) {
  return "Hi, my name is $name and I am $age";
}