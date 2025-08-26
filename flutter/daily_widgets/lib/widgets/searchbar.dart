import 'package:flutter/material.dart';

class SearchBarTest extends StatefulWidget {
  const SearchBarTest({super.key});

  @override
  State<SearchBarTest> createState() => _SearchBarTestState();
}

class _SearchBarTestState extends State<SearchBarTest> {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: Text('SearchBar & SearchAnchor')),
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const Text('Seba'),
            SearchBar(
              leading: const Icon(Icons.search),
              hintText: 'Search',
              backgroundColor: WidgetStateProperty.all(Colors.blue),
              shadowColor: WidgetStateProperty.all(Colors.black),
              elevation: WidgetStateProperty.all(4.0),
              shape: WidgetStateProperty.all(
                RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(20.0)
                )
              ),
              padding: WidgetStateProperty.all(
                EdgeInsets.symmetric(horizontal: 16.0)
              ),
              onSubmitted: (value) {
                print('Search sumbitted: $value');
              },
            ),
            const Text('Some space'),
            SearchAnchor.bar(
              barLeading: Icon(Icons.search),
              barHintText: 'Search',
              barBackgroundColor: WidgetStateProperty.all(Colors.blue),
              barElevation: WidgetStateProperty.all(4.0),
              barShape: WidgetStateProperty.all(
                RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(20.0)
                )
              ),
              barPadding: WidgetStateProperty.all(
                EdgeInsets.symmetric(horizontal: 16.0)
              ),
              suggestionsBuilder: (BuildContext context, SearchController controller) {
                final String input = controller.value.text;
                return [
                  ...List.generate(
                    5,
                    (int index) {
                      final String item = 'Something $index, $input';
                      return ListTile(
                        title: Text(item),
                        onTap: () {
                          debugPrint('You have just selected $item');
                          controller.closeView(item);
                        },
                      );
                    }
                  ),
                ];
              },
            )
          ],
        ),
      ),
    );
  }
}