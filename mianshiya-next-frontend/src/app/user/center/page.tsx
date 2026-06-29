"use client";
import { Avatar, Card, Col, Row, Segmented, Tag, message } from "antd";
import { useSelector, useDispatch } from "react-redux";
import { RootState, AppDispatch } from "@/stores";
import { setLoginUser } from "@/stores/loginUser";
import Title from "antd/es/typography/Title";
import Paragraph from "antd/es/typography/Paragraph";
import { useState } from "react";
import { CameraOutlined } from "@ant-design/icons";
import CalendarChart from "@/app/user/center/components/CalendarChart";
import "./index.css";
import UserInfo from "@/app/user/center/components/UserInfo";
import UserInfoEditForm from "@/app/user/center/components/UserInfoEditForm";
import { USER_ROLE_ENUM, USER_ROLE_TEXT_MAP } from "@/constants/user";
import { uploadFileUsingPost } from "@/api/fileController";
import { updateMyUserUsingPost } from "@/api/userController";
import dayjs from "dayjs";

export default function UserCenterPage() {
  const loginUser = useSelector((state: RootState) => state.loginUser);
  const dispatch = useDispatch<AppDispatch>();
  const user = loginUser;
  const [activeTabKey, setActiveTabKey] = useState<string>("info");
  const [currentEditState, setCurrentEditState] = useState<string>("查看信息");

  return (
    <div id="userCenterPage" className="max-width-content">
      <Row gutter={[16, 16]}>
        <Col xs={24} md={6}>
          <Card style={{ textAlign: "center" }}>
            <input
              type="file"
              accept="image/*"
              style={{ display: "none" }}
              id="avatar-upload-input"
              onChange={async (e) => {
                const file = e.target.files?.[0];
                if (!file) return;
                try {
                  const res = await uploadFileUsingPost({}, { biz: "user_avatar" }, file);
                  const url = (res as any).data || "";
                  if (url) {
                    await updateMyUserUsingPost({ userAvatar: url } as any);
                    dispatch(setLoginUser({ ...user, userAvatar: url } as any));
                    message.success("头像更新成功");
                  }
                } catch (err: any) {
                  console.error("上传失败:", err);
                  message.error("上传失败: " + (err?.message || err?.toString() || "未知错误"));
                }
                e.target.value = "";
              }}
            />
            <label htmlFor="avatar-upload-input" style={{ cursor: "pointer", display: "inline-block" }}>
              <div style={{ position: "relative", display: "inline-block" }}>
                <Avatar src={user.userAvatar} size={72} />
                <div
                  style={{
                    position: "absolute", bottom: 0, right: 0,
                    background: "rgba(0,0,0,0.45)", borderRadius: "50%",
                    width: 24, height: 24, display: "flex",
                    alignItems: "center", justifyContent: "center",
                  }}
                >
                  <CameraOutlined style={{ color: "#fff", fontSize: 12 }} />
                </div>
              </div>
            </label>
            <div style={{ marginBottom: 16 }} />
            <Card.Meta
              title={
                <Title level={4} style={{ marginBottom: 0 }}>
                  {user.userName}
                </Title>
              }
              description={
                <Paragraph type="secondary">{user.userProfile}</Paragraph>
              }
            />
            <Tag
              color={user.userRole === USER_ROLE_ENUM.ADMIN ? "gold" : "grey"}
            >
              {USER_ROLE_TEXT_MAP[user.userRole]}
            </Tag>
            <Paragraph type="secondary" style={{ marginTop: 8 }}>
              注册日期：{dayjs(user.createTime).format("YYYY-MM-DD")}
            </Paragraph>
            <Paragraph type="secondary" style={{ marginTop: 8 }} copyable={{
              text: user.id
            }}>
              我的 id：{user.id}
            </Paragraph>
          </Card>
        </Col>
        <Col xs={24} md={18}>
          <Card
            tabList={[
              { key: "info", label: "我的信息" },
              { key: "record", label: "刷题记录" },
              { key: "others", label: "其他" },
            ]}
            activeTabKey={activeTabKey}
            onTabChange={(key: string) => { setActiveTabKey(key); }}
          >
            {activeTabKey === "info" && (
              <>
                <Segmented<string>
                  options={["查看信息", "修改信息"]}
                  value={currentEditState}
                  onChange={setCurrentEditState}
                />
                {currentEditState === "查看信息" && <UserInfo user={user} />}
                {currentEditState === "修改信息" && <UserInfoEditForm user={user} />}
              </>
            )}
            {activeTabKey === "record" && <><CalendarChart /></>}
            {activeTabKey === "others" && <>bbb</>}
          </Card>
        </Col>
      </Row>
    </div>
  );
}
