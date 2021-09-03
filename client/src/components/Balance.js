import React, { useEffect, useState } from "react";
import Table from "react-bootstrap/Table";

export const Balance = ({ fetchBalance }) => {
  const [balance, setBalance] = useState([]);

  useEffect(() => {
    const updateBalance = async () => {
      let data = await fetchBalance();
      setBalance(data.data);
    };
    updateBalance();
  }, []);

  const balances = balance
    .sort((a, b) => (a.value < b.value ? 1 : -1))
    .map((x) => (
      <tr key={x.symbol}>
        <td>{x.symbol}</td>
        <td>{x.amount.toFixed(2)}</td>
        <td>{x.price.toFixed(2)} €</td>
        <td>{x.value.toFixed(2)} €</td>
      </tr>
    ));

  return (
    <Table striped bordered hover>
      <thead>
        <tr>
          <th>Symbol</th>
          <th>Amount</th>
          <th>Price</th>
          <th>Value</th>
        </tr>
      </thead>
      <tbody>{balances}</tbody>
    </Table>
  );
};

export default Balance;
