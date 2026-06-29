import { Button } from "antd";
import { PlusOutlined } from "@ant-design/icons";
import Link from "next/link";
import PostCard from "@/components/PostCard";
import PageContainer from "@/components/PageContainer";
import SectionHeader from "@/components/SectionHeader";
import "./index.css";

export default async function PostsPage({ searchParams }: { searchParams: any }) {
  const { q: searchText } = searchParams;
  let postList: any[] = [];
  const apiBase = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8101";

  try {
    const res = await fetch(`${apiBase}/api/post/search/page/vo`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        searchText,
        pageSize: 12,
        sortField: "createTime",
        sortOrder: "descend",
      }),
      cache: "no-store",
    });
    const json = await res.json();
    postList = json.data?.records ?? [];
  } catch (e) {
    console.error("获取帖子列表失败", e);
  }

  return (
    <PageContainer narrow>
      <div id="postsPage">
        <SectionHeader
          title="讨论区"
          description="技术交流与经验分享"
          extra={
            <Link href="/post/new">
              <Button type="primary" icon={<PlusOutlined />}>发布帖子</Button>
            </Link>
          }
        />
        <div className="post-list">
          {postList.map((post: any) => (
            <PostCard key={post.id} post={post} />
          ))}
          {postList.length === 0 && <div className="empty-tip">暂无帖子，快来发布第一条吧</div>}
        </div>
      </div>
    </PageContainer>
  );
}
