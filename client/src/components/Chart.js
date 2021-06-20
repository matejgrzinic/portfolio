import React, { useEffect, useState } from "react";
import "chartjs-adapter-moment";
import { Line } from "react-chartjs-2";
import TimelineButtons from "./TimelineButtons";

export const Chart = ({ fetchChart }) => {
  const [timeline, setTimeline] = useState([]);
  const [timeframe, setTimeframe] = useState("hour");

  useEffect(() => {
    const updateChart = async () => {
      let data = await fetchChart(timeframe);
      setTimeline(data.data);
    };
    updateChart();
  }, [timeframe]);

  const updateChart = (e) => {
    setTimeframe(e);
  };

  const data = {
    labels: timeline.map((x) => x.time * 1000),
    datasets: [
      {
        label: "Net Worth",
        fill: false,
        lineTension: 0.1,
        backgroundColor: "rgba(75,192,192,0.4)",
        borderColor: "rgba(75,192,192,1)",
        borderCapStyle: "butt",
        borderDash: [],
        borderDashOffset: 0.0,
        borderJoinStyle: "miter",
        pointBorderColor: "rgba(75,192,192,1)",
        pointBackgroundColor: "#fff",
        pointBorderWidth: 1,
        pointHoverRadius: 5,
        pointHoverBackgroundColor: "rgba(75,192,192,1)",
        pointHoverBorderColor: "rgba(220,220,220,1)",
        pointHoverBorderWidth: 2,
        pointRadius: 1,
        pointHitRadius: 20,
        data: timeline.map((x) => x.value.toFixed(2)),
      },
    ],
  };
  // day -> hour, week -> day, month->week year->month all->quarter
  const options = {
    //maintainAspectRatio: false,
    legend: {
      display: false,
    },
    scales: {
      y: {
        ticks: {
          fontSize: 11,
          fontColor: "#A3A3A3",
        },
      },
      x: {
        type: "time",
        time: {
          tooltipFormat: "YYYY/MM/DD HH:mm",
          unit: "hour",
          //unitStepSize: 2,
          displayFormats: {
            minute: "h:mm:ss",
            hour: "D MMM HH:mm",
            day: "D MMM",
            week: "ll",
            month: "MMM YYYY",
            quarter: "[Q]Q - YYYY",
            year: "YYYY",
          },
        },
        ticks: {
          count: 10,
          maxTicksLimit: 3,
          maxRotation: 0,
          //source: "data",
        },
      },
    },
  };

  return (
    <>
      <TimelineButtons onChange={updateChart} />
      <Line data={data} options={options} />
    </>
  );
};

export default Chart;
