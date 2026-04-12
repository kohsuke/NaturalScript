# NaturalScript

NaturalScript lets you write executable scripts in natural language.

```text
#!/bin/naturalscript
Say hello world in Italian
```

On first run (or when the prompt changes), it launches an interactive coding
agent, such as Claude Code or OpenCode, to perform the task as if it's the interpreter.
You get to watch & guide this process.

When it's done, the agent will generate the executable script from that session,
and it gets stored in the same file.

```text
#!/bin/naturalscript
Say hello world in Italian
-=-=-=-=-=-=-=-=
H4sIAAAAAAAA/wpOrFTISM3JyVcozy/KSVHIzFPwLEnMyUzMAwQAAP//JvwLdRoAAAA=
-=-=-=-=-=-=-=-=
#!/bin/bash
echo "Ciao, Mondo!"
```

On the future invocations, this generated script will run directly without
involving agents.

If the first prompt changes, the agent will be triggered again to update the script.

# Install
