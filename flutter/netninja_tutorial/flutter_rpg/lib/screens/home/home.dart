import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter/widgets.dart';
import 'package:flutter_rpg/shared/styled_button.dart';
import 'package:flutter_rpg/shared/styled_text.dart';

class Home extends StatefulWidget {
  const Home({super.key});

  @override
  State<Home> createState() => _HomeState();
}

class _HomeState extends State<Home> {

  List characters = ['mario', 'luigi', 'peach', 'toad', 'bowser', 'koopa'];

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const StyledTitle('Your Characters'),
        centerTitle: true,
      ),
      body: Container(
        padding: const EdgeInsets.all(16),
        child: Column(
          children: [
            Expanded(
              child: ListView.builder(
                itemCount: characters.length,
                itemBuilder: (_, index) {
                  return Container(
                    color: Colors.grey[800],
                    padding: const EdgeInsets.all(40),
                    margin: const EdgeInsets.only(bottom: 20),
                    child: Text(characters[index]),
                  );
                },
              ),
            ),
            StyledButton(
              onPressed: () {},
              child: const StyledHeading('Create New'),
            )
          ],
        ),
      ),
    );
  }
}