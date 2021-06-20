import React, { useState } from "react";
import Button from "react-bootstrap/Button";

export const TimelineButtons = ({ onChange }) => {
  const onChangeGen = (i) => () => onChange(i);

  return (
    <>
      <Button
        as="input"
        type="button"
        value="Hour"
        onClick={onChangeGen("hour")}
      />{" "}
      <Button
        as="input"
        type="button"
        value="Day"
        onClick={onChangeGen("day")}
      />{" "}
      <Button
        as="input"
        type="button"
        value="Week"
        onClick={onChangeGen("week")}
      />{" "}
      <Button
        as="input"
        type="button"
        value="Month"
        onClick={onChangeGen("month")}
      />{" "}
      <Button
        as="input"
        type="button"
        value="Year"
        onClick={onChangeGen("year")}
      />{" "}
      <Button
        as="input"
        type="button"
        value="All"
        onClick={onChangeGen("all")}
      />{" "}
    </>
  );
};

export default TimelineButtons;
