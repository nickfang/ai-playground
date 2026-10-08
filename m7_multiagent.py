from langchain_anthropic import ChatAnthropic
from langchain_core.messages import SystemMessage, HumanMessage, ToolMessage
from lc_tools import LC_TOOLS, TOOL_BY_NAME
import re

retriever = ChatAnthropic(model="claude-opus-4-8", max_tokens=16000).bind_tools(LC_TOOLS)
RETRIEVER_SYSTEM = SystemMessage(
    "Find and read the vault pages relevant to the question. "
    "Use list_files. search_text. and read_file. "
    "When you have read the relevant pages.  Reply with one line naming and the paths that mattered. "
    "Do NOT answer the question yourself - just locate and read the source material. "
)

def retrieve(question: str, max_turns: int = 8) -> dict[str, str]:
    """Run the retriever; return (vault path: verbatim page text) for every page it read."""
    excerpts: dict[str, str] = {}
    messages = [RETRIEVER_SYSTEM, HumanMessage(question)]
    for _ in range(max_turns):
        ai = retriever.invoke(messages)
        messages.append(ai)
        if not ai.tool_calls:
            return excerpts
        for call in ai.tool_calls:
            print(f"  [retriever] {call['name']}({call['args']})")
            output = TOOL_BY_NAME[call["name"]].invoke(call["args"])
            if call["name"] == "read_file":
                excerpts[call["args"]["path"]] = output
            messages.append(ToolMessage(content=output, tool_call_id=call["id"]))
    print("  [retriever] WARNING: hit max_turns - excerpts may be incomplete")
    return excerpts

def format_excerpts(excerpts: dict[str, str]) -> str:
    if not excerpts:
        return "NO PAGES WERE READ."
    return "\n\n".join(f"#### {path}\n{text}" for path, text in excerpts.items())

# --- Synthesizer. no tools: writes the cited answer from the excerpts. ---
synthesizer = ChatAnthropic(model="claude-opus-4-8", max_tokens=16000)

def synthesize(question: str, excerpts_text: str) -> str:
    return synthesizer.invoke([HumanMessage(
        "Answer the question using ONLY the vault excerpts below. "
        "Cite the file paths you draw from.  If the excerpts do not cover it, say so.\n\n"
        f"Question: { question}\n\nExcerpts:\n{excerpts_text}"
    )]).text

def critique(question: str, answer: str, excerpts_text: str) -> str:
    return synthesizer.invoke([HumanMessage(
        "Is every claim in the answer supported by the excerpts, and are file paths cited? "
        f"Reply PASS or FAIL with one sentence.\n\n"
        f"Question: {question}\n\nAnswer: {answer}\n\nExcerpts:\n{excerpts_text}"
    )]).text

def run(question: str) -> str:
    excerpts = retrieve(question)
    excerpts_text = format_excerpts(excerpts)
    answer = synthesize(question, excerpts_text)
    cited = re.findall("#### (.*.md)", excerpts_text)
    for path in cited:
        if path not in excerpts:
            print(f"  [citations] HALLUCINATED: {path}")
    print(f"  [critic] {critique(question, answer, excerpts_text)}")
    return answer

if __name__ == "__main__":
    print("\n=== Answer ===")
    print(run("What does the vault say about drug trials?"))

