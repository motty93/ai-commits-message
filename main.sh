#!/bin/bash

SCIRPT_DIR=$(cd $(dirname $0); pwd)

# 環境変数の読み込み
if [ -f "$SCIRPT_DIR/.env" ]; then
  set -a # 自動でexport
  . "$SCIRPT_DIR/.env" # 環境変数を読み込む
  set +a # 自動exportを無効化
fi

if ! git rev-parse --is-inside-work-tree > /dev/null 2>&1; then
  echo "Not a git repository"
  exit 1
fi

DIFF=$(git diff --cached)
if [ -z "$DIFF" ]; then
  echo "No changes to commit"
  exit 0
fi

CONTENT="あなたは優れたソフトウェアエンジニアです。"
PROMPT="以下の Git の変更差分を見て、適切なコミットメッセージを提案してください:\n\n$DIFF"
ESCAPE_CONTENT=$(printf '%s' "$CONTENT" | sed 's/"/\\"/g')
ESCAPE_PROMPT=$(printf '%s' "$PROMPT" | sed 's/"/\\"/g' | sed ':a;N;$!ba;s/\n/\\n/g')
JSON_PAYLOAD=$(printf '{
  "model": "gpt-4",
  "messages": [
    {"role": "system", "content": "%s"},
    {"role": "user", "content": "%s"}
  ],
  "temperature": 0.7
}' "$ESCAPE_CONTENT" "$ESCAPE_PROMPT")

COMMIT_MESSAGE=$(curl -s https://api.openai.com/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ${OPENAI_API_KEY}" \
  -d "$JSON_PAYLOAD" | jq -r '.choices[0].message.content')

echo "Generated Commit Message:"
echo "$COMMIT_MESSAGE"
