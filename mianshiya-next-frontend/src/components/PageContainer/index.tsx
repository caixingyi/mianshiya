import React from "react";
import "./index.css";

interface Props {
  children: React.ReactNode;
  narrow?: boolean;
  wide?: boolean;
  className?: string;
}

const PageContainer = ({ children, narrow, wide, className = "" }: Props) => {
  const cls = [
    "page-container",
    narrow ? "page-container-narrow" : "",
    wide ? "page-container-wide" : "",
    className,
  ].filter(Boolean).join(" ");

  return <main className={cls}>{children}</main>;
};

export default PageContainer;
