import React, { useState } from "react";
import Balance from "./components/Balance";
import Chart from "./components/Chart";
import Container from "react-bootstrap/Container";
import Sidebar from "./components/Sidebar";
import AddTransaction from "./components/AddTransaction";

function App() {
  const [showSideBar, setShowSideBar] = useState(true);

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

  const mainPositioning = {
    paddingLeft: showSideBar ? "300px" : "100px",
    width: "100%",
    paddingRight: "50px",
    transitionDuration: "0.5s",
  };

  return (
    <div className="App">
      <Sidebar showSideBar={showSideBar} click={setShowSideBar} />
      <Container style={mainPositioning} fluid>
        <AddTransaction />
        <Chart fetchChart={fetchChart} />
        <Balance fetchBalance={fetchBalance} />
      </Container>
    </div>
  );
}

export default App;
