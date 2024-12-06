# ai commits message
![demo](https://raw.githubusercontent.com/wiki/motty93/ai-commits-message/images/ai-commit-message-demo.gif)

## (Deprecated) shell script
```bash
sh script/ai.sh
```

## (Recommended) go run
```bash
OPENAI_API_KEY="sk-xxxxxxx" LANGUAGE="JPN" go run -ldflags "-X 'main.apiKey=${OPENAI_API_KEY}' -X 'i18n.lang=${LANGUAGE}'" main.go
```
### build
```bash
OPENAI_API_KEY="sk-xxxxxxx" LANGUAGE="JPN" go build -o ./bin/main -ldflags "-X 'main.apiKey=${OPENAI_API_KEY}' -X 'i18n.lang=${LANGUAGE}'" main.go

./bin/main
```

### debug
```bash
air
```

## vim setting
### build
```bash
OPENAI_API_KEY="sk-xxxxxxx" LANGUAGE="JPN" go build -o ./bin/main -ldflags "-X 'main.apiKey=${OPENAI_API_KEY}' -X 'i18n.lang=${LANGUAGE}'" main.go

cp ./bin/main ~/.config/generate_commit_message
```

### .vimrc
```vim
command! -nargs=0 AICommitMessage call AICommitMessage()
function! AICommitMessage()
  " コマンドの出力を取得
  let l:message = system("~/.config/generate_commit_message 2> /dev/null")

  " messageが空であれば何もしない
  if l:message == ''
    return
  endif

  " 出力のエラーハンドリング
  if v:shell_error != 0
    echohl ErrorMsg
    echo "Error running generate_commit_message"
    echohl None
    return
  endif

  " 出力結果の改行をtrim
  let l:message = substitute(l:message, '\n\+$', '', '')

  " カーソル位置に挿入（改行しない）
  call setline('.', getline('.') . l:message)
endfunction
```

## vim command
normal mode
```vim
:AICommitMessage
```
