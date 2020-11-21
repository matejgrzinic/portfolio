let myChart;
let lastTimeframe = "day";

$("#day").click(function () {
  lastTimeframe = "day";
  updateGraph("day");
});

$("#week").click(function () {
  lastTimeframe = "week";
  updateGraph("week");
});

$("#month").click(function () {
  lastTimeframe = "month";
  updateGraph("month");
});

$("#all").click(function () {
  lastTimeframe = "all";
  updateGraph("all");
});

function updateGraph(timeframe) {
  $.ajax({
    url: "/api/v1/timeframe/" + timeframe,
    success: function (result) {
      let r = jQuery.parseJSON(result);
      setupGraph(r);

      let percentChange =
        (parseFloat(r.Value[0]) / parseFloat(r.Value[r.Value.length - 1]) - 1) *
        100;
      $("#title-change").text(percentChange.toFixed(2));
    },
  });
}

function updateUsername() {
  $.ajax({
    url: "/api/v1/username",
    success: function (result) {
      let r = jQuery.parseJSON(result);
      $("#title-username").text(r);
    },
  });
}

function updateNetworth() {
  $.ajax({
    url: "/api/v1/networth",
    success: function (result) {
      let r = jQuery.parseJSON(result);
      $("#title-networth").text(r);
    },
  });
}

function clearChart() {
  try {
    myChart.destroy();
  } catch (err) {}
}

function updateAll() {
  updateGraph(lastTimeframe);
  updateDataTable();
  updateUsername();
  updateNetworth();
  updateTransactionData("gain");
  updateTransactionData("loss");
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

function updateDataTable() {
  $.ajax({
    url: "/api/v1/table/portfolio",
    success: function (result) {
      let tableData = jQuery.parseJSON(result);
      setupDataTable(tableData);
    },
  });
}

function setupDataTable(tableData) {
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
}

function updateTransactionData(id) {
  $.ajax({
    url: "/api/v1/table/" + id,
    success: function (result) {
      let tableData = jQuery.parseJSON(result);
      setupTransactionTable(tableData, id);
    },
  });
}

function setupTransactionTable(tableData, id) {
  let table = document.getElementById(id);
  let headerData = ["Symbol", "Amount", "Value", "Date"];
  $("#" + id + " tr").remove();
  generateTableHead(table, headerData);
  generateTable(table, tableData);
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

function resetTransactionModal() {
  $("#transaction-type").val("default");
  $("#transaction-currency-type").val("default");
  $("#transaction-currency").val("default");
  $("#transaction-amount").val("");
  $("#transaction-description").val("");
}

function resetTradeModal() {
  $("#trade-sell-currency-type").val("default");
  $("#trade-sell-currency").val("default");
  $("#trade-buy-currency-type").val("default");
  $("#trade-buy-currency").val("default");
  $("#trade-sell-amount").val("");
  $("#trade-buy-amount").val("");
  $("#trade-description").val("");
}

function generateTable(table, data) {
  percentAdd = ["HourChange", "DayChange", "WeekChange", "MonthChange"];
  for (let element of data) {
    let row = table.insertRow();
    for (key in element) {
      let cell = row.insertCell();
      let text = element[key];
      if (percentAdd.includes(key)) {
        let span = document.createElement("span");
        if (text[0] == "-") {
          span.style.color = "red";
        } else {
          span.style.color = "green";
        }
        cell.appendChild(span);
        span.appendChild(document.createTextNode(text));
        continue;
      }
      cell.appendChild(document.createTextNode(text));
    }
  }
}

$("#transaction-type").change(function () {
  $("#transaction-currency option[value!='default']").remove();
  addOptions($("#transaction-currency"));
});

$("#transaction-currency-type").change(function () {
  $("#transaction-currency option[value!='default']").remove();
  addOptions($("#transaction-currency"));
});

function addOptions(el) {
  $.ajax({
    url:
      "/api/v1/currencies/" +
      $("#transaction-type").val() +
      "/" +
      $("#transaction-currency-type").val(),
    success: function (result) {
      let r = jQuery.parseJSON(result);
      for (o of r) {
        el.append(new Option(o, o));
      }
    },
  });
}

$("#transaction-submit").click(function () {
  $("#transaction-success").hide();
  $("#transaction-error").hide();

  let transactionData = {
    type: $("#transaction-type").val(),
    "currency-type": $("#transaction-currency-type").val(),
    currency: $("#transaction-currency").val(),
    amount: $("#transaction-amount").val(),
    description: $("#transaction-description").val(),
  };
  console.log(transactionData);

  $.ajax({
    type: "POST",
    url: "/api/v1/transaction",
    data: JSON.stringify(transactionData),
    success: function (result) {
      let r = JSON.parse(result);
      if (r.Status == "OK") {
        $("#transaction-success").show();
        $("#transaction-success").text(r.Message);
      } else {
        $("#transaction-error").show();
        $("#transaction-error").text(r.Message);
      }
      updateAll();
      resetTransactionModal();
    },
    contentType: "json",
  });
});

$("#transaction-new").click(function () {
  $("#transaction-success").hide();
  $("#transaction-error").hide();

  $("#transactionModal").modal("show");
});

function addOptionsTrade(el, type, curType) {
  $.ajax({
    url: "/api/v1/currencies/" + type + "/" + curType,
    success: function (result) {
      let r = jQuery.parseJSON(result);
      for (o of r) {
        el.append(new Option(o, o));
      }
    },
  });
}

$("#trade-sell-currency-type").change(function () {
  $("#trade-sell-currency option[value!='default']").remove();
  addOptionsTrade(
    $("#trade-sell-currency"),
    "loss",
    $("#trade-sell-currency-type").val()
  );
});

$("#trade-buy-currency-type").change(function () {
  $("#trade-buy-currency option[value!='default']").remove();
  addOptionsTrade(
    $("#trade-buy-currency"),
    "gain",
    $("#trade-buy-currency-type").val()
  );
});

$("#trade-submit").click(function () {
  $("#trade-success").hide();
  $("#trade-error").hide();

  let tradeData = {
    "sell-type": $("#trade-sell-currency-type").val(),
    "sell-currency": $("#trade-sell-currency").val(),
    "sell-amount": $("#trade-sell-amount").val(),

    "buy-type": $("#trade-buy-currency-type").val(),
    "buy-currency": $("#trade-buy-currency").val(),
    "buy-amount": $("#trade-buy-amount").val(),

    description: $("#trade-description").val(),
  };

  $.ajax({
    type: "POST",
    url: "/api/v1/trade",
    data: JSON.stringify(tradeData),
    success: function (result) {
      let r = JSON.parse(result);
      if (r.Status == "OK") {
        $("#trade-success").show();
        $("#trade-success").text(r.Message);
      } else {
        $("#trade-error").show();
        $("#trade-error").text(r.Message);
      }
      updateAll();
      resetTradeModal();
    },
    contentType: "json",
  });
});

$("#trade-new").click(function () {
  $("#trade-success").hide();
  $("#trade-error").hide();

  $("#tradeModal").modal("show");
});
