"use client";
import { Avatar, Card, List, Typography } from "antd";
import Link from "next/link";
import "./index.css";

interface Props {
  questionBankList: API.QuestionBankVO[];
}

const QuestionBankList = (props: Props) => {
  const { questionBankList = [] } = props;

  const questionBankView = (questionBank: API.QuestionBankVO) => {
    return (
      <Card hoverable className="bank-card">
        <Link href={`/bank/${questionBank.id}`}>
          <Card.Meta
            avatar={
              <Avatar
                shape="square"
                size={48}
                src={questionBank.picture}
                style={{
                  background: "linear-gradient(135deg, #1677ff 0%, #69b1ff 100%)",
                }}
              >
                {questionBank.title?.[0] || "题"}
              </Avatar>
            }
            title={<span className="bank-card-title">{questionBank.title}</span>}
            description={
              <Typography.Paragraph
                type="secondary"
                ellipsis={{ rows: 2 }}
                style={{ marginBottom: 0, fontSize: 13 }}
              >
                {questionBank.description || "暂无描述"}
              </Typography.Paragraph>
            }
          />
        </Link>
      </Card>
    );
  };

  return (
    <div className="question-bank-list">
      <List
        grid={{
          gutter: 20,
          column: 4,
          xs: 1,
          sm: 2,
          md: 3,
          lg: 4,
        }}
        dataSource={questionBankList}
        renderItem={(item) => <List.Item>{questionBankView(item)}</List.Item>}
      />
    </div>
  );
};

export default QuestionBankList;
