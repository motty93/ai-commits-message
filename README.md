# ai commits message
## (Deprecated) shell script
```bash
sh script/ai.sh
```

## (Recommended) go run
```bash
go run main.go
```
### build
```bash
go build -o ./bin/main main.go

./bin/main
```

### debug
```bash
air
```

## vim setting
```vim
command! -nargs=0 AICommitMessage call AICommitMessage()
function! AICommitMessage()
  " コマンドの出力を取得
  let l:message = system("~/.config/generate_commit_message 2> /dev/null")

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
### vim command
normal mode
```vim
:AICommitMessage
```
