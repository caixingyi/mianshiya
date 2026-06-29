"use client";
import { AntdRegistry } from "@ant-design/nextjs-registry";
import BasicLayout from "@/layouts/BasicLayout";
import React, { useCallback, useEffect } from "react";
import { Provider, useDispatch } from "react-redux";
import { ConfigProvider } from "antd";
import store, { AppDispatch } from "@/stores";
import { getLoginUserUsingGet } from "@/api/userController";
import AccessLayout from "@/access/AccessLayout";
import { setLoginUser } from "@/stores/loginUser";
import { DEFAULT_USER } from "@/constants/user";
import "./globals.css";

/**
 * 全局初始化逻辑
 * @param children
 * @constructor
 */
const InitLayout: React.FC<
  Readonly<{
    children: React.ReactNode;
  }>
> = ({ children }) => {
  const dispatch = useDispatch<AppDispatch>();
  // 初始化全局用户状态
  const doInitLoginUser = useCallback(async () => {
    if (!localStorage.getItem("token")) {
      dispatch(setLoginUser(DEFAULT_USER));
      return;
    }
    try {
      const res = await getLoginUserUsingGet();
      if (res.data) {
        // 更新全局用户状态
        dispatch(setLoginUser(res.data));
      } else {
        dispatch(setLoginUser(DEFAULT_USER));
      }
    } catch {
      localStorage.removeItem("token");
      dispatch(setLoginUser(DEFAULT_USER));
    }
  }, [dispatch]);

  // 只执行一次
  useEffect(() => {
    doInitLoginUser();
  }, [doInitLoginUser]);
  return children;
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="zh">
      <body>
        <AntdRegistry>
          <ConfigProvider
            theme={{
              token: {
                colorPrimary: "#2563eb",
                borderRadius: 10,
                fontFamily: '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif',
              },
              components: {
                Card: {
                  borderRadiusLG: 14,
                },
                Button: {
                  borderRadius: 8,
                },
                Input: {
                  borderRadius: 8,
                },
              },
            }}
          >
            <Provider store={store}>
              <InitLayout>
                <BasicLayout>
                  <AccessLayout>{children}</AccessLayout>
                </BasicLayout>
              </InitLayout>
            </Provider>
          </ConfigProvider>
        </AntdRegistry>
      </body>
    </html>
  );
}
