"use server";
import { listQuestionBankVoByPageUsingPost } from "@/api/questionBankController";
import QuestionBankList from "@/components/QuestionBankList";
import PageContainer from "@/components/PageContainer";
import SectionHeader from "@/components/SectionHeader";
import "./index.css";

export default async function BanksPage() {
  let questionBankList = [];
  const pageSize = 100;
  try {
    const res = await listQuestionBankVoByPageUsingPost({
      pageSize,
      sortField: "createTime",
      sortOrder: "descend",
    });
    questionBankList = res.data.records ?? [];
  } catch (e) {
    console.error("获取题库列表失败", e);
  }

  return (
    <PageContainer>
      <div id="banksPage">
        <SectionHeader
          title="题库大全"
          description="系统化题库集合，帮助你按方向规划学习路线。"
        />
        <QuestionBankList questionBankList={questionBankList} />
      </div>
    </PageContainer>
  );
}
