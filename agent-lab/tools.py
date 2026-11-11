import os
from pathlib import Path

VAULT_ROOT = Path(os.environ.get("CORTEX_VAULT", Path.home() / "Cortex")).resolve()

def _safe_path(rel: str) -> Path:
   """Resolve a vault-relative path and refuse anything that escapes the vault."""
   p = (VAULT_ROOT / rel).resolve()
   if not p.is_relative_to(VAULT_ROOT):
      raise ValueError(f"Path escapes the vault: {rel!r}")
   return p

def read_file(path:str) -> str:
   p = _safe_path(path)
   if not p.is_file():
      return f"Error: no file at {path!r}"
   return p.read_text(encoding="utf-8")

READ_FILE_TOOL = {
   "name": "read_file",
   "description": (
      "Read a UTF-8 text file from the Cortex vault. "
      "Path is relative to the vault root, e.g. 'codex/ai/Prompt Engineering.md'."
   ),
   "input_schema": {
      "type": "object",
      "properties": {
         "path": {"type": "string", "description": "Vault-relative file path"},
      },
      "required": ["path"],
   },
}

def list_files(subdir: str = "") -> str:
   p = _safe_path(subdir)
   if not p.is_dir():
      return f"Error: no directory at {subdir!r}"
   names = []
   for x in sorted(p.iterdir()):
      rel = x.relative_to(VAULT_ROOT).as_posix()
      names.append(rel + "/" if x.is_dir() else rel)
   return "\n".join(names)

LIST_FILES_TOOL = {
   "name": "list_files",
   "description": "List files and folders under a vault-relative directory(default: vault root).",
   "input_schema": {
      "type": "object",
      "properties": {"subdir": {"type": "string", "description": "Vault-relative directory, or '' for root"}},
      "required": [],
   },
}

def search_text(query: str) -> str:
   """Case-insensitive substring search over all.md files. Returns up to 50 hits."""
   hits = []
   for f in VAULT_ROOT.rglob("*.md"):
      try:
         for i, line in enumerate(f.read_text(encoding="utf-8").splitlines(), 1):
            if query.lower() in line.lower():
               rel = f.relative_to(VAULT_ROOT).as_posix()
               hits.append(f"{rel}:{i}: {line.strip()}")
               if len(hits) >=50:
                  return "\n".join(hits)
      except (UnicodeDecodeError, OSError):
         continue
   return "\n".join(hits) if hits else "No matches."

SEARCH_TEXT_TOOL = {
   "name": "search_text",
   "description": "Case-insensitive substring search across all .md files.  Returns 'path:line: text' hits.",
   "input_schema": {
      "type": "object",
      "properties": {"query": { "type": "string", "description": "Text to search for"}},
      "required": ["query"],
   },
}

TOOLS = [READ_FILE_TOOL, LIST_FILES_TOOL, SEARCH_TEXT_TOOL]

def dispatch(name: str, tool_input: dict) -> str:
   """Route a tool_use to the matching function."""
   if name == "list_files":
      return list_files(**tool_input)
   if name == "read_file":
      return read_file(**tool_input)
   if name == "search_text":
      return search_text(**tool_input)
   return f"Error: unknown tool {name!r}"

