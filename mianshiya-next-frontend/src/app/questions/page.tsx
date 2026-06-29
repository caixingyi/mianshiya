"use server";
import { searchQuestionVoByPageUsingPost } from "@/api/questionController";
import QuestionTable from "@/components/QuestionTable";
import PageContainer from "@/components/PageContainer";
import SectionHeader from "@/components/SectionHeader";
import "./index.css";

export default async function QuestionsPage({ searchParams }: { searchParams: any }) {
  const { q: searchText } = searchParams;
  let questionList = [];
  let total = 0;

  try {
    const res = await searchQuestionVoByPageUsingPost({
      searchText,
      pageSize: 12,
      sortField: "createTime",
      sortOrder: "descend",
    });
    questionList = (res as any).data?.records ?? [];
    total = (res as any).data?.total ?? 0;
  } catch (e) {
    console.error("获取题目列表失败", e);
  }

  return (
    <PageContainer>
      <div id="questionsPage">
        <SectionHeader
          title="题目大全"
          description="高频面试题库，支持关键词搜索和标签筛选。"
        />
        <QuestionTable
          defaultQuestionList={questionList}
          defaultTotal={total}
          defaultSearchParams={{
            searchText,
          }}
        />
      </div>
    </PageContainer>
  );
}
