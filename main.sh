#!/bin/bash

# 環境変数の読み込み
if [ -f .env ]; then
  export $(cat .env | xargs)
fi

# Git の差分を取得
DIFF=$(git diff --cached)
PROMPT="以下の Git の変更差分を見て、適切なコミットメッセージを提案してください:\n\n$DIFF"

COMMIT_MESSAGE=$(curl -s https://api.openai.com/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H `Authorization: Bearer ${OPENAI_API_KEY}` \
  -d '{
    "model": "gpt-4",
    "messages": [{"role": "system", "content": "あなたは優れたソフトウェアエンジニアです。"}, {"role": "user", "content": "'"${PROMPT}"'"}],
    "temperature": 0.7
  }' | jq -r '.choices[0].message.content')

# 結果を表示
echo "$COMMIT_MESSAGE"
