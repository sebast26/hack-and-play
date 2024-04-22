void main() {
  var noodles = MenuItem("veg noodles", 9.99);
  var pizza = Pizza(["mushrooms"], "veg volcano", 15.99);

  print(noodles);
  print(pizza);

  var foods = Collection<MenuItem>('Menu Items', [
    noodles, pizza
  ]);
  var random = foods.randomItem();
  print(random);
}

class MenuItem {
  String title = "pizza";
  double price = 9.99;

  MenuItem(this.title, this.price);

  @override
  String toString() {
    return "$title --> $price";
  }
}

class Pizza extends MenuItem {
  List<String> toppings;

  Pizza(this.toppings, super.title, super.price);

  @override
  String toString() {
    return "pizza $title that cost $price with $toppings";
  }
}

class Collection<T> {
  String name;
  List<T> data;

  Collection(this.name, this.data);

  T randomItem() {
    data.shuffle();
    return data[0];
  }
}