"use client";
import { Card, message, Avatar } from "antd";
import {
  LikeOutlined,
  LikeFilled,
  StarOutlined,
  StarFilled,
  UserOutlined,
  ClockCircleOutlined,
} from "@ant-design/icons";
import Link from "next/link";
import { useState } from "react";
import TagList from "@/components/TagList";
import { doThumbUsingPost } from "@/api/postThumbController";
import { doPostFavourUsingPost } from "@/api/postFavourController";
import "./index.css";

interface Props {
  post: API.PostVO;
}

const PostCard = (props: Props) => {
  const { post } = props;
  const [hasThumb, setHasThumb] = useState(post.hasThumb ?? false);
  const [thumbNum, setThumbNum] = useState(post.thumbNum ?? 0);
  const [hasFavour, setHasFavour] = useState(post.hasFavour ?? false);
  const [favourNum, setFavourNum] = useState(post.favourNum ?? 0);
  const [thumbLoading, setThumbLoading] = useState(false);
  const [favourLoading, setFavourLoading] = useState(false);

  const handleThumb = async () => {
    if (thumbLoading) return;
    setThumbLoading(true);
    try {
      await doThumbUsingPost({ postId: post.id });
      setHasThumb(!hasThumb);
      setThumbNum(hasThumb ? thumbNum - 1 : thumbNum + 1);
    } catch {
      message.error("操作失败");
    }
    setThumbLoading(false);
  };

  const handleFavour = async () => {
    if (favourLoading) return;
    setFavourLoading(true);
    try {
      await doPostFavourUsingPost({ postId: post.id });
      setHasFavour(!hasFavour);
      setFavourNum(hasFavour ? favourNum - 1 : favourNum + 1);
    } catch {
      message.error("操作失败");
    }
    setFavourLoading(false);
  };

  const userName = post.user?.userName || "匿名用户";
  const avatarUrl = post.user?.userAvatar;

  return (
    <div className="post-card">
      <Card hoverable>
        <div className="post-card-header">
          <Link href={`/post/${post.id}`} className="post-title-link">
            {post.title}
          </Link>
        </div>
        <div className="post-card-body">
          {(post.content ?? "").slice(0, 200)}
          {(post.content?.length ?? 0) > 200 ? " ..." : ""}
        </div>
        <div className="post-card-footer">
          <div className="post-card-meta">
            <Avatar
              size={24}
              src={avatarUrl}
              icon={!avatarUrl ? <UserOutlined /> : undefined}
              style={{ flexShrink: 0 }}
            />
            <span className="post-author">{userName}</span>
            <span className="post-time">
              <ClockCircleOutlined style={{ marginRight: 2 }} />
              {post.createTime?.slice(0, 10)}
            </span>
          </div>
          <TagList tagList={post.tagList} />
          <div className="post-actions">
            <span
              className={`action-btn${hasThumb ? " active" : ""}`}
              onClick={handleThumb}
            >
              {hasThumb ? <LikeFilled /> : <LikeOutlined />} {thumbNum}
            </span>
            <span
              className={`action-btn${hasFavour ? " active" : ""}`}
              onClick={handleFavour}
            >
              {hasFavour ? <StarFilled /> : <StarOutlined />} {favourNum}
            </span>
          </div>
        </div>
      </Card>
    </div>
  );
};

export default PostCard;
