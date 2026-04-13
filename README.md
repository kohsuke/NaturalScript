# NaturalScript

I find it off-putting when people generate scripts with LLMs and commit
the results to Git repositories. Those scripts are verbose and hard to maintain.
It's like committing a binary file without the actual source code.

NaturalScript lets you write executable scripts in natural language.
Create a text file and pass it to `naturalscript`:

```text
$ cat check-google
Test if google.com is up or down.
If it's down, exit non-zero
$ naturalscript hello
```

It launches an interactive coding agent, such as Claude Code or OpenCode, to execute this "script".
It'll engage you in a conversation to clarify the requirements.

When it's done, the agent will generate the executable script from that session,
and it gets written back into the same file.

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

To update the script, simply edit the captured prompt section,
and re-run `naturalscript`. The agent will now try to apply the delta
between the old and new prompts into the current script.



# Install

Download the archive for your platform from the repository's
**Releases** page, then place `naturalscript` on your `PATH`.

Alternatively:

```bash
go install github.com/kohsuke/NaturalScript@latest
```
