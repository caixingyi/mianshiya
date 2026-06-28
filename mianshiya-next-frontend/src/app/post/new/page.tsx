"use client";
import { addPostUsingPost } from "@/api/postController";
import { Card, Form, Input, Select, Button, message } from "antd";
import { useRouter } from "next/navigation";
import { useState } from "react";
import "./index.css";

export default function NewPostPage() {
  const [loading, setLoading] = useState(false);
  const router = useRouter();

  const onFinish = async (values: any) => {
    setLoading(true);
    try {
      const res = await addPostUsingPost({
        title: values.title,
        content: values.content,
        tags: values.tags || [],
      });
      const r = res as any;
      if (r.code === 0) {
        message.success("发布成功");
        router.push("/posts");
      } else {
        message.error(r.message || "发布失败");
      }
    } catch {
      message.error("发布失败，请先登录");
    }
    setLoading(false);
  };

  return (
    <div id="newPostPage">
      <Card title="发布帖子">
        <Form layout="vertical" onFinish={onFinish}>
          <Form.Item
            name="title"
            label="标题"
            rules={[{ required: true, message: "请输入标题" }]}
          >
            <Input placeholder="输入帖子标题" maxLength={100} />
          </Form.Item>
          <Form.Item
            name="content"
            label="内容（支持 Markdown）"
            rules={[{ required: true, message: "请输入内容" }]}
          >
            <Input.TextArea rows={12} placeholder="用 Markdown 格式写内容..." />
          </Form.Item>
          <Form.Item name="tags" label="标签">
            <Select mode="tags" placeholder="输入标签后按回车" />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading}>
              发布
            </Button>
            <Button style={{ marginLeft: 12 }} onClick={() => router.back()}>
              取消
            </Button>
          </Form.Item>
        </Form>
      </Card>
    </div>
  );
}
