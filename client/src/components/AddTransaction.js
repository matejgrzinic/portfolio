import React, { useState } from "react";
import Button from "react-bootstrap/Button";
import Col from "react-bootstrap/esm/Col";
import Form from "react-bootstrap/Form";
import InputSelectCurrency from "./InputSelectCurrency";
import BootstrapSwitchButton from "bootstrap-switch-button-react";

export const AddTransaction = () => {
  const timeToISO = function () {
    const fixDigit = (val) => {
      return val.toString().length === 1 ? "0" + val : val;
    };
    const now = new Date();
    return `${now.getFullYear()}-${fixDigit(now.getMonth() + 1)}-${fixDigit(
      now.getDate()
    )}T${now.getHours()}:${now.getMinutes()}`;
  };

  const [gainCurrency, setGainCurrency] = useState({
    type: "",
    symbol: [],
    amount: 0.0,
  });

  const [loseCurrency, setLoseCurrency] = useState({
    type: "",
    symbol: [],
    amount: 0.0,
  });

  const [displayed, setDisplayed] = useState(false);
  const [time, setTime] = useState(timeToISO());
  const [currentTimeSwitch, setCurrentTimeSwitch] = useState(true);
  const [isTrade, setIsTrade] = useState(false);

  const insertTransaction = function () {
    let data = {
      gains: [
        {
          type: "cryptocurrency",
          symbol: "BTC",
          amount: gainCurrency.amount,
        },
      ],
      losses: [],
      time: currentTimeSwitch
        ? ~~(Date.now() / 1000)
        : ~~(Date.parse(time) / 1000),
    };

    const makeFetch = async () => {
      let resp = await fetch("/api/transaction/add", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
      });
      console.log(resp);
    };

    console.log(data);
    //makeFetch();
  };

  const contentStyle = {
    display: "block",
    border: "1px solid black",

    overflow: "hidden",
    height: displayed ? "auto" : "0",
    opacity: displayed ? "1" : "0",
    transition: "visibility 0s, opacity 0.5s linear",
  };

  const isTradeStyle = {
    paddingLeft: 20,
    paddingRight: 10,
  };

  const loseStyle = {
    height: isTrade ? "auto" : "0",
    opacity: isTrade ? "1" : "0",
    transition: "visibility 0.1s, opacity 0.1s linear",
  };

  return (
    <div>
      <Button onClick={() => setDisplayed(!displayed)}>+</Button>
      <div style={contentStyle}>
        <Form.Group>
          <Form.Row>
            <InputSelectCurrency
              obj={gainCurrency}
              updateFunc={setGainCurrency}
            />
          </Form.Row>
          <Form.Row>
            <Col>
              <Form.Control
                type="datetime-local"
                placeholder="Amount"
                value={time}
                onChange={(e) => setTime(e.target.value)}
                disabled={currentTimeSwitch ? true : false}
              />
            </Col>
            <Col>
              <BootstrapSwitchButton
                checked={currentTimeSwitch}
                onlabel="Yes"
                offlabel="No"
                onChange={() => {
                  setCurrentTimeSwitch(!currentTimeSwitch);
                  console.log(time);
                }}
              />{" "}
              <Form.Label>Use current time.</Form.Label>
            </Col>
          </Form.Row>
          <Form.Row>
            <Col>
              <Form.Label style={isTradeStyle}>
                Is this transaction a trade?
              </Form.Label>
              <BootstrapSwitchButton
                checked={isTrade}
                onlabel="Yes"
                offlabel="No"
                onChange={() => {
                  setIsTrade(!isTrade);
                }}
              />{" "}
            </Col>
          </Form.Row>
          <Form.Row style={loseStyle}>
            <InputSelectCurrency
              obj={loseCurrency}
              updateFunc={setLoseCurrency}
            />
          </Form.Row>
        </Form.Group>
      </div>
    </div>
  );
};

export default AddTransaction;
