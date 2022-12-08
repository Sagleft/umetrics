am5.ready(function() {

// Create root element
// https://www.amcharts.com/docs/v5/getting-started/#Root_element
var root = am5.Root.new("globe-indicator");


// Set themes
// https://www.amcharts.com/docs/v5/concepts/themes/
root.setThemes([
  am5themes_Animated.new(root)
]);


// Create the map chart
// https://www.amcharts.com/docs/v5/charts/map-chart/
var chart = root.container.children.push(am5map.MapChart.new(root, {
  panX: "rotateX",
  panY: "rotateY",
  projection: am5map.geoOrthographic()
}));


// Create series for background fill
// https://www.amcharts.com/docs/v5/charts/map-chart/map-polygon-series/#Background_polygon
var backgroundSeries = chart.series.push(
  am5map.MapPolygonSeries.new(root, {})
);
backgroundSeries.mapPolygons.template.setAll({
  fill: root.interfaceColors.get("alternativeBackground"),
  fillOpacity: 0.1,
  strokeOpacity: 0
});
backgroundSeries.data.push({
  geometry:
    am5map.getGeoRectangle(90, 180, -90, -180)
});


// Create main polygon series for countries
// https://www.amcharts.com/docs/v5/charts/map-chart/map-polygon-series/
var polygonSeries = chart.series.push(am5map.MapPolygonSeries.new(root, {}));
polygonSeries.data.push({
  geometry: {
    type: "Polygon",
    coordinates: [
      [
        [26.5936, 55.6676],
        [26.175, 55.0033],
        [25.8594, 54.9192],
        [25.5473, 54.3317],
        [24.7683, 53.9746],
        [23.4845, 53.9398],
        [23.37, 54.2005],
        [22.7663, 54.3568],
        [22.8311, 54.8384],
        [21.2358, 55.2641],
        [21.0462, 56.07],
        [22.0845, 56.4067],
        [24.1206, 56.2642],
        [24.9032, 56.3982],
        [26.5936, 55.6676]
      ]
    ]
  }
});

polygonSeries.mapPolygons.template.setAll({
  fill: root.interfaceColors.get("alternativeBackground"),
  fillOpacity: 0.15,
  strokeWidth: 0.5,
  stroke: root.interfaceColors.get("background")
});


// Create polygon series for projected circles
var circleSeries = chart.series.push(am5map.MapPolygonSeries.new(root, {}));
circleSeries.mapPolygons.template.setAll({
  templateField: "polygonTemplate",
  tooltipText: "{name}:{value}"
});

// Define data
var colors = am5.ColorSet.new(root, {});

var data = [
  { "id": "AF", "name": "Afghanistan", "value": 32358260, polygonTemplate: { fill: colors.getIndex(0) } },
  { "id": "AL", "name": "Albania", "value": 3215988, polygonTemplate: { fill: colors.getIndex(8) } },
  { "id": "DZ", "name": "Algeria", "value": 35980193, polygonTemplate: { fill: colors.getIndex(2) } },
  { "id": "AO", "name": "Angola", "value": 19618432, polygonTemplate: { fill: colors.getIndex(2) } },
  { "id": "AR", "name": "Argentina", "value": 40764561, polygonTemplate: { fill: colors.getIndex(3) } },
  { "id": "AM", "name": "Armenia", "value": 3100236, polygonTemplate: { fill: colors.getIndex(8) } },
];

var valueLow = Infinity;
var valueHigh = -Infinity;

for (var i = 0; i < data.length; i++) {
  var value = data[i].value;
  if (value < valueLow) {
    valueLow = value;
  }
  if (value > valueHigh) {
    valueHigh = value;
  }
}

// radius in degrees
var minRadius = 0.5;
var maxRadius = 5;

// Create circles when data for countries is fully loaded.
/*polygonSeries.events.on("datavalidated", function () {
  circleSeries.data.clear();

  for (var i = 0; i < data.length; i++) {
    var dataContext = data[i];
    var countryDataItem = polygonSeries.getDataItemById(dataContext.id);
    var countryPolygon = countryDataItem.get("mapPolygon");

    var value = dataContext.value;

    var radius = minRadius + maxRadius * (value - valueLow) / (valueHigh - valueLow);

    if (countryPolygon) {
      var geometry = am5map.getGeoCircle(countryPolygon.visualCentroid(), radius);
      circleSeries.data.push({
        name: dataContext.name,
        value: dataContext.value,
        polygonTemplate: dataContext.polygonTemplate,
        geometry: geometry
      });
    }
  }
})*/


// Make stuff animate on load
chart.appear(1000, 100);

}); // end am5.ready()
