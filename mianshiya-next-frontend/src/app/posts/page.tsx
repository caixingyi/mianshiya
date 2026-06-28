import Title from "antd/es/typography/Title";
import { Button } from "antd";
import { PlusOutlined } from "@ant-design/icons";
import Link from "next/link";
import { searchPostVoByPageUsingPost } from "@/api/postController";
import PostCard from "@/components/PostCard";
import "./index.css";

export default async function PostsPage({ searchParams }: { searchParams: any }) {
  const { q: searchText } = searchParams;
  let postList: any[] = [];

  try {
    const res = await searchPostVoByPageUsingPost({
      searchText,
      pageSize: 12,
      sortField: "createTime",
      sortOrder: "descend",
    });
    postList = (res as any).data?.records ?? [];
  } catch (e) {
    console.error("获取帖子列表失败", e);
  }

  return (
    <div id="postsPage">
      <div className="posts-header">
        <Title level={3} style={{ margin: 0 }}>讨论区</Title>
        <Link href="/post/new">
          <Button type="primary" icon={<PlusOutlined />}>发布帖子</Button>
        </Link>
      </div>
      <div className="post-list">
        {postList.map((post: any) => (
          <PostCard key={post.id} post={post} />
        ))}
        {postList.length === 0 && <div className="empty-tip">暂无帖子，快来发布第一条吧</div>}
      </div>
    </div>
  );
}
