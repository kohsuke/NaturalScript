# NaturalScript

LLMでスクリプトを生成して、その結果をGitリポジトリにそのままコミットするのは、
私はあまり好ましくないと感じています。そうしたスクリプトは冗長で保守しづらく、
実質的にソースコードなしでバイナリをコミットするのに近いからです。

NaturalScriptを使うと、自然言語で実行可能なスクリプトを**書いて、さらに保守**できます。
シェル/Python/... のスクリプトを、直接編集しない生成物として扱えるようになります。

まず、スクリプトにさせたいことを説明したテキストファイルを作成します:
```text
$ cat check-google
google.com が稼働中か停止中かを確認する。
停止していたら非ゼロで終了する
```

次に、それを `naturalscript` に渡します:
```
$ naturalscript hello
```

これにより、Claude Code や OpenCode のような対話型コーディングエージェントが起動します。
このエージェントが、このスクリプトを「実行」するためのインタープリタとして振る舞います。
必要に応じて、要件を明確にするための対話も行います。

処理が完了すると、エージェントはそのセッションから実行可能なスクリプトを生成します。
セッションを終了すると、NaturalScript はすべてをまとめて同じファイルに次のように書き戻します:

```text
#!/bin/bash
: <<'COMMENTBLOCK_FOR_NATURALSCRIPT'
Managed by NaturalScript. Edit the prompt below and run `naturalscript path/to/this/file` to regenerate.
==== NATURALSCRIPT:BEGIN ====
Test if google.com is up or down.
If it's down, exit non-zero

==== NATURALSCRIPT:END ====
H4sIAAAAAAAA/wpJLS5RyExTSM/PT89J1UvOz1XILFYoLVDIL1JIyS/P0+PyTFPILFEvBvN0FFIrMk
sU8vLzdKtSi/K5AAEAAP//uCSZCD4AAAA=

COMMENTBLOCK_FOR_NATURALSCRIPT

... (generated script follows) ...
```

これで、この生成済みスクリプトは NaturalScript なしで実行できます。

スクリプトを更新したい場合は、このファイル内のマークされたプロンプト部分を編集して、
`naturalscript` を再実行するだけです。エージェントは旧プロンプトと新プロンプトの差分を、
現在のスクリプトに反映しようとします。



# Install

[お使いのプラットフォーム向けアーカイブをダウンロード](https://github.com/kohsuke/NaturalScript/releases)し、
`naturalscript` を `PATH` に配置してください。

あるいは:

```bash
go install github.com/kohsuke/NaturalScript@latest
```
