import { Button } from "antd";
import { PlusOutlined } from "@ant-design/icons";
import Link from "next/link";
import { searchPostVoByPageUsingPost } from "@/api/postController";
import PostCard from "@/components/PostCard";
import PageContainer from "@/components/PageContainer";
import SectionHeader from "@/components/SectionHeader";
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
