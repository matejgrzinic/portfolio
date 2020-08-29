let myChart;

$("#day").click(function () {
  $.ajax({
    url: "/api/v1/timeframe/day",
    success: function (result) {
      let r = jQuery.parseJSON(result);
      setupGraph(r);
    },
  });
});

$("#week").click(function () {
  $.ajax({
    url: "/api/v1/timeframe/week",
    success: function (result) {
      let r = jQuery.parseJSON(result);
      setupGraph(r);
    },
  });
});

$("#month").click(function () {
  $.ajax({
    url: "/api/v1/timeframe/month",
    success: function (result) {
      let r = jQuery.parseJSON(result);
      setupGraph(r);
    },
  });
});

$("#all").click(function () {
  $.ajax({
    url: "/api/v1/timeframe/all",
    success: function (result) {
      let r = jQuery.parseJSON(result);
      setupGraph(r);
    },
  });
});

function clearChart() {
  try {
    myChart.destroy();
  } catch (err) {}
}

function setupGraph(data) {
  clearChart();
  const ctx = document.getElementById("myChart").getContext("2d");
  myChart = new Chart(ctx, {
    type: "line",
    data: {
      labels: data.Time,
      datasets: [
        {
          label: "net worth",
          data: data.Value,
          fill: false,
          borderColor: "rgba(255, 99, 132, 1)",
          backgroundColor: "rgba(255, 99, 132, 0.5)",
          borderWidth: 1,
        },
      ],
    },
    options: {
      scales: {
        xAxes: [
          {
            type: "time",
            time: {
              unit: data.TimeUnit,
            },
          },
        ],
      },
    },
  });
}

function setupDataTable() {
  $.ajax({
    url: "/api/v1/timeframe/alldata",
    success: function (result) {
      let tableData = jQuery.parseJSON(result);
      let table = document.getElementById("datatable");
      let headerData = [
        "Currency",
        "Symbol",
        "Amount",
        "Price",
        "Value",
        "H",
        "D",
        "W",
        "M",
      ];
      $("#datatable tr").remove();
      generateTableHead(table, headerData);
      generateTable(table, tableData);
    },
  });
}

function generateTableHead(table, data) {
  let thead = table.createTHead();
  let row = thead.insertRow();
  for (let key of data) {
    let th = document.createElement("th");
    let text = document.createTextNode(key);
    th.appendChild(text);
    row.appendChild(th);
  }
}

function generateTable(table, data) {
  eurAdd = ["Price", "Value"];
  percentAdd = ["HourChange", "DayChange", "WeekChange", "MonthChange"];
  for (let element of data) {
    let row = table.insertRow();
    for (key in element) {
      let cell = row.insertCell();

      let text = element[key];
      if (eurAdd.includes(key)) {
        text += " €";
      } else if (percentAdd.includes(key)) {
        text += " %";
      }

      let inText = document.createTextNode(text);

      cell.appendChild(inText);
    }
  }
}
