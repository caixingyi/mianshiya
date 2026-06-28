"use client";
import { getPostVoByIdUsingGet } from "@/api/postController";
import { doThumbUsingPost } from "@/api/postThumbController";
import { doPostFavourUsingPost } from "@/api/postFavourController";
import { Card, Spin, message } from "antd";
import {
  LikeOutlined,
  LikeFilled,
  StarOutlined,
  StarFilled,
} from "@ant-design/icons";
import TagList from "@/components/TagList";
import MdViewer from "@/components/MdViewer";
import { useEffect, useState } from "react";
import "./index.css";

interface Props {
  params: { postId: string };
}

export default function PostPage({ params }: Props) {
  const { postId } = params;
  const [post, setPost] = useState<API.PostVO | null>(null);
  const [loading, setLoading] = useState(true);
  const [hasThumb, setHasThumb] = useState(false);
  const [thumbNum, setThumbNum] = useState(0);
  const [hasFavour, setHasFavour] = useState(false);
  const [favourNum, setFavourNum] = useState(0);

  useEffect(() => {
    getPostVoByIdUsingGet({ id: Number(postId) })
      .then((res: any) => {
        const p = res.data;
        setPost(p as any);
        setHasThumb(p.hasThumb ?? false);
        setThumbNum(p.thumbNum ?? 0);
        setHasFavour(p.hasFavour ?? false);
        setFavourNum(p.favourNum ?? 0);
      })
      .catch(() => message.error("获取帖子失败"))
      .finally(() => setLoading(false));
  }, [postId]);

  const handleThumb = async () => {
    try {
      await doThumbUsingPost({ postId: Number(postId) });
      setHasThumb(!hasThumb);
      setThumbNum(hasThumb ? thumbNum - 1 : thumbNum + 1);
    } catch {
      message.error("请先登录");
    }
  };

  const handleFavour = async () => {
    try {
      await doPostFavourUsingPost({ postId: Number(postId) });
      setHasFavour(!hasFavour);
      setFavourNum(hasFavour ? favourNum - 1 : favourNum + 1);
    } catch {
      message.error("请先登录");
    }
  };

  if (loading) return <Spin style={{ display: "block", margin: "100px auto" }} />;
  if (!post) return <div>获取帖子详情失败，请刷新重试</div>;

  return (
    <div id="postPage">
      <Card>
        <h1 className="post-title">{post.title}</h1>
        <div className="post-meta">
          {post.user?.userName ?? "匿名用户"} · {post.createTime?.slice(0, 10)}
        </div>
        <TagList tagList={post.tagList} />
        <div style={{ marginBottom: 24 }} />
        <MdViewer value={post.content} />
        <div className="post-actions">
          <span className={`action-btn${hasThumb ? " active" : ""}`} onClick={handleThumb}>
            {hasThumb ? <LikeFilled /> : <LikeOutlined />} {thumbNum} 点赞
          </span>
          <span className={`action-btn${hasFavour ? " active" : ""}`} onClick={handleFavour}>
            {hasFavour ? <StarFilled /> : <StarOutlined />} {favourNum} 收藏
          </span>
        </div>
      </Card>
    </div>
  );
}
