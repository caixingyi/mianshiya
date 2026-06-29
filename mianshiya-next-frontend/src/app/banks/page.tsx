import QuestionBankList from "@/components/QuestionBankList";
import PageContainer from "@/components/PageContainer";
import SectionHeader from "@/components/SectionHeader";
import "./index.css";

export default async function BanksPage() {
  let questionBankList: any[] = [];
  const apiBase = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8101";

  try {
    const res = await fetch(`${apiBase}/api/questionBank/list/page/vo`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        pageSize: 100,
        sortField: "createTime",
        sortOrder: "descend",
      }),
      cache: "no-store",
    });
    const json = await res.json();
    questionBankList = json.data?.records ?? [];
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
