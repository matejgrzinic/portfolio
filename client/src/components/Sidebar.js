import React, { useState } from "react";
import Button from "react-bootstrap/Button";
import SidebarLink from "./SidebarLink";

export const Sidebar = ({ showSideBar, click }) => {
  const myStyle = {
    height: "100%",
    width: showSideBar ? "250px" : "50px",
    position: "fixed",
    zindex: 1,
    top: 1,
    left: 1,
    backgroundColor: "#111",
    overflowX: "hidden",
    paddingTop: "20px",
    transitionDuration: "0.5s",
  };

  console.log(showSideBar ? "250px" : "50px");

  const buttonStyle = {
    // position: "absolute",
    // right: "200px",
    // top: "50px",
    // "margin-bottom": "100px",
  };

  return (
    <div style={myStyle}>
      <Button
        variant="outline-primary"
        onClick={() => {
          click(!showSideBar);
        }}
      >
        {showSideBar ? "Hide" : ">"}
      </Button>{" "}
      <SidebarLink text={"Dashboard"} showSideBar={showSideBar} />
      <SidebarLink text={"Portfolio"} showSideBar={showSideBar} />
      <SidebarLink text={"History"} showSideBar={showSideBar} />
    </div>
  );
};

export default Sidebar;
