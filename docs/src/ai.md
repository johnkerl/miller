<!---  PLEASE DO NOT EDIT DIRECTLY. EDIT THE .md.in FILE PLEASE. --->
<div>
<span class="quicklinks">
Quick links:
&nbsp;
<a class="quicklink" href="../reference-main-flag-list/index.html">Flags</a>
&nbsp;
<a class="quicklink" href="../reference-verbs/index.html">Verbs</a>
&nbsp;
<a class="quicklink" href="../reference-dsl-builtin-functions/index.html">Functions</a>
&nbsp;
<a class="quicklink" href="../glossary/index.html">Glossary</a>
&nbsp;
<a class="quicklink" href="../release-docs/index.html">Release docs</a>
</span>
</div>
# Miller and AI

As of version 6.20, released in July 2026, Miller supports two ways to let agents know about it:
an **agent skill** and **MCP**. Either one works -- not sure which? Start with the Miller agent skill.

This page covers essential setup, and an example session. For more on agent skills, see [The Miller
Agent Skill](agent-skill.md); for more on MCP, see [The Miller MCP server](mcp-server.md).

## Quick start

First, you need to **install Miller 6.20 or newer** ([Installing Miller](installing-miller.md)).
Everything on this page ships inside the ordinary `mlr` binary -- there are no plugins, no separate
installs, no API keys, and nothing here makes network calls.

Examples for Claude Code (choose one):

<pre class="pre-highlight-in-pair">
<b>mlr skill install ~/.claude/skills/miller</b>
</pre>
<pre class="pre-non-highlight-in-pair">
Wrote /Users/kerl/.claude/skills/miller/SKILL.md
</pre>

<pre class="pre-highlight-in-pair">
<b>claude mcp add miller -- mlr mcp</b>
</pre>
<pre class="pre-non-highlight-in-pair">
Added stdio MCP server miller with command: mlr mcp to local config
File modified: /Users/kerl/.claude.json [project: /Users/kerl/git/johnkerl/miller]
</pre>

## Before and after: a first session with the skill installed

If you're new to Miller, or you've used Miller before but this is your first time on 6.20 or newer,
here's a worked example: install the skill, then watch what changes about talking to your AI
assistant.

One thing to be clear on before the example: you never type `mlr` yourself in this section. You
type plain English to your agent, same as always. Every `mlr` command shown below is the agent's
*own* work -- what it runs on your behalf, in the background, to answer you. They're printed here
so you can see exactly what changes, not because you'd type them.

### Before: an agent guessing at your data

Say you're looking at [example.csv](example.csv) for the first time. You type this, and nothing
else, to your AI assistant:

> **You:** In example.csv, show me the red rows.

Without the skill, a reasonable-sounding guess for the DSL might be `$color == "Red"`. Here's that
guess, run exactly as the agent would run it, behind the scenes, on your machine:

<pre class="pre-highlight-non-pair">
<b>mlr --csv filter '$color == "Red"' example.csv</b>
</pre>

Nothing comes back -- no error, no rows, no warning, exit code 0. And here's the trap: the agent
still owes you an answer, so it turns that silence into a sentence:

> **Agent:** I checked example.csv and there aren't any rows where color is red.

That's wrong -- there are four -- but it *reads* like a fact, because a wrong guess about *your
data*, unlike a wrong flag or function name, doesn't look like a failure on Miller's end. It looks
like an empty result, which could just as easily have been true.

### After: an agent that checks first

Same question, word for word, with the skill installed:

> **You:** In example.csv, show me the red rows.

The skill's playbook puts a step between your question and any guesswork: "constrain -- learn the
data before touching it." So before writing any comparison, the agent runs a `describe`, again
invisibly, on your behalf:

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --ojson describe example.csv | mlr --json filter '$field_name == "color"'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
[
{
  "field_name": "color",
  "types": {
    "string": 10
  },
  "count": 10,
  "null_count": 0,
  "distinct_count": 3,
  "min": "purple",
  "max": "yellow",
  "values": ["yellow", "red", "purple"]
}
]
</pre>

That's the real value set for `color`: lowercase, three values, no `"Red"` anywhere. The agent
corrects its own comparison before running anything user-facing:

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --opprint filter '$color == "red"' example.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color shape  flag  k index quantity rate
red   square true  2 15    79.2778  0.0130
red   circle true  3 16    13.8103  2.9010
red   square false 4 48    77.5542  7.4670
red   square false 6 64    77.1991  9.5310
</pre>

-- and only now replies to you:

> **Agent:** Four rows have color = red: rows 2, 3, 4, and 6 (all square or circle, quantity
> ranging from about 13.8 to 79.3).

Same question, same data, same underlying `mlr` binary -- the only thing that changed is that the
agent looked before it leapt, and you never saw the intermediate `describe` unless you asked to.
That one habit, *check the data before writing a comparison*, is the skill in miniature; the rest
of the playbook applies the same idea to verb and function names (discover), DSL syntax (validate),
and error messages (recover) -- see [Miller AI internals](ai-support.md) for how each of those works.

## Why AI support

Miller treats AI agents as first-class users. When an agent drives a command-line tool, the agent
can fail in predictable ways: it invents flags that don't exist, guesses values that aren't in the
data, misreads error prose, and burns whole runs discovering a typo. Miller closes off each of those
failure modes with the following structure:

* Miller's entire surface -- verbs, DSL functions, flags, keywords -- is
  available as **machine-readable JSON**, so agents ground themselves in what
  actually exists.
* Options with fixed domains carry their **complete value sets**, and input
  data can be **profiled in one pass**, so that agents copy real values instead
  of inventing them.
* DSL expressions can be **validated before running**, without reading any
  input.
* **Errors are structured** -- kind, hint, did-you-mean -- so agents branch
  on data rather than parsing English.
* A **sandbox flag** removes external-command execution, so an
  agent-constructed command line is just data processing.

Every one of those is an ordinary command-line feature, documented in [Miller AI
internals](ai-support.md): each works from any agent harness, system prompt, or script.

## Skill file or MCP: which should you use?

For day one, the short version: start with the skill; add MCP later if you want it. They aren't
exclusive; nothing stops you running both.

**Miller agent skill file:**

- Plus: One command, one static file -- no process, no client registration, nothing to reconnect.
- Plus: Works with any agent that reads Agent Skills from disk, not just MCP clients.
- Minus: No enforcement: it's advisory text, so no automatic `--no-shell` sandbox, no output caps or timeouts.
- Minus: The agent parses plain `mlr` text output and exit codes itself -- no structured JSON per call.

**Miller MCP server:**

- Plus: Structured typed calls in, structured JSON back -- no text-parsing on the agent's side.
- Plus: Sandboxed by default (`MLR_NO_SHELL=1`), output-capped, timeout-guarded.
- Minus: One more moving part: per-client registration, plus a subprocess to spawn and reconnect each session.
- Minus: Only helps agents that actually speak MCP.

In one line: the skill is less setup and the most portable, with weaker guarantees; MCP is a bit
more setup, with stronger guarantees, for a narrower set of clients.
