FROM docker.elastic.co/elasticsearch/elasticsearch:7.17.4

# 安装 ik 中文分词器
RUN bin/elasticsearch-plugin install --batch https://release.infinilabs.com/analysis-ik/stable/elasticsearch-analysis-ik-7.17.4.zip
