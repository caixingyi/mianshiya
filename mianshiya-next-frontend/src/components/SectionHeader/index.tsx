import React from "react";
import Title from "antd/es/typography/Title";
import "./index.css";

interface Props {
  title: React.ReactNode;
  description?: React.ReactNode;
  extra?: React.ReactNode;
}

const SectionHeader = ({ title, description, extra }: Props) => {
  return (
    <div className="section-header">
      <div>
        <Title level={3} className="section-header-title">
          {title}
        </Title>
        {description && <div className="section-header-desc">{description}</div>}
      </div>
      {extra && <div className="section-header-extra">{extra}</div>}
    </div>
  );
};

export default SectionHeader;
