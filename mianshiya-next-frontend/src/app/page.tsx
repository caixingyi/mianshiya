import Title from "antd/es/typography/Title";
import { Button, Divider, Space } from "antd";
import { BookOutlined, RobotOutlined, SearchOutlined } from "@ant-design/icons";
import Link from "next/link";
import { listQuestionBankVoByPageUsingPost } from "@/api/questionBankController";
import { listQuestionVoByPageUsingPost } from "@/api/questionController";
import QuestionBankList from "@/components/QuestionBankList";
import QuestionList from "@/components/QuestionList";
import PageContainer from "@/components/PageContainer";
import SectionHeader from "@/components/SectionHeader";
import "./index.css";

export const dynamic = "force-dynamic";

export default async function HomePage() {
  let questionBankList = [];
  let questionList = [];
  try {
    const res = await listQuestionBankVoByPageUsingPost({
      pageSize: 12,
      sortField: "createTime",
      sortOrder: "descend",
    });
    questionBankList = res.data.records ?? [];
  } catch (e) {
    console.error(e);
  }

  try {
    const res = await listQuestionVoByPageUsingPost({
      pageSize: 12,
      sortField: "createTime",
      sortOrder: "descend",
    });
    questionList = res.data.records ?? [];
  } catch (e) {
    console.error(e);
  }

  return (
    <PageContainer>
      <div id="homePage">
        <section className="hero-banner">
          <div className="hero-badge">AI 驱动的面试刷题平台</div>
          <Title level={1} className="hero-title">
            高效准备技术面试
          </Title>
          <p className="hero-subtitle">
            题库刷题、AI 模拟面试、讨论交流，一站式提升面试通过率。
          </p>
          <Space size={12} className="hero-actions">
            <Link href="/questions">
              <Button type="primary" size="large" icon={<SearchOutlined />}>开始刷题</Button>
            </Link>
            <Link href="/mockInterview/add">
              <Button size="large" icon={<RobotOutlined />}>AI 模拟面试</Button>
            </Link>
          </Space>
        </section>

        <section className="home-section">
          <SectionHeader
            title="最新题库"
            description="按专题系统学习，快速定位薄弱知识点。"
            extra={<Link href="/banks">查看更多 →</Link>}
          />
          <QuestionBankList questionBankList={questionBankList} />
        </section>

        <Divider className="home-divider" />

        <section className="home-section">
          <SectionHeader
            title="最新题目"
            description="精选高频面试题，覆盖后端、前端、数据库、系统设计等方向。"
            extra={<Link href="/questions">查看更多 →</Link>}
          />
          <QuestionList questionList={questionList} />
        </section>
      </div>
    </PageContainer>
  );
}
