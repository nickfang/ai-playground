from langchain_core.tools import tool
import tools

@tool
def read_file(path: str) -> str:
   """Read a UTF-8 text file from the Cortex vault.  Path is vault-relative, e.g. 'codex/ai/Prompt Engineering.md'."""
   return tools.read_file(path)

@tool
def list_files(subdir: str = "") -> str:
   """List files and folders under a vault-relative directory.  Pass '' for the vault root."""
   return tools.list_files(subdir)

@tool
def search_text(query:str) -> str:
   """Case-insensitive substring search across all .md files in the vault.  Returns 'path:line: text' hits."""
   return tools.search_text(query)

LC_TOOLS=[read_file, list_files, search_text]
TOOL_BY_NAME= {t.name: t for t in LC_TOOLS}
