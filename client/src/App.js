//import React, { useEffect, useState } from "react";
import Balance from "./components/Balance";
import Chart from "./components/Chart";
import Container from "react-bootstrap/Container";

function App() {
  // Chart fetch
  const fetchChart = async (timeframe) => {
    const res = await fetch("http://localhost:10000/api/timeline/" + timeframe);
    const data = await res.json();
    return data;
  };

  // Balance fetch
  const fetchBalance = async () => {
    const res = await fetch("http://localhost:10000/api/balance");
    const data = await res.json();
    return data;
  };

  return (
    <div className="App">
      <Container>
        <Chart fetchChart={fetchChart} />
        <Balance fetchBalance={fetchBalance} />
      </Container>
    </div>
  );
}

export default App;
