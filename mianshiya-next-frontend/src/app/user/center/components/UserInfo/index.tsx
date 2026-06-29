import React from "react";
import Paragraph from "antd/es/typography/Paragraph";
import "./index.css";

interface Props {
  user: API.LoginUserVO;
}

const UserInfo = (props: Props) => {
  const { user } = props;

  return (
    <div className="user-info">
      <div style={{ textAlign: "left" }}>
        <Paragraph type="secondary">手机号：{user.phoneNumber || "未填写"}</Paragraph>
        <Paragraph type="secondary">邮箱：{user.email || "未填写"}</Paragraph>
        <Paragraph type="secondary">年级：{user.grade || "未填写"}</Paragraph>
        <Paragraph type="secondary">工作经验：{user.workExperience || "未填写"}</Paragraph>
        <Paragraph type="secondary">擅长方向：{user.expertiseDirection || "未填写"}</Paragraph>
      </div>
    </div>
  );
};

export default UserInfo;
