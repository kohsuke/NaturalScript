# NaturalScript

I find it off-putting when people generate scripts with LLMs and commit
the results to Git repositories. Those scripts are verbose and hard to maintain.
It's like committing a binary file without the actual source code.

NaturalScript lets you write executable scripts in natural language.

```shell
#!/bin/naturalscript
Say hello world in Italian
```

On first run (or when the prompt changes), it launches an interactive coding
agent, such as Claude Code or OpenCode, to perform the task as your interpreter.
You get to watch and guide this process.

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

On future invocations, this generated script runs directly without
involving agents -- unless the prompt changes, in which case the agent is
triggered again.

# Install

Download the archive for your platform from the repository's
**Releases** page, then place `naturalscript` on your `PATH`.

Alternatively:

```bash
go install github.com/kohsuke/NaturalScript@latest
```
