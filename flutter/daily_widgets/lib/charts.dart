import 'package:fl_chart/fl_chart.dart';
import 'package:flutter/material.dart';

//
// You can find more under: https://github.com/imaNNeo/fl_chart/blob/main/repo_files/documentations/pie_chart.md
//

class ChartsTest extends StatefulWidget {
  const ChartsTest({super.key});

  @override
  State<ChartsTest> createState() => _ChartsTestState();
}

class _ChartsTestState extends State<ChartsTest> {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('fl_chart test')),
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const Text('Seba'),
            Expanded(
              child: PieChart(
                PieChartData(
                  centerSpaceRadius: double.infinity,
                  sections: [
                    PieChartSectionData(
                      value: 50,
                      title: 'Housing',
                      showTitle: true,
                      radius: 170,
                      color: Colors.purple,
                      titlePositionPercentageOffset: 0.5
                    ),
                    PieChartSectionData(
                      value: 200,
                      title: 'Internet',
                      showTitle: true,
                      radius: 20,
                      color: Colors.red,
                    ),
                    PieChartSectionData(
                      value: 200,
                      title: 'Food',
                      showTitle: true,
                      radius: 30,
                      color: Colors.yellow,
                    ),
                  ],
                ),
                duration: Duration(milliseconds: 150),
                curve: Curves.linear,
              ),
            ),
          ],
        ),
      ),
    );
  }
}
