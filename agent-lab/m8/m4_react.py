import anthropic
from jev_gates import route_folder
from tools import TOOLS, dispatch

client = anthropic.Anthropic()

SYSTEM = (
   "You are an instructor that teaches information from the personal knowledge vault (Cortex). "
   "Answer the user's question using ONLY what you find in the vault. "
   "Start by reading 'index.md' (the catalog), then open the specific pages it points to. "
   "Cite the file paths you drew from.  If the vault does not cover it, say so. "
)

def run_agent(question: str, max_turns: int = 12, hint: str | None = None) -> str:
    user = question if hint is None else f"{question}\n\n{hint}"    
    messages = [{"role": "user", "content": user}]
    total_in = 0
    for turn in range(max_turns):
        response = client.messages.create(
            model="claude-opus-4-8",
            max_tokens=16000,
            system=SYSTEM,
            tools=TOOLS,
            messages=messages,
        )
        total_in += response.usage.input_tokens
        print(f"  [turn {turn}] tokens: in={response.usage.input_tokens} out={response.usage.output_tokens}")
      
        if response.stop_reason == "end_turn":
            print(f"  [run] turns={turn + 1} total_input_tokens={total_in}")
            return next((b.text for b in response.content if b.type == "text"), "")

        messages.append({"role": "assistant", "content": response.content})

        tool_results = []
        for block in response.content:
            if block.type == "tool_use":
                print(f"  [turn {turn}] {block.name}({block.input})")
                output = dispatch(block.name, block.input)
                tool_results.append({
                    "type": "tool_result",
                    "tool_use_id": block.id,
                    "content": output,
                })
        messages.append({"role": "user", "content": tool_results})

    return "Stopped: hit the max-turns limit without finishing."

def routed(question: str) -> str:
    folder, conf = route_folder(question)
    print(f"  [route] {folder} @ {conf:.2f}")
    if folder == "other" or conf < 0.6:
        return run_agent(question)
    return run_agent(question, hint=(
        f"Hint: the answer most likely lives under 'codex/{folder}/'. "
        "Start with list_files on that folder and read from there. "
        "Fall back to index.md only if nothing there fits."
    ))

if __name__ == "__main__":
    answer = run_agent("What does the vault say about how good of a chance I have at finding a job?")
    print("\n=== Answer ===")
    print(answer)
