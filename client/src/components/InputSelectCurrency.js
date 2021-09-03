import React, { useState } from "react";
import { Typeahead } from "react-bootstrap-typeahead";
import Form from "react-bootstrap/Form";
import Col from "react-bootstrap/esm/Col";
import "react-bootstrap-typeahead/css/Typeahead.css";

export const InputSelectCurrency = ({ obj, updateFunc }) => {
  const [options, setOptions] = useState(["rok", "matej"]);

  const updateOptions = async () => {
    const fetchOptions = async () => {
      const res = await fetch("http://localhost:10000/api/balance");
      const data = await res.json();
      setOptions(data.data);
    };
    fetchOptions();
  };

  return (
    <>
      <Col>
        <Form.Control
          as="select"
          custom
          onChange={(e) => updateFunc({ ...obj, type: e.target.value })}
        >
          <option value="default">Select Currency type</option>
          <option value="fiat">Fiat</option>
          <option value="cryptocurrency">Cryptocurrency</option>
          // TODO add stock
        </Form.Control>
      </Col>
      <Typeahead
        clearButton
        id="typeahead-currency"
        onChange={(e) =>
          updateFunc({
            ...obj,
            symbol: e,
          })
        }
        options={options}
        placeholder="Choose a currency..."
        selected={obj.symbol}
      />
      <Col>
        <Form.Control
          type="number"
          placeholder="Amount"
          value={obj.amount ? obj.amount : ""}
          onChange={(e) => updateFunc({ ...obj, amount: e.target.value })}
        />
      </Col>
    </>
  );
};

export default InputSelectCurrency;
