# NaturalScript

I find it offputting when people generate scripts with LLM and commit
the results into Git repositories. Those scripts are verbose and hard to maintain.
It's like committing a binary file without the actual source code.

NaturalScript lets you write executable scripts in natural language.

```shell
#!/bin/naturalscript
Say hello world in Italian
```

On first run (or when the prompt changes), it launches an interactive coding
agent, such as Claude Code or OpenCode, to perform the task as if it's the interpreter.
You get to watch & guide this process.

When it's done, the agent will generate the executable scripRt from that session,
and it gets stored in the same file.

```text
#!/bin/naturalscript
Say hello world in Italian
-=-=-=-=-=-=-=-= GENERATED CODE BELOW: DO NOT MODIFY -=-=-=-=-=-=-=-=
H4sIAAAAAAAA/wpOrFTISM3JyVcozy/KSVHIzFPwLEnMyUzMAwQAAP//JvwLdRoAAAA=
-=-=-=-=-=-=-=-= GENERATED CODE BELOW: DO NOT MODIFY -=-=-=-=-=-=-=-=
#!/bin/bash
echo "Ciao, Mondo!"
```

On the future invocations, this generated script will run directly without
involving agents -- unless the prompt changes, in which case the agent will
be triggered again.

# Install

Download the archive for your platform from the repository's
**Releases** page, then place `naturalscript` on your `PATH`.

Alternatively,
```bash
go install github.com/kohsuke/NaturalScript@latest
```