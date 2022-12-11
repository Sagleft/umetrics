function loadGeoData(systemPeers) {
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

  // Create series for background fill
  // https://www.amcharts.com/docs/v5/charts/map-chart/map-polygon-series/#Background_polygon
  var backgroundSeries = chart.series.push(am5map.MapPolygonSeries.new(root, {}));
  backgroundSeries.mapPolygons.template.setAll({
    fill: root.interfaceColors.get("alternativeBackground"),
    fillOpacity: 0.15,
    strokeWidth: 0.2,
    stroke: root.interfaceColors.get("background")
  });

  // Add background polygon
  // https://www.amcharts.com/docs/v5/charts/map-chart/map-polygon-series/#Background_polygon
  backgroundSeries.data.push({
    geometry: am5map.getGeoRectangle(90, 180, -90, -180)
  });

  // Create main polygon series for countries
  // https://www.amcharts.com/docs/v5/charts/map-chart/map-polygon-series/
  var polygonSeries = chart.series.push(
    am5map.MapPolygonSeries.new(root, {
      geoJSON: am5geodata_worldLow
    })
  );

  // Create line series for trajectory lines
  // https://www.amcharts.com/docs/v5/charts/map-chart/map-line-series/
  var lineSeries = chart.series.push(am5map.MapLineSeries.new(root, {}));
  lineSeries.mapLines.template.setAll({
    stroke: root.interfaceColors.get("alternativeBackground"),
    strokeOpacity: 0.1
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

  // Create point series for markers
  // https://www.amcharts.com/docs/v5/charts/map-chart/map-point-series/
  var pointSeries = chart.series.push(am5map.MapPointSeries.new(root, {}));

  pointSeries.bullets.push(function () {
    var circle = am5.Circle.new(root, {
      radius: 4,
      tooltipY: 0,
      fill: am5.color(0x35c875),
      stroke: root.interfaceColors.get("background"),
      strokeWidth: 2,
      strokeOpacity: 0.75,
      tooltipText: "{title}"
    });

    return am5.Bullet.new(root, {
      sprite: circle
    });
  });

  for (var i = 0; i < systemPeers.length; i++) {
    var peer = systemPeers[i];
    addCity(peer.lon, peer.lat, peer.city);
  }

  function addCity(lon, lat, title) {
    pointSeries.data.push({
      geometry: { type: "Point", coordinates: [lon, lat] },
      title: title
    });
  }

  // Make stuff animate on load
  chart.appear(1000, 100);

}

$( document ).ready(function() {

  $.getJSON( "geoData.json", function( parsedData ) {
    //let parsedData = JSON.parse(data);
    loadGeoData(parsedData);
  });

});
