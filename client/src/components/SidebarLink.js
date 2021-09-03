import React, { useState } from "react";

export const SidebarLink = ({ text, showSideBar }) => {
  const myStyle = {
    padding: "6px 8px 6px 10px",
    textDecoration: "none",
    fontSize: "20px",
    color: "#818181",
    display: "block",
    cursor: "pointer",
  };

  return <span style={myStyle}>{showSideBar ? text : "A"}</span>;
};

export default SidebarLink;
