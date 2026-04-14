# NaturalScript

I find it off-putting when people generate scripts with LLMs and commit
the results to Git repositories. Those scripts are verbose and hard to maintain.
It's like committing a binary file without the actual source code.

NaturalScript lets you write _and maintain_ executable scripts in natural language.
It lets you treat shell/Python/... scripts as generated artifacts that you don't edit directly.

First, create a text file describing what you want the script to do:
```text
$ cat check-google
Test if google.com is up or down.
If it's down, exit non-zero
```

Then pass it to `naturalscript`:
```
$ naturalscript hello
```

This launches an interactive coding agent, such as Claude Code or OpenCode.
It'll act as an interpreter to "execute" this script.
It'll engage you in a conversation to clarify the requirements, as needed.

When it's done, the agent will generate the executable script from that session.
When you exit the session, NaturalScript packs everything and writes it back to the same file like this:

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

Now this generated script can be run without NaturalScript.

To update the script, simply edit the marked prompt section in this file,
and re-run `naturalscript`. The agent will now try to apply the delta
between the old and new prompts into the current script.



# Install

[Download the archive for your platform](https://github.com/kohsuke/NaturalScript/releases),
then place `naturalscript` on your `PATH`.

Alternatively:

```bash
go install github.com/kohsuke/NaturalScript@latest
```
